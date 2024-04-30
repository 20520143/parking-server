package client

import (
	"parking-server/conf"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var twilioClient *TwilioClient

type TwilioClient struct {
	client *twilio.RestClient
}

func init() {
	cfg := conf.GetConfig()

	twilioClient = &TwilioClient{}
	twilioClient.client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})
}

func (t *TwilioClient) SendOtp(to string) error {
	cfg := conf.GetConfig()

	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	_, err := t.client.VerifyV2.CreateVerification(cfg.TwilioServiceSID, params)

	return err
}

func (t *TwilioClient) CheckOtp(to string, code string) (bool, error) {
	cfg := conf.GetConfig()

	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := t.client.VerifyV2.CreateVerificationCheck(cfg.TwilioServiceSID, params)
	if err != nil {
		return false, err
	}

	if *resp.Status == "approved" {
		return true, err
	}

	return false, nil
}

func GetTwiliioClient() *TwilioClient {
	return twilioClient
}
