package tomate

// file generated by
// github.com/mh-cbon/ggt
// do not edit

import (
	"bytes"
	context "context"
	json "encoding/json"
	"errors"
	"net/http"
)

// RPCClient is an http-clienter of Controller.
// Controller of tomatoes.
type RPCClient struct {
	client *http.Client
}

// NewRPCClient constructs an http-clienter of Controller
func NewRPCClient(client *http.Client) *RPCClient {
	if client == nil {
		client = http.DefaultClient
	}
	ret := &RPCClient{
		client: client,
	}
	return ret
}

var xx02eba78324b686eb8959e353dd12d0255ed0c9b9 = http.StatusOK

// GetByID constructs a request to GetByID
func (t RPCClient) GetByID(routeID string) (*Tomate, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 string
		}{
			Arg0: routeID,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return nil, errors.New("todo")
		}

	}
	finalURL := "/GetByID"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return nil, errors.New("todo")
	}

	output := struct {
		Arg0 *Tomate
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}

// SimilarColor constructs a request to SimilarColor
func (t RPCClient) SimilarColor(routeColor string, getSensitive *bool) (*SimilarTomates, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 string
			Arg1 *bool
		}{
			Arg0: routeColor,
			Arg1: getSensitive,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return nil, errors.New("todo")
		}

	}
	finalURL := "/SimilarColor"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return nil, errors.New("todo")
	}

	output := struct {
		Arg0 *SimilarTomates
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}

// Create constructs a request to Create
func (t RPCClient) Create(postColor *string) (*Tomate, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 *string
		}{
			Arg0: postColor,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return nil, errors.New("todo")
		}

	}
	finalURL := "/Create"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return nil, errors.New("todo")
	}

	output := struct {
		Arg0 *Tomate
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}

// Update constructs a request to Update
func (t RPCClient) Update(routeID string, jsonReqBody *Tomate) (*Tomate, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 string
			Arg1 *Tomate
		}{
			Arg0: routeID,
			Arg1: jsonReqBody,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return nil, errors.New("todo")
		}

	}
	finalURL := "/Update"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return nil, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return nil, errors.New("todo")
	}

	output := struct {
		Arg0 *Tomate
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return nil, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}

// Remove constructs a request to Remove
func (t RPCClient) Remove(ctx context.Context, routeID string) (bool, error) {
	var reqBody bytes.Buffer

	{
		input := struct {
			Arg0 context.Context
			Arg1 string
		}{
			Arg0: ctx,
			Arg1: routeID,
		}
		encErr := json.NewEncoder(&reqBody).Encode(&input)
		if encErr != nil {
			return false, errors.New("todo")
		}

	}
	finalURL := "/Remove"
	req, reqErr := http.NewRequest("POST", finalURL, &reqBody)
	if reqErr != nil {
		return false, errors.New("todo")
	}

	res, resErr := t.client.Do(req)
	if resErr != nil {
		return false, errors.New("todo")
	}

	output := struct {
		Arg0 bool
		Arg1 error
	}{}
	{
		decErr := json.NewDecoder(res.Body).Decode(&output)
		if decErr != nil {
			return false, errors.New("todo")
		}

	}

	return output.Arg0, output.Arg1

}
