package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type FacilityData[T ResearchType | Component | EmployeeData] struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	BuildTime   int          `json:"build_time"`
	Cost        int64        `json:"cost"`
	Type        ResearchType `json:"type"`
}

var ResearchLabCreateEvent = events.NewEventType[FacilityData[ResearchType]]()
var MachineShopCreateEvent = events.NewEventType[FacilityData[Component]]()
var DockCreateEvent = events.NewEventType[FacilityData[ResearchType]]()

// type ManufacturingHub Facility[ResearchType]
// type SecurityStation Facility[ResearchType]
// type Port Facility[ResearchType]

var ResearchLab = donburi.NewComponentType[FacilityData[ResearchType]]()
var MachineShop = donburi.NewComponentType[FacilityData[Component]]()
var Dock = donburi.NewComponentType[FacilityData[ResearchType]]()

var ResearchLabTag = donburi.NewTag("ResearchLabTag")
var MachineShopTag = donburi.NewTag("MachineShopTag")
var DockTag = donburi.NewTag("DockTag")

func (f *FacilityData[ResearchType]) Eligible(w donburi.World) (string, bool) {
	switch f.Type {
	// the administration base research can only unlock the research lab
	case Administration:
		// administration eligibility is a) an admin employee and b) they player has 100 dollars to spend
		if entry, ok := UserProfile.First(w); ok {
			wallet := Wallet.Get(entry)

			if wallet.Balance() < 10 { // todo < 10000
				return "insufficient balance", false
			}
			hasAdmin := false
			Employee.Each(w, func(e *donburi.Entry) {
				emp := Employee.Get(e)
				if emp.Profession == Administration {
					hasAdmin = true
					return
				}
				if hasAdmin {
					return
				}

				for _, top := range AdministrationResearch {
					for _, research := range top {
						if research == emp.Profession {
							hasAdmin = true
						}
					}
				}
			})
			return "missing administrator employee", hasAdmin
		}
	}
	return "unknown", false
}
