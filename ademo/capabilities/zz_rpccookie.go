package capable

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	ggt "github.com/mh-cbon/ggt/lib"
	finder "github.com/mh-cbon/service-finder"
	"io"
	"net/http"
	"strconv"
	"time"
)

var xxe8631f0cd9d76afc4c738557dc9b175b8799f201 = strconv.Atoi
var xx4a5f3177cb7d6d4042f775d8455f098f0e3c23d1 = io.Copy
var xx7e0969c31e2f47ea72d3163afcb81877de90b448 = http.StatusOK

// RPCCookie is an httper of Cookie.
// Cookie ...
type RPCCookie struct {
	embed    Cookie
	Services finder.ServiceFinder
	Log      ggt.HTTPLogger
	Session  ggt.SessionStoreProvider
	Upload   ggt.Uploader
}

// NewRPCCookie constructs an httper of Cookie
func NewRPCCookie(embed Cookie) *RPCCookie {
	ret := &RPCCookie{
		embed:    embed,
		Services: finder.New(),
		Log:      &ggt.VoidLog{},
		Session:  &ggt.VoidSession{},
		Upload:   &ggt.FileProvider{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RPCCookie")
	return ret
}

// GetAll invoke Cookie.GetAll using the request body as a json payload.
// GetAll values in cookies
func (t *RPCCookie) GetAll(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "GetAll")
	var cookieValues map[string]string
	{
		for _, v := range r.Cookies() {
			cookieValues[v.Name] = v.Value
		}
	}

	t.embed.GetAll(cookieValues)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "GetAll")

}

// GetAllRaw invoke Cookie.GetAllRaw using the request body as a json payload.
// GetAllRaw  cookies
func (t *RPCCookie) GetAllRaw(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "GetAllRaw")
	cookieValues := r.Cookies()

	t.embed.GetAllRaw(cookieValues)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "GetAllRaw")

}

// GetOne invoke Cookie.GetOne using the request body as a json payload.
// GetOne value form cookies
func (t *RPCCookie) GetOne(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "GetOne")
	var cookieWhatever string
	{
		c, cookieErr := r.Cookie("whatever")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RPCCookie", "GetOne")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			cookieWhatever = c.Value
		}
	}

	t.embed.GetOne(cookieWhatever)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "GetOne")

}

// GetOneRaw invoke Cookie.GetOneRaw using the request body as a json payload.
// GetOneRaw cookie
func (t *RPCCookie) GetOneRaw(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "GetOneRaw")

	var cookieWhatever http.Cookie
	{
		c, cookieErr := r.Cookie("cookieWhatever")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RPCCookie", "GetOneRaw")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		cookieWhatever = *c
	}

	t.embed.GetOneRaw(cookieWhatever)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "GetOneRaw")

}

// MaybeGetOneRaw invoke Cookie.MaybeGetOneRaw using the request body as a json payload.
// MaybeGetOneRaw cookie
func (t *RPCCookie) MaybeGetOneRaw(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "MaybeGetOneRaw")

	var cookieWhatever *http.Cookie
	{
		c, cookieErr := r.Cookie("cookieWhatever")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RPCCookie", "MaybeGetOneRaw")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		cookieWhatever = c
	}

	t.embed.MaybeGetOneRaw(cookieWhatever)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "MaybeGetOneRaw")

}

// Write invoke Cookie.Write using the request body as a json payload.
// Write a cookie
func (t *RPCCookie) Write(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "Write")

	cookieWhatever := t.embed.Write()
	http.SetCookie(w, &cookieWhatever)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "Write")

}

// MaybeDelete invoke Cookie.MaybeDelete using the request body as a json payload.
// MaybeDelete a cookie
func (t *RPCCookie) MaybeDelete(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "MaybeDelete")

	cookieWhatever := t.embed.MaybeDelete()

	if cookieWhatever == nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "cookieWhatever",
			Expires: time.Now().Add(-time.Hour * 24 * 100),
		})
	} else {
		http.SetCookie(w, cookieWhatever)
	}

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "MaybeDelete")

}

