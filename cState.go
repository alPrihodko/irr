package main

import (
	"irrigation/irr"
	"irrigation/wsHandler"
	"log"
	"time"
)

func appStateChanged() {
	log.Println("app state change triggered")

	irr.CurrentState.GardenMode = ir01.GetMode()
	irr.CurrentState.GardenState = ir01.GetState()

	irr.CurrentState.FlowerBadMode = ir02.GetMode()
	irr.CurrentState.FlowerBadState = ir02.GetState()

	irr.CurrentState.FlowersMode = ir03.GetMode()
	irr.CurrentState.FlowersState = ir03.GetState()

	irr.CurrentState.Timestamp = int(time.Now().Unix())

	/*update UI*/
	d, errs := irr.CurrentState.ToJSON()
	if errs != nil {
		irr.ReportAlert(errs.Error(), "Cannot report Temp to socket")
		return
	}

	err := wsHandler.ReportData(d)
	if err != nil {
		irr.ReportAlert(err.Error(), "Cannot report Temp to socket")
	}

	x := irr.CurrentState

	historyData.Push(&x)
}
