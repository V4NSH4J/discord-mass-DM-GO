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

	"github.com/fatih/color"
)

func (in *Instance) SolveCaptcha(sitekey string, cookie string, rqData string, rqToken string, url string) (string, error) {
	switch true {
	case Contains([]string{"capmonster.cloud", "anti-captcha.com", "anycaptcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		return in.SolveCaptchaCapmonster(sitekey, cookie, rqData, url)
	case in.Config.CaptchaSettings.CaptchaAPI == "deathbycaptcha.com":
		return in.SolveCaptchaDeathByCaptcha(sitekey, url)
	case Contains([]string{"2captcha.com", "rucaptcha.com"}, in.Config.CaptchaSettings.CaptchaAPI):
		return in.twoCaptcha(sitekey, rqData, url)
	default:
		return "", fmt.Errorf("unsupported captcha api: %s", in.Config.CaptchaSettings.CaptchaAPI)
	}
}

func (in *Instance) SolveCaptchaCapmonster(sitekey string, cookies string, rqdata string, url string) (string, error) {
	var jsonx Pload
	if !in.Config.ProxySettings.ProxyForCaptcha || in.Config.CaptchaSettings.CaptchaAPI == "anycaptcha.com" {
		jsonx = Pload{
			ClientKey: in.Config.CaptchaSettings.ClientKey,
			Task: Task{
				Type:       "HCaptchaTaskProxyless",
				WebsiteURL: url,
				WebsiteKey: sitekey,
				UserAgent:  UserAgent,
				Cookies:    cookies,
				Data:       rqdata,
				Invisible:  true,
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
				WebsiteURL:    url,
				WebsiteKey:    sitekey,
				UserAgent:     UserAgent,
				ProxyType:     in.Config.ProxySettings.ProxyProtocol,
				ProxyAddress:  address,
				ProxyPort:     port,
				ProxyLogin:    username,
				ProxyPassword: password,
				Cookies:       cookies,
				Data:          rqdata,
				Invisible:     true,
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
	Invisible     bool   `json:"isInvisible`
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

func (in *Instance) SolveCaptchaDeathByCaptcha(sitekey, url string) (string, error) {
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
				"pageurl": "%s",
				"sitekey": "%s"
			}
		}`, username, password, proxy, proxytype, url, sitekey)
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
