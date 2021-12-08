package memory

import (
	"fmt"
	"github.com/ijingjingyang/go-frame/client"
	memCache "github.com/jjonline/go-lib-backend/memory"
	"reflect"
	"time"
)

// !!警告!!
// 本地内存缓存的过期时间并不是精确的过期，purge执行逻辑为每10分钟执行1次，所以不要依赖本地内存缓存的过期时间
// 本地内存缓存的唯一优势就是：速度快，对于不常变动的需要速度的可选用
// !!警告!!

// GetWithSetter 使用设置器setter方式获取本地內存緩存
//  -- 若缓存存在则setter不会执行，若缓存不存在则setter执行并自动设置缓存
//  -- 内存缓存是原始类型直接保存在本地内存中，所以读取出来依然是原始类型，仅仅需要 data.(Type)显式转换一下即可
func GetWithSetter(key string, setter func() (interface{}, error), timeout time.Duration) (data interface{}, err error) {
	// 先从内存缓存中读取
	data = Get(key)
	if nil != data {
		return data, nil
	}

	if setter == nil {
		return nil, fmt.Errorf("setter不得设置为nil，缓存键值key为：%s", key)
	}

	// execute getter
	data, err = setter()
	if nil != err || data == nil || reflect.ValueOf(data).IsZero() {
		return nil, err
	}

	// 直接调用Set方法设置并输出
	return data, Set(key, data, timeout)
}

// Set 设置一个memory本地内存缓存
//  -- 内存缓存是原始类型直接保存在本地内存中，所以读取出来依然是原始类型，仅仅需要 data.(Type)显式转换一下即可
func Set(key string, data interface{}, timeout time.Duration) error {
	// 参数错误：缓存的数据为空
	if data == nil || reflect.ValueOf(data).IsZero() {
		return fmt.Errorf("不得缓存零值数据，缓存键值key为：%s", key)
	}

	// 过期时间设置为0或负数即不过期
	if timeout <= 0 {
		timeout = memCache.NoExpiration
	}
	client.MemoryCache.Set(key, data, timeout)
	return nil
}

// Get memory本地内存缓存通过key读取一个缓存值
//  -- 注意读取后值非nil时需要显式转换成为设置时的类型 data.(Type)
func Get(key string) interface{} {
	data, ok := client.MemoryCache.Get(key)
	if ok {
		return data
	}
	return nil
}

// Del 删除一个本地内存缓存
func Del(key string) {
	client.MemoryCache.Delete(key)
}
