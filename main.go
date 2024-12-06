package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sony/gobreaker" // Mengimpor pustaka "gobreaker" untuk fungsi circuit breaker
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Menginisialisasi generator angka acak dengan waktu saat ini untuk memastikan hasil acak

	// Membuat circuit breaker baru dengan pengaturan tertentu
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:     "myCircuitBreaker", // Nama deskriptif untuk circuit breaker
		Timeout:  5 * time.Second,    // Waktu circuit breaker tetap dalam keadaan terbuka sebelum berpindah ke Half-Open
		Interval: 10 * time.Second,   // Interval untuk mereset penghitung kegagalan dalam keadaan Closed
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Menentukan kapan circuit breaker akan trip (berpindah ke Open)
			// Circuit breaker akan trip jika jumlah kegagalan melebihi 2
			return counts.TotalFailures > 2
		},
	})

	for i := 0; i < 20; i++ { // Menjalankan simulasi 20 kali percobaan eksekusi fungsi
		// Mencoba mengeksekusi fungsi di dalam circuit breaker
		result, err := cb.Execute(func() (interface{}, error) {
			// Mensimulasikan operasi yang berpotensi gagal
			// Gagal secara acak dengan kemungkinan 30% (jika rand.Intn(10) > 7)
			if rand.Intn(10) > 7 {
				return nil, fmt.Errorf("operation failed") // Mengembalikan error untuk mensimulasikan kegagalan
			}
			return "success", nil // Mengembalikan sukses jika tidak gagal
		})

		if err != nil {
			// Menangani kasus error
			state := cb.State() // Memeriksa keadaan circuit breaker saat ini
			if state == gobreaker.StateOpen {
				// Circuit breaker dalam keadaan Open, permintaan baru diblokir
				fmt.Println("Circuit breaker terbuka, coba lagi nanti")
			} else {
				// Operasi gagal bukan karena circuit breaker terbuka
				fmt.Println("Operasi gagal:", err)
			}
		} else {
			// Operasi berhasil
			fmt.Println("Operasi berhasil:", result)
		}

		// Jeda selama 1 detik sebelum mencoba lagi untuk mensimulasikan penggunaan dunia nyata
		time.Sleep(1 * time.Second)
	}
}
