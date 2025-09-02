package main

import (
	"bufio"
	"concurrency-examples.git/patterns"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("=== Go Concurrency Patterns Showcase ===")
	fmt.Println()
	
	for {
		showMenu()
		choice := getUserInput()
		
		switch choice {
		case 1:
			patterns.WorkerPool()
		case 2:
			patterns.FanOutFanIn()
		case 3:
			patterns.Pipeline()
		case 4:
			patterns.RateLimiter()
		case 5:
			patterns.SelectTimeout()
		case 6:
			patterns.CircuitBreakerDemo()
		case 0:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.\n")
		}
	}
}

func showMenu() {
	fmt.Println("Available Concurrency Patterns:")
	fmt.Println("1. Worker Pool")
	fmt.Println("2. Fan-out/Fan-in")
	fmt.Println("3. Pipeline")
	fmt.Println("4. Rate Limiter")
	fmt.Println("5. Select with Timeout")
	fmt.Println("6. Circuit Breaker")
	fmt.Println("0. Exit")
	fmt.Print("Select a pattern to run (0-6): ")
}

func getUserInput() int {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return -1
	}
	
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	
	return choice
}
