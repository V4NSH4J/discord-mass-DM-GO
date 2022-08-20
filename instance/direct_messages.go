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
	"math"
	"math/rand"
	"strings"
	"time"

	http "github.com/Danny-Dasilva/fhttp"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func cookieToString(cookies []*http.Cookie, cookieString string) string {
	for i := 0; i < len(cookies); i++ {
		if i == len(cookies)-1 {
			cookieString += fmt.Sprintf(`%s=%s`, cookies[i].Name, cookies[i].Value)
		} else {
			cookieString += fmt.Sprintf(`%s=%s; `, cookies[i].Name, cookies[i].Value)
		}
	}
	if !strings.Contains(cookieString, "locale=en-US; ") {
		cookieString += "; locale=en-US "
	}
	return cookieString
}

func (in *Instance) GetCookieString() (string, error) {
	if in.Config.OtherSettings.ConstantCookies && in.Cookie != "" {
		return in.Cookie, nil
	}
	cookies := []*http.Cookie{}
	link := "https://discord.com"
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", fmt.Errorf("error while making request to get cookies %v", err)
	}
	req = in.cookieHeaders(req)
	resp, err := in.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while getting response from cookies request %v", err)
	}
	defer resp.Body.Close()
	if resp.Cookies() == nil {
		utilities.LogErr("[%v] Error while getting cookies from response %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("there are no cookies in response")
	}
	if len(resp.Cookies()) == 0 {
		return "", fmt.Errorf("there are no cookies in response")
	}
	cookies = append(cookies, resp.Cookies()...)
	in.Cookie = cookieToString(cookies, "")
	return cookieToString(cookies, ""), nil
}

func (in *Instance) GetCfBm(m, r, cookies string) (string, error) {
	site := fmt.Sprintf(`https://discord.com/cdn-cgi/bm/cv/result?req_id=%s`, r)
	payload := fmt.Sprintf(
		`
		{
			"m":"%s",
			"results":["859fe3e432b90450c6ddf8fae54c9a58","460d5f1e93f296a48e3f6745675f27e2"],
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
		`, m, 60+rand.Intn(60),
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
		utilities.LogErr("[%v] Error while getting cookies from response %v", time.Now().Format("15:04:05"), err)
		return "", fmt.Errorf("there are no cookies in response")
	}
	if len(resp.Cookies()) == 0 {
		return cookies, nil
	}
	for _, cookie := range resp.Cookies() {
		cookies = cookies + cookie.Name + "=" + cookie.Value
	}
	return cookies, nil

}

func (in *Instance) OpenChannel(recepientUID string) (string, error) {
	url := "https://discord.com/api/v9/users/@me/channels"

	json_data := []byte("{\"recipients\":[\"" + recepientUID + "\"]}")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("Error while making request")
		return "", fmt.Errorf("error while making open channel request %v", err)
	}
	var cookie string
	if in.Cookie == "" {
		cookie, err = in.GetCookieString()
		if err != nil {
			return "", fmt.Errorf("error while getting cookie %v", err)
		}
	} else {
		cookie = in.Cookie
	}

	resp, err := in.Client.Do(in.OpenChannelHeaders(req, cookie))

	if err != nil {
		return "", fmt.Errorf("error while getting response from open channel request %v", err)
	}
	defer resp.Body.Close()

	body, err := utilities.ReadBody(*resp)
	if err != nil {
		return "", fmt.Errorf("error while reading body from open channel request %v", err)
	}
	if strings.Contains(string(body), "Your bot has been flagged") {
		utilities.LogErr("[%v] Token %v has been quarantined", time.Now().Format("15:04:05"), in.CensorToken())
		in.Quarantined = true 
		return "", fmt.Errorf("token has been quarantined")
	}
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		utilities.LogErr("[%v] Token %v has been locked or disabled", time.Now().Format("15:04:05"), in.CensorToken())
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

