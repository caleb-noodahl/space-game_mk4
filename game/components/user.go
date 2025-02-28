package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type UserProfileData struct {
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Authed   bool   `json:"authed"`
	Name     string `json:"name"`
	Level    int64  `json:"level"`
	Exp      int64  `json:"exp"`
}

var UserProfile = donburi.NewComponentType[UserProfileData]()
var LoginEvent = events.NewEventType[UserProfileData]()
