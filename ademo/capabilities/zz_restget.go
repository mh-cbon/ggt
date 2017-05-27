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

var xx4a8739f03dae3cc32cc07dda7581e80e3ad3a9b4 = strconv.Atoi
var xxc676b0e1df702d5b1ba91765ece2046a66e5dc6c = io.Copy
var xx3f4c105077c339d0dfc6ad93f77b6c4bdda1d4f0 = http.StatusOK

// RestGet is an httper of Get.
// Get ...
type RestGet struct {
	embed   Get
	Log     ggt.HTTPLogger
	Session ggt.SessionStoreProvider
}

// NewRestGet constructs an httper of Get
func NewRestGet(embed Get) *RestGet {
	ret := &RestGet{
		embed:   embed,
		Log:     &ggt.VoidLog{},
		Session: &ggt.VoidSession{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RestGet")
	return ret
}

// GetAll invoke Get.GetAll using the request body as a json payload.
// GetAll values in url query as a map of values
func (t *RestGet) GetAll(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetAll")

	xxURLValues := r.URL.Query()
	getValues := xxURLValues

	t.embed.GetAll(getValues)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetAll")
}

// GetAll2 invoke Get.GetAll2 using the request body as a json payload.
// GetAll2 values in url query as a map of value
func (t *RestGet) GetAll2(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetAll2")

	xxURLValues := r.URL.Query()
	getValues := map[string]string{}
	{
		for k, v := range xxURLValues {
			if len(v) > 0 {
				getValues[k] = v[0]
			}
		}
	}

	t.embed.GetAll2(getValues)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetAll2")
}

// GetOne invoke Get.GetOne using the request body as a json payload.
// GetOne arg from url query
func (t *RestGet) GetOne(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetOne")

	xxURLValues := r.URL.Query()
	var getArg1 string
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues.Get("arg1")
		t.Log.Handle(w, r, nil, "input", "get", "arg1", xxTmpgetArg1, "RestGet", "GetOne")
		getArg1 = xxTmpgetArg1
	}

	t.embed.GetOne(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetOne")
}

// GetMany invoke Get.GetMany using the request body as a json payload.
// GetMany args from url query
func (t *RestGet) GetMany(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetMany")

	xxURLValues := r.URL.Query()
	var getArg1 string
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues.Get("arg1")
		t.Log.Handle(w, r, nil, "input", "get", "arg1", xxTmpgetArg1, "RestGet", "GetMany")
		getArg1 = xxTmpgetArg1
	}
	var getArg2 string
	if _, ok := xxURLValues["arg2"]; ok {
		xxTmpgetArg2 := xxURLValues.Get("arg2")
		t.Log.Handle(w, r, nil, "input", "get", "arg2", xxTmpgetArg2, "RestGet", "GetMany")
		getArg2 = xxTmpgetArg2
	}

	t.embed.GetMany(getArg1, getArg2)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetMany")
}

// ConvertToInt invoke Get.ConvertToInt using the request body as a json payload.
// ConvertToInt an arg from url query
func (t *RestGet) ConvertToInt(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "ConvertToInt")

	xxURLValues := r.URL.Query()
	var getArg1 int
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues.Get("arg1")
		t.Log.Handle(w, r, nil, "input", "get", "arg1", xxTmpgetArg1, "RestGet", "ConvertToInt")
		{
			var err error
			getArg1, err = strconv.Atoi(xxTmpgetArg1)

			if err != nil {

				t.Log.Handle(w, r, err, "RestGet", "ConvertToInt", "get", "error", "RestGet", "ConvertToInt")
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

		}

	}

	t.embed.ConvertToInt(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "ConvertToInt")
}

