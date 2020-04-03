package main

func init() {
	initConfig()
	initSession()
	initSchema()
}

func main() {
	wsEnd := make(chan int)
	go initWebsocket(wsEnd)
	for {
		select {
		case <-wsEnd:
			go initWebsocket(wsEnd)
		}
	}
}
