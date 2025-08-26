package patterns

import (
	"fmt"
	"math/rand"
	"time"
)

func SelectTimeout() {
	fmt.Println("=== Select Statement with Timeout Pattern ===")
	fmt.Println("Non-blocking channel operations with timeouts and graceful error handling")
	fmt.Println("Use case: Service health checks with timeouts to prevent hanging")
	fmt.Println()

	// Run concurrent version
	fmt.Println("Running CONCURRENT (with timeouts) version...")
	concurrentStart := time.Now()
	runSelectTimeoutConcurrent()
	concurrentDuration := time.Since(concurrentStart)

	fmt.Printf("\nCONCURRENT (with timeouts) version took: %v\n\n", concurrentDuration)

	// Run sequential version for comparison
	fmt.Println("Running SEQUENTIAL (blocking) version for comparison...")
	sequentialStart := time.Now()
	runSelectTimeoutSequential()
	sequentialDuration := time.Since(sequentialStart)

	fmt.Printf("\nSEQUENTIAL (blocking) version took: %v\n", sequentialDuration)
	fmt.Printf("Concurrent version handles failures gracefully with timeouts!\n\n")
}

func runSelectTimeoutConcurrent() {
	
	services := []string{
		"Database Service",
		"Cache Service", 
		"Auth Service",
		"Payment Service",
		"Notification Service",
	}

	var healthyServices, timeoutServices, failedServices int

	for _, service := range services {
		// Create channels for different outcomes
		resultCh := make(chan string, 1)
		errorCh := make(chan error, 1)

		// Start health check in goroutine
		go func(svc string) {
			// Simulate variable response times and failures
			responseTime := time.Duration(rand.Intn(800)+100) * time.Millisecond
			
			// 20% chance of service being down
			if rand.Float32() < 0.2 {
				time.Sleep(responseTime)
				errorCh <- fmt.Errorf("%s is down", svc)
				return
			}

			time.Sleep(responseTime)
			resultCh <- fmt.Sprintf("%s is healthy (response time: %v)", svc, responseTime)
		}(service)

		// Use select with timeout
		select {
		case <-resultCh:
			healthyServices++

		case <-errorCh:
			failedServices++

		case <-time.After(500 * time.Millisecond):
			timeoutServices++

		// Demonstrate non-blocking select with default case
		default:
			// Wait again with timeout after showing checking message
			select {
			case <-resultCh:
				healthyServices++

			case <-errorCh:
				failedServices++

			case <-time.After(500 * time.Millisecond):
				timeoutServices++
			}
		}
	}

	fmt.Printf("Health Check Results - Healthy: %d, Failed: %d, Timeouts: %d\n", healthyServices, failedServices, timeoutServices)
}

func runSelectTimeoutSequential() {
	services := []string{
		"Database Service",
		"Cache Service", 
		"Auth Service",
		"Payment Service",
		"Notification Service",
	}

	var healthyServices, failedServices int

	for i, service := range services {
		// Simulate variable response times and failures - blocking call
		responseTime := time.Duration(rand.Intn(800)+100) * time.Millisecond
		time.Sleep(responseTime)

		// 20% chance of service being down
		if rand.Float32() < 0.2 {
			failedServices++
		} else {
			healthyServices++
		}

		// If a service hangs, this would block forever!
		// Simulate one hanging service
		if i == 2 && rand.Float32() < 0.3 {
			time.Sleep(2 * time.Second)
		}
		
		_ = service // Use the service variable
	}

	fmt.Printf("Sequential Results - Healthy: %d, Failed: %d\n", healthyServices, failedServices)
	fmt.Println("⚠️  Note: Sequential approach vulnerable to hanging services!")
}