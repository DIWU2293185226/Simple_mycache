# Simple_mycache
Simple_mycache实现了一个简易的内存缓存系统  
支持对缓存基本的设置、刷新、删除和读取的操作  
可以粗略的实现对缓存系统的内存控制，但出于现代编程语言自身的特性，内存控制在精度上有一定损失  
（学习go语言过程中的造物.jpg）  
# 包中可供使用的基本接口如下：
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
