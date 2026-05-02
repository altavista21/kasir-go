// Package services mengandung logika bisnis utama aplikasi kasir,
// termasuk manajemen keranjang belanja dan pencetakan struk.
package services

import (
	"fmt"
	"strings"

	"kasir-go/models"
	"kasir-go/utils"
)

// Kasir adalah struct utama yang menyimpan state transaksi:
// daftar produk yang tersedia dan isi keranjang belanja saat ini.
type Kasir struct {
	Produk    []models.Barang
	Keranjang []models.ItemKeranjang
}

// NewKasir membuat instance Kasir baru dengan daftar produk default yang di-load
// dari models.DaftarBarang().
func NewKasir() *Kasir {
	return &Kasir{
		Produk:    models.DaftarBarang(),
		Keranjang: []models.ItemKeranjang{},
	}
}

// CariBarang mencari produk berdasarkan kode barang (case-insensitive).
// Mengembalikan pointer ke Barang jika ditemukan, atau nil jika tidak.
func (k *Kasir) CariBarang(kode string) *models.Barang {
	for i, b := range k.Produk {
		if strings.EqualFold(b.Kode, kode) {
			return &k.Produk[i]
		}
	}
	return nil
}

// TambahKeKeranjang menambahkan barang ke keranjang berdasarkan kode dan jumlah.
// Jika barang dengan kode yang sama sudah ada, jumlahnya diakumulasi (tidak duplikat).
// Mengembalikan error jika kode barang tidak ditemukan.
func (k *Kasir) TambahKeKeranjang(kode string, jumlah int) error {
	barang := k.CariBarang(kode)
	if barang == nil {
		return fmt.Errorf("barang dengan kode '%s' tidak ditemukan", kode)
	}

	// Cek apakah barang sudah ada di keranjang → akumulasi qty
	for i, item := range k.Keranjang {
		if strings.EqualFold(item.Barang.Kode, kode) {
			k.Keranjang[i].Jumlah += jumlah
			k.Keranjang[i].Subtotal = k.Keranjang[i].Barang.Harga * k.Keranjang[i].Jumlah
			fmt.Printf("  ✓  %s diperbarui: total %d pcs\n", barang.Nama, k.Keranjang[i].Jumlah)
			return nil
		}
	}

	// Barang baru → append ke keranjang
	item := models.ItemKeranjang{
		Barang:   *barang,
		Jumlah:   jumlah,
		Subtotal: barang.Harga * jumlah,
	}
	k.Keranjang = append(k.Keranjang, item)
	fmt.Printf("  ✓  %s x%d ditambahkan ke keranjang\n", barang.Nama, jumlah)
	return nil
}

// HitungTotal menjumlahkan semua subtotal item yang ada di keranjang
// dan mengembalikan total keseluruhan transaksi.
func (k *Kasir) HitungTotal() int {
	total := 0
	for _, item := range k.Keranjang {
		total += item.Subtotal
	}
	return total
}

// KeranjangKosong memeriksa apakah keranjang belanja tidak memiliki item.
func (k *Kasir) KeranjangKosong() bool {
	return len(k.Keranjang) == 0
}

// TampilkanDaftarBarang mencetak semua produk yang tersedia ke terminal
// dalam format tabel yang mudah dibaca.
func (k *Kasir) TampilkanDaftarBarang() {
	lebar := 44
	garis := strings.Repeat("─", lebar)

	fmt.Println()
	fmt.Println("╔" + strings.Repeat("═", lebar) + "╗")
	fmt.Printf("║%-*s║\n", lebar, "  📦  DAFTAR BARANG TERSEDIA")
	fmt.Println("╠" + strings.Repeat("═", lebar) + "╣")
	fmt.Printf("║  %-6s  %-18s  %-10s  ║\n", "KODE", "NAMA BARANG", "HARGA")
	fmt.Println("╠" + garis + "╣")

	for _, b := range k.Produk {
		fmt.Printf("║  %-6s  %-18s  %-10s  ║\n",
			b.Kode, b.Nama, utils.FormatRupiah(b.Harga))
	}

	fmt.Println("╚" + strings.Repeat("═", lebar) + "╝")
	fmt.Println()
}

// CetakStruk mencetak struk pembelian lengkap ke terminal, mencakup
// daftar item, subtotal tiap item, dan total keseluruhan transaksi.
func (k *Kasir) CetakStruk() {
	lebar := 44
	garis := strings.Repeat("─", lebar)
	garisTebal := strings.Repeat("═", lebar)

	fmt.Println()
	fmt.Println("╔" + garisTebal + "╗")
	fmt.Printf("║%-*s║\n", lebar, "       🧾  STRUK PEMBELIAN")
	fmt.Println("╠" + garisTebal + "╣")

	for i, item := range k.Keranjang {
		baris := fmt.Sprintf("  %d. %s x%d", i+1, item.Barang.Nama, item.Jumlah)
		harga := utils.FormatRupiah(item.Subtotal)
		// Padding dinamis agar harga rata kanan
		padding := lebar - len(baris) - len(harga) - 1
		if padding < 1 {
			padding = 1
		}
		fmt.Printf("║%s%s%s║\n", baris, strings.Repeat(" ", padding), harga)
	}

	fmt.Println("╠" + garis + "╣")
	totalStr := utils.FormatRupiah(k.HitungTotal())
	label := "  TOTAL"
	padding := lebar - len(label) - len(totalStr) - 1
	fmt.Printf("║%s%s%s║\n", label, strings.Repeat(" ", padding), totalStr)
	fmt.Println("╚" + garisTebal + "╝")

	fmt.Println()
	fmt.Println("  Terima kasih telah berbelanja! 🙏")
	fmt.Println()
}
