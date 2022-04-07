package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "io/ioutil"
	"math/rand"
	"net/http"

	// "regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func (in *Instance) GetCookieString() (string, error) {
	url := "https://discord.com"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red("[%v] Error while making request to get cookies %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while making request to get cookie %v", err)
	}
	req = in.cookieHeaders(req)
	resp, err := in.Client.Do(req)
	if err != nil {
		color.Red("[%v] Error while getting response from cookies request %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while getting response from cookie request %v", err)
	}
	defer resp.Body.Close()

	if resp.Cookies() == nil {
		color.Red("[%v] Error while getting cookies from response %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("there are no cookies in response")
	}
	cookies := ""
	for _, cookie := range resp.Cookies() {
		cookies += fmt.Sprintf(`%s=%s; `, cookie.Name, cookie.Value)
	}
	cookies += "locale=en-US"
	// CfRay := resp.Header.Get("cf-ray")
	// if strings.Contains(CfRay, "-BOM") {
	// 	CfRay = strings.ReplaceAll(CfRay, "-BOM", "")
	// }

	// if CfRay != "" {
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		color.Red("[%v] Error while reading response body %v", time.Now().Format("15:04:05"), err)
	// 		return cookies + "locale:en-US", nil
	// 	}
	// 	m := regexp.MustCompile(`m:'(.+)'`)
	// 	match := m.FindStringSubmatch(string(body))
	// 	if match == nil {
	// 		return cookies + "locale:en-US", nil
	// 	}
	// 	finalCookies, err := in.GetCfBm(match[1], CfRay, cookies)
	// 	if err != nil {
	// 		return cookies + "locale:en-US", nil
	// 	}
	// 	fmt.Println(finalCookies)
	// 	return finalCookies, nil
	// }

	return cookies, nil
}

func (in *Instance) GetCfBm(m, r, cookies string) (string, error) {
	site := fmt.Sprintf(`https://discord.com/cdn-cgi/bm/cv/result?req_id=%s`, r)
	res := RandomResult()
	payload := fmt.Sprintf(
		`
		{
			"m":"%s",
			"results":["%s","%s"],
			"timing":%v,
			"fp":
				{
					"id":3,
					"e":{"r":[1920,1080],
					"ar":[1032,1920],
					"pr":1,
					"cd":24,
					"wb":true,
					"wp":false,
					"wn":false,
					"ch":false,
					"ws":false,
					"wd":false
				}
			}
		}
		`, m, res[0], res[1], 60+rand.Intn(60),
	)
	req, err := http.NewRequest("POST", site, strings.NewReader(payload))
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("error while making request to get cf-bm %v", err)
	}
	req = in.cfBmHeaders(req, cookies)
	resp, err := in.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while getting response from cf-bm request %v", err)
	}
	defer resp.Body.Close()
	if resp.Cookies() == nil {
		color.Red("[%v] Error while getting cookies from response %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("there are no cookies in response")
	}
	if len(resp.Cookies()) == 0 {
		return cookies, nil
	}
	cookies = cookies + "; "
	for _, cookie := range resp.Cookies() {
		cookies = cookies + cookie.Name + "=" + cookie.Value
	}
	return cookies, nil
}

type response struct {
	Fingerprint string `json:"fingerprint"`
}

func (in *Instance) GetFingerprintString(cookie string) (string, error) {
	url := "https://discord.com/api/v9/experiments"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red("[%v] Error while making request to get fingerprint %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while making request to get fingerprint %v", err)
	}
	req = in.fingerprintHeaders(req, cookie)
	resp, err := in.Client.Do(req)
	if err != nil {
		color.Red("[%v] Error while getting response from fingerprint request %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while getting response from fingerprint request %v", err)
	}

	p, err := ReadBody(*resp)
	if err != nil {
		color.Red("[%v] Error while reading body from fingerprint request %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while reading body %v", err)
	}

	var Response response

	err = json.Unmarshal(p, &Response)

	if err != nil {
		color.Red("[%v] Error while unmarshalling body from fingerprint request %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("error while unmarshalling response from fingerprint request %v", err)
	}

	return Response.Fingerprint, nil
}

