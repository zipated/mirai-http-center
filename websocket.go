package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

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

func startWebsocket(channel string) {
	for {
		if wsClosed {
			return
		}
		connectWebsocket(channel)
	}
}

func connectWebsocket(channel string) {
	if wsClosed {
		return
	}
	log.Info().Msgf(`Connecting websocket to channel "%v"...`, channel)
	wsConn, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf(
			"%s%s?authKey=%s&sessionKey=%s",
			cfg.Get("mirai.wsBaseURL").String(),
			channel,
			cfg.Get("mirai.authKey").String(),
			session,
		),
		nil,
	)
	if err != nil {
		log.Error().Msgf(`Connect websocket to channel "%v" erred. %v`, channel, err)
		return
	}
	defer wsConn.Close()
	wsConns[channel] = wsConn
	log.Info().Msgf(`Websocket connected to channel "%v".`, channel)
	readWebsocket(channel)
}

func readWebsocket(channel string) {
	wsConn := wsConns[channel]
	for {
		messageType, message, err := wsConn.ReadMessage()
		if err != nil {
			if !wsClosed {
				log.Error().Msgf(`Listen websocket from channel "%v" erred. %v`, channel, err)
			}
			return
		}
		if messageType == websocket.TextMessage {
			log.Info().Msgf(`Receive websocket message from channel "%v".`, channel)
			log.Debug().Msgf("%s", message)
			messageHandler(message, channel)
		}
	}
}
