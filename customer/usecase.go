package customer

type Usecase interface {
	Create(data CustomerEntity) error
}
