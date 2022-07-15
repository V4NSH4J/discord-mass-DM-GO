// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/gorilla/websocket"
)

type Instance struct {
	Token           string
	Email           string
	Password        string
	Proxy           string
	Cookie          string
	Fingerprint     string
	Messages        []Message
	Count           int
	LastQuery       string
	LastCount       int
	Members         []User
	AllMembers      []User
	Retry           int
	ScrapeCount     int
	ID              string
	Receiver        bool
	Config          Config
	GatewayProxy    string
	Client          *http.Client
	WG              *sync.WaitGroup
	Ws              *Connection
	fatal           chan error
	Invited         bool
	TimeServerCheck time.Time
	ChangedName     bool
	ChangedAvatar   bool
	LastID          int
	LastIDstr       string
	Mode            int
	UserAgent       string
	XSuper          string
	Reacted         []ReactInfo
	ReactChannel    chan (ReactInfo)
	MessageNumber   int
	Version         string
}

func (in *Instance) StartWS() error {
	ws, err := in.NewConnection(in.wsFatalHandler)
	if err != nil {
		return fmt.Errorf("failed to create websocket connection: %s", err)
	}
	in.Ws = ws
	return nil
}

func (in *Instance) wsFatalHandler(err error) {
	if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code == 4004 {
		in.fatal <- fmt.Errorf("websocket closed: authentication failed, try using a new token")
		return
	}
	if strings.Contains(err.Error(), "4004") {
		utilities.LogLocked("Error while opening websocket, Authentication failed %v", in.Token)
		return
	}
	utilities.LogSuccess("Websocket closed %v %v", err, in.Token)
	in.Receiver = false
	in.Ws, err = in.NewConnection(in.wsFatalHandler)
	if err != nil {
		in.fatal <- fmt.Errorf("failed to create websocket connection: %s", err)
		return
	}
	utilities.LogSuccess("Reconnected To Websocket %v", in.Token)
}

