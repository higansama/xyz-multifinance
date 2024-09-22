package dto

import (
	"time"

	"github.com/higansama/xyz-multi-finance/customer/errors"
)

type ReqCreateCustomer struct {
	ID         string    `json:"id"`
	NIK        string    `json:"nik"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	Password   string    `json:"password"`
	LegalName  string    `json:"legal_name"`
	BirthPlace string    `json:"birth_place"`
	BirthDate  time.Time `json:"birth_date"`
	Salary     int       `json:"salary"`
	PhotoID    string    `json:"photo_id"`
	Selfie     string    `json:"selfie"`
}

func (r ReqCreateCustomer) Validate() error {
	if r.NIK == "" || r.Password == "" {
		return errors.ErrorParams
	}

	if r.Salary < 0 {
		return errors.SalaryMinus
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r LoginRequest) Validate() error {
	// if r.Email == "" || r.Param == "email" {
	// 	return nil
	// }

	return nil
}

type LoginResponse struct {
	Token string `json:"token"`
}

type MyLimitResponse struct {
	Userid     string `json:"userid"`
	Email      string `json:"email"`
	OneMonth   int    `json:"one_month"`
	TwoMonth   int    `json:"two_month"`
	ThreeMonth int    `json:"three_month"`
	FourMonth  int    `json:"four_month"`
}

type CreateTransaksiRequest struct {
	NamaAsset string `json:"nama_asset" binding:"required"`
	HargaOTR  int    `json:"harga_otr" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Tenor     int    `json:"tenor" binding:"required"`
}

type CreateTransaksiResponse struct {
	NomerKontrak  string `json:"nomer_kontrak"`
	NamaAsset     string `json:"nama_asset"`
	HargaOTR      int    `json:"harga_otr"`
	AdminFee      int    `json:"admin_fee"`
	JumlahCicilan int    `json:"jumlah_cicilan"`
	CreatedAt     string `json:"created_at"`
}

type CreatePaymentRequest struct {
	NomerKontrak      string `json:"nomer_kontrak"`
	NominalPembayaran int    `json:"nominal_pembayaran"`
}

type CekPaymentContractResponse struct {
	CreateTransaksiResponse
	Payments []PaymentResponse `json:"payments"`
}

type PaymentResponse struct {
	IDBayar       string `json:"id_bayar"`
	NominalNormal string `json:"nominal_bayar"`
	Denda         string `json:"denda"`
	SisaCicilan   string `json:"sisa_cicilan"`
	CreatedAt     string `json:"created_at"`
}
