package game_frame

import (
	"fmt"
	"game_frame/src/handler"
	"net"
)

func GetHostIp() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	var ip string
	for _, address := range addrList {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}
	return ip
}

type Closure = func()

func SafeRunClosure(v interface{}, c Closure) {
	defer func() {
		if err := recover(); err != nil {
			//log.Printf("%+v: %s", err, debug.Stack())

		}
	}()

	c()
}

func RunOnRoomMgr(c chan Closure, mgr *RoomMgr, cb func(mgr *RoomMgr)) {
	c <- func() {
		cb(mgr)
	}
}

func RunOnRoom(c chan Closure, room *handler.Room, cb func(room *handler.Room)) {
	c <- func() {
		cb(room)
	}
}

func RunOnNatsService(c chan Closure, natsService *NatsService, cb func(natsService *NatsService)) {
	c <- func() {
		cb(natsService)
	}
}
