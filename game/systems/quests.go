package systems

import (
	"fmt"
	"space-game_mk4/game/components"
	qsts "space-game_mk4/game/components/quests"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/yohamta/donburi/ecs"
)

type quests struct {
	emptyDelay int
}

var Quests = &quests{}

func (q *quests) Update(e *ecs.ECS) {
	if entry, ok := components.Quests.First(e.World); ok {
		quest := components.Quests.Get(entry)
		if quest.Current == nil {
			return
		}
		if q.emptyDelay > 0 {
			q.emptyDelay--
			return
		}

		if quest.Current.Complete == nil {
			if ref, ok := lo.Find(qsts.AllQuests, func(i components.QuestItemData) bool {
				return i.Name == quest.Current.Name
			}); ok {
				quest.Current.Complete = ref.Complete
			}
		}
		if quest.Current.Complete(e.World) {
			components.UserFeedEvent.Publish(e.World, components.FeedItemData{
				ID:      "feeditem:" + uuid.NewString(),
				Message: fmt.Sprintf("Quest %s complete", quest.Current.Name),
			})
			quest.Completed = append(quest.Completed, *quest.Current)
			next := lo.Filter(qsts.AllQuests, func(i components.QuestItemData, _ int) bool {
				return i.ID == quest.Current.Next
			})
			if len(next) == 0 {
				quest.Current = nil
				return
			} else {
				quest.Current = &next[0]
			}
		}
		q.emptyDelay = 60
	}
}
