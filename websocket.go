package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

var wsConn *websocket.Conn

func initWebsocket(wsEnd chan int) {
	enableWebsocket()
	var err error
	wsConn, _, err = websocket.DefaultDialer.Dial(cfg.Get("websocket.baseURL").String()+"/all?sessionKey="+session, nil)
	if err != nil {
		log.Error().Msgf("Websocket erred. %v", err)
		log.Info().Msg("Reconnecting websocket...")
		wsEnd <- 1
		return
	}
	log.Info().Msg("Websocket connected.")
	startReadMessage(wsEnd)
}

func enableWebsocket() {
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"sessionKey":      session,
			"enableWebsocket": true,
		}).
		Post(cfg.Get("mirai.apiBaseURL").String() + "/config")
	if err != nil {
		log.Error().Msgf("Enable websocket erred. %v", err)
		return
	}
	code := gjson.GetBytes(resp.Body(), "code")
	if code.Exists() && code.Value() == 0.0 {
		log.Info().Msg("Enable websocket succeed.")
		log.Debug().Msgf("%v", resp)
		return
	}
	log.Error().Msg("Enable websocket erred.")
	log.Debug().Msgf("%v", resp)
}

func startReadMessage(wsEnd chan int) {
	for {
		messageType, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Error().Msgf("Websocket erred. %v", err)
			log.Info().Msg("Reconnecting websocket...")
			wsConn.Close()
			wsEnd <- 1
			return
		}
		if messageType == websocket.TextMessage {
			log.Info().Msg("Receive message from websocket.")
			log.Debug().Msgf("%s", message)
		}
	}
}
