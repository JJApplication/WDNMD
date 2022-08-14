/*
Create: 2022/8/14
Project: WDNMD
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"github.com/JJApplication/fushin/cron"
	"github.com/JJApplication/fushin/server/uds"
)

// 后台定时执行的任务

func healthCheck() {
	c := cron.NewGroup(cron.EveryFmt("10m"))
	_, err := c.AddFunc(func() {
		logger.Info("job [healthCheck] run")
	})
	if err != nil {
		logger.ErrorF("add job healthCheck error: %s", err.Error())
		return
	}
	c.Start()
}

// 每天早晚9点
func checkApps() {
	c := cron.NewGroup("0 0 9,21 * * ?")
	_, err := c.AddFunc(func() {
		logger.Info("job [checkApps] run")
		var allApps []string
		allApps = append(allApps, wc.Apps...)
		allApps = append(allApps, wc.ExtraApps...)
		appInfos, ok := checkProcess(allApps)
		if !ok {
			logger.Warn("job [checkApps] some apps not good")
		}

		res, err := udsc.SendWithRes(uds.Req{
			Operation: sendAlarmHtml,
			Data: mailReq(mailInfo{
				Type:     "",
				Message:  appAlarmInfo(appInfos),
				IsFile:   false,
				Subject:  TitleAppInfo,
				Attach:   nil,
				To:       []string{wc.To},
				Cc:       nil,
				Bcc:      nil,
				SyncTask: false,
				CronJob:  "",
			}),
			From: WDNMD,
			To:   []string{Hermes},
		})
		if err != nil {
			logger.ErrorF("run job checkApps error: %s", err.Error())
		}
		if res.Error != "" {
			logger.ErrorF("res job checkApps error: %s", res.Error)
		}
	})
	if err != nil {
		logger.ErrorF("add job checkApps error: %s", err.Error())
		return
	}
	c.Start()
}

// 每天早上8点
func systemCheck() {
	c := cron.NewGroup("0 0 8 * * ?")
	_, err := c.AddFunc(func() {
		logger.Info("job [systemCheck] run")
		res, err := udsc.SendWithRes(uds.Req{
			Operation: sendAlarmHtml,
			Data: mailReq(mailInfo{
				Type:     "",
				Message:  systemAlarmInfo(),
				IsFile:   false,
				Subject:  TitleSystemInfo,
				Attach:   nil,
				To:       []string{wc.To},
				Cc:       nil,
				Bcc:      nil,
				SyncTask: false,
				CronJob:  "",
			}),
			From: WDNMD,
			To:   []string{Hermes},
		})
		if err != nil {
			logger.ErrorF("run job systemCheck error: %s", err.Error())
		}
		if res.Error != "" {
			logger.ErrorF("res job systemCheck error: %s", res.Error)
		}
	})
	if err != nil {
		logger.ErrorF("add job systemCheck error: %s", err.Error())
		return
	}
	c.Start()
}
