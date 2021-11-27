// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package directmessage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)


func OpenChannel(authorization string, recepientUID string, i int, j int) (string, error) {
	url := "https://discord.com/api/v9/users/@me/channels"

	json_data := []byte("{\"recipients\":[\"" + recepientUID + "\"]}")


	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("Error while making request")
		return "", err
	}
	req.Close = true
	cookie, err := utilities.Cookies(i, j)
	if err != nil {
		fmt.Println("Error while getting cookie")
		return "", err
	}
	fingerprint, err := utilities.Fingerprint(i, j)
	if err != nil {
		fmt.Println("Error while getting fingerprint")
		return "", err
	}
	req.Header.Set("authorization", authorization)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("x-fingerprint", fingerprint)
	req.Header.Set("x-context-properties", "e30=")
	httpClient, err := utilities.SetProxy(i, j)
	if err != nil {
		fmt.Println("Error while setting proxy")
		return "", err
	}
	resp, err := httpClient.Do(utilities.CommonHeaders(req))

	if err != nil {
		fmt.Printf("Error while sending Open channel request  %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		color.Red("[%v] Token %v has been locked or disabled", time.Now().Format("15:05:04"), authorization)
		return "", fmt.Errorf("locked")
	}
	if resp.StatusCode != 200 {
		fmt.Printf("[%v]Invalid Status Code while sending request %v \n",time.Now().Format("15:05:04"), resp.StatusCode)
		return "", err
	}
	type responseBody struct {
		ID string `json:"id,omitempty"`
	}

	var channelSnowflake responseBody
	json.Unmarshal(body, &channelSnowflake)

	return channelSnowflake.ID, nil
}
