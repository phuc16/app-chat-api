package repository

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) chatColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("conversations")
}

func (r *Repo) CreateChatIndexes(ctx context.Context) (res []string, err error) {
	indexes, err := r.chatColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
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
	_, err = r.chatColl().UpdateOne(ctx, bson.M{"id": conservation.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.TokenExists()
		}
		return err
	}

	for _, userId := range conservation.ListUser {
		err = r.AddNewConversationToUser(ctx, userId, conservation.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) GetConversationById(ctx context.Context, id string) (res *entity.Conversation, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)
	var d entity.Conversation
	filter := bson.D{
		{"id", id},
	}
	if err := r.chatColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.TokenNotFound()
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repo) GetListUserInConversation(ctx context.Context, conversationId string) (res []string, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)
	var d entity.Conversation
	filter := bson.D{
		{"id", conversationId},
	}
	if err := r.chatColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.TokenNotFound()
		}
		return nil, err
	}
	return d.ListUser, nil
}

func (r *Repo) AddNewChatToConversation(ctx context.Context, chat *entity.Chat) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)
	conversation, err := r.GetConversationById(ctx, chat.ToConversationId)
	if err != nil {
		return err
	}
	chat.Timestamp = time.Now()

	conversation.Chat = append(conversation.Chat, *chat)

	filter := bson.D{{"id", chat.ToConversationId}}
	update := bson.D{
		{"$set", bson.M{"chat": conversation.Chat}},
	}

	_, err = r.chatColl().UpdateOne(ctx, filter, update)
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