// Delete invoke Cookie.Delete using the request body as a json payload.
// Delete a cookie
func (t *RPCCookie) Delete(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "Delete")

	cookieWhatever := t.embed.Delete()

	if cookieWhatever == nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "cookieWhatever",
			Expires: time.Now().Add(-time.Hour * 24 * 100),
		})
	} else {
		http.SetCookie(w, cookieWhatever)
	}

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "Delete")

}

// GetMany invoke Cookie.GetMany using the request body as a json payload.
// GetMany args from url query
func (t *RPCCookie) GetMany(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "GetMany")
	var cookieArg1 string
	{
		c, cookieErr := r.Cookie("arg1")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RPCCookie", "GetMany")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			cookieArg1 = c.Value
		}
	}
	var cookieArg2 string
	{
		c, cookieErr := r.Cookie("arg2")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RPCCookie", "GetMany")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			cookieArg2 = c.Value
		}
	}

	t.embed.GetMany(cookieArg1, cookieArg2)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "GetMany")

}

// ConvertToInt invoke Cookie.ConvertToInt using the request body as a json payload.
// ConvertToInt an arg from url query
func (t *RPCCookie) ConvertToInt(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "ConvertToInt")
	var cookieArg1 int
	{
		c, cookieErr := r.Cookie("arg1")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RPCCookie", "ConvertToInt")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			{
				var err error
				cookieArg1, err = strconv.Atoi(c.Value)

				if err != nil {

					t.Log.Handle(w, r, err, "route", "error", "RPCCookie", "ConvertToInt")
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

			}

		}
	}

	t.embed.ConvertToInt(cookieArg1)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "ConvertToInt")

}

// ConvertToBool invoke Cookie.ConvertToBool using the request body as a json payload.
// ConvertToBool an arg from url query
func (t *RPCCookie) ConvertToBool(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "ConvertToBool")
	var cookieArg1 bool
	{
		c, cookieErr := r.Cookie("arg1")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RPCCookie", "ConvertToBool")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			{
				var err error
				cookieArg1, err = strconv.ParseBool(c.Value)

				if err != nil {

					t.Log.Handle(w, r, err, "route", "error", "RPCCookie", "ConvertToBool")
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

			}

		}
	}

	t.embed.ConvertToBool(cookieArg1)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "ConvertToBool")

}

// MaybeGet invoke Cookie.MaybeGet using the request body as a json payload.
// MaybeGet an arg if it exists in url query.
func (t *RPCCookie) MaybeGet(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RPCCookie", "MaybeGet")
	var cookieArg1 *string
	{
		c, cookieErr := r.Cookie("arg1")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RPCCookie", "MaybeGet")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			cookieArg1 = &c.Value
		}
	}

	t.embed.MaybeGet(cookieArg1)

	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RPCCookie", "MaybeGet")

}

// RPCCookieDescriptor describe a *RPCCookie
type RPCCookieDescriptor struct {
	ggt.TypeDescriptor
	about                *RPCCookie
	methodGetAll         *ggt.MethodDescriptor
	methodGetAllRaw      *ggt.MethodDescriptor
	methodGetOne         *ggt.MethodDescriptor
	methodGetOneRaw      *ggt.MethodDescriptor
	methodMaybeGetOneRaw *ggt.MethodDescriptor
	methodWrite          *ggt.MethodDescriptor
	methodMaybeDelete    *ggt.MethodDescriptor
	methodDelete         *ggt.MethodDescriptor
	methodGetMany        *ggt.MethodDescriptor
	methodConvertToInt   *ggt.MethodDescriptor
	methodConvertToBool  *ggt.MethodDescriptor
	methodMaybeGet       *ggt.MethodDescriptor
}

