package login

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	json "encoding/json"
	ggt "github.com/mh-cbon/ggt/lib"
	"io"
	"net/http"
	"strconv"
	"time"
)

var xxc52d52149b9b88887c5b162b670197217023167c = strconv.Atoi
var xxdddd5c3cff35fd51b960096f826c90c635b524c8 = io.Copy
var xx50b52bbd0311cf9ec8e406218f58ce92f00e2679 = http.StatusOK

// RestController is an httper of Controller.
// Controller of user login
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

// Login invoke Controller.Login using the request body as a json payload.
// Login user by its login/password
//
// @route /login
// @methods POST
func (t *RestController) Login(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestController", "Login")

	{
		err := r.ParseForm()

		if err != nil {

			t.Log.Handle(w, r, err, "parseform", "error", "RestController", "Login")
			t.embed.Finalizer(w, r, err)

			return
		}

	}
	var postLogin string
	if _, ok := r.Form["login"]; ok {
		xxTmppostLogin := r.FormValue("login")
		t.Log.Handle(w, r, nil, "input", "form", "login", xxTmppostLogin, "RestController", "Login")
		postLogin = xxTmppostLogin
	}
	var postPassword string
	if _, ok := r.Form["password"]; ok {
		xxTmppostPassword := r.FormValue("password")
		t.Log.Handle(w, r, nil, "input", "form", "password", xxTmppostPassword, "RestController", "Login")
		postPassword = xxTmppostPassword
	}

	jsonResBody, login, err := t.embed.Login(postLogin, postPassword)

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RestController", "Login")
		t.embed.Finalizer(w, r, err)

		return
	}

	{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		encErr := json.NewEncoder(w).Encode(jsonResBody)

		if encErr != nil {

			t.Log.Handle(w, r, encErr, "res", "json", "encode", "error", "RestController", "Login")
			t.embed.Finalizer(w, r, encErr)

			return
		}

	}

	if login == nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "login",
			Expires: time.Now().Add(-time.Hour * 24 * 100),
		})
	} else {
		http.SetCookie(w, login)
	}

	t.Log.Handle(w, r, nil, "end", "RestController", "Login")
}

// Logout invoke Controller.Logout using the request body as a json payload.
// Logout user by its login/password
//
// @route /logout
// @methods POST
func (t *RestController) Logout(w http.ResponseWriter, r *http.Request) {
	t.Log.Handle(w, r, nil, "begin", "RestController", "Logout")
	var cookieLogin string
	{
		c, cookieErr := r.Cookie("login")

		if cookieErr != nil {

			t.Log.Handle(w, r, cookieErr, "req", "cookie", "error", "error", "RestController", "Logout")
			t.embed.Finalizer(w, r, cookieErr)

			return
		}

		cookieLogin = c.Value
	}

	login, err := t.embed.Logout(cookieLogin)

	if err != nil {

		t.Log.Handle(w, r, err, "business", "error", "RestController", "Logout")
		t.embed.Finalizer(w, r, err)

		return
	}

	if login == nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "login",
			Expires: time.Now().Add(-time.Hour * 24 * 100),
		})
	} else {
		http.SetCookie(w, login)
	}

	t.Log.Handle(w, r, nil, "end", "RestController", "Logout")
}

// Create invoke Controller.Create using the request body as a json payload.
// Create user by its login/password
//
// @route /create
// @methods POST
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
	var postLogin string
	if _, ok := r.Form["login"]; ok {
		xxTmppostLogin := r.FormValue("login")
		t.Log.Handle(w, r, nil, "input", "form", "login", xxTmppostLogin, "RestController", "Create")
		postLogin = xxTmppostLogin
	}
	var postPassword string
	if _, ok := r.Form["password"]; ok {
		xxTmppostPassword := r.FormValue("password")
		t.Log.Handle(w, r, nil, "input", "form", "password", xxTmppostPassword, "RestController", "Create")
		postPassword = xxTmppostPassword
	}

	jsonResBody, err := t.embed.Create(postLogin, postPassword)

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

// RestControllerDescriptor describe a *RestController
type RestControllerDescriptor struct {
	ggt.TypeDescriptor
	about        *RestController
	methodLogin  *ggt.MethodDescriptor
	methodLogout *ggt.MethodDescriptor
	methodCreate *ggt.MethodDescriptor
}

// NewRestControllerDescriptor describe a *RestController
func NewRestControllerDescriptor(about *RestController) *RestControllerDescriptor {
	ret := &RestControllerDescriptor{about: about}
	ret.methodLogin = &ggt.MethodDescriptor{
		Name:    "Login",
		Handler: about.Login,
		Route:   "/login",
		Methods: []string{"POST"},
	}
	ret.TypeDescriptor.Register(ret.methodLogin)
	ret.methodLogout = &ggt.MethodDescriptor{
		Name:    "Logout",
		Handler: about.Logout,
		Route:   "/logout",
		Methods: []string{"POST"},
	}
	ret.TypeDescriptor.Register(ret.methodLogout)
	ret.methodCreate = &ggt.MethodDescriptor{
		Name:    "Create",
		Handler: about.Create,
		Route:   "/create",
		Methods: []string{"POST"},
	}
	ret.TypeDescriptor.Register(ret.methodCreate)
	return ret
}

// Login returns a MethodDescriptor
func (t *RestControllerDescriptor) Login() *ggt.MethodDescriptor { return t.methodLogin }

// Logout returns a MethodDescriptor
func (t *RestControllerDescriptor) Logout() *ggt.MethodDescriptor { return t.methodLogout }

// Create returns a MethodDescriptor
func (t *RestControllerDescriptor) Create() *ggt.MethodDescriptor { return t.methodCreate }
