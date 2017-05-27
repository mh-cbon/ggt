package capable

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	ggt "github.com/mh-cbon/ggt/lib"
	"io"
	"net/http"
	"strconv"
	"time"
)

var xx585dc844fd22cc31825cd9defc779ff2cde4391d = strconv.Atoi
var xx82052078d95bcb749e54581546ae7824a9e37255 = io.Copy
var xx77a3b78ea2f772886b53ea6b1e988ad6bf4c61d1 = http.StatusOK

// RestCookie is an httper of Cookie.
// Cookie ...
type RestCookie struct {
	embed   Cookie
	Log     ggt.HTTPLogger
	Session ggt.SessionStoreProvider
}

// NewRestCookie constructs an httper of Cookie
func NewRestCookie(embed Cookie) *RestCookie {
	ret := &RestCookie{
		embed:   embed,
		Log:     &ggt.VoidLog{},
		Session: &ggt.VoidSession{},
	}
	ret.Log.Handle(nil, nil, nil, "constructor", "RestCookie")
	return ret
}

// GetAll invoke Cookie.GetAll using the request body as a json payload.
// GetAll values in cookies
func (t *RestCookie) GetAll(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "GetAll")
	var cookieValues map[string]string
	{
		for _, v := range r.Cookies() {
			cookieValues[v.Name] = v.Value
		}
	}

	t.embed.GetAll(cookieValues)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "GetAll")
}

// GetAllRaw invoke Cookie.GetAllRaw using the request body as a json payload.
// GetAllRaw  cookies
func (t *RestCookie) GetAllRaw(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "GetAllRaw")
	cookieValues := r.Cookies()

	t.embed.GetAllRaw(cookieValues)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "GetAllRaw")
}

// GetOne invoke Cookie.GetOne using the request body as a json payload.
// GetOne value form cookies
func (t *RestCookie) GetOne(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "GetOne")
	var cookieWhatever string
	{
		c, cookieErr := r.Cookie("whatever")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestCookie", "GetOne")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			cookieWhatever = c.Value
		}
	}

	t.embed.GetOne(cookieWhatever)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "GetOne")
}

// GetOneRaw invoke Cookie.GetOneRaw using the request body as a json payload.
// GetOneRaw cookie
func (t *RestCookie) GetOneRaw(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "GetOneRaw")
	var cookieWhatever http.Cookie
	{
		c, cookieErr := r.Cookie("whatever")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestCookie", "GetOneRaw")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			cookieWhatever = *c
		}
	}

	t.embed.GetOneRaw(cookieWhatever)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "GetOneRaw")
}

// MaybeGetOneRaw invoke Cookie.MaybeGetOneRaw using the request body as a json payload.
// MaybeGetOneRaw cookie
func (t *RestCookie) MaybeGetOneRaw(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "MaybeGetOneRaw")
	var cookieWhatever *http.Cookie
	{
		c, cookieErr := r.Cookie("whatever")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestCookie", "MaybeGetOneRaw")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		cookieWhatever = c
	}

	t.embed.MaybeGetOneRaw(cookieWhatever)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "MaybeGetOneRaw")
}

// Write invoke Cookie.Write using the request body as a json payload.
// Write a cookie
func (t *RestCookie) Write(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "Write")

	cookieWhatever := t.embed.Write()
	http.SetCookie(w, &cookieWhatever)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "Write")
}

// MaybeDelete invoke Cookie.MaybeDelete using the request body as a json payload.
// MaybeDelete a cookie
func (t *RestCookie) MaybeDelete(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "MaybeDelete")

	cookieWhatever := t.embed.MaybeDelete()

	if cookieWhatever == nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "cookieWhatever",
			Expires: time.Now().Add(-time.Hour * 24 * 100),
		})
	} else {
		http.SetCookie(w, cookieWhatever)
	}

	t.Log.Handle(w, r, nil, "end", "RestCookie", "MaybeDelete")
}

// Delete invoke Cookie.Delete using the request body as a json payload.
// Delete a cookie
func (t *RestCookie) Delete(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "Delete")

	cookieWhatever := t.embed.Delete()

	if cookieWhatever == nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "cookieWhatever",
			Expires: time.Now().Add(-time.Hour * 24 * 100),
		})
	} else {
		http.SetCookie(w, cookieWhatever)
	}

	t.Log.Handle(w, r, nil, "end", "RestCookie", "Delete")
}

