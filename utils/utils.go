package util

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
)

// 返回默认值时，用error提示
func Parse(size string) (int64, string, error) {
	rep, _ := regexp.Compile("[0-9]+")
	token := string(rep.ReplaceAll([]byte(size), []byte("")))
	num, _ := strconv.ParseInt(strings.Replace(size, token, "", 1), 10, 64)
	var size_num int64
	switch token {
	case "B":
		size_num = num * B
	case "KB":
		size_num = num * KB
	case "MB":
		size_num = num * MB
	case "GB":
		size_num = num * GB
	case "TB":
		size_num = num * TB
	default:
		size_num = 0
	}
	if size_num == 0 {
		size_str := "100MB"
		size_num = 100 * MB
		err := errors.New("输入出错，默认设置最大缓存为100MB")
		return size_num, size_str, err
	}
	size_str := strconv.FormatInt(num, 10) + token
	return size_num, size_str, nil
}

// 获取内存占用
func GetSize(val interface{}) int64 {
	// value := reflect.ValueOf(val)
	js, _ := json.Marshal(val)
	size := len(js)
	// fmt.Println(size)
	return int64(size)
}
