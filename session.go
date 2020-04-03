package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

var session string

func initSession() {
	auth()
	verify()
}

func auth() {
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"authKey": cfg.Get("mirai.authKey").String(),
		}).
		Post(cfg.Get("mirai.apiBaseURL").String() + "/auth")
	if err != nil {
		log.Fatal().Msgf("Auth failed. %v", err)
	}
	code := gjson.GetBytes(resp.Body(), "code")
	if code.Exists() && code.Value() == 0.0 {
		session = gjson.GetBytes(resp.Body(), "session").String()
		log.Info().Msg("Auth succeed.")
		log.Debug().Msgf("%v", resp)
		return
	}
	log.Debug().Msgf("%v", resp)
	log.Fatal().Msg("Auth failed.")
}

func verify() {
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"sessionKey": session,
			"qq":         cfg.Get("mirai.qq").Int(),
		}).
		Post(cfg.Get("mirai.apiBaseURL").String() + "/verify")
	if err != nil {
		log.Fatal().Msgf("Verify failed. %v", err)
	}
	code := gjson.GetBytes(resp.Body(), "code")
	if code.Exists() && code.Value() == 0.0 {
		log.Info().Msg("Verify succeed.")
		log.Debug().Msgf("%v", resp)
		return
	}
	log.Debug().Msgf("%v", resp)
	log.Fatal().Msg("Verify failed.")
}
