package customer

type Repository interface {
	Create(data CustomerEntity) error
	GetUser(email string) (CustomerEntity, error)
}