// ConvertToBool invoke Get.ConvertToBool using the request body as a json payload.
// ConvertToBool an arg from url query
func (t *RestGet) ConvertToBool(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "ConvertToBool")

	xxURLValues := r.URL.Query()
	var getArg1 bool
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues.Get("arg1")
		t.Log.Handle(w, r, nil, "input", "get", "arg1", xxTmpgetArg1, "RestGet", "ConvertToBool")
		{
			var err error
			getArg1, err = strconv.ParseBool(xxTmpgetArg1)

			if err != nil {

				t.Log.Handle(w, r, err, "RestGet", "ConvertToBool", "get", "error", "RestGet", "ConvertToBool")
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

		}

	}

	t.embed.ConvertToBool(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "ConvertToBool")
}

// ConvertToSlice invoke Get.ConvertToSlice using the request body as a json payload.
// ConvertToSlice an arg from url query
func (t *RestGet) ConvertToSlice(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "ConvertToSlice")

	xxURLValues := r.URL.Query()
	var getArg1 []bool
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues["arg1"]
		t.Log.Handle(w, r, nil, "input", "get", "arg1", "RestGet", "ConvertToSlice")

		for _, xxValueTemp := range xxTmpgetArg1 {
			var xxNewValueTemp bool
			{
				var err error
				xxNewValueTemp, err = strconv.ParseBool(xxValueTemp)

				if err != nil {

					t.Log.Handle(w, r, err, "error", "RestGet", "ConvertToSlice")
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

			}

			getArg1 = append(getArg1, xxNewValueTemp)
		}

	}

	t.embed.ConvertToSlice(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "ConvertToSlice")
}

// MaybeGet invoke Get.MaybeGet using the request body as a json payload.
// MaybeGet an arg if it exists in url query.
func (t *RestGet) MaybeGet(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "MaybeGet")

	xxURLValues := r.URL.Query()
	var getArg1 *string
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues.Get("arg1")
		t.Log.Handle(w, r, nil, "input", "get", "arg1", xxTmpgetArg1, "RestGet", "MaybeGet")
		getArg1 = &xxTmpgetArg1
	}

	t.embed.MaybeGet(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "MaybeGet")
}

// RestGetDescriptor describe a *RestGet
type RestGetDescriptor struct {
	ggt.TypeDescriptor
	about                *RestGet
	methodGetAll         *ggt.MethodDescriptor
	methodGetAll2        *ggt.MethodDescriptor
	methodGetOne         *ggt.MethodDescriptor
	methodGetMany        *ggt.MethodDescriptor
	methodConvertToInt   *ggt.MethodDescriptor
	methodConvertToBool  *ggt.MethodDescriptor
	methodConvertToSlice *ggt.MethodDescriptor
	methodMaybeGet       *ggt.MethodDescriptor
}

// NewRestGetDescriptor describe a *RestGet
func NewRestGetDescriptor(about *RestGet) *RestGetDescriptor {
	ret := &RestGetDescriptor{about: about}
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
func (t *RestGetDescriptor) GetAll() *ggt.MethodDescriptor { return t.methodGetAll }

// GetAll2 returns a MethodDescriptor
func (t *RestGetDescriptor) GetAll2() *ggt.MethodDescriptor { return t.methodGetAll2 }

// GetOne returns a MethodDescriptor
func (t *RestGetDescriptor) GetOne() *ggt.MethodDescriptor { return t.methodGetOne }

// GetMany returns a MethodDescriptor
func (t *RestGetDescriptor) GetMany() *ggt.MethodDescriptor { return t.methodGetMany }

// ConvertToInt returns a MethodDescriptor
func (t *RestGetDescriptor) ConvertToInt() *ggt.MethodDescriptor { return t.methodConvertToInt }

// ConvertToBool returns a MethodDescriptor
func (t *RestGetDescriptor) ConvertToBool() *ggt.MethodDescriptor { return t.methodConvertToBool }

// ConvertToSlice returns a MethodDescriptor
func (t *RestGetDescriptor) ConvertToSlice() *ggt.MethodDescriptor { return t.methodConvertToSlice }

// MaybeGet returns a MethodDescriptor
func (t *RestGetDescriptor) MaybeGet() *ggt.MethodDescriptor { return t.methodMaybeGet }
