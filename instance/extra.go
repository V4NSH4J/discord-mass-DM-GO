// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"net/url"

	gohttp "net/http"
	http "github.com/Danny-Dasilva/fhttp"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func GetReactions(channel string, message string, token string, emoji string, after string) ([]string, error) {
	encodedID := url.QueryEscape(emoji)
	site := "https://discord.com/api/v9/channels/" + channel + "/messages/" + message + "/reactions/" + encodedID + "?limit=100"
	if after != "" {
		site += "&after=" + after
	}

	req, err := gohttp.NewRequest("GET", site, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)

	resp, err := gohttp.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reactions []Reactionx

	fmt.Println(string(body))
	err = json.Unmarshal(body, &reactions)
	if err != nil {
		return nil, err
	}
	var UIDS []string
	for i := 0; i < len(reactions); i++ {
		UIDS = append(UIDS, reactions[i].ID)
	}

	return UIDS, nil
}

func (in *Instance) ContextProperties(invite, cookie string) (string, error) {
	site := "https://discord.com/api/v9/invites/" + invite + "?inputValue=" + invite + "&with_counts=true&with_expiration=true"
	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		return "", err
	}
	req = in.xContextPropertiesHeaders(req, cookie)
	resp, err := in.Client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error while getting invite context %v", resp.StatusCode)
	}
	body, err := utilities.ReadBody(*resp)
	if err != nil {
		return "", err
	}
	if !strings.Contains(string(body), "guild") && !strings.Contains(string(body), "id") && !strings.Contains(string(body), "channel") && !strings.Contains(string(body), "code") {
		return "", fmt.Errorf("Error while getting invite context %v", resp.StatusCode)
	}
	var guildInfo map[string]interface{}
	err = json.Unmarshal(body, &guildInfo)
	if err != nil {
		return "", err
	}
	guildID := (guildInfo["guild"].(map[string]interface{}))["id"].(string)
	channelID := (guildInfo["channel"].(map[string]interface{}))["id"].(string)
	channelType := (guildInfo["channel"].(map[string]interface{}))["type"].(float64)
	x, err := XContextGen(guildID, channelID, channelType)
	if err != nil {
		return "", err
	}
	return x, nil

}

func XContextGen(guildID string, channelID string, ChannelType float64) (string, error) {
	xcontext := XContext{
		Location:            "Join Guild",
		LocationGuildID:     guildID,
		LocationChannelID:   channelID,
		LocationChannelType: ChannelType,
	}
	jsonData, err := json.Marshal(xcontext)
	if err != nil {
		return "", err
	}
	Enc := b64.StdEncoding.EncodeToString(jsonData)
	return Enc, nil

}

func Bypass(client *http.Client, serverid string, token string, invite string) error {
	// First we require to get all the rules to send in the request
	site := "https://discord.com/api/v9/guilds/" + serverid + "/member-verification?with_guild=false&invite_code=" + invite
	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", token)
	resp, err := client.Do(CommonHeaders(req))
	if err != nil {
		return err
	}

	body, err := utilities.ReadBody(*resp)
	if err != nil {
		return err
	}
	var bypassInfo bypassInformation
	err = json.Unmarshal(body, &bypassInfo)
	if err != nil {
		return err
	}

	// Now we have all the rules, we can send the request along with our response
	for i := 0; i < len(bypassInfo.FormFields); i++ {
		// We set the response to true because we accept the terms as the good TOS followers we are
		bypassInfo.FormFields[i].Response = true
	}

	jsonData, err := json.Marshal(bypassInfo)
	if err != nil {
		return err
	}
	url := "https://discord.com/api/v9/guilds/" + serverid + "/requests/@me"

	req, err = http.NewRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		utilities.LogErr("Error while making http request %v \n", err)
		return err
	}

	req.Header.Set("authorization", token)
	resp, err = client.Do(CommonHeaders(req))
	if err != nil {
		utilities.LogErr("Error while sending HTTP request bypass %v \n", err)
		return err
	}
	body, err = utilities.ReadBody(*resp)
	if err != nil {
		utilities.LogErr("[%v] Error while reading body %v \n", time.Now().Format("15:04:05"), err)
		return err
	}

	if resp.StatusCode == 201 || resp.StatusCode == 204 {
		utilities.LogSuccess("[%v] Successfully bypassed token %v", time.Now().Format("15:04:05"), token)
	} else {
		utilities.LogErr("[%v] Failed to bypass Token %v %v %v", time.Now().Format("15:04:05"), token, resp.StatusCode, string(body))
	}
	return nil
}

