package limiter

import (
	"time"

	"github.com/xxjwxc/public/myredis"
)

type semaphore struct {
	limit       int
	nameSpace   string // 名称空间
	isTsTimeout bool   // 是否启动超时过滤

	redisClient myredis.RedisDial
}

type LimiterIFS interface {
	Acquire(timeout int) (string, error)
	Release(token string) // 释放
	GetTimeDuration(token string) (time.Duration, error)

	Init() // 重新初始化
}

func NewLimiter(ops ...Option) (lifs LimiterIFS) {
	var tmp = semaphore{}
	for _, o := range ops {
		o.apply(&tmp)
	}

	// set default
	if tmp.limit <= 0 {
		tmp.limit = 1
	}
	if len(tmp.nameSpace) == 0 {
		tmp.nameSpace = "gowp/nameSpace"
	}

	//-------------end

	if tmp.redisClient == nil { // cache模式
		//lifs = &limiterCache{tmp}
	} else { // redis sync
		lifs = &limiterRedis{semaphore: tmp}
	}
	lifs.Init()

	return
}
