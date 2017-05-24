package lib

import "net/http"

// Descriptor of a type
type Descriptor interface {
	Methods() []*MethodDescriptor
}

// TypeDescriptor is the commonalities implementation of a Descriptor.
type TypeDescriptor struct {
	Items []*MethodDescriptor
}

// Wrap every methods with w...
func (t *TypeDescriptor) Wrap(w ...Wrapper) {
	for _, m := range t.Items {
		m.Wrap(w...)
	}
}

// Register a new method on the descriptor
func (t *TypeDescriptor) Register(m *MethodDescriptor) {
	t.Items = append(t.Items, m)
}

// Methods returns the methods descriptor
func (t *TypeDescriptor) Methods() []*MethodDescriptor {
	return t.Items
}

// Wrapper wraps a controller method.
type Wrapper func(http.HandlerFunc) http.HandlerFunc

// MethodDescriptor describe a method and its property
type MethodDescriptor struct {
	Name     string
	Wrappers []Wrapper
	Handler  http.HandlerFunc
	Route    string
	Methods  []string
}

// Wrap the method handler with w...
func (m *MethodDescriptor) Wrap(w ...Wrapper) {
	m.Wrappers = append(m.Wrappers, w...)
}

// Wrapped return the final handler.
func (m *MethodDescriptor) Wrapped() http.HandlerFunc {
	ret := m.Handler
	for _, w := range m.Wrappers {
		ret = w(ret)
	}
	return ret
}
