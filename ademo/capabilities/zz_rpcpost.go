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

var xxe95288e83d8cab7b0beba9ee37d79e43f68db57d = strconv.Atoi
var xxcdce08ad494f1315b1ae677a75ff9a333d2f7a36 = io.Copy
var xx311af88f99f154b2659e5478d5663c515febcae5 = http.StatusOK

// RPCPost is an httper of Post.
// Post ...
type RPCPost struct {
	embed    Post
	Services finder.ServiceFinder
	Log      ggt.HTTPLogger
	Session  ggt.SessionStoreProvider
	Upload   ggt.Uploader
}

// NewRPCPost constructs an httper of Post
func NewRPCPost(embed Post) *RPCPost {
	ret := &RPCPost{
		embed:    embed,
		Services: finder.New(),
		Log:      &ggt.VoidLog{},
		Session:  &ggt.VoidSession{},
		Upload:   &ggt.FileProvider{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RPCPost")
	return ret
}

// GetAll invoke Post.GetAll using the request body as a json payload.
// GetAll values from the form.
func (t *RPCPost) GetAll(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCPost", "GetAll")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCPost", "GetAll")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

	}
	var postValues map[string][]string
	{
		input := struct {
			postValues map[string][]string
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCPost", "GetAll")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		postValues = input.postValues
	}

	t.embed.GetAll(postValues)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCPost", "GetAll")

}

// GetAll2 invoke Post.GetAll2 using the request body as a json payload.
// GetAll2 values from the form.
func (t *RPCPost) GetAll2(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCPost", "GetAll2")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCPost", "GetAll2")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

	}
	var postValues map[string]string
	{
		input := struct {
			postValues map[string]string
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCPost", "GetAll2")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		postValues = input.postValues
	}

	t.embed.GetAll2(postValues)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCPost", "GetAll2")

}

// GetOne invoke Post.GetOne using the request body as a json payload.
// GetOne arg form the form.
func (t *RPCPost) GetOne(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCPost", "GetOne")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCPost", "GetOne")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

	}
	var postArg1 string
	{
		input := struct {
			postArg1 string
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCPost", "GetOne")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		postArg1 = input.postArg1
	}

	t.embed.GetOne(postArg1)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCPost", "GetOne")

}

// GetMany invoke Post.GetMany using the request body as a json payload.
// GetMany args form the form.
func (t *RPCPost) GetMany(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCPost", "GetMany")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCPost", "GetMany")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

	}
	var postArg1 string
	var postArg2 string
	{
		input := struct {
			postArg1 string
			postArg2 string
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCPost", "GetMany")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		postArg1 = input.postArg1
		postArg2 = input.postArg2
	}

	t.embed.GetMany(postArg1, postArg2)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCPost", "GetMany")

}

// ConvertToInt invoke Post.ConvertToInt using the request body as a json payload.
// ConvertToInt an arg from the form.
func (t *RPCPost) ConvertToInt(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCPost", "ConvertToInt")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCPost", "ConvertToInt")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

	}
	var postArg1 int
	{
		input := struct {
			postArg1 int
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCPost", "ConvertToInt")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		postArg1 = input.postArg1
	}

	t.embed.ConvertToInt(postArg1)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCPost", "ConvertToInt")

}

// ConvertToBool invoke Post.ConvertToBool using the request body as a json payload.
// ConvertToBool an arg from the form.
func (t *RPCPost) ConvertToBool(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCPost", "ConvertToBool")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCPost", "ConvertToBool")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

	}
	var postArg1 bool
	{
		input := struct {
			postArg1 bool
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCPost", "ConvertToBool")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		postArg1 = input.postArg1
	}

	t.embed.ConvertToBool(postArg1)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCPost", "ConvertToBool")

}

// ConvertToSlice invoke Post.ConvertToSlice using the request body as a json payload.
// ConvertToSlice an arg from the form.
func (t *RPCPost) ConvertToSlice(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCPost", "ConvertToSlice")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCPost", "ConvertToSlice")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

	}
	var postArg1 []bool
	{
		input := struct {
			postArg1 []bool
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCPost", "ConvertToSlice")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		postArg1 = input.postArg1
	}

	t.embed.ConvertToSlice(postArg1)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCPost", "ConvertToSlice")

}

// MaybeGet invoke Post.MaybeGet using the request body as a json payload.
// MaybeGet an arg if it exists in the form.
func (t *RPCPost) MaybeGet(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCPost", "MaybeGet")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RPCPost", "MaybeGet")
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

	}
	var postArg1 *string
	{
		input := struct {
			postArg1 *string
		}{}
		decErr := json.NewDecoder(r.Body).Decode(&input)

		if decErr != nil {

			t.Log.Handle(w, r, decErr, "json", "decode", "input", "error", "RPCPost", "MaybeGet")
			http.Error(w, decErr.Error(), http.StatusInternalServerError)

			return
		}

		postArg1 = input.postArg1
	}

	t.embed.MaybeGet(postArg1)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCPost", "MaybeGet")

}

// RPCPostDescriptor describe a *RPCPost
type RPCPostDescriptor struct {
	ggt.TypeDescriptor
	about                *RPCPost
	methodGetAll         *ggt.MethodDescriptor
	methodGetAll2        *ggt.MethodDescriptor
	methodGetOne         *ggt.MethodDescriptor
	methodGetMany        *ggt.MethodDescriptor
	methodConvertToInt   *ggt.MethodDescriptor
	methodConvertToBool  *ggt.MethodDescriptor
	methodConvertToSlice *ggt.MethodDescriptor
	methodMaybeGet       *ggt.MethodDescriptor
}

// NewRPCPostDescriptor describe a *RPCPost
func NewRPCPostDescriptor(about *RPCPost) *RPCPostDescriptor {
	ret := &RPCPostDescriptor{about: about}
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
	ret.methodConvertToSlice = &ggt.MethodDescriptor{
		Name:    "ConvertToSlice",
		Handler: about.ConvertToSlice,
		Route:   "ConvertToSlice",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodConvertToSlice)
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
func (t *RPCPostDescriptor) GetAll() *ggt.MethodDescriptor { return t.methodGetAll }

// GetAll2 returns a MethodDescriptor
func (t *RPCPostDescriptor) GetAll2() *ggt.MethodDescriptor { return t.methodGetAll2 }

// GetOne returns a MethodDescriptor
func (t *RPCPostDescriptor) GetOne() *ggt.MethodDescriptor { return t.methodGetOne }

// GetMany returns a MethodDescriptor
func (t *RPCPostDescriptor) GetMany() *ggt.MethodDescriptor { return t.methodGetMany }

// ConvertToInt returns a MethodDescriptor
func (t *RPCPostDescriptor) ConvertToInt() *ggt.MethodDescriptor { return t.methodConvertToInt }

// ConvertToBool returns a MethodDescriptor
func (t *RPCPostDescriptor) ConvertToBool() *ggt.MethodDescriptor { return t.methodConvertToBool }

// ConvertToSlice returns a MethodDescriptor
func (t *RPCPostDescriptor) ConvertToSlice() *ggt.MethodDescriptor { return t.methodConvertToSlice }

// MaybeGet returns a MethodDescriptor
func (t *RPCPostDescriptor) MaybeGet() *ggt.MethodDescriptor { return t.methodMaybeGet }
