package utils

import "encoding/json"

func JSONString(j any) string {
	marshal, _ := json.Marshal(j)
	return string(marshal)
}

func JSONPretty(j any) string {
	marshal, _ := json.MarshalIndent(j, "", "\t")
	return string(marshal)
}
