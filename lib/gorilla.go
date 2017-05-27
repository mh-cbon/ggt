package lib

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Gorilla binds an instance on given router
func Gorilla(instance Descriptor, router *mux.Router) {
	for _, m := range instance.Methods() {
		x := m.Route
		if !strings.HasPrefix(x, "/") {
			x = "/" + x
		}
		route := router.HandleFunc(x, m.Wrapped())
		if len(m.Methods) > 0 {
			route.Methods(m.Methods...)
		}
	}
}

// GorrillaSessionStoreProvider provides gorilla session store.
type GorrillaSessionStoreProvider struct {
	Provider sessions.Store
}

// Get returns a session store.
func (v *GorrillaSessionStoreProvider) Get(r *http.Request, s string) (SessionSaver, error) {
	x, err := v.Provider.Get(r, s)
	return &GorrillaSessionSaver{x}, err
}

// GorrillaSessionSaver saves a gorilla store.
type GorrillaSessionSaver struct {
	Store *sessions.Session
}

// Save a store.
func (v *GorrillaSessionSaver) Save(r *http.Request, w http.ResponseWriter) error {
	return v.Store.Save(r, w)
}

// Get returns empty value.
func (v *GorrillaSessionSaver) Get() (map[interface{}]interface{}, error) {
	return v.Store.Values, nil
}
