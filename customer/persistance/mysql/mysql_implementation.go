package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	domain "github.com/higansama/xyz-multi-finance/customer"
	gormadapter "github.com/higansama/xyz-multi-finance/persistance/gorm"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

// GetFirstUnpaidPayment implements customer.Repository.
func (r *customerRepository) GetFirstUnpaidPayment(nomerKontrak string) (domain.Payments, error) {
	var Payment gormadapter.Payment

	err := r.db.Joins("JOIN transaksis ON transaksis.id = payments.id_transaksi").
		Where("transaksis.nomer_kontrak = ? AND payments.is_paid = ?", nomerKontrak, false).
		Order("payments.idx_transaction ASC"). // Urutkan berdasarkan waktu pembayaran (untuk mendapatkan yang pertama)
		First(&Payment).Error

	if err != nil {
		return Payment.ConvertToDomainPayment(), err
	}
	fmt.Println(Payment.IdxTransaction, "\t", Payment.ID)
	return Payment.ConvertToDomainPayment(), nil
}

// CreatePayment implements customer.Repository.
func (r *customerRepository) PayTransaction(payload domain.UpdatePayment) error {
	dataToUpdate := map[string]interface{}{
		"is_paid":       true,
		"nominal_bayar": payload.NominalBayar,
		"denda":         payload.Denda,
		"updated_at":    time.Now(),
		"paid_at":       time.Now(),
	}

	err := r.db.Model(&gormadapter.Payment{}).
		Where("id = ?", payload.PaymentID). // Kondisi berdasarkan ID payment
		Updates(dataToUpdate).Error

	return err
}

// CekTransaksiByContract implements customer.Repository.
func (r *customerRepository) CekTransaksiByContract(userid string, contractnumber string) (result domain.TransaksiAnPayments, err error) {
	var transaksi gormadapter.Transaksi
	result = domain.TransaksiAnPayments{}
	// Mengambil Transaksi dengan Payments yang terkait
	err = r.db.Preload("Payments").First(&transaksi, "nomer_kontrak = ?", contractnumber).Error
	if err != nil {
		return result, err
	}

	result.Transaksi = domain.Transaksi{
		ID:              transaksi.ID.String(),
		NomerKontrak:    transaksi.NomerKontrak,
		UserID:          userid,
		Tenor:           transaksi.Tenor,
		HargaOTR:        transaksi.HargaOTR,
		AdminFee:        transaksi.AdminFee,
		JumlahCicilan:   transaksi.JumlahCicilan,
		CicilanPerbulan: transaksi.CicilanPerbulan,
		NamaAsset:       transaksi.NamaAsset,
		CreatedAt:       transaksi.CreatedAt,
		UpdatedAt:       transaksi.UpdatedAt,
	}
	payments := make([]domain.Payments, 0)
	if len(transaksi.Payments) > 0 {
		for _, v := range transaksi.Payments {
			payments = append(payments, domain.Payments{
				ID:            v.ID.String(),
				IDUser:        v.IDUser.String(),
				NominalNormal: v.NominalNormal,
				Denda:         v.Denda,
				NominalBayar:  v.NominalBayar,
				// Interest: ,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			})
			result.Payments = payments
		}
	}

	return result, nil
}

// CreateTraksi implements customer.Repository.
func (r *customerRepository) CreateTransaksi(data domain.Transaksi) (domain.Transaksi, error) {
	transaksiData := gormadapter.Transaksi{
		ID:              uuid.New(),
		UserID:          uuid.MustParse(data.UserID),
		HargaOTR:        data.HargaOTR,
		Tenor:           data.Tenor,
		CicilanPerbulan: data.CicilanPerbulan,
		AdminFee:        data.AdminFee,
		Type:            data.Type,
		JumlahCicilan:   data.JumlahCicilan,
		NamaAsset:       data.NamaAsset,
		CreatedAt:       time.Now(),
	}

	// buat payment
	payment := make([]gormadapter.Payment, 0)
	for i := 1; i <= data.Tenor; i++ {
		t := gormadapter.Payment{
			ID:             uuid.New(),
			IDTransaksi:    transaksiData.ID,
			IDUser:         transaksiData.UserID,
			NominalNormal:  transaksiData.CicilanPerbulan,
			Denda:          0,
			NominalBayar:   0,
			Interest:       data.MontlhlyInterest(),
			IsPaid:         false,
			CreatedAt:      time.Now(),
			IdxTransaction: i,
			JatuhTempo:     time.Now().AddDate(0, i, 0),
		}
		payment = append(payment, t)
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&transaksiData).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.CreateInBatches(&payment, len(payment)).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit().Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return transaksiData.ConvertToTransactionDomain(), nil
	}

	return transaksiData.ConvertToTransactionDomain(), nil
}

// GetLimitByUser implements customer.Repository.
func (r *customerRepository) GetLimitByUser(userid string) (result domain.LoanLimit, err error) {
	var customerlimit gormadapter.CustomerWithLoanLimit
	err = r.db.Table("customers").
		Select("customers.*, loan_limits.*").
		Joins("left join loan_limits on loan_limits.userid = customers.id").
		Where("customers.id = ?", userid).
		First(&customerlimit).Error
	if err != nil {
		return result, err
	}

	result = domain.LoanLimit{
		Email:      customerlimit.Customer.Email,
		Userid:     customerlimit.Customer.ID.String(),
		OneMonth:   customerlimit.LoanLimit.OneMonth,
		TwoMonth:   customerlimit.LoanLimit.TwoMonth,
		ThreeMonth: customerlimit.LoanLimit.ThreeMonth,
		FourMonth:  customerlimit.LoanLimit.FourMonth,
	}
	return result, nil
}

func (r *customerRepository) Create(data domain.CustomerEntity) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		insertData := gormadapter.Customer{
			ID:         uuid.New(),
			NIK:        data.NIK,
			Email:      data.Email,
			FullName:   data.FullName,
			Salt:       data.Salt,
			Password:   data.Password,
			LegalName:  data.LegalName,
			BirthPlace: data.BirthPlace,
			BirthDate:  data.BirthDate,
			Salary:     data.Salary,
			PhotoID:    data.PhotoID,
			Selfie:     data.Selfie,
		}

		err := tx.Create(&insertData).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// insert to loan
		// loan := insertData.CreateLoanLimits()
		// err = tx.Create(&loan).Error
		// if err != nil {
		// 	tx.Rollback()
		// 	return err
		// }

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *customerRepository) GetUser(uniquevalue string) (domain.CustomerEntity, error) {
	var user gormadapter.Customer
	err := r.db.Where("id = ? or email = ? or nik = ?", uniquevalue, uniquevalue, uniquevalue).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.CustomerEntity{}, nil
		}
		return domain.CustomerEntity{}, err
	}
	return user.ConvertToCustomerDomain(), nil
}

// GetUserBySomething implements customer.Repository.
func (r *customerRepository) GetUserByEmail(value string) (domain.CustomerEntity, error) {
	var user gormadapter.Customer
	err := r.db.Where("`email` = ?", value).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.CustomerEntity{}, nil
		}
		return domain.CustomerEntity{}, err
	}
	return user.ConvertToCustomerDomain(), nil
}

func NewCustomerRepository(db *gorm.DB) domain.Repository {
	return &customerRepository{db: db}
}
