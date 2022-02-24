package utilities

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

const Useragent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.1013 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36"
const XSuper = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJwdGIiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC4xMDEzIiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTE1NjMzLCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ=="

// New Invite joiner
func (in *Instance) Inviter(invitationCode string, mode int, cookie string, fingerprint string) (int, string, string, error) {
	// Need X-Context-Properties
	var contextProperties string = "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6IjM5MjQyMTM5MzgwMDg4ODMyMSIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI5MDE0MTY2NjE1OTExNjI5NDIiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9"
	var err error
	// Getting cookies to set in our requests subsequently, if error while getting cookies, do not repeat requests.
	if cookie != "" {
		cookie, err = in.GetCookieString()
		if err != nil {
			return -1, "", "", fmt.Errorf("error while getting cookies %v", err)
		}
	}
	if fingerprint != "" {
		fingerprint, err = in.GetFingerprintString(cookie)
		if err != nil {
			return -1, cookie, "", fmt.Errorf("error while getting fingerprint %v", err)
		}
	}
	var guildID string
	var channelID string
	site := fmt.Sprintf(`https://ptb.discord.com/api/v10/invites/%s?inputValue=%s&with_counts=true&with_expiration=true`, invitationCode, invitationCode)
	// Return -2 error incase the error was while getting the X-Context-Properties
	if mode != -2 {
		req, err := http.NewRequest(http.MethodGet, site, nil)
		if err != nil {
			return -2, cookie, fingerprint, fmt.Errorf("error while making request for getting x-context-properties %v", err)
		}
		req = in.xContextHeaders(req, cookie, fingerprint)
		resp, err := in.Client.Do(req)
		if err != nil {
			return -2, cookie, fingerprint, fmt.Errorf("error while getting response for x-context-properties %v", err)
		}
		if resp.StatusCode != 200 {
			return -2, cookie, fingerprint, fmt.Errorf("invalid status code for x-context-properties response: %v", resp.StatusCode)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return -2, cookie, fingerprint, fmt.Errorf("error while reading xcontext response body %v", err)
		}
		if !strings.Contains(string(body), "guild") || !strings.Contains(string(body), "channel") {
			return -2, cookie, fingerprint, fmt.Errorf("xcontext response body does not contain necessary information %v", string(body))
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return -2, cookie, fingerprint, fmt.Errorf("error while unmarshalling xcontext response body: %v", err)
		}
		// Getting guild and channel ID from the response to generate the x-context-proprties
		guildID = (response["guild"].(map[string]interface{}))["id"].(string)
		channelID = (response["channel"].(map[string]interface{}))["id"].(string)
		channelType := (response["channel"].(map[string]interface{}))["type"].(float64)
		contextProperties = generateXContext(channelID, channelType, guildID)
	}
	site = fmt.Sprintf(`https://ptb.discord.com/api/v10/invites/%v`, invitationCode)
	if in.Config.MaxInvite < 2 {
		in.Config.MaxInvite = 2
	}
	var captchaKey string
	var req *http.Request
	for i := 0; i < in.Config.MaxInvite; i++ {
		if captchaKey == "" {
			req, err = http.NewRequest(http.MethodPost, site, strings.NewReader(`{}`))
		} else {
			// To Do: Check if headers change when repeating req with captcha key
			req, err = http.NewRequest(http.MethodPost, site, strings.NewReader(fmt.Sprintf(`{"captcha_key": "%s"}`, captchaKey)))
		}
		if err != nil {
			return -3, cookie, fingerprint, fmt.Errorf("error while making request for joining server %v", err)
		}
		req = in.inviteHeaders(req, cookie, fingerprint, contextProperties)
		resp, err := in.Client.Do(req)
		if err != nil {
			return -3, cookie, fingerprint, fmt.Errorf("error while getting response for joining server %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return -3, cookie, fingerprint, fmt.Errorf("error while reading response for joining server %v", err)
		}
		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return -3, cookie, fingerprint, fmt.Errorf("error while unmarshalling resposne %v", err)
		}
		if strings.Contains(string(body), "captcha_sitekey") {
			captchaSitekey := response["captcha_sitekey"].(string)
			color.Yellow("[%v] Token %v Captcha Detected [%v] [%v]", time.Now().Format("15:04:05"), in.Token, captchaSitekey, i)
			if in.Config.ClientKey == "" {
				return -3, cookie, fingerprint, fmt.Errorf("captcha detected but no api provided")
			}
			captchaKey, err = in.SolveCaptcha(captchaSitekey, cookie)
			if err != nil {
				color.Yellow("[%v] Token %v error while solving captcha %v Retrying", time.Now().Format("15:04:05"), in.Token, err)
			}
			continue
		}
		if strings.Contains(string(body), "1015") {
			return -4, cookie, fingerprint, fmt.Errorf("you're being rate limited. Cloudflare error 1015 - Use proxies or VPN")
		}
		fmt.Println(string(body))
		// Handle more status codes
		guildID = (response["guild"].(map[string]interface{}))["id"].(string)
		channelID = (response["channel"].(map[string]interface{}))["id"].(string)
		if resp.StatusCode == 200 {
			// Succesfully joint server
			if response["show_verification_form"].(bool) {
				// need to bypass community screening
				req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(`https://ptb.discord.com/api/v10/guilds/%s/member-verification?with_guild=false&invite_code=%s`, guildID, invitationCode), nil)
				if err != nil {
					return -5, cookie, fingerprint, fmt.Errorf("failed to load membership screening bypass form")
				}
				req = in.headersRules(req, cookie, fingerprint, guildID, channelID)
				resp, err := in.Client.Do(req)
				if err != nil {
					return -5, cookie, fingerprint, fmt.Errorf("failed to load membership screening bypass form")
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return -5, cookie, fingerprint, fmt.Errorf("error while loading membership screening bypass form")
				}
				if resp.StatusCode != 200 {
					return -5, cookie, fingerprint, fmt.Errorf("invalid status code while trying to get the membership screening forum %v", resp.StatusCode)
				}
				var bypassInfo bypassInformation
				err = json.Unmarshal(body, &bypassInfo)
				if err != nil {
					return -5, cookie, fingerprint, fmt.Errorf("error while unmarshalling membership screening form %v", err)
				}
				// Now we have all the rules, we can send the request along with our response
				for i := 0; i < len(bypassInfo.FormFields); i++ {
					// We set the response to true because we accept the terms as the good TOS followers we are
					bypassInfo.FormFields[i].Response = true
				}
				payloadBypass, err := json.Marshal(bypassInfo)
				if err != nil {
					return -5, cookie, fingerprint, fmt.Errorf("error while marshalling membership screening form %v", err)
				}
				req, err = http.NewRequest(http.MethodPut, fmt.Sprintf(`https://discord.com/api/v10/guilds/%s/requests/@me`, guildID), bytes.NewReader(payloadBypass))
				if err != nil {
					return -5, cookie, fingerprint, fmt.Errorf("error while making community screening bypass request %v", err)
				}
				req = in.headersRules(req, cookie, fingerprint, guildID, channelID)
				req.Header.Set("Content-Type", "application/json")
				resp, err = in.Client.Do(req)
				if err != nil {
					return -5, cookie, fingerprint, fmt.Errorf("error while getting response for community screening bypass %v", err)
				}
				if resp.StatusCode == 201 {
					return resp.StatusCode, cookie, fingerprint, nil
				} else {
					return -5, cookie, fingerprint, fmt.Errorf("invalid status code while trying to bypass community screening %v", resp.StatusCode)
				}
			} else {
				return resp.StatusCode, "", "", nil
			}
		}
	}
	return -6, "", "", fmt.Errorf("captcha max retries exceeded")
}

