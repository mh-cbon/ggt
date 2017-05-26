package login

import (
	"errors"
	"net/http"
	"time"
)

// Controller of user login
type Controller struct {
	backend HashedUsersContract
}

// NewController creates a new tomates controller
func NewController(backend HashedUsersContract) Controller {
	return Controller{backend}
}

// Login user by its login/password
//
// @route /login
// @methods POST
func (t Controller) Login(postLogin, postPassword string) (jsonResBody *User, login *http.Cookie, err error) {
	h := Hash(postLogin + postPassword + "s")
	t.backend.Transact(func(backend *HashedUsers) {
		res := backend.
			Filter(FilterHashedUsers.ByHash(h)).
			Map(SetterHashedUsers.SetLastLogin(time.Now())).
			First()
		if res != nil {
			jsonResBody = &res.User
			login = &http.Cookie{Value: postLogin, Expires: time.Now().Add(time.Hour * 24 * 365)}
		}
	})
	if jsonResBody == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, login, err
}

// Logout user by its login/password
//
// @route /logout
// @methods POST
func (t Controller) Logout(cookieLogin string) (login *http.Cookie, err error) {
	if cookieLogin != "" {
		t.backend.Transact(func(backend *HashedUsers) {
			backend.
				Filter(FilterHashedUsers.ByLogin(cookieLogin)).
				Map(SetterHashedUsers.SetLastLogout(time.Now()))
		})
	}
	return nil, nil
}

// Create user by its login/password
//
// @route /create
// @methods POST
func (t Controller) Create(postLogin, postPassword string) (jsonResBody *User, err error) {
	t.backend.Transact(func(backend *HashedUsers) {
		if !backend.Filter(FilterHashedUsers.ByLogin(postLogin)).Empty() {
			err = &NotFoundError{errors.New("Login not found")}
			return
		}
		res := &HashedUser{
			User: User{postLogin, postPassword},
			Hash: Hash(postLogin + postPassword + "s"),
		}
		if !backend.Filter(FilterHashedUsers.ByHash(res.Hash)).Empty() {
			err = &NotFoundError{errors.New("Collision detected")}
			return
		}
		backend.Push(res)
		jsonResBody = &res.User
	})
	if jsonResBody == nil && err == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, err
}

// NotFoundError is an error thrown when a value is not found
type NotFoundError struct {
	error
}

// UserInputError is an error thrown when the user input is incomplete or invalid
type UserInputError struct {
	error
}

// Finalizer behave appropriately by error types
func (t Controller) Finalizer(w http.ResponseWriter, r *http.Request, err error) {
	if _, ok := err.(*UserInputError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if _, ok := err.(*NotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
