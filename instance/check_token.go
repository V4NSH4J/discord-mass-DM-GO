// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	http "github.com/Danny-Dasilva/fhttp"
)

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

func (in *Instance) CheckTokenNew() (int, error) {
	url := "https://discord.com/api/v9/users/@me/affinities/guilds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return 0, err
	}
	resp, err := in.Client.Do(in.AtMeHeaders(req, cookie))
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil

}

func (in *Instance) AtMe() (int, TokenInfo, error) {
	url := "https://discord.com/api/v9/users/@me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, TokenInfo{}, fmt.Errorf("error while making request %v", err)
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, TokenInfo{}, fmt.Errorf("error while getting cookie %v", err)
	}
	req = in.AtMeHeaders(req, cookie)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, TokenInfo{}, fmt.Errorf("error while sending request %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, TokenInfo{}, fmt.Errorf("error while reading response %v", err)
	}
	var info TokenInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		return -1, TokenInfo{}, fmt.Errorf("error while unmarshalling response %v", err)
	}
	return resp.StatusCode, info, nil
}

func (in *Instance) Guilds() (int, int, []string, error) {
	url := "https://discord.com/api/v9/users/@me/guilds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while making request %v", err)
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while getting cookie %v", err)
	}
	req = in.AtMeHeaders(req, cookie)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while sending request %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while reading response %v", err)
	}
	var info []Guilds
	err = json.Unmarshal(body, &info)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while unmarshalling response %v", err)
	}
	var guilds []string
	for i := 0; i < len(info); i++ {
		guilds = append(guilds, info[i].ID)
	}
	return resp.StatusCode, len(info), guilds, nil
}

func (in *Instance) Channels() (int, int, []Guilds, error) {
	url := "https://discord.com/api/v9/users/@me/channels"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while making request %v", err)
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while getting cookie %v", err)
	}
	req = in.AtMeHeaders(req, cookie)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while sending request %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while reading response %v", err)
	}
	var info []Guilds
	err = json.Unmarshal(body, &info)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while unmarshalling response %v", err)
	}
	return resp.StatusCode, len(info), info, nil
}

func (in *Instance) Relationships() (int, int, int, int, int, []Guilds, error) {
	url := "https://discord.com/api/v9/users/@me/relationships"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, -1, -1, -1, -1, nil, fmt.Errorf("error while making request %v", err)
	}
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, -1, -1, -1, -1, nil, fmt.Errorf("error while getting cookie %v", err)
	}
	req = in.AtMeHeaders(req, cookie)
	resp, err := in.Client.Do(req)
	if err != nil {
		return -1, -1, -1, -1, -1, nil, fmt.Errorf("error while sending request %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, -1, -1, -1, nil, fmt.Errorf("error while reading response %v", err)
	}
	var info []Guilds
	err = json.Unmarshal(body, &info)
	if err != nil {
		return -1, -1, -1, -1, -1, nil, fmt.Errorf("error while unmarshalling response %v", err)
	}
	var friend, blocked, incoming, outgoing int
	for i := 0; i < len(info); i++ {
		if info[i].Type == 1 {
			friend++
		} else if info[i].Type == 2 {
			blocked++
		} else if info[i].Type == 3 {
			incoming++
		} else if info[i].Type == 4 {
			outgoing++
		}
	}
	return resp.StatusCode, friend, blocked, incoming, outgoing, info, nil

}
