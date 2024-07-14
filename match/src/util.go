package match

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"math/rand"
)

func JsonToServices(result string) (*[]model.Service, error) {
	var services []model.Service
	err := json.Unmarshal([]byte(result), &services)
	if err != nil {
		return nil, err
	}
	//if len(service.Hosts) == 0 {
	//	logger.Warnf("instance list is empty,json string:%s", result)
	//}
	return &services, nil

}
func ToJsonString(object interface{}) string {
	js, _ := json.Marshal(object)
	return string(js)
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

func RunOnMatch(c chan Closure, mgr *GameMatch, cb func(mgr *GameMatch)) {
	c <- func() {
		cb(mgr)
	}
}

func RandomInt(rand *rand.Rand, min int64, maxPlus1 int64) int64 {
	return min + rand.Int63n((maxPlus1 - min))
}
