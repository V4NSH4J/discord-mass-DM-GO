// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

type ImageEmbed struct {
	Url       string `json:"url"`
	Thumbnail bool   `json:"thumbnail"`
}
type ProviderEmbed struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}
type AuthorEmbed struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}
type Embed struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Color       string        `json:"color"`
	Redirect    string        `json:"redirect"`
	Author      AuthorEmbed   `json:"author"`
	Image       ImageEmbed    `json:"image"`
	Provider    ProviderEmbed `json:"provider"`
}

type EmbedJSONResponse struct {
	EmbedData Embed  `json:"embed"`
	Timestamp string `json:"timestamp"`
	Id        string `json:"id"`
	Link      string `json:"link"`
}

func LanuchEmbed() {
	utilities.LogWarn("This feature is provided and hosted by a 3rd party Entity. Use at your own discretion. Contact https://github.com/itschasa for more information.")
	var embeddata []byte
	var err string
	embeddata, err = utilities.GetEmbed()
	if err == "" {
		responseBody := bytes.NewBuffer(embeddata)
		resp, err := http.Post("https://e.chasa.wtf/api/v1/embed", "application/json", responseBody)
		if err != nil {
			utilities.LogErr("Error while getting response from 3rd Party Embed API", err)
		} else {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				utilities.LogErr("Error while reading response body", err)
			} else {
				var respdata EmbedJSONResponse
				err := json.Unmarshal(bodyBytes, &respdata)
				if err != nil {
					utilities.LogErr("Error while unmarshalling response body", err)
				} else {
					if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
						utilities.LogInfo("Created Embed Link, use the link below and add it to your message in message.json.")
						utilities.LogInfo(respdata.Link)
						utilities.LogInfo("Make sure to restart DMDGO after editing message.json")
						utilities.LogInfo("Service provided with <3 by chasa (https://github.com/itschasa)")
					} else {
						utilities.LogErr("Unexpected response from server: %v [Try again or the API might be down] Response :%v", fmt.Sprint(resp.StatusCode), string(bodyBytes))
					}
				}
			}
		}
	} else {
		utilities.LogErr("Error while getting embed data from embed.json %v", err)
	}
}
