package components

import (
	"github.com/yohamta/donburi/features/events"
)

type GameStateData struct {
	XP            XPData                      `json:"xp"`
	UserProfile   UserProfileData             `json:"user_profile"`
	Station       StationData                 `json:"station"`
	Wallet        WalletData                  `json:"wallet"`
	Employees     []EmployeeData              `json:"employees"`
	EmployeeTasks []TaskData                  `json:"employee_tasks"`
	Environmental EnvironmentalData           `json:"environmental"`
	ServerTime    int64                       `json:"server_time"`
	TickIterator  int64                       `json:"tick_iterator"`
	Research      ResearchData                `json:"research"`
	ResearchLab   *FacilityData[ResearchType] `json:"research_facility"`
	MachineShop   *FacilityData[Component]    `json:"machine_shop"`
	Dock          *FacilityData[ResearchType] `json:"dock"`
	Quests        *QuestData                  `json:"quests"`
}

var GameStatePublish = events.NewEventType[GameStateData]()
