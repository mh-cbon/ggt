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
    - [$ tomate/zz_tomatessync.go](#-tomatezz_tomatessyncgo)
    - [$ tomate/zz_tomates.go](#-tomatezz_tomatesgo)
  - [an http rpc implementation](#an-http-rpc-implementation)
    - [$ tomate/zz_rpccontroller.go](#-tomatezz_rpccontrollergo)
    - [$ tomate/zz_rpcclient.go](#-tomatezz_rpcclientgo)
  - [an http rest implementation](#an-http-rest-implementation)
    - [$ tomate/zz_restcontroller.go](#-tomatezz_restcontrollergo)
    - [$ tomate/zz_restclient.go](#-tomatezz_restclientgo)
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
//go:generate ggt -mode rpc http-provider Controller:RPCController

//go:generate ggt -mode route http-consumer Controller:RestClient
//go:generate ggt -mode rpc http-consumer Controller:RPCClient
```

# The code for free

## a backend in-memory
#### $ tomate/zz_tomatessync.go
```go
package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

// TomatesSync is channeled.
type TomatesSync struct {
	embed Tomates
	ops   chan func()
	stop  chan bool
	tick  chan bool
}

// NewTomatesSync constructs a channeled version of Tomates
func NewTomatesSync() *TomatesSync {
	ret := &TomatesSync{
		ops:  make(chan func()),
		tick: make(chan bool),
		stop: make(chan bool),
	}
	go ret.Start()
	return ret
}

// Push is channeled
func (t *TomatesSync) Push(x ...*Tomate) *Tomates {
	var retVar0 *Tomates
	t.ops <- func() {
		retVar0 = t.embed.Push(x...)
	}
	<-t.tick
	return retVar0
}

// Unshift is channeled
func (t *TomatesSync) Unshift(x ...*Tomate) *Tomates {
	var retVar1 *Tomates
	t.ops <- func() {
		retVar1 = t.embed.Unshift(x...)
	}
	<-t.tick
	return retVar1
}

// Pop is channeled
func (t *TomatesSync) Pop() *Tomate {
	var retVar2 *Tomate
	t.ops <- func() {
		retVar2 = t.embed.Pop()
	}
	<-t.tick
	return retVar2
}

// Shift is channeled
func (t *TomatesSync) Shift() *Tomate {
	var retVar3 *Tomate
	t.ops <- func() {
		retVar3 = t.embed.Shift()
	}
	<-t.tick
	return retVar3
}

// Index is channeled
func (t *TomatesSync) Index(s *Tomate) int {
	var retVar4 int
	t.ops <- func() {
		retVar4 = t.embed.Index(s)
	}
	<-t.tick
	return retVar4
}

// Contains is channeled
func (t *TomatesSync) Contains(s *Tomate) bool {
	var retVar5 bool
	t.ops <- func() {
		retVar5 = t.embed.Contains(s)
	}
	<-t.tick
	return retVar5
}

// RemoveAt is channeled
func (t *TomatesSync) RemoveAt(i int) bool {
	var retVar6 bool
	t.ops <- func() {
		retVar6 = t.embed.RemoveAt(i)
	}
	<-t.tick
	return retVar6
}

// Remove is channeled
func (t *TomatesSync) Remove(s *Tomate) bool {
	var retVar7 bool
	t.ops <- func() {
		retVar7 = t.embed.Remove(s)
	}
	<-t.tick
	return retVar7
}

// InsertAt is channeled
func (t *TomatesSync) InsertAt(i int, s *Tomate) *Tomates {
	var retVar8 *Tomates
	t.ops <- func() {
		retVar8 = t.embed.InsertAt(i, s)
	}
	<-t.tick
	return retVar8
}

// Splice is channeled
func (t *TomatesSync) Splice(start int, length int, s ...*Tomate) []*Tomate {
	var retVar9 []*Tomate
	t.ops <- func() {
		retVar9 = t.embed.Splice(start, length, s...)
	}
	<-t.tick
	return retVar9
}

// Slice is channeled
func (t *TomatesSync) Slice(start int, length int) []*Tomate {
	var retVar10 []*Tomate
	t.ops <- func() {
		retVar10 = t.embed.Slice(start, length)
	}
	<-t.tick
	return retVar10
}

// Reverse is channeled
func (t *TomatesSync) Reverse() *Tomates {
	var retVar11 *Tomates
	t.ops <- func() {
		retVar11 = t.embed.Reverse()
	}
	<-t.tick
	return retVar11
}

// Len is channeled
func (t *TomatesSync) Len() int {
	var retVar12 int
	t.ops <- func() {
		retVar12 = t.embed.Len()
	}
	<-t.tick
	return retVar12
}

// Set is channeled
func (t *TomatesSync) Set(x []*Tomate) *Tomates {
	var retVar13 *Tomates
	t.ops <- func() {
		retVar13 = t.embed.Set(x)
	}
	<-t.tick
	return retVar13
}

// Get is channeled
func (t *TomatesSync) Get() []*Tomate {
	var retVar14 []*Tomate
	t.ops <- func() {
		retVar14 = t.embed.Get()
	}
	<-t.tick
	return retVar14
}

// At is channeled
func (t *TomatesSync) At(i int) *Tomate {
	var retVar15 *Tomate
	t.ops <- func() {
		retVar15 = t.embed.At(i)
	}
	<-t.tick
	return retVar15
}

// Filter is channeled
func (t *TomatesSync) Filter(filters ...func(*Tomate) bool) *Tomates {
	var retVar16 *Tomates
	t.ops <- func() {
		retVar16 = t.embed.Filter(filters...)
	}
	<-t.tick
	return retVar16
}

// Map is channeled
func (t *TomatesSync) Map(mappers ...func(*Tomate) *Tomate) *Tomates {
	var retVar17 *Tomates
	t.ops <- func() {
		retVar17 = t.embed.Map(mappers...)
	}
	<-t.tick
	return retVar17
}

// First is channeled
func (t *TomatesSync) First() *Tomate {
	var retVar18 *Tomate
	t.ops <- func() {
		retVar18 = t.embed.First()
	}
	<-t.tick
	return retVar18
}

// Last is channeled
func (t *TomatesSync) Last() *Tomate {
	var retVar19 *Tomate
	t.ops <- func() {
		retVar19 = t.embed.Last()
	}
	<-t.tick
	return retVar19
}

// Empty is channeled
func (t *TomatesSync) Empty() bool {
	var retVar20 bool
	t.ops <- func() {
		retVar20 = t.embed.Empty()
	}
	<-t.tick
	return retVar20
}

// UnmarshalJSON is channeled
func (t *TomatesSync) UnmarshalJSON(b []byte) error {
	var retVar21 error
	t.ops <- func() {
		retVar21 = t.embed.UnmarshalJSON(b)
	}
	<-t.tick
	return retVar21
}

// MarshalJSON is channeled
func (t *TomatesSync) MarshalJSON() ([]byte, error) {
	var retVar22 []byte
	var retVar23 error
	t.ops <- func() {
		retVar22, retVar23 = t.embed.MarshalJSON()
	}
	<-t.tick
	return retVar22, retVar23
}

// Transact execute one op.
func (t *TomatesSync) Transact(f func(*Tomates)) {
	ref := &t.embed
	f(ref)
	t.embed = *ref
}

// Start the main loop
func (t *TomatesSync) Start() {
	for {
		select {
		case op := <-t.ops:
			op()
			t.tick <- true
		case <-t.stop:
			return
		}
	}
}

// Stop the main loop
func (t *TomatesSync) Stop() {
	t.stop <- true
}
```
#### $ tomate/zz_tomates.go
```go
package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	json "encoding/json"
)

// Tomates implements a typed slice of *Tomate
type Tomates struct{ items []*Tomate }

// NewTomates creates a new typed slice of *Tomate
func NewTomates() *Tomates {
	return &Tomates{items: []*Tomate{}}
}

// Push appends every *Tomate
func (t *Tomates) Push(x ...*Tomate) *Tomates {
	t.items = append(t.items, x...)
	return t
}

// Unshift prepends every *Tomate
func (t *Tomates) Unshift(x ...*Tomate) *Tomates {
	t.items = append(x, t.items...)
	return t
}

// Pop removes then returns the last *Tomate.
func (t *Tomates) Pop() *Tomate {
	var ret *Tomate
	if len(t.items) > 0 {
		ret = t.items[len(t.items)-1]
		t.items = append(t.items[:0], t.items[len(t.items)-1:]...)
	}
	return ret
}

// Shift removes then returns the first *Tomate.
func (t *Tomates) Shift() *Tomate {
	var ret *Tomate
	if len(t.items) > 0 {
		ret = t.items[0]
		t.items = append(t.items[:0], t.items[1:]...)
	}
	return ret
}

// Index of given *Tomate. It must implements Ider interface.
func (t *Tomates) Index(s *Tomate) int {
	ret := -1
	for i, item := range t.items {
		if s.GetID() == item.GetID() {
			ret = i
			break
		}
	}
	return ret
}

// Contains returns true if s in is t.
func (t *Tomates) Contains(s *Tomate) bool {
	return t.Index(s) > -1
}

// RemoveAt removes a *Tomate at index i.
func (t *Tomates) RemoveAt(i int) bool {
	if i >= 0 && i < len(t.items) {
		t.items = append(t.items[:i], t.items[i+1:]...)
		return true
	}
	return false
}

// Remove removes given *Tomate
func (t *Tomates) Remove(s *Tomate) bool {
	if i := t.Index(s); i > -1 {
		t.RemoveAt(i)
		return true
	}
	return false
}

// InsertAt adds given *Tomate at index i
func (t *Tomates) InsertAt(i int, s *Tomate) *Tomates {
	if i < 0 || i >= len(t.items) {
		return t
	}
	res := []*Tomate{}
	res = append(res, t.items[:0]...)
	res = append(res, s)
	res = append(res, t.items[i:]...)
	t.items = res
	return t
}

// Splice removes and returns a slice of *Tomate, starting at start, ending at start+length.
// If any s is provided, they are inserted in place of the removed slice.
func (t *Tomates) Splice(start int, length int, s ...*Tomate) []*Tomate {
	var ret []*Tomate
	for i := 0; i < len(t.items); i++ {
		if i >= start && i < start+length {
			ret = append(ret, t.items[i])
		}
	}
	if start >= 0 && start+length <= len(t.items) && start+length >= 0 {
		t.items = append(
			t.items[:start],
			append(s,
				t.items[start+length:]...,
			)...,
		)
	}
	return ret
}

// Slice returns a copied slice of *Tomate, starting at start, ending at start+length.
func (t *Tomates) Slice(start int, length int) []*Tomate {
	var ret []*Tomate
	if start >= 0 && start+length <= len(t.items) && start+length >= 0 {
		ret = t.items[start : start+length]
	}
	return ret
}

// Reverse the slice.
func (t *Tomates) Reverse() *Tomates {
	for i, j := 0, len(t.items)-1; i < j; i, j = i+1, j-1 {
		t.items[i], t.items[j] = t.items[j], t.items[i]
	}
	return t
}

// Len of the slice.
func (t *Tomates) Len() int {
	return len(t.items)
}

// Set the slice.
func (t *Tomates) Set(x []*Tomate) *Tomates {
	t.items = append(t.items[:0], x...)
	return t
}

// Get the slice.
func (t *Tomates) Get() []*Tomate {
	return t.items
}

// At return the item at index i.
func (t *Tomates) At(i int) *Tomate {
	return t.items[i]
}

// Filter return a new Tomates with all items satisfying f.
func (t *Tomates) Filter(filters ...func(*Tomate) bool) *Tomates {
	ret := NewTomates()
	for _, i := range t.items {
		ok := true
		for _, f := range filters {
			ok = ok && f(i)
			if !ok {
				break
			}
		}
		if ok {
			ret.Push(i)
		}
	}
	return ret
}

// Map return a new Tomates of each items modified by f.
func (t *Tomates) Map(mappers ...func(*Tomate) *Tomate) *Tomates {
	ret := NewTomates()
	for _, i := range t.items {
		val := i
		for _, m := range mappers {
			val = m(val)
			if val == nil {
				break
			}
		}
		if val != nil {
			ret.Push(val)
		}
	}
	return ret
}

// First returns the first value or default.
func (t *Tomates) First() *Tomate {
	var ret *Tomate
	if len(t.items) > 0 {
		ret = t.items[0]
	}
	return ret
}

// Last returns the last value or default.
func (t *Tomates) Last() *Tomate {
	var ret *Tomate
	if len(t.items) > 0 {
		ret = t.items[len(t.items)-1]
	}
	return ret
}

// Empty returns true if the slice is empty.
func (t *Tomates) Empty() bool {
	return len(t.items) == 0
}

// Transact execute one op.
func (t *Tomates) Transact(f func(*Tomates)) {
	f(t)
}

//UnmarshalJSON JSON unserializes Tomates
func (t *Tomates) UnmarshalJSON(b []byte) error {
	var items []*Tomate
	if err := json.Unmarshal(b, &items); err != nil {
		return err
	}
	t.items = items
	return nil
}

//MarshalJSON JSON serializes Tomates
func (t *Tomates) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.items)
}

// TomatesContract are the requirements of Tomates
type TomatesContract interface {
	Push(x ...*Tomate) *Tomates
	Unshift(x ...*Tomate) *Tomates
	Pop() *Tomate
	Shift() *Tomate
	Index(s *Tomate) int
	Contains(s *Tomate) bool
	RemoveAt(i int) bool
	Remove(s *Tomate) bool
	InsertAt(i int, s *Tomate) *Tomates
	Splice(start int, length int, s ...*Tomate) []*Tomate
	Slice(start int, length int) []*Tomate
	Reverse() *Tomates
	Set(x []*Tomate) *Tomates
	Get() []*Tomate
	At(i int) *Tomate
	Filter(filters ...func(*Tomate) bool) *Tomates
	Map(mappers ...func(*Tomate) *Tomate) *Tomates
	First() *Tomate
	Last() *Tomate
	Transact(func(*Tomates))
	Len() int
	Empty() bool
}

// FilterTomates provides filters for a struct.
var FilterTomates = struct {
	ByID    func(string) func(*Tomate) bool
	ByColor func(string) func(*Tomate) bool
}{
	ByID:    func(v string) func(*Tomate) bool { return func(o *Tomate) bool { return o.ID == v } },
	ByColor: func(v string) func(*Tomate) bool { return func(o *Tomate) bool { return o.Color == v } },
}

// SetterTomates provides sets properties.
var SetterTomates = struct {
	SetID    func(string) func(*Tomate) *Tomate
	SetColor func(string) func(*Tomate) *Tomate
}{
	SetID: func(v string) func(*Tomate) *Tomate {
		return func(o *Tomate) *Tomate {
			o.ID = v
			return o
		}
	},
	SetColor: func(v string) func(*Tomate) *Tomate {
		return func(o *Tomate) *Tomate {
			o.Color = v
			return o
		}
	},
}
```

## an http rpc implementation
#### $ tomate/zz_rpccontroller.go
```go
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
```
#### $ tomate/zz_rpcclient.go
```go
package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	"bytes"
	"context"
	json "encoding/json"
	"errors"
	"net/http"
)

// RPCClient is an http-clienter of Controller.
// Controller of tomatoes.
type RPCClient struct {
	client *http.Client
}

// NewRPCClient constructs an http-clienter of Controller
func NewRPCClient(client *http.Client) *RPCClient {
	if client == nil {
		client = http.DefaultClient
	}
	ret := &RPCClient{
		client: client,
	}
	return ret
}

// GetByID constructs a request to GetByID
func (t RPCClient) GetByID(routeID string) (*Tomate, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 string
		}{
			Arg0: routeID,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return nil, errors.New("todo")
		}

	}
	finalURL := "/GetByID"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return nil, errors.New("todo")
	}

	output := struct {
		Arg0 *Tomate
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}

// Create constructs a request to Create
func (t RPCClient) Create(postColor *string) (*Tomate, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 *string
		}{
			Arg0: postColor,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return nil, errors.New("todo")
		}

	}
	finalURL := "/Create"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return nil, errors.New("todo")
	}

	output := struct {
		Arg0 *Tomate
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}

