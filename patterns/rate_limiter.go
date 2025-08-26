package patterns

import (
	"fmt"
	"time"
)

func RateLimiter() {
	fmt.Println("=== Rate Limiter Pattern ===")
	fmt.Println("Controlling the rate of operations to prevent overwhelming resources")
	fmt.Println("Use case: API client making requests with rate limiting to avoid being blocked")
	fmt.Println()

	// Run concurrent version
	fmt.Println("Running CONCURRENT (rate-limited) version...")
	concurrentStart := time.Now()
	runRateLimiterConcurrent()
	concurrentDuration := time.Since(concurrentStart)

	fmt.Printf("\nCONCURRENT (rate-limited) version took: %v\n\n", concurrentDuration)

	// Run sequential version for comparison
	fmt.Println("Running SEQUENTIAL (unlimited) version for comparison...")
	sequentialStart := time.Now()
	runRateLimiterSequential()
	sequentialDuration := time.Since(sequentialStart)

	fmt.Printf("\nSEQUENTIAL (unlimited) version took: %v\n", sequentialDuration)
	fmt.Printf("Note: Rate limiter adds controlled delay vs unlimited requests\n")
	fmt.Printf("Rate limiter prevents resource exhaustion and API blocks!\n\n")
}

func runRateLimiterConcurrent() {
	
	// Create rate limiter: 3 requests per second
	rateLimiter := time.NewTicker(time.Second / 3)
	defer rateLimiter.Stop()

	// Burst limiter: allow up to 2 requests immediately
	burstLimiter := make(chan struct{}, 2)
	for i := 0; i < cap(burstLimiter); i++ {
		burstLimiter <- struct{}{}
	}

	// Refill burst limiter periodically
	go func() {
		ticker := time.NewTicker(time.Second / 3)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case burstLimiter <- struct{}{}:
			default:
			}
		}
	}()

	// Simulate API requests
	requests := []string{
		"GET /api/users",
		"POST /api/users",
		"GET /api/posts",
		"PUT /api/users/1",
		"DELETE /api/posts/5",
		"GET /api/comments",
		"POST /api/posts",
		"GET /api/analytics",
		"PUT /api/settings",
		"GET /api/dashboard",
	}

	var completed int
	for _, request := range requests {
		// Wait for burst token or rate limit
		select {
		case <-burstLimiter:
			// Use burst token immediately
		default:
			// Wait for rate limiter
			<-rateLimiter.C
		}

		// Simulate API call processing time
		time.Sleep(50 * time.Millisecond)
		completed++
		_ = request // Use the request variable
	}

	fmt.Printf("Completed %d rate-limited requests\n", completed)
}

func runRateLimiterSequential() {
	requests := []string{
		"GET /api/users",
		"POST /api/users",
		"GET /api/posts",
		"PUT /api/users/1",
		"DELETE /api/posts/5",
		"GET /api/comments",
		"POST /api/posts",
		"GET /api/analytics",
		"PUT /api/settings",
		"GET /api/dashboard",
	}

	for _, request := range requests {
		// Simulate API call processing time (same as concurrent)
		time.Sleep(50 * time.Millisecond)
		_ = request // Use the request variable
	}

	fmt.Printf("Completed %d unlimited requests\n", len(requests))
	fmt.Println("⚠️  Warning: This approach might get blocked by API rate limits!")
}