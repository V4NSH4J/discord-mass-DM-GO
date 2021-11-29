// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fatih/color"
)

type guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type joinresponse struct {
	VerificationForm bool  `json:"show_verification_form"`
	GuildObj         guild `json:"guild"`
}

func Bypass(serverid string, token string) error {
	url := "https://discord.com/api/v9/guilds/" + serverid + "/requests/@me"
	json_data := "{\"response\":true}"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(json_data)))
	if err != nil {
		color.Red("Error while making http request %v \n", err)
		return err
	}
	req.Close = true
	req.Header.Set("authorization", token)
	httpClient := http.Client{}
	resp, err := httpClient.Do(CommonHeaders(req))
	if err != nil {
		color.Red("Error while sending HTTP request bypass %v \n", err)
		return err
	}
	if resp.StatusCode == 201 || resp.StatusCode == 204 {
		color.Green("Successfully bypassed token")
	} else {
		color.Red("Failed to bypass Token %v", resp.StatusCode)
	}
	return nil
}

func Invite(Code string, token string) error {
	url := "https://discord.com/api/v9/invites/" + Code

	cookie, err := Cookies()
	if err != nil {
		color.Red("Error while getting cookies %v \n", err)
		return err
	}
	var headers struct{}
	requestBytes, _ := json.Marshal(headers)

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
	if err != nil {
		color.Red("Error while making http request %v \n", err)
		return err
	}
	fingerprint, err := Fingerprint()
	if err != nil {
		color.Red("Error while getting fingerprint %v \n", err)
		return err
	}

	req.Close = true

	req.Header.Set("authorization", token)
	req.Header.Set("cookie", cookie)
	req.Header.Set("x-fingerprint", fingerprint)
	req.Header.Set("x-context-properties", "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6IjkxMzQ2MDQxNzUzMDA2OTAyMiIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI5MTM0NjA0MTc1MzAwNjkwMjUiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(CommonHeaders(req))
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
		color.Green("[%v] %v joint guild", time.Now().Format("15:05:04"), token)
		if Join.VerificationForm {
			if len(Join.GuildObj.ID) != 0 {
				Bypass(Join.GuildObj.ID, token)
			}
		}
	}
	if resp.StatusCode != 200 {
		color.Red("[%v] %v Failed to join guild", time.Now().Format("15:05:04"), resp.StatusCode)
	}
	return nil

}
