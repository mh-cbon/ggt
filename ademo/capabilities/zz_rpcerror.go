package capable

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	json "encoding/json"
	ggt "github.com/mh-cbon/ggt/lib"
	finder "github.com/mh-cbon/service-finder"
	"io"
	"net/http"
	"strconv"
)

var xx79a9ccdea5519d9dec806c6294339ebc25f553a2 = strconv.Atoi
var xxabc532f31125328928a200ef387ff0a396b548dc = io.Copy
var xxdd4c23991e2ccabab6f9ca2621445d8cc42d9020 = http.StatusOK

// RPCError is an httper of Error.
// Error ...
type RPCError struct {
	embed    Error
	Services finder.ServiceFinder
	Log      ggt.HTTPLogger
	Session  ggt.SessionStoreProvider
	Upload   ggt.Uploader
}

// NewRPCError constructs an httper of Error
func NewRPCError(embed Error) *RPCError {
	ret := &RPCError{
		embed:    embed,
		Services: finder.New(),
		Log:      &ggt.VoidLog{},
		Session:  &ggt.VoidSession{},
		Upload:   &ggt.FileProvider{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RPCError")
	return ret
}

// TheMethod invoke Error.TheMethod using the request body as a json payload.
// TheMethod ...
func (t *RPCError) TheMethod(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCError", "TheMethod")

	err := t.embed.TheMethod()
	output := struct {
		Arg0 error
	}{
		Arg0: err,
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(output)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RPCError", "TheMethod")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RPCError", "TheMethod")

}

// RPCErrorDescriptor describe a *RPCError
type RPCErrorDescriptor struct {
	ggt.TypeDescriptor
	about           *RPCError
	methodTheMethod *ggt.MethodDescriptor
}

// NewRPCErrorDescriptor describe a *RPCError
func NewRPCErrorDescriptor(about *RPCError) *RPCErrorDescriptor {
	ret := &RPCErrorDescriptor{about: about}
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
func (t *RPCErrorDescriptor) TheMethod() *ggt.MethodDescriptor { return t.methodTheMethod }
