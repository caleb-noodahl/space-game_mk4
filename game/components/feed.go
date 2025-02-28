package components

import (
	"bytes"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type FeedItemData struct {
	ID       string `json:"id"`
	SourceID string `json:"source_id"`
	Message  string `json:"message"`
}

type FeedData struct {
	buf         bytes.Buffer
	current     string
	Items       []string
	feedBuf     []rune
	bufIndex    int
	LastPublish int64
	LastHash    string
	incrementor int
}

func (f *FeedData) Tick() {
	// Delay for 30 ticks.
	f.incrementor++
	if f.incrementor < 5 {
		return
	}
	f.incrementor = 0

	// If there's no current notification (nothing in feedBuf),
	// load the next one if available.
	if len(f.feedBuf) == 0 {
		if len(f.Items) == 0 {
			return // Nothing to do.
		}
		// Load the next message.
		f.current = f.Items[0]
		f.Items = f.Items[1:]
		// Convert once to runes.
		f.feedBuf = []rune(f.current[0:1])
		return
	}
	// Showing the text
	if f.bufIndex < int(len(f.current)-1) {
		f.bufIndex++
		f.feedBuf = append(f.feedBuf, []rune(f.current[f.bufIndex:f.bufIndex+1])...)
		return
	}

	// pop one rune from the beginning of feedBuf.
	f.feedBuf = f.feedBuf[1:]

	// If we've removed all the runes, clear the current message.
	if len(f.feedBuf) == 0 {
		f.current = ""
		f.LastHash = ""
		f.bufIndex = 0
	}
}

func (f *FeedData) CurrentString() string {
	for _, r := range f.feedBuf {
		f.buf.WriteString(string(r))
	}
	defer f.buf.Truncate(0)
	return f.buf.String()
}

var Feed = donburi.NewComponentType[FeedData]()
var StationFeedEvent = events.NewEventType[FeedItemData]()
var UserFeedEvent = events.NewEventType[FeedItemData]()