// Update constructs a request to Update
func (t RPCClient) Update(routeID string, jsonReqBody *Tomate) (*Tomate, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 string
			Arg1 *Tomate
		}{
			Arg0: routeID,
			Arg1: jsonReqBody,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return nil, errors.New("todo")
		}

	}
	finalURL := "/Update"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return nil, errors.New("todo")
	}

	output := struct {
		Arg0 *Tomate
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}

// Remove constructs a request to Remove
func (t RPCClient) Remove(ctx context.Context, routeID string) (bool, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 context.Context
			Arg1 string
		}{
			Arg0: ctx,
			Arg1: routeID,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return false, errors.New("todo")
		}

	}
	finalURL := "/Remove"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return false, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return false, errors.New("todo")
	}

	output := struct {
		Arg0 bool
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return false, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}
```

## an http rest implementation
#### $ tomate/zz_restcontroller.go
```go
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

var xxc52d52149b9b88887c5b162b670197217023167c = strconv.Atoi
var xxdddd5c3cff35fd51b960096f826c90c635b524c8 = io.Copy
var xx50b52bbd0311cf9ec8e406218f58ce92f00e2679 = http.StatusOK

// RestController is an httper of Controller.
// Controller of tomatoes.
type RestController struct {
	embed Controller
	Log   ggt.HTTPLogger
}

