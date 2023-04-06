package cache

import (
	util "go_study/Simple_mycache/utils"
	"log"
	"sync"
	"time"
)

type CacheManger struct {
	//最大缓存
	maxMemory int64
	//最大缓存（字符串）
	maxMemorystr string
	//数据
	val map[string]CacheConsumer
	//当前缓存
	memory int64
	//读写锁
	lock sync.RWMutex
}

type CacheConsumer struct {
	value interface{}
	//过期时间
	expire time.Time
	//过期时长
	expireDur time.Duration
	//占用缓存
	size int64
}

func NewMemCache() *CacheManger {
	cm := &CacheManger{
		maxMemorystr: "100MB",
		maxMemory:    100 * util.MB,
		val:          make(map[string]CacheConsumer, 0),
	}
	go cm.TickerClear()
	return cm
}

//size:1KB 100KB 1MB 2MB 1GB
func (cm *CacheManger) SetMaxMemory(size string) bool {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	var err error
	//解析size。且在输入不规范时，为了不panic，给最大缓存赋默认值
	cm.maxMemory, cm.maxMemorystr, err = util.Parse(size)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// 将value写入缓存
func (cm *CacheManger) Set(key string, val interface{}, expire time.Duration) bool {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.dele(key)
	cm.add(key, val, expire)

	return true
}

// 根据Key值获取value
func (cm *CacheManger) Get(key string) (interface{}, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	val, ok := cm.val[key]
	if !ok {
		log.Println("缓存中不存在", key)
		return nil, false
	}
	if time.Now().After(val.expire) {
		log.Println(key, ":数据已过期")
		cm.dele(key)
		return nil, false
	}
	return val.value, ok
}

// 删除Key值
func (cm *CacheManger) Del(key string) bool {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	ok := cm.dele(key)
	return ok
}

// 判断Key值是否存在
func (cm *CacheManger) Exists(key string) bool {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	_, ok := cm.val[key]
	return ok
}

// 清空所有的Key
func (cm *CacheManger) Flush() bool {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.val = make(map[string]CacheConsumer, 0)
	cm.memory = 0
	return true
}

// 获取缓存中所有key的数量
func (cm *CacheManger) Keys() int64 {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	size := len(cm.val)
	return int64(size)
}

// 写入缓存
func (cm *CacheManger) add(key string, val interface{}, expire time.Duration) bool {
	Consumer := &CacheConsumer{
		value:     val,
		expireDur: expire,
		expire:    time.Now().Add(expire),
		size:      util.GetSize(val),
	}
	cm.memory += Consumer.size
	//超缓存，清理过期缓存
	if cm.memory >= cm.maxMemory {
		cm.ClearExpire()
		//仍然超缓存，按淘汰策略
		if cm.memory >= cm.maxMemory {
			// ok := cm.ZipCache(Consumer.size)
			// if !ok {
			log.Println("缓存写入失败，缓存信息：", Consumer)
			return false
			// }
		}
	}
	cm.val[key] = *Consumer
	return true
}

// 移除缓存
func (cm *CacheManger) dele(key string) bool {
	delete(cm.val, key)
	cm.memory -= cm.val[key].size
	return true
}

// 清理过期缓存
func (cm *CacheManger) ClearExpire() bool {
	for k, Consunmer := range cm.val {
		if time.Now().After(Consunmer.expire) {
			cm.dele(k)
		}
	}
	return true
}

// 按过期时间淘汰缓存
func (cm *CacheManger) ZipCache(size int64) bool {

	return true
}

// 周期性淘汰过期缓存
func (cm *CacheManger) TickerClear() {
	T := time.Tick(time.Minute * 20)
	for {
		select {
		case <-T:
			cm.lock.Lock()
			cm.ClearExpire()
			cm.lock.Unlock()
		}
	}
}
