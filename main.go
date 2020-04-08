package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

var wsConns map[string]*websocket.Conn
var wsClosed bool

func init() {
	initConfig()
	initSchema()
	initSession()
	wsConns = make(map[string]*websocket.Conn)
	wsClosed = false
}

func main() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	enableWebsocket()
	go startWebsocket("/all")
	go startWebsocket("/command")
	go initHTTP()
	<-sig
	wsClosed = true
	for _, wsConn := range wsConns {
		wsConn.Close()
	}
}
