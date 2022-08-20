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
	var solution string 
	var err error
	switch true {
	case in.Config.CaptchaSettings.Self != "":
		solution, err = in.self(sitekey, rqData)
	case in.Config.CaptchaSettings.CaptchaAPI == "invisifox.com":
		solution, err = in.invisifox(sitekey, cookie, rqData)
	case in.Config.CaptchaSettings.CaptchaAPI == "captchaai.io":
		solution, err = in.captchaAI(sitekey, rqData)
	case utilities.Contains([]string{"capmonster.cloud", "anti-captcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		solution, err = in.Capmonster(sitekey, url, rqData, cookie)
	case utilities.Contains([]string{"2captcha.com", "rucaptcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		solution, err = in.twoCaptcha(sitekey, rqData, url)
	case in.Config.CaptchaSettings.CaptchaAPI == "capcat.xyz":
		solution, err = in.CapCat(sitekey, rqData)
	default:
		return "", fmt.Errorf("unsupported captcha api: %s", in.Config.CaptchaSettings.CaptchaAPI)
	}
	if err != nil {
		return "", err
	}
	utilities.CaptchaSolved(in.CensorToken(), solution)
	return solution, nil
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
	q.Set("pageurl", "https://discord.com")
	q.Set("userAgent", in.UserAgent)
	q.Set("json", "1")
	q.Set("soft_id", "3359")
	if rqdata != "" {
		q.Set("data", rqdata)
		q.Set("invisible", "0")
	}
	if in.Config.ProxySettings.ProxyForCaptcha {
		q.Set("proxy", in.Proxy)
		q.Set("proxytype", "http")
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
	in.LastIDstr = inResponse.Request
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
			submitCaptcha.Task.ProxyType = "http"
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
	submitCaptcha.Task.WebsiteURL, submitCaptcha.Task.WebsiteKey, submitCaptcha.Task.UserAgent = "https://discord.com", sitekey, in.UserAgent
	if rqdata != "" && in.Config.CaptchaSettings.CaptchaAPI == "capmonster.cloud" {
		submitCaptcha.Task.Data = rqdata
		// Try with true too
		submitCaptcha.Task.IsInvisible = true
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
	site := "https://api.anti-captcha.com/reportIncorrectHcaptcha"
	payload := CapmonsterPayload{
		ClientKey: in.Config.CaptchaSettings.ClientKey,
		TaskId:    in.LastID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error while marshalling payload %v", err)
	}
	resp, err := http.Post(site, "application/json", bytes.NewBuffer(payloadBytes))
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

func (in *Instance) CapCat(sitekey, rqdata string) (string, error) {
	postURL := "http://capcat.xyz/api/tasks"
	x := CapCat{
		SiteKey: sitekey,
		RqData:  rqdata,
		ApiKey:  in.Config.CaptchaSettings.ClientKey,
	}
	ipAPI := "https://api.myip.com"
	req, err := http.NewRequest("GET", ipAPI, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request [%v]", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response [%v]", err)
	}
	if !strings.Contains(string(body), "ip") {
		return "", fmt.Errorf("error invalid response [%v]", string(body))
	}
	var ipResponse map[string]interface{}
	err = json.Unmarshal(body, &ipResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response [%v]", err)
	}
	x.IP = ipResponse["ip"].(string)
	payload, err := json.Marshal(x)
	if err != nil {
		return "", fmt.Errorf("error while marshalling payload %v", err)
	}
	req, err = http.NewRequest(http.MethodPost, postURL, strings.NewReader(string(payload)))
	if err != nil {
		return "", fmt.Errorf("error creating request [%v]", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response [%v]", err)
	}
	var outResponse CapCatResponse
	err = json.Unmarshal(body, &outResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response [%v]", err)
	}
	if outResponse.ID == 0 {
		return "", fmt.Errorf("error %v %v", outResponse.Msg, string(body))
	}
	t := time.Now()
	for {
		time.Sleep(5 * time.Second)
		if time.Since(t).Seconds() >= float64(in.Config.CaptchaSettings.Timeout) || time.Since(t).Seconds() >= 300 {
			return "", fmt.Errorf("timedout - increase timeout in config to wait longer")
		}
		getURL := "http://capcat.xyz/api/result/"
		y := CapCat{
			ID:     fmt.Sprintf("%v", outResponse.ID),
			ApiKey: in.Config.CaptchaSettings.ClientKey,
		}
		payload, err = json.Marshal(y)
		if err != nil {
			return "", fmt.Errorf("error while marshalling payload %v", err)
		}
		req, err = http.NewRequest(http.MethodPost, getURL, strings.NewReader(string(payload)))
		if err != nil {
			return "", fmt.Errorf("error creating request [%v]", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("error sending request [%v]", err)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading response [%v]", err)
		}
		var outResponse CapCatResponse
		err = json.Unmarshal(body, &outResponse)
		if err != nil {
			return "", fmt.Errorf("error unmarshalling response [%v]", err)
		}
		if strings.Contains(string(body), "working") {
			continue
		} else if outResponse.Code == 1 && outResponse.Data != "" {
			return outResponse.Data, nil
		} else {
			return "", fmt.Errorf("error %v", string(body))
		}
	}
}

func (in *Instance) self(sitekey, rqData string) (string, error) {
	var solution string
	var err error
	link := in.Config.CaptchaSettings.Self
	if link == "" {
		return "", fmt.Errorf("self captcha not configured")
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	selfPayload := SelfRequest{
		Sitekey:   sitekey,
		RqData:    rqData,
		Host:      "discord.com",
		Proxy:     in.Proxy,
		Username:  in.Config.CaptchaSettings.SelfUsername,
		Password:  in.Config.CaptchaSettings.SelfPassword,
		ProxyType: "http",
	}
	payloadBytes, err := json.Marshal(selfPayload)
	if err != nil {
		return "", fmt.Errorf("error marshalling payload %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, link, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request [%v]", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response [%v]", err)
	}
	var outResponse SelfResponse
	err = json.Unmarshal(body, &outResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response [%v]", err)
	}
	if outResponse.Answer != "" {
		solution = outResponse.Answer
	} else {
		return "", fmt.Errorf("error %v", string(body))
	}
	return solution, err
}

type CapCat struct {
	ApiKey  string `json:"apikey"`
	SiteKey string `json:"sitkey"`
	RqData  string `json:"rqdata"`
	IP      string `json:"ip"`
	ID      string `json:"id,omitempty"`
}

type CapCatResponse struct {
	ID   int    `json:"id,omitempty"`
	Msg  string `json:"mess,omitempty"`
	Code int    `json:"code,omitempty"`
	Data string `json:"data,omitempty"`
}

type SelfRequest struct {
	Sitekey   string `json:"sitekey"`
	RqData    string `json:"rqdata"`
	Proxy     string `json:"proxy"`
	Host      string `json:"host"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	ProxyType string `json:"proxytype"`
}

type SelfResponse struct {
	Answer string `json:"generated_pass_UUID"`
}

/*
	invisifox
*/

func (in *Instance) invisifox(sitekey, cookie, rqdata string) (string, error) {
	site := "http://localhost:8888/solve"
	inURL, err := url.Parse(site)
	if err != nil {
		return "", fmt.Errorf("error parsing url [%v]", err)
	}
	q := inURL.Query()
	q.Set("sitekey", sitekey)
	q.Set("host", "discord.com")
	if rqdata != "" {
		q.Set("rqdata", rqdata)
	}
	if in.Config.ProxySettings.ProxyForCaptcha {
		if strings.Contains(in.Proxy, "http://") {
			q.Set("proxy", strings.Split(in.Proxy, "http://")[1])
		} else {
			q.Set("proxy", in.Proxy)
		}
	}
	inURL.RawQuery = q.Encode()
	finalSite := inURL.String()
	req, err := http.NewRequest(http.MethodPost, finalSite, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request [%v]", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response [%v]", err)
	}
	switch resp.StatusCode {
	case 200:
		return string(body), nil
	case 500:
		return "", fmt.Errorf("an unknown error has occured")
	case 404:
		return "", fmt.Errorf("trial attempts exceeded")
	case 405:
		return "", fmt.Errorf("got flagged images, your proxies or useragent might be overused")
	case 401:
		return "", fmt.Errorf("error while communicating with hcaptcha. it may be an issue with the proxy timing out or being rate limited")
	case 408:
		return "", fmt.Errorf("solver instance timedout. could not load/solve captcha in time")
	default:
		return "", fmt.Errorf("error %v %s", resp.StatusCode, string(body))
	}

}

type invisifoxRequest struct {
	Sitekey string `json:"sitekey"`
	Proxy   string `json:"proxy,omitempty"`
	Host    string `json:"host"`
	Rqdata  string `json:"rqdata,omitempty"`
}

/*
	captchaai
*/

func (in *Instance) captchaAI(sitekey, rqdata string) (string, error) {
	var captchaSolution string
	var err error
	submitURL := "https://api.captchaai.io/createTask"
	var submitPayload captchaAIpayload
	submitPayload.ClientKey = in.Config.CaptchaSettings.ClientKey
	if in.Config.ProxySettings.ProxyForCaptcha {
		submitPayload.Task.Type = "HCaptchaTask"
		submitPayload.Task.ProxyType = "http"
		var onlyProxy string 
		if strings.Contains(in.Proxy, "http://") {
			onlyProxy = strings.Split(in.Proxy, "http://")[1]
		} else {
			onlyProxy = in.Proxy
		}
		if strings.Contains(onlyProxy, "@") {
			auth, proxy := strings.Split(onlyProxy, "@")[0], strings.Split(onlyProxy, "@")[1]
			if !strings.Contains(auth, ":") {
				return captchaSolution, fmt.Errorf("proxy auth is not in format user:pass")
			}
			submitPayload.Task.ProxyLogin = strings.Split(auth, ":")[0]
			submitPayload.Task.ProxyPassword = strings.Split(auth, ":")[1]
			if !strings.Contains(proxy, ":") {
				return captchaSolution, fmt.Errorf("proxy is not in format host:port")
			}
			submitPayload.Task.ProxyAddress = strings.Split(proxy, ":")[0]
			p, err := strconv.Atoi(strings.Split(proxy, ":")[1])
			if err != nil {
				return captchaSolution, fmt.Errorf("proxy port is not a number")
			}
			submitPayload.Task.ProxyPort = p
		} else {
			if !strings.Contains(onlyProxy, ":") {
				return captchaSolution, fmt.Errorf("proxy is not in format host:port")
			}
			submitPayload.Task.ProxyAddress = strings.Split(onlyProxy, ":")[0]
			p, err := strconv.Atoi(strings.Split(onlyProxy, ":")[1])
			if err != nil {
				return captchaSolution, fmt.Errorf("proxy port is not a number")
			}
			submitPayload.Task.ProxyPort = p
		}
	} else {
		submitPayload.Task.Type = "HCaptchaTaskProxyless"
	}
	if rqdata == "" {
		submitPayload.Task.IsEnterprise = false
	} else {
		submitPayload.Task.IsEnterprise = true
		submitPayload.Task.EnterprisePayload.RqData = rqdata
	}
	submitPayload.Task.WebsiteURL = "https://discord.com"
	submitPayload.Task.WebsiteKey = sitekey
	submitPayload.Task.UserAgent = in.UserAgent
	submitPayload.Task.IsInvisible = false
	bytes, err := json.Marshal(submitPayload)
	if err != nil {
		return captchaSolution, fmt.Errorf("error while marshalling payload [%v]", err)
	}
	req, err := http.NewRequest(http.MethodPost, submitURL, strings.NewReader(string(bytes)))
	if err != nil {
		return captchaSolution, fmt.Errorf("error creating request [%v]", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return captchaSolution, fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return captchaSolution, fmt.Errorf("error reading response [%v]", err)
	}
	var submitResponse CaptchaAISubmitResponse
	err = json.Unmarshal(body, &submitResponse)
	if err != nil {
		return captchaSolution, fmt.Errorf("error unmarshalling response [%v]", err)
	}
	if submitResponse.ErrorID != 0 {
		return captchaSolution, fmt.Errorf("CaptchaAI error while submitting payload, errorID %d errroCode %s errorDescription %s", submitResponse.ErrorID, submitResponse.ErrorCode, submitResponse.ErrorDescription)
	}
	if submitResponse.TaskID == "" {
		return captchaSolution, fmt.Errorf("CaptchaAI error while submitting payload, no taskID")
	}
	now := time.Now()
	for {
		if time.Since(now) > time.Duration(in.Config.CaptchaSettings.Timeout)*time.Second {
			return captchaSolution, fmt.Errorf("captcha timeout")
		}
		time.Sleep(15 * time.Second)
		result := "https://api.captchaai.io/getTaskResult"
		var resultPayload captchaAIpayload
		resultPayload.ClientKey = in.Config.CaptchaSettings.ClientKey
		resultPayload.TaskID = submitResponse.TaskID
		bytes, err := json.Marshal(resultPayload)
		if err != nil {
			return captchaSolution, fmt.Errorf("error while marshalling payload [%v]", err)
		}
		req, err := http.NewRequest(http.MethodPost, result, strings.NewReader(string(bytes)))
		if err != nil {
			return captchaSolution, fmt.Errorf("error creating request [%v]", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return captchaSolution, fmt.Errorf("error sending request [%v]", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return captchaSolution, fmt.Errorf("error reading response [%v]", err)
		}
		var resultResponse CaptchaAISubmitResponse
		err = json.Unmarshal(body, &resultResponse)
		if err != nil {
			return captchaSolution, fmt.Errorf("error unmarshalling response [%v]", err)
		}
		if resultResponse.ErrorID != 0 {
			return captchaSolution, fmt.Errorf("CaptchaAI error while submitting payload, errorID %d errroCode %s errorDescription %s", resultResponse.ErrorID, resultResponse.ErrorCode, resultResponse.ErrorDescription)
		}
		if resultResponse.Status == "processing" {
			continue
		}
		return resultResponse.Solution.CaptchaResponse, nil
	}

}

type captchaAIpayload struct {
	ClientKey string        `json:"clientKey"`
	Task      CaptchaAiTask `json:"task,omitempty"`
	TaskID    string        `json:"taskId,omitempty"`
}

type CaptchaAiTask struct {
	Type              string            `json:"type"`
	WebsiteURL        string            `json:"websiteURL"`
	WebsiteKey        string            `json:"websiteKey"`
	ProxyType         string            `json:"proxyType,omitempty"`
	ProxyAddress      string            `json:"proxyAddress,omitempty"`
	ProxyPort         int            `json:"proxyPort,omitempty"`
	ProxyLogin        string            `json:"proxyLogin,omitempty"`
	ProxyPassword     string            `json:"proxyPassword,omitempty"`
	UserAgent         string            `json:"userAgent"`
	IsInvisible       bool              `json:"isInvisible"`
	IsEnterprise      bool              `json:"isEnterprise"`
	EnterprisePayload EnterprisePayload `json:"enterprisePayload,omitempty"`
}

type EnterprisePayload struct {
	RqData string `json:"rqdata,omitempty"`
}

type CaptchaAISubmitResponse struct {
	ErrorID          int      `json:"errorId"`
	TaskID           string   `json:"taskId"`
	Status           string   `json:"status"`
	CreateTime       int   `json:"createTime"`
	ErrorCode        string   `json:"errorCode"`
	ErrorDescription string   `json:"errorDescription"`
	Solution         Solution `json:"solution"`
}
