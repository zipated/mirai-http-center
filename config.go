package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

var cfg gjson.Result

func initConfig() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	configPath, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "config.json"))
	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		log.Fatal().Msgf("Open config file failed. %v", err)
	}
	cfgBytes, _ := ioutil.ReadAll(file)
	if gjson.ValidBytes(cfgBytes) {
		cfg = gjson.ParseBytes(cfgBytes)
		setGlobalLogLevel()
		log.Info().Msg("Load config file succeed.")
		log.Debug().Msgf("Config:\n%s", cfgBytes)
	} else {
		log.Fatal().Msg("Config file format error.")
	}
}

func setGlobalLogLevel() {
	switch strings.ToLower(cfg.Get("log.level").String()) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
