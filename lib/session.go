package lib

import "net/http"

// SessionStoreProvider of session container.
type SessionStoreProvider interface {
	Get(*http.Request, string) (SessionSaver, error)
}

// SessionSaver saves container.
type SessionSaver interface {
	Save(*http.Request, http.ResponseWriter) error
	Get() (map[interface{}]interface{}, error)
}

// VoidSession does not store session.
type VoidSession struct{}

// Get does not returns a session store.
func (v *VoidSession) Get(r *http.Request, s string) (SessionSaver, error) { return &VoidSaver{}, nil }

// VoidSaver does not save session.
type VoidSaver struct{}

// Save does not save.
func (v *VoidSaver) Save(r *http.Request, w http.ResponseWriter) error { return nil }

// Get returns empty value.
func (v *VoidSaver) Get() (map[interface{}]interface{}, error) {
	return map[interface{}]interface{}{}, nil
}
