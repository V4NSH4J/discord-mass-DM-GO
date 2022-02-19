package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

const Useragent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9004 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36"
const XSuper = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRGlzY29yZCBDbGllbnQiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfdmVyc2lvbiI6IjEuMC45MDA0Iiwib3NfdmVyc2lvbiI6IjEwLjAuMjIwMDAiLCJvc19hcmNoIjoieDY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTE1NjMzLCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ=="
// Unused invite joiner
func (in *Instance) Inviter(invite string, captchaSolution string, cookie string, fingerprint string) (int, error) {

	site := fmt.Sprintf(`https://discord.com/api/v9/invites/%s`, invite)
	var payload invitePayload
	if captchaSolution == "" {
		payload = invitePayload{}
	} else {
		// Setting captcha in payload if captchaSolution is not empty
		payload = invitePayload{
			CaptchaKey: captchaSolution,
		}
	}
	payloadFinal, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("could not marshal payload %v", err)
	}
	req, err := http.NewRequest("POST", site, strings.NewReader(string(payloadFinal)))
	if err != nil {
		return -1, fmt.Errorf("error while creating invite request: %v", err)
	}
	// Getting cookies to set in the headers
	if cookie == "" {
		cookie, err = in.GetCookieString()
		if err != nil {
			return -1, fmt.Errorf("error while getting cookie for joining server: %v", err)
		}
	}
	if fingerprint == "" {
		fingerprint, err = in.GetFingerprintString()
		if err != nil {
			return -1, fmt.Errorf("error while getting fingerprint for joining server: %v", err)
		}
	}

	req = in.InviterHeaders(req, cookie, fingerprint)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, fmt.Errorf("error while getting response from http request %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("error while reading response body %v", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		if strings.Contains(string(body), "1015") {
			return -1, fmt.Errorf("you're being rate limited. use proxies or vpn")
		}
		return -1, fmt.Errorf("error while unmarshalling response body %v [%v]", err, string(body))
	}
	// Checking if the server join requires a captcha
	if resp.StatusCode == 400 {
		if strings.Contains(string(body), "captcha_sitekey") {
			cap := response["captcha_sitekey"].(string)
			// Solve captcha using pseudo api wrapper defined in captcha.go and repeat the function
			if in.Config.CaptchaAPI == "" {
				return -1, fmt.Errorf("captcha required but no api key provided")
			}
			color.Yellow("[%v] Captcha Detected [%v] %v Solving", time.Now().Format("15:04:05"), cap, in.Token)
			solvedKey, err := in.SolveCaptcha(cap, cookie)
			if err != nil {
				return -1, fmt.Errorf("error while solving captcha %v", err)
			}
			if captchaSolution != "" {
				return -1, fmt.Errorf("failed to join server with solved captcha")
			}
			in.Inviter(invite, solvedKey, cookie, fingerprint)
		}
	}
	if resp.StatusCode != 200 {
		return -1, fmt.Errorf("error while joining server %v %v", resp.StatusCode, string(body))
	}
	// Checking if token requires to bypass community screening; returning -2 subsequently signifies that the guild was joint but bypass failed
	if strings.Contains(string(body), "show_verification_form") {
		if !strings.Contains(string(body), "guild") {
			return -2, fmt.Errorf("join response does not contain guild id")
		}
		// Getting ID from guild object in invite join response
		guildID := (response["guild"].(map[string]interface{}))["id"].(string)
		if !strings.Contains(string(body), "channel") {
			return -2, fmt.Errorf("join response does not contain channel id")
		}
		// Getting ID from channel object in invite join response
		channelID := (response["channel"].(map[string]interface{}))["id"].(string)
		// Getting verification form fields
		site = fmt.Sprintf(`https://discord.com/api/v9/guilds/%s/member-verification?with_guild=false&invite_code=%s`, guildID, invite)
		req, err = http.NewRequest("GET", site, nil)
		if err != nil {
			return -2, fmt.Errorf("error while creating verification form request: %v", err)
		}
		req = in.CommunityScreeningHeaders(req, cookie, guildID, channelID, fingerprint)
		resp, err = in.Client.Do(req)
		if err != nil {
			return -2, fmt.Errorf("error while getting response from community screening request %v", err)
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return -2, fmt.Errorf("error while reading response body from community screening response %v", err)
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return -2, fmt.Errorf("error while unmarshalling response body from community screening response %v [%v]", err, string(body))
		}
		if resp.StatusCode != 200 {
			return -2, fmt.Errorf("invalid status code while accepting community screening rules %v %v", resp.StatusCode, string(body))
		}
		// From the response, we have to edit the form fields array, set the response to true and resend it to another endpoint.
		if !strings.Contains(string(body), "form_fields") {
			return -2, fmt.Errorf("community screening response body does not contain form fields to submit")
		}
		formFields := response["form_fields"].([](map[string]interface{}))
		for i := range formFields {
			formFields[i]["response"] = true
		}
		response["form_fields"] = formFields
		site = fmt.Sprintf(`https://discord.com/api/v9/guilds/%s/requests/@me`, guildID)
		payloadBytes, err := json.Marshal(response)
		if err != nil {
			return -2, fmt.Errorf("error while marshalling community screening payload %v", err)
		}
		req, err := http.NewRequest("POST", site, bytes.NewReader(payloadBytes))
		if err != nil {
			return -2, fmt.Errorf("error while creating community screening request %v", err)
		}
		req = in.CommunityScreeningHeaders(req, cookie, guildID, channelID, fingerprint)
		resp, err = in.Client.Do(req)
		if err != nil {
			return -2, fmt.Errorf("error while getting response from community screening request %v", err)
		}
		if resp.StatusCode != 201 {
			return -2, fmt.Errorf("invalid status code while accepting community screening rules %v %v", resp.StatusCode, string(body))
		}
	}
	return resp.StatusCode, nil
}