// Inputs the Channel snowflake and sends them the message; outputs the response code for error handling.
func (in *Instance) SendMessage(channelSnowflake string, memberid string) (int, []byte, error) {
	// Sending a random message incase there are multiple.
	index := rand.Intn(len(in.Messages))
	if in.Config.DirectMessage.MultipleMessages {
		if in.MessageNumber < len(in.Messages) {
			index = in.MessageNumber
		} else {
			return 0, nil, fmt.Errorf("sent all messages")
		}
	}
	message := in.Messages[index]
	x := message.Content
	if strings.Contains(message.Content, "<user>") {
		ping := "<@" + memberid + ">"
		x = strings.ReplaceAll(message.Content, "<user>", ping)
	}

	payload, err := json.Marshal(&map[string]interface{}{
		"content": x,
		"tts":     false,
		"nonce":   utilities.Snowflake(),
	})
	if err != nil {
		return -1, nil, fmt.Errorf("error while marshalling message %v %v ", index, err)
	}

	url := "https://discord.com/api/v9/channels/" + channelSnowflake + "/messages"
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))

	if err != nil {
		return -1, nil, fmt.Errorf("error while making request to send message %v", err)
	}
	var cookie string
	if in.Cookie == "" {
		cookie, err = in.GetCookieString()
		if err != nil {
			return -1, nil, fmt.Errorf("error while getting cookie %v", err)
		}
	} else {
		cookie = in.Cookie
	}

	if in.Config.SuspicionAvoidance.Typing {
		dur := typingSpeed(x, in.Config.SuspicionAvoidance.TypingVariation, in.Config.SuspicionAvoidance.TypingSpeed, in.Config.SuspicionAvoidance.TypingBase)
		if dur != 0 {
			iterations := int((int64(dur) / int64(time.Second*10))) + 1
			for i := 0; i < iterations; i++ {
				if err := in.typing(channelSnowflake, cookie); err != nil {
					continue
				}
				s := time.Second * 10
				if i == iterations-1 {
					s = dur % time.Second * 10
				}
				time.Sleep(s)
			}
		}
	}

	res, err := in.Client.Do(in.SendMessageHeaders(req, cookie, channelSnowflake))
	if err != nil {
		fmt.Printf("[%v]Error while sending http request %v \n", time.Now().Format("15:04:05"), err)
		return -1, nil, fmt.Errorf("error while getting send message response %v", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, fmt.Errorf("error while reading body %v", err)
	}
	t := res.StatusCode
	if res.StatusCode == 200 || res.StatusCode == 204 {
		if in.Config.DirectMessage.MultipleMessages {
			go func() {
				for {
					in.MessageNumber++
					if in.MessageNumber < len(in.Messages) {
						time.Sleep(time.Second * time.Duration(in.Config.DirectMessage.DelayBetweenMultipleMessages))
						status, body, err := in.SendMessage(channelSnowflake, memberid)
						if err != nil {
							utilities.LogFailed("%v Error while sending message %v \n", in.CensorToken(), err)
							continue
						} else {
							utilities.LogSuccess("%v Message #%v sent successfully %v %v\n", in.CensorToken(), in.MessageNumber+1, status, body)
						}
					}
				}
			}()
		}
	}
	if res.StatusCode == 400 || res.StatusCode == 403 {
		if !strings.Contains(string(body), "captcha") {
			return res.StatusCode, body, nil
		}
		if in.Config.CaptchaSettings.ClientKey == "" && in.Config.CaptchaSettings.CaptchaAPI != "invisifox.com" {
			return res.StatusCode, body, fmt.Errorf("captcha detected but no client key set")
		}
		var captchaDetect captchaDetected
		err = json.Unmarshal(body, &captchaDetect)
		if err != nil {
			return res.StatusCode, body, fmt.Errorf("error while unmarshalling captcha %v", err)
		}
		utilities.CaptchaDetected(in.CensorToken(), captchaDetect.Sitekey)
		solved, err := in.SolveCaptcha(captchaDetect.Sitekey, cookie, captchaDetect.RqData, captchaDetect.RqToken, fmt.Sprintf("https://discord.com/channels/@me/%s", channelSnowflake))
		if err != nil {
			return res.StatusCode, body, fmt.Errorf("error while solving captcha %v", err)
		}
		payload, err = json.Marshal(&map[string]interface{}{
			"content":         x,
			"tts":             false,
			"nonce":           utilities.Snowflake(),
			"captcha_key":     solved,
			"captcha_rqtoken": captchaDetect.RqToken,
		})
		if err != nil {
			return res.StatusCode, body, fmt.Errorf("error while marshalling message %v %v ", index, err)
		}
		req, err = http.NewRequest("POST", url, strings.NewReader(string(payload)))
		if err != nil {
			return res.StatusCode, body, fmt.Errorf("error while making request to send message %v", err)
		}
		res, err = in.Client.Do(in.SendMessageHeaders(req, cookie, channelSnowflake))
		if err != nil {
			return t, body, fmt.Errorf("error while getting send message response %v", err)
		}
	}
	in.Count++
	return res.StatusCode, body, nil
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

	resp, err := in.Client.Do(in.AtMeHeaders(req, cookie))
	if err != nil {
		return UserInf{}, err
	}

	body, err := utilities.ReadBody(*resp)
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

func (in *Instance) Ring(snowflake string) (int, error) {

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

	req.Header.Set("Authorization", in.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := in.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := utilities.ReadBody(*resp)
	if err != nil {
		return 0, err
	}
	fmt.Println(string(body))
	return resp.StatusCode, nil

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
	resp, err := in.Client.Do(in.AtMeHeaders(req, cookie))
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
	resp, err := in.Client.Do(in.AtMeHeaders(req, cookie))
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

func (in *Instance) typing(channelID, cookie string) error {
	reqURL := fmt.Sprintf(`https://discord.com/api/v9/channels/%s/typing`, channelID)
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return err
	}
	req = in.TypingHeaders(req, cookie, channelID)
	resp, err := in.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf(`invalid status code while sending dm%v`, resp.StatusCode)
	}
	return nil
}

func typingSpeed(msg string, TypingVariation, TypingSpeed, TypingBase int) time.Duration {
	msPerKey := int(math.Round((1.0 / float64(TypingSpeed)) * 60000))
	d := TypingBase
	d += len(msg) * msPerKey
	if TypingVariation > 0 {
		d += rand.Intn(TypingVariation)
	}
	return time.Duration(d) * time.Millisecond
}

func (in *Instance) Call(snowflake string) error {
	if in.Ws == nil {
		return fmt.Errorf("websocket is not initialized")
	}
	e := CallEvent{
		Op: 4,
		Data: CallData{
			ChannelId: snowflake,
			GuildId:   nil,
			SelfDeaf:  false,
			SelfMute:  false,
			SelfVideo: false,
		},
	}
	err := in.Ws.WriteRaw(e)
	if err != nil {
		return fmt.Errorf("failed to write to websocket: %s", err)
	}

	return nil
}
