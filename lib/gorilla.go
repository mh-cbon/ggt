package lib

import (
	"strings"

	"github.com/gorilla/mux"
)

// Gorilla binds an instance on given router
func Gorilla(instance Descriptor, router *mux.Router) {
	for _, m := range instance.Methods() {
		x := m.Route
		if !strings.HasPrefix(x, "/") {
			x = "/" + x
		}
		router.HandleFunc(x, m.Wrapped())
	}
}
