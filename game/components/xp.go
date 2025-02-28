package components

import "github.com/yohamta/donburi"

type XPData struct {
	XP    int64 `json:"xp"`
	Level int   `json:"level"`
}

var XP = donburi.NewComponentType[XPData]()
