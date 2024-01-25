package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func worker(ctx context.Context, taskID int) {
	duration := time.Duration(rand.Intn(5)) * time.Second

	select {
	case <-time.After(duration):
		fmt.Printf("Task %d: Completed\n", taskID)
	case <-ctx.Done():
		fmt.Printf("Task %d: Canceled\n", taskID)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	numTasks := 5

	for i := 1; i <= numTasks; i++ {
		go worker(ctx, i)
	}

	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Program: Timeout reached. Canceling remaining tasks")
		cancel()
	case <-ctx.Done():
		fmt.Println("Program: All Tasks complete or cancelled")
	}

	time.Sleep(1 * time.Second)
}
