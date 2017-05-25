# A demo

A demo of ggt capabilities to create a service to create/read/update/delete `tomatoes`.


# TOC
- [The model](#the-model)
- [The controller](#the-controller)
- [The main](#the-main)
- [The test](#the-test)

# The model

$ tomate/model.go
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

$ tomate/controller.go
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

$ main.go
```go
// A demo of ggt capabilities to create a service to create/read/update/delete `tomatoes`.
package main

import (
	"fmt"
	"log"
	"net/http"

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

```sh
[mh-cbon@pc4 ademo] $ go generate *go
2017/05/24 15:43:01 no initial packages were loaded
2017/05/24 15:43:01 no initial packages were loaded
2017/05/24 15:43:01 no initial packages were loaded

[mh-cbon@pc4 ademo] $ go run *go &
[1] 5833

[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=0
...
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 24 May 2017 13:50:28 GMT
< Content-Length: 25
<
{"ID":"0","Color":"Red"}

[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=10
...
< HTTP/1.1 404 Not Found
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Wed, 24 May 2017 13:52:26 GMT
< Content-Length: 17
<
Tomate not found

[mh-cbon@pc4 ademo] $ curl -v --data "color=blue" http://localhost:8080/Create
...
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 24 May 2017 13:49:58 GMT
< Content-Length: ...
<
{"ID":"1","Color":"blue"}

[mh-cbon@pc4 ademo] $ curl --data "color=blue" http://localhost:8080/Create
...
< HTTP/1.1 400 Bad Request
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Wed, 24 May 2017 13:49:15 GMT
< Content-Length: 21
<
color must be unique

[mh-cbon@pc4 ademo] $ curl --data "color=" http://localhost:8080/Create
...
< HTTP/1.1 400 Bad Request
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Wed, 24 May 2017 13:48:46 GMT
< Content-Length: 24
<
color must not be empty

[mh-cbon@pc4 ademo] $ curl --data "color=green" http://localhost:8080/Create
{"ID":"2","Color":"green"}

[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=1
{"ID":"1","Color":"blue"}

[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=2
{"ID":"2","Color":"green"}

[mh-cbon@pc4 ademo] $ fg
go run *go
^Csignal: interrupt
```
