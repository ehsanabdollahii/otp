package services

import (
	"github.com/sarulabs/di"
	"gorm.io/gorm"
	"otp/services/component/smsGateway"
	"otp/services/registry"
)

var container di.Container

func SetupServices(services ...*di.Def) {
	builder, _ := di.NewBuilder()

	for _, service := range services {
		err := builder.Add(*service)
		if err != nil {
			panic("error on add definition to the container")
		}
	}

	container = builder.Build()
}

func GetOrmService() *gorm.DB {
	return container.Get(registry.OrmServiceDefinition).(*gorm.DB)
}

func GetsmsGatewayService() *smsGateway.SmsGateway {
	return container.Get(registry.SmsGatewayServiceDefinition).(*smsGateway.SmsGateway)
}
