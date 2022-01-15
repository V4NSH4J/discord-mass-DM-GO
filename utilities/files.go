// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
)

func ReadLines(filename string) ([]string, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/input/" + filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func WriteLines(filename string, line string) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	ex = filepath.ToSlash(ex)
	f, err := os.OpenFile(path.Join(path.Dir(ex)+"/input/"+filename), os.O_RDWR|os.O_APPEND, 0660)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	_, err2 := f.WriteString(line + "\n")
	if err2 != nil {
		log.Fatal(err2)
	}
	return nil

}

func TruncateLines(filename string, line []string) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	ex = filepath.ToSlash(ex)
	f, err := os.OpenFile(path.Join(path.Dir(ex)+"/input/"+filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	for i := 0; i < len(line); i++ {
		_, err2 := f.WriteString(line[i] + "\n")
		if err2 != nil {
			log.Fatal(err2)
		}
	}
	return nil

}

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
	Content   string        `json:"content,omitempty"`
	Embeds    []Embed       `json:"embeds,omitempty"`
	Reactions []Reaction    `json:"reactions,omitempty"`
	Author    User `json:"author,omitempty"`
	GuildID string `json:"guild_id,omitempty"`
}

func GetMessage() ([]Message, error) {
	var messages []Message
	ex, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		return []Message{}, err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "message.json"))
	if err != nil {
		color.Red("Error while Opening message.json")
		fmt.Println(err)
		return []Message{}, err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	errr := json.Unmarshal(bytes, &messages)
	if errr != nil {
		fmt.Println(errr)

		return []Message{}, errr
	}

	return messages, nil
}

type Config struct {
	Delay         int    `json:"individual_delay"`
	LongDelay     int    `json:"rate_limit_delay"`
	Offset        int    `json:"offset"`
	Skip          bool   `json:"skip_completed"`
	Proxy         string `json:"proxy"`
	Call          bool   `json:"call"`
	Remove        bool   `json:"remove_dead_tokens"`
	RemoveM       bool   `json:"remove_completed_members"`
	Stop          bool   `json:"stop_dead_tokens"`
	Mutual        bool   `json:"check_mutual"`
	Friend        bool   `json:"friend_before_DM"`
	Websocket     bool   `json:"online_tokens"`
	SleepSc       int    `json:"online_scraper_delay"`
	ProxyFromFile bool   `json:"proxy_from_file"`
	MaxDMS        int    `json:"max_dms_per_token"`
	Receive 	  bool   `json:"receive_messages"`
	GatewayProxy bool    `json:"use_proxy_for_gateway"`
}

func GetConfig() (Config, error) {
	var config Config
	ex, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		return Config{-1, -1, -1, false, "", false, false, false, false, false, false, false, -1, false, -1, false, false}, err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "config.json"))
	if err != nil {
		color.Red("Error while Opening config.json")
		return Config{-1, -1, -1, false, "", false, false, false, false, false, false, false, -1, false, -1, false, false}, err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	errr := json.Unmarshal(bytes, &config)
	if errr != nil {
		fmt.Println(err)
		return Config{-1, -1, -1, false, "", false, false, false, false, false, false, false, -1, false, -1, false, false}, err
	}

	return Config{config.Delay, config.LongDelay, config.Offset, config.Skip, config.Proxy, config.Call, config.Remove, config.RemoveM, config.Stop, config.Mutual, config.Friend, config.Websocket, config.SleepSc, config.ProxyFromFile, config.MaxDMS, config.Receive, config.GatewayProxy}, nil
}
