package gormadapter

import (
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	domain "github.com/higansama/xyz-multi-finance/customer"
	"github.com/higansama/xyz-multi-finance/internal/utils"
)

type Customer struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	NIK        string    `gorm:"type:varchar(191);unique"`
	Email      string    `gorm:"type:varchar(191);unique"`
	FullName   string    `gorm:"type:varchar(191)"`
	Salt       string    `gorm:"type:varchar(191)"`
	Password   string    `gorm:"type:text"`
	LegalName  string    `gorm:"type:varchar(191)"`
	BirthPlace string    `gorm:"type:varchar(191)"`
	BirthDate  time.Time `gorm:"type:datetime(3)"`
	Salary     int       `gorm:"type:int(11)"`
	PhotoID    string    `gorm:"type:varchar(191)"`
	Selfie     string    `gorm:"type:varchar(191)"`
	CreatedAt  time.Time `gorm:"type:datetime(3);autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:datetime(3);autoUpdateTime"`
}

// Hook to generate UUID before creating a record
func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

// convert
func (c *Customer) ConvertToCustomerDomain() domain.CustomerEntity {
	return domain.CustomerEntity{
		ID:         c.ID.String(),
		NIK:        c.NIK,
		FullName:   c.FullName,
		Salt:       c.Salt,
		Password:   c.Password,
		LegalName:  c.LegalName,
		BirthPlace: c.BirthPlace,
		BirthDate:  c.BirthDate,
		Salary:     c.Salary,
		PhotoID:    c.PhotoID,
		Selfie:     c.Selfie,
		Email:      c.Email,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

type LoanLimit struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Userid     uuid.UUID `gorm:"type:uuid;not null" gorm:"foreignKey:CustomerID"` // Foreign Key to Customer
	Customer   Customer  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OneMonth   float64
	TwoMonth   float64
	ThreeMonth float64
	FourMonth  float64
}

type CustomerWithLoanLimit struct {
	Customer  Customer  `gorm:"embeded"`
	LoanLimit LoanLimit `gorm:"embeded"`
}

type Transaksi struct {
	ID              uuid.UUID `gorm:"type:char(36);primaryKey"`
	NomerKontrak    string    `gorm:"type:varchar(50);not null;unique"`
	UserID          uuid.UUID `gorm:"type:char(36);not null"` // Foreign key ke Customer
	HargaOTR        int       `gorm:"type:int(11);not null"`
	Tenor           int       `gorm:"type:int(11);not null"`
	AdminFee        int       `gorm:"type:int(11);not null"`
	Type            string    `gorm:"type:varchar(11);not null"`
	JumlahCicilan   int       `gorm:"type:int(11);not null"`
	CicilanPerbulan int       `gorm:"type:int(11);not null"`
	NamaAsset       string    `gorm:"type:varchar(50);not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"` // Default timestamp current_timestamp()
	UpdatedAt       time.Time `gorm:"autoUpdateTime"` // Default timestamp current_timestamp() on update
	Customer        Customer  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Payments        []Payment `gorm:"foreignKey:IDTransaksi;references:ID"`
}

// Hook untuk generate UUID sebelum create
func (t *Transaksi) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New() // Generate UUID jika nil
	}
	if t.NomerKontrak == "" {
		t.NomerKontrak = strings.Join(utils.GenerateKontrakCode(), "") // Generate random Nomer Kontrak
	}
	return nil
}

func (t *Transaksi) ConvertToTransactionDomain() domain.Transaksi {
	return domain.Transaksi{
		ID:              t.ID.String(),
		NomerKontrak:    t.NomerKontrak,
		UserID:          t.UserID.String(),
		Tenor:           t.Tenor,
		HargaOTR:        t.HargaOTR,
		AdminFee:        t.AdminFee,
		JumlahCicilan:   t.JumlahCicilan,
		CicilanPerbulan: t.CicilanPerbulan,
		NamaAsset:       t.NamaAsset,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}

}

type Payment struct {
	ID             uuid.UUID `gorm:"type:varchar(50);primaryKey;column:id"`
	IDTransaksi    uuid.UUID `gorm:"type:varchar(50);column:id_transaksi"`
	IDUser         uuid.UUID `gorm:"type:varchar(50);column:id_user"`
	IdxTransaction int       `gorm:"type:int(11);column:idx_transaction"`
	NominalNormal  int       `gorm:"type:int(11);column:nominal_normal"`
	Denda          int       `gorm:"type:int(11)"`
	NominalBayar   int       `gorm:"type:int(11)"`
	Interest       int       `gorm:"type:int(11)"`
	// SisaCicilan    int       `gorm:"type:int(11);column:sisa_cicilan"`
	IsPaid     bool      `gorm:"type:tinyint(1);not null;default:0"`
	CreatedAt  time.Time `gorm:"type:timestamp;default:current_timestamp();column:created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp;default:current_timestamp() on update current_timestamp();column:updated_at"`
	JatuhTempo time.Time `gorm:"type:timestamp;default:current_timestamp() on update current_timestamp();column:jatuh_tempo"`
	PaidAt     time.Time `gorm:"type:timestamp;column:paid_at"`
	Transaksi  Transaksi `gorm:"foreignKey:IDTransaksi;references:ID"`
}

// Hook untuk generate UUID sebelum create
func (t *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New() // Generate UUID jika nil
	}
	return nil
}

func (p *Payment) ConvertToDomainPayment() domain.Payments {
	return domain.Payments{
		ID:            p.ID.String(),
		IDTransaksi:   p.IDTransaksi.String(),
		IDUser:        p.IDUser.String(),
		NominalNormal: p.NominalNormal,
		Denda:         p.Denda,
		NominalBayar:  p.NominalBayar,
		Interest:      p.Interest,
		// SisaCicilan:   p.SisaCicilan,
		IsPaid:     p.IsPaid,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		JatuhTempo: p.JatuhTempo,
		IdxPayment: p.IdxTransaction,
		PaidAt:     p.PaidAt,
	}
}
