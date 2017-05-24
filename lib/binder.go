package lib

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"
)

type Method struct {
	Wrappers []Wrapper
	Getter   func(interface{}) http.HandlerFunc
	Name     string
	Route    string
	Methods  []string
}

func (m Method) Wrap(s http.HandlerFunc) http.HandlerFunc {
	for _, w := range m.Wrappers {
		s = w(s)
	}
	return s
}

type MethodSet struct {
	items []Method
}

type Wrapper func(http.HandlerFunc) http.HandlerFunc

func NewMethodSet() MethodSet { return MethodSet{} }

func (m MethodSet) Register(getter func(interface{}) http.HandlerFunc, name, route string, methods []string) MethodSet {
	m.items = append(m.items, Method{nil, getter, name, route, methods})
	return m
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type GorillaBinder struct {
	Set MethodSet
}

func NewGorillaBinder(set MethodSet) GorillaBinder {
	return GorillaBinder{set}
}

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

func (g GorillaBinder) Apply(router *mux.Router, instance interface{}) {
	for _, item := range g.Set.items {
		handler := item.Getter(instance)
		handler = item.Wrap(handler)
		router.HandleFunc("/"+item.Route, handler)
	}
}
