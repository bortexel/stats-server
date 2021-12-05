package models

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

type Player struct {
	ID           string         `json:"id" bson:"_id"`
	UUID         string         `json:"uuid" bson:"uuid"`
	Name         string         `json:"name" bson:"name"`
	Stats        *Stats         `json:"stats" bson:"stats"`
	Advancements []*Advancement `json:"advancements" bson:"advancements"`
}

type InputPlayer struct {
	UUID         string         `json:"uuid" bson:"uuid"`
	Name         string         `json:"name" bson:"name"`
	Stats        *StatsInput    `json:"stats" bson:"stats"`
	Advancements []*Advancement `json:"advancements" bson:"advancements"`
}

type Stat struct {
	Key   string `json:"key" bson:"key"`
	Value int    `json:"value" bson:"value"`
}

type Stats struct {
	Custom   map[string]interface{} `json:"custom" bson:"custom"`
	Mined    map[string]interface{} `json:"mined" bson:"mined"`
	Broken   map[string]interface{} `json:"broken" bson:"broken"`
	Crafted  map[string]interface{} `json:"crafted" bson:"crafted"`
	Used     map[string]interface{} `json:"used" bson:"used"`
	PickedUp map[string]interface{} `json:"picked_up" bson:"picked_up"`
	Dropped  map[string]interface{} `json:"dropped" bson:"dropped"`
	Killed   map[string]interface{} `json:"killed" bson:"killed"`
	KilledBy map[string]interface{} `json:"killed_by" bson:"killed_by"`
}

type StatInput struct {
	Key   string `json:"key" bson:"key"`
	Value int    `json:"value" bson:"value"`
}

type StatsInput struct {
	Custom   map[string]interface{} `json:"custom" bson:"custom"`
	Mined    map[string]interface{} `json:"mined" bson:"mined"`
	Broken   map[string]interface{} `json:"broken" bson:"broken"`
	Crafted  map[string]interface{} `json:"crafted" bson:"crafted"`
	Used     map[string]interface{} `json:"used" bson:"used"`
	PickedUp map[string]interface{} `json:"picked_up" bson:"picked_up"`
	Dropped  map[string]interface{} `json:"dropped" json:"dropped"`
	Killed   map[string]interface{} `json:"killed" json:"killed"`
	KilledBy map[string]interface{} `json:"killed_by" bson:"killed_by"`
}
