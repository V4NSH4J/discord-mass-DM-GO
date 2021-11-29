// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"
)

type Config struct {
	Delay     int    `json:"individual_delay"`
	LongDelay int    `json:"rate_limit_delay"`
	Offset    int    `json:"offset"`
	Skip      bool   `json:"skip_completed"`
	Proxy     string `json:"proxy"`
	Call      bool   `json:"call"`
	Remove    bool   `json:"remove_dead_tokens"`
	RemoveM   bool   `json:"remove_completed_members"`
	Stop      bool   `json:"stop_dead_tokens"`
	Bypass    bool   `json:"bypass_tos"`
	Mutual    bool   `json:"check_mutual"`
}

func GetConfig() (Config, error) {
	var config Config
	ex, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		return Config{-1, -1, -1, false, "", false, false, false, false, false, false}, err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "config.json"))
	if err != nil {
		color.Red("Error while Opening config.json")
		return Config{-1, -1, -1, false, "", false, false, false, false, false, false}, err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	errr := json.Unmarshal(bytes, &config)
	if errr != nil {
		fmt.Println(err)
		return Config{-1, -1, -1, false, "", false, false, false, false, false, false}, err
	}

	return Config{config.Delay, config.LongDelay, config.Offset, config.Skip, config.Proxy, config.Call, config.Remove, config.RemoveM, config.Stop, config.Bypass, config.Mutual}, nil
}
