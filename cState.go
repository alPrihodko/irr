package main

import (
	"irrigation/irr"
	"irrigation/wsHandler"
	"log"
	"time"
)

func appStateChanged() {
	log.Println("app state change triggered")

	irr.CurrentState.GardenName = ir01.Relay.Name()
	log.Println("registering mode: " + ir01.GetMode())
	irr.CurrentState.GardenMode = ir01.GetMode()
	irr.CurrentState.GardenState = ir01.GetState()

	irr.CurrentState.FlowerBadMode = ir02.GetMode()
	irr.CurrentState.FlowerBadState = ir02.GetState()

	irr.CurrentState.GrapesMode = ir03.GetMode()
	irr.CurrentState.GrapesState = ir03.GetState()

	irr.CurrentState.Timestamp = int(time.Now().Unix())

	/*update UI*/
	d, errs := irr.CurrentState.ToJSON()
	if errs != nil {
		irr.ReportAlert(errs.Error(), "Cannot report current state")
		return
	}

	err := wsHandler.ReportData(d)
	if err != nil {
		irr.ReportAlert(err.Error(), "Cannot report relay state")
	}

	x := irr.CurrentState

	historyData.Push(&x)
}
