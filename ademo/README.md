# A demo

A demo of ggt capabilities to create a service to create/read/update/delete `tomatoes`.


# TOC
- [The main](#the-main)
- [The controller](#the-controller)
- [The test](#the-test)

# The main

main.go
```go
// A demo of ggt capabilities to create a service to create/read/update/delete `tomatoes`.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mh-cbon/ggt/ademo/controller"
	"github.com/mh-cbon/ggt/ademo/controllergen"
	"github.com/mh-cbon/ggt/ademo/model"
	"github.com/mh-cbon/ggt/ademo/slicegen"
	"github.com/mh-cbon/ggt/lib"
)

//go:generate ggt -c slicer model/*Tomate:slicegen/Tomates
//go:generate ggt chaner slicegen/Tomates:slicegen/TomatesSync

//go:generate ggt http-provider controller/Tomates:controllergen/TomatesController

func main() {

	router := mux.NewRouter()

	backend := slicegen.NewTomatesSync()
	backend.Transact(func(b *slicegen.Tomates) {
		b.Push(&model.Tomate{ID: fmt.Sprintf("%v", b.Len()), Color: "Red"})
	})

	controller := controllergen.NewTomatesController(
		controller.NewTomates(
			backend,
		),
	)

	desc := controllergen.NewTomatesControllerDescriptor(controller)
	// desc.Create().WrapMethod(logReq)
	desc.WrapMethod(logReq)

	lib.Gorilla(desc, router)

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

# The controller

controller/tomate.go
```go
package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mh-cbon/ggt/ademo/model"
	"github.com/mh-cbon/ggt/ademo/slicegen"
)

// Tomates controller.
type Tomates struct {
	backend slicegen.TomatesContract
}

// NewTomates creates a new tomates controller
func NewTomates(backend slicegen.TomatesContract) Tomates {
	return Tomates{backend}
}

// GetByID read the Tomate of given ID
func (t Tomates) GetByID(getID string) (jsonResBody *model.Tomate, err error) {
	err = &NotFoundError{errors.New("Tomate not found")}
	t.backend.Map(func(x *model.Tomate) *model.Tomate {
		if x.ID == getID {
			jsonResBody = x
			err = nil
		}
		return x
	})
	return jsonResBody, err
}

// Create a new Tomate
func (t Tomates) Create(postColor *string) (jsonResBody *model.Tomate, err error) {
	if postColor == nil {
		return nil, &UserInputError{errors.New("Missing color parameter")}
	}
	color := strings.TrimSpace(*postColor)
	if color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	t.backend.Transact(func(backend *slicegen.Tomates) {
		exist := backend.Filter(slicegen.FilterTomates.ByColor(color)).Len()
		if exist > 0 {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		jsonResBody = &model.Tomate{ID: fmt.Sprintf("%v", backend.Len()), Color: color}
		backend.Push(jsonResBody)
	})
	return jsonResBody, err
}

// Update an existing Tomate
//
// @route /write/{id:[0-9]+}
func (t Tomates) Update(routeID string, jsonReqBody *model.Tomate) (jsonResBody *model.Tomate, err error) {
	jsonReqBody.Color = strings.TrimSpace(jsonReqBody.Color)
	if jsonReqBody.Color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	t.backend.Transact(func(backend *slicegen.Tomates) {
		byID := backend.Filter(slicegen.FilterTomates.ByID(routeID))
		if byID.Len() == 0 {
			err = &NotFoundError{errors.New("ID does not exists")}
			return
		}
		byColor := backend.Filter(slicegen.FilterTomates.ByColor(jsonReqBody.Color))
		if byColor.Len() > 0 && byID.First().ID != byColor.First().ID {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		backend.Map(func(x *model.Tomate) *model.Tomate {
			if x.ID == routeID {
				x.Color = jsonReqBody.Color
			}
			return x
		})
	})
	jsonResBody = t.backend.Filter(slicegen.FilterTomates.ByID(routeID)).First()
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
func (t Tomates) Finalizer(w http.ResponseWriter, r *http.Request, err error) {
	if _, ok := err.(*UserInputError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if _, ok := err.(*NotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
