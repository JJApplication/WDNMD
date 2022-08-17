/*
Create: 2022/8/17
Project: WDNMD
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

// 告警存储
// 数据库wdnmd会存储当前系统的告警信息
// 默认存储一个月

type Message struct {
	ID         string `json:"id" bson:"id"`                   // uuid
	CreateTime int64  `json:"create_time" bson:"create_time"` // 时间戳 精确到s
	Level      string `json:"level" bson:"level"`             // 告警级别
	Info       string `json:"info" bson:"info"`
	Source     string `json:"source" bson:"source"` // 告警来源
}

const (
	SourceApp       = "app"
	SourceSystem    = "system"
	SourceContainer = "container"
	SourceOther     = "other"
)

const (
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)
