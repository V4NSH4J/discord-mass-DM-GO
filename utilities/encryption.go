// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"bytes"
	"fmt"
	"io/ioutil"

	http "github.com/Danny-Dasilva/fhttp"

	"compress/zlib"

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

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipreader, err := zlib.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		gzipbody, err := ioutil.ReadAll(gzipreader)
		if err != nil {
			return nil, err
		}
		return gzipbody, nil
	}

	if resp.Header.Get("Content-Encoding") == "br" {
		brreader := brotli.NewReader(bytes.NewReader(body))
		brbody, err := ioutil.ReadAll(brreader)
		if err != nil {
			fmt.Println(string(brbody))
			return nil, err
		}

		return brbody, nil
	}
	return body, nil
}
