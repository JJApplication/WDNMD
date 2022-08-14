/*
Create: 2022/8/14
Project: WDNMD
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"path/filepath"
	"strings"

	"github.com/JJApplication/fushin/utils/env"
	"github.com/JJApplication/octopus_meta"
)

// 配置 从环境变量加载

// v1 仅监控微服务的运行状态
// 默认监控jjapplication
// 基于octopus-meta 未发布的服务跳过
type wdnmdConfig struct {
	Apps        []string // 过滤掉未发布态后的jjapplication
	ExtraApps   []string // 额外监控的服务
	UnixAddress string   // 监听地址
	Talker      string   // 接收方地址
	To          string   // 收件人
	AppRoot     string
}

const (
	SymbolDot   = "."
	SymbolComma = ","
)

var wc wdnmdConfig

func init() {
	envLoader := env.EnvLoader{}
	wc = wdnmdConfig{
		Apps:        getApps(envLoader.Get("APP_ROOT").Raw()),
		ExtraApps:   strings.Split(strings.TrimSpace(envLoader.Get("ExtraApps").Raw()), SymbolComma),
		UnixAddress: envLoader.Get("UnixAddress").Raw(),
		Talker:      envLoader.Get("Talker").Raw(),
		To:          envLoader.Get("To").Raw(),
		AppRoot:     envLoader.Get("APP_ROOT").Raw(),
	}
}

func getApps(appRoot string) []string {
	var apps []string
	appsMap, err := octopus_meta.AutoLoad()
	if err != nil {
		return apps
	}
	for k, v := range appsMap {
		if v.ReleaseStatus == octopus_meta.Published {
			apps = append(apps, filepath.Join(appRoot, k))
		}
	}
	return apps
}
