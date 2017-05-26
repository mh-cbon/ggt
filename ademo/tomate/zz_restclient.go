package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	"bytes"
	json "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RestClient is an http-clienter of Controller.
// Controller of tomatoes.
type RestClient struct {
	router *mux.Router
	Base   string
	Client *http.Client
}

// NewRestClient constructs an http-clienter of Controller
func NewRestClient(router *mux.Router) *RestClient {
	ret := &RestClient{
		router: router,
		Client: http.DefaultClient,
	}
	return ret
}

// GetByID constructs a request to /read/{id:[0-9]+}
func (t RestClient) GetByID(routeID string) (jsonResBody *Tomate, err error) {
	sReqURL := "/read/{id:[0-9]+}"
	sReqURL = strings.Replace(sReqURL, "{id:[0-9]+}", fmt.Sprintf("%v", routeID), 1)
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return jsonResBody, err
	}
	finalURL := reqURL.String()
	finalURL = fmt.Sprintf("%v%v", t.Base, finalURL)

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return jsonResBody, err
	}

	{
		res, resErr := t.Client.Do(req)
		if resErr != nil {
			return jsonResBody, err
		}

		decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
		if decErr != nil {
			return jsonResBody, err
		}

	}

	return jsonResBody, err
}

// Create constructs a request to /create
func (t RestClient) Create(postColor *string) (jsonResBody *Tomate, err error) {
	sReqURL := "/create"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return jsonResBody, err
	}
	form := url.Values{}
	form.Add("color", *postColor)
	finalURL := reqURL.String()
	finalURL = fmt.Sprintf("%v%v", t.Base, finalURL)

	req, reqErr := http.NewRequest("GET", finalURL, strings.NewReader(form.Encode()))
	if reqErr != nil {
		return jsonResBody, err
	}

	{
		res, resErr := t.Client.Do(req)
		if resErr != nil {
			return jsonResBody, err
		}

		decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
		if decErr != nil {
			return jsonResBody, err
		}

	}

	return jsonResBody, err
}

// Update constructs a request to /write/{id:[0-9]+}
func (t RestClient) Update(routeID string, jsonReqBody *Tomate) (jsonResBody *Tomate, err error) {

	var body io.ReadWriter
	{
		var b bytes.Buffer
		body = &b
		encErr := json.NewEncoder(body).Encode(jsonReqBody)
		if encErr != nil {
			return jsonResBody, err
		}

	}
	sReqURL := "/write/{id:[0-9]+}"
	sReqURL = strings.Replace(sReqURL, "{id:[0-9]+}", fmt.Sprintf("%v", routeID), 1)
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return jsonResBody, err
	}
	finalURL := reqURL.String()
	finalURL = fmt.Sprintf("%v%v", t.Base, finalURL)

	req, reqErr := http.NewRequest("GET", finalURL, body)
	if reqErr != nil {
		return jsonResBody, err
	}

	{
		res, resErr := t.Client.Do(req)
		if resErr != nil {
			return jsonResBody, err
		}

		decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
		if decErr != nil {
			return jsonResBody, err
		}

	}

	return jsonResBody, err
}

// Remove constructs a request to /remove/{id:[0-9]+}
func (t RestClient) Remove(routeID string) (jsonResBody bool, err error) {
	sReqURL := "/remove/{id:[0-9]+}"
	sReqURL = strings.Replace(sReqURL, "{id:[0-9]+}", fmt.Sprintf("%v", routeID), 1)
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return jsonResBody, err
	}
	finalURL := reqURL.String()
	finalURL = fmt.Sprintf("%v%v", t.Base, finalURL)

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return jsonResBody, err
	}

	{
		res, resErr := t.Client.Do(req)
		if resErr != nil {
			return jsonResBody, err
		}

		decErr := json.NewDecoder(res.Body).Decode(jsonResBody)
		if decErr != nil {
			return jsonResBody, err
		}

	}

	return jsonResBody, err
}