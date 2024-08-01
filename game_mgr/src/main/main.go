package main

import (
	"game_mgr/src"
	"game_mgr/src/config"
	"game_mgr/src/internal"
	log2 "game_mgr/src/log"
	"log"
	"os"
	"os/signal"
)

func main() {
	logName := "/Users/wuchuangeng/game/logs/" + "game_mgr.log"
	//日志名 + 文件大小（M为单位） + 打印标志 + 线程数量 （未启动） + 工作协程长度（未启动） + 深度
	internal.GLog = log2.NewLogger2(logName, 1024*2, log.LstdFlags|log.Lshortfile, 8, 1024, 2)
	initLog := internal.GLog
	globalConfig, err := config.NewGlobalConfig()
	if err != nil {
		initLog.Info("LoadConfig err ")
		panic(err)
		return
	}

	initLog.Info("config %+v", globalConfig)
	server, err := game_mgr.NewHttpService(globalConfig)
	if err != nil {
		initLog.Error("new server err, %+v", err)
		return
	}
	server.Run()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		initLog.Warn("exit svr by signal ...")
		server.Quit()
	}()
}
