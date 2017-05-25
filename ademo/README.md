# A demo

A demo of ggt capabilities to create a service to create/read/update/delete `tomatoes`.


# TOC
- [The model](#the-model)
  - [$ tomate/model.go](#-tomatemodelgo)
- [The controller](#the-controller)
  - [$ tomate/controller.go](#-tomatecontrollergo)
- [The main](#the-main)
  - [$ main.go](#-maingo)
- [The test](#the-test)
  - [$ sh test.sh](#-sh-testsh)

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
func (t Controller) GetByID(getID string) (jsonResBody *Tomate, err error) {
	t.backend.Transact(func(backend *Tomates) {
		jsonResBody = backend.
			Filter(FilterTomates.ByID(getID)).
			First()
	})
	if jsonResBody == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, err
}

// Create a new Tomate
func (t Controller) Create(postColor *string) (jsonResBody *Tomate, err error) {
	if postColor == nil {
		return nil, &UserInputError{errors.New("Missing color parameter")}
	}
	color := strings.TrimSpace(*postColor)
	if color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	t.backend.Transact(func(backend *Tomates) {
		exist := backend.Filter(FilterTomates.ByColor(color)).Len()
		if exist > 0 {
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
	if jsonResBody == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
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

# The test

#### $ sh test.sh
```sh
+ go generate tomate/gen.go
+ CURL='curl -s -D -'
+ sleep 1
+ go run main.go
+ curl -s -D - 'http://localhost:8080/GetByID?id=0'
2017/05/25 15:17:50 handling  GetByID GetByID
2017-05-25 15:17:50.040292094 +0200 CEST [begin RestController GetByID] <nil>
2017-05-25 15:17:50.040372712 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 25

{"ID":"0","Color":"Red"}
+ curl -s -D - 'http://localhost:8080/GetByID?id=10'
2017/05/25 15:17:50 handling  GetByID GetByID
2017-05-25 15:17:50.05647081 +0200 CEST [begin RestController GetByID] <nil>
2017-05-25 15:17:50.056620853 +0200 CEST [business error RestController GetByID] Tomate not found
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 17

Tomate not found
+ curl -s -D - --data color=blue http://localhost:8080/Create
2017/05/25 15:17:50 handling  Create Create
2017-05-25 15:17:50.071306762 +0200 CEST [begin RestController Create] <nil>
2017-05-25 15:17:50.071418859 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 26

{"ID":"1","Color":"blue"}
+ curl -s -D - --data color=blue http://localhost:8080/Create
2017/05/25 15:17:50 handling  Create Create
2017-05-25 15:17:50.086844932 +0200 CEST [begin RestController Create] <nil>
2017-05-25 15:17:50.086996435 +0200 CEST [business error RestController Create] color must be unique
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 21

color must be unique
+ curl -s -D - --data color= http://localhost:8080/Create
2017/05/25 15:17:50 handling  Create Create
2017-05-25 15:17:50.101724464 +0200 CEST [begin RestController Create] <nil>
2017-05-25 15:17:50.101786766 +0200 CEST [business error RestController Create] color must not be empty
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 24

color must not be empty
+ curl -s -D - --data color=green http://localhost:8080/Create
2017/05/25 15:17:50 handling  Create Create
2017-05-25 15:17:50.118776217 +0200 CEST [begin RestController Create] <nil>
2017-05-25 15:17:50.118996412 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 27

{"ID":"2","Color":"green"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/25 15:17:50 handling  Update /write/{id:[0-9]+}
2017-05-25 15:17:50.132018446 +0200 CEST [begin RestController Update] <nil>
2017-05-25 15:17:50.132158275 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - 'http://localhost:8080/GetByID?id=1'
2017/05/25 15:17:50 handling  GetByID GetByID
2017-05-25 15:17:50.146456017 +0200 CEST [begin RestController GetByID] <nil>
2017-05-25 15:17:50.146544555 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/0
2017/05/25 15:17:50 handling  Update /write/{id:[0-9]+}
2017-05-25 15:17:50.157744032 +0200 CEST [begin RestController Update] <nil>
2017-05-25 15:17:50.157901834 +0200 CEST [business error RestController Update] Tomate not found
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 17

Tomate not found
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/25 15:17:50 handling  Update /write/{id:[0-9]+}
2017-05-25 15:17:50.170158889 +0200 CEST [begin RestController Update] <nil>
2017-05-25 15:17:50.170302657 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - 'http://localhost:8080/GetByID?id=0'
2017/05/25 15:17:50 handling  GetByID GetByID
2017-05-25 15:17:50.182788484 +0200 CEST [begin RestController GetByID] <nil>
2017-05-25 15:17:50.182888884 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 25 May 2017 13:17:50 GMT
Content-Length: 25

{"ID":"0","Color":"Red"}
+ killall main
signal: terminated
```
