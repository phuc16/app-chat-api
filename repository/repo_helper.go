package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryParams struct {
	Filter    map[string]any
	Limit     int64
	Skip      int64
	SortField string
	SortType  int
	Search    string
}

func (p QueryParams) SkipLimitSortPipeline() (pipeline mongo.Pipeline) {
	if p.SortField != "" {
		pipeline = append(pipeline, sortPipeline(p.SortField, p.SortType))
	}
	pipeline = append(pipeline, mongo.Pipeline{skipPipeline(p.Skip), limitPipeline(p.Limit)}...)
	return pipeline
}

func (p QueryParams) SkipLimitPipeline() mongo.Pipeline {
	pipeline := mongo.Pipeline{skipPipeline(p.Skip), limitPipeline(p.Limit)}
	return pipeline
}

var sortPipeline = func(sortField string, sortType int) bson.D { return bson.D{{"$sort", bson.M{sortField: sortType}}} }
var skipPipeline = func(skip int64) bson.D { return bson.D{{"$skip", skip}} }
var limitPipeline = func(limit int64) bson.D { return bson.D{{"$limit", limit}} }

var matchPipeline = func(value interface{}) bson.D {
	return bson.D{{"$match", value}}
}
var matchRegexFieldPipeline = func(field string, value string) bson.D {
	return bson.D{{"$match", bson.M{field: primitive.Regex{
		Pattern: value,
		Options: "g",
	}}}}
}
var matchTextSearchPipeline = func(value string) bson.D {
	return bson.D{{"$match", bson.M{
		"$text": bson.M{
			"$search": value,
		},
	}}}
}
var matchFieldPipeline = func(field string, value interface{}) bson.D {
	return bson.D{{"$match", bson.M{field: value}}}
}
var matchInListPipeline = func(fieldName string, list []string) bson.D {
	return bson.D{{"$match", bson.M{fieldName: bson.M{"$in": list}}}}
}

var collationAggregateOption = &options.AggregateOptions{
	Collation: &options.Collation{
		Locale: "en",
	},
}

var filterField = func(field string, value interface{}) bson.D {
	return bson.D{{field, value}}
}
var filterRegexField = func(field string, value string) bson.D {
	return bson.D{{field, bson.M{"$regex": primitive.Regex{
		Pattern: value,
		Options: "g",
	}}}}
}

var textSearchPipeline = func(value string) bson.D {
	return bson.D{{"$match", bson.M{
		"$text": bson.M{
			"$search": value,
		},
	}}}
}

var partialMatchingSearchPipeline = func(fields []string, value string) []bson.D {
	var pipeline mongo.Pipeline

	match := bson.A{}

	for _, field := range fields {
		matchStage := bson.D{{field, bson.M{"$regex": primitive.Regex{Pattern: value, Options: "i"}}}}
		match = append(match, matchStage)
	}

	matchStage := bson.D{{"$match", bson.D{{"$or", match}}}}
	pipeline = append(pipeline, matchStage)
	return pipeline
}

var friendsUnwindPipeline = bson.D{{"$unwind", bson.M{
	"path":                       "$friends",
	"preserveNullAndEmptyArrays": true,
}}}
var friendsLookupPipeline = bson.D{
	{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "friend_ids"},
		{"foreignField", "id"},
		{"as", "friends"},
	}},
}

var friendRequestsUnwindPipeline = bson.D{{"$unwind", bson.M{
	"path":                       "$friend_requests",
	"preserveNullAndEmptyArrays": true,
}}}
var friendRequestsLookupPipeline = bson.D{
	{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "friend_request_ids"},
		{"foreignField", "id"},
		{"as", "friend_requests"},
	}},
}
