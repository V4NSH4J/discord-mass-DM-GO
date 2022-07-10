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
				MinVersion:    tls.VersionTLS12,
				// Renegotiation: tls.RenegotiateOnceAsClient,
				// SessionTicketsDisabled: false,
				CipherSuites: []uint16{
					tls.TLS_AES_128_GCM_SHA256,
					tls.TLS_AES_256_GCM_SHA384,
					tls.TLS_CHACHA20_POLY1305_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
				// NextProtos: []string{"h2" ,"http/1.1"},
				InsecureSkipVerify: true,
				CurvePreferences: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				},
			},
			DisableKeepAlives: cfg.OtherSettings.DisableKL,
			ForceAttemptHTTP2: true,
			Proxy:             http.ProxyURL(proxyURL),
		},
	}
	return Client, nil

}
