package data

import (
	"time"

	"golang.org/x/exp/rand"
)

// Example first names (inspired by characters or naming style in The Expanse)
var firstNames = []string{
	"James",
	"Naomi",
	"Amos",
	"Alex",
	"Chrisjen",
	"Bobbie",
	"Camina",
	"Josephus",
	"Paolo",
	"Clarissa",
	"Elvi",
	"Filip",
	"Julie",
	"Marco",
	"Michio",
	"Roberta",
	"Santiago",
	"Thiago",
	"Rocinante", // Just for fun, though it’s a ship name!
	"Martinez",
}

// Example surnames (inspired by characters or naming style in The Expanse)
var surnames = []string{
	"Holden",
	"Nagata",
	"Burton",
	"Kamal",
	"Avasarala",
	"Draper",
	"Drummer",
	"Johnson",
	"Cortázar",
	"Mao",
	"Okoye",
	"Inaros",
	"Pa",
	"Sanchez",
	"Smith",
	"Baker",
	"Sundrani",
	"Ashford",
	"Monica",
	"Wei",
}

func GetRandomName() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	fn := firstNames[rand.Intn(len(firstNames))]
	sn := surnames[rand.Intn(len(surnames))]
	return fn + " " + sn
}
