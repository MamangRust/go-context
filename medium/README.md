## Studi Kasus Golang dengan Penggunaan Paket `context` dan `Goroutine` (Medium)

### Deskripsi:

Buatlah program pengelolaan sederhana yang menunjukkan penggunaan paket context dan goroutine. Program ini mensimulasikan pengunduhan beberapa file dari internet dengan batasan waktu.

### Persyaratan:

- Program harus dapat mengunduh beberapa file secara konkuren.
- Gunakan paket context untuk mengatur batasan waktu dan membatalkan operasi unduhan jika waktu eksekusi melebihi batas.
- Tampilkan informasi tentang file yang berhasil diunduh dan file yang dibatalkan.
- Implementasikan penggunaan goroutine untuk mengelola unduhan file secara konkuren.
- Program harus berhenti setelah semua unduhan selesai atau setelah batas waktu tertentu tercapai.

## Contoh:

```go
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

```

### Penjelasan:

- Program membuat beberapa objek Download untuk mensimulasikan file yang akan diunduh, masing-masing dengan ID, nama, dan ukuran.
- Konteks dengan batasan waktu 10 detik dibuat menggunakan context.WithTimeout.
- Setiap unduhan dijalankan sebagai goroutine yang menggunakan select untuk memeriksa apakah waktu eksekusi sudah melebihi batas atau pembatalan telah diminta.
- Program menggunakan time.Sleep untuk memberi cukup waktu kepada goroutine untuk menyelesaikan output mereka sebelum program berakhir.
- Program mencetak pesan setelah semua unduhan selesai atau setelah batas waktu tertentu tercapai.
