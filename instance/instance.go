// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
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
	color.Red("Websocket closed %v %v", err, in.Token)
	in.Receiver = false
	in.Ws, err = in.NewConnection(in.wsFatalHandler)
	if err != nil {
		in.fatal <- fmt.Errorf("failed to create websocket connection: %s", err)
		return
	}
	color.Green("Reconnected To Websocket")
}

func GetEverything() (Config, []Instance, error) {
	var cfg Config
	var instances []Instance
	var err error
	var tokens []string
	var proxies []string
	var proxy string

	// Load config
	cfg, err = GetConfig()
	if err != nil {
		return cfg, instances, err
	}
	supportedProtocols := []string{"http", "https", "socks4", "socks5"}
	if cfg.ProxySettings.ProxyProtocol != "" && !utilities.Contains(supportedProtocols, cfg.ProxySettings.ProxyProtocol) {
		color.Red("[!] You're using an unsupported proxy protocol. Assuming http by default")
		cfg.ProxySettings.ProxyProtocol = "http"
	}
	if cfg.ProxySettings.ProxyProtocol == "https" {
		cfg.ProxySettings.ProxyProtocol = "http"
	}
	if cfg.CaptchaSettings.CaptchaAPI == "" {
		color.Red("[!] You're not using a Captcha API, some functionality like invite joining might be unavailable")
	}
	if cfg.ProxySettings.Proxy != "" && os.Getenv("HTTPS_PROXY") == "" {
		os.Setenv("HTTPS_PROXY", cfg.ProxySettings.ProxyProtocol+"://"+cfg.ProxySettings.Proxy)
	}
	if !cfg.ProxySettings.ProxyFromFile && cfg.ProxySettings.ProxyForCaptcha {
		color.Red("[!] You must enabe proxy_from_file to use proxy_for_captcha")
		cfg.ProxySettings.ProxyForCaptcha = false
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
		instances = append(instances, Instance{Client: client, Token: instanceToken, Proxy: proxy, Config: cfg, GatewayProxy: Gproxy, Email: email, Password: password})
	}
	if len(instances) == 0 {
		color.Red("[!] You may be using 0 tokens")
	}
	var empty Config
	if cfg == empty {
		color.Red("[!] You may be using a malformed config.json")
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
