// Kasir-Go: Aplikasi kasir sederhana berbasis CLI.
// Menjalankan orkestrasi alur belanja dari input barang hingga cetak struk.
package main

import (
	"fmt"

	"kasir-go/services"
	"kasir-go/utils"
)

func main() {
	k := services.NewKasir()

	cetakHeader()
	k.TampilkanDaftarBarang()
	jalankanSesiPembelian(k)
}

// cetakHeader menampilkan banner selamat datang di awal program.
func cetakHeader() {
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║         🛒  KASIR-GO  •  Point of Sale        ║")
	fmt.Println("║         Sistem Kasir Sederhana (CLI)          ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
}

// jalankanSesiPembelian mengelola loop utama pembelian:
// input kode barang, jumlah, dan konfirmasi lanjut/selesai.
// Setelah selesai, struk dicetak ke terminal.
func jalankanSesiPembelian(k *services.Kasir) {
	fmt.Println("─── Mulai Input Pembelian ───────────────────────")
	fmt.Println()

	for {
		// Input kode barang dengan validasi
		kode := inputKodeBarang(k)

		// Input jumlah dengan validasi angka positif
		jumlah := utils.BacaInt("  Jumlah       : ")

		// Tambahkan ke keranjang
		err := k.TambahKeKeranjang(kode, jumlah)
		if err != nil {
			// Seharusnya tidak sampai sini karena sudah divalidasi, tapi jaga-jaga
			fmt.Println(" ⚠  Error:", err)
		}

		fmt.Println()

		// Tanya apakah ingin menambah item lagi
		lanjut := utils.TanyaYaTidak("  Tambah barang lagi? (y/n): ")
		fmt.Println()
		if !lanjut {
			break
		}
	}

	// Pastikan keranjang tidak kosong sebelum cetak struk
	if k.KeranjangKosong() {
		fmt.Println("  Keranjang kosong. Tidak ada transaksi.")
		return
	}

	k.CetakStruk()
}

// inputKodeBarang meminta user memasukkan kode barang dan memvalidasinya.
// Jika kode tidak ditemukan, user diminta mengulang hingga kode valid dimasukkan.
func inputKodeBarang(k *services.Kasir) string {
	for {
		kode := utils.BacaString("  Kode Barang  : ")
		if kode == "" {
			fmt.Println("  ⚠  Kode barang tidak boleh kosong.")
			continue
		}
		if k.CariBarang(kode) == nil {
			fmt.Printf("  ⚠  Kode '%s' tidak ditemukan. Cek daftar barang di atas.\n", kode)
			continue
		}
		return kode
	}
}
