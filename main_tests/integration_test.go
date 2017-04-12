package main_test

import (
	"encoding/json"
	"fmt"
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
	call(home.HOST+"/control/garden?state="+"unkn", t)
	ret := callCmp(home.HOST+"/control/garden", irRelay.ON, t)
	if !ret {
		t.Fail()
	}
	call(home.HOST+"/control/garden?mode="+irRelay.ON, t)
	ret = callCmp(home.HOST+"/control/garden", irRelay.OFF, t)
	if !ret {
		t.Fail()
	}

}

func call(url string, t *testing.T) {
	_, err := http.Get(url)

	if err != nil {
		t.Log(err.Error())
		t.Fatal()
	}
}

func callCmp(url string, cmp string, t *testing.T) bool {
	res, err := http.Get(url)

	if err != nil {
		t.Log(err.Error())
		t.Fatal()
	}

	//rb, err := ioutil.ReadAll(res.Body)
	//res.Body.Close()
	//if err != nil {
	//	t.Fatal()
	//}

	bd := irRelay.Ir{}
	json.NewDecoder(res.Body).Decode(&bd)
	fmt.Println(bd)
	if bd.RelayMode != cmp {
		return false
	}
	return true
}
