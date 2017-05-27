package capable

import (
	"errors"
	"fmt"
	"net/http"
)

// Error ...
type Error struct{}

// TheMethod ...
func (c Error) TheMethod() (err error) {
	err = &WhateverError{errors.New("you decide")}
	fmt.Printf(`err %v
    `, err)
	return err
}

// Finalizer handle system error return
func (c Error) Finalizer(w http.ResponseWriter, r *http.Request, err error) {
	if _, ok := err.(*WhateverError); ok {
		http.Error(w, err.Error(), 666) // antechist in western world
	} else {
		http.Error(w, err.Error(), 444) // death in eastern world
	}
}

// WhateverError is your error type.
type WhateverError struct {
	error
}
