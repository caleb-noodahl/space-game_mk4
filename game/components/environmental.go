package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type EnvironmentalData struct {
	OxygenGenerator Machine[Material, Material] `json:"oxygen_generator"`
	WaterRecycler   Machine[Material, Material] `json:"water_recycler"`
	OxygenStorage   Storage[Material]           `json:"oxygen_storage"`
	WaterStorage    Storage[Material]           `json:"water_storage"`
}

type EnvironmentalEventData struct {
	EventText string `json:"event_text"`
	Action    func(*EnvironmentalData) *EnvironmentalData
}

func NewEnvironmentalSystem() *EnvironmentalData {
	return &EnvironmentalData{
		OxygenGenerator: Machine[Material, Material]{
			Level:      1,
			Name:       "O2 Generator",
			Inputs:     []Material{},
			Outputs:    []Material{},
			Rate:       4,
			Durability: 23,
		},
		WaterRecycler: Machine[Material, Material]{
			Level:      1,
			Name:       "Water Recycler",
			Inputs:     []Material{},
			Outputs:    []Material{},
			Rate:       2,
			Durability: 100,
		},
		OxygenStorage: Storage[Material]{
			Contents: Oxygen,
			Current:  100,
			Max:      100,
		},
		WaterStorage: Storage[Material]{
			Contents: Water,
			Current:  100,
			Max:      100,
		},
	}
}

var Environmental = donburi.NewComponentType[EnvironmentalData]()
var EnvironmentalEvent = events.NewEventType[EnvironmentalEventData]()
