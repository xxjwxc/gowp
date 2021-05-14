package limiter

type limiterCache struct {
	semaphore
}

func (l *limiterCache) Init() {

}
