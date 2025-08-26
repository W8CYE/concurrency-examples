package patterns

import (
	"fmt"
	"sync"
	"time"
)

func WorkerPool() {
	fmt.Println("=== Worker Pool Pattern ===")
	fmt.Println("Multiple workers processing jobs from a shared channel")
	fmt.Println()

	// Run concurrent version
	fmt.Println("Running CONCURRENT version...")
	concurrentStart := time.Now()
	runWorkerPoolConcurrent()
	concurrentDuration := time.Since(concurrentStart)

	fmt.Printf("\nCONCURRENT version took: %v\n\n", concurrentDuration)

	// Run sequential version for comparison
	fmt.Println("Running SEQUENTIAL version for comparison...")
	sequentialStart := time.Now()
	runWorkerPoolSequential()
	sequentialDuration := time.Since(sequentialStart)

	fmt.Printf("\nSEQUENTIAL version took: %v\n", sequentialDuration)
	fmt.Printf("Speedup: %.2fx faster with concurrency!\n\n", float64(sequentialDuration)/float64(concurrentDuration))
}

func runWorkerPoolConcurrent() {
	
	const numWorkers = 3
	const numJobs = 10
	
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	var wg sync.WaitGroup
	
	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}
	
	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Wait for all workers to finish
	wg.Wait()
	close(results)
	
	// Count completed jobs
	var completed int
	for range results {
		completed++
	}
	
	fmt.Printf("Completed %d jobs with %d workers\n", completed, numWorkers)
}

func runWorkerPoolSequential() {
	const numJobs = 10
	
	for j := 1; j <= numJobs; j++ {
		time.Sleep(100 * time.Millisecond) // Same work simulation as concurrent version
	}
	
	fmt.Printf("Completed %d jobs sequentially\n", numJobs)
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		time.Sleep(100 * time.Millisecond) // Simulate work
		results <- job
	}
}