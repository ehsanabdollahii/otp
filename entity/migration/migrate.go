package migration

import (
	log "github.com/sirupsen/logrus"
	"otp/entity"
	"otp/services"
)

func MigrateDB() error {
	db := services.GetOrmService()

	log.Debug("migrating database ...")

	err := db.AutoMigrate(
		&entity.User{},
		&entity.OneTimePassword{},
		&entity.TextMessage{},
		&entity.LoggedInDevice{},
	)

	if err != nil {
		log.Infof("migration failed: %s", err.Error())
		return err
	}

	log.Info("migration successful")

	return nil
}
