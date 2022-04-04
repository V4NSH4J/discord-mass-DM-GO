package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/fatih/color"
)

func (in *Instance) SolveCaptcha(sitekey string, cookie string, rqData string, rqToken string) (string, error) {
	switch true {
	case Contains([]string{"capmonster.cloud", "anti-captcha.com", "anycaptcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		return in.SolveCaptchaCapmonster(sitekey, cookie, rqData)
	case Contains([]string{"rucaptcha.com", "azcaptcha.com", "solvecaptcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		return in.SolveCaptchaRucaptcha(sitekey, rqData, rqToken)
	case in.Config.CaptchaSettings.CaptchaAPI == "deathbycaptcha.com":
		return in.SolveCaptchaDeathByCaptcha(sitekey)
	case in.Config.CaptchaSettings.CaptchaAPI == "2captcha.com":
		return in.twoCaptcha(sitekey, rqData)
	default:
		return "", fmt.Errorf("unsupported captcha api: %s", in.Config.CaptchaSettings.CaptchaAPI)
	}
}

// Function to use a captcha solving service and return a solved captcha key
func (in *Instance) SolveCaptchaCapmonster(sitekey string, cookies string, rqdata string) (string, error) {
	var jsonx Pload
	if !in.Config.ProxySettings.ProxyForCaptcha || in.Config.CaptchaSettings.CaptchaAPI == "anycaptcha.com" {
		jsonx = Pload{
			ClientKey: in.Config.CaptchaSettings.ClientKey,
			Task: Task{
				Type:       "HCaptchaTaskProxyless",
				WebsiteURL: "https://discord.com/channels/@me",
				WebsiteKey: sitekey,
				Cookies:    cookies,
				UserAgent:  UserAgent,
				Data:       rqdata,
			},
		}
	} else {
		var address string
		var port int
		var username string
		var password string
		var err error
		// Proxies with user-pass AUTH
		if strings.Contains(in.Proxy, "@") {
			proxyParts := strings.Split(in.Proxy, "@")
			username, password, address = strings.Split(proxyParts[0], ":")[0], strings.Split(proxyParts[0], ":")[1], strings.Split(proxyParts[1], ":")[0]
			port, err = strconv.Atoi(strings.Split(proxyParts[1], ":")[1])
			if err != nil {
				return "", fmt.Errorf("could not parse proxy port %v", err)
			}
		} else {
			// IP AUTH proxies
			address = strings.Split(in.Proxy, ":")[0]
			port, err = strconv.Atoi(strings.Split(in.Proxy, ":")[1])
			if err != nil {
				return "", fmt.Errorf("could not parse proxy port %v", err)
			}
		}
		jsonx = Pload{
			ClientKey: in.Config.CaptchaSettings.ClientKey,
			Task: Task{
				Type:          "HCaptchaTask",
				WebsiteURL:    "https://discord.com/channels/@me",
				WebsiteKey:    sitekey,
				UserAgent:     UserAgent,
				ProxyType:     in.Config.ProxySettings.ProxyProtocol,
				ProxyAddress:  address,
				ProxyPort:     port,
				ProxyLogin:    username,
				ProxyPassword: password,
				Cookies:       cookies,
				Data:          rqdata,
			},
		}
	}

	bytes, err := json.Marshal(jsonx)
	if err != nil {
		return "", fmt.Errorf("error marshalling json [%v]", err)

	}
	// Almost all solving services have similar API, so we can use the same function and replace the domain.
	resp, err := http.Post("https://api."+in.Config.CaptchaSettings.CaptchaAPI+"/createTask", "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		return "", fmt.Errorf("error creating the request for captcha [%v]", err)
	}
	defer resp.Body.Close()
	// Get taskID from response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response body [%v]", err)
	}
	var response Resp
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling the response body [%v]", err)
	}
	switch response.ErrorID {
	case 0:
		// Poling server for the solved captcha
		jsonx = Pload{
			ClientKey: in.Config.CaptchaSettings.ClientKey,
			TaskID:    response.TaskID,
		}
		y, err := json.Marshal(jsonx)
		if err != nil {
			return "", fmt.Errorf("error marshalling json [%v]", err)
		}
		// Anti Captcha documentation prescribes to use a delay of 5 seconds before requesting the captcha and 3 seconds delays after that.
		time.Sleep(5 * time.Second)
		p := 0
		for {
			if p > 100 {
				// Max retries
				break
			}
			resp, err := http.Post("https://api."+in.Config.CaptchaSettings.CaptchaAPI+"/getTaskResult", "application/json", strings.NewReader(string(y)))
			if err != nil {
				return "", fmt.Errorf("error creating the request for captcha [%v]", err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("error reading the response body [%v]", err)
			}
			var response Resp
			err = json.Unmarshal(body, &response)
			if err != nil {
				return "", fmt.Errorf("error unmarshalling the response body [%v]", err)
			}
			if response.Status == "ready" {
				return response.Solution.Ans, nil

			} else if response.Status == "processing" {
				p++ // Incrementing the counter
				time.Sleep(3 * time.Second)
				continue
			}
			if response.ErrorID != 0 {
				return "", fmt.Errorf("ErrorID: %s, ErrorCode: %s, ErrorDescription: %s", response.ErrorID, response.ErrorCode, response.ErrorDesc)
			}

		}
		return "", fmt.Errorf("max captcha retries reached [%v]", err)
	case 1:
		return "", fmt.Errorf("No captcha API key or Incorrect Captcha API key (Can happen if you've specified a different service but using a different service's key)")
	case 2:
		color.Red("No Available Captcha Workers - Increase your Maximum Bid in your Captcha API settings. Also try reducing the number of threads. Sleeping 10 seconds")
		time.Sleep(10 * time.Second)
		return "", fmt.Errorf("no captcha workers were available, retrying")
	case 3:
		return "", fmt.Errorf("the size of the captcha you are uploading is less than 100 bytes.")
	case 4:
		return "", fmt.Errorf("the size of the captcha you are uploading is greater than 500,000 bytes.")

	case 10:
		return "", fmt.Errorf("you have zero or negative captcha API balance")
	case 11:
		return "", fmt.Errorf("captcha was unsolvable.")
	default:
		return "", fmt.Errorf("unknown error [%v]", response.ErrorID)
	}
}

type Pload struct {
	ClientKey string `json:"clientKey"`
	Task      Task   `json:"task"`
	ErrorID   int    `json:"ErrorId"`
	TaskID    int    `json:"taskId"`
}

type Task struct {
	Type          string `json:"type"`
	WebsiteURL    string `json:"websiteURL"`
	WebsiteKey    string `json:"websiteKey"`
	ProxyType     string `json:"proxyType"`
	ProxyAddress  string `json:"proxyAddress"`
	ProxyPort     int    `json:"proxyPort"`
	ProxyLogin    string `json:"proxyLogin"`
	Data          string `json:"data"`
	ProxyPassword string `json:"proxyPassword"`
	UserAgent     string `json:"userAgent"`
	Cookies       string `json:"cookies`
}

type Resp struct {
	TaskID    int    `json:"taskID"`
	ErrorID   int    `json:"ErrorId"`
	Status    string `json:"status"`
	Solution  Sol    `json:"solution"`
	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDescription"`
}

type Sol struct {
	Ans string `json:"gRecaptchaResponse"`
}

func (in *Instance) SolveCaptcha2Captcha(sitekey string) (string, error) {
	client := api2captcha.NewClient(in.Config.CaptchaSettings.ClientKey)
	client.DefaultTimeout = 120
	client.PollingInterval = 22

	cap := api2captcha.HCaptcha{
		SiteKey: sitekey,
		Url:     "https://discord.com/channels/@me",
	}
	var proxyType string
	if in.Config.ProxySettings.ProxyProtocol == "socks5" {
		proxyType = "SOCKS5"
	} else if in.Config.ProxySettings.ProxyProtocol == "socks4" {
		proxyType = "SOCKS4"
	} else if in.Config.ProxySettings.ProxyProtocol == "http" {
		proxyType = "HTTPS"
	}

	req := cap.ToRequest()
	if in.Config.ProxySettings.ProxyForCaptcha {
		req.SetProxy(proxyType, in.Proxy)
	}
	req.Params["userAgent"] = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9003 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36"

	code, err := client.Solve(req)
	if err != nil {
		if err == api2captcha.ErrTimeout {
			return "", fmt.Errorf("Timeout")
		} else if err == api2captcha.ErrApi {
			return "", fmt.Errorf("API error")
		} else if err == api2captcha.ErrNetwork {
			return "", fmt.Errorf("Network error")
		} else {
			return "", fmt.Errorf("Unknown error %v", err)
		}
	}
	return code, nil
}

// Incomplete
func (in *Instance) SolveCaptchaDeathByCaptcha(sitekey string) (string, error) {
	// Authentication can be a user:pass combination or with a 2fa key.
	var username string
	var password string
	var authtoken string
	var proxy string
	var proxytype string

	if strings.Contains(in.Config.CaptchaSettings.ClientKey, ":") {
		credentials := strings.Split(in.Config.CaptchaSettings.ClientKey, ":")
		username, password = credentials[0], credentials[1]
	} else {
		authtoken = in.Config.CaptchaSettings.ClientKey
	}
	captchaPostEndpoint := "http://api.dbcapi.me/api/captcha"

	fmt.Println(authtoken)
	payload := fmt.Sprintf(
		`{
			"username": "%s",
			"password": "%s",
			"type": 7, 
			"token_params": {
				"proxy": "%s",
				"proxytype": "%s",
				"pageurl": "http://discord.com",
				"sitekey": "%s"
			}
		}`, username, password, proxy, proxytype, sitekey)
	fmt.Println(payload)
	req, err := http.NewRequest(http.MethodPost, captchaPostEndpoint, strings.NewReader(payload))
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
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)

	return "", nil

}

func (in *Instance) SolveCaptchaRucaptcha(sitekey string, rqData string, rqToken string) (string, error) {
	encUa := `Mozilla%2F5.0%20%28Windows%20NT%2010.0%3B%20Win64%3B%20x64%3B%20rv%3A83.0%29%20Gecko%2F20100101%20Firefox%2F83.0`
	var submitEndpoint string
	if !in.Config.ProxySettings.ProxyForCaptcha {
		if in.Config.CaptchaSettings.CaptchaAPI == "2captcha.com" {
			submitEndpoint = fmt.Sprintf("http://%s/in.php?key=%s&method=hcaptcha&sitekey=%s&pageurl=%s&userAgent=%s&json=1&soft_id=12368652", in.Config.CaptchaSettings.CaptchaAPI, in.Config.CaptchaSettings.ClientKey, sitekey, "https://discord.com/channels/@me", encUa)
		} else {
			submitEndpoint = fmt.Sprintf("http://%s/in.php?key=%s&method=hcaptcha&sitekey=%s&pageurl=%s&userAgent=%sjson=1&soft_id=13615286", in.Config.CaptchaSettings.CaptchaAPI, in.Config.CaptchaSettings.ClientKey, sitekey, "https://discord.com/channels/@me", encUa)
		}

	} else {
		var proxyType string
		if in.Config.ProxySettings.ProxyProtocol == "socks5" {
			proxyType = "SOCKS5"
		} else if in.Config.ProxySettings.ProxyProtocol == "socks4" {
			proxyType = "SOCKS4"
		} else if in.Config.ProxySettings.ProxyProtocol == "http" {
			proxyType = "HTTPS"
		}
		if in.Config.CaptchaSettings.CaptchaAPI == "2captcha.com" {
			submitEndpoint = fmt.Sprintf("http://%s/in.php?key=%s&method=hcaptcha&sitekey=%s&pageurl=%s&userAgent=%s&proxy=%s&proxy_type=%s&json=1&soft_id=12368652", in.Config.CaptchaSettings.CaptchaAPI, in.Config.CaptchaSettings.ClientKey, sitekey, "https://discord.com/channels/@me", encUa, in.Proxy, proxyType)
		} else {
			submitEndpoint = fmt.Sprintf("http://%s/in.php?key=%s&method=hcaptcha&sitekey=%s&pageurl=%s&userAgent=%s&proxy=%s&proxy_type=%s&json=1&soft_id=13615286", in.Config.CaptchaSettings.CaptchaAPI, in.Config.CaptchaSettings.ClientKey, sitekey, "https://discord.com/channels/@me", encUa, in.Proxy, proxyType)
		}
	}
	if rqData != "" {
		submitEndpoint = fmt.Sprintf("%s&data=%s", submitEndpoint, rqData)
	}
	req, err := http.NewRequest(http.MethodGet, submitEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request [%v]", err)
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request [%v]", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response [%v]", err)
	}
	fmt.Println(string(body))
	if !strings.Contains(string(body), "status") {
		return "", fmt.Errorf("cannot proccess response, it does not contain status [%v] %v", err, string(body))
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response [%v]", err)
	}
	if response["status"].(float64) != 1 {
		return "", fmt.Errorf("error %v", response["request"])
	}
	var captchaIDfloat int
	var captchaIDString string
	var captchaGetEndpoint string
	if in.Config.CaptchaSettings.CaptchaAPI == "azcaptcha.com" {
		captchaIDfloat = int(response["request"].(float64))
		fmt.Println(captchaIDfloat)
		captchaGetEndpoint = fmt.Sprintf("https://%s/res.php?key=%s&action=get&action=get&id=%s&json=1", in.Config.CaptchaSettings.CaptchaAPI, in.Config.CaptchaSettings.ClientKey, strconv.Itoa(captchaIDfloat))
	} else {
		captchaIDString = response["request"].(string)
		captchaGetEndpoint = fmt.Sprintf("https://%s/res.php?key=%s&action=get&action=get&id=%s&json=1", in.Config.CaptchaSettings.CaptchaAPI, in.Config.CaptchaSettings.ClientKey, captchaIDString)
	}
	fmt.Println(captchaGetEndpoint)
	// time recommended in rucaptcha documentation
	time.Sleep(15 * time.Second)

	for i := 0; i < 100; i++ {
		req, err := http.NewRequest(http.MethodGet, captchaGetEndpoint, nil)
		if err != nil {
			return "", fmt.Errorf("error creating request [%v]", err)
		}
		req.Header.Set("User-Agent", UserAgent)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", fmt.Errorf("error sending request [%v]", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading response [%v]", err)
		}
		fmt.Println(string(body))
		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", fmt.Errorf("error unmarshalling response [%v]", err)
		}
		if response["request"] == "CAPCHA_NOT_READY" {
			time.Sleep(10 * time.Second)
			continue
		} else {
			return response["request"].(string), nil
		}
	}
	return "", fmt.Errorf("max retries exceeded")
}

type twoCaptchaSubmitResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

func (in *Instance) twoCaptcha(sitekey, rqdata string) (string, error) {
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
	q.Set("pageurl", "https://discord.com/channels/@me")
	q.Set("userAgent", UserAgent)
	q.Set("json", "1")
	q.Set("soft_id", "3359")
	if rqdata != "" {
		q.Set("data", rqdata)
	}
	if in.Config.ProxySettings.ProxyForCaptcha {
		q.Set("proxy", in.Proxy)
		q.Set("proxytype", in.Config.ProxySettings.ProxyProtocol)
	}
	inURL.RawQuery = q.Encode()
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
