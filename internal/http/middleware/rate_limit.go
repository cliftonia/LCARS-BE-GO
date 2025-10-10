package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// rateLimiter holds rate limiters for each IP address
type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// newRateLimiter creates a new rate limiter
func newRateLimiter(requestsPerMinute int, burst int) *rateLimiter {
	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerMinute) / 60, // Convert to per-second rate
		burst:    burst,
	}
}

// getLimiter returns the rate limiter for a given IP
func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter

		// Start cleanup goroutine for this IP (clean up after 3 minutes of inactivity)
		go func() {
			time.Sleep(3 * time.Minute)
			rl.mu.Lock()
			delete(rl.limiters, ip)
			rl.mu.Unlock()
		}()
	}

	return limiter
}

// RateLimit creates a rate limiting middleware
// requestsPerMinute: number of requests allowed per minute
// burst: maximum burst size
func RateLimit(requestsPerMinute, burst int) func(http.Handler) http.Handler {
	rl := newRateLimiter(requestsPerMinute, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP from request
			ip := r.RemoteAddr
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				ip = forwardedFor
			}

			// Get rate limiter for this IP
			limiter := rl.getLimiter(ip)

			// Check if request is allowed
			if !limiter.Allow() {
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
