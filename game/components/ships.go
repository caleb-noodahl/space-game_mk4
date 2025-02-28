package components

type ShipStatus string

type ShipData struct {
	Name      string               `json:"name"`
	Status    ShipStatus           `json:"dock_status"`
	MatCargo  []Storage[Material]  `json:"mat_cargo"`
	CompCargo []Storage[Component] `json:"comp_cargo"`
	Depart    int64                `json:"depart"`
	Arrival   int64                `json:"arrival"`
}

const (
	Docked    ShipStatus = "docked"
	Departing ShipStatus = "departing"
	Underway  ShipStatus = "underway"
	OnTask    ShipStatus = "on_task"
)