// NewRestController constructs an httper of Controller
func NewRestController(embed Controller) *RestController {
	ret := &RestController{
		embed: embed,
		Log:   &ggt.VoidLog{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RestController")
	return ret
}

// GetByID invoke Controller.GetByID using the request body as a json payload.
// GetByID read the Tomate of given ID
//
// @route /read/{id:[0-9]+}
func (t *RestController) GetByID(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestController", "GetByID")

	xxRouteVars := mux.Vars(r)
	var routeID string
	if _, ok := xxRouteVars["id"]; ok {
		xxTmprouteID := xxRouteVars["id"]
		routeID = xxTmprouteID
	}

	jsonResBody, err := t.embed.GetByID(routeID)

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RestController", "GetByID")
		t.embed.Finalizer(w, r, err)

		return
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(jsonResBody)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RestController", "GetByID")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RestController", "GetByID")
}

// Create invoke Controller.Create using the request body as a json payload.
// Create a new Tomate
//
// @route /create
func (t *RestController) Create(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestController", "Create")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RestController", "Create")
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

		t.Log.Handle(w, r, err, "business", "error", "RestController", "Create")
		t.embed.Finalizer(w, r, err)

		return
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(jsonResBody)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RestController", "Create")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RestController", "Create")
}

// Update invoke Controller.Update using the request body as a json payload.
// Update an existing Tomate
//
// @route /write/{id:[0-9]+}
func (t *RestController) Update(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestController", "Update")

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

			t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RestController", "Update")
			t.embed.Finalizer(w, r, decErr)

			return
		}

		defer r.Body.Close()
	}

	jsonResBody, err := t.embed.Update(routeID, jsonReqBody)

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RestController", "Update")
		t.embed.Finalizer(w, r, err)

		return
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(jsonResBody)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RestController", "Update")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RestController", "Update")
}

