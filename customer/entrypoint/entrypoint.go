package customerpackage

import (
	"github.com/gin-gonic/gin"

	customerhandler "github.com/higansama/xyz-multi-finance/customer/delivery/http"
	mysqlpersistance "github.com/higansama/xyz-multi-finance/customer/persistance/mysql"
	usecasev1 "github.com/higansama/xyz-multi-finance/customer/usecase/v1"
	"github.com/higansama/xyz-multi-finance/internal/infrastructure"
)

func RegisterModuleCustomer(infra *infrastructure.Infrastructure, engine *gin.Engine) error {
	customerRepo := mysqlpersistance.NewCustomerRepository(infra.MysqlGorm)
	usecasev1instance, err := usecasev1.NewCustomerUsecase(infra, customerRepo)
	if err != nil {
		return err
	}

	// err = delivery.NewHandler(engine, *infra, u)
	// if err != nil {
	// 	return err
	// }

	customerhandler.NewCustomerHandler(engine, *infra, usecasev1instance)
	return nil
}
