package main

import (
	"gateway/src"
	llog "gateway/src/log"
	"log"
)

func main() {
	logName := "/data/llog/" + "gateway.log"
	//日志名 + 文件大小（M为单位） + 打印标志 + 线程数量 （未启动） + 工作协程长度（未启动） + 深度
	locLog := llog.NewLogger2(logName, 1024*2, log.LstdFlags|log.Lshortfile, 8, 1024, 2)

	server, err := gateway.NewServer(locLog)
	if err != nil {
		locLog.Error("new server err, %+v", err)
	}

	server.Run()
}