// Remove invoke Controller.Remove using the request body as a json payload.
// Remove an existing Tomate
//
// @route /remove/{id:[0-9]+}
func (t *RestController) Remove(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestController", "Remove")

	xxRouteVars := mux.Vars(r)
	ctx := r.Context()
	var routeID string
	if _, ok := xxRouteVars["id"]; ok {
		xxTmprouteID := xxRouteVars["id"]
		routeID = xxTmprouteID
	}

	jsonResBody, err := t.embed.Remove(ctx, routeID)

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RestController", "Remove")
		t.embed.Finalizer(w, r, err)

		return
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(jsonResBody)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RestController", "Remove")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RestController", "Remove")
}

// RestControllerDescriptor describe a *RestController
type RestControllerDescriptor struct {
	ggt.TypeDescriptor
	about         *RestController
	methodGetByID *ggt.MethodDescriptor
	methodCreate  *ggt.MethodDescriptor
	methodUpdate  *ggt.MethodDescriptor
	methodRemove  *ggt.MethodDescriptor
}

// NewRestControllerDescriptor describe a *RestController
func NewRestControllerDescriptor(about *RestController) *RestControllerDescriptor {
	ret := &RestControllerDescriptor{about: about}
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
func (t *RestControllerDescriptor) GetByID() *ggt.MethodDescriptor { return t.methodGetByID }

// Create returns a MethodDescriptor
func (t *RestControllerDescriptor) Create() *ggt.MethodDescriptor { return t.methodCreate }

// Update returns a MethodDescriptor
func (t *RestControllerDescriptor) Update() *ggt.MethodDescriptor { return t.methodUpdate }

// Remove returns a MethodDescriptor
func (t *RestControllerDescriptor) Remove() *ggt.MethodDescriptor { return t.methodRemove }
```
#### $ tomate/zz_restclient.go
```go
package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	"bytes"
	"context"
	json "encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RestClient is an http-clienter of Controller.
// Controller of tomatoes.
type RestClient struct {
	router *mux.Router
	client *http.Client
}

// NewRestClient constructs an http-clienter of Controller
func NewRestClient(router *mux.Router, client *http.Client) *RestClient {
	if client == nil {
		client = http.DefaultClient
	}
	ret := &RestClient{
		router: router,
		client: client,
	}
	return ret
}

// GetByID constructs a request to /read/{id:[0-9]+}
func (t RestClient) GetByID(routeID string) (jsonResBody *Tomate, err error) {
	sReqURL := "/read/{id:[0-9]+}"
	sReqURL = strings.Replace(sReqURL, "{id:[0-9]+}", fmt.Sprintf("%v", routeID), 1)
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return nil, errors.New("todo")
	}
	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	{
		res, resErr := t.client.Do(req)
		if resErr != nil {
			return nil, errors.New("todo")
		}

		decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return jsonResBody, err
}

// Create constructs a request to /create
func (t RestClient) Create(postColor *string) (jsonResBody *Tomate, err error) {
	sReqURL := "/create"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return nil, errors.New("todo")
	}
	form := url.Values{}
	form.Add("color", *postColor)
	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, strings.NewReader(form.Encode()))
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	{
		res, resErr := t.client.Do(req)
		if resErr != nil {
			return nil, errors.New("todo")
		}

		decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return jsonResBody, err
}

