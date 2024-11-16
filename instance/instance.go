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
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	gohttp "net/http"

	http "github.com/Danny-Dasilva/fhttp"
	"github.com/V4NSH4J/discord-mass-dm-GO/client"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/gorilla/websocket"
)

type Instance struct {
	Token           string
	Proxy           string
	Cookie          string
	Client          *http.Client
	UserAgent       string
	XSuper          string
	JA3             string
	ProxyProt       string
	Email           string
	Password        string
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
	Reacted         []ReactInfo
	ReactChannel    chan (ReactInfo)
	MessageNumber   int
	Version         string
	Quarantined     bool
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
	var Instances []Instance
	var err error
	var proxies []string
	var tokens []string
	var fingerprints []Fingerprints

	// Getting the config
	cfg, err = GetConfig()
	if err != nil {
		return cfg, Instances, fmt.Errorf(`error while getting config %s`, err)
	}
	if cfg.CaptchaSettings.CaptchaAPI == "invisifox.com" {
		utilities.LogWarn("You're using invisifox's downloadable solver. Make sure you're using port 8888 (default) and that you're running the application. Otherwise this wouldn't work.")
	}
	// Getting proxies & removing empty lines
	if cfg.ProxySettings.ProxyFromFile {
		proxies, err = utilities.ReadLines("proxies.txt")
		if err != nil {
			return cfg, Instances, fmt.Errorf(`error while getting proxies %s`, err)
		}
		var p []string
		for j := 0; j < len(proxies); j++ {
			if proxies[j] != "" {
				p = append(p, proxies[j])
			}
		}
		proxies = p
		if len(proxies) == 0 {
			cfg.ProxySettings.ProxyFromFile = false
			utilities.LogWarn("No proxies found in proxies.txt, disabling proxy from file")
		}
	}
	// Getting the tokens
	tokens, err = utilities.ReadLines("tokens.txt")
	if err != nil {
		return cfg, Instances, err
	}
	var v []string
	for j := 0; j < len(tokens); j++ {
		if tokens[j] != "" {
			v = append(v, tokens[j])
		}
	}
	tokens = v
	if len(tokens) == 0 {
		return cfg, Instances, fmt.Errorf("no tokens found in tokens.txt")
	}
	// Getting fingerprints - Prioritizing from file so multiple can be loaded.
	fingerprints, err = GetFingerprints()
	if err != nil {
		utilities.LogWarn("Error while getting fingerprints %s", err)
	}
	if len(fingerprints) == 0 {
		// Getting fingerprint from config if none are found in file
		if cfg.OtherSettings.JA3 != "" && cfg.OtherSettings.XSuperProperties != "" && cfg.OtherSettings.Useragent != "" {
			fingerprints = append(fingerprints, Fingerprints{
				JA3:              cfg.OtherSettings.JA3,
				XSuperProperties: cfg.OtherSettings.XSuperProperties,
				Useragent:        cfg.OtherSettings.Useragent,
			})
			// Getting fingerprint from Dolfies' API if none are set in config
		} else {
			xsuper, ua, _, err := DolfiesXsuper()
			if err != nil {
				// Hardcoding a fingerprint if Dolfies' API is down
				xsuper, ua = "eyJvcyI6Ik1hYyBPUyBYIiwiYnJvd3NlciI6IkZpcmVmb3giLCJkZXZpY2UiOiIiLCJzeXN0ZW1fbG9jYWxlIjoiZW4tVVMiLCJicm93c2VyX3VzZXJfYWdlbnQiOiJNb3ppbGxhLzUuMCAoTWFjaW50b3NoOyBJbnRlbCBNYWMgT1MgWCAxMC4xNTsgcnY6MTAzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvMTAzLjAiLCJicm93c2VyX3ZlcnNpb24iOiIxMDMuMCIsIm9zX3ZlcnNpb24iOiIxMC4xNSIsInJlZmVycmVyIjoiIiwicmVmZXJyaW5nX2RvbWFpbiI6IiIsInJlZmVycmVyX2N1cnJlbnQiOiIiLCJyZWZlcnJpbmdfZG9tYWluX2N1cnJlbnQiOiIiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfYnVpbGRfbnVtYmVyIjoxNDAwOTEsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:103.0) Gecko/20100101 Firefox/103.0"
			}
			fingerprints = append(fingerprints, Fingerprints{
				JA3:              "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0",
				XSuperProperties: xsuper,
				Useragent:        ua,
			})
		}
	}

	// Checking empty JA3s & setting according to same priority order as fingerprints as having a JA3 is important
	var ja3s []string
	for i := 0; i < len(fingerprints); i++ {
		if fingerprints[i].JA3 != "" {
			ja3s = append(ja3s, fingerprints[i].JA3)
		}
	}
	var j []string
	for i := 0; i < len(ja3s); i++ {
		if ja3s[i] != "" {
			j = append(j, ja3s[i])
		}
	}
	ja3s = j
	for i := 0; i < len(fingerprints); i++ {
		if fingerprints[i].JA3 == "" {
			if cfg.OtherSettings.JA3 != "" {
				fingerprints[i].JA3 = cfg.OtherSettings.JA3
			} else if len(ja3s) != 0 {
				fingerprints[i].JA3 = ja3s[rand.Intn(len(ja3s))]
			} else {
				fingerprints[i].JA3 = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0"
			}
		}
	}
	r := regexp.MustCompile(`(.+):(.+):(.+)`)
	for i := 0; i < len(tokens); i++ {
		var email, password, token string
		if r.MatchString(tokens[i]) {
			p := strings.Split(tokens[i], ":")
			email = p[0]
			password = p[1]
			token = p[2]
		} else {
			token = tokens[i]
		}
		var proxy, Gproxy, proxyProt string
		if cfg.ProxySettings.ProxyFromFile {
			proxy = proxies[rand.Intn(len(proxies))]
			Gproxy = proxy
			proxyProt = "http://" + proxy
		} else {
			proxy = ""
			proxyProt = ""
		}
		index := rand.Intn(len(fingerprints))
		httpclient, err := client.NewClient(client.Browser{JA3: fingerprints[index].JA3, UserAgent: fingerprints[index].Useragent, Cookies: nil}, cfg.ProxySettings.Timeout, false, fingerprints[index].Useragent, proxyProt)
		if err != nil {
			return cfg, Instances, err
		}
		// proxy is put in struct only to be used by gateway. If proxy for gateway is disabled, it will be empty
		if !cfg.ProxySettings.GatewayProxy {
			Gproxy = ""
		}
		Instances = append(Instances, Instance{Client: httpclient, Token: token, Proxy: proxyProt, Config: cfg, GatewayProxy: Gproxy, Email: email, Password: password, UserAgent: fingerprints[index].Useragent, XSuper: fingerprints[index].XSuperProperties, JA3: fingerprints[index].JA3, ProxyProt: proxyProt})
	}
	if len(Instances) == 0 {
		utilities.LogErr(" You may be using 0 tokens")
	}

	return cfg, Instances, nil

}

