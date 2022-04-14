package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (in *Instance) SolveCaptcha(sitekey string, cookie string, rqData string, rqToken string, url string) (string, error) {
	switch true {
	case Contains([]string{"capmonster.cloud", "anti-captcha.com", "anycaptcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		return in.Capmonster(sitekey, url, rqData, cookie)
	case Contains([]string{"2captcha.com", "rucaptcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		return in.twoCaptcha(sitekey, rqData, url)
	default:
		return "", fmt.Errorf("unsupported captcha api: %s", in.Config.CaptchaSettings.CaptchaAPI)
	}
}

/*
	2Captcha/RuCaptcha
*/

type twoCaptchaSubmitResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

func (in *Instance) twoCaptcha(sitekey, rqdata, site string) (string, error) {
	var solvedKey string
	inEndpoint := "https://2captcha.com/in.php"
	inURL, err := url.Parse(inEndpoint)
	if err != nil {
		return solvedKey, fmt.Errorf("error while parsing url %v", err)
	}
	q := inURL.Query()
	if in.Config.CaptchaSettings.ClientKey == "" {
		return solvedKey, fmt.Errorf("client key is empty")
	}
	q.Set("key", in.Config.CaptchaSettings.ClientKey)
	q.Set("method", "hcaptcha")
	q.Set("sitekey", sitekey)
	// Page URL same as referer in headers
	q.Set("pageurl", site)
	q.Set("userAgent", UserAgent)
	q.Set("json", "1")
	q.Set("soft_id", "3359")
	if rqdata != "" {
		q.Set("data", rqdata)
		q.Set("invisible", "1")
	}
	if in.Config.ProxySettings.ProxyForCaptcha {
		q.Set("proxy", in.Proxy)
		q.Set("proxytype", in.Config.ProxySettings.ProxyProtocol)
	}
	inURL.RawQuery = q.Encode()
	if in.Config.CaptchaSettings.CaptchaAPI == "2captcha.com" {
		inURL.Host = "2captcha.com"
	} else if in.Config.CaptchaSettings.CaptchaAPI == "rucaptcha.com" {
		inURL.Host = "rucaptcha.com"
	}
	inEndpoint = inURL.String()
	req, err := http.NewRequest(http.MethodGet, inEndpoint, nil)
	if err != nil {
		return solvedKey, fmt.Errorf("error creating request [%v]", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return solvedKey, fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return solvedKey, fmt.Errorf("error reading response [%v]", err)
	}
	var inResponse twoCaptchaSubmitResponse
	err = json.Unmarshal(body, &inResponse)
	if err != nil {
		return solvedKey, fmt.Errorf("error unmarshalling response [%v]", err)
	}
	if inResponse.Status != 1 {
		return solvedKey, fmt.Errorf("error %v", inResponse.Request)
	}
	outEndpoint := "https://2captcha.com/res.php"
	outURL, err := url.Parse(outEndpoint)
	if err != nil {
		return solvedKey, fmt.Errorf("error while parsing url %v", err)
	}
	q = outURL.Query()
	q.Set("key", in.Config.CaptchaSettings.ClientKey)
	q.Set("action", "get")
	q.Set("id", inResponse.Request)
	q.Set("json", "1")
	if in.Config.CaptchaSettings.CaptchaAPI == "2captcha.com" {
		outURL.Host = "2captcha.com"
	} else if in.Config.CaptchaSettings.CaptchaAPI == "rucaptcha.com" {
		outURL.Host = "rucaptcha.com"
	}
	outURL.RawQuery = q.Encode()
	outEndpoint = outURL.String()
	req, err = http.NewRequest(http.MethodGet, outEndpoint, nil)
	if err != nil {
		return solvedKey, fmt.Errorf("error creating request [%v]", err)
	}
	time.Sleep(10 * time.Second)
	now := time.Now()
	for {
		if time.Since(now) > time.Duration(in.Config.CaptchaSettings.Timeout)*time.Second {
			return solvedKey, fmt.Errorf("captcha response from 2captcha timedout")
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return solvedKey, fmt.Errorf("error sending request [%v]", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return solvedKey, fmt.Errorf("error reading response [%v]", err)
		}
		var outResponse twoCaptchaSubmitResponse
		err = json.Unmarshal(body, &outResponse)
		if err != nil {
			return solvedKey, fmt.Errorf("error unmarshalling response [%v]", err)
		}
		if outResponse.Request == "CAPCHA_NOT_READY" {
			time.Sleep(5 * time.Second)
			continue
		} else if strings.Contains(string(body), "ERROR") {
			return solvedKey, fmt.Errorf("error %v", outResponse.Request)
		} else {
			solvedKey = outResponse.Request
			break
		}
	}
	return solvedKey, nil
}

/*
	Capmonster
*/

type CapmonsterPayload struct {
	ClientKey string `json:"clientKey,omitempty"`
	Task      Task   `json:"task,omitempty"`
	TaskId    int    `json:"taskId,omitempty"`
}

type Task struct {
	CaptchaType   string `json:"type,omitempty"`
	WebsiteURL    string `json:"websiteURL,omitempty"`
	WebsiteKey    string `json:"websiteKey,omitempty"`
	IsInvisible   bool   `json:"isInvisible,omitempty"`
	Data          string `json:"data,omitempty"`
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	Cookies       string `json:"cookies,omitempty"`
}

type CapmonsterSubmitResponse struct {
	ErrorID int `json:"errorId,omitempty"`
	TaskID  int `json:"taskId,omitempty"`
}

type CapmonsterOutResponse struct {
	ErrorID   int      `json:"errorId,omitempty"`
	ErrorCode string   `json:"errorCode,omitempty"`
	Status    string   `json:"status,omitempty"`
	Solution  Solution `json:"solution"`
}

type Solution struct {
	CaptchaResponse string `json:"gRecaptchaResponse,omitempty"`
}

func (in *Instance) Capmonster(sitekey, website, rqdata, cookies string) (string, error) {
	var solvedKey string
	inEndpoint, outEndpoint := fmt.Sprintf("https://api.%s/createTask", in.Config.CaptchaSettings.CaptchaAPI), fmt.Sprintf("https://api.%s/getTaskResult", in.Config.CaptchaSettings.CaptchaAPI)
	var submitCaptcha CapmonsterPayload
	if in.Config.CaptchaSettings.ClientKey == "" {
		return solvedKey, fmt.Errorf("no client key provided in config")
	} else {
		submitCaptcha.ClientKey = in.Config.CaptchaSettings.ClientKey
	}
	if in.Config.ProxySettings.ProxyForCaptcha && in.Proxy != "" {
		submitCaptcha.Task.CaptchaType = "HCaptchaTask"
		proxyURL, err := url.Parse(in.Proxy)
		if err != nil {
			return solvedKey, fmt.Errorf("error while parsing proxy url %v", err)
		}
		submitCaptcha.Task.ProxyType = in.Config.ProxySettings.ProxyProtocol
		submitCaptcha.Task.ProxyAddress = proxyURL.Hostname()
		submitCaptcha.Task.ProxyPort, err = strconv.Atoi(proxyURL.Port())
		if err != nil {
			return solvedKey, fmt.Errorf("error while parsing proxy port %v", err)
		}
		submitCaptcha.Task.ProxyLogin = proxyURL.User.Username()
		pwd, setPwd := proxyURL.User.Password()
		if setPwd {
			submitCaptcha.Task.ProxyPassword = pwd
		}
	} else {
		submitCaptcha.Task.CaptchaType = "HCaptchaTaskProxyless"
	}
	submitCaptcha.Task.WebsiteURL, submitCaptcha.Task.WebsiteKey, submitCaptcha.Task.Data, submitCaptcha.Task.Cookies, submitCaptcha.Task.IsInvisible, submitCaptcha.Task.UserAgent = website, sitekey, rqdata, cookies, true, UserAgent
	payload, err := json.Marshal(submitCaptcha)
	if err != nil {
		return solvedKey, fmt.Errorf("error while marshalling payload %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, inEndpoint, strings.NewReader(string(payload)))
	if err != nil {
		return solvedKey, fmt.Errorf("error creating request [%v]", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return solvedKey, fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return solvedKey, fmt.Errorf("error reading response [%v]", err)
	}
	var inResponse CapmonsterSubmitResponse
	err = json.Unmarshal(body, &inResponse)
	if err != nil {
		return solvedKey, fmt.Errorf("error unmarshalling response [%v]", err)
	}
	if inResponse.ErrorID != 0 {
		return solvedKey, fmt.Errorf("error %v %v", inResponse.ErrorID, string(body))
	}
	var retrieveCaptcha CapmonsterPayload
	retrieveCaptcha.ClientKey = in.Config.CaptchaSettings.ClientKey
	retrieveCaptcha.TaskId = inResponse.TaskID
	payload, err = json.Marshal(retrieveCaptcha)
	if err != nil {
		return solvedKey, fmt.Errorf("error while marshalling payload %v", err)
	}
	req, err = http.NewRequest(http.MethodPost, outEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return solvedKey, fmt.Errorf("error creating request [%v]", err)
	}
	req.Header.Set("Content-Type", "application/json")
	t := time.Now()
	for i := 0; i < 120; i++ {
		if time.Since(t).Seconds() >= float64(in.Config.CaptchaSettings.Timeout) {
			return solvedKey, fmt.Errorf("timedout - increase timeout in config to wait longer")
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return solvedKey, fmt.Errorf("error sending request [%v]", err)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return solvedKey, fmt.Errorf("error reading response [%v]", err)
		}
		var outResponse CapmonsterOutResponse
		err = json.Unmarshal(body, &outResponse)
		if err != nil {
			return solvedKey, fmt.Errorf("error unmarshalling response [%v]", err)
		}
		if outResponse.ErrorID != 0 {
			return solvedKey, fmt.Errorf("error %v %v", outResponse.ErrorID, string(body))
		}
		if outResponse.Status == "ready" {
			solvedKey = outResponse.Solution.CaptchaResponse
			break
		} else if outResponse.Status == "processing" {
			time.Sleep(5 * time.Second)
			continue
		} else {
			return solvedKey, fmt.Errorf("error invalid status %v %v", outResponse.ErrorID, string(body))
		}

	}
	return solvedKey, nil
}
