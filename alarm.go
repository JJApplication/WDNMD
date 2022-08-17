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
	TitleAppInfo    = "[AOO INFO]服务状态"
	TitleAppAlarm   = "[APP ALARM]服务告警"
	TitleSysAlarm   = "[SYSTEM ALARM]系统告警"
)

// 生成系统信息
// 非告警 定时发送
func systemAlarmInfo() string {
	p, f, v := getPlatform()
	return fmt.Sprintf(
		"<h4>环境信息</h4>操作系统: %s</br>系统家族: %s</br>系统版本: %s</br>内核版本: %s</br>上次启动时间: %s</br>CPU核心数: %d</br>CPU使用率: %v</br>内存使用率: %v</br>可用内存: %v</br>运行进程数: %d</br>",
		p,
		f,
		v,
		getKernel(),
		getBoot(),
		getCpuCount(),
		fmt.Sprintf("%f%%", getCpu()),
		fmt.Sprintf("%f%%", getMemUsed()),
		fmt.Sprintf("%f bytes", getMemAvail()),
		getProcessCount(),
	)
}

func systemAlarmLoopCheckInfo() string {
	return fmt.Sprintf("<h4>系统告警</h4><p>内存占用达到阈值</p></br><p>内存占用: %f%%</br>空闲内存: %f bytes</p>", getMemUsed(), getMemAvail())
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
	return fmt.Sprintf("%s%s", header, body)
}

// 服务监控告警信息
// 上报出错的微服务
func appAlarmLoopInfo(appInfos []appInfo) string {
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
	return fmt.Sprintf("%s%s", header, body)
}
