package repository

import (
	"app/config"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) chatColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("chat")
}

func (r *Repo) CreateChatIndexes(ctx context.Context) ([]string, error) {
	indexes, err := r.otpColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{"id", 1},
		}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return nil, err
	}
	return indexes, nil
}
