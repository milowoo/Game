package main

import (
	"game_mgr/src"
	log2 "game_mgr/src/log"
	"log"
	"os"
	"os/signal"
)

func main() {
	logName := "/data/llog/" + "gateway.log"
	//日志名 + 文件大小（M为单位） + 打印标志 + 线程数量 （未启动） + 工作协程长度（未启动） + 深度
	initLog := log2.NewLogger2(logName, 1024*2, log.LstdFlags|log.Lshortfile, 8, 1024, 2)
	config, err := game_mgr.NewGlobalConfig(initLog)
	if err != nil {
		initLog.Info("LoadConfig err ")
		panic(err)
		return
	}
	server, err := game_mgr.NewHttpService(config, initLog)
	if err != nil {
		initLog.Error("new server err, %+v", err)
	}

	server.Run()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		initLog.Warn("exit svr by signal ...")
		server.OnQuit()
	}()
}
