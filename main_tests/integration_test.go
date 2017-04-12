package main_test

import (
	"io/ioutil"
	"irrigation/home"
	"irrigation/irRelay"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestGarden(t *testing.T) {
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	call(home.HOST+"/control/garden?sate="+"unkn", t)
	callCmp(home.HOST+"/control/garden", irRelay.OFF, t)
	call(home.HOST+"/control/garden?state="+irRelay.ON, t)
	callCmp(home.HOST+"/control/garden", irRelay.ON, t)
}

func call(url string, t *testing.T) {
	res, err := http.Get(url)

	if err != nil {
		t.Log(err.Error())
		t.Fatal()
	}

	rb, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal()
	}

	t.Log(string(rb))
}

func callCmp(url string, cmp string, t *testing.T) {
	res, err := http.Get(url)

	if err != nil {
		t.Log(err.Error())
		t.Fatal()
	}

	rb, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal()
	}

	t.Log(string(rb))
	if string(rb) != cmp {
		t.Fail()
	}
}
