package components

import "github.com/yohamta/donburi"

type ServerTimeData struct {
	Time int64
}

var ServerTime = donburi.NewComponentType[ServerTimeData]()
