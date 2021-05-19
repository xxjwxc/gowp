package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/xxjwxc/public/message"
	"github.com/xxjwxc/public/mycache"
	"github.com/xxjwxc/public/myglobal"
	"golang.org/x/sync/semaphore"
)

type limiterCache struct {
	mySemaphore
	*semaphore.Weighted
	queueName       string // hash表中的队列名字
	tokenTsHashName string // 超时hash队列设置
	weighted        int64  // 权重
	cache           *mycache.MyCache
}

func (l *limiterCache) Init() {
	l.queueName = l.nameSpace + "_" + "queue"
	l.tokenTsHashName = l.nameSpace + "_" + "hash"
	l.cache = mycache.NewCache(fmt.Sprintf("_%v_cache", l.nameSpace))
	l.weighted = 1

	l.Weighted = semaphore.NewWeighted(int64(l.limit))

}

func (l *limiterCache) Acquire(timeout int) (string, error) {
	if timeout == 0 { // 不超时
		if l.Weighted.TryAcquire(l.weighted) {
			return l.generateToken(), nil
		}
		return "", message.GetError(message.NotFindError)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	err := l.Weighted.Acquire(ctx, l.weighted)
	if err != nil {
		return "", err
	}

	return l.generateToken(), nil
}

func (l *limiterCache) Release(token string) { // 释放
	l.Weighted.Release(l.weighted)
	if l.isTsTimeout {
		l.cache.Delete(token)
	}
}

func (l *limiterCache) GetTimeDuration(token string) (time.Duration, error) {
	if !l.isTsTimeout {
		return 0, nil
	}

	var ts int64
	err := l.cache.Value(token, &ts)
	if err != nil {
		return 0, err
	}
	return time.Since(time.Unix(ts, 0)), err
}

func (l *limiterCache) generateToken() string {
	str := myglobal.GetNode().GetIDStr()
	if l.isTsTimeout {
		l.cache.Add(str, time.Now().Unix(), 0)
	}
	return str
}
