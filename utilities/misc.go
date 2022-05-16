// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func Snowflake() int64 {
	snowflake := strconv.FormatInt((time.Now().UTC().UnixNano()/1000000)-1420070400000, 2) + "0000000000000000000000"
	nonce, _ := strconv.ParseInt(snowflake, 2, 64)
	return nonce
}

func Contains(s []string, e string) bool {
	defer HandleOutOfBounds()
	if len(s) == 0 {
		return false
	}
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Inputs 2 slices of strings and returns a slice of strings which does not contain elements from the second slice
func RemoveSubset(s []string, r []string) []string {
	var n []string
	for _, v := range s {
		if !Contains(r, v) {
			n = append(n, v)
		}
	}
	return n
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func HandleOutOfBounds() {
	if r := recover(); r != nil {
		fmt.Printf("Recovered from Panic %v", r)
	}
}

func ExitSafely() {
	color.Red("Press ENTER to EXIT")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)
}

func VersionCheck(version string) {
	link := "https://pastebin.com/raw/CCaVBSPv"
	client := &http.Client{Timeout: time.Second * 15}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return 
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 204 {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	v := response["version"].(string)
	message := response["message"].(string)
	if v != version {
		color.Red("[!] You're using DMDGO V%v, but the latest version is V%v. Consider updating at https://github.com/V4NSH4J/discord-mass-DM-GO/releases", version, v)
	} else {
		color.Green("[O] You're Up-to-Date! You're using DMDGO V%v", version)
	}
	if message != "" {
		color.Yellow("[!] %v", message)
	}

	link = "https://pastebin.com/CCaVBSPv"
	req, err = http.NewRequest("GET", link, nil)
	if err != nil {
		return
	}
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 204 {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	r := regexp.MustCompile(`<div class="visits" title="Unique visits to this paste">\n(.+)<\/div>`)
	matches := r.FindStringSubmatch(string(body))
	if len(matches) == 0 {
		return
	}
	views := strings.ReplaceAll(matches[1], " ", "")
	color.Green("[O] DMDGO Users: %v [21-February-2022 - %v]", views, time.Now().Format("02-January-2006"))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
