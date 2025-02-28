package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type EmployeeData struct {
	ID         string       `json:"id"`
	XP         int64        `json:"xp"`
	Name       string       `json:"name"`
	Age        int          `json:"age"`
	Level      int          `json:"level"`
	Profession ResearchType `json:"profession"`
	Joined     int64        `json:"joined"`
	Salary     int64        `json:"salary"`
	BaseSalary int64        `json:"base_salary"`
}

var Employee = donburi.NewComponentType[EmployeeData]()
var EmployeeMarketBuyEvent = events.NewEventType[EmployeeData]()
var EmployeeMarketSellEvent = events.NewEventType[EmployeeData]()
