package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type Material string

type Component struct {
	Name   string           `json:"name"`
	Recipe map[Material]int `json:"recipe"`
}

type OrderType string

const (
	Buy  OrderType = "buy"
	Sell OrderType = "sell"
)

type Order[T Material | Component | EmployeeData] struct {
	ID      string    `json:"id"`
	OwnerID string    `json:"owner_id"`
	Created int64     `json:"created"`
	Expires int64     `json:"expires"`
	Amount  int64     `json:"amount"`
	Price   int64     `json:"price"`
	Type    OrderType `json:"type"`
	Item    T         `json:"item"`
}

type MarketData[T EmployeeData | Component | Material] struct {
	Buys  []Order[T] `json:"buys"`
	Sells []Order[T] `json:"sells"`
}

var MarketTag = donburi.NewTag("Market")
var EmployeeMarket = donburi.NewComponentType[MarketData[EmployeeData]]()
var ComponentMarket = donburi.NewComponentType[MarketData[Component]]()
var MaterialMarket = donburi.NewComponentType[MarketData[Material]]()

var MarketsMaterialsBuyEvent = events.NewEventType[Order[Material]]()
var MarketsMaterialsSellEvent = events.NewEventType[Order[Material]]()

var MarketsComponentsBuyEvent = events.NewEventType[Order[Component]]()
var MarketsComponentsSellEvent = events.NewEventType[Order[Component]]()
