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
	t.backend.Transact(func(backend *Tomates) {
		backend.Map(func(t *Tomate) *Tomate {
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
+ go run main.go
+ sleep 1
+ curl -s -D - http://localhost:8080/read/0
2017/05/27 15:50:47 handling  GetByID /read/{id:[0-9]+}
2017-05-27 15:50:47.480852185 +0200 CEST [begin RestController GetByID] <nil>
2017-05-27 15:50:47.480892846 +0200 CEST [input route id 0 RestController GetByID] <nil>
2017-05-27 15:50:47.481040041 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 25

{"ID":"0","Color":"Red"}
+ curl -s -D - http://localhost:8080/read/10
2017/05/27 15:50:47 handling  GetByID /read/{id:[0-9]+}
2017-05-27 15:50:47.496816299 +0200 CEST [begin RestController GetByID] <nil>
2017-05-27 15:50:47.496936521 +0200 CEST [input route id 10 RestController GetByID] <nil>
2017-05-27 15:50:47.496979688 +0200 CEST [business error RestController GetByID] Tomate not found
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 17

Tomate not found
+ curl -s -D - --data color=blue http://localhost:8080/create
2017/05/27 15:50:47 handling  Create /create
2017-05-27 15:50:47.507386617 +0200 CEST [begin RestController Create] <nil>
2017-05-27 15:50:47.507519006 +0200 CEST [input form color blue RestController Create] <nil>
2017-05-27 15:50:47.507595034 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 26

{"ID":"1","Color":"blue"}
+ curl -s -D - --data color=blue http://localhost:8080/create
2017/05/27 15:50:47 handling  Create /create
2017-05-27 15:50:47.521541237 +0200 CEST [begin RestController Create] <nil>
2017-05-27 15:50:47.521677289 +0200 CEST [input form color blue RestController Create] <nil>
2017-05-27 15:50:47.521724408 +0200 CEST [business error RestController Create] color must be unique
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 21

color must be unique
+ curl -s -D - --data color= http://localhost:8080/create
2017/05/27 15:50:47 handling  Create /create
2017-05-27 15:50:47.5382096 +0200 CEST [begin RestController Create] <nil>
2017-05-27 15:50:47.538354419 +0200 CEST [input form color  RestController Create] <nil>
2017-05-27 15:50:47.538375541 +0200 CEST [business error RestController Create] color must not be empty
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 24

color must not be empty
+ curl -s -D - --data color=green http://localhost:8080/create
2017/05/27 15:50:47 handling  Create /create
2017-05-27 15:50:47.548777277 +0200 CEST [begin RestController Create] <nil>
2017-05-27 15:50:47.549011795 +0200 CEST [input form color green RestController Create] <nil>
2017-05-27 15:50:47.549136118 +0200 CEST [end RestController Create] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 27

{"ID":"2","Color":"green"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/27 15:50:47 handling  Update /write/{id:[0-9]+}
2017-05-27 15:50:47.565246246 +0200 CEST [begin RestController Update] <nil>
2017-05-27 15:50:47.565355727 +0200 CEST [input route id 1 RestController Update] <nil>
2017-05-27 15:50:47.565543037 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - http://localhost:8080/read/1
2017/05/27 15:50:47 handling  GetByID /read/{id:[0-9]+}
2017-05-27 15:50:47.580409474 +0200 CEST [begin RestController GetByID] <nil>
2017-05-27 15:50:47.58050009 +0200 CEST [input route id 1 RestController GetByID] <nil>
2017-05-27 15:50:47.580564246 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/0
2017/05/27 15:50:47 handling  Update /write/{id:[0-9]+}
2017-05-27 15:50:47.596177642 +0200 CEST [begin RestController Update] <nil>
2017-05-27 15:50:47.596238055 +0200 CEST [input route id 0 RestController Update] <nil>
2017-05-27 15:50:47.596307174 +0200 CEST [business error RestController Update] color must be unique
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 21

color must be unique
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/27 15:50:47 handling  Update /write/{id:[0-9]+}
2017-05-27 15:50:47.608282392 +0200 CEST [begin RestController Update] <nil>
2017-05-27 15:50:47.608345903 +0200 CEST [input route id 1 RestController Update] <nil>
2017-05-27 15:50:47.608421842 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 28

{"ID":"1","Color":"yellow"}
+ curl -s -D - http://localhost:8080/read/0
2017/05/27 15:50:47 handling  GetByID /read/{id:[0-9]+}
2017-05-27 15:50:47.622746194 +0200 CEST [begin RestController GetByID] <nil>
2017-05-27 15:50:47.622819561 +0200 CEST [input route id 0 RestController GetByID] <nil>
2017-05-27 15:50:47.622868029 +0200 CEST [end RestController GetByID] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 25

{"ID":"0","Color":"Red"}
+ curl -s -D - 'http://localhost:8080/similar/color/r?sensitive=false'
2017/05/27 15:50:47 handling  SimilarColor /similar/color/{color}
2017-05-27 15:50:47.635852052 +0200 CEST [begin RestController SimilarColor] <nil>
2017-05-27 15:50:47.635985266 +0200 CEST [input route color r RestController SimilarColor] <nil>
2017-05-27 15:50:47.635992526 +0200 CEST [input get sensitive false RestController SimilarColor] <nil>
2017-05-27 15:50:47.63619748 +0200 CEST [end RestController SimilarColor] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 120

[{"ID":"0","Color":"Red","Similarity":0.33333333333333337},{"ID":"2","Color":"green","Similarity":0.19999999999999996}]
+ curl -s -D - http://localhost:8080/similar/color/ll
2017/05/27 15:50:47 handling  SimilarColor /similar/color/{color}
2017-05-27 15:50:47.648583097 +0200 CEST [begin RestController SimilarColor] <nil>
2017-05-27 15:50:47.648680258 +0200 CEST [input route color ll RestController SimilarColor] <nil>
2017-05-27 15:50:47.648795283 +0200 CEST [end RestController SimilarColor] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 27 May 2017 13:50:47 GMT
Content-Length: 63

[{"ID":"1","Color":"yellow","Similarity":0.33333333333333337}]
+ killall main
signal: terminated
```
