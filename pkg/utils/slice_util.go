package utils

import (
	"reflect"
	"strings"
	"time"
)

func IsExistTime(slice []time.Time, search time.Time) bool {
	for _, t := range slice {
		if t.Unix() == search.Unix() {
			return true
		}
	}
	return false
}

func IsExist(slice []string, search string) bool {
	for _, s := range slice {
		if s == search {
			return true
		}
	}
	return false
}

func FindString(slice []string, search string) int {
	for i, s := range slice {
		if s == search {
			return i
		}
	}
	return -1
}

func FindDiff(slice1 []string, slice2 []string) []string {
	result := make([]string, 0)
	for _, item := range slice1 {
		if !IsExist(slice2, item) {
			result = append(result, item)
		}
	}
	return result
}

func GetTagList(u interface{}) []string {
	res := []string{}
	result := make(map[string]string)
	value := reflect.ValueOf(u)
	getTypeTags(value.Type(), "", result)
	for _, v := range result {
		res = append(res, v)
	}
	return res
}

func getTypeTags(t reflect.Type, prefix string, tags map[string]string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		pre := ""
		if prefix != "" {
			pre = prefix + "."
		}
		if tag != "" {
			tag = strings.ReplaceAll(tag, ",omitempty", "")
			tags[pre+field.Name] = pre + tag
		}
		if field.Type.Kind() == reflect.Struct {
			getTypeTags(field.Type, pre+tag, tags)
		}
	}
}

func RemoveItemInSlice(slice []string, item string) []string {
	result := make([]string, 0)
	for i := 0; i < len(slice); i++ {
		if slice[i] != item {
			result = append(result, slice[i])
		}
	}
	return result
}