func (in *Instance) OpenChannel(recepientUID string) (string, error) {
	url := "https://discord.com/api/v9/users/@me/channels"

	json_data := []byte("{\"recipients\":[\"" + recepientUID + "\"]}")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("Error while making request")
		return "", fmt.Errorf("error while making open channel request %v", err)
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return "", fmt.Errorf("error while getting cookie %v", err)
	}

	resp, err := in.Client.Do(in.OpenChannelHeaders(req, cookie))
	if err != nil {
		return "", fmt.Errorf("error while getting response from open channel request %v", err)
	}
	defer resp.Body.Close()

	body, err := ReadBody(*resp)
	if err != nil {
		return "", fmt.Errorf("error while reading body from open channel request %v", err)
	}
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		color.Red("[%v] Token %v has been locked or disabled", time.Now().Format("15:04:05"), in.Token)
		return "", fmt.Errorf("token has been locked or disabled")
	}
	if resp.StatusCode != 200 {
		fmt.Printf("[%v]Invalid Status Code while sending request %v \n", time.Now().Format("15:04:05"), resp.StatusCode)
		return "", fmt.Errorf("invalid status code while sending request %v", resp.StatusCode)
	}
	type responseBody struct {
		ID string `json:"id,omitempty"`
	}

	var channelSnowflake responseBody
	errx := json.Unmarshal(body, &channelSnowflake)
	if errx != nil {
		return "", errx
	}

	return channelSnowflake.ID, nil
}

type captchaDetected struct {
	CaptchaKey []string `json:"captcha_key"`
	Sitekey    string   `json:"captcha_sitekey"`
	Service    string   `json:"captcha_service"`
	RqData     string   `json:"captcha_rqdata"`
	RqToken    string   `json:"captcha_rqtoken"`
}

func (in *Instance) SendMessage(channelSnowflake string, memberid string) (http.Response, error) {
	// Sending a random message incase there are multiple.
	index := rand.Intn(len(in.Messages))
	message := in.Messages[index]
	x := message.Content
	if strings.Contains(message.Content, "<user>") {
		ping := "<@" + memberid + ">"
		x = strings.ReplaceAll(message.Content, "<user>", ping)
	}

	body, err := json.Marshal(&map[string]interface{}{
		"content": x,
		"tts":     false,
		"nonce":   Snowflake(),
	})
	if err != nil {
		return http.Response{}, fmt.Errorf("error while marshalling message %v %v ", index, err)
	}

	url := "https://discord.com/api/v9/channels/" + channelSnowflake + "/messages"

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return http.Response{}, fmt.Errorf("error while making request to send message %v", err)
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return http.Response{}, fmt.Errorf("error while getting cookie %v", err)
	}
	res, err := in.Client.Do(in.SendMessageHeaders(req, cookie, channelSnowflake))
	if err != nil {
		fmt.Printf("[%v]Error while sending http request %v \n", time.Now().Format("15:04:05"), err)
		return http.Response{}, fmt.Errorf("error while getting send message response %v", err)
	}
	if res.StatusCode == 400 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return http.Response{}, fmt.Errorf("error while reading body %v", err)
		}
		if strings.Contains(string(body), "captcha") {
			color.Yellow("[%v] Captcha detected %v Solving", time.Now().Format("15:04:05"), in.Token)
		}
		if in.Config.CaptchaSettings.ClientKey == "" {
			return http.Response{}, fmt.Errorf("captcha detected but no client key set")
		}
		var captchaDetect captchaDetected
		err = json.Unmarshal(body, &captchaDetect)
		if err != nil {
			return http.Response{}, fmt.Errorf("error while unmarshalling captcha %v", err)
		}
		solved, err := in.SolveCaptcha(captchaDetect.Sitekey, cookie, captchaDetect.RqData, captchaDetect.RqToken, fmt.Sprintf("https://discord.com/channels/@me/%s", channelSnowflake))
		if err != nil {
			return http.Response{}, fmt.Errorf("error while solving captcha %v", err)
		}
		body, err = json.Marshal(&map[string]interface{}{
			"content":         x,
			"tts":             false,
			"nonce":           Snowflake(),
			"captcha_key":     solved,
			"captcha_rqtoken": captchaDetect.RqToken,
		})
		if err != nil {
			return http.Response{}, fmt.Errorf("error while marshalling message %v %v ", index, err)
		}
		req, err = http.NewRequest("POST", url, strings.NewReader(string(body)))
		if err != nil {
			return http.Response{}, fmt.Errorf("error while making request to send message %v", err)
		}
		res, err = in.Client.Do(in.SendMessageHeaders(req, cookie, channelSnowflake))
		if err != nil {
			return http.Response{}, fmt.Errorf("error while getting send message response %v", err)
		}
	}
	in.Count++
	return *res, nil
}

type UserInf struct {
	User   User     `json:"user"`
	Mutual []Guilds `json:"mutual_guilds"`
}

type Guilds struct {
	ID string `json:"id"`
}

