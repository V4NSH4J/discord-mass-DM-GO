# discord-mass-DM-GO
 A tool written in GO to demonstrate how bad actors utilize requests to spam Discord Users and launch large unsolicited DM Advertisement Campaigns
 
## Overview 🔍
 This program is a multi-threaded Discord Direct Message Spammer. It has 2 modes - Single and Multi. In Single mode, multiple tokens send messages to One discord account they share a mutual server with. In Multi mode, multiple discord tokens send messages to multiple discord accounts scraped from a Public discord server. 

 ![Feature preview - Discord Mass DM GO](https://i.imgur.com/DH9qMsl.png)
 
 
## Disclaimer ⚠️
 The automation of User Discord accounts also known as self-bots is a violation of Discord Terms of Service & Community guidelines and will result in your account(s) being terminated. Discretion is adviced. I will not be responsible for your actions. Please do not use my programs for raiding/ Spamming/ Harassment/ Unsolicited Advertisement . This program was solely written to check a discord server's security measures and to document the relative ease with which bad actors function on Discord.
 
## Features ✅
  - Proxyless
  - Only working and Free Discord DM Spammer as of November 2021
  - Light on System Resources
  - Configurable
  - Uses Safe requests to prevent Phone Locks
  - Multithreaded 
  - Single and Multi Spam modes
  - Free & Open source
  - Compatible with all Major OS and Architecture

![Mass DM in action](https://i.imgur.com/oCAz1GB.gif)


[Single DM in action](https://imgur.com/uXKKGyB.gif)


## Usage 💻
 - Build from Source or Download from [releases](https://github.com/V4NSH4J/discord-mass-DM-GO/releases)
 - Input your tokens in "input/tokens.txt"
 - [Scrape the UIDs](https://github.com/Merubokkusu/Discord-S.C.U.M/blob/master/examples/gettingGuildMembers.py) of a server for Multi DM mode. 
 - Add UID's of discord Users who you want to message in "input/memberids.txt"
 - Decide the delay and the message by setting your config file "config.json"
 - Run the binary
 - Follow the instructions on the Binary
 
## Building from Source 🚧
 - [Install Golang](https://golang.org) and verify your installation
 - Open up a terminal window 
 - Navigate to the directory of the source code
 - Type "go build" into your console and a Binary should pop up

## Configuration

Name | Type | Description
---- | ---- | ----
`mode` | int | Mode 0 for spamming a Single account. Mode 1 for Mass spamming Discord accounts
`message` | string | The message to be sent to the Discord User
`delay` | int | Duration in seconds between 2 consecutive messages from a single discord token

## Donations 🪙
I spend quite a lot of time in making High Quality & Open Source discord tools because hundreds of people get ripped-off everyday searching for this stuff. If this helped you out even in the slightest, Buy me a coffee and make my day! 
BTC: bc1qfmk95sqtw6sw2xc3kyaemcnltwcr5cs2phg2gh
