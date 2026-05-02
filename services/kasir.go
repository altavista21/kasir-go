// Package services mengandung seluruh logika bisnis aplikasi kasir:
// manajemen keranjang, checkout, penyimpanan JSON, pembuatan struk file,
// dan laporan harian.
package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"kasir-go/models"
	"kasir-go/utils"
)

const (
	namaToko  = "TOKO SERBA ADA MAJU JAYA"
	alamatoko = "Jl. Contoh No. 1, Jakarta"
	lebarStruk = 44
)

// ── STRUCT KASIR ──────────────────────────────────────────────────────────────

// Kasir menyimpan state transaksi aktif: produk tersedia dan keranjang belanja.
type Kasir struct {
	Produk    []models.Barang
	Keranjang []models.ItemKeranjang
}

// NewKasir membuat instance Kasir baru dengan daftar produk default.
func NewKasir() *Kasir {
	return &Kasir{
		Produk:    models.DaftarBarang(),
		Keranjang: []models.ItemKeranjang{},
	}
}

// ── MANAJEMEN KERANJANG ───────────────────────────────────────────────────────

// CariBarang mencari produk berdasarkan kode (case-insensitive).
// Mengembalikan pointer ke Barang jika ditemukan, nil jika tidak.
func (k *Kasir) CariBarang(kode string) *models.Barang {
	for i, b := range k.Produk {
		if strings.EqualFold(b.Kode, kode) {
			return &k.Produk[i]
		}
	}
	return nil
}

// TambahKeKeranjang menambahkan item ke keranjang.
// Jika kode sudah ada, qty diakumulasi; jika baru, di-append.
func (k *Kasir) TambahKeKeranjang(kode string, jumlah int) error {
	barang := k.CariBarang(kode)
	if barang == nil {
		return fmt.Errorf("barang dengan kode '%s' tidak ditemukan", kode)
	}
	for i, item := range k.Keranjang {
		if strings.EqualFold(item.Barang.Kode, kode) {
			k.Keranjang[i].Jumlah += jumlah
			k.Keranjang[i].Subtotal = k.Keranjang[i].Barang.Harga * k.Keranjang[i].Jumlah
			fmt.Printf("  ✓  %s diperbarui → total %d pcs\n", barang.Nama, k.Keranjang[i].Jumlah)
			return nil
		}
	}
	item := models.ItemKeranjang{
		Barang:   *barang,
		Jumlah:   jumlah,
		Subtotal: barang.Harga * jumlah,
	}
	k.Keranjang = append(k.Keranjang, item)
	fmt.Printf("  ✓  %s x%d ditambahkan ke keranjang\n", barang.Nama, jumlah)
	return nil
}

// HapusDariKeranjang menghapus item dari keranjang berdasarkan kode barang.
func (k *Kasir) HapusDariKeranjang(kode string) error {
	for i, item := range k.Keranjang {
		if strings.EqualFold(item.Barang.Kode, kode) {
			nama := item.Barang.Nama
			k.Keranjang = append(k.Keranjang[:i], k.Keranjang[i+1:]...)
			fmt.Printf("  🗑  %s dihapus dari keranjang\n", nama)
			return nil
		}
	}
	return fmt.Errorf("barang dengan kode '%s' tidak ada di keranjang", kode)
}

// UpdateQtyKeranjang mengubah jumlah item di keranjang.
// Jika qty baru <= 0, item akan dihapus dari keranjang.
func (k *Kasir) UpdateQtyKeranjang(kode string, qtyBaru int) error {
	for i, item := range k.Keranjang {
		if strings.EqualFold(item.Barang.Kode, kode) {
			if qtyBaru <= 0 {
				return k.HapusDariKeranjang(kode)
			}
			k.Keranjang[i].Jumlah = qtyBaru
			k.Keranjang[i].Subtotal = item.Barang.Harga * qtyBaru
			fmt.Printf("  ✓  %s diperbarui → %d pcs\n", item.Barang.Nama, qtyBaru)
			return nil
		}
	}
	return fmt.Errorf("barang dengan kode '%s' tidak ada di keranjang", kode)
}

// ResetKeranjang mengosongkan keranjang belanja untuk transaksi baru.
func (k *Kasir) ResetKeranjang() {
	k.Keranjang = []models.ItemKeranjang{}
}

// KeranjangKosong memeriksa apakah keranjang tidak memiliki item.
func (k *Kasir) KeranjangKosong() bool {
	return len(k.Keranjang) == 0
}

// ── PERHITUNGAN ───────────────────────────────────────────────────────────────

// HitungTotal menjumlahkan semua subtotal item di keranjang.
func (k *Kasir) HitungTotal() int {
	total := 0
	for _, item := range k.Keranjang {
		total += item.Subtotal
	}
	return total
}

