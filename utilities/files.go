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
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	
	"gopkg.in/yaml.v2"
)

func ReadLines(filename string) ([]string, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.OpenFile(path.Join(path.Dir(ex)+"/input/"+filename), os.O_RDWR, 0660)
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
	Content   string     `json:"content,omitempty"`
	Embeds    []Embed    `json:"embeds,omitempty"`
	Reactions []Reaction `json:"reactions,omitempty"`
	Author    User       `json:"author,omitempty"`
	GuildID   string     `json:"guild_id,omitempty"`
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
	Delay           int    `yaml:"individual_delay"`
	LongDelay       int    `yaml:"rate_limit_delay"`
	Offset          int    `yaml:"offset"`
	Skip            bool   `yaml:"skip_completed"`
	Proxy           string `yaml:"proxy"`
	Call            bool   `yaml:"call"`
	Remove          bool   `yaml:"remove_dead_tokens"`
	RemoveM         bool   `yaml:"remove_completed_members"`
	Stop            bool   `yaml:"stop_dead_tokens"`
	Mutual          bool   `yaml:"check_mutual"`
	Friend          bool   `yaml:"friend_before_DM"`
	Websocket       bool   `yaml:"online_tokens"`
	SleepSc         int    `yaml:"online_scraper_delay"`
	ProxyFromFile   bool   `yaml:"proxy_from_file"`
	ProxyProtocol   string `yaml:"proxy_protocol"`
	MaxDMS          int    `yaml:"max_dms_per_token"`
	Receive         bool   `yaml:"receive_messages"`
	GatewayProxy    bool   `yaml:"use_proxy_for_gateway"`
	Timeout         int    `yaml:"timeout"`
	SkipFailed      bool   `yaml:"skip_failed"`
	ClientKey       string `yaml:"captcha_api_key"`
	CaptchaAPI      string `yaml:"captcha_api"`
	MaxInvite       int    `yaml:"max_attempt_invite_rejoin"`
	DisableKL       bool   `yaml:"disable_keep_alives"`
	ScrapeUsernames bool   `yaml:"scrape_usernames"`
	ScrapeAvatars   bool   `yaml:"scrape_avatars"`
	ProxyForCaptcha bool   `yaml:"proxy_for_captcha"`
	Block           bool   `yaml:"block_after_dm"`
	Close           bool   `yaml:"close_dm_after_message"`
}

func GetConfig() (Config, error) {
	ex, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		return Config{}, err
	}
	ex = filepath.ToSlash(ex)
	file, err1 := os.Open(path.Join(path.Dir(ex) + "/" + "config.yml"))
	if err1 != nil {
		color.Red("Error while Opening Config (tried .json + .yml)")
		return Config{}, err
	} else {
		defer file.Close()
		var config Config
		bytes, _ := io.ReadAll(file)
		errr := yaml.Unmarshal(bytes, &config)
		if errr != nil {
			fmt.Println(err)
			return Config{}, err
		}
		return config, nil
	}
}

func ProcessAvatar(av string, memberid string) error {
	if strings.Contains(av, "a_") {
		// Nitro Avatar
		return nil
	}
	link := "https://cdn.discordapp.com/avatars/" + memberid + "/" + av + ".png"
	nameFile := "input/pfps/" + av + ".png"

	err := processFiles(link, nameFile)
	if err != nil {
		return err
	}

	return nil
}

func processFiles(url string, nameFile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected http status code while downloading avatar%d", resp.StatusCode)
	}
	file, err := os.Create(nameFile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
