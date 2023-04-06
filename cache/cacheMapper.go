package cache

import "time"

type Cache interface {
	//size:1KB 100KB 1MB 2MB 1GB
	SetMaxMemory(size string) bool
	//将value写入缓存
	Set(key string, val interface{}, expire time.Duration) bool
	//根据Key值获取value
	Get(key string) (interface{}, bool)
	//删除Key值
	Del(key string) bool
	//判断Key值是否存在
	Exists(key string) bool
	//清空所有的Key
	Flush() bool
	//获取缓存中所有key的数量
	Keys() int64
}
