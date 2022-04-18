// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

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

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func (in *Instance) SolveCaptcha(sitekey string, cookie string, rqData string, rqToken string, url string) (string, error) {
	switch true {
	case utilities.Contains([]string{"capmonster.cloud", "anti-captcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		return in.Capmonster(sitekey, url, rqData, cookie)
	case utilities.Contains([]string{"2captcha.com", "rucaptcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		return in.twoCaptcha(sitekey, rqData, url)
	default:
		return "", fmt.Errorf("unsupported captcha api: %s", in.Config.CaptchaSettings.CaptchaAPI)
	}
}

/*
	2Captcha/RuCaptcha
*/

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
		q.Set("invisible", "0")
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

	time.Sleep(10 * time.Second)
	now := time.Now()
	for {
		if time.Since(now) > time.Duration(in.Config.CaptchaSettings.Timeout)*time.Second {
			return solvedKey, fmt.Errorf("captcha response from 2captcha timedout")
		}
		req, err = http.NewRequest(http.MethodGet, outEndpoint, nil)
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

func (in *Instance) Capmonster(sitekey, website, rqdata, cookies string) (string, error) {
	var solvedKey string
	inEndpoint, outEndpoint := fmt.Sprintf("https://api.%s/createTask", in.Config.CaptchaSettings.CaptchaAPI), fmt.Sprintf("https://api.%s/getTaskResult", in.Config.CaptchaSettings.CaptchaAPI)
	var submitCaptcha CapmonsterPayload
	if in.Config.CaptchaSettings.ClientKey == "" {
		return solvedKey, fmt.Errorf("no client key provided in config")
	} else {
		submitCaptcha.ClientKey = in.Config.CaptchaSettings.ClientKey
	}
	if in.Config.CaptchaSettings.CaptchaAPI == "anti-captcha.com" {
		submitCaptcha.SoftID = 1021
	}
	if in.Config.ProxySettings.ProxyForCaptcha && in.Proxy != "" {
		submitCaptcha.Task.CaptchaType = "HCaptchaTask"
		if strings.Contains(in.Proxy, "@") {
			// User:pass authenticated proxy
			parts := strings.Split(in.Proxy, "@")
			userPass, ipPort := parts[0], parts[1]
			if !strings.Contains(ipPort, ":") || !strings.Contains(userPass, ":") {
				return solvedKey, fmt.Errorf("invalid proxy format")
			}
			submitCaptcha.Task.ProxyType = in.Config.ProxySettings.ProxyProtocol
			submitCaptcha.Task.ProxyLogin, submitCaptcha.Task.ProxyPassword = strings.Split(userPass, ":")[0], strings.Split(userPass, ":")[1]
			port := strings.Split(ipPort, ":")[1]
			var err error
			submitCaptcha.Task.ProxyPort, err = strconv.Atoi(port)
			if err != nil {
				return solvedKey, fmt.Errorf("invalid proxy format")
			}
			submitCaptcha.Task.ProxyAddress = strings.Split(ipPort, ":")[0]
		} else {
			if !strings.Contains(in.Proxy, ":") {
				return solvedKey, fmt.Errorf("invalid proxy format")
			}
			submitCaptcha.Task.ProxyAddress = strings.Split(in.Proxy, ":")[0]
			port := strings.Split(in.Proxy, ":")[1]
			var err error
			submitCaptcha.Task.ProxyPort, err = strconv.Atoi(port)
			if err != nil {
				return solvedKey, fmt.Errorf("invalid proxy format")
			}
		}
	} else {
		submitCaptcha.Task.CaptchaType = "HCaptchaTaskProxyless"
	}
	submitCaptcha.Task.WebsiteURL, submitCaptcha.Task.WebsiteKey, submitCaptcha.Task.Cookies, submitCaptcha.Task.UserAgent = website, sitekey, cookies, UserAgent
	if rqdata != "" && in.Config.CaptchaSettings.CaptchaAPI == "capmonster.cloud" {
		submitCaptcha.Task.Data = rqdata
		// Try with true too
		submitCaptcha.Task.IsInvisible = false
	} else if rqdata != "" && in.Config.CaptchaSettings.CaptchaAPI == "anti-captcha.com" {
		submitCaptcha.Task.IsInvisible = false
		submitCaptcha.Task.Enterprise.RqData = rqdata
		submitCaptcha.Task.Enterprise.Sentry = true
	}
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
	in.LastID = inResponse.TaskID
	payload, err = json.Marshal(retrieveCaptcha)
	if err != nil {
		return solvedKey, fmt.Errorf("error while marshalling payload %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	time.Sleep(5 * time.Second)
	t := time.Now()
	for i := 0; i < 120; i++ {
		if time.Since(t).Seconds() >= float64(in.Config.CaptchaSettings.Timeout) {
			return solvedKey, fmt.Errorf("timedout - increase timeout in config to wait longer")
		}
		req, err = http.NewRequest(http.MethodPost, outEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			return solvedKey, fmt.Errorf("error creating request [%v]", err)
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

func (in *Instance) ReportIncorrectRecaptcha() error {
	site := "https://api.anti-captcha.com/reportIncorrectRecaptcha"
	payload := CapmonsterPayload{
		ClientKey: in.Config.CaptchaSettings.ClientKey,
		TaskId:    in.LastID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error while marshalling payload %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, site, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request [%v]", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response [%v]", err)
	}
	var outResponse CapmonsterOutResponse
	err = json.Unmarshal(body, &outResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling response [%v]", err)
	}
	if outResponse.Status != "success" {
		return fmt.Errorf("error %v ", outResponse.ErrorID)
	}

	return nil
}
