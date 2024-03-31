package repository

import (
	"app/config"
	"app/entity"
	"app/errors"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) tokenColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("tokens")
}

func (r *Repo) CreateTokenIndexes(ctx context.Context) ([]string, error) {
	indexes, err := r.tokenColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{"id", 1},
		}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func (r *Repo) CreateToken(ctx context.Context, token *entity.Token) error {
	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", token},
	}
	_, err := r.tokenColl().UpdateOne(ctx, bson.M{"id": token.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.TokenExists()
		}
		return err
	}
	return nil
}

func (r *Repo) GetTokenById(ctx context.Context, id string) (*entity.Token, error) {
	var d entity.Token
	filter := bson.D{
		{"id", id},
	}
	if err := r.tokenColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.TokenNotFound()
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repo) GetTokenList(ctx context.Context, params QueryParams) ([]*entity.Token, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) UpdateToken(ctx context.Context, token *entity.Token) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) DeleteToken(ctx context.Context, token *entity.Token) error {
	filter := bson.D{{"id", token.ID}}
	_, err := r.tokenColl().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
