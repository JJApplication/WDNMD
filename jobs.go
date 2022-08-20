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

// 每天早9点
func checkApps() {
	c := cron.NewGroup("0 0 9 * * ?")
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

// 系统定时检查
// 暂时不使用 因为wdnmd运行时 cpu一定是高占用的
// 每2小时运行一次
// 内存基线60%
func systemLoopCheck() {
	c := cron.NewGroup("0 0 0/2 * * ?")
	_, err := c.AddFunc(func() {
		logger.Info("job [systemLoopCheck] run")
		memInfo := getMemUsed()
		if memInfo <= 0.5*100 {
			return
		}
		res, err := udsc.SendWithRes(uds.Req{
			Operation: sendAlarmHtml,
			Data: mailReq(mailInfo{
				Type:     "",
				Message:  systemAlarmAlert(),
				IsFile:   false,
				Subject:  TitleSysAlarm,
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
			logger.ErrorF("run job systemLoopCheck error: %s", err.Error())
		}
		if res.Error != "" {
			logger.ErrorF("res job systemLoopCheck error: %s", res.Error)
		}
	})
	if err != nil {
		logger.ErrorF("add job systemLoopCheck error: %s", err.Error())
		return
	}
	c.Start()
}

// 服务检查 一小时一次
func checkAppsLoop() {
	c := cron.NewGroup("0 0 0/1 * * ?")
	_, err := c.AddFunc(func() {
		logger.Info("job [checkAppsLoop] run")
		var allApps []string
		allApps = append(allApps, wc.Apps...)
		allApps = append(allApps, wc.ExtraApps...)
		appInfos, ok := checkProcess(allApps)
		if !ok {
			logger.Warn("job [checkAppsLoop] some apps not good")
		}

		// 过滤app
		var badApps []appInfo
		for _, app := range appInfos {
			if app.Status == StatusBad {
				badApps = append(badApps, app)
			}
		}
		if len(badApps) <= 0 {
			return
		}
		res, err := udsc.SendWithRes(uds.Req{
			Operation: sendAlarmHtml,
			Data: mailReq(mailInfo{
				Type:     "",
				Message:  appAlarmAlert(badApps),
				IsFile:   false,
				Subject:  TitleAppAlarm,
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
			logger.ErrorF("run job checkAppsLoop error: %s", err.Error())
		}
		if res.Error != "" {
			logger.ErrorF("res job checkAppsLoop error: %s", res.Error)
		}
	})
	if err != nil {
		logger.ErrorF("add job checkAppsLoop error: %s", err.Error())
		return
	}
	c.Start()
}
