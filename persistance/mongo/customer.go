package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID         primitive.ObjectID `bson:"_id"`
	NIK        string             `bson:"nik"`
	FullName   string             `bson:"full_name"`
	Salt       string             `bson:"salt"`
	Password   string             `bson:"password"`
	LegalName  string             `bson:"legal_name"`
	BirthPlace string             `bson:"birth_place"`
	BirthDate  time.Time          `bson:"birth_date"`
	Salary     float64            `bson:"salary"`
	PhotoID    string             `bson:"photo_id"`
	Selfie     string             `bson:"selfie"`
	LoanType   Loan               `bson:"loan_type"`
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

type Loan struct {
	TypeName string
	Max      int
	Tenor    int
}
