package apiConn

import (
	"context"
	"fmt"
	"sort"
	"time"

	"golang.org/x/time/rate"
)

// type APIConn struct {
// 	rateLimiter *rate.Limiter
// }

// func Open() *APIConn {
// 	return &APIConn{
// 		rateLimiter: rate.NewLimiter(rate.Every(time.Second/4), 4),
// 	}
// }

// func (c *APIConn) Read(ctx context.Context) (string, error) {
// 	if err := c.rateLimiter.Wait(ctx); err != nil {
// 		return "", err
// 	}

// 	//DoWork
// 	return "Read", nil
// }

// func (c *APIConn) Resolve(ctx context.Context) error {
// 	if err := c.rateLimiter.Wait(ctx); err != nil {
// 		return err
// 	}

// 	//DoWork
// 	return nil
// }

type APIConn struct {
	apiLimit,
	dbLimit RateLimiter
}

type RateLimiter interface {
	Allow() bool
	Wait(context.Context) error
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimiter
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}

	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

func Open() *APIConn {
	return &APIConn{
		apiLimit: MultiLimiter(
			rate.NewLimiter(Per(1, 3*time.Second), 1),
			rate.NewLimiter(Per(1, time.Minute), 5),
		),
		dbLimit: MultiLimiter(
			rate.NewLimiter(rate.Every(time.Second*5), 1),
		),
	}
}

func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

func (l *multiLimiter) Allow() bool {
	for _, l := range l.limiters {
		if !l.Allow() {
			return false
		}
	}

	return true
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

func (c *APIConn) ReadAllow() (string, bool) {
	if !c.dbLimit.Allow() {
		return "", false
	}

	return "Read", true
}

func (c *APIConn) Read(ctx context.Context) (string, error) {
	fmt.Printf("%v Limit : %v\n", time.Now().Format("15:04:05"), c.apiLimit.Limit())
	if err := c.dbLimit.Wait(ctx); err != nil {
		return "", err
	}

	//DoWork
	return "Read", nil
}

func (c *APIConn) Resolve(ctx context.Context) error {
	if err := c.apiLimit.Wait(ctx); err != nil {
		return err
	}

	//DoWork
	return nil
}
