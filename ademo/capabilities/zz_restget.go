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
	embed Get
	Log   ggt.HTTPLogger
}

// NewRestGet constructs an httper of Get
func NewRestGet(embed Get) *RestGet {
	ret := &RestGet{
		embed: embed,
		Log:   &ggt.VoidLog{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RestGet")
	return ret
}

// GetAll invoke Get.GetAll using the request body as a json payload.
// GetAll ...
func (t *RestGet) GetAll(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetAll")

	xxURLValues := r.URL.Query()
	getValues := xxURLValues

	t.embed.GetAll(getValues)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetAll")
}

// GetAll2 invoke Get.GetAll2 using the request body as a json payload.
// GetAll2 ...
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
// GetOne ...
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
// GetMany ...
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

// GetConvertedToInt invoke Get.GetConvertedToInt using the request body as a json payload.
// GetConvertedToInt ...
func (t *RestGet) GetConvertedToInt(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetConvertedToInt")

	xxURLValues := r.URL.Query()
	var getArg1 int
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues.Get("arg1")
		t.Log.Handle(w, r, nil, "input", "get", "arg1", xxTmpgetArg1, "RestGet", "GetConvertedToInt")
		{
			var err error
			getArg1, err = strconv.Atoi(xxTmpgetArg1)

			if err != nil {

				t.Log.Handle(w, r, err, "RestGet", "GetConvertedToInt", "get", "error", "RestGet", "GetConvertedToInt")
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

		}

	}

	t.embed.GetConvertedToInt(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetConvertedToInt")
}

// GetConvertedToBool invoke Get.GetConvertedToBool using the request body as a json payload.
// GetConvertedToBool ...
func (t *RestGet) GetConvertedToBool(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetConvertedToBool")

	xxURLValues := r.URL.Query()
	var getArg1 bool
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues.Get("arg1")
		t.Log.Handle(w, r, nil, "input", "get", "arg1", xxTmpgetArg1, "RestGet", "GetConvertedToBool")
		{
			var err error
			getArg1, err = strconv.ParseBool(xxTmpgetArg1)

			if err != nil {

				t.Log.Handle(w, r, err, "RestGet", "GetConvertedToBool", "get", "error", "RestGet", "GetConvertedToBool")
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

		}

	}

	t.embed.GetConvertedToBool(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetConvertedToBool")
}

// GetConvertedToSlice invoke Get.GetConvertedToSlice using the request body as a json payload.
// GetConvertedToSlice ...
func (t *RestGet) GetConvertedToSlice(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetConvertedToSlice")

	xxURLValues := r.URL.Query()
	var getArg1 []bool
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues["arg1"]
		t.Log.Handle(w, r, nil, "input", "get", "arg1", "RestGet", "GetConvertedToSlice")

		for _, xxValueTemp := range xxTmpgetArg1 {
			var xxNewValueTemp bool
			{
				var err error
				xxNewValueTemp, err = strconv.ParseBool(xxValueTemp)

				if err != nil {

					t.Log.Handle(w, r, err, "error", "RestGet", "GetConvertedToSlice")
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

			}

			getArg1 = append(getArg1, xxNewValueTemp)
		}

	}

	t.embed.GetConvertedToSlice(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetConvertedToSlice")
}

// GetMaybe invoke Get.GetMaybe using the request body as a json payload.
// GetMaybe ...
func (t *RestGet) GetMaybe(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestGet", "GetMaybe")

	xxURLValues := r.URL.Query()
	var getArg1 *string
	if _, ok := xxURLValues["arg1"]; ok {
		xxTmpgetArg1 := xxURLValues.Get("arg1")
		t.Log.Handle(w, r, nil, "input", "get", "arg1", xxTmpgetArg1, "RestGet", "GetMaybe")
		getArg1 = &xxTmpgetArg1
	}

	t.embed.GetMaybe(getArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestGet", "GetMaybe")
}

// RestGetDescriptor describe a *RestGet
type RestGetDescriptor struct {
	ggt.TypeDescriptor
	about                     *RestGet
	methodGetAll              *ggt.MethodDescriptor
	methodGetAll2             *ggt.MethodDescriptor
	methodGetOne              *ggt.MethodDescriptor
	methodGetMany             *ggt.MethodDescriptor
	methodGetConvertedToInt   *ggt.MethodDescriptor
	methodGetConvertedToBool  *ggt.MethodDescriptor
	methodGetConvertedToSlice *ggt.MethodDescriptor
	methodGetMaybe            *ggt.MethodDescriptor
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
	ret.methodGetConvertedToInt = &ggt.MethodDescriptor{
		Name:    "GetConvertedToInt",
		Handler: about.GetConvertedToInt,
		Route:   "GetConvertedToInt",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetConvertedToInt)
	ret.methodGetConvertedToBool = &ggt.MethodDescriptor{
		Name:    "GetConvertedToBool",
		Handler: about.GetConvertedToBool,
		Route:   "GetConvertedToBool",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetConvertedToBool)
	ret.methodGetConvertedToSlice = &ggt.MethodDescriptor{
		Name:    "GetConvertedToSlice",
		Handler: about.GetConvertedToSlice,
		Route:   "GetConvertedToSlice",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetConvertedToSlice)
	ret.methodGetMaybe = &ggt.MethodDescriptor{
		Name:    "GetMaybe",
		Handler: about.GetMaybe,
		Route:   "GetMaybe",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetMaybe)
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

// GetConvertedToInt returns a MethodDescriptor
func (t *RestGetDescriptor) GetConvertedToInt() *ggt.MethodDescriptor {
	return t.methodGetConvertedToInt
}

// GetConvertedToBool returns a MethodDescriptor
func (t *RestGetDescriptor) GetConvertedToBool() *ggt.MethodDescriptor {
	return t.methodGetConvertedToBool
}

// GetConvertedToSlice returns a MethodDescriptor
func (t *RestGetDescriptor) GetConvertedToSlice() *ggt.MethodDescriptor {
	return t.methodGetConvertedToSlice
}

// GetMaybe returns a MethodDescriptor
func (t *RestGetDescriptor) GetMaybe() *ggt.MethodDescriptor { return t.methodGetMaybe }
