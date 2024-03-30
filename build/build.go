package build

import "encoding/json"

var AppName string = "app-chat-api"
var Version string = "1.0.0"
var CommitID string
var Date string
var User string

type info struct {
	AppName  string `json:"app_name"`
	Version  string `json:"version"`
	CommitID string `json:"commit_id"`
	Date     string `json:"date"`
	User     string `json:"user"`
}

func (i info) String() string {
	bytes, _ := json.Marshal(i)
	return string(bytes)
}

func Info() info {
	info := info{
		AppName:  AppName,
		Version:  Version,
		CommitID: CommitID,
		Date:     Date,
		User:     User,
	}
	return info
}
