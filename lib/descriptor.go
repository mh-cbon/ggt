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

// WrapMethod every methods with w...
func (t *TypeDescriptor) WrapMethod(w ...MethodWrapper) {
	for _, m := range t.Items {
		m.WrapMethod(w...)
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

// MethodWrapper wraps a controller method.
type MethodWrapper func(*MethodDescriptor) Wrapper

// MethodDescriptor describe a method and its property
type MethodDescriptor struct {
	Name     string
	Wrappers []interface{}
	Handler  http.HandlerFunc
	Route    string
	Methods  []string
}

// Wrap the method handler with w...
func (m *MethodDescriptor) Wrap(w ...Wrapper) {
	for _, e := range w {
		m.Wrappers = append(m.Wrappers, e)
	}
}

// WrapMethod the method handler with w...
func (m *MethodDescriptor) WrapMethod(w ...MethodWrapper) {
	for _, e := range w {
		m.Wrappers = append(m.Wrappers, e)
	}
}

// Wrapped return the final handler.
func (m *MethodDescriptor) Wrapped() http.HandlerFunc {
	ret := m.Handler
	for _, w := range m.Wrappers {
		if x, ok := w.(Wrapper); ok {
			ret = x(ret)

		} else if y, ok := w.(MethodWrapper); ok {
			ret = y(m)(ret)
		}
	}
	return ret
}
