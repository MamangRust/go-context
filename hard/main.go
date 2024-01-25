package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Task struct {
	ID       int
	Name     string
	Status   string
	Cancel   context.CancelFunc
	Canceled bool
}

type TaskResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Canceled bool   `json:"canceled"`
}

var (
	tasks     []Task
	taskIDSeq int
	mux       sync.Mutex
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	task.ID = taskIDSeq
	task.Name = fmt.Sprintf("Task %d", task.ID)
	task.Status = "Queued"
	task.Cancel = cancel

	mux.Lock()
	tasks = append(tasks, task)
	taskIDSeq++
	mux.Unlock()

	go processTask(ctx, task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, "Task added successfully. Task ID: %d\n", task.ID)

}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	w.Header().Set("Content-Type", "application/json")

	if len(tasks) == 0 {
		http.Error(w, "No tasks found", http.StatusNotFound)
		return
	}

	var simplifiedTasks []TaskResponse
	for _, task := range tasks {
		simplifiedTask := TaskResponse{
			ID:       task.ID,
			Name:     task.Name,
			Status:   task.Status,
			Canceled: task.Canceled,
		}
		simplifiedTasks = append(simplifiedTasks, simplifiedTask)
	}

	json.NewEncoder(w).Encode(simplifiedTasks)
}

func cancelTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/cancel/"):]
	mux.Lock()
	defer mux.Unlock()

	for i, task := range tasks {
		if fmt.Sprint(task.ID) == taskID && task.Status == "Queued" {
			task.Cancel()
			tasks[i].Canceled = true
			tasks[i].Status = "Canceled"
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Task canceled successfully\n")
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func processTask(ctx context.Context, task Task) {
	defer func() {
		mux.Lock()
		defer mux.Unlock()
		task.Status = "Completed"
	}()

	duration := time.Duration(5) * time.Second

	select {
	case <-time.After(duration):

		fmt.Println("Task completed")

	case <-ctx.Done():

		mux.Lock()
		defer mux.Unlock()
		task.Status = "Canceled"
	}

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskHandler(w, r)
		case http.MethodGet:
			getTasksHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/cancel/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			cancelTaskHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	fmt.Println("Server is running on :8080")

	stop := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("\nShutting down server...")
		server.Shutdown(context.Background())
		close(stop)
	}()

	<-stop
	fmt.Println("Server gracefully stopped.")
}
