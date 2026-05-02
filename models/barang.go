// Package models mendefinisikan struktur data utama aplikasi kasir.
package models

// Barang merepresentasikan satu produk yang tersedia di toko.
type Barang struct {
	Kode  string
	Nama  string
	Harga int
}

// ItemKeranjang merepresentasikan satu baris dalam keranjang belanja,
// berisi data barang beserta jumlah yang dibeli.
type ItemKeranjang struct {
	Barang   Barang
	Jumlah   int
	Subtotal int
}

// DaftarBarang mengembalikan slice berisi semua produk yang tersedia (hardcode).
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
