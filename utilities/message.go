// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
)

type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
type Embed struct {
	Title string `json:"title,omitempty"`

	// The type of embed. Always EmbedTypeRich for webhook embeds.
	Type        string             `json:"type,omitempty"`
	Description string             `json:"description,omitempty"`
	URL         string             `json:"url,omitempty"`
	Image       *MessageEmbedImage `json:"image,omitempty"`

	// The color code of the embed.
	Color     int                    `json:"color,omitempty"`
	Footer    EmbedFooter            `json:"footer,omitempty"`
	Thumbnail *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Provider  EmbedProvider          `json:"provider,omitempty"`
	Author    EmbedAuthor            `json:"author,omitempty"`
	Fields    []EmbedField           `json:"fields,omitempty"`
}
type Emoji struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}
type Reaction struct {
	Emojis Emoji `json:"emoji,omitempty"`
	Count  int   `json:"count,omitempty"`
}

type Message struct {
	Content   string     `json:"content,omitempty"`
	Embeds    []Embed    `json:"embeds,omitempty"`
	Reactions []Reaction `json:"reactions,omitempty"`
}

func GetMessage() (Message, error) {
	var message Message
	ex, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		return Message{}, err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "message.json"))
	if err != nil {
		color.Red("Error while Opening message.json")
		fmt.Println(err)
		return Message{}, err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	errr := json.Unmarshal(bytes, &message)
	if errr != nil {
		fmt.Println(err)

		return Message{}, err
	}

	return message, nil
}
