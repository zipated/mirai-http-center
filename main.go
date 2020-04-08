package main

import (
	"github.com/gorilla/websocket"
)

var wsConns map[string]*websocket.Conn

func init() {
	initConfig()
	initSchema()
	initSession()
	wsConns = make(map[string]*websocket.Conn)
}

func main() {
	done := make(chan bool)
	enableWebsocket()
	go startWebsocket("/all")
	go startWebsocket("/command")
	go initHTTP()
	<-done
}
