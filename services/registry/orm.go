package registry

import (
	gorm_logrus "github.com/onrik/gorm-logrus"
	"github.com/sarulabs/di"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"otp/config"
)

func OrmService() *di.Def {
	return &di.Def{
		Name: OrmServiceDefinition,
		Build: func(ctn di.Container) (interface{}, error) {
			dsn := config.Database.Uri

			log.Debug("connection to database ...")

			newLogger := gorm_logrus.New()

			gormConfig := &gorm.Config{}

			if config.Server.Debug {
				gormConfig.Logger = newLogger
			}

			db, err := gorm.Open(mysql.Open(dsn), gormConfig)

			if err != nil {
				log.Infof("could not connect to database: %s", err.Error())
				return nil, err
			}

			log.Info("successfully connected to database")

			return db, nil
		},
	}
}
