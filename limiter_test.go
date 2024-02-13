package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_GetVisitor(t *testing.T) {
	rps := 2.0
	burst := 5
	ttl := 5 * time.Second

	rl := NewRateLimiter(rps, burst, ttl)

	limiter1 := rl.GetVisitor("user1")
	assert.NotNil(t, limiter1)

	limiter2 := rl.GetVisitor("user1")
	assert.Equal(t, limiter1, limiter2)

	limiter3 := rl.GetVisitor("user2")
	assert.NotEqual(t, limiter3, nil)
}

func TestRateLimiter_CleanupVisitors(t *testing.T) {
	rps := 2.0
	burst := 5
	ttl := 2 * time.Second

	rl := NewRateLimiter(rps, burst, ttl)

	time.Sleep(3 * time.Second)

	rl.mu.RLock()
	_, exists := rl.visitors["user1"]
	rl.mu.RUnlock()
	assert.False(t, exists)
}
