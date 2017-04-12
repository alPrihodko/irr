package main

import (
	"io"
	"irrigation/irr"
	"irrigation/wsHandler"
	"log"
	"net/http"
	"time"
)

func cState(w http.ResponseWriter, r *http.Request) {

	d, errs := currentState.ToJSON()
	if errs != nil {
		log.Println(errs.Error())
		http.Error(w, errs.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(d))
}

func appStateChanged() {
	log.Println("app state change triggered")

	currentState.GardenMode = ir01.GetMode()
	currentState.GardenState = ir01.GetState()

	currentState.FlowerBadMode = ir02.GetMode()
	currentState.FlowerBadState = ir02.GetState()

	currentState.FlowersMode = ir03.GetMode()
	currentState.FlowersState = ir03.GetState()

	currentState.Timestamp = int(time.Now().Unix())

	/*update UI*/
	d, errs := currentState.ToJSON()
	if errs != nil {
		irr.ReportAlert(errs.Error(), "Cannot report Temp to socket")
		return
	}

	err := wsHandler.ReportData(d)
	if err != nil {
		irr.ReportAlert(err.Error(), "Cannot report Temp to socket")
	}

	x := currentState

	historyData.Push(&x)
}
