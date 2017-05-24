package lib

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
)

// Method describe a controller method.
type Method struct {
	Wrappers []Wrapper
	Getter   func(interface{}) http.HandlerFunc
	Name     string
	Route    string
	Methods  []string
}

// Wrap the method controller.
func (m Method) Wrap(s http.HandlerFunc) http.HandlerFunc {
	for _, w := range m.Wrappers {
		s = w(s)
	}
	return s
}

// MethodSet gathers multiple methods.
type MethodSet struct {
	items []Method
}

// Wrapper wraps a controller method.
type Wrapper func(http.HandlerFunc) http.HandlerFunc

// NewMethodSet is a constructor.
func NewMethodSet() MethodSet { return MethodSet{} }

// Register a method to route on the transport.
func (m MethodSet) Register(getter func(interface{}) http.HandlerFunc, name, route string, methods []string) MethodSet {
	m.items = append(m.items, Method{nil, getter, name, route, methods})
	return m
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// GorillaBinder binds a method set to a girlla router.
type GorillaBinder struct {
	Set MethodSet
}

// NewGorillaBinder isa  constructor.
func NewGorillaBinder(set MethodSet) GorillaBinder {
	return GorillaBinder{set}
}

// Wrap given controller method with w.
func (g GorillaBinder) Wrap(x interface{}, w ...Wrapper) GorillaBinder {
	methodName := getFunctionName(x)
	for i, item := range g.Set.items {
		if item.Name == methodName {
			item.Wrappers = append(item.Wrappers, w...)
			g.Set.items[i] = item
		}
	}
	return g
}

// Apply given controller instance to a gorilla router.
func (g GorillaBinder) Apply(router *mux.Router, instance interface{}) {
	for _, item := range g.Set.items {
		handler := item.Getter(instance)
		handler = item.Wrap(handler)
		x := item.Route
		if !strings.HasPrefix(x, "/") {
			x = "/" + x
		}
		router.HandleFunc(x, handler)
	}
}
