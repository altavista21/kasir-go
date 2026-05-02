// Package utils menyediakan fungsi-fungsi pembantu (helper) untuk
// keperluan input/output dan validasi di seluruh aplikasi.
package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// reader adalah scanner global agar tidak kehilangan buffer antar pemanggilan.
var reader = bufio.NewReader(os.Stdin)

// BacaString membaca satu baris input dari terminal dan mengembalikannya
// sebagai string yang sudah di-trim whitespace-nya.
func BacaString(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// BacaInt membaca input integer dari terminal. Jika input tidak valid,
// pengguna akan terus diminta mengulang hingga memasukkan angka yang benar.
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

// TanyaYaTidak menampilkan pertanyaan ya/tidak dan mengembalikan true jika
// pengguna menjawab 'y' atau 'Y', false jika 'n' atau 'N'.
// Input selain itu akan meminta pengguna mengulang.
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

// FormatRupiah memformat angka integer menjadi string mata uang Rupiah
// dengan pemisah ribuan, contoh: 15000 → "Rp 15.000".
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