func (in *Instance) Invite(Code string) error {
	var solvedKey string
	var payload invitePayload
	var rqData string
	var rqToken string
	var j int
	var reported []string
	for i := 0; i < in.Config.CaptchaSettings.MaxCaptchaInv; i++ {
		if solvedKey == "" || in.Config.CaptchaSettings.CaptchaAPI == "" {
			payload = invitePayload{}
		} else {
			payload = invitePayload{
				CaptchaKey: solvedKey,
				RqToken:    rqToken,
			}
		}
		payload, err := json.Marshal(payload)
		if err != nil {
			utilities.LogErr("error while marshalling payload %v", err)
			continue
		}

		cookie, err := in.GetCookieString()
		if err != nil {
			utilities.LogErr("[%v] Error while Getting cookies: %v", time.Now().Format("15:04:05"), err)
			continue
		}
		XContext, err := in.ContextProperties(Code, cookie)
		if err != nil {
			XContext = ""
		}
		url := fmt.Sprintf("https://discord.com/api/v9/invites/%s", Code)
		req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
		if err != nil {
			utilities.LogErr("Error while making http request %v \n", err)
			continue
		}

		req = in.inviteHeaders(req, cookie, XContext)

		resp, err := in.Client.Do(req)
		if err != nil {
			utilities.LogErr("Error while sending HTTP request %v \n", err)
			continue
		}

		body, err := utilities.ReadBody(*resp)
		if err != nil {
			utilities.LogErr("Error while reading body %v \n", err)
			continue
		}
		if strings.Contains(string(body), "captcha_sitekey") {
			if j > 1 {
				if in.Config.CaptchaSettings.CaptchaAPI == "anti-captcha.com" && in.LastID != 0 && !utilities.Contains(reported, string(in.LastID)) {
					reported = append(reported, string(in.LastID))
					err := in.ReportIncorrectRecaptcha()
					if err != nil {
						utilities.LogErr("[%v] Error while reporting incorrect hcaptcha: %v", time.Now().Format("15:04:05"), err)
					} else {
						utilities.LogSuccess("[%v] Succesfully reported incorrect hcaptcha [%v]", time.Now().Format("15:04:05"), in.LastID)
					}
				}
			}
			var resp map[string]interface{}
			err = json.Unmarshal(body, &resp)
			if err != nil {
				utilities.LogErr("[%v] Error while Unmarshalling body: %v", time.Now().Format("15:04:05"), err)
				continue
			}
			cap := resp["captcha_sitekey"].(string)
			if strings.Contains(string(body), "captcha_rqdata") {
				rqData = resp["captcha_rqdata"].(string)
			}
			if strings.Contains(string(body), "captcha_rqtoken") {
				rqToken = resp["captcha_rqtoken"].(string)
			}
			if in.Config.CaptchaSettings.CaptchaAPI == "" && in.Config.CaptchaSettings.CaptchaAPI != "invisifox.com"{
				utilities.LogErr("Captcha detected but no API key provided %v", in.CensorToken())
				break
			} else {
				utilities.CaptchaDetected(in.CensorToken(), cap)
			}
			solvedKey, err = in.SolveCaptcha(cap, cookie, rqData, rqToken, "https://discord.com/channels/@me")
			if err != nil {
				utilities.LogErr("[%v] Error while Solving Captcha: %v", time.Now().Format("15:04:05"), err)
				continue
			}
			j++
			continue
		}
		if resp.StatusCode == 429 && strings.Contains(string(body), "1015") {
			utilities.LogErr("Cloudflare Error 1015 - Your IP is being Rate Limited. Use proxies. If you already are, make sure proxy_from_file is enabled in your config")
			break
		}

		var Join joinresponse
		err = json.Unmarshal(body, &Join)
		if err != nil {
			utilities.LogErr("Error while unmarshalling body %v %v\n", err, string(body))
			return err
		}
		if resp.StatusCode == 200 {
			utilities.LogSuccess("%v joined guild %v", in.CensorToken(), Code)
			if Join.VerificationForm {
				if len(Join.GuildObj.ID) != 0 {
					Bypass(in.Client, Join.GuildObj.ID, in.Token, Code)
				}
			}
		}
		if resp.StatusCode != 200 {
			utilities.LogErr("[%v] %v Failed to join guild %v", time.Now().Format("15:04:05"), resp.StatusCode, string(body))
		}
		return nil

	}
	return fmt.Errorf("max retries exceeded")

}

