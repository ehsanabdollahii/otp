package registry

import (
	"github.com/sarulabs/di"
	"otp/config"
	"otp/services/component/smsGateway"
)

func SmsGatewayService() *di.Def {
	return &di.Def{
		Name: SmsGatewayServiceDefinition,
		Build: func(ctn di.Container) (interface{}, error) {
			smsGatewayService := smsGateway.SmsGateway{Token: config.Server.SmsGatewayToken}
			return &smsGatewayService, nil
		},
	}
}
