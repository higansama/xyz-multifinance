package usecasev1

import (
	"errors"
	"time"

	domain "github.com/higansama/xyz-multi-finance/customer"
	"github.com/higansama/xyz-multi-finance/customer/dto"
	derror "github.com/higansama/xyz-multi-finance/customer/errors"
	"github.com/higansama/xyz-multi-finance/internal/utils"
)

// BuatTransaksi implements customer.Usecase.
func (c *CustomerUsecase) BuatTransaksi(userid string, payload dto.CreateTransaksiRequest) (response dto.CreateTransaksiResponse, err error) {
	// get user
	customer, err := c.repo.GetUser(userid)
	if err != nil {
		return response, err
	}

	// cek nilai produk apakah nilai produk ada diatas maksimal limit
	limit := customer.MyLoanLimit()
	if payload.HargaOTR > int(limit.FourMonth) {
		return response, derror.CicilanDiatasLimit
	}

	if payload.HargaOTR < int(limit.OneMonth) {
		return response, errors.New("minimal transaksi anda adalah " + utils.FormatRupiah(int(limit.OneMonth)))
	}

	if payload.Tenor > 5 {
		return response, derror.TenorNotAvailable
	}

	switch payload.Tenor {
	case 4:
		if payload.HargaOTR > int(limit.FourMonth) {
			return response, errors.New("Harga OTR diatas limit ke-4 , minimal cicilan anda di 4 bulan adalah " + utils.FormatRupiah(int(limit.FourMonth)))
		}
	case 3:
		if payload.HargaOTR > int(limit.ThreeMonth) {
			return response, errors.New("Harga OTR diatas limit ke-3 , minimal cicilan anda di 3 bulan adalah " + utils.FormatRupiah(int(limit.ThreeMonth)))
		}
	case 2:
		if payload.HargaOTR > int(limit.TwoMonth) {
			return response, errors.New("Harga OTR diatas limit ke-2 , minimal cicilan anda di 2 bulan adalah " + utils.FormatRupiah(int(limit.TwoMonth)))
		}
	case 1:
		if payload.HargaOTR > int(limit.OneMonth) {
			return response, errors.New("Harga OTR diatas limit ke-1 , minimal cicilan anda di 1 bulan adalah " + utils.FormatRupiah(int(limit.OneMonth)))
		}

	}

	adminFee := utils.CreateAdminFee(payload.HargaOTR)

	// total cicilan
	totlCicilan := payload.HargaOTR + adminFee
	// hitung cicilan perbulan
	cicilanPerbulan := float64(totlCicilan/payload.Tenor) + (float64(totlCicilan/payload.Tenor) * 0.02)

	// buat data transaksi
	transaksi := domain.Transaksi{
		UserID:          customer.ID,
		HargaOTR:        payload.HargaOTR,
		AdminFee:        adminFee,
		JumlahCicilan:   totlCicilan,
		CicilanPerbulan: int(cicilanPerbulan),
		Tenor:           payload.Tenor,
		NamaAsset:       payload.NamaAsset,
		CreatedAt:       time.Now(),
	}

	// buat transaksi
	transaksi, err = c.repo.CreateTransaksi(transaksi)
	if err != nil {
		return response, err
	}

	response = dto.CreateTransaksiResponse{
		NomerKontrak:  transaksi.NomerKontrak,
		NamaAsset:     transaksi.NamaAsset,
		HargaOTR:      transaksi.HargaOTR,
		AdminFee:      adminFee,
		JumlahCicilan: transaksi.JumlahCicilan,
		CreatedAt:     transaksi.CreatedAt.String(),
	}
	return response, nil
}
