// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"net/http"
)

func CheckToken(auth string) int {
	url := "https://discord.com/api/v9/users/@me/affinities/guilds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1
	}
	req.Close = true
	req.Header.Set("authorization", auth)
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(CommonHeaders(req))
	if err != nil {
		return -1
	}

	return resp.StatusCode

}
