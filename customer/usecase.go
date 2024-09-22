package customer

import (
	"github.com/higansama/xyz-multi-finance/customer/dto"
)

type Usecase interface {
	Create(data dto.ReqCreateCustomer) error
	Login(payload dto.LoginRequest) (response dto.LoginResponse, err error)
	GetMyLimit(userid string) (response dto.MyLimitResponse, err error)
	BuatTransaksi(email string, payload dto.CreateTransaksiRequest) (response dto.CreateTransaksiResponse, err error)
	CekKontrak(userid, idkontrak string) (result dto.CekPaymentContractResponse, err error)
	CreatePayment(userid string, payload dto.CreatePaymentRequest) (err error)
}
