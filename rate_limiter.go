package amo

import "time"

type RateLimiter interface {
	WaitForRequest()
}

var defaultRTLimiter = newDefaultRateLimiter()

type empty struct{}

type defaultRateLimiter struct {
	limiter chan empty
}

func newDefaultRateLimiter() *defaultRateLimiter {
	limiter := &defaultRateLimiter{}
	limiter.limiter = make(chan empty)
	go func() {
		for {
			limiter.limiter <- empty{}
			time.Sleep(time.Second)
		}
	}()
	return limiter
}

func (rl *defaultRateLimiter) WaitForRequest() {
	<-rl.limiter
}
