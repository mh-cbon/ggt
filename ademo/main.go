package main

import (
	"fmt"
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
	desc.Create().Wrap(func(in http.HandlerFunc) http.HandlerFunc { return in })
	desc.Wrap(func(in http.HandlerFunc) http.HandlerFunc { return in })

	lib.Gorilla(desc, router)

	http.ListenAndServe(":8080", router)
}
