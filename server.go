package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bortexel/stats-server/data"
	"github.com/bortexel/stats-server/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActionHandler func(r *http.Request, input []byte) (output any, err error, status int)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	var next ActionHandler

	switch r.Method {
	case http.MethodOptions: // CORS
		next = HandleRoot
	case http.MethodGet: // Root (health checks, etc.)
		next = HandleRoot
	case http.MethodPost: // Request leaderboard
		next = HandlePlayerInfo
	case http.MethodPatch: // Update player info
		next = ConfiguredAuthorizationMiddleware(HandleUpdatePlayer)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseData, err, status := next(r, body)
	if err != nil {
		w.WriteHeader(status)

		if status >= 500 {
			log.Println("Error serving", r.Method, "request:", err)
		}

		return
	}

	if responseData != nil {
		handleData(w, status, responseData)
	} else {
		w.WriteHeader(status)
	}
}

func HandleRoot(_ *http.Request, _ []byte) (any, error, int) {
	return nil, nil, http.StatusNoContent
}

func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
}

type SortDirection string

func (d SortDirection) getValue() int {
	if d == SortDirectionAscending {
		return 1
	}

	return -1
}

const (
	SortDirectionAscending  = "ascending"
	SortDirectionDescending = "descending"
)

type SortOptions struct {
	Field     StatField     `json:"field"`
	Direction SortDirection `json:"direction"`
}

func (o SortOptions) GetDirection() SortDirection {
	if o.Direction == SortDirectionAscending {
		return SortDirectionAscending
	}

	return SortDirectionDescending
}

type ServerIdentifier struct {
	ServerName string `json:"serverName"`
	Season     int    `json:"season"`
}

func (i ServerIdentifier) String() string {
	return fmt.Sprintf("%s_%d", i.ServerName, i.Season)
}

func (r LeaderboardRequest) ShouldSort() bool {
	return r.Sort.Field.FieldName != "" && r.Sort.Field.GroupName != ""
}

type StatField struct {
	GroupName string `json:"groupName"`
	FieldName string `json:"fieldName"`
}

func (f *StatField) RemoveSpecialCharacters() {
	f.GroupName = strings.ReplaceAll(f.GroupName, ".", "")
	f.GroupName = strings.ReplaceAll(f.GroupName, "$", "")
	f.FieldName = strings.ReplaceAll(f.FieldName, "$", "")
	f.FieldName = strings.ReplaceAll(f.FieldName, "$", "")
}

func (f *StatField) GetFullPath() string {
	f.RemoveSpecialCharacters()
	return fmt.Sprintf("stats.%s.%s", f.GroupName, f.FieldName)
}

type LeaderboardRequest struct {
	Sort               SortOptions      `json:"sort"`
	Server             ServerIdentifier `json:"server"`
	PlayerUUID         string           `json:"playerUUID"`
	PlayerName         string           `json:"playerName"`
	StatsFilter        []StatField      `json:"filter"`
	ReturnAdvancements bool             `json:"returnAdvancements"`
	LimitExpansionKey  string           `json:"limitExpansionKey"`
}

const MaxRecords int64 = 100

func (r LeaderboardRequest) getRecordLimit() int64 {
	if expectedKey, ok := os.LookupEnv("LIMIT_EXPANSION_KEY"); ok && expectedKey == r.LimitExpansionKey {
		return 0
	}

	return MaxRecords
}

func (r LeaderboardRequest) makeProjection() bson.D {
	projection := bson.D{{"name", 1}, {"uuid", 1}}

	if r.ReturnAdvancements {
		projection = append(projection, bson.E{Key: "advancements", Value: 1})
	}

	for _, field := range r.StatsFilter {
		projection = append(projection, bson.E{Key: field.GetFullPath(), Value: 1})
	}

	return projection
}

func (r LeaderboardRequest) makeFilter() bson.D {
	filter := bson.D{}

	if r.PlayerUUID != "" {
		filter = append(filter, bson.E{Key: "uuid", Value: r.PlayerUUID})
	}

	if r.PlayerName != "" {
		filter = append(filter, bson.E{Key: "name", Value: r.PlayerName})
	}

	return filter
}

