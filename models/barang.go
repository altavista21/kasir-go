// Package models mendefinisikan seluruh struktur data aplikasi kasir.
package models

import "time"

// ── PRODUK ────────────────────────────────────────────────────────────────────

// Barang merepresentasikan satu produk yang tersedia di toko.
type Barang struct {
	Kode  string
	Nama  string
	Harga int
}

// ItemKeranjang merepresentasikan satu baris dalam keranjang belanja.
type ItemKeranjang struct {
	Barang   Barang
	Jumlah   int
	Subtotal int
}

// DaftarBarang mengembalikan semua produk yang tersedia (data hardcode).
func DaftarBarang() []Barang {
	return []Barang{
		{Kode: "001", Nama: "Indomie Goreng", Harga: 3000},
		{Kode: "002", Nama: "Aqua Botol", Harga: 3000},
		{Kode: "003", Nama: "Roti Tawar", Harga: 5000},
		{Kode: "004", Nama: "Teh Botol", Harga: 4000},
		{Kode: "005", Nama: "Kopi Sachet", Harga: 2000},
		{Kode: "006", Nama: "Kerupuk", Harga: 1500},
	}
}

// ── TRANSAKSI ─────────────────────────────────────────────────────────────────

// ItemTransaksi adalah snapshot satu baris item saat transaksi berlangsung,
// digunakan untuk disimpan ke JSON (tidak menyimpan pointer ke Barang).
type ItemTransaksi struct {
	Nama     string `json:"nama"`
	Qty      int    `json:"qty"`
	Harga    int    `json:"harga"`
	Subtotal int    `json:"subtotal"`
}

// Transaksi menyimpan seluruh data satu sesi pembelian yang sudah selesai.
type Transaksi struct {
	IDTransaksi string          `json:"id_transaksi"`
	Waktu       time.Time       `json:"waktu"`
	Items       []ItemTransaksi `json:"items"`
	Total       int             `json:"total"`
	Bayar       int             `json:"bayar"`
	Kembalian   int             `json:"kembalian"`
}
