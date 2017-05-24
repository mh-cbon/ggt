# A demo

A demo of ggt capabilities to create a service to read/create `tomatoes`.

# The main

main.go
```go
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

//go:generate ggt -c slicer model/Tomate:slicegen/Tomates
//go:generate ggt chaner slicegen/Tomates:slicegen/TomatesSync

//go:generate ggt http-provider controller/Tomates:controllergen/TomatesController

func main() {

	router := mux.NewRouter()

	backend := slicegen.NewTomatesSync()
	backend.Transact(func(b *slicegen.Tomates) {
		b.Push(model.Tomate{ID: fmt.Sprintf("%v", b.Len()), Color: "Red"})
	})
	log.Println("backend", backend)

	controller := controllergen.NewTomatesController(
		controller.NewTomates(
			backend,
		),
	)

	binder := lib.NewGorillaBinder(controllergen.TomatesControllerMethodSet(controller))
	binder.Wrap(controller.GetById, func(in http.HandlerFunc) http.HandlerFunc { return in })
	binder.Apply(router, controller)

	http.ListenAndServe(":8080", router)
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

type Tomates struct {
	backend slicegen.TomatesContract
}

func NewTomates(backend slicegen.TomatesContract) Tomates {
	return Tomates{backend}
}

func (t Tomates) GetById(getID string) (jsonResBody model.Tomate, err error) {
	t.backend.Map(func(x model.Tomate) model.Tomate {
		if x.ID == getID {
			jsonResBody = x
		}
		return x
	})
	return jsonResBody, err
}

func (t Tomates) Create(postColor string) (jsonResBody *model.Tomate, err error) {
	postColor = strings.TrimSpace(postColor)
	if postColor == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	tomate := model.Tomate{Color: postColor}
	t.backend.Transact(func(backend *slicegen.Tomates) {
		exist := backend.Filter(slicegen.FilterTomates.ByColor(postColor)).Len()
		if exist > 0 {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		tomate.ID = fmt.Sprintf("%v", backend.Len())
		backend.Push(model.Tomate{ID: fmt.Sprintf("%v", backend.Len()), Color: postColor})
		jsonResBody = &tomate
	})
	return jsonResBody, err
}

type UserInputError struct {
	error
}

func (t Tomates) Finalizer(w http.ResponseWriter, r *http.Request, err error) {
	if _, ok := err.(*UserInputError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
[mh-cbon@pc4 ademo] $ 2017/05/24 15:43:06 backend &{ {[{0 Red}] } 0xc4200781e0 0xc4200782a0 0xc420078240}
curl http://localhost:8080/GetById?id=0
{"ID":"0","Color":"Red"}
[mh-cbon@pc4 ademo] $ curl --data "color=blue" http://localhost:8080/Create
{"ID":"1","Color":"blue"}
[mh-cbon@pc4 ademo] $ curl --data "color=blue" http://localhost:8080/Create
2017/05/24 15:43:14 http: multiple response.WriteHeader calls
color must be unique
null
[mh-cbon@pc4 ademo] $ curl --data "color=" http://localhost:8080/Create
2017/05/24 15:43:17 http: multiple response.WriteHeader calls
color must not be empty
null
[mh-cbon@pc4 ademo] $ curl --data "color=green" http://localhost:8080/Create
{"ID":"2","Color":"green"}
[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=1
{"ID":"1","Color":"blue"}
[mh-cbon@pc4 ademo] $ curl http://localhost:8080/GetById?id=2
{"ID":"2","Color":"green"}
[mh-cbon@pc4 ademo] $ fg
go run *go
^Csignal: interrupt
[mh-cbon@pc4 ademo] $
```
