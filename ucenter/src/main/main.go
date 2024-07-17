package main

import (
	"log"
	"ucenter/src"
	log2 "ucenter/src/log"
)

func main() {
	logName := "/data/llog/" + "ucenter.log"
	//日志名 + 文件大小（M为单位） + 打印标志 + 线程数量 （未启动） + 工作协程长度（未启动） + 深度
	initLog := log2.NewLogger2(logName, 1024*2, log.LstdFlags|log.Lshortfile, 8, 1024, 2)
	server, err := ucenter.NewServer(initLog)
	if err != nil {
		initLog.Error("new server err, %+v", err)
	}

	server.Run()
}