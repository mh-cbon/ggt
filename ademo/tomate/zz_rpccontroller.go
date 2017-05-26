package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	"context"
	json "encoding/json"
	ggt "github.com/mh-cbon/ggt/lib"
	"io"
	"net/http"
	"strconv"
)

var xxe54b0b93f158f759b6dd1585bc54a7965806c533 = strconv.Atoi
var xx72debe57433968dd2f490b9bde2dee8c155e06bc = io.Copy
var xx0921489b227b5f67966848e3dde3a344935f476d = http.StatusOK

// RPCController is an httper of Controller.
// Controller of tomatoes.
type RPCController struct {
	embed Controller
	Log   ggt.HTTPLogger
}

// NewRPCController constructs an httper of Controller
func NewRPCController(embed Controller) *RPCController {
	ret := &RPCController{
		embed: embed,
		Log:   &ggt.VoidLog{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RPCController")
	return ret
}

// GetByID invoke Controller.GetByID using the request body as a json payload.
// GetByID read the Tomate of given ID
//
// @route /read/{id:[0-9]+}
func (t *RPCController) GetByID(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCController", "GetByID")
	input := struct {
		Arg0 string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCController", "GetByID")
		t.embed.Finalizer(w, r, decErr)

		return
	}

	jsonResBody, err := t.embed.GetByID(input.Arg0)

	{
		output := struct {
			Arg0 *Tomate
			Arg1 error
		}{
			Arg0: jsonResBody,
			Arg1: err,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(output)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RPCController", "GetByID")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RPCController", "GetByID")
}

// Create invoke Controller.Create using the request body as a json payload.
// Create a new Tomate
//
// @route /create
func (t *RPCController) Create(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCController", "Create")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCController", "Create")
			t.embed.Finalizer(w, r, err)

			return
		}

	}
	input := struct {
		Arg0 *string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCController", "Create")
		t.embed.Finalizer(w, r, decErr)

		return
	}

	jsonResBody, err := t.embed.Create(input.Arg0)

	{
		output := struct {
			Arg0 *Tomate
			Arg1 error
		}{
			Arg0: jsonResBody,
			Arg1: err,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(output)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RPCController", "Create")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RPCController", "Create")
}

// Update invoke Controller.Update using the request body as a json payload.
// Update an existing Tomate
//
// @route /write/{id:[0-9]+}
func (t *RPCController) Update(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCController", "Update")
	input := struct {
		Arg0 string
		Arg1 *Tomate
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCController", "Update")
		t.embed.Finalizer(w, r, decErr)

		return
	}

	jsonResBody, err := t.embed.Update(input.Arg0, input.Arg1)

	{
		output := struct {
			Arg0 *Tomate
			Arg1 error
		}{
			Arg0: jsonResBody,
			Arg1: err,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(output)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RPCController", "Update")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RPCController", "Update")
}

// Remove invoke Controller.Remove using the request body as a json payload.
// Remove an existing Tomate
//
// @route /remove/{id:[0-9]+}
func (t *RPCController) Remove(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCController", "Remove")
	input := struct {
		Arg0 context.Context
		Arg1 string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCController", "Remove")
		t.embed.Finalizer(w, r, decErr)

		return
	}

	jsonResBody, err := t.embed.Remove(input.Arg0, input.Arg1)

	{
		output := struct {
			Arg0 bool
			Arg1 error
		}{
			Arg0: jsonResBody,
			Arg1: err,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(output)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RPCController", "Remove")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RPCController", "Remove")
}

// RPCControllerDescriptor describe a *RPCController
type RPCControllerDescriptor struct {
	ggt.TypeDescriptor
	about         *RPCController
	methodGetByID *ggt.MethodDescriptor
	methodCreate  *ggt.MethodDescriptor
	methodUpdate  *ggt.MethodDescriptor
	methodRemove  *ggt.MethodDescriptor
}

// NewRPCControllerDescriptor describe a *RPCController
func NewRPCControllerDescriptor(about *RPCController) *RPCControllerDescriptor {
	ret := &RPCControllerDescriptor{about: about}
	ret.methodGetByID = &ggt.MethodDescriptor{
		Name:    "GetByID",
		Handler: about.GetByID,
		Route:   "/read/{id:[0-9]+}",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetByID)
	ret.methodCreate = &ggt.MethodDescriptor{
		Name:    "Create",
		Handler: about.Create,
		Route:   "/create",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodCreate)
	ret.methodUpdate = &ggt.MethodDescriptor{
		Name:    "Update",
		Handler: about.Update,
		Route:   "/write/{id:[0-9]+}",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodUpdate)
	ret.methodRemove = &ggt.MethodDescriptor{
		Name:    "Remove",
		Handler: about.Remove,
		Route:   "/remove/{id:[0-9]+}",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodRemove)
	return ret
}

// GetByID returns a MethodDescriptor
func (t *RPCControllerDescriptor) GetByID() *ggt.MethodDescriptor { return t.methodGetByID }

// Create returns a MethodDescriptor
func (t *RPCControllerDescriptor) Create() *ggt.MethodDescriptor { return t.methodCreate }

// Update returns a MethodDescriptor
func (t *RPCControllerDescriptor) Update() *ggt.MethodDescriptor { return t.methodUpdate }

// Remove returns a MethodDescriptor
func (t *RPCControllerDescriptor) Remove() *ggt.MethodDescriptor { return t.methodRemove }
