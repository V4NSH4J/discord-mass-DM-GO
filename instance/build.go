// Credits: https://github.com/bytixo/

package instance

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

var (
	buildNumber = make(map[string]string)
	//buildHash   map[string]string
)

// UpdateDiscordBuildInfo fetches the latest build info for every discord builds
// such as PTB, Canary and Stable. see https://github.com/KiyonoKara/Discord-Build-Info-PY/blob/main/discord_build_info_py/clientInfo.py
func UpdateDiscordBuildInfo() error {
	jsFileRegex := regexp.MustCompile(`([a-zA-z0-9]+)\.js`)
	buildInfoRegex := regexp.MustCompile(`Build Number: [0-9]+, Version Hash: [A-Za-z0-9]+`)

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Get("https://discord.com/app")
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	r := jsFileRegex.FindAllString(string(body), -1)
	asset := r[len(r)-1]

	resp, err := client.Get("https://discord.com/assets/" + asset)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	z := buildInfoRegex.FindAllString(string(b), -1)
	e := strings.ReplaceAll(z[0], " ", "")
	buildInfos := strings.Split(e, ",")

	buildNum := strings.Split(buildInfos[0], ":")
	buildNumber["stable"] = buildNum[len(buildNum)-1]

	utilities.LogInfo("Fetched Latest Build Info")

	return nil
}

// GetDiscordBuildNumber returns the current buildNumber for the specified version
func GetDiscordBuildNumber(discord string) string {
	return buildNumber[discord]
}

/*
// GetDiscordBuildHash returns the current buildHash for the specified version
func GetDiscordBuildHash(discord string) string {
	return buildHash[discord]
}
*/

// func main() {
// 	err := UpdateDiscordBuildInfo()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(GetDiscordBuildNumber("stable"))
// }
