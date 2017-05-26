package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	json "encoding/json"
	"github.com/gorilla/mux"
	ggt "github.com/mh-cbon/ggt/lib"
	"io"
	"net/http"
	"strconv"
)

var xxStrconvAtoi = strconv.Atoi
var xxIoCopy = io.Copy
var xxHTTPOk = http.StatusOK

// RpcController is an httper of Controller.
// Controller of tomatoes.
type RpcController struct {
	embed Controller
	Log   ggt.HTTPLogger
}

// NewRpcController constructs an httper of Controller
func NewRpcController(embed Controller) *RpcController {
	ret := &RpcController{
		embed: embed,
		Log:   &ggt.VoidLog{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RpcController")
	return ret
}

// GetByID invoke Controller.GetByID using the request body as a json payload.
// GetByID read the Tomate of given ID
func (t *RpcController) GetByID(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RpcController", "GetByID")

	xxURLValues := r.URL.Query()
	var getID string
	if _, ok := xxURLValues["id"]; ok {
		xxTmpgetID := xxURLValues.Get("id")
		getID = xxTmpgetID
	}

	jsonResBody, err := t.embed.GetByID(getID)

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RpcController", "GetByID")
		t.embed.Finalizer(w, r, err)

		return
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(jsonResBody)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RpcController", "GetByID")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RpcController", "GetByID")
}

// Create invoke Controller.Create using the request body as a json payload.
// Create a new Tomate
func (t *RpcController) Create(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RpcController", "Create")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RpcController", "Create")
			t.embed.Finalizer(w, r, err)

			return
		}

	}
	var postColor *string
	if _, ok := r.Form["color"]; ok {
		xxTmppostColor := r.FormValue("color")
		postColor = &xxTmppostColor
	}

	jsonResBody, err := t.embed.Create(postColor)

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RpcController", "Create")
		t.embed.Finalizer(w, r, err)

		return
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(jsonResBody)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RpcController", "Create")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RpcController", "Create")
}

// Update invoke Controller.Update using the request body as a json payload.
// Update an existing Tomate
//
// @route /write/{id:[0-9]+}
func (t *RpcController) Update(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RpcController", "Update")

	xxRouteVars := mux.Vars(r)
	var routeID string
	if _, ok := xxRouteVars["id"]; ok {
		xxTmprouteID := xxRouteVars["id"]
		routeID = xxTmprouteID
	}
	var jsonReqBody *Tomate
	{
		jsonReqBody = &Tomate{}
		decErr := json.NewDecoder(r.Body).Decode(jsonReqBody)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RpcController", "Update")
			t.embed.Finalizer(w, r, decErr)

			return
		}

		defer r.Body.Close()
	}

	jsonResBody, err := t.embed.Update(routeID, jsonReqBody)

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RpcController", "Update")
		t.embed.Finalizer(w, r, err)

		return
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(jsonResBody)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RpcController", "Update")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RpcController", "Update")
}

// RpcControllerDescriptor describe a *RpcController
type RpcControllerDescriptor struct {
	ggt.TypeDescriptor
	about         *RpcController
	methodGetByID *ggt.MethodDescriptor
	methodCreate  *ggt.MethodDescriptor
	methodUpdate  *ggt.MethodDescriptor
}

// NewRpcControllerDescriptor describe a *RpcController
func NewRpcControllerDescriptor(about *RpcController) *RpcControllerDescriptor {
	ret := &RpcControllerDescriptor{about: about}
	ret.methodGetByID = &ggt.MethodDescriptor{
		Name:    "GetByID",
		Handler: about.GetByID,
		Route:   "GetByID",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetByID)
	ret.methodCreate = &ggt.MethodDescriptor{
		Name:    "Create",
		Handler: about.Create,
		Route:   "Create",
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
	return ret
}

// GetByID returns a MethodDescriptor
func (t *RpcControllerDescriptor) GetByID() *ggt.MethodDescriptor { return t.methodGetByID }

// Create returns a MethodDescriptor
func (t *RpcControllerDescriptor) Create() *ggt.MethodDescriptor { return t.methodCreate }

// Update returns a MethodDescriptor
func (t *RpcControllerDescriptor) Update() *ggt.MethodDescriptor { return t.methodUpdate }
