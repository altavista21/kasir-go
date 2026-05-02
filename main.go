// Kasir-Go v2: Sistem kasir CLI dengan penyimpanan JSON, struk file, dan laporan.
package main

import (
	"fmt"
	"strings"
	"time"

	"kasir-go/models"
	"kasir-go/services"
	"kasir-go/utils"
)

func main() {
	cetakHeader()
	menuUtama()
}

// cetakHeader menampilkan banner selamat datang saat program dimulai.
func cetakHeader() {
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║       🛒  KASIR-GO v2  •  Point of Sale       ║")
	fmt.Println("║         Sistem Kasir Sederhana (CLI)          ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
}

// menuUtama menampilkan menu pilihan dan mengarahkan ke fitur yang dipilih.
func menuUtama() {
	for {
		fmt.Println()
		fmt.Println("  ┌─────────────────────────────────┐")
		fmt.Println("  │         MENU UTAMA              │")
		fmt.Println("  ├─────────────────────────────────┤")
		fmt.Println("  │  1. Transaksi Baru              │")
		fmt.Println("  │  2. Lihat Laporan Hari Ini      │")
		fmt.Println("  │  3. Keluar                      │")
		fmt.Println("  └─────────────────────────────────┘")

		pilihan := utils.BacaString("  Pilih menu (1-3): ")

		switch pilihan {
		case "1":
			jalankanTransaksiBaru()
		case "2":
			services.TampilkanLaporanHarian()
		case "3":
			fmt.Println()
			fmt.Println("  Sampai jumpa! 👋")
			fmt.Println()
			return
		default:
			fmt.Println("  ⚠  Pilihan tidak valid. Masukkan 1, 2, atau 3.")
		}
	}
}

// jalankanTransaksiBaru mengelola satu sesi transaksi penuh:
// input barang → edit keranjang → checkout → simpan → cetak struk.
func jalankanTransaksiBaru() {
	k := services.NewKasir()

	k.TampilkanDaftarBarang()
	fmt.Println("─── Input Pembelian ─────────────────────────────")
	fmt.Println()

	// Loop input item ke keranjang
	for {
		menuKeranjang(k)

		lanjut := utils.TanyaYaTidak("  Tambah barang lagi? (y/n): ")
		fmt.Println()
		if !lanjut {
			break
		}
	}

	// Keranjang harus ada isinya
	if k.KeranjangKosong() {
		fmt.Println("  Keranjang kosong. Transaksi dibatalkan.")
		return
	}

	k.TampilkanKeranjang()

	// Tawari edit keranjang sebelum checkout
	if utils.TanyaYaTidak("  Edit keranjang (hapus/ubah qty)? (y/n): ") {
		menuEditKeranjang(k)
	}

	if k.KeranjangKosong() {
		fmt.Println("  Keranjang kosong setelah diedit. Transaksi dibatalkan.")
		return
	}

	// Proses checkout
	checkout(k)
}

// menuKeranjang menangani input satu item: kode barang dan jumlah.
func menuKeranjang(k *services.Kasir) {
	kode := inputKodeBarang(k)
	jumlah := utils.BacaInt("  Jumlah       : ")

	if err := k.TambahKeKeranjang(kode, jumlah); err != nil {
		fmt.Println("  ⚠  Error:", err)
	}
	fmt.Println()
}

// menuEditKeranjang menampilkan sub-menu untuk hapus item atau ubah qty.
func menuEditKeranjang(k *services.Kasir) {
	for {
		k.TampilkanKeranjang()
		fmt.Println("  ┌───────────────────────────┐")
		fmt.Println("  │  EDIT KERANJANG           │")
		fmt.Println("  ├───────────────────────────┤")
		fmt.Println("  │  1. Hapus item            │")
		fmt.Println("  │  2. Ubah jumlah item      │")
		fmt.Println("  │  3. Selesai edit          │")
		fmt.Println("  └───────────────────────────┘")

		pilihan := utils.BacaString("  Pilih (1-3): ")

		switch pilihan {
		case "1":
			kode := utils.BacaString("  Kode barang yang dihapus: ")
			if err := k.HapusDariKeranjang(kode); err != nil {
				fmt.Println("  ⚠ ", err)
			}
		case "2":
			kode := utils.BacaString("  Kode barang yang diubah: ")
			qty := utils.BacaInt("  Jumlah baru             : ")
			if err := k.UpdateQtyKeranjang(kode, qty); err != nil {
				fmt.Println("  ⚠ ", err)
			}
		case "3":
			return
		default:
			fmt.Println("  ⚠  Masukkan 1, 2, atau 3.")
		}
		fmt.Println()

		if k.KeranjangKosong() {
			fmt.Println("  Keranjang sudah kosong.")
			return
		}
	}
}

// checkout menangani proses pembayaran, penyimpanan data, dan pencetakan struk.
func checkout(k *services.Kasir) {
	total := k.HitungTotal()
	fmt.Println("─── Checkout ────────────────────────────────────")
	fmt.Printf("  Total belanja  : %s\n", utils.FormatRupiah(total))

	bayar := utils.BacaIntMin("  Uang bayar    : Rp ", total)
	kembalian := services.HitungKembalian(total, bayar)

	fmt.Printf("  Kembalian      : %s\n\n", utils.FormatRupiah(kembalian))

	// Susun data transaksi
	idTrx := services.GenerateIDTransaksi()
	now := time.Now()

	var items []models.ItemTransaksi
	for _, item := range k.Keranjang {
		items = append(items, models.ItemTransaksi{
			Nama:     item.Barang.Nama,
			Qty:      item.Jumlah,
			Harga:    item.Barang.Harga,
			Subtotal: item.Subtotal,
		})
	}

	trx := models.Transaksi{
		IDTransaksi: idTrx,
		Waktu:       now,
		Items:       items,
		Total:       total,
		Bayar:       bayar,
		Kembalian:   kembalian,
	}

	// Simpan ke JSON
	if err := services.SaveTransaksiToFile(trx); err != nil {
		fmt.Println("  ⚠  Gagal simpan transaksi:", err)
	} else {
		fmt.Println("  💾  Data transaksi disimpan ke data/")
	}

	// Simpan struk ke file
	if err := services.GenerateStrukFile(trx); err != nil {
		fmt.Println("  ⚠  Gagal simpan struk:", err)
	}

	// Cetak struk ke terminal
	services.CetakStrukTerminal(trx)

	fmt.Println("  " + strings.Repeat("─", 44))
	fmt.Println("  Transaksi selesai. Kembali ke menu utama...")
}

// inputKodeBarang meminta dan memvalidasi kode barang dari pengguna.
// Loop sampai kode yang dimasukkan valid.
func inputKodeBarang(k *services.Kasir) string {
	for {
		kode := utils.BacaString("  Kode Barang  : ")
		if kode == "" {
			fmt.Println("  ⚠  Kode tidak boleh kosong.")
			continue
		}
		if k.CariBarang(kode) == nil {
			fmt.Printf("  ⚠  Kode '%s' tidak ditemukan.\n", kode)
			continue
		}
		return kode
	}
}
