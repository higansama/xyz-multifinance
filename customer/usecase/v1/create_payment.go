package usecasev1

import (
	"errors"
	"time"

	domain "github.com/higansama/xyz-multi-finance/customer"
	derror "github.com/higansama/xyz-multi-finance/customer/errors"

	"github.com/higansama/xyz-multi-finance/customer/dto"
)

// CreatePayment implements customer.Usecase.
func (c *CustomerUsecase) CreatePayment(userid string, payload dto.CreatePaymentRequest) (err error) {
	transaction, err := c.repo.CekTransaksiByContract(userid, payload.NomerKontrak)
	if err != nil {
		return err
	}

	// get first payment
	paymentToUpdate, err := c.repo.GetFirstUnpaidPayment(transaction.NomerKontrak)
	if err != nil {
		return err
	}

	if paymentToUpdate.ID == "" {
		return errors.New("tidak ada transaksi untuk dibayar")
	}

	// validasi
	denda := paymentToUpdate.HitungDenda()

	if payload.NominalPembayaran != paymentToUpdate.NominalNormal+denda {
		return derror.PembayaranKurang
	}

	// update the payment

	updateData := domain.UpdatePayment{
		PaymentID:    paymentToUpdate.ID,
		NominalBayar: payload.NominalPembayaran,
		Denda:        denda,
		IsPaid:       true,
		PaidAt:       time.Now(),
		Idx:          paymentToUpdate.IdxPayment,
	}
	c.repo.PayTransaction(updateData)

	return nil
}
