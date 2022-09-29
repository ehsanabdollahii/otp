package middleware

import (
	"github.com/gin-gonic/gin"
	"otp/model/auth"
)

func ApplyClientRouter(router *gin.Engine) {
	versionOne := router.Group("/api/v1")
	{
		authGroup := versionOne.Group("/auth")
		{
			authGroup.POST("/send_otp", auth.SendOTP)
			authGroup.POST("/verify_otp", auth.VerifyOTP)
		}

	}
}