func HandlePlayerInfo(_ *http.Request, body []byte) (any, error, int) {
	var request LeaderboardRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		return nil, err, http.StatusUnprocessableEntity
	}

	opts := options.Find()
	if request.ShouldSort() {
		opts.SetSort(bson.D{{request.Sort.Field.GetFullPath(), request.Sort.GetDirection().getValue()}})
	}

	opts.SetProjection(request.makeProjection())
	opts.SetLimit(request.getRecordLimit())

	cursor, err := database.Database.Collection(request.Server.String()).
		Find(context.Background(), request.makeFilter(), opts)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	var results []*database.StoredPlayer
	err = cursor.All(context.Background(), &results)

	if err == mongo.ErrNoDocuments || results == nil || len(results) == 0 {
		return nil, nil, http.StatusNotFound
	}

	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return results, nil, http.StatusOK
}

type UpdatePlayerRequest struct {
	Server       ServerIdentifier        `json:"server"`
	UUID         string                  `json:"uuid"`
	Name         string                  `json:"name"`
	Stats        database.StatsContainer `json:"stats"`
	Advancements []*AdvancementInput     `json:"advancements"`
}

type AdvancementInput struct {
	Key  string `json:"key"`
	Done bool   `json:"done"`
}

func HandleUpdatePlayer(_ *http.Request, body []byte) (any, error, int) {
	var request UpdatePlayerRequest
	request.Stats = database.MakeStatsContainer()
	err := json.Unmarshal(body, &request)
	if err != nil {
		return nil, err, http.StatusUnprocessableEntity
	}

	stats := request.Stats
	advancements := FormatAdvancements(request.Advancements)
	AppendTotalStats(stats, len(advancements))

	var player database.StoredPlayer
	newPlayer := database.Player{
		UUID:         request.UUID,
		Name:         request.Name,
		Stats:        stats,
		Advancements: advancements,
	}

	collection := database.Database.Collection(request.Server.String())
	err = collection.FindOne(context.Background(), bson.D{{"uuid", request.UUID}}).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create a new player
			result, err := collection.InsertOne(context.Background(), newPlayer)
			if err != nil {
				return nil, err, http.StatusInternalServerError
			}

			err = collection.FindOne(context.Background(), bson.D{{"_id", result.InsertedID}}).Decode(&player)
			if err != nil {
				return nil, err, http.StatusInternalServerError
			}

			return player, nil, http.StatusOK
		} else {
			return nil, err, http.StatusInternalServerError
		}
	} else {
		_, err := collection.UpdateOne(context.Background(), bson.D{{"uuid", request.UUID}}, bson.M{"$set": newPlayer})
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}

		err = collection.FindOne(context.Background(), bson.D{{"uuid", request.UUID}}).Decode(&player)
		if err != nil {
			return nil, err, http.StatusInternalServerError
		}

		return player, nil, http.StatusOK
	}
}

func AppendTotalStats(stats database.StatsContainer, advancementsCount int) {
	var deaths int64
	var playTime int64

	if actualDeaths, ok := stats[database.StatCustom]["minecraft:deaths"]; ok {
		deaths = int64(actualDeaths.(float64))
	}

	if actualPlayTime, ok := stats[database.StatCustom]["minecraft:play_time"]; ok {
		playTime = int64(actualPlayTime.(float64) / 20)
	}

	stats[database.StatTotals]["bortexel:deaths"] = deaths
	stats[database.StatTotals]["bortexel:play_time"] = playTime
	stats[database.StatTotals]["bortexel:blocks_placed"] = stats[database.StatUsed].GetValueSum(data.IsBlock)
	stats[database.StatTotals]["bortexel:blocks_broken"] = stats[database.StatMined].GetValueSum(database.EmptyPredicate)
	stats[database.StatTotals]["bortexel:advancements_done"] = advancementsCount
}

func FormatAdvancements(inputAdvancements []*AdvancementInput) []*database.Advancement {
	advancements := make([]*database.Advancement, 0)

	for _, advancement := range inputAdvancements {
		if !advancement.Done {
			continue
		}

		advancements = append(advancements, &database.Advancement{
			Key: advancement.Key,
		})
	}

	return advancements
}

func handleData(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	body, err := json.Marshal(data)
	if err != nil {
		return
	}

	_, _ = w.Write(body)
}
