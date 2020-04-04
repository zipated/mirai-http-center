package main

func init() {
	initConfig()
	initSchema()
	initSession()
}

func main() {
	wsEnd := make(chan int)
	go initWebsocket(wsEnd)
	go initHTTP()
	for {
		select {
		case <-wsEnd:
			go initWebsocket(wsEnd)
		}
	}
}
