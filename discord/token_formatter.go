// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"strings"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func LaunchTokenFormatter() {
	color.Cyan("Email:Password:Token to Token")
	Tokens, err := utilities.ReadLines("tokens.txt")
	if err != nil {
		color.Red("Error while opening tokens.txt: %v", err)
		utilities.ExitSafely()
		return
	}
	if len(Tokens) == 0 {
		color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:04:05"))
		utilities.ExitSafely()
		return
	}
	var onlytokens []string
	for i := 0; i < len(Tokens); i++ {
		if strings.Contains(Tokens[i], ":") {
			token := strings.Split(Tokens[i], ":")[2]
			onlytokens = append(onlytokens, token)
		}
	}
	t := utilities.TruncateLines("tokens.txt", onlytokens)
	if t != nil {
		color.Red("[%v] Error while truncating tokens.txt: %v", time.Now().Format("15:04:05"), t)
		utilities.ExitSafely()
		return
	}
}
