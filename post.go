package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

func messageHandler(message []byte) {
	messageLoader := gojsonschema.NewBytesLoader(message)
	for _, schema := range schemas {
		if result, err := schema.schema.Validate(messageLoader); err != nil {
			log.Warn().Msgf(`Validate message with schema "%v" erred. %v`, schema.name, err)
		} else if result.Valid() {
			log.Info().Msgf(`Message validated with schema "%v".`, schema.name)
			postMessage(message, schema.postURL)
		}
	}
}

func postMessage(message []byte, postURL string) {
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(message).
		Post(postURL)
	if err != nil {
		log.Error().Msgf("Post message erred. %v", err)
		return
	}
	log.Info().Msgf(`Post message to "%v", return code %v.`, postURL, resp.StatusCode())
	log.Debug().Msgf("%v", resp)
}
