# 🛒 Kasir-Go

Aplikasi kasir sederhana berbasis CLI (terminal) menggunakan Golang.

## Fitur
- Tampil daftar barang
- Input barang ke keranjang
- Akumulasi qty jika barang sama
- Hitung total otomatis
- Cetak struk pembelian

## Struktur Project
kasir-go/
├── main.go
├── go.mod
├── models/
│   └── barang.go
├── services/
│   └── kasir.go
└── utils/
└── helper.go
## Cara Menjalankan
```bash
git clone https://github.com/altavista21/kasir-go.git
cd kasir-go
go run main.go

Teknologi
Go (Golang)
CLI / Terminal based
In-memory data (tanpa database)
