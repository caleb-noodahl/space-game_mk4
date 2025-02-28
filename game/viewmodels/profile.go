package viewmodels

import (
	"fmt"
	"image"
	"space-game_mk4/game/components"

	"github.com/ebitengine/debugui"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
)

type profileViewModel struct {
	World      donburi.World
	showWallet bool
	showQuests bool
}

var ProfileVM = &profileViewModel{}

func (l *profileViewModel) ProfileSummary(ctx *debugui.Context, w donburi.World) {
	if entry, ok := components.UserProfile.First(w); ok {
		profile := components.UserProfile.Get(entry)
		xp := components.XP.Get(entry)
		ctx.Window("Profile", image.Rect(40, 40, 340, 240), func(res debugui.Response, layout debugui.Layout) {
			ctx.SetLayoutRow([]int{60, 120, 40, -1}, 20)
			ctx.Label(profile.Username)
			ctx.Label(fmt.Sprintf("lvl: %v exp: %v", xp.Level, xp.XP))
			if ctx.Button("save") != 0 {
				components.GameStatePublish.Publish(l.World, components.GameStateData{})
				components.StationFeedEvent.Publish(w, components.FeedItemData{
					ID:       "feeditem:" + uuid.NewString(),
					SourceID: profile.ID,
					Message:  "Game Saved!",
				})
			}
			ctx.SetLayoutRow([]int{64, 64, 64, 64 - 1}, 20)
			if ctx.Button("Profile") != 0 {

			}
			if ctx.Button("Wallet") != 0 {
				l.showWallet = !l.showWallet
			}

			if ctx.Button("Quests") != 0 {
				l.showQuests = !l.showQuests
			}

			if ctx.Button("Metrics") != 0 {

			}

			userfeed := components.Feed.Get(entry)
			ctx.SetLayoutRow([]int{-1}, 20)
			ctx.Label(userfeed.CurrentString())

		})

		if l.showWallet {
			l.WalletSummary(ctx, w, entry)
		}
		if l.showQuests {
			l.QuestsSummary(ctx, w, entry)
		}
	}
}

func (l *profileViewModel) WalletSummary(ctx *debugui.Context, w donburi.World, entry *donburi.Entry) {
	wallet := components.Wallet.Get(entry)
	ctx.Window("Wallet", image.Rect(40, 240, 560, 560), func(res debugui.Response, layout debugui.Layout) {
		if ctx.Header("Statement Summary", true) != 0 {
			ctx.SetLayoutRow([]int{64, -1}, 0)
			ctx.Label("Balance")
			ctx.Label(wallet.BalanceDisplay())
			ctx.SetLayoutRow([]int{-1}, -1)
			ctx.LayoutColumn(func() {
				ctx.TreeNode("Latest Transactions", func(res debugui.Response) {
					ctx.Label(wallet.LatestTransactionsTable(3))
				})
			})
		}
	})
}

func (l *profileViewModel) QuestsSummary(ctx *debugui.Context, w donburi.World, entry *donburi.Entry) {
	ctx.Window("Quests", image.Rect(40, 240, 560, 560), func(res debugui.Response, layout debugui.Layout) {
		if ctx.Header("Quests Summary", true) != 0 {
			ctx.SetLayoutRow([]int{64, 64, 64, 64, -1}, 20)
			if entry, ok := components.Quests.First(w); ok {
				quest := components.Quests.Get(entry)
				ctx.Label("Completed")
				ctx.Label(fmt.Sprintf("%v", len(quest.Completed)))

				ctx.Label("Available")
				ctx.Label(fmt.Sprintf("%v", len(quest.Available)))

				ctx.SetLayoutRow([]int{-1}, 20)
				if quest.Current != nil {
					ctx.Label(quest.Current.Name)
					ctx.SetLayoutRow([]int{-1}, 35)
					ctx.Label(quest.Current.Description)
					ctx.SetLayoutRow([]int{40, -1}, 20)
					for _, q := range quest.CurrentSteps {
						ctx.Label("")
						ctx.Checkbox(q.Description, lo.ToPtr(q.Complete))
					}
				}

			}

		}
	})
}