func OldGetEverything() (Config, []Instance, error) {
	var cfg Config
	var instances []Instance
	var err error
	var tokens []string
	var proxies []string
	var proxy string
	var xsuper string
	var ua string
	var v string
	var ja3 string
	var proxyProt string

	// Load config
	cfg, err = GetConfig()
	if err != nil {
		return cfg, instances, err
	}
	if !cfg.ProxySettings.ProxyFromFile && cfg.ProxySettings.ProxyForCaptcha {
		utilities.LogErr(" You must enable proxy_from_file to use proxy_for_captcha")
		cfg.ProxySettings.ProxyForCaptcha = false
	}
	if cfg.OtherSettings.XSuperProperties == "" && cfg.OtherSettings.Useragent == "" {
		xsuper, ua, v, err = DolfiesXsuper()
		if err != nil {
			utilities.LogErr(" Failed to get useragent and xsuper %s Turn off Dolfies mode or try again, otherwise program will be continued with hardcoded chrome emulation", err)
			xsuper, ua = "eyJvcyI6Ik1hYyBPUyBYIiwiYnJvd3NlciI6IkZpcmVmb3giLCJkZXZpY2UiOiIiLCJzeXN0ZW1fbG9jYWxlIjoiZW4tVVMiLCJicm93c2VyX3VzZXJfYWdlbnQiOiJNb3ppbGxhLzUuMCAoTWFjaW50b3NoOyBJbnRlbCBNYWMgT1MgWCAxMC4xNTsgcnY6MTAzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvMTAzLjAiLCJicm93c2VyX3ZlcnNpb24iOiIxMDMuMCIsIm9zX3ZlcnNpb24iOiIxMC4xNSIsInJlZmVycmVyIjoiIiwicmVmZXJyaW5nX2RvbWFpbiI6IiIsInJlZmVycmVyX2N1cnJlbnQiOiIiLCJyZWZlcnJpbmdfZG9tYWluX2N1cnJlbnQiOiIiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfYnVpbGRfbnVtYmVyIjoxNDAwOTEsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:103.0) Gecko/20100101 Firefox/103.0"
			v = "103"
		} else {
			utilities.LogSuccess("Successfully obtained build number, useragent and latest chrome version")
		}
	} else {
		xsuper, ua = cfg.OtherSettings.XSuperProperties, cfg.OtherSettings.Useragent
	}
	if cfg.OtherSettings.JA3 == "" {
		ja3 = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0"
	} else {
		ja3 = cfg.OtherSettings.JA3
	}
	ja32 := "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0"
	ja3s := []string{ja32, ja3}
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
			proxyProt = "http://" + proxy
		} else {
			proxy = ""
			proxyProt = ""
		}
		httpclient, err := client.NewClient(client.Browser{JA3: ja3, UserAgent: ua, Cookies: nil}, cfg.ProxySettings.Timeout, false, ua, proxyProt)
		if err != nil {
			return cfg, instances, err
		}
		// proxy is put in struct only to be used by gateway. If proxy for gateway is disabled, it will be empty
		if !cfg.ProxySettings.GatewayProxy {
			Gproxy = ""
		}
		instances = append(instances, Instance{Client: httpclient, Token: instanceToken, Proxy: proxyProt, Config: cfg, GatewayProxy: Gproxy, Email: email, Password: password, UserAgent: ua, XSuper: xsuper, Version: v, JA3: ja3s[rand.Intn(len(ja3s))], ProxyProt: proxyProt})
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

