package main

import (
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

type schema struct {
	name    string
	schema  *gojsonschema.Schema
	postURL string
	block   bool
}

var schemasMap map[string][]*schema

func initSchema() {
	schemasMap = make(map[string][]*schema)
	loadSchema("/all")
	loadSchema("/command")
}

func loadSchema(channel string) {
	var schemas []*schema
	for _, item := range cfg.Get("schemas." + channel).Array() {
		loader := gojsonschema.NewStringLoader(item.Get("schema").Raw)
		s, err := gojsonschema.NewSchema(loader)
		if err != nil {
			log.Warn().Msgf(`Load schema "%v" for channel "%v" error. %v`, item.Get("name").String(), channel, err)
		} else {
			schema := &schema{
				name:    item.Get("name").String(),
				schema:  s,
				postURL: item.Get("postURL").String(),
				block:   item.Get("block").Bool(),
			}
			schemas = append(schemas, schema)
			log.Info().Msgf(`Load schema "%v" for channel "%v" succeed.`, item.Get("name").String(), channel)
		}
	}
	schemasMap[channel] = schemas
	log.Info().Msgf(`Loaded %v schemas for channel "%v".`, len(schemas), channel)
}
