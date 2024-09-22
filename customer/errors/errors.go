package errors

import ierrors "github.com/higansama/xyz-multi-finance/internal/errors"

var ErrorParams = ierrors.NewDomainError("Customer", "invalid parameter")
var SalaryMinus = ierrors.NewDomainError("Customer", "salary minus")
var EmailAtauNIKTerdaftar = ierrors.NewDomainError("Customer", "Email atau NIK terdaftar")
var UserNotFound = ierrors.NewDomainError("Customer", "user not found")
var LoginFailed = ierrors.NewDomainError("customer", "your username or password is invalid")
var CicilanDiatasLimit = ierrors.NewDomainError("customer", "cicilan diatas limit")
var TenorNotAvailable = ierrors.NewDomainError("customer", "tenor not available")
var PembayaranKurang = ierrors.NewDomainError("installment", "pembayaran tidak sesuai")
