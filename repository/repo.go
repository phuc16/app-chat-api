package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewRepo(db *mongo.Client) *Repo {
	return &Repo{db: db}
}

type Repo struct {
	db *mongo.Client
}

func (r *Repo) ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error) {
	session, err := r.db.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	res, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		return fn(ctx)
	})
	if err != nil {
		return nil, err
	}
	return res, err
}

func (r *Repo) InitIndex(ctx context.Context) (err error) {
	_, err = r.CreateUserIndexes(ctx)
	if err != nil {
		return err
	}
	_, err = r.CreateTokenIndexes(ctx)
	if err != nil {
		return err
	}
	_, err = r.CreateOtpIndexes(ctx)
	if err != nil {
		return err
	}
	return
}
