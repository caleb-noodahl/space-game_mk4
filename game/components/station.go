package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type StationData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	Faction   string `json:"faction"`
	Employees int    `json:"employees"`
	Lifetime  int64  `json:"lifetime"`
}

var Station = donburi.NewComponentType[StationData]()
var StationCreateEvent = events.NewEventType[StationData]()
var StationEnvironmentalEvent = events.NewEventType[EnvironmentalEventData]()
