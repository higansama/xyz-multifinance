package customer

import (
	"time"

	"github.com/higansama/xyz-multi-finance/customer/errors"
	"github.com/higansama/xyz-multi-finance/internal/auth"
	"github.com/higansama/xyz-multi-finance/utils"
)

type CustomerEntity struct {
	ID         string
	NIK        string
	FullName   string
	Email      string
	Salt       string
	Password   string
	LegalName  string
	BirthPlace string
	BirthDate  time.Time
	Salary     int
	PhotoID    string
	Selfie     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (c *CustomerEntity) HashPassword(rawpassword string) {
	c.Salt = auth.GenerateSalt()
	hashpassword := auth.HashPassword(rawpassword, auth.HashOptions{
		Salt: c.Salt,
		Type: "pbkdf2",
	})
	c.Password = hashpassword
}

func (c *CustomerEntity) ValidatePassword(rawpassword string) error {
	valid := auth.VerifyPassword(rawpassword, auth.CompareOptions{
		Hash: c.Password,
		Salt: c.Salt,
		Type: "pbkdf2",
	})
	if !valid {
		return errors.LoginFailed
	}
	return nil
}

func (c *CustomerEntity) MyLoanLimit() LoanLimit {
	// Rasio cicilan terhadap gaji, berdasarkan data Annisa
	oneMonthRatio := 0.3    // 33.3%
	twoMonthsRatio := 0.4   // 40%
	threeMonthsRatio := 0.5 // 50%
	fourMonthsRatio := 0.6

	oneMonthLimit := oneMonthRatio * float64(c.Salary)
	twoMonthsLimit := utils.RoundToNearest500k(twoMonthsRatio * float64(c.Salary))
	threeMonthsLimit := utils.RoundToNearest500k(threeMonthsRatio * float64(c.Salary))
	fourMonthsLimit := utils.RoundToNearest500k(fourMonthsRatio * float64(c.Salary))
	if twoMonthsLimit == oneMonthLimit {
		twoMonthsLimit += 500000
	}

	if threeMonthsLimit == twoMonthsLimit {
		threeMonthsLimit += 500000
	}

	if fourMonthsLimit == threeMonthsLimit {
		fourMonthsLimit += 500000
	}

	return LoanLimit{
		Userid:     c.ID,
		OneMonth:   oneMonthLimit,
		TwoMonth:   twoMonthsLimit,
		ThreeMonth: threeMonthsLimit,
		FourMonth:  fourMonthsLimit,
	}
}

type LoanLimit struct {
	Userid     string
	Email      string
	OneMonth   float64
	TwoMonth   float64
	ThreeMonth float64
	FourMonth  float64
}

type Transaksi struct {
	ID              string
	NomerKontrak    string
	UserID          string
	Type            string
	Tenor           int
	HargaOTR        int
	AdminFee        int
	JumlahCicilan   int
	CicilanPerbulan int
	NamaAsset       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (t *Transaksi) CreatePayment(NominalBayar, idxPayment int) Payments {
	// denda := 0

	return Payments{
		IDUser:        t.UserID,
		NominalNormal: NominalBayar,
		// Denda:         ,
		NominalBayar: NominalBayar,
		// SisaCicilan:  0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func (t *Transaksi) MontlhlyInterest() int {
	return int(float64(t.CicilanPerbulan) * 0.02)
}

type Payments struct {
	ID            string
	IDTransaksi   string
	IDUser        string
	NominalNormal int
	Denda         int
	Interest      int
	NominalBayar  int
	// SisaCicilan   int
	IdxPayment int
	IsPaid     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	JatuhTempo time.Time
	PaidAt     time.Time
}

func (p *Payments) HitungDenda() int {
	var denda interface{}
	if time.Now().After(p.JatuhTempo) {
		jumlahHariTerlewat := p.JatuhTempo.Sub(time.Now()) / 24
		if (jumlahHariTerlewat - 30) > 0 {
			denda = float64(p.NominalNormal) * float64(0.05) * float64(jumlahHariTerlewat)
		}
	}
	if denda != nil {
		return denda.(int)
	}
	return 0
}

type TransaksiAnPayments struct {
	Transaksi
	Payments []Payments `json:"payments"`
}

type UpdatePayment struct {
	PaymentID    string
	NominalBayar int
	Denda        int
	IsPaid       bool
	PaidAt       time.Time
	Idx          int
}
