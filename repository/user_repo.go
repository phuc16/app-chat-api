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

func (r *Repo) userColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("users")
}

func (r *Repo) CreateUserIndexes(ctx context.Context) ([]string, error) {
	indexes, err := r.userColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{"id", 1},
			{"deleted_at", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"username", 1},
			{"deleted_at", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"email", 1},
			{"deleted_at", 1},
		}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{
			{"name", "text"},
			{"username", "text"},
			{"email", "text"},
		}},
	})
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func (r *Repo) SaveUser(ctx context.Context, user *entity.User) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", user},
	}
	_, err = r.userColl().UpdateOne(ctx, bson.M{"id": user.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.UserExists()
		}
		return err
	}
	return nil
}

func (r *Repo) GetUserById(ctx context.Context, id string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d []*entity.User

	pipeLine := mongo.Pipeline{}
	pipeLine = append(pipeLine, matchFieldPipeline("id", id))
	pipeLine = append(pipeLine, matchFieldPipeline("deleted_at", nil))
	pipeLine = append(pipeLine, limitPipeline(1))
	pipeLine = append(pipeLine, friendsLookupPipeline)
	pipeLine = append(pipeLine, friendRequestsLookupPipeline)
	// pipeLine = append(pipeLine, friendsUnwindPipeline)

	cursor, err := r.userColl().Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &d); err != nil {
		return nil, err
	}
	if len(d) <= 0 {
		return nil, errors.UserNotFound()
	}
	if !d[0].IsActive {
		return nil, errors.UserInactive()
	}
	return d[0], nil
}
func (r *Repo) GetUserByEmail(ctx context.Context, email string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d []*entity.User

	pipeLine := mongo.Pipeline{}
	pipeLine = append(pipeLine, matchFieldPipeline("email", email))
	pipeLine = append(pipeLine, matchFieldPipeline("deleted_at", nil))
	pipeLine = append(pipeLine, limitPipeline(1))
	pipeLine = append(pipeLine, friendsLookupPipeline)
	pipeLine = append(pipeLine, friendRequestsLookupPipeline)
	// pipeLine = append(pipeLine, friendsUnwindPipeline)

	cursor, err := r.userColl().Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &d); err != nil {
		return nil, err
	}
	if len(d) <= 0 {
		return nil, errors.UserNotFound()
	}
	if !d[0].IsActive {
		return nil, errors.UserInactive()
	}
	return d[0], nil
}

func (r *Repo) GetUserByUserName(ctx context.Context, username string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d []*entity.User

	pipeLine := mongo.Pipeline{}
	pipeLine = append(pipeLine, matchFieldPipeline("username", username))
	pipeLine = append(pipeLine, matchFieldPipeline("deleted_at", nil))
	pipeLine = append(pipeLine, limitPipeline(1))
	pipeLine = append(pipeLine, friendsLookupPipeline)
	pipeLine = append(pipeLine, friendRequestsLookupPipeline)
	// pipeLine = append(pipeLine, friendsUnwindPipeline)

	cursor, err := r.userColl().Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &d); err != nil {
		return nil, err
	}
	if len(d) <= 0 {
		return nil, errors.UserNotFound()
	}
	if !d[0].IsActive {
		return nil, errors.UserInactive()
	}
	return d[0], nil
}

func (r *Repo) GetUserByUserNameOrEmail(ctx context.Context, username string, email string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d entity.User
	filter := bson.D{{"$or", []interface{}{
		bson.D{{"username", username}, {"deleted_at", nil}},
		bson.D{{"email", email}, {"deleted_at", nil}},
	}}}
	if err := r.userColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.UserNotFound()
		}
		return nil, err
	}
	if !d.IsActive {
		return nil, errors.UserInactive()
	}
	return &d, nil
}

func (r *Repo) GetInactiveUser(ctx context.Context, email string) (res *entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d entity.User
	filter := bson.D{{"$and", []interface{}{
		bson.D{{"email", email}},
		bson.D{{"is_active", false}},
		bson.D{{"deleted_at", nil}},
	}}}
	if err := r.userColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.UserNotFound()
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repo) CheckUserNameAndEmailExist(ctx context.Context, username string, email string) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d entity.User
	filter := bson.D{
		{"username", username},
		{"deleted_at", nil},
	}
	if err2 := r.userColl().FindOne(ctx, filter).Decode(&d); err2 != nil {
		if err2 != mongo.ErrNoDocuments {
			return err2
		}
	}
	if d.ID != "" {
		return errors.UserNameExists()
	}
	filter = bson.D{
		{"email", email},
		{"deleted_at", nil},
	}
	if err2 := r.userColl().FindOne(ctx, filter).Decode(&d); err2 != nil {
		if err2 != mongo.ErrNoDocuments {
			return err2
		}
	}
	if d.ID != "" {
		return errors.UserEmailExists()
	}
	return
}

