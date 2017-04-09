package irRelay

import (
	"errors"
	"io"
	"irrigation/wsHandler"
	"log"
	"net/http"

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

	/*TESTS emanbles test mode*/
	//TESTS = true
)

//"github.com/hybridgroup/gobot"
var r = raspi.NewRaspiAdaptor("raspi")

type irrigationRelay struct {
	relayMode string
	Relay     *gpio.LedDriver
	Wh        *wsHandler.WsHandler
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
		log.Println("Wwitch off relay:" + r.Relay.Name())
	}
}

/*New - returns new relay instance */
func New(name string, pin string, w *wsHandler.WsHandler) Ir {
	rel := Ir{OFF, gpio.NewLedDriver(r, name, pin), w}
	rel.Relay.On()
	//rel.Relay =
	relays[pin] = rel
	//rel.Wh = w
	http.HandleFunc("/control/"+rel.Relay.Name(), rel.RelayHandler)
	return rel
}

/*
SetMode sets the behavior for the relay
*/
func (r *Ir) SetMode(str string) error {
	if str != ON && str != OFF && str != AUTO {
		return errors.New("Wrong parameter: " + str + " constant ON/OFF/AUTO expected")
	}

	if str == ON {
		err := r.Relay.Off()
		if err != nil {
			return err
		}
	}

	if str == OFF || str == AUTO {
		err := r.Relay.On()
		if err != nil {
			return err
		}
	}

	r.relayMode = str
	return nil
}

/*
GetMode sets the behavior for the relay
*/
func (r *Ir) GetMode() string {
	return r.relayMode
}

/*RelayHandler - http handler for simple rest */
func (r *Ir) RelayHandler(w http.ResponseWriter, re *http.Request) {
	//defer reportPump()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	state := re.FormValue("state")
	//log.Println(state)

	if len(state) == 0 {
		//log.Println("state requested:")
		io.WriteString(w, ":"+r.GetMode())
		return
	}

	errr := r.SetMode(state)
	if errr != nil {
		http.Error(w, errr.Error(), http.StatusInternalServerError)
		return
	}

	r.Wh.ReportWsEvent("relayStateChanged", r.Relay.Name())

}
