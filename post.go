package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

func messageHandler(message []byte, channel string) {
	schemas := schemasMap[channel]
	messageLoader := gojsonschema.NewBytesLoader(message)
	for _, schema := range schemas {
		if result, err := schema.schema.Validate(messageLoader); err != nil {
			log.Warn().Msgf(`Validate message from channel "%v" with schema "%v" erred. %v`, channel, schema.name, err)
		} else if result.Valid() {
			log.Info().Msgf(`Message from channel "%v" validated with schema "%v".`, channel, schema.name)
			postMessage(message, channel, schema.postURL)
		}
	}
}

func postMessage(message []byte, channel, postURL string) {
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(message).
		Post(postURL)
	if err != nil {
		log.Error().Msgf(`Post message from channel "%v" erred. %v`, channel, err)
		return
	}
	log.Info().Msgf(`Post message from channel "%v" to "%v", return code %v.`, channel, postURL, resp.StatusCode())
	log.Debug().Msgf("%v", resp)
}
