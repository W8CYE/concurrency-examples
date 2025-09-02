package patterns

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type CircuitState int

const (
	CLOSED CircuitState = iota
	OPEN
	HALF_OPEN
)

func (cs CircuitState) String() string {
	switch cs {
	case CLOSED:
		return "🟢 CLOSED"
	case OPEN:
		return "🔴 OPEN"
	case HALF_OPEN:
		return "🟡 HALF_OPEN"
	default:
		return "❓ UNKNOWN"
	}
}

type CircuitBreaker struct {
	state          CircuitState
	failureCount   int
	lastFailure    time.Time
	failureThreshold int
	timeout        time.Duration
	mutex          sync.RWMutex
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            CLOSED,
		failureThreshold: threshold,
		timeout:          timeout,
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == OPEN {
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = HALF_OPEN
			cb.failureCount = 0
		} else {
			return fmt.Errorf("circuit breaker is OPEN")
		}
	}

	err := fn()
	if err != nil {
		cb.failureCount++
		
		if cb.state == HALF_OPEN {
			cb.state = OPEN
			cb.lastFailure = time.Now()
		} else {
			cb.lastFailure = time.Now()
			if cb.failureCount >= cb.failureThreshold {
				cb.state = OPEN
			}
		}
		return err
	}

	// Success case
	if cb.state == HALF_OPEN {
		cb.state = CLOSED
	}
	cb.failureCount = 0
	return nil
}

func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

func CircuitBreakerDemo() {
	fmt.Println("=== Circuit Breaker Pattern ===")
	fmt.Println("Preventing cascading failures by monitoring service health")
	fmt.Println("Use case: External API calls with automatic failure detection")
	fmt.Println()

	for {
		fmt.Println("Circuit Breaker Demo Options:")
		fmt.Println("1. 🟢 CLOSED state demo (healthy service)")
		fmt.Println("2. 🔴 OPEN state demo (failing service)")
		fmt.Println("3. 🟡 HALF_OPEN state demo (recovery attempt)")
		fmt.Println("4. ❌ No Circuit Breaker (comparison)")
		fmt.Println("5. 🔄 Full Lifecycle Demo")
		fmt.Println("0. Back to main menu")
		fmt.Print("Select demo (0-5): ")

		var choice int
		fmt.Scanf("%d", &choice)
		fmt.Println()

		switch choice {
		case 1:
			runClosedStateDemo()
		case 2:
			runOpenStateDemo()
		case 3:
			runHalfOpenStateDemo()
		case 4:
			runNoCircuitBreakerDemo()
		case 5:
			runFullLifecycleDemo()
		case 0:
			return
		default:
			fmt.Println("Invalid choice. Please try again.\n")
		}
		
		fmt.Println("\nPress Enter to continue...")
		fmt.Scanf("\n")
		fmt.Println()
	}
}

