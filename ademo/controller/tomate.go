package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mh-cbon/ggt/ademo/model"
	"github.com/mh-cbon/ggt/ademo/slicegen"
)

type Tomates struct {
	backend slicegen.TomatesContract
}

func NewTomates(backend slicegen.TomatesContract) Tomates {
	return Tomates{backend}
}

func (t Tomates) GetById(getID string) (jsonResBody model.Tomate, err error) {
	err = &NotFoundError{errors.New("Tomate not found")}
	t.backend.Map(func(x model.Tomate) model.Tomate {
		if x.ID == getID {
			jsonResBody = x
			err = nil
		}
		return x
	})
	return jsonResBody, err
}

func (t Tomates) Create(postColor string) (jsonResBody *model.Tomate, err error) {
	postColor = strings.TrimSpace(postColor)
	if postColor == "" {
		return nil, &UserInputError{errors.New("color must not be empty")}
	}
	tomate := model.Tomate{Color: postColor}
	t.backend.Transact(func(backend *slicegen.Tomates) {
		exist := backend.Filter(slicegen.FilterTomates.ByColor(postColor)).Len()
		if exist > 0 {
			err = &UserInputError{errors.New("color must be unique")}
			return
		}
		tomate.ID = fmt.Sprintf("%v", backend.Len())
		backend.Push(model.Tomate{ID: fmt.Sprintf("%v", backend.Len()), Color: postColor})
		jsonResBody = &tomate
	})
	return jsonResBody, err
}

type NotFoundError struct {
	error
}

type UserInputError struct {
	error
}

func (t Tomates) Finalizer(w http.ResponseWriter, r *http.Request, err error) {
	if _, ok := err.(*UserInputError); ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else if _, ok := err.(*NotFoundError); ok {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
