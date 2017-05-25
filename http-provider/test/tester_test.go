package tester

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetValueParam(t *testing.T) {
	g := GetTest{T: t, WantValueParam: "10"}
	gg := NewRestGetTest(g)
	req, err := http.NewRequest("GET", "/ValueParam?param=10", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	gg.ValueParam(rr, req)
}

func TestGetValueParamMissing(t *testing.T) {
	g := GetTest{T: t, WantValueParam: ""}
	gg := NewRestGetTest(g)
	req, err := http.NewRequest("GET", "/ValueParam", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	gg.ValueParam(rr, req)
}

func TestGetRefParam(t *testing.T) {
	h := "10"
	g := GetTest{T: t, WantRefParam: &h}
	gg := NewRestGetTest(g)
	req, err := http.NewRequest("GET", "/RefParam?param=10", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	gg.RefParam(rr, req)
}

func TestGetRefParamMissing(t *testing.T) {
	g := GetTest{T: t, WantRefParam: nil}
	gg := NewRestGetTest(g)
	req, err := http.NewRequest("GET", "/RefParam", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	gg.RefParam(rr, req)
}
