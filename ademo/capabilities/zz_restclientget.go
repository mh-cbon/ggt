package capable

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// RestClientGet is an http-clienter of Get.
// Get ...
type RestClientGet struct {
	router *mux.Router
	client *http.Client
}

// NewRestClientGet constructs an http-clienter of Get
func NewRestClientGet(router *mux.Router, client *http.Client) *RestClientGet {
	if client == nil {
		client = http.DefaultClient
	}
	ret := &RestClientGet{
		router: router,
		client: client,
	}
	return ret
}

var xx0324d1f07eaec31c9f34207f0bfc8d9aa69e69cf = bytes.MinRead
var xx4da51200447f6306d6b0d6539a4eb45ec2bc971f = fmt.Println
var xx225c65374f4b2a9f6f2ef056219c3bdb3c932dc5 = url.PathEscape
var xx19c36ba5a408b3b3e4901c6dbc58a948e67028cf = strings.ToUpper
var xx3bde715e679ca5a1af594823df08339cd10fadf4 = context.Canceled
var xxef66dc4d974896e9db8bbe408596bc595b77fd5e = mux.Vars
var xx5cff34fbc483de2d38b5dc90e9c358cb5d66cb72 = io.Copy
var xx290bdb6a7e347a56f5865dba31b5788bccbad206 = http.StatusOK

// GetAll constructs a request to getall
func (t RestClientGet) GetAll(getValues map[string][]string) {
	sReqURL := "getall"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return
	}

	for k, vv := range getValues {
		for _, v := range vv {
			reqURL.Query().Add(k, v)
		}
	}

	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return
	}
	_, resErr := t.client.Do(req)
	if resErr != nil {
		return
	}

	return
}

// GetAll2 constructs a request to getall2
func (t RestClientGet) GetAll2(getValues map[string]string) {
	sReqURL := "getall2"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return
	}

	for k, v := range getValues {
		reqURL.Query().Add(k, v)
	}

	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return
	}
	_, resErr := t.client.Do(req)
	if resErr != nil {
		return
	}

	return
}

// GetOne constructs a request to getone
func (t RestClientGet) GetOne(getArg1 string) {
	sReqURL := "getone"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return
	}
	var xxgetArg1 string
	xxgetArg1 = getArg1
	reqURL.Query().Add("arg1", xxgetArg1)

	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return
	}
	_, resErr := t.client.Do(req)
	if resErr != nil {
		return
	}

	return
}

// GetMany constructs a request to getmany
func (t RestClientGet) GetMany(getArg1 string, getArg2 string) {
	sReqURL := "getmany"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return
	}
	var xxgetArg1 string
	xxgetArg1 = getArg1
	reqURL.Query().Add("arg1", xxgetArg1)
	var xxgetArg2 string
	xxgetArg2 = getArg2
	reqURL.Query().Add("arg2", xxgetArg2)

	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return
	}
	_, resErr := t.client.Do(req)
	if resErr != nil {
		return
	}

	return
}

// ConvertToInt constructs a request to converttoint
func (t RestClientGet) ConvertToInt(getArg1 int) {
	sReqURL := "converttoint"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return
	}
	var xxgetArg1 string
	xxgetArg1 = fmt.Sprintf("%v", getArg1)

	reqURL.Query().Add("arg1", xxgetArg1)

	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return
	}
	_, resErr := t.client.Do(req)
	if resErr != nil {
		return
	}

	return
}

// ConvertToBool constructs a request to converttobool
func (t RestClientGet) ConvertToBool(getArg1 bool) {
	sReqURL := "converttobool"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return
	}
	var xxgetArg1 string
	xxgetArg1 = "false"
	if getArg1 {
		xxgetArg1 = "true"
	}

	reqURL.Query().Add("arg1", xxgetArg1)

	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return
	}
	_, resErr := t.client.Do(req)
	if resErr != nil {
		return
	}

	return
}

// ConvertToSlice constructs a request to converttoslice
func (t RestClientGet) ConvertToSlice(getArg1 []bool) {
	sReqURL := "converttoslice"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return
	}

	for _, itemgetArg1 := range getArg1 {
		var xxgetArg1 string
		xxgetArg1 = "false"
		if itemgetArg1 {
			xxgetArg1 = "true"
		}

		reqURL.Query().Add("arg1", xxgetArg1)
	}

	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return
	}
	_, resErr := t.client.Do(req)
	if resErr != nil {
		return
	}

	return
}

// MaybeGet constructs a request to maybeget
func (t RestClientGet) MaybeGet(getArg1 *string) {
	sReqURL := "maybeget"
	reqURL, URLerr := url.ParseRequestURI(sReqURL)
	if URLerr != nil {
		return
	}

	if getArg1 != nil {
		var xxgetArg1 string
		xxgetArg1 = *getArg1
		reqURL.Query().Add("arg1", xxgetArg1)
	}

	finalURL := reqURL.String()

	req, reqErr := http.NewRequest("GET", finalURL, nil)
	if reqErr != nil {
		return
	}
	_, resErr := t.client.Do(req)
	if resErr != nil {
		return
	}

	return
}
