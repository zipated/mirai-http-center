package main

import (
	"github.com/rs/zerolog/log"
	"github.com/xeipuuv/gojsonschema"
)

type schema struct {
	name    string
	schema  *gojsonschema.Schema
	postURL string
}

var schemas []*schema

func initSchema() {
	for _, item := range cfg.Get("schemas").Array() {
		loader := gojsonschema.NewStringLoader(item.Get("schema").Raw)
		s, err := gojsonschema.NewSchema(loader)
		if err != nil {
			log.Warn().Msgf(`Load schema "%v" error. %v`, item.Get("name").String(), err)
		} else {
			schema := &schema{
				name:    item.Get("name").String(),
				schema:  s,
				postURL: item.Get("postURL").String(),
			}
			schemas = append(schemas, schema)
			log.Info().Msgf(`Load schema "%v" succeed.`, item.Get("name").String())
		}
	}
	log.Info().Msgf("Loaded %v schemas.", len(schemas))
}
