# A demo

[![GoDoc](https://godoc.org/github.com/mh-cbon/ggt?status.svg)](http://godoc.org/github.com/mh-cbon/ggt)

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

// SimilarTomate indiicates tomate similarity to a value
type SimilarTomate struct {
	Tomate
	Similarity float64
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

	"github.com/agext/levenshtein"
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
	byRouteID := FilterTomates.ByID(routeID)
	t.backend.Transact(func(backend *Tomates) {
		jsonResBody = backend.Filter(byRouteID).First()
	})
	if jsonResBody == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, err
}

// SimilarColor returns colors similar to the given input
//
// @route /similar/color/{color}
func (t Controller) SimilarColor(routeColor string, getSensitive *bool) (jsonResBody *SimilarTomates, err error) {
	sensitive := getSensitive != nil && *getSensitive != false
	rVal := routeColor
	p := levenshtein.NewParams()
	if sensitive == false {
		rVal = strings.ToLower(rVal)
	}
	jsonResBody = NewSimilarTomates()
	t.backend.Map(func(t *Tomate) *Tomate {
		lVal := t.Color
		if sensitive == false {
			lVal = strings.ToLower(lVal)
		}
		res := levenshtein.Similarity(lVal, rVal, p)
		if res > 0.1 {
			jsonResBody.Push(&SimilarTomate{*t, res})
		}
		return t
	})
	return jsonResBody, err
}

// Create a new Tomate
//
// @route /create
// @methods POST
func (t Controller) Create(postColor *string) (jsonResBody *Tomate, err error) {
	if postColor == nil {
		return nil, &UserInputError{errors.New("Missing color parameter")}
	}
	color := strings.TrimSpace(*postColor)
	if color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	byPostColor := FilterTomates.ByColor(color)
	t.backend.Transact(func(backend *Tomates) {
		if backend.Filter(byPostColor).NotEmpty() {
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
// @methods POST
func (t Controller) Update(routeID string, jsonReqBody *Tomate) (jsonResBody *Tomate, err error) {
	jsonReqBody.Color = strings.TrimSpace(jsonReqBody.Color)
	if jsonReqBody.Color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	byRouteID := FilterTomates.ByID(routeID)
	notRouteID := FilterTomates.NotID(routeID)
	byBodyColor := FilterTomates.ByColor(jsonReqBody.Color)
	updateColor := SetterTomates.SetColor(jsonReqBody.Color)
	t.backend.Transact(func(backend *Tomates) {
		byID := backend.Filter(byRouteID)
		if byID.Empty() {
			err = &NotFoundError{errors.New("ID does not exists")}
			return
		}
		if backend.Filter(byBodyColor, notRouteID).NotEmpty() {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		jsonResBody = byID.Map(updateColor).First()
	})
	if jsonResBody == nil && err == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, err
}

// Remove an existing Tomate
//
// @route /remove/{id:[0-9]+}
// @methods POST
func (t Controller) Remove(ctx context.Context, routeID string) (jsonResBody bool, err error) {
	byRouteID := FilterTomates.ByID(routeID)
	t.backend.Transact(func(backend *Tomates) {
		byID := backend.Filter(byRouteID)
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
//go:generate ggt -c slicer *SimilarTomate:SimilarTomates
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
2017/05/27 20:31:11 handling  GetByID /read/{id:[0-9]+}
2017-05-27 20:31:11.53815307 +0200 CEST [begin RestController GetByID] <nil>
2017-05-27 20:31:11.538216728 +0200 CEST [input route id 0 RestController GetByID] <nil>
2017-05-27 20:31:11.538337176 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 25

{"ID":"0","Color":"Red"}
+ curl -s -D - http://localhost:8080/read/10
2017/05/27 20:31:11 handling  GetByID /read/{id:[0-9]+}
2017-05-27 20:31:11.548527088 +0200 CEST [begin RestController GetByID] <nil>
2017-05-27 20:31:11.548630467 +0200 CEST [input route id 10 RestController GetByID] <nil>
2017-05-27 20:31:11.548674527 +0200 CEST [business error RestController GetByID] Tomate not found
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 17

Tomate not found
+ curl -s -D - --data color=blue http://localhost:8080/create
2017/05/27 20:31:11 handling  Create /create
2017-05-27 20:31:11.563432526 +0200 CEST [begin RestController Create] <nil>
2017-05-27 20:31:11.563579166 +0200 CEST [input form color blue RestController Create] <nil>
2017-05-27 20:31:11.563644979 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 26

{"ID":"1","Color":"blue"}
+ curl -s -D - --data color=blue http://localhost:8080/create
2017/05/27 20:31:11 handling  Create /create
2017-05-27 20:31:11.57801844 +0200 CEST [begin RestController Create] <nil>
2017-05-27 20:31:11.578157491 +0200 CEST [input form color blue RestController Create] <nil>
2017-05-27 20:31:11.578206942 +0200 CEST [business error RestController Create] color must be unique
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 21

color must be unique
+ curl -s -D - --data color= http://localhost:8080/create
2017/05/27 20:31:11 handling  Create /create
2017-05-27 20:31:11.593481731 +0200 CEST [begin RestController Create] <nil>
2017-05-27 20:31:11.593607477 +0200 CEST [input form color  RestController Create] <nil>
2017-05-27 20:31:11.59362383 +0200 CEST [business error RestController Create] color must not be empty
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 24

color must not be empty
+ curl -s -D - --data color=green http://localhost:8080/create
2017/05/27 20:31:11 handling  Create /create
2017-05-27 20:31:11.608433886 +0200 CEST [begin RestController Create] <nil>
2017-05-27 20:31:11.608597304 +0200 CEST [input form color green RestController Create] <nil>
2017-05-27 20:31:11.608704329 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 27

{"ID":"2","Color":"green"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/27 20:31:11 handling  Update /write/{id:[0-9]+}
2017-05-27 20:31:11.621296112 +0200 CEST [begin RestController Update] <nil>
2017-05-27 20:31:11.621388077 +0200 CEST [input route id 1 RestController Update] <nil>
2017-05-27 20:31:11.621551725 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - http://localhost:8080/read/1
2017/05/27 20:31:11 handling  GetByID /read/{id:[0-9]+}
2017-05-27 20:31:11.632576402 +0200 CEST [begin RestController GetByID] <nil>
2017-05-27 20:31:11.632675404 +0200 CEST [input route id 1 RestController GetByID] <nil>
2017-05-27 20:31:11.632760243 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/0
2017/05/27 20:31:11 handling  Update /write/{id:[0-9]+}
2017-05-27 20:31:11.647757713 +0200 CEST [begin RestController Update] <nil>
2017-05-27 20:31:11.647840825 +0200 CEST [input route id 0 RestController Update] <nil>
2017-05-27 20:31:11.64802311 +0200 CEST [business error RestController Update] color must be unique
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 21

color must be unique
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/27 20:31:11 handling  Update /write/{id:[0-9]+}
2017-05-27 20:31:11.658733177 +0200 CEST [begin RestController Update] <nil>
2017-05-27 20:31:11.658813852 +0200 CEST [input route id 1 RestController Update] <nil>
2017-05-27 20:31:11.658898308 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - http://localhost:8080/read/0
2017/05/27 20:31:11 handling  GetByID /read/{id:[0-9]+}
2017-05-27 20:31:11.669750909 +0200 CEST [begin RestController GetByID] <nil>
2017-05-27 20:31:11.669869725 +0200 CEST [input route id 0 RestController GetByID] <nil>
2017-05-27 20:31:11.669974572 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 25

{"ID":"0","Color":"Red"}
+ curl -s -D - -X POST http://localhost:8080/remove/2
2017/05/27 20:31:11 handling  Remove /remove/{id:[0-9]+}
2017-05-27 20:31:11.681084314 +0200 CEST [begin RestController Remove] <nil>
2017-05-27 20:31:11.681155057 +0200 CEST [input route id 2 RestController Remove] <nil>
2017-05-27 20:31:11.681220245 +0200 CEST [end RestController Remove] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 5

true
+ curl -s -D - -X POST http://localhost:8080/remove/2
2017/05/27 20:31:11 handling  Remove /remove/{id:[0-9]+}
2017-05-27 20:31:11.69758235 +0200 CEST [begin RestController Remove] <nil>
2017-05-27 20:31:11.697676458 +0200 CEST [input route id 2 RestController Remove] <nil>
2017-05-27 20:31:11.697727176 +0200 CEST [business error RestController Remove] ID does not exists
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 19

ID does not exists
+ curl -s -D - --data color=green http://localhost:8080/create
2017/05/27 20:31:11 handling  Create /create
2017-05-27 20:31:11.708080633 +0200 CEST [begin RestController Create] <nil>
2017-05-27 20:31:11.708231859 +0200 CEST [input form color green RestController Create] <nil>
2017-05-27 20:31:11.708316332 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 27

{"ID":"2","Color":"green"}
+ curl -s -D - 'http://localhost:8080/similar/color/r?sensitive=false'
2017/05/27 20:31:11 handling  SimilarColor /similar/color/{color}
2017-05-27 20:31:11.72518365 +0200 CEST [begin RestController SimilarColor] <nil>
2017-05-27 20:31:11.725294173 +0200 CEST [input route color r RestController SimilarColor] <nil>
2017-05-27 20:31:11.725321103 +0200 CEST [input get sensitive false RestController SimilarColor] <nil>
2017-05-27 20:31:11.72555877 +0200 CEST [end RestController SimilarColor] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 120

[{"ID":"0","Color":"Red","Similarity":0.33333333333333337},{"ID":"2","Color":"green","Similarity":0.19999999999999996}]
+ curl -s -D - http://localhost:8080/similar/color/ll
2017/05/27 20:31:11 handling  SimilarColor /similar/color/{color}
2017-05-27 20:31:11.739745538 +0200 CEST [begin RestController SimilarColor] <nil>
2017-05-27 20:31:11.739876488 +0200 CEST [input route color ll RestController SimilarColor] <nil>
2017-05-27 20:31:11.740049389 +0200 CEST [end RestController SimilarColor] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 18:31:11 GMT
Content-Length: 63

[{"ID":"1","Color":"yellow","Similarity":0.33333333333333337}]
+ killall main
signal: terminated
```
