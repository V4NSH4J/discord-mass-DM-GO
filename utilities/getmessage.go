// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func FindMessage(channel string, messageid string, token string) (string, error) {
	url := "https://discord.com/api/v9/channels/" + channel + "/messages?limit=1&around=" + messageid
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Close = true
	Cookie, err := Cookies()
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Cookie", Cookie)

	client := &http.Client{}
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
	req.Close = true
	Cookie, err := Cookies()
	if err != nil {
		return Message{}, err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Cookie", Cookie)

	client := &http.Client{}
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
