// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crypto/tls"
)

func InitClient(proxy string, cfg Config) (*http.Client, error) {
	// If proxy is empty, return a default client (if proxy from file is false)
	unproxiedClient := &http.Client{
		Timeout: time.Second * time.Duration(cfg.ProxySettings.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion:         tls.VersionTLS13,
				CipherSuites:       cipherSuites(),
				InsecureSkipVerify: true,
				CurvePreferences:   curves(),
			},
			ForceAttemptHTTP2: true,
			DisableKeepAlives: cfg.OtherSettings.DisableKL,
		},
	}
	if proxy == "" {
		return unproxiedClient, nil
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
		return unproxiedClient, err
	}
	// Creating a client and modifying the transport.

	Client := &http.Client{
		Timeout: time.Second * time.Duration(cfg.ProxySettings.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS13,
				CipherSuites:       cipherSuites(),
				InsecureSkipVerify: true,
				CurvePreferences:   curves(),
			},
			ForceAttemptHTTP2: true,
			DisableKeepAlives: cfg.OtherSettings.DisableKL,
			Proxy:             http.ProxyURL(proxyURL),
		},
	}
	return Client, nil
}

func cipherSuites() []uint16 {
	var cipherSuites []uint16
	cipherSuites = append(cipherSuites, pickRandom([]uint16{0x1301, 0x1302, 0x1303}, 1)...)
	cipherSuites = append(cipherSuites, pickRandom([]uint16{0x0005, 0x000a, 0x002f, 0x0035, 0x003c, 0x009c, 0x009d, 0xc007, 0xc009, 0xc00a, 0xc011, 0xc012, 0xc013, 0xc014, 0xc023, 0xc027, 0xc02f, 0xc02b, 0xc030, 0xc02c, 0xcca8, 0xcca9}, 3)...)
	return cipherSuites
}

func curves() []tls.CurveID {
	x := pickRandom([]uint16{0x001d, 0x0017, 0x0018, 0x0019}, 1)
	var p []tls.CurveID
	for i := 0; i < len(x); i++ {
		p = append(p, tls.CurveID(x[i]))
	}
	return p
}

func pickRandom(array []uint16, minimum int) []uint16 {
	var results []uint16
	s := shuffle(array)
	var newArray []uint16
	for i := 0; i < len(array); i++ {
		if i < minimum {
			results = append(results, s[i])
		} else {
			newArray = append(newArray, s[i])
		}
	}
	if minimum == len(array) {
		return results
	}
	r := rand.Intn(len(newArray))
	for i := 0; i < r; i++ {
		results = append(results, newArray[i])
	}
	return results
}

func shuffle(array []uint16) []uint16 {
	for i := len(array) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
	return array
}
