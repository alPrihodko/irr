package main_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
)

var host = "http://sasha123.ddns.ukrtel.net:1235"

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestR01(t *testing.T) {
	resp := httptest.NewRecorder()

	uri := "r01?"
	unlno := "Off"
	param := make(url.Values)
	param["state"] = []string{unlno}
	req, err := http.NewRequest("GET", host+"/"+uri+param.Encode(), nil)
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		t.Log(string(p))
		if resp.Code != 200 {
			t.Error("Error code: " + strconv.Itoa(resp.Code))
			t.Fail()
		}
		//if strings.Contains(string(p), "Error") {
		//        t.Errorf("header response shouldn't return error: %s", p)
		//} else if !strings.Contains(string(p), `expected result`) {
		//        t.Errorf("header response doen't match:\n%s", p)
		//}
	}
}
