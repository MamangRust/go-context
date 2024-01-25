package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Download struct {
	ID   int
	Name string
	Size int
}

func downloadFile(ctx context.Context, download Download) {
	duration := time.Duration(rand.Intn(5)) * time.Second

	select {
	case <-time.After(duration):
		fmt.Printf("Download %s (ID: %d, Size: %d MB): Completed\n", download.Name, download.ID, download.Size)

	case <-ctx.Done():
		fmt.Printf("Download %s (ID: %d, Size: %d MB): Canceled\n", download.Name, download.ID, download.Size)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	downloads := []Download{
		{ID: 1, Name: "File_A.txt", Size: 5},
		{ID: 2, Name: "File_B.txt", Size: 8},
		{ID: 3, Name: "File_C.txt", Size: 3},
	}

	for _, download := range downloads {
		go downloadFile(ctx, download)
	}

	time.Sleep(6 * time.Second)

	defer cancel()

	fmt.Println("Program: All Downloads complete or cancelled")
}
