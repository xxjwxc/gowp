package limiter

import "github.com/xxjwxc/public/myredis"

type Option interface {
	apply(*semaphore)
}

type optionFunc func(*semaphore)

func (f optionFunc) apply(o *semaphore) {
	f(o)
}

// WithLimit 设置最大并发数
func WithLimit(limit int) Option {
	return optionFunc(func(s *semaphore) {
		s.limit = limit
	})
}

// WithNamespace 设置命名空间
func WithNamespace(namespace string) Option {
	return optionFunc(func(s *semaphore) {
		s.nameSpace = namespace
	})
}

// WithRedis 设置默认redis
func WithRedis(redisClient myredis.RedisDial) Option {
	return optionFunc(func(s *semaphore) {
		s.redisClient = redisClient
	})
}

// WithRedis 是否超时记录
func WithTsTimeout(isTsTimeout bool) Option {
	return optionFunc(func(s *semaphore) {
		s.isTsTimeout = isTsTimeout
	})
}
