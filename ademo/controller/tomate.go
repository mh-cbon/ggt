package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mh-cbon/ggt/ademo/model"
	"github.com/mh-cbon/ggt/ademo/slicegen"
)

// Tomates controller.
type Tomates struct {
	backend slicegen.TomatesContract
}

// NewTomates creates a new tomates controller
func NewTomates(backend slicegen.TomatesContract) Tomates {
	return Tomates{backend}
}

// GetByID read the Tomate of given ID
func (t Tomates) GetByID(getID string) (jsonResBody *model.Tomate, err error) {
	err = &NotFoundError{errors.New("Tomate not found")}
	t.backend.Map(func(x *model.Tomate) *model.Tomate {
		if x.ID == getID {
			jsonResBody = x
			err = nil
		}
		return x
	})
	return jsonResBody, err
}

// Create a new Tomate
func (t Tomates) Create(postColor *string) (jsonResBody *model.Tomate, err error) {
	if postColor == nil {
		return nil, &UserInputError{errors.New("Missing color parameter")}
	}
	color := strings.TrimSpace(*postColor)
	if color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	t.backend.Transact(func(backend *slicegen.Tomates) {
		exist := backend.Filter(slicegen.FilterTomates.ByColor(color)).Len()
		if exist > 0 {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		jsonResBody = &model.Tomate{ID: fmt.Sprintf("%v", backend.Len()), Color: color}
		backend.Push(jsonResBody)
	})
	return jsonResBody, err
}

// Update an existing Tomate
//
// @route /write/{id:[0-9]+}
func (t Tomates) Update(routeID string, jsonReqBody *model.Tomate) (jsonResBody *model.Tomate, err error) {
	jsonReqBody.Color = strings.TrimSpace(jsonReqBody.Color)
	if jsonReqBody.Color == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	t.backend.Transact(func(backend *slicegen.Tomates) {
		byID := backend.Filter(slicegen.FilterTomates.ByID(routeID))
		if byID.Len() == 0 {
			err = &NotFoundError{errors.New("ID does not exists")}
			return
		}
		byColor := backend.Filter(slicegen.FilterTomates.ByColor(jsonReqBody.Color))
		if byColor.Len() > 0 && byID.First().ID != byColor.First().ID {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		backend.Map(func(x *model.Tomate) *model.Tomate {
			if x.ID == routeID {
				x.Color = jsonReqBody.Color
			}
			return x
		})
	})
	jsonResBody = t.backend.Filter(slicegen.FilterTomates.ByID(routeID)).First()
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
func (t Tomates) Finalizer(w http.ResponseWriter, r *http.Request, err error) {
	if _, ok := err.(*UserInputError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if _, ok := err.(*NotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
