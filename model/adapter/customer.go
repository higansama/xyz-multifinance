package entity

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"` // UUID sebagai primary key
	NIK        string    `gorm:"unique;not null"`
	FullName   string
	Salt       string
	Password   string
	LegalName  string
	BirthPlace string
	BirthDate  time.Time
	Salary     float64
	PhotoID    string
	Selfie     string
	// LoanType   Loan
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
