package migrate

import (
	"github.com/higansama/xyz-multi-finance/config"
	gormadapter "github.com/higansama/xyz-multi-finance/persistance/gorm"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMigrate(cfg config.Config) {
	err := config.InitConfig("")
	if err != nil {
		log.Panic().Err(err).Send()
	}

	// cfg := config.Cfg
	db, err := gorm.Open(
		mysql.Open(cfg.DB.MysqlUri),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}
	// Migrasi tabel
	err = db.AutoMigrate(&gormadapter.Customer{})
	if err != nil {
		panic("gagal melakukan migrasi")
	}
}
