package repository_test

import (
	"errors"
	"testing"

	mysqlrepo "github.com/higansama/xyz-multi-finance/customer/persistance/mysql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Create mock for the customer repository and use the GetFirstUnpaidPayment method
func TestGetFirstUnpaidPayment_Success(t *testing.T) {
	// Setup mock DB and SQLMock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	// Use the GORM postgres driver with the mock DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	// Create a mock response for the SQL query
	rows := sqlmock.NewRows([]string{"id", "id_transaksi", "id_user", "nominal_normal", "denda", "nominal_bayar", "sisa_cicilan", "idx_transaction"}).
		AddRow(uuid.New(), uuid.New(), uuid.New(), 1000000, 0, 0, 5000000, 1)

	// Mock the expected query
	mock.ExpectQuery("JOIN transaksis ON transaksis.id = payments.id_transaksi").
		WithArgs("12345", false).
		WillReturnRows(rows)

	// Create a repository and call the GetFirstUnpaidPayment method
	repo := mysqlrepo.NewCustomerRepository(gormDB)
	payment, err := repo.GetFirstUnpaidPayment("12345")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1000000, payment.NominalNormal)
	assert.Equal(t, 1, payment.IdxPayment)

	// Check that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetFirstUnpaidPayment_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	// Mengantisipasi query inisialisasi GORM
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.0.25"))

	// Menginisialisasi GORM dengan mock DB
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	// Ekspektasi query
	mock.ExpectQuery("JOIN transaksis ON transaksis.id = payments.id_transaksi").
		WithArgs("12345", false).
		WillReturnError(gorm.ErrRecordNotFound)

	repo := mysqlrepo.NewCustomerRepository(gormDB) // Added closing parenthesis

	_, err = repo.GetFirstUnpaidPayment("12345")
	if !errors.Is(err, gorm.ErrRecordNotFound) { // Use errors.Is
		t.Errorf("expected error 'record not found', got %v", err)
	}

	// Pastikan semua ekspektasi terpenuhi
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
