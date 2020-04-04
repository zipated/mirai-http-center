package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func initHTTP() {
	authKey := cfg.Get("http.authKey").String()
	e := echo.New()
	e.Logger.SetOutput(ioutil.Discard)
	e.Use(middleware.Recover())
	e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == authKey, nil
	}))

	router(e)

	log.Fatal().Msgf("Start HTTP server failed. %v", e.Start(cfg.Get("http.host").String()))
}

func router(e *echo.Echo) {
	// 发送好友消息
	e.POST("/sendFriendMessage", handleStandardPostJSONRequest)

	// 发送群消息
	e.POST("/sendGroupMessage", handleStandardPostJSONRequest)

	// 发送图片消息（通过URL）
	e.POST("/sendImageMessage", handleStandardPostJSONRequest)

	// 撤回消息
	e.POST("/recall", handleStandardPostJSONRequest)

	// 获取Bot收到的消息和事件
	e.GET("/fetchMessage", handleStandardGetRequest)

	// 通过messageId获取一条被缓存的消息
	e.GET("/messageFromId", handleStandardGetRequest)

	// 获取好友列表
	e.GET("/friendList", handleStandardGetRequest)

	// 获取群列表
	e.GET("/groupList", handleStandardGetRequest)

	// 获取群成员列表
	e.GET("/memberList", handleStandardGetRequest)

	// 群全体禁言
	e.POST("/muteAll", handleStandardPostJSONRequest)

	// 群解除全体禁言
	e.POST("/unmuteAll", handleStandardPostJSONRequest)

	// 群禁言群成员
	e.POST("/mute", handleStandardPostJSONRequest)

	// 群解除群成员禁言
	e.POST("/unmute", handleStandardPostJSONRequest)

	// 移除群成员
	e.POST("/kick", handleStandardPostJSONRequest)

	// 群设置
	e.POST("/groupConfig", handleStandardPostJSONRequest)

	// 获取群设置
	e.GET("/groupConfig", handleStandardGetRequest)

	// 修改群员资料
	e.POST("/memberInfo", handleStandardPostJSONRequest)

	// 获取群员资料
	e.GET("/memberInfo", handleStandardGetRequest)
}

func handleStandardPostJSONRequest(ctx echo.Context) error {
	log.Info().Msgf(`Receive post http request from "%v" to "%v".`, ctx.RealIP(), ctx.Path())
	bodyBytes, _ := ioutil.ReadAll(ctx.Request().Body)
	log.Debug().Msgf("%s", bodyBytes)
	if gjson.ValidBytes(bodyBytes) {
		data, setBytesErr := sjson.SetBytes(bodyBytes, "sessionKey", session)
		if setBytesErr != nil {
			log.Warn().Msg("Http request received is not standard json.")
			return ctx.NoContent(http.StatusBadRequest)
		}
		client := resty.New()
		client.SetCloseConnection(true)
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(data).
			Post(cfg.Get("mirai.apiBaseURL").String() + ctx.Path())
		if err != nil {
			log.Error().Msgf("Forward http request erred. %v", err)
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		if gjson.ValidBytes(resp.Body()) {
			log.Info().Msgf(`Forward http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
			log.Debug().Msgf("%v", resp)
			return ctx.JSON(resp.StatusCode(), gjson.ParseBytes(resp.Body()).Value())
		}
		log.Info().Msgf(`Forward http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
		log.Debug().Msgf("%v", resp)
		return ctx.String(resp.StatusCode(), string(resp.Body()))
	}
	log.Warn().Msg("Http request received is not standard json.")
	return ctx.NoContent(http.StatusBadRequest)
}

func handleStandardGetRequest(ctx echo.Context) error {
	log.Info().Msgf(`Receive get http request from "%v" to "%v".`, ctx.RealIP(), ctx.Path())
	log.Debug().Msgf("%v", ctx.QueryString())
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetQueryParam("sessionKey", session).
		Get(cfg.Get("mirai.apiBaseURL").String() + ctx.Path() + "?" + ctx.QueryString())
	if err != nil {
		log.Error().Msgf("Forward http request erred. %v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if gjson.ValidBytes(resp.Body()) {
		log.Info().Msgf(`Forward http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
		log.Debug().Msgf("%v", resp)
		return ctx.JSON(resp.StatusCode(), gjson.ParseBytes(resp.Body()).Value())
	}
	log.Info().Msgf(`Forward http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
	log.Debug().Msgf("%v", resp)
	return ctx.String(resp.StatusCode(), string(resp.Body()))
}