func GetEverything() (Config, []Instance, error) {
	var cfg Config
	var instances []Instance
	var err error
	var tokens []string
	var proxies []string
	var proxy string
	var xsuper string
	var ua string
	var v string

	// Load config
	cfg, err = GetConfig()
	if err != nil {
		return cfg, instances, err
	}
	supportedProtocols := []string{"http", "https", "socks4", "socks5"}
	if cfg.ProxySettings.ProxyProtocol != "" && !utilities.Contains(supportedProtocols, cfg.ProxySettings.ProxyProtocol) {
		utilities.LogErr(" You're using an unsupported proxy protocol. Assuming http by default")
		cfg.ProxySettings.ProxyProtocol = "http"
	}
	if cfg.ProxySettings.ProxyProtocol == "https" {
		cfg.ProxySettings.ProxyProtocol = "http"
	}
	if cfg.CaptchaSettings.CaptchaAPI == "" {
		utilities.LogErr(" You're not using a Captcha API, some functionality like invite joining might be unavailable")
	}
	if cfg.ProxySettings.Proxy != "" && os.Getenv("HTTPS_PROXY") == "" {
		os.Setenv("HTTPS_PROXY", cfg.ProxySettings.ProxyProtocol+"://"+cfg.ProxySettings.Proxy)
	}
	if !cfg.ProxySettings.ProxyFromFile && cfg.ProxySettings.ProxyForCaptcha {
		utilities.LogErr(" You must enabe proxy_from_file to use proxy_for_captcha")
		cfg.ProxySettings.ProxyForCaptcha = false
	}

	xsuper, ua, v, err = DolfiesXsuper()
	if err != nil {
		utilities.LogErr(" Failed to get useragent and xsuper %s Turn off Dolfies mode or try again, otherwise program will be continued with hardcoded chrome emulation", err)
		xsuper, ua = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiQ2hyb21lIiwiZGV2aWNlIjoiIiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiYnJvd3Nlcl91c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV2luNjQ7IHg2NCkgQXBwbGVXZWJLaXQvNTM3LjM2IChLSFRNTCwgbGlrZSBHZWNrbykgQ2hyb21lLzEwMy4wLjAuMCBTYWZhcmkvNTM3LjM2IiwiYnJvd3Nlcl92ZXJzaW9uIjoiMTAzLjAuMC4wIiwib3NfdmVyc2lvbiI6IjEwIiwicmVmZXJyZXIiOiIiLCJyZWZlcnJpbmdfZG9tYWluIjoiIiwicmVmZXJyZXJfY3VycmVudCI6Imh0dHBzOi8vZGlzY29yZC5jb20vIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiZGlzY29yZC5jb20iLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfYnVpbGRfbnVtYmVyIjoxMzY5MjEsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
		v =  "103"
	} else {
		utilities.LogSuccess("Successfully obtained build number, useragent and latest chrome version")
	}

	// Load instances
	tokens, err = utilities.ReadLines("tokens.txt")
	if err != nil {
		return cfg, instances, err
	}
	if len(tokens) == 0 {
		return cfg, instances, fmt.Errorf("no tokens found in tokens.txt")
	}

	if cfg.ProxySettings.ProxyFromFile {
		proxies, err = utilities.ReadLines("proxies.txt")
		if err != nil {
			return cfg, instances, err
		}
		if len(proxies) == 0 {
			return cfg, instances, fmt.Errorf("no proxies found in proxies.txt")
		}
	}
	var Gproxy string
	var instanceToken string
	var email string
	var password string
	reg := regexp.MustCompile(`(.+):(.+):(.+)`)
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "" {
			continue
		}
		if reg.MatchString(tokens[i]) {
			parts := strings.Split(tokens[i], ":")
			instanceToken = parts[2]
			email = parts[0]
			password = parts[1]
		} else {
			instanceToken = tokens[i]
		}
		if cfg.ProxySettings.ProxyFromFile {
			proxy = proxies[rand.Intn(len(proxies))]
			Gproxy = proxy
		} else {
			proxy = ""
		}
		client, err := InitClient(proxy, cfg)
		if err != nil {
			return cfg, instances, fmt.Errorf("couldn't initialize client: %v", err)
		}
		// proxy is put in struct only to be used by gateway. If proxy for gateway is disabled, it will be empty
		if !cfg.ProxySettings.GatewayProxy {
			Gproxy = ""
		}
		instances = append(instances, Instance{Client: client, Token: instanceToken, Proxy: proxy, Config: cfg, GatewayProxy: Gproxy, Email: email, Password: password, UserAgent: ua, XSuper: xsuper, Version: v})
	}
	if len(instances) == 0 {
		utilities.LogErr(" You may be using 0 tokens")
	}

	return cfg, instances, nil

}

func SetMessages(instances []Instance, messages []Message) error {
	var err error
	if len(messages) == 0 {
		messages, err = GetMessage()
		if err != nil {
			return err
		}
		if len(messages) == 0 {
			return fmt.Errorf("no messages found in messages.txt")
		}
		for i := 0; i < len(instances); i++ {
			instances[i].Messages = messages
		}
	} else {
		for i := 0; i < len(instances); i++ {
			instances[i].Messages = messages
		}
	}

	return nil
}

func (in *Instance) CensorToken() string {
	if len(in.Token) == 0 {
		return ""
	}
	if in.Config.OtherSettings.CensorToken {
		var censored string
		l := len(in.Token)
		uncensoredPart := int(2 * l / 3)
		for i := 0; i < l; i++ {
			if i < uncensoredPart {
				censored += string(in.Token[i])
			} else {
				censored += "*"
			}
		}
		return censored
	} else {
		return in.Token
	}

}

func (in *Instance) WriteInstanceToFile(path string) {
	var line string
	if in.Email != "" && in.Password != "" {
		line = fmt.Sprintf("%s:%s:%s", in.Email, in.Password, in.Token)
	} else {
		line = in.Token
	}
	_ = utilities.WriteLinesPath(path, line)
}

