package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"otp/api/middleware"
	"otp/config"
	"otp/entity/migration"
	"otp/model/auth"
	"otp/services"
	"otp/services/registry"
)

// Server is the cli command that runs our main web server
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Starts the web server",
		Before: func(c *cli.Context) error {

			if config.Server.Debug {
				log.SetLevel(log.DebugLevel)
				gin.SetMode(gin.DebugMode)
			}

			services.SetupServices(registry.OrmService(), registry.SmsGatewayService())

			migrationErr := migration.MigrateDB()
			if migrationErr != nil {
				return migrationErr
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "mysql-uri",
				Value:       "",
				Usage:       "Mysql database uri",
				EnvVars:     []string{"MYSQL_URI"},
				Destination: &config.Database.Uri,
			},
			&cli.StringFlag{
				Name:        "asynq-redis",
				Value:       "",
				Usage:       "Asynq redis",
				EnvVars:     []string{"ASYNQ_REDIS"},
				Destination: &config.Database.AsynqRedis,
			},
			&cli.StringFlag{
				Name:        "static-path",
				Value:       "./front/build",
				Usage:       "Static files path",
				EnvVars:     []string{"STATIC_PATH"},
				Destination: &config.Server.StaticPath,
			},
			&cli.IntFlag{
				Name:        "server-port",
				Value:       8066,
				Usage:       "Web server port",
				EnvVars:     []string{"SERVER_PORT", "PORT"},
				Destination: &config.Server.ServerPort,
			},
			&cli.StringFlag{
				Name:        "server-host",
				Value:       "0.0.0.0",
				Usage:       "Web server host",
				EnvVars:     []string{"SERVER_HOST"},
				Destination: &config.Server.ServerHost,
			},
			&cli.StringFlag{
				Name:        "smsGateway-token",
				Value:       "",
				Usage:       "smsGateway_TOKEN",
				EnvVars:     []string{"smsGateway_TOKEN"},
				Destination: &config.Server.SmsGatewayToken,
			},
		},
		Action: func(c *cli.Context) error {
			r := gin.Default()

			r.Use(static.Serve("/", static.LocalFile(config.Server.StaticPath, true)))

			r.Use(gin.Logger())
			r.Use(auth.Authenticate)

			middleware.ApplyClientRouter(r)

			err := r.Run(fmt.Sprintf("%s:%d", config.Server.ServerHost, config.Server.ServerPort))
			if err != nil {
				panic("failed running gin")
			}

			return nil
		},
	}
}
