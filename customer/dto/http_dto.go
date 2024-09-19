package dto

import "time"

type CreateCustomer struct {
	ID         string    `json:"id"`
	NIK        string    `json:"nik"`
	FullName   string    `json:"full_name"`
	Password   string    `json:"password"`
	LegalName  string    `json:"legal_name"`
	BirthPlace string    `json:"birth_place"`
	BirthDate  time.Time `json:"birth_date"`
	Salary     float64   `json:"salary"`
	PhotoID    string    `json:"photo_id"`
	Selfie     string    `json:"selfie"`
}