func (in *Instance) UserInfo(userid string) (UserInf, error) {
	url := "https://discord.com/api/v9/users/" + userid + "/profile?with_mutual_guilds=true"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return UserInf{}, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return UserInf{}, fmt.Errorf("error while getting cookie %v", err)
	}
	fingerprint, err := in.GetFingerprintString(cookie)
	if err != nil {
		return UserInf{}, fmt.Errorf("error while getting fingerprint %v", err)
	}
	req.Header.Set("Authorization", in.Token)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("x-fingerprint", fingerprint)
	req.Header.Set("host", "discord.com")

	resp, err := in.Client.Do(CommonHeaders(req))
	if err != nil {
		return UserInf{}, err
	}

	body, err := ReadBody(*resp)
	if err != nil {
		return UserInf{}, err
	}

	if body == nil {
		return UserInf{}, fmt.Errorf("body is nil")
	}

	var info UserInf
	errx := json.Unmarshal(body, &info)
	if errx != nil {
		fmt.Println(string(body))
		return UserInf{}, errx
	}
	return info, nil
}

type RingData struct {
	Recipients interface{} `json:"recipients"`
}

func Ring(httpClient *http.Client, auth string, snowflake string) (int, error) {
	url := "https://discord.com/api/v9/channels/" + snowflake + "/call"

	p := RingData{
		Recipients: nil,
	}
	jsonx, err := json.Marshal(&p)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonx)))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ReadBody(*resp)
	if err != nil {
		return 0, err
	}
	fmt.Println(string(body))
	return resp.StatusCode, nil
}

func Snowflake() int64 {
	snowflake := strconv.FormatInt((time.Now().UTC().UnixNano()/1000000)-1420070400000, 2) + "0000000000000000000000"
	nonce, _ := strconv.ParseInt(snowflake, 2, 64)
	return nonce
}

func CommonHeaders(req *http.Request) *http.Request {
	req.Header.Set("X-Super-Properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDAzIiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTA0OTY3LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	req.Header.Set("sec-fetch-dest", "empty")
	// req.Header.Set("Connection", "keep-alive")
	req.Header.Set("x-debug-options", "bugReporterEnabled")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("X-Discord-Locale", "en-US")
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("accept-language", "en-US")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0")
	req.Header.Set("TE", "trailers")
	return req
}

func RegisterHeaders(req *http.Request) *http.Request {
	req.Header.Set("accept", "*/*")
	req.Header.Set("authority", "discord.com")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/api/v9/auth/register")
	req.Header.Set("scheme", "https")
	// req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-Discord-Locale", "en-US")
	req.Header.Set("origin", "discord.com")
	req.Header.Set("referer", "discord.com/register")
	req.Header.Set("x-debug-options", "bugReporterEnabled")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-Type", "application/json")
	// Imitating Discord Desktop Client
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9003 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36")
	req.Header.Set("x-super-properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDAzIiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTA0OTY3LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")

	return req
}

func (in *Instance) CloseDMS(snowflake string) (int, error) {
	site := "https://discord.com/api/v9/channels/" + snowflake
	req, err := http.NewRequest("DELETE", site, nil)
	if err != nil {
		return -1, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, err
	}
	fingerprint, err := in.GetFingerprintString(cookie)
	if err != nil {
		return -1, err
	}
	req.Header.Set("cookie", cookie)
	req.Header.Set("X-Fingerprint", fingerprint)
	req.Header.Set("Authorization", in.Token)
	req = CommonHeaders(req)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, err
	}
	return resp.StatusCode, nil
}

func (in *Instance) BlockUser(userid string) (int, error) {
	site := "https://discord.com/api/v9/users/@me/relationships/" + userid
	payload := `{"type":2}`
	req, err := http.NewRequest("PUT", site, strings.NewReader(payload))
	if err != nil {
		return -1, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, err
	}
	fingerprint, err := in.GetFingerprintString(cookie)
	if err != nil {
		return -1, err
	}
	req.Header.Set("cookie", cookie)
	req.Header.Set("X-Fingerprint", fingerprint)
	req.Header.Set("Authorization", in.Token)
	req = CommonHeaders(req)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, err
	}
	return resp.StatusCode, nil
}

func (in *Instance) greet(channelid, cookie, fingerprint string) (string, error) {
	site := fmt.Sprintf(`https://discord.com/api/v9/channels/%s/greet`, channelid)
	payload := `{"sticker_ids":["749054660769218631"]}`
	req, err := http.NewRequest("POST", site, strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	req = in.SendMessageHeaders(req, cookie, channelid)
	resp, err := in.Client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf(`invalid status code while sending dm %v`, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	var msgid string
	if strings.Contains(string(body), "id") {
		msgid = response["id"].(string)
	} else {
		return "", fmt.Errorf(`invalid response %v`, string(body))
	}
	return msgid, nil
}

func (in *Instance) ungreet(channelid, cookie, fingerprint, msgid string) error {
	site := fmt.Sprintf(`https://discord.com/api/v9/channels/%s/messages/%s`, channelid, msgid)
	req, err := http.NewRequest("DELETE", site, nil)
	if err != nil {
		return err
	}
	req = in.SendMessageHeaders(req, cookie, channelid)
	resp, err := in.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf(`invalid status code while sending dm%v`, resp.StatusCode)
	}
	return nil
}
