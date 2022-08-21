/*
Create: 2022/8/6
Project: WDNMD
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	client "github.com/JJApplication/fushin/client/uds"
	"github.com/JJApplication/fushin/db/mongo"
	"github.com/JJApplication/fushin/log"
	"github.com/JJApplication/fushin/server/uds"
	"github.com/JJApplication/fushin/utils/json"
)

var udsc *client.UDSClient
var logger *log.Logger
var mongoC *mongo.Mongo

func main() {
	logger = NewLogger()
	udsc = NewClient()
	wdnmdServer := NewServer()
	wdnmdServer.AddFunc("ping", func(c *uds.UDSContext, req uds.Req) {
		_ = c.Response(uds.Res{
			Error: "",
			Data:  "",
			From:  WDNMD,
			To:    nil,
		})
	})

	wdnmdServer.AddFunc("push", func(c *uds.UDSContext, req uds.Req) {
		var data AlarmBase
		if err := json.Json.UnmarshalFromString(req.Data, &data); err != nil {
			_ = c.Response(uds.Res{
				Error: err.Error(),
				Data:  "",
				From:  WDNMD,
				To:    nil,
			})
		} else {
			go func() {
				err := CreateOneAlarm(data.Title, data.Level, data.Message)
				if err != nil {
					logger.ErrorF("push message to mongo error: %s", err.Error())
				}
			}()
		}
	})

	wdnmdServer.AddFunc("pull", func(c *uds.UDSContext, req uds.Req) {

	})

	wdnmdServer.AddFunc("purge", func(c *uds.UDSContext, req uds.Req) {

	})

	// 初始化数据库
	mongoC = NewMongo()
	err := mongoC.Init()
	if err != nil {
		logger.ErrorF("%s mongo client init error: %s", WDNMD, err.Error())
	}
	// 初始化客户端
	err = udsc.Dial()
	if err != nil {
		logger.ErrorF("%s uds client dial error: %s", WDNMD, err.Error())
	}
	// 初始化任务
	InitJobs()
	if err = wdnmdServer.Listen(); err != nil {
		logger.ErrorF("%s server start error: %s", WDNMD, err.Error())
	}
}
