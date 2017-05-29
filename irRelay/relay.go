package irRelay

import (
	"encoding/json"
	"errors"
	"io"
	"irrigation/wsHandler"
	"log"
	"net/http"
	"time"

	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

const (
	/*ON Constant */
	ON = "On"
	/*OFF constant*/
	OFF = "Off"
	/*AUTO constant*/
	AUTO = "Auto"

	/*INTERVAL - default interval to keep relay on*/
	INTERVAL = 1
	/*TESTS emanbles test mode*/
	//TESTS = true
)

type fn func()

//"github.com/hybridgroup/gobot"
var r = raspi.NewRaspiAdaptor("raspi")

type irrigationRelay struct {
	RelayMode    string `json:"RelayMode, string"`
	Relay        *gpio.LedDriver
	Wh           *wsHandler.WsHandler
	RelayState   bool `json:"State, boolean"`
	stateChanged fn
	stop         chan bool
}

/*Ir irrigation relay type */
type Ir irrigationRelay

/*Relays initiated relays*/
var relays map[string]Ir

func init() {
	relays = make(map[string]Ir)
}

/*
Stop - Set relays to default position
*/
func Stop() {
	log.Println("Set relays to default state")
	for _, r := range relays {
		err := r.Relay.On()
		if err != nil {
			log.Println(err.Error())
		}
		r.stateChanged()
		r.Wh.ReportWsEvent("relayStateChanged", r.Relay.Name())
		log.Println("Wwitch off relay:" + r.Relay.Name())
	}
}

/*New - returns new relay instance */
func New(name string, pin string, w *wsHandler.WsHandler, f fn) Ir {
	rel := Ir{"", gpio.NewLedDriver(r, name, pin), w, false, f, nil}
	rel.Relay.On()
	//rel.Relay =
	relays[pin] = rel
	//rel.Wh = w
	//http.HandleFunc("/control/"+rel.Relay.Name(), rel.RelayHandler)
	return rel
}

/*
SetMode sets the behavior for the relay
*/
func (r *Ir) SetMode(str string) error {
	log.Println("irRelay.SetMode: " + r.Relay.Name())
	if str != ON && str != OFF && str != AUTO {
		log.Println("irRelay.SetMode: Wrong parameter")
		return errors.New("Wrong parameter: " + str + " constant ON/OFF/AUTO expected")
	}

	if str == ON {
		log.Println("irRelay.SetMode: On")
		err := r.Relay.Off()
		if err != nil {
			return err
		}
		log.Println("irRelay.SetMode: ", r.Relay.State())
	}

	if str == OFF || str == AUTO {
		log.Println("irRelay.SetMode: On")
		err := r.Relay.On()
		if err != nil {
			return err
		}
		log.Println("irRelay.SetMode: ", r.Relay.State())
	}

	r.RelayMode = str
	log.Println("irRelay.SetMode: set to ", r.RelayMode, " : ", r.GetMode())

	if r.stop != nil {
		log.Println("Dropping timer")
		close(r.stop)
	}

	r.stateChanged()
	r.Wh.ReportWsEvent("relayStateChanged", r.Relay.Name())

	if r.GetMode() == ON {
		log.Println("Try to set scheduler")
		r.stop = r.scheduleRelayAuto(turnoff, time.Duration(INTERVAL*60)*time.Second)
	}

	log.Println("All done")
	return nil
}

/*
GetMode sets the behavior for the relay
*/
func (r *Ir) GetMode() string {
	log.Println("return mode for relay: ", r.RelayState, " name: ", r.Relay.Name())
	return r.RelayMode
}

/*
GetState sets the behavior for the relay
*/
func (r *Ir) GetState() bool {
	r.RelayState = r.Relay.State()
	return !r.RelayState
}

/*RelayHandler - http handler for simple rest */
func (r *Ir) RelayHandler(w http.ResponseWriter, re *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	st := re.FormValue("mode")
	log.Println("Mode: " + st)

	//set or get
	if len(st) == 0 {
		//log.Println("state requested:")
		r.RelayState = r.GetState()
		b, err := r.ToJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, string(b))
		return
	}

	//set
	errr := r.SetMode(st)
	if errr != nil {
		http.Error(w, errr.Error(), http.StatusInternalServerError)
		return
	}

}

func (r *Ir) scheduleRelayAuto(what func(r *Ir), delay time.Duration) chan bool {
	stop := make(chan bool)
	log.Println("Relay timer scheduled")

	go func() {
		select {
		case <-time.After(delay):
			log.Println("returning to AUTO")
			what(r)
			log.Println("returned auto mode")
			return
		case <-stop:
			return
		}
	}()

	log.Println("Scheduler is prepared...")
	return stop
}

func turnoff(r *Ir) {
	log.Println("Switching to Auto...")
	r.SetMode(AUTO)
}

/*
ToJSON returns serialized date
*/
func (r *Ir) ToJSON() (d []byte, err error) {

	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	return b, nil
}
