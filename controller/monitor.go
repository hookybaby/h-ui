package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"h-ui/model/vo"
	"h-ui/service"
	"h-ui/util"
	"regexp"
)

func MonitorSystem(c *gin.Context) {
	cpuPercent, err := util.GetCpuPercent()
	if err != nil {
		vo.Fail("cpu query failed", c)
		return
	}
	memPercent, err := util.GetMemPercent()
	if err != nil {
		vo.Fail("mem query failed", c)
		return
	}
	diskPercent, err := util.GetDiskPercent()
	if err != nil {
		vo.Fail("disk query failed", c)
		return
	}
	vo.Success(vo.SystemMonitorVo{
		CpuPercent:  cpuPercent,
		MemPercent:  memPercent,
		DiskPercent: diskPercent,
	}, c)
}

func MonitorHysteria2(c *gin.Context) {
	onlineUsers, err := service.Hysteria2Online()
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	var hysteria2MonitorVo vo.Hysteria2MonitorVo
	if len(onlineUsers) > 0 {
		hysteria2MonitorVo.UserTotal = int64(len(onlineUsers))
		var deviceTotal int64 = 0
		for _, value := range onlineUsers {
			deviceTotal += value
		}
		hysteria2MonitorVo.DeviceTotal = deviceTotal
	}

	content, err := util.Exec(fmt.Sprintf("%s version", util.GetHysteria2BinPath()))
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}

	pattern := `v\d+\.\d+\.\d+`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(content, -1)
	version := ""
	if len(matches) > 0 {
		version = matches[0]
		hysteria2MonitorVo.Version = version
	}

	running := service.Hysteria2Status()
	hysteria2MonitorVo.Running = running

	vo.Success(hysteria2MonitorVo, c)
	return
}
