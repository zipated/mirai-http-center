package main

func init() {
	initConfig()
	initSchema()
	initSession()
}

func main() {
	wsAllEnd := make(chan int)
	wsCommandEnd := make(chan int)
	enableWebsocket()
	go initWebsocket(wsAllEnd, "/all")
	go initWebsocket(wsCommandEnd, "/command")
	go initHTTP()
	for {
		select {
		case <-wsAllEnd:
			go initWebsocket(wsAllEnd, "/all")
		case <-wsCommandEnd:
			go initWebsocket(wsCommandEnd, "/command")
		}
	}
}
