package usecasev1

import (
	"github.com/higansama/xyz-multi-finance/customer/dto"
	"github.com/higansama/xyz-multi-finance/internal/utils"
)

// CekKontrak implements customer.Usecase.
func (c *CustomerUsecase) CekKontrak(userid string, idkontrak string) (result dto.CekPaymentContractResponse, err error) {
	dbresult, err := c.repo.CekTransaksiByContract(userid, idkontrak)
	if err != nil {
		return result, err
	}

	payments := make([]dto.PaymentResponse, 0)
	if len(dbresult.Payments) > 0 {
		for _, v := range dbresult.Payments {
			payments = append(payments, dto.PaymentResponse{
				IDBayar:       v.ID,
				NominalNormal: utils.FormatRupiah(v.NominalBayar),
				Denda:         utils.FormatRupiah(v.Denda),
				// SisaCicilan:   utils.FormatRupiah(v.SisaCicilan),
				CreatedAt: v.CreatedAt.String(),
			})
		}
	}

	result = dto.CekPaymentContractResponse{
		CreateTransaksiResponse: dto.CreateTransaksiResponse{
			NomerKontrak:  dbresult.NomerKontrak,
			NamaAsset:     dbresult.NamaAsset,
			HargaOTR:      dbresult.HargaOTR,
			AdminFee:      dbresult.AdminFee,
			JumlahCicilan: dbresult.JumlahCicilan,
			CreatedAt:     dbresult.CreatedAt.String(),
		},
		Payments: payments,
	}
	return result, nil
}
