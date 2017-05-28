package capable

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	json "encoding/json"
	ggt "github.com/mh-cbon/ggt/lib"
	"io"
	"net/http"
	"strconv"
)

var xx75dfadaae5b48de33b7ebc82ffab547b3bb246c4 = strconv.Atoi
var xxd1a2e54d5453fde3f436876c82267a9ca80342d9 = io.Copy
var xxd52c045e402fb2d4669649eac947db0b10904609 = http.StatusOK

// RPCURL is an httper of URL.
// URL is a merge of route, url
type RPCURL struct {
	embed   URL
	Log     ggt.HTTPLogger
	Session ggt.SessionStoreProvider
	Upload  ggt.Uploader
}

// NewRPCURL constructs an httper of URL
func NewRPCURL(embed URL) *RPCURL {
	ret := &RPCURL{
		embed:   embed,
		Log:     &ggt.VoidLog{},
		Session: &ggt.VoidSession{},
		Upload:  &ggt.FileProvider{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RPCURL")
	return ret
}

// GetAll invoke URL.GetAll using the request body as a json payload.
// GetAll  return a merged map of route, url
func (t *RPCURL) GetAll(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCURL", "GetAll")
	input := struct {
		Arg0 map[string][]string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCURL", "GetAll")
		http.Error(w, decErr.Error(), http.StatusInternalServerError)

		return
	}

	t.embed.GetAll(input.Arg0)

	t.Log.Handle(w, r, nil, "end", "RPCURL", "GetAll")
}

// GetAll2 invoke URL.GetAll2 using the request body as a json payload.
// GetAll2 return a merged map of route, url
func (t *RPCURL) GetAll2(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCURL", "GetAll2")
	input := struct {
		Arg0 map[string]string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCURL", "GetAll2")
		http.Error(w, decErr.Error(), http.StatusInternalServerError)

		return
	}

	t.embed.GetAll2(input.Arg0)

	t.Log.Handle(w, r, nil, "end", "RPCURL", "GetAll2")
}

// GetOne invoke URL.GetOne using the request body as a json payload.
// GetOne return the first value in route, url
func (t *RPCURL) GetOne(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCURL", "GetOne")
	input := struct {
		Arg0 string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCURL", "GetOne")
		http.Error(w, decErr.Error(), http.StatusInternalServerError)

		return
	}

	t.embed.GetOne(input.Arg0)

	t.Log.Handle(w, r, nil, "end", "RPCURL", "GetOne")
}

// GetMany invoke URL.GetMany using the request body as a json payload.
// GetMany return the first value of each parameter in route, url
func (t *RPCURL) GetMany(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCURL", "GetMany")
	input := struct {
		Arg0 string
		Arg1 string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCURL", "GetMany")
		http.Error(w, decErr.Error(), http.StatusInternalServerError)

		return
	}

	t.embed.GetMany(input.Arg0, input.Arg1)

	t.Log.Handle(w, r, nil, "end", "RPCURL", "GetMany")
}

// ConvertToInt invoke URL.ConvertToInt using the request body as a json payload.
// ConvertToInt an arg
func (t *RPCURL) ConvertToInt(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCURL", "ConvertToInt")
	input := struct {
		Arg0 int
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCURL", "ConvertToInt")
		http.Error(w, decErr.Error(), http.StatusInternalServerError)

		return
	}

	t.embed.ConvertToInt(input.Arg0)

	t.Log.Handle(w, r, nil, "end", "RPCURL", "ConvertToInt")
}

// ConvertToBool invoke URL.ConvertToBool using the request body as a json payload.
// ConvertToBool an arg
func (t *RPCURL) ConvertToBool(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCURL", "ConvertToBool")
	input := struct {
		Arg0 bool
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCURL", "ConvertToBool")
		http.Error(w, decErr.Error(), http.StatusInternalServerError)

		return
	}

	t.embed.ConvertToBool(input.Arg0)

	t.Log.Handle(w, r, nil, "end", "RPCURL", "ConvertToBool")
}

// MaybeGet invoke URL.MaybeGet using the request body as a json payload.
// MaybeGet an arg if it exists.
func (t *RPCURL) MaybeGet(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCURL", "MaybeGet")
	input := struct {
		Arg0 *string
	}{}
	decErr := json.NewDecoder(r.Body).Decode(&input)

	if decErr != nil {

		t.Log.Handle(w, r, decErr, "req", "json", "decode", "error", "RPCURL", "MaybeGet")
		http.Error(w, decErr.Error(), http.StatusInternalServerError)

		return
	}

	t.embed.MaybeGet(input.Arg0)

	t.Log.Handle(w, r, nil, "end", "RPCURL", "MaybeGet")
}

// RPCURLDescriptor describe a *RPCURL
type RPCURLDescriptor struct {
	ggt.TypeDescriptor
	about               *RPCURL
	methodGetAll        *ggt.MethodDescriptor
	methodGetAll2       *ggt.MethodDescriptor
	methodGetOne        *ggt.MethodDescriptor
	methodGetMany       *ggt.MethodDescriptor
	methodConvertToInt  *ggt.MethodDescriptor
	methodConvertToBool *ggt.MethodDescriptor
	methodMaybeGet      *ggt.MethodDescriptor
}

// NewRPCURLDescriptor describe a *RPCURL
func NewRPCURLDescriptor(about *RPCURL) *RPCURLDescriptor {
	ret := &RPCURLDescriptor{about: about}
	ret.methodGetAll = &ggt.MethodDescriptor{
		Name:    "GetAll",
		Handler: about.GetAll,
		Route:   "GetAll",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetAll)
	ret.methodGetAll2 = &ggt.MethodDescriptor{
		Name:    "GetAll2",
		Handler: about.GetAll2,
		Route:   "GetAll2",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetAll2)
	ret.methodGetOne = &ggt.MethodDescriptor{
		Name:    "GetOne",
		Handler: about.GetOne,
		Route:   "GetOne",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetOne)
	ret.methodGetMany = &ggt.MethodDescriptor{
		Name:    "GetMany",
		Handler: about.GetMany,
		Route:   "GetMany",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetMany)
	ret.methodConvertToInt = &ggt.MethodDescriptor{
		Name:    "ConvertToInt",
		Handler: about.ConvertToInt,
		Route:   "ConvertToInt",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodConvertToInt)
	ret.methodConvertToBool = &ggt.MethodDescriptor{
		Name:    "ConvertToBool",
		Handler: about.ConvertToBool,
		Route:   "ConvertToBool",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodConvertToBool)
	ret.methodMaybeGet = &ggt.MethodDescriptor{
		Name:    "MaybeGet",
		Handler: about.MaybeGet,
		Route:   "MaybeGet",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodMaybeGet)
	return ret
}

// GetAll returns a MethodDescriptor
func (t *RPCURLDescriptor) GetAll() *ggt.MethodDescriptor { return t.methodGetAll }

// GetAll2 returns a MethodDescriptor
func (t *RPCURLDescriptor) GetAll2() *ggt.MethodDescriptor { return t.methodGetAll2 }

// GetOne returns a MethodDescriptor
func (t *RPCURLDescriptor) GetOne() *ggt.MethodDescriptor { return t.methodGetOne }

// GetMany returns a MethodDescriptor
func (t *RPCURLDescriptor) GetMany() *ggt.MethodDescriptor { return t.methodGetMany }

// ConvertToInt returns a MethodDescriptor
func (t *RPCURLDescriptor) ConvertToInt() *ggt.MethodDescriptor { return t.methodConvertToInt }

// ConvertToBool returns a MethodDescriptor
func (t *RPCURLDescriptor) ConvertToBool() *ggt.MethodDescriptor { return t.methodConvertToBool }

// MaybeGet returns a MethodDescriptor
func (t *RPCURLDescriptor) MaybeGet() *ggt.MethodDescriptor { return t.methodMaybeGet }