// NewRPCCookieDescriptor describe a *RPCCookie
func NewRPCCookieDescriptor(about *RPCCookie) *RPCCookieDescriptor {
	ret := &RPCCookieDescriptor{about: about}
	ret.methodGetAll = &ggt.MethodDescriptor{
		Name:    "GetAll",
		Handler: about.GetAll,
		Route:   "GetAll",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetAll)
	ret.methodGetAllRaw = &ggt.MethodDescriptor{
		Name:    "GetAllRaw",
		Handler: about.GetAllRaw,
		Route:   "GetAllRaw",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetAllRaw)
	ret.methodGetOne = &ggt.MethodDescriptor{
		Name:    "GetOne",
		Handler: about.GetOne,
		Route:   "GetOne",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetOne)
	ret.methodGetOneRaw = &ggt.MethodDescriptor{
		Name:    "GetOneRaw",
		Handler: about.GetOneRaw,
		Route:   "GetOneRaw",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodGetOneRaw)
	ret.methodMaybeGetOneRaw = &ggt.MethodDescriptor{
		Name:    "MaybeGetOneRaw",
		Handler: about.MaybeGetOneRaw,
		Route:   "MaybeGetOneRaw",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodMaybeGetOneRaw)
	ret.methodWrite = &ggt.MethodDescriptor{
		Name:    "Write",
		Handler: about.Write,
		Route:   "Write",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodWrite)
	ret.methodMaybeDelete = &ggt.MethodDescriptor{
		Name:    "MaybeDelete",
		Handler: about.MaybeDelete,
		Route:   "MaybeDelete",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodMaybeDelete)
	ret.methodDelete = &ggt.MethodDescriptor{
		Name:    "Delete",
		Handler: about.Delete,
		Route:   "Delete",
		Methods: []string{},
	}
	ret.TypeDescriptor.Register(ret.methodDelete)
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
func (t *RPCCookieDescriptor) GetAll() *ggt.MethodDescriptor { return t.methodGetAll }

// GetAllRaw returns a MethodDescriptor
func (t *RPCCookieDescriptor) GetAllRaw() *ggt.MethodDescriptor { return t.methodGetAllRaw }

// GetOne returns a MethodDescriptor
func (t *RPCCookieDescriptor) GetOne() *ggt.MethodDescriptor { return t.methodGetOne }

// GetOneRaw returns a MethodDescriptor
func (t *RPCCookieDescriptor) GetOneRaw() *ggt.MethodDescriptor { return t.methodGetOneRaw }

// MaybeGetOneRaw returns a MethodDescriptor
func (t *RPCCookieDescriptor) MaybeGetOneRaw() *ggt.MethodDescriptor { return t.methodMaybeGetOneRaw }

// Write returns a MethodDescriptor
func (t *RPCCookieDescriptor) Write() *ggt.MethodDescriptor { return t.methodWrite }

// MaybeDelete returns a MethodDescriptor
func (t *RPCCookieDescriptor) MaybeDelete() *ggt.MethodDescriptor { return t.methodMaybeDelete }

// Delete returns a MethodDescriptor
func (t *RPCCookieDescriptor) Delete() *ggt.MethodDescriptor { return t.methodDelete }

// GetMany returns a MethodDescriptor
func (t *RPCCookieDescriptor) GetMany() *ggt.MethodDescriptor { return t.methodGetMany }

// ConvertToInt returns a MethodDescriptor
func (t *RPCCookieDescriptor) ConvertToInt() *ggt.MethodDescriptor { return t.methodConvertToInt }

// ConvertToBool returns a MethodDescriptor
func (t *RPCCookieDescriptor) ConvertToBool() *ggt.MethodDescriptor { return t.methodConvertToBool }

// MaybeGet returns a MethodDescriptor
func (t *RPCCookieDescriptor) MaybeGet() *ggt.MethodDescriptor { return t.methodMaybeGet }
