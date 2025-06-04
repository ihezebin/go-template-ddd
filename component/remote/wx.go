package remote

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/ihezebin/olympus/httpclient"
	"github.com/pkg/errors"
)

type WXClient struct {
	*resty.Client
	appID     string
	secretKey string
	host      string
}

var wxClient *WXClient

func InitWX(host string, appID string, secretKey string) error {
	wxClient = &WXClient{
		Client:    httpclient.NewClient(httpclient.WithHost(host)),
		appID:     appID,
		secretKey: secretKey,
		host:      host,
	}

	_, err := wxClient.Token()
	if err != nil {
		return errors.Wrap(err, "init wx client error")
	}

	return nil
}

func WX() *WXClient {
	return wxClient
}

type WXUserPhoneNumberResp struct {
	PhoneInfo struct {
		// 用户绑定的手机号（国外手机号会有区号）
		PhoneNumber string `json:"phoneNumber"`
		// 没有区号的手机号
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int    `json:"timestamp"`
			AppID     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}

func (c *WXClient) GetUserPhoneNumber(token string, code string) (*WXUserPhoneNumberResp, error) {
	resp, err := c.NewRequest().SetQueryParams(map[string]string{
		"access_token": token,
	}).SetBody(map[string]interface{}{
		"code": code,
	}).Post("/wxa/business/getuserphonenumber")
	if err != nil {
		return nil, errors.Wrap(err, "get wx user phone number error")
	}
	if resp.IsError() {
		return nil, errors.Errorf("get wx user phone number resp status error: %d, body: %s", resp.StatusCode(), resp.String())
	}

	phoneResp := &WXUserPhoneNumberResp{}
	if err := json.Unmarshal(resp.Body(), phoneResp); err != nil {
		return nil, errors.Wrap(err, "unmarshal wx user phone number resp error")
	}

	return phoneResp, nil
}

type WXTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *WXClient) Token() (*WXTokenResp, error) {
	resp, err := c.NewRequest().SetQueryParams(map[string]string{
		"appid":      c.appID,
		"secret":     c.secretKey,
		"grant_type": "client_credential",
	}).Get("/cgi-bin/token")
	if err != nil {
		return nil, errors.Wrap(err, "get wx token error")
	}
	if resp.IsError() {
		return nil, errors.Errorf("get wx token resp status error: %d, body: %s", resp.StatusCode(), resp.String())
	}

	tokenResp := &WXTokenResp{}
	if err := json.Unmarshal(resp.Body(), tokenResp); err != nil {
		return nil, errors.Wrap(err, "unmarshal wx token resp error")
	}

	return tokenResp, nil
}

func (c *WXClient) GetWXACodeUnlimit(token string, scene string) (string, error) {
	resp, err := c.NewRequest().SetQueryParams(map[string]string{
		"access_token": token,
	}).SetBody(map[string]interface{}{
		"scene":      scene,
		"check_path": true,
		"page":       "pages/login/login",
	}).Post("/wxa/getwxacodeunlimit")
	if err != nil {
		return "", errors.Wrap(err, "get wx acode unlimit error")
	}
	if resp.IsError() {
		return "", errors.Errorf("get wx acode unlimit resp status error: %d, body: %s", resp.StatusCode(), resp.String())
	}

	return resp.String(), nil
}

type WXJsCode2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

func (c *WXClient) JsCode2Session(code string) (*WXJsCode2SessionResp, error) {
	loginResp := WXJsCode2SessionResp{}

	resp, err := c.NewRequest().SetQueryParams(map[string]string{
		"appid":      c.appID,
		"secret":     c.secretKey,
		"js_code":    code,
		"grant_type": "authorization_code",
	}).Get("/sns/jscode2session")

	if err != nil {
		return nil, errors.Wrap(err, "get wx login resp error")
	}

	if resp.IsError() {
		return nil, errors.Errorf("get wx login resp status error: %d, body: %s", resp.StatusCode(), resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &loginResp); err != nil {
		return nil, errors.Wrap(err, "unmarshal wx login resp error")
	}

	return &loginResp, nil
}
