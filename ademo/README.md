# A demo

A demo of ggt capabilities to create a service to create/read/update/delete `tomatoes`.


# TOC
- [The main](#the-main)
  - [$ main.go](#-maingo)
- [The model](#the-model)
  - [$ tomate/model.go](#-tomatemodelgo)
- [The controller](#the-controller)
  - [$ tomate/controller.go](#-tomatecontrollergo)
- [The gen](#the-gen)
  - [$ tomate/gen.go](#-tomategengo)
- [The code for free](#the-code-for-free)
  - [a backend in-memory](#a-backend-in-memory)
  - [an http rpc implementation](#an-http-rpc-implementation)
  - [an http rest implementation](#an-http-rest-implementation)
- [The test](#the-test)
  - [$ sh test.sh](#-sh-testsh)

# The main

#### $ main.go
```go
// A demo of ggt capabilities to create a service to create/read/update/delete `tomatoes`.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mh-cbon/ggt/ademo/tomate"
	"github.com/mh-cbon/ggt/lib"
)

func main() {

	// create the routes handler
	router := mux.NewRouter()

	// create a storage backend, in memory for current example.
	backend := tomate.NewTomatesSync()

	// populate the backend for testing
	backend.Transact(func(b *tomate.Tomates) {
		b.Push(&tomate.Tomate{ID: fmt.Sprintf("%v", b.Len()), Color: ""})
	})

	// for the fun, demonstrates generator capabilities :D
	backend.
		Filter(tomate.FilterTomates.ByID("0")).
		Map(tomate.SetterTomates.SetColor("Red"))

	// make a complete controller (transport+business+storage)
	controller := tomate.NewRestController(
		tomate.NewController(
			backend,
		),
	)
	controller.Log = &lib.WriteLog{Sink: os.Stderr}

	// create a descriptor of the controller exposed methods
	desc := tomate.NewRestControllerDescriptor(controller)

	// manipulates the handlers to wrap them
	// desc.Create().WrapMethod(logReq)
	desc.WrapMethod(logReq)

	// bind the route handlers to the routes handler
	lib.Gorilla(desc, router)

	// beer time!
	http.ListenAndServe(":8080", router)
}

func logReq(m *lib.MethodDescriptor) lib.Wrapper {
	return func(in http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("handling ", m.Name, m.Route)
			in(w, r)
		}
	}
}
```

# The model

#### $ tomate/model.go
```go
package tomate

// Tomate is a model of tomatoes
type Tomate struct {
	ID    string
	Color string
}

// GetID is useful for identity check.
func (t Tomate) GetID() string {
	return t.ID
}
```

# The controller

#### $ tomate/controller.go
```go
package tomate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Controller of tomatoes.
type Controller struct {
	backend TomatesContract
}

// NewController creates a new tomates controller
func NewController(backend TomatesContract) Controller {
	return Controller{backend}
}

// GetByID read the Tomate of given ID
//
// @route /read/{id:[0-9]+}
func (t Controller) GetByID(routeID string) (jsonResBody *Tomate, err error) {
	t.backend.Transact(func(backend *Tomates) {
		jsonResBody = backend.
			Filter(FilterTomates.ByID(routeID)).
			First()
	})
	if jsonResBody == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, err
}

// Create a new Tomate
//
// @route /create
func (t Controller) Create(postColor *string) (jsonResBody *Tomate, err error) {
	if postColor == nil {
		return nil, &UserInputError{errors.New("Missing color parameter")}
	}
	color := strings.TrimSpace(*postColor)
	if color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	t.backend.Transact(func(backend *Tomates) {
		if !backend.Filter(FilterTomates.ByColor(color)).Empty() {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		jsonResBody = &Tomate{ID: fmt.Sprintf("%v", backend.Len()), Color: color}
		backend.Push(jsonResBody)
	})
	return jsonResBody, err
}

// Update an existing Tomate
//
// @route /write/{id:[0-9]+}
func (t Controller) Update(routeID string, jsonReqBody *Tomate) (jsonResBody *Tomate, err error) {
	jsonReqBody.Color = strings.TrimSpace(jsonReqBody.Color)
	if jsonReqBody.Color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	t.backend.Transact(func(backend *Tomates) {
		byID := backend.Filter(FilterTomates.ByID(routeID))
		if byID.Len() == 0 {
			err = &NotFoundError{errors.New("ID does not exists")}
			return
		}
		byColor := backend.Filter(FilterTomates.ByColor(jsonReqBody.Color))
		if byColor.Len() > 0 && byID.First().ID != byColor.First().ID {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		jsonResBody = backend.
			Filter(FilterTomates.ByID(routeID)).
			Map(SetterTomates.SetColor(jsonReqBody.Color)).
			First()
	})
	if jsonResBody == nil && err == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, err
}

// Remove an existing Tomate
//
// @route /remove/{id:[0-9]+}
func (t Controller) Remove(ctx context.Context, routeID string) (jsonResBody bool, err error) {
	t.backend.Transact(func(backend *Tomates) {
		byID := backend.Filter(FilterTomates.ByID(routeID))
		if byID.Empty() {
			err = &NotFoundError{errors.New("ID does not exists")}
			return
		}
		jsonResBody = backend.Remove(byID.First())
	})
	return jsonResBody, err
}

// NotFoundError is an error thrown when a value is not found
type NotFoundError struct {
	error
}

// UserInputError is an error thrown when the user input is incomplete or invalid
type UserInputError struct {
	error
}

// Finalizer behave appropriately by error types
func (t Controller) Finalizer(w http.ResponseWriter, r *http.Request, err error) {
	if _, ok := err.(*UserInputError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if _, ok := err.(*NotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
```

# The gen

#### $ tomate/gen.go
```go
package tomate

//go:generate ggt -c slicer *Tomate:Tomates
//go:generate ggt chaner Tomates:TomatesSync

//go:generate ggt -mode route http-provider Controller:RestController
//go:generate ggt -mode route http-consumer Controller:RestClient

//go:generate ggt -mode rpc http-provider Controller:RPCController
//go:generate ggt -mode rpc http-consumer Controller:RPCClient
```

# The code for free

## a backend in-memory

- [tomate/zz_tomatessync.go](tomate/zz_tomatessync.go)
- [tomate/zz_tomates.go](tomate/zz_tomates.go)

## an http rpc implementation

- [tomate/zz_rpccontroller.go](tomate/zz_rpccontroller.go)
- [tomate/zz_rpcclient.go](tomate/zz_rpcclient.go)

## an http rest implementation

- [tomate/zz_restclient.go](tomate/zz_restclient.go)
- [tomate/zz_restcontroller.go](tomate/zz_restcontroller.go)

# The test

#### $ sh test.sh
```sh
+ go generate tomate/gen.go
+ CURL='curl -s -D -'
+ sleep 1
+ go run main.go
+ curl -s -D - http://localhost:8080/read/0
2017/05/26 17:06:48 handling  GetByID /read/{id:[0-9]+}
2017-05-26 17:06:48.667478223 +0200 CEST [begin RestController GetByID] <nil>
2017-05-26 17:06:48.667545705 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 25

{"ID":"0","Color":"Red"}
+ curl -s -D - http://localhost:8080/read/10
2017/05/26 17:06:48 handling  GetByID /read/{id:[0-9]+}
2017-05-26 17:06:48.677345385 +0200 CEST [begin RestController GetByID] <nil>
2017-05-26 17:06:48.67741992 +0200 CEST [business error RestController GetByID] Tomate not found
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 17

Tomate not found
+ curl -s -D - --data color=blue http://localhost:8080/create
2017/05/26 17:06:48 handling  Create /create
2017-05-26 17:06:48.687722718 +0200 CEST [begin RestController Create] <nil>
2017-05-26 17:06:48.68781314 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 26

{"ID":"1","Color":"blue"}
+ curl -s -D - --data color=blue http://localhost:8080/create
2017/05/26 17:06:48 handling  Create /create
2017-05-26 17:06:48.697714722 +0200 CEST [begin RestController Create] <nil>
2017-05-26 17:06:48.697794639 +0200 CEST [business error RestController Create] color must be unique
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 21

color must be unique
+ curl -s -D - --data color= http://localhost:8080/create
2017/05/26 17:06:48 handling  Create /create
2017-05-26 17:06:48.708094987 +0200 CEST [begin RestController Create] <nil>
2017-05-26 17:06:48.708167699 +0200 CEST [business error RestController Create] color must not be empty
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 24

color must not be empty
+ curl -s -D - --data color=green http://localhost:8080/create
2017/05/26 17:06:48 handling  Create /create
2017-05-26 17:06:48.718767944 +0200 CEST [begin RestController Create] <nil>
2017-05-26 17:06:48.718885813 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 27

{"ID":"2","Color":"green"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/26 17:06:48 handling  Update /write/{id:[0-9]+}
2017-05-26 17:06:48.72879706 +0200 CEST [begin RestController Update] <nil>
2017-05-26 17:06:48.728898323 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - http://localhost:8080/read/1
2017/05/26 17:06:48 handling  GetByID /read/{id:[0-9]+}
2017-05-26 17:06:48.74520684 +0200 CEST [begin RestController GetByID] <nil>
2017-05-26 17:06:48.745339264 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/0
2017/05/26 17:06:48 handling  Update /write/{id:[0-9]+}
2017-05-26 17:06:48.760152446 +0200 CEST [begin RestController Update] <nil>
2017-05-26 17:06:48.760229136 +0200 CEST [business error RestController Update] color must be unique
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 21

color must be unique
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/26 17:06:48 handling  Update /write/{id:[0-9]+}
2017-05-26 17:06:48.776641419 +0200 CEST [begin RestController Update] <nil>
2017-05-26 17:06:48.776729174 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - http://localhost:8080/read/0
2017/05/26 17:06:48 handling  GetByID /read/{id:[0-9]+}
2017-05-26 17:06:48.792532726 +0200 CEST [begin RestController GetByID] <nil>
2017-05-26 17:06:48.792591968 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 26 May 2017 15:06:48 GMT
Content-Length: 25

{"ID":"0","Color":"Red"}
+ killall main
signal: terminated
```
