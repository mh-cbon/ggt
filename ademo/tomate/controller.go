package tomate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/agext/levenshtein"
)

// Controller of tomatoes.
type Controller struct {
	backend TomatesContract
}

// NewController creates a new tomates controller
func NewController(backend TomatesContract) Controller {
	return Controller{backend}
}

// GetByID read the Tomate of given ID
//
// @route /read/{id:[0-9]+}
func (t Controller) GetByID(routeID string) (jsonResBody *Tomate, err error) {
	byRouteID := FilterTomates.ByID(routeID)
	t.backend.Transact(func(backend *Tomates) {
		jsonResBody = backend.Filter(byRouteID).First()
	})
	if jsonResBody == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, err
}

// SimilarColor returns colors similar to the given input
//
// @route /similar/color/{color}
func (t Controller) SimilarColor(routeColor string, getSensitive *bool) (jsonResBody *SimilarTomates, err error) {
	sensitive := getSensitive != nil && *getSensitive != false
	rVal := routeColor
	p := levenshtein.NewParams()
	if sensitive == false {
		rVal = strings.ToLower(rVal)
	}
	jsonResBody = NewSimilarTomates()
	t.backend.Transact(func(backend *Tomates) {
		backend.Map(func(t *Tomate) *Tomate {
			lVal := t.Color
			if sensitive == false {
				lVal = strings.ToLower(lVal)
			}
			res := levenshtein.Similarity(lVal, rVal, p)
			if res > 0.1 {
				jsonResBody.Push(&SimilarTomate{*t, res})

			}
			return t
		})
	})
	return jsonResBody, err
}

// Create a new Tomate
//
// @route /create
// @methods POST
func (t Controller) Create(postColor *string) (jsonResBody *Tomate, err error) {
	if postColor == nil {
		return nil, &UserInputError{errors.New("Missing color parameter")}
	}
	color := strings.TrimSpace(*postColor)
	if color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	byPostColor := FilterTomates.ByColor(color)
	t.backend.Transact(func(backend *Tomates) {
		if backend.Filter(byPostColor).NotEmpty() {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		jsonResBody = &Tomate{ID: fmt.Sprintf("%v", backend.Len()), Color: color}
		backend.Push(jsonResBody)
	})
	return jsonResBody, err
}

// Update an existing Tomate
//
// @route /write/{id:[0-9]+}
// @methods POST
func (t Controller) Update(routeID string, jsonReqBody *Tomate) (jsonResBody *Tomate, err error) {
	jsonReqBody.Color = strings.TrimSpace(jsonReqBody.Color)
	if jsonReqBody.Color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	byRouteID := FilterTomates.ByID(routeID)
	notRouteID := FilterTomates.NotID(routeID)
	byBodyColor := FilterTomates.ByColor(jsonReqBody.Color)
	updateColor := SetterTomates.SetColor(jsonReqBody.Color)
	t.backend.Transact(func(backend *Tomates) {
		byID := backend.Filter(byRouteID)
		if byID.Empty() {
			err = &NotFoundError{errors.New("ID does not exists")}
			return
		}
		if backend.Filter(byBodyColor, notRouteID).NotEmpty() {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		jsonResBody = byID.Map(updateColor).First()
	})
	if jsonResBody == nil && err == nil {
		err = &NotFoundError{errors.New("Tomate not found")}
	}
	return jsonResBody, err
}

// Remove an existing Tomate
//
// @route /remove/{id:[0-9]+}
// @methods POST
func (t Controller) Remove(ctx context.Context, routeID string) (jsonResBody bool, err error) {
	byRouteID := FilterTomates.ByID(routeID)
	t.backend.Transact(func(backend *Tomates) {
		byID := backend.Filter(byRouteID)
		if byID.Empty() {
			err = &NotFoundError{errors.New("ID does not exists")}
			return
		}
		jsonResBody = backend.Remove(byID.First())
	})
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