func (in *Instance) Leave(serverid string) int {
	url := "https://discord.com/api/v9/users/@me/guilds/" + serverid
	json_data := "{\"lurking\":false}"
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte(json_data)))
	if err != nil {
		utilities.LogErr("Error: %s", err)
		return 0
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return 0
	}
	req.Header.Set("authorization", in.Token)
	req.Header.Set("Cookie", cookie)
	resp, errq := in.Client.Do(CommonHeaders(req))
	if errq != nil {
		fmt.Println(errq)
		return 0
	}
	return resp.StatusCode
}

func (in *Instance) React(channelID string, MessageID string, Emoji string) error {
	encodedID := url.QueryEscape(Emoji)
	site := "https://discord.com/api/v9/channels/" + channelID + "/messages/" + MessageID + "/reactions/" + encodedID + "/@me"
	req, err := http.NewRequest("PUT", site, nil)
	if err != nil {
		return err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return fmt.Errorf("error while getting cookie %v", err)
	}
	resp, err := in.Client.Do(in.AtMeHeaders(req, cookie))
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("%s", resp.Status)
}

func (in *Instance) Friend(Username string, Discrim int) (*http.Response, error) {

	url := "https://discord.com/api/v9/users/@me/relationships"

	fr := friendRequest{Username, Discrim}
	jsonx, err := json.Marshal(&fr)
	if err != nil {
		return &http.Response{}, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonx)))
	if err != nil {
		return &http.Response{}, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return &http.Response{}, fmt.Errorf("error while getting cookie %v", err)
	}

	resp, err := in.Client.Do(in.AtMeHeaders(req, cookie))

	if err != nil {
		return &http.Response{}, err
	}
	return resp, nil

}

func FindMessage(channel string, messageid string, token string) (string, error) {
	url := "https://discord.com/api/v9/channels/" + channel + "/messages?limit=1&around=" + messageid
	req, err := gohttp.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", token)

	client := gohttp.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var message []Message
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &message)
	if err != nil {
		fmt.Println(string(body))
		return "", err
	}
	msg, err := json.Marshal(message[0])
	if err != nil {
		fmt.Println(string(body))
		return "", err
	}
	return string(msg), nil
}

func GetRxn(channel string, messageid string, token string) (Message, error) {
	url := "https://discord.com/api/v9/channels/" + channel + "/messages?limit=1&around=" + messageid
	req, err := gohttp.NewRequest("GET", url, nil)
	if err != nil {
		return Message{}, err
	}

	req.Header.Set("Authorization", token)

	client := gohttp.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	var message []Message
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return Message{}, err
	}

	err = json.Unmarshal(body, &message)
	if err != nil {
		return Message{}, err
	}

	return message[0], nil
}

func (in *Instance) ServerCheck(serverid string) (int, error) {
	url := "https://discord.com/api/v9/guilds/" + serverid
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}

	req.Header.Set("Authorization", in.Token)

	client := in.Client
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (in *Instance) EndRelation(snowflake string) (int, error) {
	link := fmt.Sprintf("https://discord.com/api/v9/users/@me/relationships/%s", snowflake)
	req, err := http.NewRequest("DELETE", link, nil)
	if err != nil {
		return -1, err
	}
	cookies, err := in.GetCookieString()
	if err != nil {
		return -1, err
	}
	req = in.AtMeHeaders(req, cookies)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, err
	}
	return resp.StatusCode, nil
}

func (in *Instance) PressButton(actionRow, button int, guildID string, msg Message) (int, error) {
	site := "https://discord.com/api/v9/interactions"
	data := map[string]interface{}{"component_type": msg.Components[actionRow].Buttons[button].Type, "custom_id": msg.Components[actionRow].Buttons[button].CustomID, "hash": msg.Components[actionRow].Buttons[button].Hash}
	values := map[string]interface{}{"application_id": msg.Author.ID, "channel_id": msg.ChannelID, "type": 3, "data": data, "guild_id": guildID, "message_flags": msg.Flags, "message_id": msg.MessageId, "nonce": utilities.Snowflake(), "session_id": in.Ws.sessionID}
	jsonData, err := json.Marshal(values)
	if err != nil {
		return -1, err
	}
	req, err := http.NewRequest("POST", site, strings.NewReader(string(jsonData)))
	if err != nil {
		return -1, err
	}
	cookies, err := in.GetCookieString()
	if err != nil {
		return -1, err
	}
	req = in.AtMeHeaders(req, cookies)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}
	fmt.Println(string(body))
	return resp.StatusCode, nil

}
