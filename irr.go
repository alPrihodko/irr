package main

import (
	"flag"
	"irrigation/irRelay"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	//_ "github.com/icattlecoder/godaemon"

	"github.com/hybridgroup/gobot"

	"golang.org/x/net/websocket"
)

const configFileName = "/etc/irrigation.conf"

/*
HISTORYDATASERIAL file which contains history data for my home
*/
const HISTORYDATASERIAL = "/home/pi/goIrrigationData.b64"

/*
INTERVAL  Check sensors status with interval
*/
var INTERVAL int

//var err error

var conf Config

type socketConns struct {
	ws   map[int32]*websocket.Conn
	lock *sync.Mutex
}

var conns socketConns
var rconns socketConns

//relays
var ir01 irRelay.Ir
var ir02 irRelay.Ir
var ir03 irRelay.Ir

//var currentState home.HData
//var historyData home.HistoryData

func main() {

	err := conf.loadConfig()
	if err != nil {
		log.Println("Likely use default configuration")
	}

	gbot := gobot.NewGobot()

	conns = socketConns{make(map[int32]*websocket.Conn), &sync.Mutex{}}
	rconns = socketConns{make(map[int32]*websocket.Conn), &sync.Mutex{}}

	ir01 = irRelay.New("garden", "19")
	ir02 = irRelay.New("flowerbad", "21")
	ir03 = irRelay.New("flowers", "23")

	//currentState = home.HData{}
	//currentState.Index = 2
	//log.Println(currentState.Index)
	flag.IntVar(&INTERVAL, "timeout", 60, "Timeout?")
	flag.Parse()

	log.Println("Timeout interval to track sensors: ", INTERVAL)
	//stop = scheduleT(reportFloat, 10*time.Second, "temp1", 10)
	//historyData.RestoreFromFile(HISTORYDATASERIAL)

	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/relays", websocket.Handler(relHandler))

	http.Handle("/", http.FileServer(http.Dir("/home/pi/w/go/src/irrigation")))
	//http.HandleFunc("/control/currentState", cState)
	http.HandleFunc("/control/config", configHandler)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, os.Kill)
	signal.Notify(c, syscall.SIGABRT)

	go func() {
		<-c
		log.Println("Save history data...")
		//historyData.SerializeToFile(HISTORYDATASERIAL)
		irRelay.Stop()
		os.Exit(1)
	}()

	//stop = schedule(reportSensors, time.Duration(INTERVAL)*time.Second, sensors)
	//      led.Toggle()

	/*
		work := func() {
			//defer home.Stop()
			gobot.Every(time.Duration(INTERVAL)*time.Second, func() {
				//log.Println("gobot heartbeat")

				//if SENSORS {
				//	reportSensors(sensors)
				//}

			})

		}
	*/
	//robot := gobot.NewRobot("blinkBot",
	//	[]gobot.Connection{home.GetRelayAdaptor()},
	//	[]gobot.Device{home.GetRelHeat(), home.GetRelHeatPump(),
	//		home.SmokeAlarmSauna, home.SmokeAlarmKitchen},
	//	work,
	//)

	//gbot.AddRobot(robot)

	go gbot.Start()

	//stop := scheduleBackup(backupHistoryData, time.Duration(INTERVAL*60)*time.Second, &historyData, HISTORYDATASERIAL)

	err = http.ListenAndServe(":1235", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

	/*
		if stop != nil {
			stop <- true
		}
	*/

}

/*
func scheduleBackup(what func(*home.HistoryData, string), delay time.Duration,
	q *home.HistoryData, l string) chan bool {
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

func backupHistoryData(q *home.HistoryData, local string) {
	historyData.SerializeToFile(local)

	if _, err := DB.UploadFile(local, "/backup/goHome.b64", true, ""); err != nil {
		log.Printf("Error uploading %s: %s\n", local, err)
	} else {
		log.Printf("File %s successfully uploaded\n", local)
	}
}
*/
