package customer

type Repository interface {
	Create(data CustomerEntity) error
	GetUser(uniquevalue string) (CustomerEntity, error)
	GetUserByEmail(value string) (CustomerEntity, error)
	GetLimitByUser(userid string) (LoanLimit, error)
	CreateTransaksi(data Transaksi) (Transaksi, error)
	CekTransaksiByContract(userid, contractnumber string) (TransaksiAnPayments, error)
	PayTransaction(payload UpdatePayment) error
	GetFirstUnpaidPayment(transactionID string) (Payments, error)
}