func (in *Instance) xContextHeaders(req *http.Request, cookie, fingerprint string) *http.Request {

	for k, v := range map[string]string{
		"Host":               "ptb.discord.com",
		"Connection":         "keep-alive",
		"X-Super-Properties": XSuper,
		"X-Discord-Locale":   "en-US",
		"X-Debug-Options":    "bugReporterEnabled",
		"X-Fingerprint":      fingerprint,
		"Accept-Language":    "en-US,en-IN;q=0.9,zh-Hans-CN;q=0.8",
		"Authorization":      in.Token,
		"User-Agent":         Useragent,
		"Accept":             "*/*",
		"Sec-Fetch-Site":     "same-origin",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Dest":     "empty",
		"Referer":            "https://ptb.discord.com/channels/@me",
		"Cookie":             fmt.Sprintf(cookie+"%s%s", fakeStripeMid(), fakeCFruid()),
		//"Cookie": "__dcfduid=23b89500957e11eca6dff1afbbb1a2bb; __sdcfduid=23b89501957e11eca6dff1afbbb1a2bbec44de1e329adfb0cd13bb8fb79436d8682a43efa5b029585307695946cda945",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) inviteHeaders(req *http.Request, cookie string, fingerprint string, xcontext string) *http.Request {

	for k, v := range map[string]string{
		"Host":                 "ptb.discord.com",
		"Connection":           "keep-alive",
		"X-Super-Properties":   XSuper,
		"X-Context-Properties": xcontext,
		"X-Debug-Options":      "bugReporterEnabled",
		"X-Fingerprint":        fingerprint,
		"Accept-Language":      "en-US,en-IN;q=0.9,zh-Hans-CN;q=0.8",
		"Authorization":        in.Token,
		"Content-Type":         "application/json",
		"User-Agent":           Useragent,
		"X-Discord-Locale":     "en-US",
		"Accept":               "*/*",
		"Origin":               "https://ptb.discord.com",
		"Sec-Fetch-Site":       "same-origin",
		"Sec-Fetch-Mode":       "cors",
		"Sec-Fetch-Dest":       "empty",
		"Referer":              "https://ptb.discord.com/channels/@me",
		"Cookie":               fmt.Sprintf(cookie+"%s%s", fakeStripeMid(), fakeCFruid()),
		//"Cookie": "__dcfduid=23b89500957e11eca6dff1afbbb1a2bb; __sdcfduid=23b89501957e11eca6dff1afbbb1a2bbec44de1e329adfb0cd13bb8fb79436d8682a43efa5b029585307695946cda945",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) headersRules(req *http.Request, cookie string, fingerprint string, serverID string, channelID string) *http.Request {

	for k, v := range map[string]string{
		"Host":               "ptb.discord.com",
		"Connection":         "keep-alive",
		"X-Super-Properties": XSuper,
		"X-Discord-Locale":   "en-US",
		"X-Debug-Options":    "bugReporterEnabled",
		"X-Fingerprint":      fingerprint,
		"Accept-Language":    "en-US,en-IN;q=0.9,zh-Hans-CN;q=0.8",
		"Authorization":      in.Token,
		"User-Agent":         Useragent,
		"Accept":             "*/*",
		"Sec-Fetch-Site":     "same-origin",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Dest":     "empty",
		"Cookie":             fmt.Sprintf(cookie+"%s%s", fakeStripeMid(), fakeCFruid()),
		//"Cookie": "__dcfduid=23b89500957e11eca6dff1afbbb1a2bb; __sdcfduid=23b89501957e11eca6dff1afbbb1a2bbec44de1e329adfb0cd13bb8fb79436d8682a43efa5b029585307695946cda945",
	} {
		req.Header.Set(k, v)
	}
	if channelID != "" && serverID != "" {
		req.Header.Set("Referer", fmt.Sprintf("https://ptb.discord.com/channels/%s/%s", serverID, channelID))
	}
	return req
}

func fakeStripeMid() string {
	return fmt.Sprintf(` __stripe_mid=%s-%s-%s-%s-%s`, randomString(8), randomString(4), randomString(4), randomString(4), randomString(18))
}
func fakeCFruid() string {
	return fmt.Sprintf(`; __cfruid=%s-%s`, randomString(40), randomIntegerString(10))
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	var rndm string
	for i := 0; i < length; i++ {
		rndm += string(charset[rand.Intn(len(charset))])
	}
	return rndm
}
func randomIntegerString(length int) string {
	numset := "0123456789"
	var rndm string
	for i := 0; i < length; i++ {
		rndm += string(numset[rand.Intn(len(numset))])
	}
	return rndm
}

func generateXContext(channelID string, channelType float64, guildID string) string {
	dec := fmt.Sprintf(`{"location":"Join Guild","location_guild_id":"%s","location_channel_id":"%s","location_channel_type":%s}`, guildID, channelID, channelType)
	return base64.StdEncoding.EncodeToString([]byte(dec))
}
