// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func cookieFromheaders(cookie string, headers map[string]string) string {
	v := headers["Set-Cookie"]
	if v == "" {
		return cookie
	}
	if !strings.Contains(v, ";") && !strings.Contains(v, "/,/") {
		return v
	}
	cookies := strings.Split(v, "/,/")
	for i := 0; i < len(cookies); i++ {
		if i == len(cookies)-1 {
			cookie += fmt.Sprintf(`%s `, strings.Split(cookies[i], ";")[0])
		} else {
			cookie += fmt.Sprintf(`%s; `, strings.Split(cookies[i], ";")[0])
		}
	}
	if !strings.Contains(cookie, "locale=en-US; ") {
		cookie += "; locale=en-US "
	}
	return cookie
}

func (in *Instance) GetCookieString() (string, error) {
	if in.Config.OtherSettings.ConstantCookies && in.Cookie != "" {
		return in.Cookie, nil
	}
	link := "https://discord.com"
	resp, err := in.Client.Do(link, in.CycleOptions("", in.cookieHeaders()), "GET")
	if err != nil {
		return "", fmt.Errorf("error while getting response from cookies request %v", err)
	}
	cookie := cookieFromheaders("", resp.Headers)
	if cookie == "" {
		return "", fmt.Errorf("error while getting cookie from response")
	}
	return cookie, nil
}

func (in *Instance) OpenChannel(recepientUID string) (string, error) {
	url := "https://discord.com/api/v9/users/@me/channels"

	json_data := []byte("{\"recipients\":[\"" + recepientUID + "\"]}")
	var cookie string
	var err error
	if in.Cookie == "" {
		cookie, err = in.GetCookieString()
		if err != nil {
			return "", fmt.Errorf("error while getting cookie %v", err)
		}
	} else {
		cookie = in.Cookie
	}
	resp, err := in.Client.Do(url, in.CycleOptions(string(json_data), in.OpenChannelHeaders(cookie)), "POST")

	if err != nil {
		return "", fmt.Errorf("error while getting response from open channel request %v", err)
	}
	body := resp.Body
	if err != nil {
		return "", fmt.Errorf("error while reading body from open channel request %v", err)
	}
	if resp.Status == 401 || resp.Status == 403 {
		utilities.LogErr("[%v] Token %v has been locked or disabled", time.Now().Format("15:04:05"), in.CensorToken())
		return "", fmt.Errorf("token has been locked or disabled")
	}
	if resp.Status != 200 {
		fmt.Printf("[%v]Invalid Status Code while sending request %v \n", time.Now().Format("15:04:05"), resp.Status)
		return "", fmt.Errorf("invalid status code while sending request %v", resp.Status)
	}
	type responseBody struct {
		ID string `json:"id,omitempty"`
	}

	var channelSnowflake responseBody
	errx := json.Unmarshal([]byte(body), &channelSnowflake)
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

	res, err := in.Client.Do(url, in.CycleOptions(string(payload), in.SendMessageHeaders(cookie, channelSnowflake)), "POST")
	if err != nil {
		fmt.Printf("[%v]Error while sending http request %v \n", time.Now().Format("15:04:05"), err)
		return -1, nil, fmt.Errorf("error while getting send message response %v", err)
	}
	body := res.Body
	t := res.Status
	if res.Status == 200 || res.Status == 204 {
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
	if res.Status == 400 {
		if !strings.Contains(string(body), "captcha") {
			return res.Status, []byte(body), nil
		}
		if in.Config.CaptchaSettings.ClientKey == "" {
			return res.Status, []byte(body), fmt.Errorf("captcha detected but no client key set")
		}
		var captchaDetect captchaDetected
		err = json.Unmarshal([]byte(body), &captchaDetect)
		if err != nil {
			return res.Status, []byte(body), fmt.Errorf("error while unmarshalling captcha %v", err)
		}
		utilities.LogWarn("Captcha detected %v [%v]", in.CensorToken(), captchaDetect.Sitekey)
		solved, err := in.SolveCaptcha(captchaDetect.Sitekey, cookie, captchaDetect.RqData, captchaDetect.RqToken, fmt.Sprintf("https://discord.com/channels/@me/%s", channelSnowflake))
		if err != nil {
			return res.Status, []byte(body), fmt.Errorf("error while solving captcha %v", err)
		}
		payload, err = json.Marshal(&map[string]interface{}{
			"content":         x,
			"tts":             false,
			"nonce":           utilities.Snowflake(),
			"captcha_key":     solved,
			"captcha_rqtoken": captchaDetect.RqToken,
		})
		if err != nil {
			return res.Status, []byte(body), fmt.Errorf("error while marshalling message %v %v ", index, err)
		}
		res, err = in.Client.Do(url, in.CycleOptions(string(payload), in.SendMessageHeaders(cookie, channelSnowflake)), "POST")
		if err != nil {
			return t, []byte(body), fmt.Errorf("error while getting send message response %v", err)
		}
	}
	in.Count++
	return res.Status, []byte(body), nil
}

func (in *Instance) UserInfo(userid string) (UserInf, error) {
	url := "https://discord.com/api/v9/users/" + userid + "/profile?with_mutual_guilds=true"
	cookie, err := in.GetCookieString()
	if err != nil {
		return UserInf{}, fmt.Errorf("error while getting cookie %v", err)
	}
	resp, err := in.Client.Do(url, in.CycleOptions("", in.UserInfoHeaders(cookie)), "GET")
	if err != nil {
		return UserInf{}, err
	}
	body := resp.Body
	var info UserInf
	errx := json.Unmarshal([]byte(body), &info)
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
	resp, err := in.Client.Do(url, in.CycleOptions(string(jsonx), in.AtMeHeaders("")), "POST")
	if err != nil {
		return 0, err
	}

	body := resp.Body
	if err != nil {
		return 0, err
	}
	fmt.Println(string(body))
	return resp.Status, nil

}

func (in *Instance) CloseDMS(snowflake string) (int, error) {
	site := "https://discord.com/api/v9/channels/" + snowflake
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, err
	}
	resp, err := in.Client.Do(site, in.CycleOptions("", in.AtMeHeaders(cookie)), "DELETE")
	if err != nil {
		return -1, err
	}
	return resp.Status, nil
}

func (in *Instance) BlockUser(userid string) (int, error) {
	site := "https://discord.com/api/v9/users/@me/relationships/" + userid
	payload := `{"type":2}`
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, err
	}
	resp, err := in.Client.Do(site, in.CycleOptions(payload, in.AtMeHeaders(cookie)), "PUT")
	if err != nil {
		return -1, err
	}
	return resp.Status, nil
}

func (in *Instance) typing(channelID, cookie string) error {
	reqURL := fmt.Sprintf(`https://discord.com/api/v9/channels/%s/typing`, channelID)
	resp, err := in.Client.Do(reqURL, in.CycleOptions("", in.TypingHeaders(cookie, channelID)), "POST")
	if err != nil {
		return err
	}
	if resp.Status != 204 {
		return fmt.Errorf(`invalid status code while sending dm%v`, resp.Status)
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
