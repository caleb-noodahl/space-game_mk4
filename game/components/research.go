package components

import (
	_ "embed"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type ResearchType string

func (r ResearchType) String() string {
	return string(r)
}

const (
	Administration ResearchType = "administration"
	Fabrication    ResearchType = "fabrication"
	Logistics      ResearchType = "logistics"
	Security       ResearchType = "security"
	//administration
	Personnel  ResearchType = "personnel"
	Accounting ResearchType = "accounting"
	Law        ResearchType = "law"
	//personnel
	Recruitment ResearchType = "recruitment"
	Training    ResearchType = "training"
	Morale      ResearchType = "morale"
	//accounting
	ResourceAllocation ResearchType = "resource allocation"
	ProfitOptimization ResearchType = "profit optimization"
	CostReduction      ResearchType = "cost reduction"
	//law
	RegulationCompliance ResearchType = "regulation compliance"
	Safety               ResearchType = "safety"
	LegalAdvocacy        ResearchType = "legal advocacy"
	//recruitment
	TalentScouting  ResearchType = "talent scouting"
	DiversityHiring ResearchType = "diversity hiring"
	//training
	Simulations      ResearchType = "simulations"
	SkillEnhancement ResearchType = "skill enhancement"
	//morale
	Entertainment ResearchType = "entertainment"
	MentalWelfare ResearchType = "mental welfare"
	//fabrication
	Construction  ResearchType = "construction"
	Manufacturing ResearchType = "manufacturing"
	RD            ResearchType = "r&d"

	StructuralEng    ResearchType = "structural eng"
	ModularDesign    ResearchType = "modular design"
	Automation       ResearchType = "automation"
	Processing       ResearchType = "processing"
	AdvMaterials     ResearchType = "adv materials"
	NanoTech         ResearchType = "nano tech"
	Innovation       ResearchType = "innovation"
	Prototyping      ResearchType = "prototyping"
	TechIntegration  ResearchType = "tech integration"
	LoadManagement   ResearchType = "load management"
	SpaceExpansion   ResearchType = "space expansion"
	Parts            ResearchType = "parts"
	RapidAssembly    ResearchType = "rapid assembly"
	RoboticBuilders  ResearchType = "robotic builders"
	AI               ResearchType = "ai"
	Recycling        ResearchType = "recycling"
	EnergyEfficiency ResearchType = "energy efficiency"
	CompositeAlloys  ResearchType = "composite alloys"
	SmartMaterials   ResearchType = "smart materials"
	//logistics
	SupplyChain         ResearchType = "supply chain"
	Transport           ResearchType = "transport"
	Harvesting          ResearchType = "harvesting"
	InventoryControl    ResearchType = "inventory control"
	SupplyRoutes        ResearchType = "supply routes"
	DemandForecasting   ResearchType = "demand forecasting"
	CargoHandling       ResearchType = "cargo handling"
	FleetManagement     ResearchType = "fleet management"
	RouteOptimization   ResearchType = "route optimization"
	AstroidMining       ResearchType = "asteroid mining"
	EnergyCollection    ResearchType = "energy collection"
	WaterExtraction     ResearchType = "water extraction"
	JustInTime          ResearchType = "just in time"
	AutomatedInventory  ResearchType = "automated inventory"
	RouteSecurity       ResearchType = "route security"
	AlternatePathways   ResearchType = "alternate pathways"
	PredictiveAnalytics ResearchType = "predictive analytics"
	SeasonalPlanning    ResearchType = "seasonal planning"
	AutomatedLoaders    ResearchType = "automated loaders"
	MassOptimization    ResearchType = "mass optimization"
	MaintenanceSch      ResearchType = "maintenance sch"
	FuelEfficiency      ResearchType = "fuel efficiency"
	RealTimeTracking    ResearchType = "real-time tracking"
	MiningAutomation    ResearchType = "mining automation"
	SolarOptimization   ResearchType = "solar optimization"
	EnergyStorage       ResearchType = "energy storage"
	IceMining           ResearchType = "ice mining"
	Purification        ResearchType = "purification"
	//security
	StationDefense        ResearchType = "station defense"
	InternalSecurity      ResearchType = "internal security"
	ExternalRelations     ResearchType = "external relations"
	ShieldTechnology      ResearchType = "shield technology"
	WeaponSystems         ResearchType = "weapons systems"
	ThreatDetection       ResearchType = "threat detection"
	Surveillance          ResearchType = "surveillance"
	CounterIntel          ResearchType = "counter intel"
	CrisisManagement      ResearchType = "crisis management"
	Diplomacy             ResearchType = "diplomacy"
	TradeNegotiation      ResearchType = "trade negotiation"
	Espionage             ResearchType = "espionage"
	EnergyShields         ResearchType = "energy shields"
	KineticBarriers       ResearchType = "kinetic barriers"
	KineticTurrets        ResearchType = "kinetic turrets"
	LaserTurrets          ResearchType = "laser turrets"
	MissileDefense        ResearchType = "missile defense"
	SensorArrays          ResearchType = "sensor arrays"
	EarlyWarning          ResearchType = "early warning"
	AIMonitoring          ResearchType = "ai monitoring"
	PrivacyManagement     ResearchType = "privacy management"
	SpyDetection          ResearchType = "spy detection"
	DataSecurity          ResearchType = "data security"
	EvacuationProtocols   ResearchType = "evacuation protocols"
	ContainmentProcedures ResearchType = "containment procedures"
	TreatyNegotiation     ResearchType = "treaty negotiation"
	CulturalAwareness     ResearchType = "cultural awareness"
	Bargaining            ResearchType = "bargaining"
	SanctionsManagement   ResearchType = "sanctions management"
	CovertOps             ResearchType = "covert ops"
	AssetRecruitment      ResearchType = "asset recruitment"
)

var TopLevelResearchTypes = []ResearchType{
	Administration, Fabrication, Logistics, Security,
}

var SecurityResearch = [][]ResearchType{
	{StationDefense, InternalSecurity, ExternalRelations},
	{ShieldTechnology, WeaponSystems, ThreatDetection, Surveillance, CounterIntel, CrisisManagement,
		Diplomacy, TradeNegotiation, Espionage},
	{EnergyShields, KineticBarriers, KineticTurrets, LaserTurrets, MissileDefense,
		SensorArrays, EarlyWarning, AIMonitoring, PrivacyManagement, SpyDetection, DataSecurity,
		EvacuationProtocols, ContainmentProcedures, TreatyNegotiation, CulturalAwareness, Bargaining,
		SanctionsManagement, CovertOps, AssetRecruitment},
}

var AdministrationResearch = [][]ResearchType{
	{Personnel, Accounting, Law},
	{Recruitment, Training, Morale, ResourceAllocation,
		ProfitOptimization, CostReduction, RegulationCompliance, Safety, LegalAdvocacy},
	{TalentScouting, DiversityHiring, Simulations, SkillEnhancement, Entertainment, MentalWelfare},
}

var FabricationResearch = [][]ResearchType{
	{Construction, Manufacturing, RD},
	{StructuralEng, ModularDesign, Automation, Processing, AdvMaterials, NanoTech,
		Innovation, Prototyping, TechIntegration},
	{LoadManagement, SpaceExpansion,
		Parts, RapidAssembly,
		RoboticBuilders, AI,
		Recycling, EnergyEfficiency,
		CompositeAlloys, SmartMaterials,
	},
}

var LogisticsResearch = [][]ResearchType{
	{SupplyChain, Transport, Harvesting},
	{InventoryControl, SupplyRoutes, DemandForecasting, CargoHandling,
		FleetManagement, RouteOptimization, AstroidMining,
		EnergyCollection, WaterExtraction},
	{JustInTime, AutomatedInventory, RouteSecurity, AlternatePathways,
		PredictiveAnalytics, SeasonalPlanning, AutomatedLoaders, MassOptimization,
		MaintenanceSch, FuelEfficiency, RealTimeTracking,
		MiningAutomation, SolarOptimization, EnergyStorage, IceMining, Purification},
}

func AllResearch() []ResearchType {
	out := []ResearchType{}
	for _, types := range append(AdministrationResearch,
		append(FabricationResearch,
			append(SecurityResearch, LogisticsResearch...)...)...) {
		out = append(out, types...)
		out = append(out, TopLevelResearchTypes...)
	}
	return out
}

type ResearchItemData struct {
	Type  ResearchType
	Start int64
	End   int64
	Level int
}

type ResearchData struct {
	Current   *ResearchType
	Start     int64
	End       int64
	Completed map[ResearchType]int
}

func NewResearch() *ResearchData {
	return &ResearchData{
		Completed: map[ResearchType]int{},
	}
}

var Research = donburi.NewComponentType[ResearchData]()
var ResearchStartEvent = events.NewEventType[ResearchItemData]()
var ResearchEndEvent = events.NewEventType[ResearchItemData]()
