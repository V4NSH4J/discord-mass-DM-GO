// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

func SetProxy(i int, j int) (*http.Client, error) {
	cfg, _ := GetConfig()
	if cfg.Proxy {
		proxies, err := ReadLines("proxy.txt")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		var proxyUrl *url.URL
		if i == -1 && j == -1 {
			proxyUrl, err = url.Parse("http://" + proxies[rand.Intn(len(proxies))])
			if err != nil {
				return nil, err
			}
		} else {
			if len(proxies) > j {
				proxies = proxies[:j]
			}
			tpp := len(proxies) / j
			proxyUrl, err = url.Parse("http://" + proxies[int(i*tpp)])
			if err != nil {
				return nil, err
			}
		}

		myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		return myClient, nil
	}
	return &http.Client{}, nil
}
