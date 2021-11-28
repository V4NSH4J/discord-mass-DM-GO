// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package directmessage

import (
        "encoding/json"
        "fmt"
        "net/http"
        "time"

        "strings"

        "github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

// Inputs the Channel snowflake and sends them the message; outputs the response code for error handling.
func SendMessage(authorization string, channelSnowflake string, message *utilities.Message, memberid string) (*http.Response, error) {
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
                return nil, err
        }

        url := "https://discord.com/api/v9/channels/" + channelSnowflake + "/messages"
        cookie, err := utilities.Cookies()
        if err != nil {
                fmt.Println("Error while getting cookie")
                return nil, err
        }
        fingerprint, err := utilities.Fingerprint()
        if err != nil {
                fmt.Println("Error while getting fingerprint")
                return nil, err
        }

        req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
        req.Close = true
        if err != nil {
                return nil, err
        }

        req.Header.Add("Authorization", authorization)
        req.Header.Add("referer", "https://discord.com/channels/@me/"+channelSnowflake)
        req.Header.Set("Cookie", cookie)
        req.Header.Set("x-fingerprint", fingerprint)
        req.Header.Set("host", "discord.com")
        req.Header.Set("origin", "https://discord.com")
        httpClient := http.DefaultClient

        res, err := httpClient.Do(utilities.CommonHeaders(req))

        if err != nil {
                fmt.Printf("[%v]Error while sending http request %v \n", time.Now().Format("15:05:04"), err)
                return nil, err
        }

        return res, nil
}