// Function to set the headers for joining servers
func (in *Instance) InviterHeaders(req *http.Request, cookie string, fingerprint string) *http.Request {
	for k, v := range map[string]string{
		"Origin":             "https://discord.com",
		"Connection":         "keep-alive",
		"Content-Type":       "application/json",
		"X-Super-Properties": XSuper,
		// Ideally context properties should not be a constant
		// "X-Context-Properties": xcontext,
		"X-Debug-Options":      "bugReporterEnabled",
		"X-Fingerprint":        fingerprint,
		"Authorization":        in.Token,
		"User-Agent":           Useragent,
		"X-Discord-Locale":     "en-US",
		"Accept":               "*/*",
		"Sec-Fetch-Site":       "same-origin",
		"Sec-Fetch-Mode":       "cors",
		"Sec-Fetch-Dest":       "empty",
		"Sec-ch-ua-mobile":     "?0",
		"Sec-ch-ua-platform":   `"Windows"`,
		"Referer":              "https://discord.com/channels/@me",
		"Accept-Language":      "en-GB,en;q=0.9",
		"Cookie":               cookie,
	} {
		req.Header.Set(k, v)
	}
	return req
}

// Function to set headers for community screening acceptance
func (in *Instance) CommunityScreeningHeaders(req *http.Request, cookie string, guildID string, channelID string, fingerprint string) *http.Request {
	for k, v := range map[string]string{
		"Host":               "discord.com",
		"Connection":         "keep-alive",
		"X-Super-Properties": XSuper,
		"X-Debug-Options":    "bugReporterEnabled",
		"X-Fingerprint":      fingerprint,
		"Authorization":      in.Token,
		"User-Agent":         Useragent,
		"X-Discord-Locale":   "en-US",
		"Accept":             "*/*",
		"Sec-Fetch-Site":     "same-origin",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Dest":     "empty",
		"Referer":            fmt.Sprintf(`https://discord.com/channels/%s/%s`, guildID, channelID),
		"Accept-Language":    "en-GB,en;q=0.9",
		"Cookie":             cookie,
	} {
		req.Header.Set(k, v)
	}
	return req
}
