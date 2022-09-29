package auth

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"otp/dto/api"
	authdto "otp/dto/auth"
	"otp/entity"
	"otp/services"
	"time"
)

func SendOTP(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Debugf("SendOTP, failed reading request body: %s", err.Error())
		c.JSON(500, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode500InternalServerError,
				Message: api.ErrorMessage500InternalServerError,
			}},
		})
		return
	}

	request := &authdto.SendOTPRequest{}

	err = json.Unmarshal(body, request)
	if err != nil {
		log.Debugf("SendOTP, failed parsing json: %s", err.Error())
		c.JSON(500, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode400InvalidJson,
				Message: api.ErrorMessage400InvalidJson,
			}},
		})
		return
	}

	validate := validator.New()

	err = validate.Struct(request)
	if err != nil {
		log.Debugf("SendOTP, failed validating data: %s", err.Error())
		c.JSON(500, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode400InvalidData,
				Message: api.ErrorMessage400InvalidData,
			}},
		})
		return
	}

	db := services.GetOrmService()

	deleteAllLoginOTPsForSource(db, request.Phone)

	otp := &entity.OneTimePassword{
		Type:      entity.OTPTypeLogin,
		Source:    request.Phone,
		ExpiresAt: time.Now().Add(time.Minute * 5),
		CreatedAt: time.Now(),
	}

	otp.SetCode(5, true)

	result := db.Create(otp)

	if result.Error != nil {
		log.Error("Failed inserting OTP into database", result.Error)
		c.JSON(500, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode500InternalServerError,
				Message: api.ErrorMessage500InternalServerError,
			}},
		})
		return
	}

	sg := services.GetsmsGatewayService()
	err = sg.SendOTP(request.Phone, "verify", otp.Code)

	if err != nil {
		log.Debugf("SendOTP, failed sending sms smsGateway: %s", err.Error())
		c.JSON(500, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode500InternalServerError,
				Message: api.ErrorMessage500InternalServerError,
			}},
		})
		return
	}

	c.JSON(201, nil)
	return
}

func VerifyOTP(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Debugf("VerifyOTP, failed reading request body: %s", err.Error())
		c.JSON(500, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode500InternalServerError,
				Message: api.ErrorMessage500InternalServerError,
			}},
		})
		return
	}

	request := &authdto.VerifyOTPRequest{}

	err = json.Unmarshal(body, request)
	if err != nil {
		log.Debugf("VerifyOTP, failed parsing json: %s", err.Error())
		c.JSON(500, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode400InvalidJson,
				Message: api.ErrorMessage400InvalidJson,
			}},
		})
		return
	}

	validate := validator.New()

	err = validate.Struct(request)
	if err != nil {
		log.Debugf("VerifyOTP, failed validating data: %s", err.Error())
		c.JSON(500, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode400InvalidData,
				Message: api.ErrorMessage400InvalidData,
			}},
		})
		return
	}

	db := services.GetOrmService()

	result := db.Where(&entity.OneTimePassword{Code: request.Code, Type: entity.OTPTypeLogin}).First(&entity.OneTimePassword{})

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(400, api.Response{
			OK:       false,
			Response: nil,
			Errors: []api.Error{{
				Code:    api.ErrorCode400InvalidData,
				Message: api.ErrorMessage400InvalidData,
			}},
		})
		return
	}

	user := &entity.User{Phone: request.Phone}

	usersResult := db.Where(&entity.User{Phone: request.Phone}).First(user)

	if usersResult.Error != nil && errors.Is(usersResult.Error, gorm.ErrRecordNotFound) {
		// user is not found , so we create one with this phone number
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.IsActive = true
		result = db.Create(&user)
	}

	LoggedInDevice := &entity.LoggedInDevice{User: *user, ExpiresAt: time.Now().AddDate(0, 1, 0)}
	LoggedInDevice.SetToken()

	result = db.Create(&LoggedInDevice)

	deleteAllLoginOTPsForSource(db, request.Phone)

	c.JSON(200, api.Response{
		OK:       true,
		Response: map[string]string{"token": LoggedInDevice.Token, "expires": LoggedInDevice.ExpiresAt.Format(time.RFC3339)},
		Errors:   []api.Error{},
	})
	return
}

func deleteAllLoginOTPsForSource(db *gorm.DB, source string) {
	db.Unscoped().Where(&entity.OneTimePassword{Source: source, Type: entity.OTPTypeLogin}).Delete(&entity.OneTimePassword{})
}