// HitungKembalian mengembalikan selisih antara uang bayar dan total belanja.
func HitungKembalian(total, bayar int) int {
	return bayar - total
}

// ── NOMOR TRANSAKSI ───────────────────────────────────────────────────────────

// GenerateIDTransaksi membuat ID transaksi format TRX-YYYYMMDD-XXX.
// Nomor urut dihitung berdasarkan jumlah transaksi yang sudah ada hari ini.
func GenerateIDTransaksi() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	filePath := filepath.Join("data", "transaksi-"+utils.TanggalFile(now)+".json")

	transaksiHariIni, _ := loadTransaksiDariFile(filePath)
	nomor := len(transaksiHariIni) + 1

	return fmt.Sprintf("TRX-%s-%03d", dateStr, nomor)
}

// ── PENYIMPANAN JSON ──────────────────────────────────────────────────────────

// SaveTransaksiToFile menyimpan data transaksi ke file JSON harian.
// Jika file sudah ada, data di-append (tidak overwrite).
func SaveTransaksiToFile(trx models.Transaksi) error {
	if err := utils.EnsureDir("data"); err != nil {
		return fmt.Errorf("gagal membuat folder data: %w", err)
	}

	filePath := filepath.Join("data", "transaksi-"+utils.TanggalFile(trx.Waktu)+".json")

	// Load existing data
	existing, _ := loadTransaksiDariFile(filePath)
	existing = append(existing, trx)

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Errorf("gagal encode JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("gagal simpan file: %w", err)
	}
	return nil
}

// loadTransaksiDariFile membaca slice Transaksi dari file JSON.
// Mengembalikan slice kosong jika file belum ada (bukan error).
func loadTransaksiDariFile(filePath string) ([]models.Transaksi, error) {
	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return []models.Transaksi{}, nil
	}
	if err != nil {
		return nil, err
	}
	var list []models.Transaksi
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// LoadTransaksiHarian membaca semua transaksi untuk tanggal hari ini.
func LoadTransaksiHarian() ([]models.Transaksi, error) {
	filePath := filepath.Join("data", "transaksi-"+utils.TanggalFile(time.Now())+".json")
	return loadTransaksiDariFile(filePath)
}

// ── STRUK FILE ────────────────────────────────────────────────────────────────

// GenerateStrukFile membuat dan menyimpan file struk .txt ke folder struk/.
func GenerateStrukFile(trx models.Transaksi) error {
	if err := utils.EnsureDir("struk"); err != nil {
		return fmt.Errorf("gagal membuat folder struk: %w", err)
	}

	filePath := filepath.Join("struk", trx.IDTransaksi+".txt")
	konten := buildKontenStruk(trx)

	if err := os.WriteFile(filePath, []byte(konten), 0644); err != nil {
		return fmt.Errorf("gagal simpan struk: %w", err)
	}
	fmt.Printf("  💾  Struk disimpan → %s\n", filePath)
	return nil
}

// buildKontenStruk menyusun string isi struk yang akan dicetak dan disimpan.
func buildKontenStruk(trx models.Transaksi) string {
	lebar := lebarStruk
	garis := strings.Repeat("-", lebar)
	garisTebal := strings.Repeat("=", lebar)
	sb := &strings.Builder{}

	// Header toko
	tulis := func(s string) { sb.WriteString(s + "\n") }
	tengah := func(s string) {
		pad := (lebar - len(s)) / 2
		if pad < 0 {
			pad = 0
		}
		tulis(strings.Repeat(" ", pad) + s)
	}

	tulis(garisTebal)
	tengah(namaToko)
	tengah(alamatoko)
	tulis(garis)
	tulis(fmt.Sprintf("No  : %s", trx.IDTransaksi))
	tulis(fmt.Sprintf("Tgl : %s", utils.FormatWaktuLengkap(trx.Waktu)))
	tulis(garisTebal)

	// Header kolom
	tulis(fmt.Sprintf("%-20s %5s %8s %9s", "Item", "Qty", "Harga", "Subtotal"))
	tulis(garis)

	// Baris item
	for _, item := range trx.Items {
		nama := item.Nama
		if len(nama) > 20 {
			nama = nama[:20]
		}
		tulis(fmt.Sprintf("%-20s %5d %8s %9s",
			nama,
			item.Qty,
			utils.PadKiri(utils.FormatRupiah(item.Harga), 8),
			utils.PadKiri(utils.FormatRupiah(item.Subtotal), 9),
		))
	}

	tulis(garis)
	tulis(fmt.Sprintf("%s%s", utils.PadKanan("TOTAL", lebar-14), utils.PadKiri(utils.FormatRupiah(trx.Total), 14)))
	tulis(fmt.Sprintf("%s%s", utils.PadKanan("BAYAR", lebar-14), utils.PadKiri(utils.FormatRupiah(trx.Bayar), 14)))
	tulis(fmt.Sprintf("%s%s", utils.PadKanan("KEMBALIAN", lebar-14), utils.PadKiri(utils.FormatRupiah(trx.Kembalian), 14)))
	tulis(garisTebal)
	tengah("Terima kasih telah berbelanja!")
	tengah("Selamat datang kembali :)")
	tulis(garisTebal)

	return sb.String()
}

