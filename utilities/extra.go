// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/fatih/color"
)

type Reactionx struct {
	ID string `json:"id"`
}

func GetReactions(channel string, message string, token string, emoji string, after string) ([]string, error) {
	encodedID := url.QueryEscape(emoji)
	site := "https://discord.com/api/v9/channels/" + channel + "/messages/" + message + "/reactions/" + encodedID + "?limit=100"
	if after != "" {
		site += "&after=" + after
	}

	req, err := http.NewRequest("GET", site, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(CommonHeaders(req))
	if err != nil {
		return nil, err
	}
	body, err := ReadBody(*resp)
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

type guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type joinresponse struct {
	VerificationForm bool  `json:"show_verification_form"`
	GuildObj         guild `json:"guild"`
}

func Bypass(client *http.Client, serverid string, token string) error {
	url := "https://discord.com/api/v9/guilds/" + serverid + "/requests/@me"
	json_data := "{\"response\":true}"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(json_data)))
	if err != nil {
		color.Red("Error while making http request %v \n", err)
		return err
	}

	req.Header.Set("authorization", token)
	resp, err := client.Do(CommonHeaders(req))
	if err != nil {
		color.Red("Error while sending HTTP request bypass %v \n", err)
		return err
	}
	body, err := ReadBody(*resp)
	if err != nil {
		color.Red("[%v] Error while reading body %v \n", time.Now().Format("15:04:05"), err)
		return err
	}

	if resp.StatusCode == 201 || resp.StatusCode == 204 {
		color.Green("[%v] Successfully bypassed token %v", time.Now().Format("15:04:05"), token)
	} else {
		color.Red("[%v] Failed to bypass Token %v %v %v", time.Now().Format("15:04:05"), token, resp.StatusCode, string(body))
	}
	return nil
}

func (in *Instance) Invite(Code string) error {
	url := "https://discord.com/api/v9/invites/" + Code

	var headers struct{}
	requestBytes, _ := json.Marshal(headers)

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
	if err != nil {
		color.Red("Error while making http request %v \n", err)
		return err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return fmt.Errorf("error while getting cookie %v", err)
	}
	fingerprint, err := in.GetFingerprintString()
	if err != nil {
		return fmt.Errorf("error while getting fingerprint %v", err)
	}
	req.Header.Set("authorization", in.Token)
	req.Header.Set("cookie", cookie)
	req.Header.Set("x-fingerprint", fingerprint)
	// Not constant but discord doesn't care. (yet)
	req.Header.Set("x-context-properties", "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6IjkxMzQ2MDQxNzUzMDA2OTAyMiIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI5MTM0NjA0MTc1MzAwNjkwMjUiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9")

	resp, err := in.Client.Do(CommonHeaders(req))
	if err != nil {
		color.Red("Error while sending HTTP request %v \n", err)
		return err
	}

	body, err := ReadBody(*resp)
	if err != nil {
		color.Red("Error while reading body %v \n", err)
		return err
	}

	var Join joinresponse
	err = json.Unmarshal(body, &Join)
	if err != nil {
		color.Red("Error while unmarshalling body %v \n", err)
		return err
	}
	if resp.StatusCode == 200 {
		color.Green("[%v] %v joint guild", time.Now().Format("15:04:05"), in.Token)
		if Join.VerificationForm {
			if len(Join.GuildObj.ID) != 0 {
				Bypass(in.Client, Join.GuildObj.ID, in.Token)
			}
		}
	}
	if resp.StatusCode != 200 {
		color.Red("[%v] %v Failed to join guild %v", time.Now().Format("15:04:05"), resp.StatusCode, string(body))
	}
	return nil

}

func (in *Instance) Leave(serverid string) int {
	url := "https://discord.com/api/v9/users/@me/guilds/" + serverid
	json_data := "{\"lurking\":false}"
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte(json_data)))
	if err != nil {
		color.Red("Error: %s", err)
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
	req.Header.Set("Authorization", in.Token)
	req.Header.Set("Cookie", cookie)

	resp, err := in.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("%s", resp.Status)
}

type friendRequest struct {
	Username string `json:"username"`
	Discrim  int    `json:"discriminator"`
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
	fingerprint, err := in.GetFingerprintString()
	if err != nil {
		return &http.Response{}, fmt.Errorf("error while getting fingerprint %v", err)
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("x-fingerprint", fingerprint)
	req.Header.Set("Authorization", in.Token)

	resp, err := in.Client.Do(CommonHeaders(req))

	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil

}

func (in *Instance) CheckToken() int {
	url := "https://discord.com/api/v9/users/@me/affinities/guilds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1
	}
	req.Header.Set("authorization", in.Token)

	resp, err := in.Client.Do(CommonHeaders(req))
	if err != nil {
		return -1
	}
	return resp.StatusCode

}

func FindMessage(channel string, messageid string, token string) (string, error) {
	url := "https://discord.com/api/v9/channels/" + channel + "/messages?limit=1&around=" + messageid
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", token)

	client := http.DefaultClient
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
		return "", err
	}
	msg, err := json.Marshal(message[0])
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

func GetRxn(channel string, messageid string, token string) (Message, error) {
	url := "https://discord.com/api/v9/channels/" + channel + "/messages?limit=1&around=" + messageid
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Message{}, err
	}

	req.Header.Set("Authorization", token)

	client := http.DefaultClient
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
