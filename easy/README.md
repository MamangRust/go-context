## Studi Kasus Golang dengan Penggunaan Paket `context`

### Deskripsi:

Buatlah sebuah program sederhana dalam bahasa Go yang menggunakan paket `context` untuk mengelola operasi yang memerlukan batasan waktu. Program ini dapat mengeksekusi beberapa tugas konkuren dan membatalkannya jika waktu eksekusi melebihi batas yang ditentukan.

### Persyaratan:

- Program harus memiliki tugas konkuren yang dapat dibatalkan.
- Setiap tugas harus menerima konteks sebagai parameter dan memeriksa apakah pembatalan telah diminta.
- Gunakan fungsi-fungsi dari paket context seperti WithCancel atau WithTimeout untuk mengatur batasan waktu dan pembatalan.
- Program harus memberikan informasi yang jelas tentang tugas mana yang berhasil diselesaikan dan tugas mana yang dibatalkan.

### Contoh:

```go
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

```

### Penjelasan:

- Program ini membuat beberapa tugas konkuren yang masing-masing mensimulasikan tugas yang membutuhkan waktu acak untuk diselesaikan.
- Konteks dengan batasan waktu 3 detik dibuat menggunakan context.WithTimeout.
- Setiap tugas menggunakan select untuk memeriksa apakah waktu eksekusi sudah melebihi batas atau pembatalan telah diminta.
- Program menggunakan select lagi untuk menunggu semua tugas selesai atau pembatalan konteks.
- Jika waktu eksekusi melebihi batas, program membatalkan konteks dan mencetak pesan yang sesuai. Jika semua tugas selesai, program mencetak pesan bahwa semua tugas telah diselesaikan.
