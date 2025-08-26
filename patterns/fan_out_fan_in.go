package patterns

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func FanOutFanIn() {
	fmt.Println("=== Fan-out/Fan-in Pattern ===")
	fmt.Println("Distributing work to multiple goroutines, then collecting results")
	fmt.Println()

	// Run concurrent version
	fmt.Println("Running CONCURRENT version...")
	concurrentStart := time.Now()
	runFanOutFanInConcurrent()
	concurrentDuration := time.Since(concurrentStart)

	fmt.Printf("\nCONCURRENT version took: %v\n\n", concurrentDuration)

	// Run sequential version for comparison
	fmt.Println("Running SEQUENTIAL version for comparison...")
	sequentialStart := time.Now()
	runFanOutFanInSequential()
	sequentialDuration := time.Since(sequentialStart)

	fmt.Printf("\nSEQUENTIAL version took: %v\n", sequentialDuration)
	fmt.Printf("Speedup: %.2fx faster with concurrency!\n\n", float64(sequentialDuration)/float64(concurrentDuration))
}

func runFanOutFanInConcurrent() {
	
	// Input data
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	
	// Fan-out: distribute work
	input := make(chan int)
	
	// Start multiple workers (fan-out)
	const numWorkers = 3
	var outputs []<-chan int
	
	for i := 0; i < numWorkers; i++ {
		output := make(chan int)
		outputs = append(outputs, output)
		go fanOutWorker(i+1, input, output)
	}
	
	// Send input data
	go func() {
		defer close(input)
		for _, num := range numbers {
			input <- num
		}
	}()
	
	// Fan-in: collect results from all workers
	results := fanIn(outputs...)
	
	// Count processed results
	var processed int
	for range results {
		processed++
	}
	
	fmt.Printf("Processed %d numbers with %d workers\n", processed, numWorkers)
}

func runFanOutFanInSequential() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	
	for _, num := range numbers {
		// Simulate processing with same average delay as concurrent version
		processingTime := time.Duration(rand.Intn(200)+50) * time.Millisecond
		time.Sleep(processingTime)
		
		_ = num * num // Square the number
	}
	
	fmt.Printf("Processed %d numbers sequentially\n", len(numbers))
}

func fanOutWorker(id int, input <-chan int, output chan<- int) {
	defer close(output)
	for num := range input {
		// Simulate processing with random delay
		processingTime := time.Duration(rand.Intn(200)+50) * time.Millisecond
		time.Sleep(processingTime)
		
		result := num * num // Square the number
		output <- result
	}
}

func fanIn(inputs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	output := make(chan int)
	
	// Start a goroutine for each input channel
	for _, input := range inputs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				output <- val
			}
		}(input)
	}
	
	// Close output channel when all input channels are done
	go func() {
		wg.Wait()
		close(output)
	}()
	
	return output
}