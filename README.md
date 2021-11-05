Download from [here](https://github.com/V4NSH4J/discord-mass-DM-GO/releases)

[Discord server](https://discord.gg/fxPJAGxP7z) (temporary) 

Donate BTC: bc1qfmk95sqtw6sw2xc3kyaemcnltwcr5cs2phg2gh
# discord-mass-DM-GO
 A tool written in GO to demonstrate how bad actors utilize requests to spam Discord Users and launch large unsolicited DM Advertisement Campaigns
 
## Overview 🔍
 This program is a multi-threaded Discord Direct Message Spammer. It has 2 modes - Single and Multi. In Single mode, multiple tokens send messages to One discord account they share a mutual server with. In Multi mode, multiple discord tokens send messages to multiple discord accounts scraped from a Public discord server. 

 ![Feature preview - Discord Mass DM GO](https://i.imgur.com/DH9qMsl.png)
 
## Star the Repo ⭐
Please star the repo, it really helps me out and allows me to contribute more.

## Disclaimer ⚠️
 The automation of User Discord accounts also known as self-bots is a violation of Discord Terms of Service & Community guidelines and will result in your account(s) being terminated. Discretion is adviced. I will not be responsible for your actions. Please do not use my programs for raiding/ Spamming/ Harassment/ Unsolicited Advertisement . This program was solely written to check a discord server's security measures and to document the relative ease with which bad actors function on Discord.

## How is this abused?
If you've been part of big discord servers, I'm sure you've at some point recieved a DM from one of such bots. Discord is a very large market of gamers with 150 million+ Monthly active users which is why this is such a big issue. People send Crypto exchange scams where they claim you won a fortune in a crypto currency and have to make an account on their website and make a deposit. Second type is Nitro Scams, where they either sent you a token logger binary or link you to a phishing website where they steal your credentials from either QR codes or login. After access of a user's account, their account is also used in a similar spam and their payment method is abused. Third people use to advertise their servers or their NFTs or their crypto to either Pump & dump or just make it popular 

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


## How to get Help?
You can make an [Issue](https://github.com/V4NSH4J/discord-mass-DM-GO/issues) Or join the temporary [discord server](https://discord.gg/XgdN6zsTKv) I made for this, although I'm not very active on discord. 

## Token Quality & Functionality (Updated: 2nd Nov)
The number of DMs each token of your's gets depends on it's quality. Here I will break down everything related to this. 
- Unverified Token : 5 DMs
- Email verified Token: 20+ DMs (Can be more or less depending on quality)
- Phone verified Tokens: 20+ DMs [Can be more or less depending on quality]


What happens when this limit is crossed? Unverified and Email verified tokens get phone locked (Meaning it requires a phone number to unlock them) And Phone verified tokens get disabled by discord for "Suspicious activity" and you need to reset their password to access them again. But for commerical purposes, tokens are one time use. 

About functionality, servers have a minimum verification level which server administrators can set. If the minimum server verification is set to none, then you can easily use Unverified tokens to DM it's members. But if it's set to email verified, your unverified tokens won't be able to DM anyone. Same goes with email verified tokens in Phone verification required servers. 

You do not need to do any of those crappy verifications (Click the check mark to continue, etc) to DM members. You don't even need to do verifications by bots like Alt Identifier, although they will kick your accounts in 10 minutes and you won't be able to DM anyone after that. Keep this in mind while using the program.
 
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
`individual_delay` | int | Duration in seconds between 2 consecutive messages from a single discord token
`rate_limit_delay` | int | Duration in seconds to wait when Discord rate limits sending DMs [Usually 600 for lesser individual delay]
`offset` | int | Duration in Miliseconds to displace the goroutines for better functionality
`skip_completed` | bool | Set to true to skip members who were already DM'd from completed.txt

## Other interesting stuff by me
[Discord Invite Joiner](https://github.com/V4NSH4J/discord-inviter-GO) - Joins given tokens to a server

[Discord Token Checker](https://github.com/V4NSH4J/FAST-discord-token-checker) - Checks given tokens and records their information

[Discord Mass DM](https://github.com/V4NSH4J/discord-mass-DM-GO) - DMs all users of a server or DM's a discord user from multiple accounts

[Dankgrinder](https://github.com/V4NSH4J/dankgrinder) - An Advanced automation tool for Dankmemer

## Donations 🪙
I spend quite a lot of time in making High Quality & Open Source discord tools because hundreds of people get ripped-off everyday searching for this stuff. If this helped you out even in the slightest, Buy me a coffee and make my day! 
BTC: bc1qfmk95sqtw6sw2xc3kyaemcnltwcr5cs2phg2gh
