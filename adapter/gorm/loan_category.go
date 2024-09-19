package adapter

import "gorm.io/gorm"

type Loan struct {
	gorm.Model
	TypeName string
	Max      int
	Tenor    int
}
