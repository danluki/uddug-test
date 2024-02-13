package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu sync.RWMutex

	visitors map[string]*visitor
	limit    rate.Limit
	burst    int
	ttl      time.Duration
}

func NewRateLimiter(rps float64, burst int, ttl time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    rate.Limit(rps),
		burst:    burst,
		ttl:      ttl,
	}

	go rl.cleanupVisitors()

	return rl
}

func (l *RateLimiter) GetVisitor(id string) *rate.Limiter {
	l.mu.RLock()
	v, exists := l.visitors[id]
	l.mu.RUnlock()

	if !exists {
		limiter := rate.NewLimiter(l.limit, l.burst)
		l.mu.Lock()
		l.visitors[id] = &visitor{limiter, time.Now()}
		l.mu.Unlock()

		return limiter
	}

	v.lastSeen = time.Now()

	return v.limiter
}

func (l *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(30 * time.Second)

		l.mu.Lock()
		for ip, visitor := range l.visitors {
			if time.Since(visitor.lastSeen) > l.ttl {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}
