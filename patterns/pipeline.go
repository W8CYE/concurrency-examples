package patterns

import (
	"fmt"
	"strings"
	"time"
)

func Pipeline() {
	fmt.Println("=== Pipeline Pattern ===")
	fmt.Println("Processing data through multiple concurrent stages")
	fmt.Println("Use case: Text processing pipeline (clean -> transform -> analyze)")
	fmt.Println()

	// Run concurrent version
	fmt.Println("Running CONCURRENT version...")
	concurrentStart := time.Now()
	runPipelineConcurrent()
	concurrentDuration := time.Since(concurrentStart)

	fmt.Printf("\nCONCURRENT version took: %v\n\n", concurrentDuration)

	// Run sequential version for comparison
	fmt.Println("Running SEQUENTIAL version for comparison...")
	sequentialStart := time.Now()
	runPipelineSequential()
	sequentialDuration := time.Since(sequentialStart)

	fmt.Printf("\nSEQUENTIAL version took: %v\n", sequentialDuration)
	fmt.Printf("Speedup: %.2fx faster with concurrency!\n\n", float64(sequentialDuration)/float64(concurrentDuration))
}

func runPipelineConcurrent() {
	
	// Sample data to process
	rawData := []string{
		"  Hello World!!!  ",
		"  Go is AWESOME  ",
		"  Concurrency ROCKS!!!  ",
		"  Programming is FUN  ",
		"  Pipelines are COOL  ",
		"  Channels are GREAT  ",
		"  Goroutines RULE  ",
		"  Synchronization MATTERS  ",
	}

	// Stage 1: Clean data (trim whitespace, remove extra punctuation)
	cleaned := cleanStage(generator(rawData))
	
	// Stage 2: Transform data (convert to lowercase, add prefix)
	transformed := transformStage(cleaned)
	
	// Stage 3: Analyze data (count words, measure length)
	analyzed := analyzeStage(transformed)

	// Count results
	var processed int
	for range analyzed {
		processed++
	}
	
	fmt.Printf("Processed %d items through 3-stage pipeline\n", processed)
}

func runPipelineSequential() {
	rawData := []string{
		"  Hello World!!!  ",
		"  Go is AWESOME  ",
		"  Concurrency ROCKS!!!  ",
		"  Programming is FUN  ",
		"  Pipelines are COOL  ",
		"  Channels are GREAT  ",
		"  Goroutines RULE  ",
		"  Synchronization MATTERS  ",
	}

	for _, data := range rawData {
		// Stage 1: Clean
		time.Sleep(50 * time.Millisecond) // Simulate cleaning work
		cleaned := strings.TrimSpace(data)
		cleaned = strings.ReplaceAll(cleaned, "!!!", "!")

		// Stage 2: Transform  
		time.Sleep(30 * time.Millisecond) // Simulate transform work
		transformed := "processed: " + strings.ToLower(cleaned)

		// Stage 3: Analyze
		time.Sleep(40 * time.Millisecond) // Simulate analysis work
		wordCount := len(strings.Fields(transformed))
		_ = fmt.Sprintf("%s (words: %d, length: %d)", transformed, wordCount, len(transformed))
	}

	fmt.Printf("Processed %d items sequentially through all stages\n", len(rawData))
}

func generator(data []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, item := range data {
			out <- item
		}
	}()
	return out
}

func cleanStage(input <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for data := range input {
			// Simulate cleaning work
			time.Sleep(50 * time.Millisecond)
			
			cleaned := strings.TrimSpace(data)
			cleaned = strings.ReplaceAll(cleaned, "!!!", "!")
			out <- cleaned
		}
	}()
	return out
}

func transformStage(input <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for data := range input {
			// Simulate transformation work
			time.Sleep(30 * time.Millisecond)
			
			transformed := "processed: " + strings.ToLower(data)
			out <- transformed
		}
	}()
	return out
}

func analyzeStage(input <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for data := range input {
			// Simulate analysis work
			time.Sleep(40 * time.Millisecond)
			
			wordCount := len(strings.Fields(data))
			result := fmt.Sprintf("%s (words: %d, length: %d)", data, wordCount, len(data))
			out <- result
		}
	}()
	return out
}