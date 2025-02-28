package systems

import (
	"fmt"
	"space-game_mk4/game/components"
	"space-game_mk4/utils"

	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type feed struct {
}

var Feed = &feed{}

func (s *feed) StationFeedEventHandler(w donburi.World, item components.FeedItemData) {
	var time int64
	if entry, ok := components.ServerTime.First(w); ok {
		t := components.ServerTime.Get(entry)
		time = t.Time
	}

	if entry, ok := components.Station.First(w); ok {
		hash := utils.Hash(fmt.Sprintf("%s%s", item.SourceID, item.Message))
		feed := components.Feed.Get(entry)

		if feed.LastHash != hash && !lo.ContainsBy(feed.Items, func(s string) bool {
			return item.Message == s
		}) {
			feed.Items = append(feed.Items, item.Message)
			feed.LastPublish = time
			feed.LastHash = hash
		}
	}
}

func (f *feed) UserFeedEventHandler(w donburi.World, item components.FeedItemData) {
	var time int64
	if entry, ok := components.ServerTime.First(w); ok {
		t := components.ServerTime.Get(entry)
		time = t.Time
	}

	if entry, ok := components.UserProfile.First(w); ok {
		hash := utils.Hash(fmt.Sprintf("%s%s", item.SourceID, item.Message))
		feed := components.Feed.Get(entry)
		if feed.LastHash != hash {
			feed.Items = append(feed.Items, item.Message)
			feed.LastPublish = time
			feed.LastHash = hash
		}
	}
}

func (f *feed) Update(e *ecs.ECS) {
	components.StationFeedEvent.ProcessEvents(e.World)
	components.UserFeedEvent.ProcessEvents(e.World)

	components.Feed.Each(e.World, func(e *donburi.Entry) {
		feed := components.Feed.Get(e)
		feed.Tick()
	})
}
