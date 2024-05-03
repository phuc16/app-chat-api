package repository

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) conversationColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("conversations")
}

func (r *Repo) CreateChatIndexes(ctx context.Context) (res []string, err error) {
	indexes, err := r.conversationColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{
			{"id", 1},
		}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func (r *Repo) NewConversation(ctx context.Context, conservation *entity.Conversation) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)
	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", conservation},
	}
	_, err = r.conversationColl().UpdateOne(ctx, bson.M{"id": conservation.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.DuplicateConversationId()
		}
		return err
	}
	return nil
}

func (r *Repo) GetConversationById(ctx context.Context, id string) (res *entity.Conversation, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)
	var d []*entity.Conversation
	pipeLine := mongo.Pipeline{}
	pipeLine = append(pipeLine, matchFieldPipeline("id", id))
	pipeLine = append(pipeLine, conversationUsersLookupPipeline)

	cursor, err := r.conversationColl().Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return res, err
	}
	if err = cursor.All(ctx, &d); err != nil {
		return res, err
	}
	if len(d) == 0 {
		return res, errors.ConversationNotFound()
	}
	return d[0], nil
}

func (r *Repo) GetChatByConversationId(ctx context.Context, conversationId string, params *QueryParams) (res []entity.Chat, total int64, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	coll := r.conversationColl()

	pipeLine := mongo.Pipeline{}
	pipeLine = append(pipeLine, matchFieldPipeline("id", conversationId))

	cursor, err := coll.Aggregate(ctx, pipeLine, collationAggregateOption)
	if err != nil {
		return res, 0, err
	}
	conversation := []*entity.Conversation{}
	if err = cursor.All(ctx, &conversation); err != nil {
		return res, 0, err
	}
	if len(conversation) == 0 {
		return res, 0, errors.ConversationNotFound()
	}
	res, total = getMatchingValues(conversation[0].Chat, params)
	return
}

func (r *Repo) GetListIDUserInConversation(ctx context.Context, conversationId string) (res []string, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)
	var d entity.Conversation
	filter := bson.D{
		{"id", conversationId},
	}
	if err := r.conversationColl().FindOne(ctx, filter).Decode(&d); err != nil {
		return nil, errors.CanNotGetListIDUserInConversation()
	}
	return d.ListUser, nil
}

func (r *Repo) AddNewChatToConversation(ctx context.Context, chat *entity.Chat) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)
	chat.Timestamp = time.Now()
	chat.ID = utils.NewID()

	filter := bson.D{{"id", chat.ToConversationId}}
	update := bson.M{"$addToSet": bson.M{"chat": chat}}

	_, err = r.conversationColl().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) AddNewConversationToUser(ctx context.Context, userID string, conversationID string) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)
	filter := bson.D{{"id", userID}}
	update := bson.M{"$addToSet": bson.M{"conversation_ids": conversationID}}
	_, err = r.userColl().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateMessage(ctx context.Context, conversation *entity.Conversation) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	filter := bson.D{{"id", conversation.ID}, {"chat.id", conversation.Chat[0].ID}}
	update := bson.M{"$set": bson.M{"chat.$.msg": conversation.Chat[0].Msg}}

	_, err = r.conversationColl().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func getMatchingValues(slice []entity.Chat, params *QueryParams) (res []entity.Chat, total int64) {
	pattern := regexp.QuoteMeta(params.Search)

	pattern = ".*" + pattern + ".*"

	regex := regexp.MustCompile(pattern)

	for _, value := range slice {
		if regex.MatchString(value.Msg.(string)) {
			res = append(res, value)
		}
	}

	total = int64(len(res))

	if params.SortType == -1 {
		length := len(res)
		for i := 0; i < length/2; i++ {
			res[i], res[length-i-1] = res[length-i-1], res[i]
		}
	}

	endIndex := params.Skip + params.Limit
	if endIndex > int64(len(res)) {
		endIndex = int64(len(res))
	}

	return res[params.Skip:endIndex], total
}
