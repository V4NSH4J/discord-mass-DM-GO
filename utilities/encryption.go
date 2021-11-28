// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/andybalholm/brotli"
)

// Decoding brotli encrypted responses
func DecodeBr(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	br := brotli.NewReader(r)
	return ioutil.ReadAll(br)
}

// Function to handle all sorts of accepted-encryptions
func ReadBody(resp http.Response) ([]byte, error) {

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.Header.Get("Content-Encoding") == "br" {
		bodybr, err := DecodeBr(body)
		if err != nil {
			return nil, err
		}
		return bodybr, nil
	}
	return body, nil

}
