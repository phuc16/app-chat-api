package utils

import (
	"fmt"
	"runtime"
)

func GetFuncName(v ...int) string {
	skip := 1
	if len(v) > 0 {
		skip = v[0]
	}
	pc, _, _, _ := runtime.Caller(skip)
	return fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
}

func GetCurrentFuncName() string {
	return GetFuncName(2)
}
