// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/discord"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"

	"github.com/fatih/color"
)

var CaptchaServices []string

func main() {
	version := "1.9.1"
	CaptchaServices = []string{"capmonster.cloud", "2captcha.com", "rucaptcha.com", "anti-captcha.com"}
	rand.Seed(time.Now().UTC().UnixNano())
	color.Blue(logo + " v" + version + "\n")
	color.Green("Made by https://github.com/V4NSH4J\nStar repository on github for updates!")
	utilities.VersionCheck(version)
	Options()
}

// Options menu
func Options() {
	color.White("Menu:\n |- 01) Invite Joiner [Token]\n |- 02) Mass DM advertiser [Token]\n |- 03) Single DM spam [Token]\n |- 04) Reaction Adder [Token]\n |- 05) Get message [Input]\n |- 06) Email:Pass:Token to Token [Email:Password:Token]\n |- 07) Token Checker [Token]\n |- 08) Guild Leaver [Token]\n |- 09) Token Onliner [Token]\n |- 10) Scraping Menu [Input]\n |- 11) Name Changer [Email:Password:Token]\n |- 12) Profile Picture Changer [Token]\n |- 13) Token Servers Check [Token]\n |- 14) Bio Changer [Token]\n |- 15) DM on React\n |- 16) Hypesquad Changer\n |- 17) Mass token changer\n |- 18) Credits & Info\n |- 19) Exit")
	color.White("\nEnter your choice: ")
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	default:
		color.Red("Invalid choice!")
		Options()
	case 1:
		discord.LaunchinviteJoiner()
	case 2:
		discord.LaunchMassDM()
	case 3:
		discord.LaunchSingleDM()
	case 4:
		discord.LaunchReactionAdder()
	case 5:
		discord.LaunchGetMessage()
	case 6:
		discord.LaunchTokenFormatter()
	case 7:
		discord.LaunchTokenChecker()
	case 8:
		discord.LaunchGuildLeaver()
	case 9:
		discord.LaunchTokenOnliner()
	case 10:
		discord.LaunchScraperMenu()
	case 11:
		discord.LaunchNameChanger()
	case 12:
		discord.LaunchAvatarChanger()
	case 13:
		discord.LaunchServerChecker()
	case 14:
		discord.LaunchBioChanger()
	case 15:
		discord.LaunchDMReact()
	case 16:
		discord.LaunchHypeSquadChanger()
	case 17:
		discord.LaunchTokenChanger()

	case 18:
		color.Blue("Made with <3 by github.com/V4NSH4J - Check out the github page for detailed documentation")
	case 19:
		os.Exit(0)
	}
	time.Sleep(1 * time.Second)
	Options()

}

const logo = "\r\n\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2557 \u2588\u2588\u2588\u2557   \u2588\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2588\u2588\u2557  \u2588\u2588\u2588\u2588\u2588\u2588\u2557  \u2588\u2588\u2588\u2588\u2588\u2588\u2557 \r\n\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2557 \u2588\u2588\u2588\u2588\u2551\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2550\u2550\u255D \u2588\u2588\u2554\u2550\u2550\u2550\u2588\u2588\u2557\r\n\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2554\u2588\u2588\u2588\u2588\u2554\u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2588\u2557\u2588\u2588\u2551   \u2588\u2588\u2551\r\n\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551\u255A\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551   \u2588\u2588\u2551\u2588\u2588\u2551   \u2588\u2588\u2551\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551 \u255A\u2550\u255D \u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\r\n\u255A\u2550\u2550\u2550\u2550\u2550\u255D \u255A\u2550\u255D     \u255A\u2550\u255D\u255A\u2550\u2550\u2550\u2550\u2550\u255D  \u255A\u2550\u2550\u2550\u2550\u2550\u255D  \u255A\u2550\u2550\u2550\u2550\u2550\u255D \r\nDISCORD MASS DM GO"
