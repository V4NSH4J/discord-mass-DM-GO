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

// Inputs the Channel snowflake and sends them the message; outputs the response code for error handling.
func SendMessage(authorization string, channelSnowflake string, message string) *http.Response {


	body, err := json.Marshal(&map[string]interface{}{
		"content": message,
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
		log.Panicf("[%v]Error while sending http request %v \n",time.Now().Format("15:05:04"), err)
	}



	return res
}
