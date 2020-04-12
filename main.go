package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

var flagVersion bool
var wsConns map[string]*websocket.Conn
var wsClosed bool
var (
	version   = ""
	goVersion = ""
	buildTime = ""
	commitID  = ""
)

func init() {
	flag.BoolVar(&flagVersion, "v", false, "Print version information")
	wsConns = make(map[string]*websocket.Conn)
	wsClosed = false
}

func main() {
	flag.Parse()
	if flagVersion {
		fmt.Println("Version:", version)
		fmt.Println("Go Version:", goVersion)
		fmt.Println("Build Time:", buildTime)
		fmt.Println("Git Commit ID:", commitID)
		return
	}
	initConfig()
	initSchema()
	initSession()
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
