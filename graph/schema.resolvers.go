package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/bortexel/stats-server/database"
	"github.com/bortexel/stats-server/graph/generated"
	"github.com/bortexel/stats-server/graph/model"
	"github.com/bortexel/stats-server/models"
	"github.com/bortexel/stats-server/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *advancementResolver) Display(context.Context, *models.Advancement) (*model.AdvancementDisplay, error) {
	return &model.AdvancementDisplay{
		Tab:   &[]string{"test"}[0],
		Type:  &[]string{"advancement"}[0],
		Icon:  &[]string{"minecraft:iron_pickaxe"}[0],
		Title: &[]string{"Test"}[0],
	}, nil
}

func (r *mutationResolver) UpdatePlayer(ctx context.Context, input model.UpdatePlayer) (*models.Player, error) {
	if !util.IsAuthorized(ctx) {
		return nil, errors.New("unauthorized")
	}

	advancements := make([]*models.Advancement, 0)
	stats := input.Stats

	for _, advancement := range input.Advancements {
		if !advancement.Done {
			continue
		}

		advancements = append(advancements, &models.Advancement{
			Key: advancement.Key,
		})
	}

	blocksBroken := 0
	for _, value := range input.Stats.Mined {
		blocksBroken += value.(int)
	}

	stats.Custom["blocks_broken"] = blocksBroken
	stats.Custom["advancements_done"] = len(advancements)

	var player models.Player
	newPlayer := models.InputPlayer{
		UUID:         input.UUID,
		Name:         input.Name,
		Stats:        stats,
		Advancements: advancements,
	}

	collection := database.Database.Collection(input.Server)
	err := collection.FindOne(context.Background(), bson.D{{"uuid", input.UUID}}).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Create a new player
			result, err := collection.InsertOne(context.Background(), newPlayer)
			if err != nil {
				return nil, err
			}

			err = collection.FindOne(context.Background(), bson.D{{"_id", result.InsertedID}}).Decode(&player)
			return &player, err
		} else {
			return nil, err
		}
	} else {
		_, err := collection.UpdateOne(context.Background(), bson.D{{"uuid", input.UUID}}, bson.M{"$set": newPlayer})
		if err != nil {
			return nil, err
		}

		err = collection.FindOne(context.Background(), bson.D{{"uuid", input.UUID}}).Decode(&player)
		return &player, err
	}
}

func (r *queryResolver) Players(_ context.Context, server string, sort *string) ([]*models.Player, error) {
	opts := options.Find()
	if sort != nil {
		if sortData, ok := util.Sorts[*sort]; ok {
			opts.SetSort(sortData)
		}
	}

	cursor, err := database.Database.Collection(server).Find(context.Background(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	var results []*models.Player
	err = cursor.All(context.Background(), &results)
	return results, nil
}

// Advancement returns generated.AdvancementResolver implementation.
func (r *Resolver) Advancement() generated.AdvancementResolver { return &advancementResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type advancementResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
