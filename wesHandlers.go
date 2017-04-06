package main

import (
	"bytes"
	"encoding/json"
	"goHome/home"
	"log"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

/*
MsgCommand - simple comm protocol to deliver infor to server
*/
type MsgCommand struct {
	sync.Mutex
	Command string `json:"Command, string"`
	Object  string `json:"Object, string"`
	Param1  string `json:"Param1, string"`
	Param2  string `json:"Param2, string"`
}

func echoHandler(ws *websocket.Conn) {

	var id = int32(time.Now().Unix())

	conns.lock.Lock()
	conns.ws[id] = ws
	conns.lock.Unlock()
	defer func() {
		conns.lock.Lock()
		delete(conns.ws, id)
		conns.lock.Unlock()
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

	rconns.lock.Lock()
	rconns.ws[id] = ws
	rconns.lock.Unlock()
	defer func() {
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
