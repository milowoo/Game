package ucenter

import (
	"math/rand"
	"ucenter/src/handler"
)

type Closure = func()

func SafeRunClosure(v interface{}, c Closure) {
	defer func() {
		if err := recover(); err != nil {
			//log.Printf("%+v: %s", err, debug.Stack())

		}
	}()

	c()
}

func RunOnHandler(c chan Closure, mgr *handler.HandlerMgr, cb func(mgr *handler.HandlerMgr)) {
	c <- func() {
		cb(mgr)
	}
}

func RandomInt(rand *rand.Rand, min int64, maxPlus1 int64) int64 {
	return min + rand.Int63n((maxPlus1 - min))
}
