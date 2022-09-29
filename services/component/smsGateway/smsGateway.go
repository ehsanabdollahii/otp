package smsGateway

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type SmsGateway struct {
	Token string
}

func (s *SmsGateway) SendOTP(phone, template string, tokens ...string) error {
	url := ``

	url = fmt.Sprintf(url, s.Token, phone, tokens[0], template)

	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}
	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	apiRes := &APIResponse{}

	err = json.Unmarshal(body, apiRes)
	if err != nil {
		return err
	}

	if apiRes.Return.Status == 200 {
		log.Infof("sent otp to %s", phone)
	} else {
		errorText := fmt.Sprintf("failed sending otp to %s", phone) + apiRes.Return.Message
		log.Error(errorText)
		return errors.New(errorText)
	}

	return nil
}

type APIResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
}
