package components

const (
	Carbon     Material = "carbon"
	CO2        Material = "co2"
	Oxygen     Material = "oxygen"
	Water      Material = "water"
	Silicon    Material = "silicon"
	Nanofiber  Material = "nanofiber"
	Iron       Material = "iron"
	Titanium   Material = "titanium"
	Zeolite    Material = "zeolite"
	SteelPanel Material = "steel panel"

	SolarPanel          Material = "solar panel"
	Circuit             Material = "circuit"
	PolymerGel          Material = "polymer gel"
	NanofiberInsulation Material = "nanofiber insulation"
	TitaniumAlloy       Material = "titanium alloy"
)

func AllMaterials() []Material {
	return []Material{
		Carbon,
		CO2,
		Oxygen,
		Water,
		Silicon,
		Nanofiber,
		Iron,
		Titanium,
	}
}

type Storage[T Material | Component] struct {
	Contents T   `json:"contents"`
	Current  int `json:"current"`
	Max      int `json:"max"`
}

type Machine[InputType Material | Component, OutputType Material | Component] struct {
	Level      int          `json:"level"`
	Name       string       `json:"name"`
	Inputs     []InputType  `json:"inputs"`
	Outputs    []OutputType `json:"outputs"`
	Rate       int          `json:"rate"`
	Durability int          `json:"durability"`
	Alert      bool         `json:"alert"`
}

var (
	CO2Scrubber = Component{
		Name: "co2 scrubber",
		Recipe: map[Material]int{
			Carbon:     2,
			Zeolite:    1,
			SolarPanel: 1,
			Circuit:    1,
		},
	}
	OxygenGenerator = Component{
		Name: "oxygen generator",
		Recipe: map[Material]int{
			SteelPanel: 2,
			Circuit:    2,
			SolarPanel: 1,
		},
	}
	WaterRecyclingSystem = Component{
		Name: "water recycling system",
		Recipe: map[Material]int{
			PolymerGel:          2,
			Circuit:             5,
			NanofiberInsulation: 20,
		},
	}
	ReinforcedHullPanel = Component{
		Name: "Reinforced Hull Panel",
		Recipe: map[Material]int{
			SteelPanel:    6,
			TitaniumAlloy: 6,
		},
	}
	SolarArray = Component{
		Name: "Solar Array",
		Recipe: map[Material]int{
			SolarPanel:    1,
			Circuit:       1,
			TitaniumAlloy: 1,
		},
	}
)

func AllComponents() []Component {
	return []Component{
		CO2Scrubber,
		OxygenGenerator,
		WaterRecyclingSystem,
		ReinforcedHullPanel,
		SolarArray,
	}
}
