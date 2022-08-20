/*
Create: 2022/8/17
Project: WDNMD
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"github.com/JJApplication/fushin/db/mongo"
)

// 告警存储
// 数据库wdnmd会存储当前系统的告警信息
// 默认存储一个月

type Alarm struct {
	mongo.MetaModel `bson:",inline"`
	Title           string `json:"title" bson:"title"`
	Level           string `json:"level" bson:"level"`
	Message         string `json:"message" bson:"message"`
}

type Message struct {
	ID         string `json:"id" bson:"id"`                   // uuid
	CreateTime int64  `json:"create_time" bson:"create_time"` // 时间戳 精确到s
	Level      string `json:"level" bson:"level"`             // 告警级别
	Info       string `json:"info" bson:"info"`
	Source     string `json:"source" bson:"source"` // 告警来源
}

func (alarm *Alarm) CollectionName() string {
	return "alarm"
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

// CreateOneAlarm 信息以纯文本的方式存储
// mongo存储失败时不做任何事
func CreateOneAlarm(title, level, message string) error {
	return mongoC.Coll(&Alarm{}).Create(&Alarm{
		Title:   title,
		Level:   level,
		Message: message,
	})
}
