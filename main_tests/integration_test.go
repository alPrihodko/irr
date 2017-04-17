package main_test

import (
	"encoding/json"
	"fmt"
	"irrigation/irRelay"
	"irrigation/irr"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestGardenOff(t *testing.T) {
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	call(irr.HOST+"/control/garden?mode="+irRelay.ON, t)
	call(irr.HOST+"/control/garden?mode="+irRelay.OFF, t)
	callCmp(irr.HOST+"/control/garden", irRelay.OFF, t)
}

func TestGardenOn(t *testing.T) {
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	call(irr.HOST+"/control/garden?mode="+irRelay.OFF, t)
	call(irr.HOST+"/control/garden?mode="+irRelay.ON, t)
	callCmp(irr.HOST+"/control/garden", irRelay.ON, t)
}

func call(url string, t *testing.T) {
	_, err := http.Get(url)

	if err != nil {
		t.Log(err.Error())
		t.Fatal(err.Error())
	}
}

func callCmp(url string, cmp string, t *testing.T) bool {
	res, err := http.Get(url)

	if err != nil {
		t.Log(err.Error())
		t.Fatal(err.Error())
	}

	bd := irRelay.Ir{}
	json.NewDecoder(res.Body).Decode(&bd)
	fmt.Println(bd)
	if bd.RelayMode != cmp {
		t.Error("Expected: ", cmp, ", got:", bd.RelayMode)
		return false
	}
	return true
}
