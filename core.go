/*
Create: 2022/8/14
Project: WDNMD
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	client "github.com/JJApplication/fushin/client/uds"
	"github.com/JJApplication/fushin/log"
	"github.com/JJApplication/fushin/server/uds"
)

// NewServer 新建uds服务器 用于心跳
func NewServer() *uds.UDSServer {
	s := uds.Default(wc.UnixAddress)
	s.Option.AutoRecover = true
	s.Option.AutoCheck = false
	logger.InfoF("%s uds server run @ [%s]", WDNMD, s.Name)
	return s
}

// NewClient 新建uds客户端
func NewClient() *client.UDSClient {
	logger.InfoF("%s uds client dial @ [%s]", WDNMD, wc.Talker)
	return &client.UDSClient{
		Addr:        wc.Talker,
		MaxRecvSize: 1 << 20,
	}
}

func NewLogger() *log.Logger {
	return log.Default(WDNMD)
}

// InitJobs 初始化定时任务
func InitJobs() {
	healthCheck()
	checkApps()
	systemCheck()
}
