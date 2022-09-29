package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"otp/dto/api"
	"otp/entity"
	"otp/services"
)

func GetUserID(c *gin.Context) int {
	userID, exists := c.Get("userID")

	if exists {
		userIDNumber := userID.(int)
		return userIDNumber
	}

	return 0
}

func LoginRequired(c *gin.Context) {
	if GetUserID(c) == 0 {
		c.JSON(403, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode403Forbidden,
				Message: api.ErrorMessage403Forbidden,
			}},
		})
		return
	}
}

func Authenticate(c *gin.Context) {
	token := ""

	log.Debug("trying to get the token cookie")
	token, _ = c.Cookie("Token")

	log.Debug("trying to get the token header")
	token = c.GetHeader("Token")

	userID := ResolveUserIDByToken(token)

	c.Set("userID", userID)

}

func ResolveUserIDByToken(token string) int {
	if token == "" {
		return 0
	}

	db := services.GetOrmService()

	var loggedInDeviceEntity entity.LoggedInDevice

	res := db.Where(&entity.LoggedInDevice{Token: token}).First(&loggedInDeviceEntity)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0
	}

	if loggedInDeviceEntity.IsExpired() {
		return 0
	}

	return loggedInDeviceEntity.UserID
}
