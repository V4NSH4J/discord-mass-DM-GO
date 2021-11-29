// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"fmt"
	"net/http"
	"net/url"
)

func React(token string, channelID string, MessageID string, Emoji string) error {
	encodedID := url.QueryEscape(Emoji)
	site := "https://discord.com/api/v9/channels/" + channelID + "/messages/" + MessageID + "/reactions/" + encodedID + "/@me"

	req, err := http.NewRequest("PUT", site, nil)
	if err != nil {
		return err
	}
	req.Close = true
	cookie, err := Cookies()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Cookie", cookie)

	httpClient := http.DefaultClient

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("%s", resp.Status)
}
