package repository

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) otpColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("otps")
}

func (r *Repo) CreateOtpIndexes(ctx context.Context) ([]string, error) {
	indexes, err := r.otpColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{"email", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"code", 1},
		}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func (r *Repo) SaveOtp(ctx context.Context, otp *entity.Otp) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", otp},
	}
	_, err = r.otpColl().UpdateOne(ctx, bson.M{"id": otp.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.OtpExists()
		}
		return err
	}
	return nil
}

func (r *Repo) GetOtp(ctx context.Context, otp *entity.Otp) (res *entity.Otp, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d entity.Otp
	filter := bson.D{{"$and", []interface{}{
		bson.D{{"email", otp.Email}},
		bson.D{{"code", otp.Code}},
	}}}
	if err := r.otpColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.OtpNotFound()
		}
		return nil, err
	}
	return &d, nil
}
func (r *Repo) DeleteOtp(ctx context.Context, otp *entity.Otp) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", otp.ID}}
	_, err = r.otpColl().DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
