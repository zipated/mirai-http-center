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

	// 图片文件上传
	e.POST("/uploadImage", handleUploadImage)

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

	// 获取指定Session的配置
	e.GET("/config", handleStandardGetRequest)

	// 设置指定Session的配置
	e.POST("/config", handleStandardPostJSONRequest)

	// 注册指令
	e.POST("/command/register", handleAuthKeyPostJSONRequest)

	// 发送指令
	e.POST("/command/send", handleAuthKeyPostJSONRequest)

	// 获取Mangers
	e.GET("/managers", handleGetRequest)
}

// post + sessionKey
func handleStandardPostJSONRequest(ctx echo.Context) error {
	log.Info().Msgf(`Receive post http request from "%v" to "%v".`, ctx.RealIP(), ctx.Path())
	bodyBytes, _ := ioutil.ReadAll(ctx.Request().Body)
	log.Debug().Msgf("%s", bodyBytes)
	if gjson.ValidBytes(bodyBytes) {
		data, setBytesErr := sjson.SetBytes(bodyBytes, "sessionKey", session)
		if setBytesErr != nil {
			log.Warn().Msg("Http request received is not object json.")
			log.Debug().Msgf("%v", setBytesErr)
			return ctx.NoContent(http.StatusBadRequest)
		}
		client := resty.New()
		client.SetCloseConnection(true)
		resp, err := client.R().
			SetHeader("Content-Type", "application/json;charset=UTF-8").
			SetBody(data).
			Post(cfg.Get("mirai.apiBaseURL").String() + ctx.Path())
		if err != nil {
			log.Error().Msgf("Forward http request erred. %v", err)
			log.Debug().Msgf("%v", err)
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		if gjson.ValidBytes(resp.Body()) {
			log.Info().Msgf(`Forward post http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
			log.Debug().Msgf("%v", resp)
			return ctx.JSON(resp.StatusCode(), gjson.ParseBytes(resp.Body()).Value())
		}
		log.Info().Msgf(`Forward post http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
		log.Debug().Msgf("%v", resp)
		return ctx.String(resp.StatusCode(), string(resp.Body()))
	}
	log.Warn().Msg("Http post request received is not standard json.")
	return ctx.NoContent(http.StatusBadRequest)
}

// post + authKey
func handleAuthKeyPostJSONRequest(ctx echo.Context) error {
	log.Info().Msgf(`Receive post http request from "%v" to "%v".`, ctx.RealIP(), ctx.Path())
	bodyBytes, _ := ioutil.ReadAll(ctx.Request().Body)
	log.Debug().Msgf("%s", bodyBytes)
	if gjson.ValidBytes(bodyBytes) {
		data, setBytesErr := sjson.SetBytes(bodyBytes, "authKey", cfg.Get("mirai.authKey").String())
		if setBytesErr != nil {
			log.Warn().Msg("Http request received is not object json.")
			log.Debug().Msgf("%v", setBytesErr)
			return ctx.NoContent(http.StatusBadRequest)
		}
		client := resty.New()
		client.SetCloseConnection(true)
		resp, err := client.R().
			SetHeader("Content-Type", "application/json;charset=UTF-8").
			SetBody(data).
			Post(cfg.Get("mirai.apiBaseURL").String() + ctx.Path())
		if err != nil {
			log.Error().Msgf("Forward post http request erred. %v", err)
			log.Debug().Msgf("%v", err)
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		if gjson.ValidBytes(resp.Body()) {
			log.Info().Msgf(`Forward post http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
			log.Debug().Msgf("%v", resp)
			return ctx.JSON(resp.StatusCode(), gjson.ParseBytes(resp.Body()).Value())
		}
		log.Info().Msgf(`Forward post http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
		log.Debug().Msgf("%v", resp)
		return ctx.String(resp.StatusCode(), string(resp.Body()))
	}
	log.Warn().Msg("Http request received is not standard json.")
	return ctx.NoContent(http.StatusBadRequest)
}

// get + sessionKey
func handleStandardGetRequest(ctx echo.Context) error {
	log.Info().Msgf(`Receive get http request from "%v" to "%v".`, ctx.RealIP(), ctx.Path())
	log.Debug().Msgf("%v", ctx.QueryString())
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		SetQueryParam("sessionKey", session).
		Get(cfg.Get("mirai.apiBaseURL").String() + ctx.Path() + "?" + ctx.QueryString())
	if err != nil {
		log.Error().Msgf("Forward get http request erred. %v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if gjson.ValidBytes(resp.Body()) {
		log.Info().Msgf(`Forward get http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
		log.Debug().Msgf("%v", resp)
		return ctx.JSON(resp.StatusCode(), gjson.ParseBytes(resp.Body()).Value())
	}
	log.Info().Msgf(`Forward get http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
	log.Debug().Msgf("%v", resp)
	return ctx.String(resp.StatusCode(), string(resp.Body()))
}

// get
func handleGetRequest(ctx echo.Context) error {
	log.Info().Msgf(`Receive get http request from "%v" to "%v".`, ctx.RealIP(), ctx.Path())
	log.Debug().Msgf("%v", ctx.QueryString())
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		Get(cfg.Get("mirai.apiBaseURL").String() + ctx.Path() + "?" + ctx.QueryString())
	if err != nil {
		log.Error().Msgf("Forward get http request erred. %v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if gjson.ValidBytes(resp.Body()) {
		log.Info().Msgf(`Forward get http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
		log.Debug().Msgf("%v", resp)
		return ctx.JSON(resp.StatusCode(), gjson.ParseBytes(resp.Body()).Value())
	}
	log.Info().Msgf(`Forward get http request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
	log.Debug().Msgf("%v", resp)
	return ctx.String(resp.StatusCode(), string(resp.Body()))
}

func handleUploadImage(ctx echo.Context) error {
	log.Info().Msgf(`Receive upload image request from "%v" to "%v".`, ctx.RealIP(), ctx.Path())
	t := ctx.FormValue("type")
	if t == "" {
		log.Error().Msgf("No image type specified.")
		return ctx.NoContent(http.StatusBadRequest)
	}
	file, err := ctx.FormFile("img")
	if err != nil {
		log.Error().Msgf("Read upload image file erred. %v", err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		log.Error().Msgf("Open upload image file erred. %v", err)
		return err
	}
	defer src.Close()
	client := resty.New()
	client.SetCloseConnection(true)
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data;charset=UTF-8").
		SetFormData(map[string]string{
			"sessionKey": session,
			"type":       t,
		}).
		SetFileReader("img", file.Filename, src).
		Post(cfg.Get("mirai.apiBaseURL").String() + ctx.Path())
	if err != nil {
		log.Error().Msgf("Forward upload image request erred. %v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if gjson.ValidBytes(resp.Body()) {
		log.Info().Msgf(`Forward upload image request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
		log.Debug().Msgf("%v", resp)
		return ctx.JSON(resp.StatusCode(), gjson.ParseBytes(resp.Body()).Value())
	}
	log.Info().Msgf(`Forward upload image request to "%v", return code %v.`, ctx.Path(), resp.StatusCode())
	log.Debug().Msgf("%v", resp)
	return ctx.String(resp.StatusCode(), string(resp.Body()))
}
