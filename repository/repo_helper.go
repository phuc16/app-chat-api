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

var creatorUnwindPipeline = bson.D{{"$unwind", bson.M{
	"path":                       "$creator",
	"preserveNullAndEmptyArrays": true,
}}}
var creatorLookupPipeline = bson.D{
	{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "creator_id"},
		{"foreignField", "id"},
		{"as", "creator"},
	}},
}
var updaterUnwindPipeline = bson.D{{"$unwind", bson.M{
	"path":                       "$updater",
	"preserveNullAndEmptyArrays": true,
}}}
var updaterLookupPipeline = bson.D{
	{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "updater_id"},
		{"foreignField", "id"},
		{"as", "updater"},
	}},
}

var reviewerUnwindPipeline = bson.D{{"$unwind", bson.M{
	"path":                       "$reviewer",
	"preserveNullAndEmptyArrays": true,
}}}
var reviewerLookupPipeline = bson.D{
	{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "reviewer_id"},
		{"foreignField", "id"},
		{"as", "reviewer"},
	}},
}
var approverUnwindPipeline = bson.D{{"$unwind", bson.M{
	"path":                       "$approver",
	"preserveNullAndEmptyArrays": true,
}}}
var approverLookupPipeline = bson.D{
	{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "approver_id"},
		{"foreignField", "id"},
		{"as", "approver"},
	}},
}

var subTicketsUnwindPipeline = bson.D{{"$unwind", bson.M{
	"path":                       "$sub_tickets",
	"preserveNullAndEmptyArrays": true,
}}}
var subTicketsLookupPipeline = bson.D{
	{"$lookup", bson.D{
		{"from", "tickets"},
		{"localField", "sub_ticket_ids"},
		{"foreignField", "id"},
		{"as", "sub_tickets"},
	}},
}

var ruleConnectorsPipeline = []bson.D{
	bson.D{{"$unwind", bson.D{
		{"path", "$connectors"},
		{"preserveNullAndEmptyArrays", true},
	}}},
	bson.D{
		{"$lookup", bson.D{
			{"from", "connectors"},
			{"localField", "connectors.connector_id"},
			{"foreignField", "id"},
			{"as", "connectors.connector"},
		}},
	},
	bson.D{{"$unwind", bson.D{
		{"path", "$connectors.connector"},
		{"preserveNullAndEmptyArrays", true},
	}}},
	bson.D{
		{"$group", bson.D{
			{"_id", "$_id"},
			{"id", bson.D{{"$first", "$id"}}},
			{"name", bson.D{{"$first", "$name"}}},
			{"active", bson.D{{"$first", "$active"}}},
			{"status", bson.D{{"$first", "$status"}}},
			{"content", bson.D{{"$first", "$content"}}},
			{"tag", bson.D{{"$first", "$tag"}}},
			{"connectors", bson.D{{"$push", "$connectors"}}},
			{"creator_id", bson.D{{"$first", "$creator_id"}}},
			{"updater_id", bson.D{{"$first", "$updater_id"}}},
			{"creator", bson.D{{"$first", "$creator"}}},
			{"updater", bson.D{{"$first", "$updater"}}},
			{"created_at", bson.D{{"$first", "$created_at"}}},
			{"updated_at", bson.D{{"$first", "$updated_at"}}},
		}},
	},
}
