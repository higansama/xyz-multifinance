package customer

import (
	"time"

	"github.com/higansama/xyz-multi-finance/internal/auth"
)

type CustomerEntity struct {
	ID         string
	NIK        string
	FullName   string
	Salt       string
	Password   string
	LegalName  string
	BirthPlace string
	BirthDate  time.Time
	Salary     float64
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
