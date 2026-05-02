// Package utils menyediakan fungsi-fungsi pembantu untuk input/output,
// validasi, formatting, dan manajemen file di seluruh aplikasi.
package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// reader adalah scanner global agar buffer tidak hilang antar pemanggilan.
var reader = bufio.NewReader(os.Stdin)

// ── INPUT ─────────────────────────────────────────────────────────────────────

// BacaString membaca satu baris input dari terminal dan mengembalikannya
// sebagai string yang sudah di-trim whitespace-nya.
func BacaString(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// BacaInt membaca input integer positif dari terminal.
// Jika input tidak valid atau <= 0, pengguna diminta mengulang.
func BacaInt(prompt string) int {
	for {
		input := BacaString(prompt)
		nilai, err := strconv.Atoi(input)
		if err != nil || nilai <= 0 {
			fmt.Println("  ⚠  Input tidak valid. Masukkan angka bulat positif.")
			continue
		}
		return nilai
	}
}

// BacaIntMin membaca input integer dengan nilai minimum tertentu.
// Digunakan untuk input pembayaran agar tidak kurang dari total belanja.
func BacaIntMin(prompt string, minimum int) int {
	for {
		input := BacaString(prompt)
		nilai, err := strconv.Atoi(input)
		if err != nil || nilai <= 0 {
			fmt.Println("  ⚠  Input tidak valid. Masukkan angka bulat positif.")
			continue
		}
		if nilai < minimum {
			fmt.Printf("  ⚠  Uang kurang! Minimal pembayaran: %s\n", FormatRupiah(minimum))
			continue
		}
		return nilai
	}
}

// TanyaYaTidak menampilkan pertanyaan y/n dan mengembalikan true jika 'y'.
func TanyaYaTidak(prompt string) bool {
	for {
		jawaban := strings.ToLower(BacaString(prompt))
		switch jawaban {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("  ⚠  Masukkan 'y' untuk ya atau 'n' untuk tidak.")
		}
	}
}

// ── FORMATTING ────────────────────────────────────────────────────────────────

// FormatRupiah memformat angka integer menjadi string Rupiah dengan pemisah ribuan.
// Contoh: 15000 → "Rp 15.000"
func FormatRupiah(nominal int) string {
	s := strconv.Itoa(nominal)
	result := ""
	for i, ch := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += "."
		}
		result += string(ch)
	}
	return "Rp " + result
}

// FormatTanggal mengembalikan string tanggal format "02 Jan 2006".
func FormatTanggal(t time.Time) string {
	return t.Format("02 Jan 2006")
}

// FormatWaktuLengkap mengembalikan string waktu lengkap "02/01/2006 15:04:05".
func FormatWaktuLengkap(t time.Time) string {
	return t.Format("02/01/2006 15:04:05")
}

// TanggalFile mengembalikan string tanggal untuk nama file "YYYY-MM-DD".
func TanggalFile(t time.Time) string {
	return t.Format("2006-01-02")
}

// PadKanan menambahkan spasi di kanan string hingga mencapai panjang n.
func PadKanan(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat(" ", n-len(s))
}

// PadKiri menambahkan spasi di kiri string hingga mencapai panjang n.
func PadKiri(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return strings.Repeat(" ", n-len(s)) + s
}

// ── FILE SYSTEM ───────────────────────────────────────────────────────────────

// EnsureDir memastikan direktori ada, membuatnya jika belum ada (rekursif).
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
