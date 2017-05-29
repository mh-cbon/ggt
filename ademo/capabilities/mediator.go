package capable

import "net/http"

// Mediator so called middleware
type Mediator struct{}

// AddSomeHeader ...
func (c Mediator) AddSomeHeader(handler http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	handler(w, r)
}

// WriteResponse ...
func (c Mediator) WriteResponse() (jsonResBody *string) {
	res := ""
	jsonResBody = &res
	return
}
