package main

func init() {
	initConfig()
	initSession()
}

func main() {
	wsEnd := make(chan int)
	initWebsocket(wsEnd)
	for {
		select {
		case <-wsEnd:
			initWebsocket(wsEnd)
		}
	}
}