func (r *Repo) CheckDuplicateUserNameAndEmail(ctx context.Context, user *entity.User, username string, email string) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var d entity.User
	filter := bson.D{{"$and", []interface{}{
		bson.D{{"id", bson.M{"$ne": user.ID}}},
		bson.D{{"username", username}},
		bson.D{{"deleted_at", nil}},
	}}}
	if err2 := r.userColl().FindOne(ctx, filter).Decode(&d); err2 != nil {
		if err2 != mongo.ErrNoDocuments {
			return err2
		}
	}
	if d.ID != "" {
		return errors.UserNameExists()
	}
	filter = bson.D{{"$and", []interface{}{
		bson.D{{"id", bson.M{"$ne": user.ID}}},
		bson.D{{"email", email}},
		bson.D{{"deleted_at", nil}},
	}}}
	if err2 := r.userColl().FindOne(ctx, filter).Decode(&d); err2 != nil {
		if err2 != mongo.ErrNoDocuments {
			return err2
		}
	}
	if d.ID != "" {
		return errors.UserEmailExists()
	}
	return
}

func (r *Repo) GetUserList(ctx context.Context, params *QueryParams) (res []*entity.User, total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	coll := r.userColl()

	pipeLine := mongo.Pipeline{}
	if params.Search != "" {
		pipeLine = append(pipeLine, partialMatchingSearchPipeline([]string{"name", "username", "email"}, params.Search)...)
	}
	for k, v := range params.Filter {
		pipeLine = append(pipeLine, matchFieldPipeline(k, v))
	}
	pipeLine = append(pipeLine, matchFieldPipeline("deleted_at", nil))
	pipeLine = append(pipeLine, matchFieldPipeline("is_active", true))

	cursor, err := coll.Aggregate(ctx, append(pipeLine, bson.D{{"$count", "total_count"}}))
	if err != nil {
		return nil, 0, err
	}
	result := bson.M{}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, 0, err
		}
	}
	totalCount, _ := result["total_count"].(int32)
	total = int64(totalCount)

	pipeLine = append(pipeLine, params.SkipLimitSortPipeline()...)

	cursor, err = coll.Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return res, 0, err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return res, 0, err
	}
	return res, total, nil
}

func (r *Repo) GetAllUsers(ctx context.Context) (res []*entity.User, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	coll := r.userColl()

	pipeLine := mongo.Pipeline{}
	pipeLine = append(pipeLine, matchFieldPipeline("deleted_at", nil))
	pipeLine = append(pipeLine, matchFieldPipeline("is_active", true))
	pipeLine = append(pipeLine, friendsLookupPipeline)

	cursor, err := coll.Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return res, err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return res, err
	}
	return res, nil
}

func (r *Repo) UpdateUser(ctx context.Context, user *entity.User) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", user.ID}}
	update := bson.M{"$set": user}
	_, err = r.userColl().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteUser(ctx context.Context, user *entity.User) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", user.ID}}
	update := bson.M{"$set": user}
	_, err = r.userColl().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) CountUser(ctx context.Context) (total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	var filter = filterField("deleted_at", nil)
	total, err = r.userColl().CountDocuments(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (r *Repo) AddFriendRequest(ctx context.Context, user *entity.User, friend *entity.User) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", user.ID}}
	update := bson.M{"$addToSet": bson.M{"friend_request_ids": friend.ID}}
	_, err = r.userColl().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) RemoveFriendRequest(ctx context.Context, user *entity.User, friend *entity.User) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", user.ID}}
	update := bson.M{"$pull": bson.M{"friend_request_ids": friend.ID}}
	_, err = r.userColl().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) AddFriend(ctx context.Context, user *entity.User, friend *entity.User) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", user.ID}}
	update := bson.M{"$addToSet": bson.M{"friend_ids": friend.ID}}
	_, err = r.userColl().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) RemoveFriend(ctx context.Context, user *entity.User, friend *entity.User) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", user.ID}}
	update := bson.M{"$pull": bson.M{"friend_ids": friend.ID}}
	_, err = r.userColl().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) AddNewConversationToUser(ctx context.Context, userID string) (err error) {
	// TODO: Implement this
	return nil
}
