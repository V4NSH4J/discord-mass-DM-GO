// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func LaunchGetMessage() {
	// Uses ?around & ?limit parameters to discord's REST API to get messages to get the exact message needed
	color.Cyan("Get Message - This will get the message from Discord which you want to send.")
	color.White("Enter your token: \n")
	var token string
	fmt.Scanln(&token)
	color.White("Enter the channelID: \n")
	var channelID string
	fmt.Scanln(&channelID)
	color.White("Enter the messageID: \n")
	var messageID string
	fmt.Scanln(&messageID)
	message, err := instance.FindMessage(channelID, messageID, token)
	if err != nil {
		color.Red("Error while finding message: %v", err)
		utilities.ExitSafely()
		return
	}
	color.Green("[%v] Message: %v", time.Now().Format("15:04:05"), message)
}
