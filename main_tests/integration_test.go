package main_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var host = "http://192.168.1.46:1235"

func TestInit(t *testing.T) {
	resp := httptest.NewRecorder()

	uri := "r01?"
	unlno := "Off"

	param := make(url.Values)
	param["state"] = []string{unlno}
	t.Log("start")
	req, err := http.NewRequest("GET", host+"/"+uri+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		t.Log(string(p))
		//if strings.Contains(string(p), "Error") {
		//        t.Errorf("header response shouldn't return error: %s", p)
		//} else if !strings.Contains(string(p), `expected result`) {
		//        t.Errorf("header response doen't match:\n%s", p)
		//}
	}
}