// Update constructs a request to /write/{id:[0-9]+}
func (t RestClient) Update(routeID string, jsonReqBody *Tomate) (jsonResBody *Tomate, err error) {

	var body io.ReadWriter
	{
		var b bytes.Buffer
		body = &b
		encErr := json.NewEncoder(body).Encode(jsonReqBody)
		if encErr != nil {
			return nil, errors.New("todo")
		}

	}
	sReqURL := "/write/{id:[0-9]+}"
	sReqURL = strings.Replace(sReqURL, "{id:[0-9]+}", fmt.Sprintf("%v", routeID), 1)
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return nil, errors.New("todo")
	}
	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, body)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	{
		res, resErr := t.client.Do(req)
		if resErr != nil {
			return nil, errors.New("todo")
		}

		decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return jsonResBody, err
}

// Remove constructs a request to /remove/{id:[0-9]+}
func (t RestClient) Remove(ctx context.Context, routeID string) (jsonResBody bool, err error) {
	sReqURL := "/remove/{id:[0-9]+}"
	sReqURL = strings.Replace(sReqURL, "{id:[0-9]+}", fmt.Sprintf("%v", routeID), 1)
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return false, errors.New("todo")
	}
	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return false, errors.New("todo")
	}

	{
		res, resErr := t.client.Do(req)
		if resErr != nil {
			return false, errors.New("todo")
		}

		decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
		if decErr != nil {
			return false, errors.New("todo")
		}

	}

	return jsonResBody, err
}
```

# The test

#### $ sh test.sh
```sh
+ go generate tomate/gen.go
+ CURL='curl -s -D -'
+ go run main.go
+ sleep 1
+ curl -s -D - 'http://localhost:8080/GetByID?id=0'
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

