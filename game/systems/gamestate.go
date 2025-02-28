package systems

import (
	"encoding/json"
	"space-game_mk4/game/components"
	qst "space-game_mk4/game/components/quests"
	"space-game_mk4/game/components/tasks"
	"strconv"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type gameStateRequest struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
}

func (g *gameStateRequest) Bytes() []byte {
	out, _ := json.Marshal(g)
	return out
}

type gamestate struct {
	World      donburi.World
	once       sync.Once
	serverTime int64
	timetick   int64
}

var GameState *gamestate = &gamestate{}

func (s *gamestate) init() {
}

func (s *gamestate) GameStatePublishEvent(w donburi.World, gsd components.GameStateData) {
	if entry, ok := components.UserProfile.First(w); ok {
		user := components.UserProfile.Get(entry)
		wallet := components.Wallet.Get(entry)
		quests := components.Quests.Get(entry)
		xp := components.XP.Get(entry)
		gsd.UserProfile = *user
		gsd.Wallet = *wallet
		gsd.Quests = quests
		gsd.XP = *xp

		if entry, ok := components.Station.First(w); ok {
			station := components.Station.Get(entry)
			environmental := components.Environmental.Get(entry)
			research := components.Research.Get(entry)
			gsd.Station = *station
			gsd.Environmental = *environmental
			gsd.Research = *research
		}

		components.Employee.Each(w, func(e *donburi.Entry) {
			gsd.Employees = append(gsd.Employees, *components.Employee.Get(e))
			gsd.EmployeeTasks = append(gsd.EmployeeTasks, *components.Task.Get(e))
		})

		if entry, ok := components.ResearchLab.First(w); ok {
			lab := components.ResearchLab.Get(entry)
			gsd.ResearchLab = lab
		}

		if entry, ok := components.MachineShop.First(w); ok {
			msshop := components.MachineShop.Get(entry)
			gsd.MachineShop = msshop
		}

		if entry, ok := components.Dock.First(w); ok {
			dock := components.Dock.Get(entry)
			gsd.Dock = dock
		}

		if entry, ok := components.Quests.First(w); ok {
			quests := components.Quests.Get(entry)
			if len(quests.Completed) == 0 {
				quests.Current = &qst.TutorialChain[0]
			}
			gsd.Quests = quests
		}

		if mq, ok := components.MQTT.First(w); ok {
			mqc := components.MQTT.Get(mq)
			mqc.Publish("gamestate/create/"+user.ID, gsd)
		}
	}
}

func (s *gamestate) GameStateMessageEvent(client mqtt.Client, msg mqtt.Message) {
	gsd := components.GameStateData{}
	err := json.Unmarshal(msg.Payload(), &gsd)
	if err != nil {
		return //todo handle errors better
	}

	if entry, ok := components.Wallet.First(s.World); ok {
		gsd.Wallet.Init() //buffer inits, kinda gross
		components.Wallet.Set(entry, &gsd.Wallet)
	}

	if gsd.Quests != nil {
		if entry, ok := components.Quests.First(s.World); ok {
			components.Quests.Set(entry, gsd.Quests)
		}
	}

	if entry, ok := components.XP.First(s.World); ok {
		components.XP.Set(entry, &gsd.XP)
	}

	if _, ok := components.Station.First(s.World); !ok {
		entity := s.World.Create(components.Station, components.Research, components.Environmental, components.Feed)
		components.Station.Set(s.World.Entry(entity), &gsd.Station)
		components.Environmental.Set(s.World.Entry(entity), &gsd.Environmental)
		components.Feed.Set(s.World.Entry(entity), &components.FeedData{})
		components.Research.Set(s.World.Entry(entity), &gsd.Research)
	}

	if gsd.ResearchLab != nil {
		entity := s.World.Create(components.ResearchLab)
		components.ResearchLab.Set(s.World.Entry(entity), gsd.ResearchLab)
	}
	if gsd.MachineShop != nil {
		entity := s.World.Create(components.MachineShop)
		components.MachineShop.Set(s.World.Entry(entity), gsd.MachineShop)
	}
	if gsd.Dock != nil {
		entity := s.World.Create(components.Dock)
		components.Dock.Set(s.World.Entry(entity), gsd.Dock)
	}

	for i, emp := range s.World.CreateMany(len(gsd.Employees), components.Employee, components.Task) {
		components.Employee.Set(s.World.Entry(emp), &gsd.Employees[i])
		switch gsd.EmployeeTasks[i].Type {
		case components.RepairTask:
			switch gsd.EmployeeTasks[i].Name {
			case "O2 Gen Repair":
				components.Task.Set(s.World.Entry(emp), tasks.New02GenRepairTask(s.World, gsd.EmployeeTasks[i].Duration))
			case "Build Research Laboratory":
				components.Task.Set(s.World.Entry(emp), tasks.NewResearchLabBuildTask(s.World, components.FacilityData[components.ResearchType]{}))
			}

		case components.ResearchTask:
			tmp := gsd.EmployeeTasks[i]
			tmplvl := 1 //todo parse the tmp level

			components.Task.Set(s.World.Entry(emp), tasks.NewResearchTask(s.World, tmp.Name, tmp.Duration, tmp.Difficulty, tmplvl, tmp.Item.(components.ResearchType)))
		}
	}

}

func (s *gamestate) SystemTimeEventHandler(client mqtt.Client, msg mqtt.Message) {
	time, _ := strconv.ParseInt(string(msg.Payload()), 0, 64)
	s.serverTime = time
	s.timetick = 0
}

func (s *gamestate) ServerTime() int64 {
	return s.serverTime
}

func (g *gamestate) ProfileMessageEventHandler(client mqtt.Client, msg mqtt.Message) {
	profile := components.UserProfileData{}
	err := json.Unmarshal(msg.Payload(), &profile)
	if err != nil {
		return
	}
	if entry, ok := components.UserInterface.First(g.World); ok { //todo decouple this
		ui := components.UserInterface.Get(entry)
		ui.SetContentValue(msg.Topic(), profile.ClientID)

		if entry, ok := components.UserProfile.First(g.World); ok {
			auth := components.UserProfile.Get(entry)
			auth.ID = profile.ID
			auth.ClientID = profile.ClientID
			auth.Authed = true
			components.UserProfile.Set(entry, auth)
			ui.SetAuth(auth) //anti-pattern

			if entry, ok := components.MQTT.First(g.World); ok {
				mqtt := components.MQTT.Get(entry)
				// the server will publish our previous save data to this topic if exists
				mqtt.Client.Subscribe("gamestate/"+auth.ClientID, g.GameStateMessageEvent)
				req := gameStateRequest{UserID: auth.ID, ClientID: auth.ClientID}
				mqtt.Client.Publish("gamestate/fetch", req.Bytes())

			}
		}
	}
}

func (s *gamestate) Update(e *ecs.ECS) {
	components.GameStatePublish.ProcessEvents(e.World)
	//hack to get around 10/sec delay in time publishing from server
	if s.timetick == 60 {
		s.serverTime++
		s.timetick = 0
	} else {
		s.timetick++
	}
	if entry, ok := components.ServerTime.First(e.World); ok {
		st := components.ServerTime.Get(entry)
		st.Time = s.serverTime
		components.ServerTime.Set(entry, st)
	}
}
