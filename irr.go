package main

import (
	"flag"
	"irrigation/irRelay"
	"irrigation/irr"
	"irrigation/wsHandler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	//_ "github.com/icattlecoder/godaemon"

	"github.com/hybridgroup/gobot"

	"golang.org/x/net/websocket"
)

const configFileName = "/etc/irrigation.conf"

/*
HISTORYDATASERIAL file which contains history data for my home
*/
const HISTORYDATASERIAL = "goIrrigationData.b64"

/*
INTERVAL  Check sensors status with interval
*/
var INTERVAL int

//var err error

var conf Config

//var conns socketConns

//relays
var ir01 irRelay.Ir
var ir02 irRelay.Ir
var ir03 irRelay.Ir

var wh wsHandler.WsHandler

//var currentState irr.HData
var historyData irr.HistoryData

func main() {

	err := conf.loadConfig()
	if err != nil {
		log.Println("Likely use default configuration")
	}

	gbot := gobot.NewGobot()

	//conns = socketConns{make(map[int32]*websocket.Conn), &sync.Mutex{}}

	ir01 = irRelay.New("garden", "33", &wh, appStateChanged)
	http.HandleFunc("/control/"+ir01.Relay.Name(), ir01.RelayHandler)
	ir02 = irRelay.New("flowerbad", "35", &wh, appStateChanged)
	http.HandleFunc("/control/"+ir02.Relay.Name(), ir02.RelayHandler)
	ir03 = irRelay.New("grapes", "31", &wh, appStateChanged)
	http.HandleFunc("/control/"+ir03.Relay.Name(), ir03.RelayHandler)

	//ir03 = irRelay.New("grapes", "29", &wh, appStateChanged)
	//http.HandleFunc("/control/"+ir03.Relay.Name(), ir03.RelayHandler)

	//currentState = irr.HData{}
	//currentState.Index = 2
	//log.Println(currentState.Index)
	flag.IntVar(&INTERVAL, "timeout", 60, "Timeout?")
	flag.Parse()

	log.Println("Timeout interval to track sensors: ", INTERVAL)
	historyData.RestoreFromFile(HISTORYDATASERIAL)
	irr.CurrentState = historyData.Last()
	http.Handle("/relays", websocket.Handler(relHandler))

	http.Handle("/", http.FileServer(http.Dir("ui")))
	http.HandleFunc("/control/currentState", irr.CurrentStateHandler)
	http.HandleFunc("/control/hdata", historyData.HistoryDataHandler)
	http.HandleFunc("/control/config", configHandler)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGABRT)

	go func() {
		<-c
		log.Println("Save history data...")
		historyData.SerializeToFile(HISTORYDATASERIAL)
		irRelay.Stop()
		os.Exit(1)
	}()

	go gbot.Start()

	stop := scheduleBackup(backupHistoryData, time.Duration(INTERVAL*60)*time.Second, &historyData, HISTORYDATASERIAL)

	err = http.ListenAndServe(":1235", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	if stop != nil {
		stop <- true
	}

}

func scheduleBackup(what func(*irr.HistoryData, string), delay time.Duration,
	q *irr.HistoryData, l string) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			what(q, l)
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func backupHistoryData(q *irr.HistoryData, local string) {
	historyData.SerializeToFile(local)

	if _, err := DB.UploadFile(local, "/backup/irrigation.b64", true, ""); err != nil {
		log.Printf("Error uploading %s: %s\n", local, err)
	} else {
		log.Printf("File %s successfully uploaded\n", local)
	}
}
