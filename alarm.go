/*
Create: 2022/8/14
Project: WDNMD
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"fmt"
	"time"
)

// 生成告警信息

const (
	TitleSystemInfo = "[SYSTEM INFO]系统状态"
	TitleAppInfo    = "[APP INFO]服务状态"
	TitleAppAlarm   = "[APP ALARM]服务告警"
	TitleSysAlarm   = "[SYSTEM ALARM]系统告警"
)

// 生成系统信息
// 非告警 定时发送
func systemAlarmInfo() string {
	p, f, v := getPlatform()
	sysInfo := SystemAlarmInfo{
		System:       p,
		Family:       f,
		Version:      v,
		Kernel:       getKernel(),
		BootTime:     getBoot(),
		CpuCount:     getCpuCount(),
		CpuPercent:   getCpu(),
		MemUsed:      getMemUsed(),
		MemAvail:     getMemAvail(),
		ProcessCount: getProcessCount(),
	}
	go mongoSystemInfo(sysInfo)
	return fmt.Sprintf(SystemInfoTemplate,
		sysInfo.System,
		sysInfo.Family,
		sysInfo.Version,
		sysInfo.Kernel,
		sysInfo.BootTime,
		sysInfo.CpuCount,
		fmt.Sprintf("%f%%", sysInfo.CpuPercent),
		fmt.Sprintf("%f%%", sysInfo.MemUsed),
		fmt.Sprintf("%f bytes", sysInfo.MemAvail),
		sysInfo.ProcessCount,
	)
}

func systemAlarmAlert() string {
	sysAlert := SystemAlarmAlert{
		MemUsed:  getMemUsed(),
		MemAvail: getMemAvail(),
	}
	go mongoSystemAlert(sysAlert)
	return fmt.Sprintf(SystemAlertTemplate, sysAlert.MemUsed, sysAlert.MemAvail)
}

// 服务监控定时信息
func appAlarmInfo(appInfos []appInfo) string {
	header := `<h4 style="color: #378de5">微服务状态</h4>`
	body := ""
	for _, app := range appInfos {
		if app.Status == StatusBad {
			appBody := `<div style="font-size: 0.85rem;margin-bottom: 0.5rem">%s</div>`
			pre := fmt.Sprintf(`<strong style="font-size: 1rem;color: #dc4905">[*] %s</strong>`, app.App)
			infoBody := fmt.Sprintf(appBody, pre)
			body = body + infoBody
			continue
		}
		appBody := `<div style="font-size: 0.85rem;margin-bottom: 0.5rem">%s%s%s%s%s%s%s%s%s</div>`
		pre := fmt.Sprintf(`<strong style="font-size: 1rem;color: #30a24c">[*] %s</strong>`, app.App)
		t := fmt.Sprintf(`<p style="margin: 4px 0"><strong>创建时间:</strong> %s</p>`, time.Unix(app.CreateTime/1000, 0).Local().Format("2006-01-02 15:04:05"))
		pid := fmt.Sprintf(`<p style="margin: 4px 0"><strong>PID:</strong> %d</p>`, app.Pid)
		cpu := fmt.Sprintf(`<p style="margin: 4px 0"><strong>CPU使用:</strong> %f%%</p>`, app.CpuPercent)
		mem := fmt.Sprintf(`<p style="margin: 4px 0"><strong>内存使用:</strong> %f%%</p>`, app.MemPercent)
		memRss := fmt.Sprintf(`<p style="margin: 4px 0"><strong>内存占用:</strong> %d bytes</p>`, app.MemRss)
		conn := fmt.Sprintf(`<p style="margin: 4px 0"><strong>连接数:</strong> %d</p>`, app.Connections)
		thread := fmt.Sprintf(`<p style="margin: 4px 0"><strong>线程数:</strong> %d</p>`, app.NumberThreads)
		io := fmt.Sprintf(`<p style="margin: 4px 0"><strong>IO使用:</strong> <span style="padding: 0 4px">读次数: %d</span><span style="padding: 0 4px">写次数: %d</span><span style="padding: 0 4px">读字节: %d</span><span style="padding: 0 4px">写字节: %d</span></p>`,
			app.IOReadCount,
			app.IOWriteCount,
			app.IOReadBytes,
			app.IOWriteBytes,
		)

		infoBody := fmt.Sprintf(appBody, pre, t, pid, cpu, mem, memRss, conn, thread, io)
		body = body + infoBody
	}
	go mongoAppInfo(appInfos)
	return fmt.Sprintf("%s%s", header, body)
}

// 服务监控告警信息
// 上报出错的微服务
func appAlarmAlert(appInfos []appInfo) string {
	header := `<h4 style="color: #dc4905">微服务状态异常</h4>`
	body := ""
	for _, app := range appInfos {
		if app.Status == StatusBad {
			appBody := `<div style="font-size: 0.85rem;margin-bottom: 0.5rem">%s</div>`
			pre := fmt.Sprintf(`<strong style="font-size: 1rem;color: #dc4905">[*] %s</strong>`, app.App)
			infoBody := fmt.Sprintf(appBody, pre)
			body = body + infoBody
			continue
		}
	}
	go mongoAppAlert(appInfos)
	return fmt.Sprintf("%s%s", header, body)
}

// 存储上报的通知

func mongoSystemInfo(sysInfo SystemAlarmInfo) {
	body := fmt.Sprintf(SystemInfoTemplateMongo,
		sysInfo.System,
		sysInfo.Family,
		sysInfo.Version,
		sysInfo.Kernel,
		sysInfo.BootTime,
		sysInfo.CpuCount,
		fmt.Sprintf("%f%%", sysInfo.CpuPercent),
		fmt.Sprintf("%f%%", sysInfo.MemUsed),
		fmt.Sprintf("%f bytes", sysInfo.MemAvail),
		sysInfo.ProcessCount)
	err := CreateOneAlarm(TitleSystemInfo, LevelInfo, body)
	if err != nil {
		logger.ErrorF("[SystemInfo] push message to mongo error: %s", err.Error())
	}
}

func mongoSystemAlert(sysAlert SystemAlarmAlert) {
	body := fmt.Sprintf(SystemAlertTemplateMongo, sysAlert.MemUsed, sysAlert.MemAvail)
	err := CreateOneAlarm(TitleSysAlarm, LevelError, body)
	if err != nil {
		logger.ErrorF("[SystemAlert] push message to mongo error: %s", err.Error())
	}
}

func mongoAppInfo(appInfos []appInfo) {
	var hasBad bool
	for _, app := range appInfos {
		if app.Status == StatusBad {
			hasBad = true
			break
		}
	}
	if hasBad {
		body := "微服务状态正常"
		err := CreateOneAlarm(TitleAppAlarm, LevelInfo, body)
		if err != nil {
			logger.ErrorF("[AppInfo] push message to mongo error: %s", err.Error())
		}
	}
}

func mongoAppAlert(appInfos []appInfo) {
	body := "微服务状态异常\n"
	for _, app := range appInfos {
		if app.Status == StatusBad {
			appBody := fmt.Sprintf(`[*] %s\n`, app.App)
			body = body + appBody
			continue
		}
	}
	err := CreateOneAlarm(TitleAppAlarm, LevelError, body)
	if err != nil {
		logger.ErrorF("[AppAlert] push message to mongo error: %s", err.Error())
	}
}
