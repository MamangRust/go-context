## Studi Kasus Golang dengan Penggunaan Paket `context` dan `Goroutine` (Hard)

### Deskripsi:

Buatlah program web sederhana untuk mengelola antrian tugas. Pengguna dapat menambahkan tugas, melihat daftar tugas, dan membatalkan tugas yang masih dalam antrian. Program ini akan menggunakan paket `context` dan `goroutine` untuk mengelola tugas-tugas secara konkuren.

### Persyaratan:

- Program harus memiliki server web menggunakan framework seperti Gin atau Echo.
- Implementasikan API endpoint untuk menambahkan tugas baru.
- Implementasikan API endpoint untuk melihat daftar tugas.
- Setiap tugas harus memiliki konteks yang dapat dibatalkan.
- Pengguna dapat membatalkan tugas yang masih dalam antrian menggunakan API endpoint khusus.
- Gunakan goroutine untuk mengelola setiap tugas secara konkuren.
- Tampilkan informasi yang jelas jika tugas berhasil ditambahkan, tampilan daftar tugas, dan jika tugas berhasil dibatalkan.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

var (
	tasks     []Task
	taskIDSeq int
	mux       sync.Mutex
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	task.ID = taskIDSeq
	task.Status = "Queued"
	task.Cancel = cancel

	mux.Lock()
	tasks = append(tasks, task)
	taskIDSeq++
	mux.Unlock()

	go processTask(ctx, task)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Task added successfully. Task ID: %d\n", task.ID)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
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

	http.Error(w, "Task not found or cannot be canceled", http.StatusNotFound)
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

```

### Penjelasan:

- Program ini menggunakan framework web Gin untuk membuat server web.
- Endpoint /tasks digunakan untuk menambahkan tugas baru (POST) dan melihat daftar tugas (GET).
- Endpoint /tasks/:id digunakan untuk membatalkan tugas yang masih dalam antrian (DELETE).
- Setiap tugas memiliki ID, nama, status, fungsi pembatalan (CancelFunc), dan status pembatalan (Canceled).
- Fungsi processTask mensimulasikan tugas yang membutuhkan waktu dan membatalkan tugas jika konteks dibatalkan.
- Program menggunakan goroutine untuk mengelola setiap tugas secara konkuren.
- Server dapat dihentikan dengan aman dengan menanggapi sinyal SIGINT atau SIGTERM.
