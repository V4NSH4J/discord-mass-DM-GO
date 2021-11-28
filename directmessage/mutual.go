// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package directmessage

import (
	"encoding/json"

	"net/http"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

type UserInf struct {
	User   User     `json:"user"`
	Mutual []Guilds `json:"mutual_guilds"`
}
type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}
type Guilds struct {
	ID string `json:"id"`
}

func UserInfo(token string, userid string, i int, j int) (UserInf, error) {
	url := "https://discord.com/api/v9/users/" + userid + "/profile?with_mutual_guilds=true"
	cookie, err := utilities.Cookies(i, j)
	if err != nil {
		return UserInf{}, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return UserInf{}, err
	}
	fingerprint, err := utilities.Fingerprint(i, j)
	if err != nil {
		return UserInf{}, err
	}
	req.Close = true
	req.Header.Set("Authorization", token)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("x-fingerprint", fingerprint)

	client, err := utilities.SetProxy(i, j)
	if err != nil {
		return UserInf{}, err
	}

	resp, err := client.Do(utilities.CommonHeaders(req))
	if err != nil {
		return UserInf{}, err
	}
	body, err := utilities.ReadBody(*resp)
	if err != nil {
		return UserInf{}, err
	}

	var info UserInf
	errx := json.Unmarshal(body, &info)
	if errx != nil {
		return UserInf{}, errx
	}
	return info, nil
}
