package main

import (
	"bytes"
	"encoding/json"
	"goHome/home"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

func execute(c *MsgCommand) error {
	defer c.Unlock()
	c.Lock()

	if c.Object == home.CommandOnPumpr1 {
		return home.OnHeat()
	}

	if c.Object == home.CommandOffPumpr1 {
		return home.OffHeat()
	}

	return nil
}

func relHandler(ws *websocket.Conn) {

	var id = int32(time.Now().Unix())

	wh := WsHandler.New(id, ws)
	defer func() {
		wh
		rconns.lock.Lock()
		delete(rconns.ws, id)
		rconns.lock.Unlock()
	}()

	msg := make([]byte, 512)
	for {
		n, err := ws.Read(msg)
		if err != nil {
			log.Println(err)
			return
		}

		x := new(MsgCommand)
		log.Printf("Receive: %s\n", msg[:n])
		if err := json.NewDecoder(bytes.NewReader(msg)).Decode(x); err == nil {
			execute(x)
		} else {
			log.Println(err)
		}
	}
}