// ── TAMPILAN TERMINAL ─────────────────────────────────────────────────────────

// TampilkanDaftarBarang mencetak tabel produk ke terminal.
func (k *Kasir) TampilkanDaftarBarang() {
	lebar := lebarStruk
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

// TampilkanKeranjang mencetak isi keranjang belanja saat ini ke terminal.
func (k *Kasir) TampilkanKeranjang() {
	if k.KeranjangKosong() {
		fmt.Println("  Keranjang masih kosong.")
		return
	}
	lebar := lebarStruk
	garis := strings.Repeat("─", lebar)
	fmt.Println()
	fmt.Println("  🛒  ISI KERANJANG:")
	fmt.Println("  " + garis)
	for i, item := range k.Keranjang {
		fmt.Printf("  %d. %-18s x%-3d = %s\n",
			i+1, item.Barang.Nama, item.Jumlah, utils.FormatRupiah(item.Subtotal))
	}
	fmt.Println("  " + garis)
	fmt.Printf("  %-24s       %s\n", "TOTAL", utils.FormatRupiah(k.HitungTotal()))
	fmt.Println()
}

// CetakStrukTerminal menampilkan struk pembelian di terminal setelah checkout.
func CetakStrukTerminal(trx models.Transaksi) {
	fmt.Println()
	fmt.Println(buildKontenStruk(trx))
}

// ── LAPORAN HARIAN ────────────────────────────────────────────────────────────

// TampilkanLaporanHarian membaca file JSON hari ini dan menampilkan ringkasan
// total pendapatan, jumlah transaksi, dan item terlaris.
func TampilkanLaporanHarian() {
	lebar := lebarStruk
	garisTebal := strings.Repeat("═", lebar)

	transaksi, err := LoadTransaksiHarian()

	fmt.Println()
	fmt.Println("╔" + garisTebal + "╗")
	fmt.Printf("║%-*s║\n", lebar, "  📊  LAPORAN HARIAN — "+utils.FormatTanggal(time.Now()))
	fmt.Println("╠" + garisTebal + "╣")

	if err != nil || len(transaksi) == 0 {
		fmt.Printf("║%-*s║\n", lebar, "  Belum ada transaksi hari ini.")
		fmt.Println("╚" + garisTebal + "╝")
		fmt.Println()
		return
	}

	// Hitung total pendapatan & frekuensi item
	totalPendapatan := 0
	frekuensi := make(map[string]int)
	for _, trx := range transaksi {
		totalPendapatan += trx.Total
		for _, item := range trx.Items {
			frekuensi[item.Nama] += item.Qty
		}
	}

	fmt.Printf("║  %-20s : %-18s║\n", "Jumlah Transaksi", fmt.Sprintf("%d transaksi", len(transaksi)))
	fmt.Printf("║  %-20s : %-18s║\n", "Total Pendapatan", utils.FormatRupiah(totalPendapatan))
	fmt.Println("╠" + strings.Repeat("─", lebar) + "╣")
	fmt.Printf("║%-*s║\n", lebar, "  🏆  Item Terlaris:")

	// Cari item dengan qty terbanyak
	terlaris := itemTerlaris(frekuensi)
	for i, entry := range terlaris {
		if i >= 5 {
			break
		}
		baris := fmt.Sprintf("  %d. %-20s → %d pcs", i+1, entry[0], entry[1])
		fmt.Printf("║%-*s║\n", lebar, baris)
	}

	fmt.Println("╠" + strings.Repeat("─", lebar) + "╣")
	fmt.Printf("║%-*s║\n", lebar, "  📋  Riwayat Transaksi:")
	for _, trx := range transaksi {
		baris := fmt.Sprintf("  %-22s  %s", trx.IDTransaksi, utils.FormatRupiah(trx.Total))
		fmt.Printf("║%-*s║\n", lebar, baris)
	}

	fmt.Println("╚" + garisTebal + "╝")
	fmt.Println()
}

// itemTerlaris mengurutkan map frekuensi menjadi slice [nama, qty] secara descending.
func itemTerlaris(freq map[string]int) [][2]interface{} {
	result := make([][2]interface{}, 0, len(freq))
	for nama, qty := range freq {
		result = append(result, [2]interface{}{nama, qty})
	}
	// Bubble sort sederhana (tanpa import sort)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j][1].(int) > result[i][1].(int) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
