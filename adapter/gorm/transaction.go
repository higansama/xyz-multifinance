package adapter

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"` // UUID sebagai primary key
	CustomerID  uint      `gorm:"not null"`                                         // Foreign key ke tabel Customer
	ContractNo  string    `gorm:"not null"`                                         // Nomor kontrak transaksi
	OTR         float64   `gorm:"not null"`                                         // Harga On The Road
	AdminFee    float64   `gorm:"not null"`                                         // Biaya admin
	Installment float64   `gorm:"not null"`                                         // Jumlah cicilan
	Interest    float64   `gorm:"not null"`                                         // Bunga transaksi
	AssetName   string    `gorm:"not null"`                                         // Nama aset
	CreatedAt   time.Time // Waktu pembuatan transaksi
	UpdatedAt   time.Time // Waktu pembaruan transaksi
}