func GetUseragentXSuper(locale string) (string, string, error) {
	err := UpdateDiscordBuildInfo()
	if err != nil {
		fmt.Println(err)
	}
	buildNo := GetDiscordBuildNumber("stable")
	if buildNo == "" {
		utilities.LogErr("Couldn't get build number")
		buildNo = "136240"
	}
	req, err := http.NewRequest("GET", "https://pastebin.com/raw/pGgMQGiJ", nil)
	if err != nil {
		return "", "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	var Info PastebinResponse
	err = json.Unmarshal(body, &Info)
	if err != nil {
		return "", "", err
	}
	r := rand.Intn(len(Info))
	super, err := XSuper(locale, buildNo, Info[r].Xs)
	if err != nil {
		return "", "", err
	}
	return Info[r].Ua, super, nil

}

func XSuper(locale string, buildNumber string, Xsuper string) (string, error) {
	// Decode Xsuper from base 64
	decoded, err := base64.StdEncoding.DecodeString(Xsuper)
	if err != nil {
		return "", err
	}
	decodedstr := string(decoded)
	build, err := strconv.Atoi(buildNumber)
	if err != nil {
		return "", err
	}
	decodedstr = strings.Replace(decodedstr, "xyzabc", locale, -1)
	decodedstr = strings.Replace(decodedstr, "12345", strconv.Itoa(build), -1)
	return base64.StdEncoding.EncodeToString([]byte(decodedstr)), nil
}

func DolfiesXsuper() (string, string, string, error) {
	apiLink := "https://discord-user-api.cf/api/v1/properties/web"
	req, err := http.NewRequest("GET", apiLink, nil)
	if err != nil {
		return "", "","", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "","", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "","", err
	}
	var Info DolfiesResponse
	err = json.Unmarshal(body, &Info)
	if err != nil {
		return "", "","", err
	}
	if Info.ChromeUserAgent == "" || Info.ChromeVersion == "" || Info.ClientBuildNumber == 0 {
		return "", "","", fmt.Errorf("couldn't get xsuper from Dolfies")
	}
	r := regexp.MustCompile(`Chrome/(.+) Safari`)
	chromeVersion := r.FindStringSubmatch(Info.ChromeUserAgent)
	if len(chromeVersion) == 0 {
		return "", "","", fmt.Errorf("couldn't get xsuper from Dolfies")
	}
	xsuper := fmt.Sprintf(`{"os":"Windows","browser":"Chrome","device":"","system_locale":"en-US","browser_user_agent":"%s","browser_version":"%s","os_version":"10","referrer":"","referring_domain":"","referrer_current":"https://discord.com/","referring_domain_current":"discord.com","release_channel":"stable","client_build_number":%d,"client_event_source":null}`, Info.ChromeUserAgent, chromeVersion[1], Info.ClientBuildNumber)
	return base64.StdEncoding.EncodeToString([]byte(xsuper)), Info.ChromeUserAgent,strings.Split(chromeVersion[1], ".")[0], nil

}

type XSuperProperties struct {
	OS                     string `json:"os,omitempty"`
	Browser                string `json:"browser,omitempty"`
	Device                 string `json:"device,omitempty"`
	SystemLocale           string `json:"system_locale,omitempty"`
	BrowserVersion         string `json:"browser_version,omitempty"`
	BrowserUserAgent       string `json:"browser_user_agent,omitempty"`
	OSVersion              string `json:"os_version,omitempty"`
	Referrer               string `json:"referrer,omitempty"`
	ReferringDomain        string `json:"referring_domain,omitempty"`
	ReferrerCurrent        string `json:"referrer_current,omitempty"`
	ReferringDomainCurrent string `json:"referring_domain_current,omitempty"`
	ReleaseChannel         string `json:"release_channel,omitempty"`
	ClientBuildNumber      int    `json:"client_build_number,omitempty"`
}

type PastebinResponse []struct {
	Ua string `json:"ua,omitempty"`
	Xs string `json:"xs,omitempty"`
}

type DolfiesResponse struct {
	ChromeUserAgent   string `json:"chrome_user_agent,omitempty"`
	ChromeVersion     string `json:"chrome_version,omitempty"`
	ClientBuildNumber int    `json:"client_build_number,omitempty"`
}