// GetMany invoke Cookie.GetMany using the request body as a json payload.
// GetMany args from url query
func (t *RestCookie) GetMany(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "GetMany")
	var cookieArg1 string
	{
		c, cookieErr := r.Cookie("arg1")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestCookie", "GetMany")
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

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestCookie", "GetMany")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			cookieArg2 = c.Value
		}
	}

	t.embed.GetMany(cookieArg1, cookieArg2)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "GetMany")
}

// ConvertToInt invoke Cookie.ConvertToInt using the request body as a json payload.
// ConvertToInt an arg from url query
func (t *RestCookie) ConvertToInt(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "ConvertToInt")
	var cookieArg1 int
	{
		c, cookieErr := r.Cookie("arg1")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestCookie", "ConvertToInt")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			{
				var err error
				cookieArg1, err = strconv.Atoi(c.Value)

				if err != nil {

					t.Log.Handle(w, r, err, "route", "error", "RestCookie", "ConvertToInt")
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

			}

		}
	}

	t.embed.ConvertToInt(cookieArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "ConvertToInt")
}

// ConvertToBool invoke Cookie.ConvertToBool using the request body as a json payload.
// ConvertToBool an arg from url query
func (t *RestCookie) ConvertToBool(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "ConvertToBool")
	var cookieArg1 bool
	{
		c, cookieErr := r.Cookie("arg1")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestCookie", "ConvertToBool")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			{
				var err error
				cookieArg1, err = strconv.ParseBool(c.Value)

				if err != nil {

					t.Log.Handle(w, r, err, "route", "error", "RestCookie", "ConvertToBool")
					http.Error(w, err.Error(), http.StatusInternalServerError)

					return
				}

			}

		}
	}

	t.embed.ConvertToBool(cookieArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "ConvertToBool")
}

// MaybeGet invoke Cookie.MaybeGet using the request body as a json payload.
// MaybeGet an arg if it exists in url query.
func (t *RestCookie) MaybeGet(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestCookie", "MaybeGet")
	var cookieArg1 *string
	{
		c, cookieErr := r.Cookie("arg1")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestCookie", "MaybeGet")
			http.Error(w, cookieErr.Error(), http.StatusInternalServerError)

			return
		}

		if c != nil {
			cookieArg1 = &c.Value
		}
	}

	t.embed.MaybeGet(cookieArg1)
	w.WriteHeader(200)

	t.Log.Handle(w, r, nil, "end", "RestCookie", "MaybeGet")
}

// RestCookieDescriptor describe a *RestCookie
type RestCookieDescriptor struct {
	ggt.TypeDescriptor
	about                *RestCookie
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

// NewRestCookieDescriptor describe a *RestCookie
func NewRestCookieDescriptor(about *RestCookie) *RestCookieDescriptor {
	ret := &RestCookieDescriptor{about: about}
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
func (t *RestCookieDescriptor) GetAll() *ggt.MethodDescriptor { return t.methodGetAll }

// GetAllRaw returns a MethodDescriptor
func (t *RestCookieDescriptor) GetAllRaw() *ggt.MethodDescriptor { return t.methodGetAllRaw }

// GetOne returns a MethodDescriptor
func (t *RestCookieDescriptor) GetOne() *ggt.MethodDescriptor { return t.methodGetOne }

// GetOneRaw returns a MethodDescriptor
func (t *RestCookieDescriptor) GetOneRaw() *ggt.MethodDescriptor { return t.methodGetOneRaw }

// MaybeGetOneRaw returns a MethodDescriptor
func (t *RestCookieDescriptor) MaybeGetOneRaw() *ggt.MethodDescriptor { return t.methodMaybeGetOneRaw }

// Write returns a MethodDescriptor
func (t *RestCookieDescriptor) Write() *ggt.MethodDescriptor { return t.methodWrite }

// MaybeDelete returns a MethodDescriptor
func (t *RestCookieDescriptor) MaybeDelete() *ggt.MethodDescriptor { return t.methodMaybeDelete }

// Delete returns a MethodDescriptor
func (t *RestCookieDescriptor) Delete() *ggt.MethodDescriptor { return t.methodDelete }

// GetMany returns a MethodDescriptor
func (t *RestCookieDescriptor) GetMany() *ggt.MethodDescriptor { return t.methodGetMany }

// ConvertToInt returns a MethodDescriptor
func (t *RestCookieDescriptor) ConvertToInt() *ggt.MethodDescriptor { return t.methodConvertToInt }

// ConvertToBool returns a MethodDescriptor
func (t *RestCookieDescriptor) ConvertToBool() *ggt.MethodDescriptor { return t.methodConvertToBool }

// MaybeGet returns a MethodDescriptor
func (t *RestCookieDescriptor) MaybeGet() *ggt.MethodDescriptor { return t.methodMaybeGet }
