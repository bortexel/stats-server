package database

type Advancement struct {
	Key   string `json:"key" bson:"key"`
	Tab   string `json:"tab" bson:"-"`
	Type  string `json:"type" bson:"-"`
	Icon  string `json:"icon" bson:"-"`
	Title string `json:"title" bson:"-"`
}

type Criteria struct {
	Key  string `json:"key" bson:"key"`
	Done string `json:"done" bson:"done"`
}

type StoredPlayer struct {
	ID           string         `json:"-" bson:"_id"`
	UUID         string         `json:"uuid" bson:"uuid"`
	Name         string         `json:"name" bson:"name"`
	Stats        StatsContainer `json:"stats,omitempty" bson:"stats"`
	Advancements []*Advancement `json:"advancements,omitempty" bson:"advancements"`
}

type Player struct {
	UUID         string         `json:"uuid" bson:"uuid"`
	Name         string         `json:"name" bson:"name"`
	Stats        StatsContainer `json:"stats" bson:"stats"`
	Advancements []*Advancement `json:"advancements" bson:"advancements"`
}

type Stat struct {
	Key   string `json:"key" bson:"key"`
	Value int    `json:"value" bson:"value"`
}

type StatsMap map[string]any
type StatsContainer map[StatGroupName]StatsMap

type StatGroupName string

const (
	StatBroken   StatGroupName = "minecraft:broken"
	StatCrafted  StatGroupName = "minecraft:crafted"
	StatCustom   StatGroupName = "minecraft:custom"
	StatDropped  StatGroupName = "minecraft:dropped"
	StatKilledBy StatGroupName = "minecraft:killed_by"
	StatKilled   StatGroupName = "minecraft:killed"
	StatMined    StatGroupName = "minecraft:mined"
	StatPickedUp StatGroupName = "minecraft:picked_up"
	StatUsed     StatGroupName = "minecraft:used"
	StatHelpers  StatGroupName = "bortexel:helpers"
)

var defaultStatGroups = []StatGroupName{
	StatBroken,
	StatCrafted,
	StatCustom,
	StatDropped,
	StatKilledBy,
	StatKilled,
	StatMined,
	StatPickedUp,
	StatUsed,
	StatHelpers,
}

func MakeStatsContainer() StatsContainer {
	container := make(StatsContainer)

	for _, groupName := range defaultStatGroups {
		if _, ok := container[groupName]; !ok {
			container[groupName] = make(StatsMap)
		}
	}

	return container
}
