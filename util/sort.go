package util

import "go.mongodb.org/mongo-driver/bson"

var Sorts = map[string]bson.D{
	"advancements_done": {{"stats.custom.advancements_done", -1}},
	"total_played":      {{"stats.custom.play_time", -1}},
	"blocks_broken":     {{"stats.custom.blocks_broken", -1}},
	"deaths":            {{"stats.custom.deaths", -1}},
}
