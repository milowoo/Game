package game_frame

import "game_frame/src/handler"

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