func DolfiesXsuper() (string, string, string, error) {
	apiLink := "https://discord-user-api.cf/api/v1/properties/web"
	req, err := gohttp.NewRequest("GET", apiLink, nil)
	if err != nil {
		return "", "", "", err
	}
	resp, err := gohttp.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}
	fmt.Println(string(body))
	var Info DolfiesResponse
	err = json.Unmarshal(body, &Info)
	if err != nil {
		return "", "", "", err
	}
	if Info.ChromeUserAgent == "" || Info.ChromeVersion == "" || Info.ClientBuildNumber == 0 {
		return "", "", "", fmt.Errorf("couldn't get xsuper from Dolfies")
	}
	r := regexp.MustCompile(`Chrome/(.+) Safari`)
	chromeVersion := r.FindStringSubmatch(Info.ChromeUserAgent)
	if len(chromeVersion) == 0 {
		return "", "", "", fmt.Errorf("couldn't get xsuper from Dolfies")
	}
	locales := []string{"af", "af-NA", "af-ZA", "agq", "agq-CM", "ak", "ak-GH", "am", "am-ET", "ar", "ar-001", "ar-AE", "ar-BH", "ar-DJ", "ar-DZ", "ar-EG", "ar-EH", "ar-ER", "ar-IL", "ar-IQ", "ar-JO", "ar-KM", "ar-KW", "ar-LB", "ar-LY", "ar-MA", "ar-MR", "ar-OM", "ar-PS", "ar-QA", "ar-SA", "ar-SD", "ar-SO", "ar-SS", "ar-SY", "ar-TD", "ar-TN", "ar-YE", "as", "as-IN", "asa", "asa-TZ", "ast", "ast-ES", "az", "az-Cyrl", "az-Cyrl-AZ", "az-Latn", "az-Latn-AZ", "bas", "bas-CM", "be", "be-BY", "bem", "bem-ZM", "bez", "bez-TZ", "bg", "bg-BG", "bm", "bm-ML", "bn", "bn-BD", "bn-IN", "bo", "bo-CN", "bo-IN", "br", "br-FR", "brx", "brx-IN", "bs", "bs-Cyrl", "bs-Cyrl-BA", "bs-Latn", "bs-Latn-BA", "ca", "ca-AD", "ca-ES", "ca-FR", "ca-IT", "ccp", "ccp-BD", "ccp-IN", "ce", "ce-RU", "cgg", "cgg-UG", "chr", "chr-US", "ckb", "ckb-IQ", "ckb-IR", "cs", "cs-CZ", "cy", "cy-GB", "da", "da-DK", "da-GL", "dav", "dav-KE", "de", "de-AT", "de-BE", "de-CH", "de-DE", "de-IT", "de-LI", "de-LU", "dje", "dje-NE", "dsb", "dsb-DE", "dua", "dua-CM", "dyo", "dyo-SN", "dz", "dz-BT", "ebu", "ebu-KE", "ee", "ee-GH", "ee-TG", "el", "el-CY", "el-GR", "en", "en-001", "en-150", "en-AG", "en-AI", "en-AS", "en-AT", "en-AU", "en-BB", "en-BE", "en-BI", "en-BM", "en-BS", "en-BW", "en-BZ", "en-CA", "en-CC", "en-CH", "en-CK", "en-CM", "en-CX", "en-CY", "en-DE", "en-DG", "en-DK", "en-DM", "en-ER", "en-FI", "en-FJ", "en-FK", "en-FM", "en-GB", "en-GD", "en-GG", "en-GH", "en-GI", "en-GM", "en-GU", "en-GY", "en-HK", "en-IE", "en-IL", "en-IM", "en-IN", "en-IO", "en-JE", "en-JM", "en-KE", "en-KI", "en-KN", "en-KY", "en-LC", "en-LR", "en-LS", "en-MG", "en-MH", "en-MO", "en-MP", "en-MS", "en-MT", "en-MU", "en-MW", "en-MY", "en-NA", "en-NF", "en-NG", "en-NL", "en-NR", "en-NU", "en-NZ", "en-PG", "en-PH", "en-PK", "en-PN", "en-PR", "en-PW", "en-RW", "en-SB", "en-SC", "en-SD", "en-SE", "en-SG", "en-SH", "en-SI", "en-SL", "en-SS", "en-SX", "en-SZ", "en-TC", "en-TK", "en-TO", "en-TT", "en-TV", "en-TZ", "en-UG", "en-UM", "en-US", "en-US-POSIX", "en-VC", "en-VG", "en-VI", "en-VU", "en-WS", "en-ZA", "en-ZM", "en-ZW", "eo", "es", "es-419", "es-AR", "es-BO", "es-BR", "es-BZ", "es-CL", "es-CO", "es-CR", "es-CU", "es-DO", "es-EA", "es-EC", "es-ES", "es-GQ", "es-GT", "es-HN", "es-IC", "es-MX", "es-NI", "es-PA", "es-PE", "es-PH", "es-PR", "es-PY", "es-SV", "es-US", "es-UY", "es-VE", "et", "et-EE", "eu", "eu-ES", "ewo", "ewo-CM", "fa", "fa-AF", "fa-IR", "ff", "ff-CM", "ff-GN", "ff-MR", "ff-SN", "fi", "fi-FI", "fil", "fil-PH", "fo", "fo-DK", "fo-FO", "fr", "fr-BE", "fr-BF", "fr-BI", "fr-BJ", "fr-BL", "fr-CA", "fr-CD", "fr-CF", "fr-CG", "fr-CH", "fr-CI", "fr-CM", "fr-DJ", "fr-DZ", "fr-FR", "fr-GA", "fr-GF", "fr-GN", "fr-GP", "fr-GQ", "fr-HT", "fr-KM", "fr-LU", "fr-MA", "fr-MC", "fr-MF", "fr-MG", "fr-ML", "fr-MQ", "fr-MR", "fr-MU", "fr-NC", "fr-NE", "fr-PF", "fr-PM", "fr-RE", "fr-RW", "fr-SC", "fr-SN", "fr-SY", "fr-TD", "fr-TG", "fr-TN", "fr-VU", "fr-WF", "fr-YT", "fur", "fur-IT", "fy", "fy-NL", "ga", "ga-IE", "gd", "gd-GB", "gl", "gl-ES", "gsw", "gsw-CH", "gsw-FR", "gsw-LI", "gu", "gu-IN", "guz", "guz-KE", "gv", "gv-IM", "ha", "ha-GH", "ha-NE", "ha-NG", "haw", "haw-US", "he", "he-IL", "hi", "hi-IN", "hr", "hr-BA", "hr-HR", "hsb", "hsb-DE", "hu", "hu-HU", "hy", "hy-AM", "id", "id-ID", "ig", "ig-NG", "ii", "ii-CN", "is", "is-IS", "it", "it-CH", "it-IT", "it-SM", "it-VA", "ja", "ja-JP", "jgo", "jgo-CM", "jmc", "jmc-TZ", "ka", "ka-GE", "kab", "kab-DZ", "kam", "kam-KE", "kde", "kde-TZ", "kea", "kea-CV", "khq", "khq-ML", "ki", "ki-KE", "kk", "kk-KZ", "kkj", "kkj-CM", "kl", "kl-GL", "kln", "kln-KE", "km", "km-KH", "kn", "kn-IN", "ko", "ko-KP", "ko-KR", "kok", "kok-IN", "ks", "ks-IN", "ksb", "ksb-TZ", "ksf", "ksf-CM", "ksh", "ksh-DE", "kw", "kw-GB", "ky", "ky-KG", "lag", "lag-TZ", "lb", "lb-LU", "lg", "lg-UG", "lkt", "lkt-US", "ln", "ln-AO", "ln-CD", "ln-CF", "ln-CG", "lo", "lo-LA", "lrc", "lrc-IQ", "lrc-IR", "lt", "lt-LT", "lu", "lu-CD", "luo", "luo-KE", "luy", "luy-KE", "lv", "lv-LV", "mas", "mas-KE", "mas-TZ", "mer", "mer-KE", "mfe", "mfe-MU", "mg", "mg-MG", "mgh", "mgh-MZ", "mgo", "mgo-CM", "mk", "mk-MK", "ml", "ml-IN", "mn", "mn-MN", "mr", "mr-IN", "ms", "ms-BN", "ms-MY", "ms-SG", "mt", "mt-MT", "mua", "mua-CM", "my", "my-MM", "mzn", "mzn-IR", "naq", "naq-NA", "nb", "nb-NO", "nb-SJ", "nd", "nd-ZW", "nds", "nds-DE", "nds-NL", "ne", "ne-IN", "ne-NP", "nl", "nl-AW", "nl-BE", "nl-BQ", "nl-CW", "nl-NL", "nl-SR", "nl-SX", "nmg", "nmg-CM", "nn", "nn-NO", "nnh", "nnh-CM", "nus", "nus-SS", "nyn", "nyn-UG", "om", "om-ET", "om-KE", "or", "or-IN", "os", "os-GE", "os-RU", "pa", "pa-Arab", "pa-Arab-PK", "pa-Guru", "pa-Guru-IN", "pl", "pl-PL", "ps", "ps-AF", "pt", "pt-AO", "pt-BR", "pt-CH", "pt-CV", "pt-GQ", "pt-GW", "pt-LU", "pt-MO", "pt-MZ", "pt-PT", "pt-ST", "pt-TL", "qu", "qu-BO", "qu-EC", "qu-PE", "rm", "rm-CH", "rn", "rn-BI", "ro", "ro-MD", "ro-RO", "rof", "rof-TZ", "ru", "ru-BY", "ru-KG", "ru-KZ", "ru-MD", "ru-RU", "ru-UA", "rw", "rw-RW", "rwk", "rwk-TZ", "sah", "sah-RU", "saq", "saq-KE", "sbp", "sbp-TZ", "se", "se-FI", "se-NO", "se-SE", "seh", "seh-MZ", "ses", "ses-ML", "sg", "sg-CF", "shi", "shi-Latn", "shi-Latn-MA", "shi-Tfng", "shi-Tfng-MA", "si", "si-LK", "sk", "sk-SK", "sl", "sl-SI", "smn", "smn-FI", "sn", "sn-ZW", "so", "so-DJ", "so-ET", "so-KE", "so-SO", "sq", "sq-AL", "sq-MK", "sq-XK", "sr", "sr-Cyrl", "sr-Cyrl-BA", "sr-Cyrl-ME", "sr-Cyrl-RS", "sr-Cyrl-XK", "sr-Latn", "sr-Latn-BA", "sr-Latn-ME", "sr-Latn-RS", "sr-Latn-XK", "sv", "sv-AX", "sv-FI", "sv-SE", "sw", "sw-CD", "sw-KE", "sw-TZ", "sw-UG", "ta", "ta-IN", "ta-LK", "ta-MY", "ta-SG", "te", "te-IN", "teo", "teo-KE", "teo-UG", "tg", "tg-TJ", "th", "th-TH", "ti", "ti-ER", "ti-ET", "to", "to-TO", "tr", "tr-CY", "tr-TR", "tt", "tt-RU", "twq", "twq-NE", "tzm", "tzm-MA", "ug", "ug-CN", "uk", "uk-UA", "ur", "ur-IN", "ur-PK", "uz", "uz-Arab", "uz-Arab-AF", "uz-Cyrl", "uz-Cyrl-UZ", "uz-Latn", "uz-Latn-UZ", "vai", "vai-Latn", "vai-Latn-LR", "vai-Vaii", "vai-Vaii-LR", "vi", "vi-VN", "vun", "vun-TZ", "wae", "wae-CH", "wo", "wo-SN", "xog", "xog-UG", "yav", "yav-CM", "yi", "yi-001", "yo", "yo-BJ", "yo-NG", "yue", "yue-Hans", "yue-Hans-CN", "yue-Hant", "yue-Hant-HK", "zgh", "zgh-MA", "zh", "zh-Hans", "zh-Hans-CN", "zh-Hans-HK", "zh-Hans-MO", "zh-Hans-SG", "zh-Hant", "zh-Hant-HK", "zh-Hant-MO", "zh-Hant-TW", "zu", "zu-ZA"}
	xsuper := fmt.Sprintf(`{"os":"Windows","browser":"Chrome","device":"","system_locale":"%s","browser_user_agent":"%s","browser_version":"%s","os_version":"10","referrer":"","referring_domain":"","referrer_current":"","referring_domain_current":"","release_channel":"stable","client_build_number":%d,"client_event_source":null}`, locales[rand.Intn(len(locales))], Info.ChromeUserAgent, chromeVersion[1], Info.ClientBuildNumber)
	return base64.StdEncoding.EncodeToString([]byte(xsuper)), Info.ChromeUserAgent, strings.Split(chromeVersion[1], ".")[0], nil

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

func GetFingerprints() ([]Fingerprints, error) {
	file, err := os.Open("fingerprints.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var fingerprints []Fingerprints
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &fingerprints)
	if err != nil {
		return nil, err
	}
	return fingerprints, nil
}
