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
	Locale          string
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

	locales := []string{"de-AT", "de-DE", "de-IT", "de-LI", "de-LU", "en-AG", "en-AI", "en-AT", "en-AU", "en-BB", "en-CA", "en-BS", "en-CH", "en-DE", "en-FI", "en-GB", "en-HK", "en-IN", "en-MY", "en-SG", "en-US", "fr-CA", "fr-FR"}
	locale := locales[rand.Intn(len(locales))]
	ua, xsuper, err = GetUseragentXSuper(locale)
	if err != nil {
		fmt.Println(err)
		ua , xsuper = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36", "eyJvcyI6Ik1hYyBPUyBYIiwiYnJvd3NlciI6IkNocm9tZSIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1HQiIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMDMuMC4wLjAgU2FmYXJpLzUzNy4zNiIsImJyb3dzZXJfdmVyc2lvbiI6IjEwMy4wLjAuMCIsIm9zX3ZlcnNpb24iOiIxMC4xNS43IiwicmVmZXJyZXIiOiJodHRwczovL3d3dy5nb29nbGUuY29tLyIsInJlZmVycmluZ19kb21haW4iOiJ3d3cuZ29vZ2xlLmNvbSIsInNlYXJjaF9lbmdpbmUiOiJnb29nbGUiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6OTk5OSwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbH0=" 
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
		instances = append(instances, Instance{Client: client, Token: instanceToken, Proxy: proxy, Config: cfg, GatewayProxy: Gproxy, Email: email, Password: password, UserAgent: ua, XSuper: xsuper, Locale: locale})
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
