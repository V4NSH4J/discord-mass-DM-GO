package directmessage

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
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
type Message struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

// Inputs the Channel snowflake and sends them the message; outputs the response code for error handling.
func SendMessage(authorization string, channelSnowflake string, message *Message, memberid string) *http.Response {
	x := message.Content
	if strings.Contains(message.Content, "<user>") {
		ping := "<@" + memberid + ">"
		x = strings.ReplaceAll(message.Content, "<user>", ping)
	}

	body, err := json.Marshal(&map[string]interface{}{
		"content": x,
		"embeds":  message.Embeds,
		"tts":     false,
		"nonce":   utilities.Snowflake(),
	})

	if err != nil {
		log.Panicln("Error while marshalling message content")
	}

	url := "https://discord.com/api/v9/channels/" + channelSnowflake + "/messages"

	Cookie := utilities.GetCookie()
	if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
		fmt.Println("ERR: Empty cookie")
	}

	Cookies := "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + " locale=us" + "; __cfruid=d2f75b0a2c63c38e6b3ab5226909e5184b1acb3e-1634536904"

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))

	if err != nil {
		log.Panicf("Error while making HTTP request")
	}

	req.Header.Add("Authorization", authorization)
	req.Header.Add("referer", "ttps://discord.com/channels/@me/"+channelSnowflake)
	req.Header.Set("Cookie", Cookies)
	req.Header.Set("x-fingerprint", utilities.GetFingerprint())
	res, err := http.DefaultClient.Do(utilities.CommonHeaders(req))

	if err != nil {
		log.Panicf("[%v]Error while sending http request %v \n", time.Now().Format("15:05:04"), err)
	}

	return res
}