404 page not found
+ curl -s -D - 'http://localhost:8080/GetByID?id=10'
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

404 page not found
+ curl -s -D - --data color=blue http://localhost:8080/Create
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

404 page not found
+ curl -s -D - --data color=blue http://localhost:8080/Create
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

404 page not found
+ curl -s -D - --data color= http://localhost:8080/Create
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

404 page not found
+ curl -s -D - --data color=green http://localhost:8080/Create
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

404 page not found
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/26 16:59:50 handling  Update /write/{id:[0-9]+}
2017-05-26 16:59:50.648310101 +0200 CEST [begin RestController Update] <nil>
2017-05-26 16:59:50.648496012 +0200 CEST [business error RestController Update] ID does not exists
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

ID does not exists
+ curl -s -D - 'http://localhost:8080/GetByID?id=1'
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

404 page not found
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/0
2017/05/26 16:59:50 handling  Update /write/{id:[0-9]+}
2017-05-26 16:59:50.674293249 +0200 CEST [begin RestController Update] <nil>
2017-05-26 16:59:50.67461622 +0200 CEST [end RestController Update] <nil>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 28

{"ID":"0","Color":"yellow"}
+ curl -s -D - -H 'Content-Type: application/json' -X POST -d '{"color":"yellow"}' http://localhost:8080/write/1
2017/05/26 16:59:50 handling  Update /write/{id:[0-9]+}
2017-05-26 16:59:50.687599503 +0200 CEST [begin RestController Update] <nil>
2017-05-26 16:59:50.687766253 +0200 CEST [business error RestController Update] ID does not exists
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

ID does not exists
+ curl -s -D - 'http://localhost:8080/GetByID?id=0'
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 26 May 2017 14:59:50 GMT
Content-Length: 19

404 page not found
+ killall main
signal: terminated
```
