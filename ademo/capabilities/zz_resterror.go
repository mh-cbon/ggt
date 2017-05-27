package capable

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	ggt "github.com/mh-cbon/ggt/lib"
	"io"
	"net/http"
	"strconv"
)

var xxa959b75ba7e5fe2e3d3096f41ee0ff87f57dff2a = strconv.Atoi
var xx2d99e18666caa4551a04d4555592cabe892b93c1 = io.Copy
var xxc7840b50abe1844351cb21a3e77009c3343ead3b = http.StatusOK

// RestError is an httper of Error.
// Error ...
type RestError struct {
	embed   Error
	Log     ggt.HTTPLogger
	Session ggt.SessionStoreProvider
}

// NewRestError constructs an httper of Error
func NewRestError(embed Error) *RestError {
	ret := &RestError{
		embed:   embed,
		Log:     &ggt.VoidLog{},
		Session: &ggt.VoidSession{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RestError")
	return ret
}

// TheMethod invoke Error.TheMethod using the request body as a json payload.
// TheMethod ...
func (t *RestError) TheMethod(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestError", "TheMethod")

	err := t.embed.TheMethod()

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RestError", "TheMethod")
		t.embed.Finalizer(w, r, err)

		return
	}

	t.Log.Handle(w, r, nil, "end", "RestError", "TheMethod")
}

// RestErrorDescriptor describe a *RestError
type RestErrorDescriptor struct {
	ggt.TypeDescriptor
	about           *RestError
	methodTheMethod *ggt.MethodDescriptor
}

// NewRestErrorDescriptor describe a *RestError
func NewRestErrorDescriptor(about *RestError) *RestErrorDescriptor {
	ret := &RestErrorDescriptor{about: about}
	ret.methodTheMethod = &ggt.MethodDescriptor{
		Name:    "TheMethod",
		Handler: about.TheMethod,
		Route:   "TheMethod",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodTheMethod)
	return ret
}

// TheMethod returns a MethodDescriptor
func (t *RestErrorDescriptor) TheMethod() *ggt.MethodDescriptor { return t.methodTheMethod }
