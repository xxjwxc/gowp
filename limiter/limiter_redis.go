package limiter

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

type limiterRedis struct {
	mySemaphore
	queueName       string // hash表中的队列名字
	lockName        string // 初始化的时候只允许单个初始化，(锁住当前初始化)
	tokenTsHashName string // 超时hash队列设置
}

// Init 初始化token
func (l *limiterRedis) Init() {
	l.queueName = l.nameSpace + "_" + "queue"
	l.lockName = l.nameSpace + "_" + "lock"
	l.tokenTsHashName = l.nameSpace + "_" + "hash"

	ok, _ := l.tryLock(0)
	if !ok {
		fmt.Println("lock failed")
		return
	}

	// clean old token list
	// to do: pipeline
	l.redisClient.Do("DEL", l.queueName)
	for i := 1; i <= l.limit; i++ {
		l.push(fmt.Sprintf("token_seq_%v", i))
		// l.Tokens = append(l.Tokens, tmp_token)
	}
}

// Acquire 获取一个
func (l *limiterRedis) Acquire(timeout int) (string, error) {
	var token string
	var err error

	//l.ScanTimeoutToken()
	if timeout > 0 {
		token, err = l.popBlock(timeout)
	} else {
		token, err = l.pop()
	}
	return token, err
}

// Release 释放一个
func (l *limiterRedis) Release(token string) {
	l.push(token)
}

// GetTimeDuration 获取已超时时间
func (l *limiterRedis) GetTimeDuration(token string) (time.Duration, error) {
	if !l.isTsTimeout {
		return 0, nil
	}

	res, err := redis.String(l.redisClient.Do("HGET", l.tokenTsHashName, token))
	if err != nil {
		return 0, err
	}
	ts, _ := strconv.Atoi(res)
	return time.Since(time.Unix(int64(ts), 0)), err
}

func (l *limiterRedis) tryLock(timeout int) (bool, error) {
	var err error

	if timeout == 0 {
		_, err = l.redisClient.Do("SET", l.lockName, "locked", "NX")
	} else {
		_, err = l.redisClient.Do("SET", l.lockName, "locked", "EX", 30, "NX")
	}

	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

// 还回去
func (l *limiterRedis) push(body string) (int, error) {
	// to do: pipeline
	res, err := redis.Int(l.redisClient.Do("RPUSH", l.queueName, body)) // 还回去
	if l.isTsTimeout {
		l.redisClient.Do("HDEL", l.tokenTsHashName, body) // 删除超时
	}

	return res, err
}

// pop 获取一个
func (l *limiterRedis) pop() (string, error) {
	res, err := redis.String(l.redisClient.Do("LPOP", l.queueName))
	// 允许队列为空值
	if err == redis.ErrNil {
		err = nil
	}
	if l.isTsTimeout {
		l.redisClient.Do("HSET", l.tokenTsHashName, res, time.Now().Unix())
	}

	return res, err
}

func (l *limiterRedis) popBlock(timeout int) (string, error) {
	// refer: https://www.runoob.com/redis/lists-blpop.html
	res_map, err := redis.StringMap(l.redisClient.Do("BLPOP", l.queueName, timeout))
	// 允许队列为空值
	if err == redis.ErrNil {
		err = nil
	} else if err != nil {
		fmt.Println(err)
	}

	res, ok := res_map[l.queueName]
	if ok && l.isTsTimeout {
		l.redisClient.Do("HSET", l.tokenTsHashName, res, time.Now().Unix())
	}

	// if !ok {
	// 	return "", err
	// }

	return res, err
}

// func (l *limiterRedis) ScanTimeoutToken() []string {
// 	expire_tokens := []string{}

// 	l.ScanLock.Lock()

// 	if !l.ScanIsContinue() {
// 		l.ScanLock.Unlock()
// 		return expire_tokens
// 	}

// 	res, _ := redis.StringMap(l.redisClient.Do("HGETALL", l.TokenTsHashName))

// 	for token, ts_s := range res {
// 		ts, _ := strconv.Atoi(ts_s)
// 		diff_ts := time.Since(time.Unix(int64(ts), 0)) //time.Now().Sub(time.Unix(int64(ts), 0))
// 		if int(diff_ts.Seconds()) > l.ScanTimeout {
// 			expire_tokens = append(expire_tokens, token)
// 		}
// 	}
// 	l.LastScanTs = time.Now()
// 	l.ScanLock.Unlock()

// 	return expire_tokens
// }

// func (l *limiterRedis) ScanIsContinue() bool {
// 	now := time.Now()
// 	return int(now.Sub(l.LastScanTs).Seconds()) >= l.ScanInterval
// }
