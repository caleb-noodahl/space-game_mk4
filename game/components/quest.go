package components

import (
	"github.com/yohamta/donburi"
)

type QuestItemData struct {
	ID          string                   `json:"id"`
	Next        string                   `json:"next"`
	Previous    string                   `json:"prev"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Chain       string                   `json:"chain"`
	Complete    func(donburi.World) bool `json:"-"`
}

type QuestStep struct {
	Description string `json:"description"`
	Complete    bool   `json:"complete"`
}

type QuestChainData struct {
	Quests []QuestItemData
}

type QuestData struct {
	Completed    []QuestItemData `json:"completed"`
	Available    []QuestItemData `json:"available"`
	Current      *QuestItemData  `json:"current"`
	CurrentSteps []QuestStep     `json:"current_steps"`
}

var Quests = donburi.NewComponentType[QuestData]()
