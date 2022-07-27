// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/discord"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/gookit/color"
)

func main() {
	version := "1.10.15"
	rand.Seed(time.Now().UTC().UnixNano())
	color.Blue.Printf(logo + " v" + version + "\n")
	color.Green.Printf("Made by https://github.com/V4NSH4J\nStar repository on github for updates!\n")
	utilities.VersionCheck(version)
	Options()
}

// Options menu
func Options() {
	utilities.PrintMenu([]string{"Invite Joiner", "Mass DM", "Single DM", "Reaction Adder", "Email:Password:Token to Token", "Token Checker", "Guild Leaver", "Token Onliner", "Scraping Menu", "Name Changer", "Avatar Changer", "Token Server Checker", "Bio Changer", "DM on Reaction", "Hypesquad Changer", "Token Password Changer", "Embed Maker", "Login into Token", "Token Nuker", "Button Presser", "Server Nickname Changer", "Friend Request Spammer", "Friends Mass DM [BETA]","Credits & Help", "Exit"})
	choice := utilities.UserInputInteger("Enter your choice!")
	switch choice {
	default:
		color.Red.Printf("Invalid choice!\n")
		Options()
	case 1:
		color.Cyan.Printf("Invite Joiner\n")
		discord.LaunchinviteJoiner()
	case 2:
		color.Cyan.Printf("Mass DM advertiser\n")
		discord.LaunchMassDM()
	case 3:
		color.Cyan.Printf("Single DM spam\n")
		discord.LaunchSingleDM()
	case 4:
		color.Cyan.Printf("Reaction Adder\n")
		discord.LaunchReactionAdder()
	case 5:
		color.Cyan.Printf("Email:Pass:Token to Token\n")
		discord.LaunchTokenFormatter()
	case 6:
		color.Cyan.Printf("Token Checker\n")
		discord.LaunchTokenChecker()
	case 7:
		color.Cyan.Printf("Guild Leaver\n")
		discord.LaunchGuildLeaver()
	case 8:
		color.Cyan.Printf("Token Onliner\n")
		discord.LaunchTokenOnliner()
	case 9:
		color.Cyan.Printf("Scraping Menu\n")
		discord.LaunchScraperMenu()
	case 10:
		color.Cyan.Printf("Name Changer\n")
		discord.LaunchNameChanger()
	case 11:
		color.Cyan.Printf("Profile Picture Changer\n")
		discord.LaunchAvatarChanger()
	case 12:
		color.Cyan.Printf("Token Servers Check\n")
		discord.LaunchServerChecker()
	case 13:
		color.Cyan.Printf("Bio Changer\n")
		discord.LaunchBioChanger()
	case 14:
		color.Cyan.Printf("DM on React\n")
		discord.LaunchDMReact()
	case 15:
		color.Cyan.Printf("Hypesquad Changer\n")
		discord.LaunchHypeSquadChanger()
	case 16:
		color.Cyan.Printf("Mass token changer\n")
		discord.LaunchTokenChanger()
	case 17:
		color.Cyan.Printf("Create Embed\n")
		discord.LanuchEmbed()
	case 18:
		color.Cyan.Printf("Login into Token\n")
		discord.LaunchTokenLogin()
	case 19:
		color.Cyan.Printf("Token Nuker\n")
		discord.LaunchTokenNuker()
	case 20:
		color.Cyan.Printf("Button Press\n")
		discord.LaunchButtonClicker()
	case 21:
		color.Cyan.Printf("Server Nickname Changer\n")
		discord.LaunchServerNicknameChanger()
	case 22:
		color.Cyan.Printf("Friend Request Spammer\n")
		discord.LaunchFriendRequestSpammer()
	case 23:
		color.Cyan.Printf("Friends Mass DM\n")
		discord.LaunchFriendSpammer()
	case 24:
		color.Blue.Printf("Made with <3 by github.com/V4NSH4J - Check out the github page for detailed documentation\n")
	case 25:
		os.Exit(0)
	}
	time.Sleep(1 * time.Second)
	Options()

}

const logo = "\r\n\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2557 \u2588\u2588\u2588\u2557   \u2588\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2588\u2588\u2557  \u2588\u2588\u2588\u2588\u2588\u2588\u2557  \u2588\u2588\u2588\u2588\u2588\u2588\u2557 \r\n\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2557 \u2588\u2588\u2588\u2588\u2551\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2550\u2550\u255D \u2588\u2588\u2554\u2550\u2550\u2550\u2588\u2588\u2557\r\n\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2554\u2588\u2588\u2588\u2588\u2554\u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2588\u2557\u2588\u2588\u2551   \u2588\u2588\u2551\r\n\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551\u255A\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551   \u2588\u2588\u2551\u2588\u2588\u2551   \u2588\u2588\u2551\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551 \u255A\u2550\u255D \u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\r\n\u255A\u2550\u2550\u2550\u2550\u2550\u255D \u255A\u2550\u255D     \u255A\u2550\u255D\u255A\u2550\u2550\u2550\u2550\u2550\u255D  \u255A\u2550\u2550\u2550\u2550\u2550\u255D  \u255A\u2550\u2550\u2550\u2550\u2550\u255D \r\nDISCORD MASS DM GO"
