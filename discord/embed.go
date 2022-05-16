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
	"github.com/fatih/color"
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
	color.Yellow("This feature is provided and hosted by a 3rd party Entity. Use at your own discretion. Contact https://github.com/itschasa for more information.")
	var embeddata []byte
	var err string
	embeddata, err = utilities.GetEmbed()
	if err == "" {
		responseBody := bytes.NewBuffer(embeddata)
		resp, err := http.Post("https://e.chasa.wtf/api/v1/embed", "application/json", responseBody)
		if err != nil {
			color.Red("Error making HTTP request to server")
			color.Red(err.Error())
		} else {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				color.Red("Error unmarhsalling HTTP request")
				color.Red(err.Error())
			} else {
				var respdata EmbedJSONResponse
				err := json.Unmarshal(bodyBytes, &respdata)
				if err != nil {
					color.Red("Error unmarhsalling HTTP request")
					color.Red(err.Error())
				} else {
					if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
						color.Green("Created Embed Link, use the link below and add it to your message in message.json.")
						color.Green(respdata.Link)
						color.Green("Make sure to restart DMDGO after editing message.json")
						color.Green("Service provided with <3 by chasa (https://github.com/itschasa)")
					} else {
						color.Red("Unexpected response from server: %v", fmt.Sprint(resp.StatusCode))
						color.Red(string(bodyBytes))
					}
				}
			}
		}
	} else {
		color.Red(err)
	}
}
