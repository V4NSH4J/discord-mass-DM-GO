// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func InitClient(proxy string, cfg Config) (*http.Client, error) {
	// If proxy is empty, return a default client (if proxy from file is false)
	if proxy == "" {
		return http.DefaultClient, nil
	}
	switch cfg.ProxySettings.ProxyProtocol {
	case "http":
		if !strings.Contains(proxy, "http://") {
			proxy = "http://" + proxy
		}
	case "socks5":
		if !strings.Contains(proxy, "socks5://") {
			proxy = "socks5://" + proxy
		}
	case "socks4":
		if !strings.Contains(proxy, "socks4://") {
			proxy = "socks4://" + proxy
		}
	}
	// Error while converting proxy string to url.url would result in default client being returned
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return http.DefaultClient, err
	}
	// Creating a client and modifying the transport.

	Client := &http.Client{
		Timeout: time.Second * time.Duration(cfg.ProxySettings.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				CipherSuites:       []uint16{0x1301, 0x1303, 0x1302, 0xc02b, 0xc02f, 0xcca9, 0xcca8, 0xc02c, 0xc030, 0xc00a, 0xc009, 0xc013, 0xc014, 0x009c, 0x009d, 0x002f, 0x0035},
				InsecureSkipVerify: true,
				CurvePreferences:   []tls.CurveID{tls.CurveID(0x001d), tls.CurveID(0x0017), tls.CurveID(0x0018), tls.CurveID(0x0019), tls.CurveID(0x0100), tls.CurveID(0x0101)},
			},
			DisableKeepAlives: cfg.OtherSettings.DisableKL,
			ForceAttemptHTTP2: true,
			Proxy:             http.ProxyURL(proxyURL),
		},
	}
	return Client, nil

}
