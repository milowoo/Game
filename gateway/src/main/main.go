package main

import (
	"gateway/src"
	"gateway/src/internal"
	llog "gateway/src/log"
	"log"
)

func main() {
	logName := "/Users/wuchuangeng/game/logs/" + "gateway.log"
	//日志名 + 文件大小（M为单位） + 打印标志 + 线程数量 （未启动） + 工作协程长度（未启动） + 深度
	initLog := llog.NewLogger2(logName, 1024*2, log.LstdFlags|log.Lshortfile, 8, 1024, 2)
	initLog.Info("main begin ....")
	internal.GLog = initLog

	server, err := gateway.NewServer()
	if err != nil {
		initLog.Error("new server err, %+v", err)
		return
	}

	server.Run()
}
