// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"encoding/json"
	"fmt"
)

func (in *Instance) CheckToken() int {
	url := "https://discord.com/api/v9/users/@me/affinities/guilds"
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1
	}
	resp, err := in.Client.Do(url, in.CycleOptions("", in.AtMeHeaders(cookie)), "GET")
	if err != nil {
		return -1
	}
	return resp.Status

}

func (in *Instance) CheckTokenNew() (int, error) {
	url := "https://discord.com/api/v9/users/@me/affinities/guilds"
	cookie, err := in.GetCookieString()
	if err != nil {
		return 0, err
	}
	resp, err := in.Client.Do(url, in.CycleOptions("", in.AtMeHeaders(cookie)), "GET")
	if err != nil {
		return 0, err
	}
	return resp.Status, nil

}

func (in *Instance) AtMe() (int, TokenInfo, error) {
	url := "https://discord.com/api/v9/users/@me"
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, TokenInfo{}, fmt.Errorf("error while getting cookie %v", err)
	}
	resp, err := in.Client.Do(url, in.CycleOptions("", in.AtMeHeaders(cookie)), "GET")
	if err != nil {
		return -1, TokenInfo{}, fmt.Errorf("error while sending request %v", err)
	}
	body := resp.Body
	var info TokenInfo
	err = json.Unmarshal([]byte(body), &info)
	if err != nil {
		return -1, TokenInfo{}, fmt.Errorf("error while unmarshalling response %v", err)
	}
	return resp.Status, info, nil
}

func (in *Instance) Guilds() (int, int, []string, error) {
	url := "https://discord.com/api/v9/users/@me/guilds"
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while getting cookie %v", err)
	}
	resp, err := in.Client.Do(url, in.CycleOptions("", in.AtMeHeaders(cookie)), "GET")
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while sending request %v", err)
	}
	body := resp.Body
	var info []Guilds
	err = json.Unmarshal([]byte(body), &info)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while unmarshalling response %v", err)
	}
	var guilds []string
	for i := 0; i < len(info); i++ {
		guilds = append(guilds, info[i].ID)
	}
	return resp.Status, len(info), guilds, nil
}

func (in *Instance) Channels() (int, int, []Guilds, error) {
	url := "https://discord.com/api/v9/users/@me/channels"
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while getting cookie %v", err)
	}
	resp, err := in.Client.Do(url, in.CycleOptions("", in.AtMeHeaders(cookie)), "GET")
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while sending request %v", err)
	}
	body := resp.Body
	var info []Guilds
	err = json.Unmarshal([]byte(body), &info)
	if err != nil {
		return -1, -1, nil, fmt.Errorf("error while unmarshalling response %v", err)
	}
	return resp.Status, len(info), info, nil
}

func (in *Instance) Relationships() (int, int, int, int, int, []Guilds, error) {
	url := "https://discord.com/api/v9/users/@me/relationships"
	cookie, err := in.GetCookieString()
	if err != nil {
		return -1, -1, -1, -1, -1, nil, fmt.Errorf("error while getting cookie %v", err)
	}
	resp, err := in.Client.Do(url, in.CycleOptions("", in.AtMeHeaders(cookie)), "GET")
	if err != nil {
		return -1, -1, -1, -1, -1, nil, fmt.Errorf("error while sending request %v", err)
	}
	body := resp.Body
	var info []Guilds
	err = json.Unmarshal([]byte(body), &info)
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
	return resp.Status, friend, blocked, incoming, outgoing, info, nil

}