func runClosedStateDemo() {
	fmt.Println("🟢 === CLOSED State Demo ===")
	fmt.Println("Circuit is closed - all requests pass through normally")
	fmt.Println()

	cb := NewCircuitBreaker(3, 5*time.Second)
	var successful, failed int

	for i := 1; i <= 10; i++ {
		fmt.Printf("Request %d: ", i)
		
		err := cb.Call(func() error {
			return simulateHealthyService()
		})

		if err != nil {
			failed++
			fmt.Printf("❌ Failed - %v\n", err)
		} else {
			successful++
			fmt.Printf("✅ Success (State: %s)\n", cb.GetState())
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("\n📊 Results: %d successful, %d failed\n", successful, failed)
	fmt.Printf("🔧 Circuit remained CLOSED - all requests processed\n")
}

func runOpenStateDemo() {
	fmt.Println("🔴 === OPEN State Demo ===")
	fmt.Println("Circuit is open - requests are blocked to protect failing service")
	fmt.Println()

	cb := NewCircuitBreaker(3, 5*time.Second)
	var successful, failed, blocked int

	// First, trigger the circuit to open by simulating failures
	fmt.Println("Triggering circuit to open with failures...")
	for i := 1; i <= 3; i++ {
		cb.Call(func() error {
			return fmt.Errorf("service unavailable")
		})
	}

	// Now show blocked requests
	for i := 1; i <= 8; i++ {
		fmt.Printf("Request %d: ", i)
		
		err := cb.Call(func() error {
			return simulateHealthyService()
		})

		if err != nil {
			if err.Error() == "circuit breaker is OPEN" {
				blocked++
				fmt.Printf("🛑 BLOCKED by circuit breaker (State: %s)\n", cb.GetState())
			} else {
				failed++
				fmt.Printf("❌ Failed - %v\n", err)
			}
		} else {
			successful++
			fmt.Printf("✅ Success (State: %s)\n", cb.GetState())
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("\n📊 Results: %d successful, %d failed, %d blocked\n", successful, failed, blocked)
	fmt.Printf("🛡️  Circuit breaker protected the failing service from %d requests\n", blocked)
}

func runHalfOpenStateDemo() {
	fmt.Println("🟡 === HALF_OPEN State Demo ===")
	fmt.Println("Circuit allows ONE test request to check service recovery")
	fmt.Println()

	cb := NewCircuitBreaker(3, 2*time.Second)
	var successful, failed, blocked int

	// Trigger circuit to open
	fmt.Println("Opening circuit with failures...")
	for i := 1; i <= 3; i++ {
		cb.Call(func() error {
			return fmt.Errorf("service down")
		})
	}
	fmt.Printf("Circuit State: %s\n\n", cb.GetState())

	// Wait for timeout to allow half-open
	fmt.Println("⏰ Waiting for timeout to allow recovery test...")
	time.Sleep(2100 * time.Millisecond)
	
	// First cycle: Failed recovery test
	fmt.Printf("Circuit State: %s (timeout expired, ready for test)\n", cb.GetState())
	fmt.Println("→ Next request will transition to HALF_OPEN for testing")
	
	fmt.Print("Test Request 1: ")
	err := cb.Call(func() error {
		return fmt.Errorf("service still failing")
	})
	
	if err != nil {
		if err.Error() == "circuit breaker is OPEN" {
			blocked++
			fmt.Printf("🛑 BLOCKED")
		} else {
			failed++
			fmt.Printf("❌ Failed - %v", err)
		}
	} else {
		successful++
		fmt.Printf("✅ Success!")
	}
	fmt.Printf(" (State after call: %s)\n", cb.GetState())
	fmt.Println("→ Test failed, circuit returned to OPEN\n")
	
	// Show blocking during OPEN
	for i := 2; i <= 4; i++ {
		fmt.Printf("Request %d: ", i)
		err := cb.Call(func() error {
			return simulateHealthyService()
		})
		
		if err != nil && err.Error() == "circuit breaker is OPEN" {
			blocked++
			fmt.Printf("🛑 BLOCKED (State: %s)\n", cb.GetState())
		}
		time.Sleep(200 * time.Millisecond)
	}
	
	// Second cycle: Successful recovery
	fmt.Println("\n⏰ Waiting for next recovery window...")
	time.Sleep(2100 * time.Millisecond)
	
	fmt.Printf("Circuit State: %s (timeout expired, ready for test)\n", cb.GetState())
	fmt.Println("→ Next request will transition to HALF_OPEN for testing")
	
	fmt.Print("Test Request 5: ")
	err = cb.Call(func() error {
		return simulateHealthyService() // This will succeed
	})
	
	if err != nil {
		if err.Error() == "circuit breaker is OPEN" {
			blocked++
			fmt.Printf("🛑 BLOCKED")
		} else {
			failed++
			fmt.Printf("❌ Failed - %v", err)
		}
	} else {
		successful++
		fmt.Printf("✅ Success!")
	}
	fmt.Printf(" (State after call: %s)\n", cb.GetState())
	fmt.Println("→ Test succeeded, circuit is now CLOSED and healthy!\n")

	fmt.Printf("📊 Results: %d successful, %d failed, %d blocked\n", successful, failed, blocked)
	fmt.Printf("🔄 HALF_OPEN allows exactly ONE test request to determine recovery\n")
}

func runNoCircuitBreakerDemo() {
	fmt.Println("❌ === No Circuit Breaker Demo ===")
	fmt.Println("Direct calls to failing service - shows the problem circuit breakers solve")
	fmt.Println()

	var successful, failed int
	
	for i := 1; i <= 10; i++ {
		fmt.Printf("Request %d: ", i)
		
		err := simulateFailingService()
		if err != nil {
			failed++
			fmt.Printf("❌ Failed - %v (wasted resources!)\n", err)
		} else {
			successful++
			fmt.Printf("✅ Success\n")
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("\n📊 Results: %d successful, %d failed\n", successful, failed)
	fmt.Printf("⚠️  Without circuit breaker: %d requests wasted on failing service!\n", failed)
	fmt.Printf("🔥 This could cause cascading failures in production!\n")
}

func runFullLifecycleDemo() {
	fmt.Println("🔄 === Full Circuit Breaker Lifecycle ===")
	fmt.Println("Watch circuit breaker automatically handle service degradation and recovery")
	fmt.Println()

	cb := NewCircuitBreaker(3, 3*time.Second)
	var successful, failed, blocked int

	// Phase 1: Healthy service (CLOSED)
	fmt.Println("📡 Phase 1: Healthy service...")
	for i := 1; i <= 5; i++ {
		fmt.Printf("Request %d: ", i)
		err := cb.Call(simulateHealthyService)
		if err != nil {
			failed++
			fmt.Printf("❌ Failed (State: %s)\n", cb.GetState())
		} else {
			successful++
			fmt.Printf("✅ Success (State: %s)\n", cb.GetState())
		}
		time.Sleep(300 * time.Millisecond)
	}

	// Phase 2: Service starts failing (CLOSED → OPEN)
	fmt.Println("\n💥 Phase 2: Service degrading...")
	for i := 6; i <= 10; i++ {
		fmt.Printf("Request %d: ", i)
		err := cb.Call(simulateFailingService)
		if err != nil {
			if err.Error() == "circuit breaker is OPEN" {
				blocked++
				fmt.Printf("🛑 BLOCKED (State: %s)\n", cb.GetState())
			} else {
				failed++
				fmt.Printf("❌ Failed (State: %s)\n", cb.GetState())
			}
		} else {
			successful++
			fmt.Printf("✅ Success (State: %s)\n", cb.GetState())
		}
		time.Sleep(300 * time.Millisecond)
	}

	// Phase 3: Wait and try recovery (OPEN → HALF_OPEN)
	fmt.Println("\n⏰ Phase 3: Waiting for recovery window...")
	time.Sleep(3100 * time.Millisecond)

	for i := 11; i <= 15; i++ {
		fmt.Printf("Request %d: ", i)
		err := cb.Call(simulateRecoveringService)
		if err != nil {
			if err.Error() == "circuit breaker is OPEN" {
				blocked++
				fmt.Printf("🛑 BLOCKED (State: %s)\n", cb.GetState())
			} else {
				failed++
				fmt.Printf("❌ Failed (State: %s)\n", cb.GetState())
			}
		} else {
			successful++
			fmt.Printf("✅ Success! (State: %s)\n", cb.GetState())
		}
		time.Sleep(400 * time.Millisecond)
	}

	fmt.Printf("\n📊 Final Results: %d successful, %d failed, %d blocked\n", successful, failed, blocked)
	fmt.Printf("🛡️  Circuit breaker prevented %d requests to failing service\n", blocked)
	fmt.Printf("⚡ Automatic recovery detection enabled graceful service restoration\n")
}

func simulateHealthyService() error {
	time.Sleep(50 * time.Millisecond)
	return nil
}

func simulateFailingService() error {
	time.Sleep(100 * time.Millisecond)
	return fmt.Errorf("service unavailable")
}

func simulateRecoveringService() error {
	time.Sleep(75 * time.Millisecond)
	// 70% chance of success during recovery
	if rand.Float32() < 0.7 {
		return nil
	}
	return fmt.Errorf("service still unstable")
}