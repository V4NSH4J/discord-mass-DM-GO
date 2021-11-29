// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Cookies are required for legitimate looking requests, a GET request to discord.com has these required cookies in it's response along with the website HTML
// We can use this to get the cookies & arrange them in a string
func Cookies() (string, error) {

	url := "https://discord.com"

	req, err := http.NewRequest("GET", url, nil)
	req.Close = true

	if err != nil {
		return "", err
	}

	httpClient := http.DefaultClient

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.Cookies() == nil {
		return "", fmt.Errorf("no cookies found")
	}
	var cookies string
	for _, cookie := range resp.Cookies() {
		cookies = cookies + cookie.Name + "=" + cookie.Value + "; "
	}

	return cookies + "locale=us", nil

}

type response struct {
	Fingerprint string `json:"fingerprint"`
}

// Getting Fingerprint to use in our requests for more legitimate seeming requests.
func Fingerprint() (string, error) {
	url := "https://discord.com/api/v9/experiments"

	req, err := http.NewRequest("GET", url, nil)
	req.Close = true

	if err != nil {
		return "", err
	}

	httpClient := http.DefaultClient

	resp, err := httpClient.Do(RegisterHeaders(req))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	p, m := DecodeBr(body)
	if m != nil {

		return "", m
	}

	var Response response

	err = json.Unmarshal(p, &Response)

	if err != nil {
		return "", err
	}

	return Response.Fingerprint, nil
}
