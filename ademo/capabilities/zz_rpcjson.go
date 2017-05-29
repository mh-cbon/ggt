package capable

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	"encoding/json"
	ggt "github.com/mh-cbon/ggt/lib"
	finder "github.com/mh-cbon/service-finder"
	"io"
	"net/http"
	"strconv"
)

var xxe18936dae8f8e2659fd300da0180855c3befecdb = strconv.Atoi
var xx37fcd18fcdc2c8eb09f5e6e8048550932a9f0ec8 = io.Copy
var xxa0ff83a6bdd0b4cbeee232916a203d81132d06c8 = http.StatusOK

// RPCJSON is an httper of JSON.
// JSON ...
type RPCJSON struct {
	embed    JSON
	Services finder.ServiceFinder
	Log      ggt.HTTPLogger
	Session  ggt.SessionStoreProvider
	Upload   ggt.Uploader
}

// NewRPCJSON constructs an httper of JSON
func NewRPCJSON(embed JSON) *RPCJSON {
	ret := &RPCJSON{
		embed:    embed,
		Services: finder.New(),
		Log:      &ggt.VoidLog{},
		Session:  &ggt.VoidSession{},
		Upload:   &ggt.FileProvider{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RPCJSON")
	return ret
}

// ReadJSONBody invoke JSON.ReadJSONBody using the request body as a json payload.
// ReadJSONBody ...
func (t *RPCJSON) ReadJSONBody(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCJSON", "ReadJSONBody")
	var jsonReqBody Whatever
	{
		input := struct {
			jsonReqBody Whatever
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCJSON", "ReadJSONBody")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		jsonReqBody = input.jsonReqBody
	}

	t.embed.ReadJSONBody(jsonReqBody)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCJSON", "ReadJSONBody")

}

// WriteJSONBody invoke JSON.WriteJSONBody using the request body as a json payload.
// WriteJSONBody ...
func (t *RPCJSON) WriteJSONBody(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCJSON", "WriteJSONBody")

	jsonResBody := t.embed.WriteJSONBody()
	output := struct {
		Arg0 Whatever
	}{
		Arg0: jsonResBody,
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(output)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RPCJSON", "WriteJSONBody")
			http.Error(w, encErr.Error(), http.StatusInternalServerError)

			return
		}

	}

	t.Log.Handle(w, r, nil, "end", "RPCJSON", "WriteJSONBody")

}

// RPCJSONDescriptor describe a *RPCJSON
type RPCJSONDescriptor struct {
	ggt.TypeDescriptor
	about               *RPCJSON
	methodReadJSONBody  *ggt.MethodDescriptor
	methodWriteJSONBody *ggt.MethodDescriptor
}

// NewRPCJSONDescriptor describe a *RPCJSON
func NewRPCJSONDescriptor(about *RPCJSON) *RPCJSONDescriptor {
	ret := &RPCJSONDescriptor{about: about}
	ret.methodReadJSONBody = &ggt.MethodDescriptor{
		Name:    "ReadJSONBody",
		Handler: about.ReadJSONBody,
		Route:   "ReadJSONBody",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodReadJSONBody)
	ret.methodWriteJSONBody = &ggt.MethodDescriptor{
		Name:    "WriteJSONBody",
		Handler: about.WriteJSONBody,
		Route:   "WriteJSONBody",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodWriteJSONBody)
	return ret
}

// ReadJSONBody returns a MethodDescriptor
func (t *RPCJSONDescriptor) ReadJSONBody() *ggt.MethodDescriptor { return t.methodReadJSONBody }

// WriteJSONBody returns a MethodDescriptor
func (t *RPCJSONDescriptor) WriteJSONBody() *ggt.MethodDescriptor { return t.methodWriteJSONBody }
