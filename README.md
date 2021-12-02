## Quick Information
Download from [here](https://github.com/V4NSH4J/discord-mass-DM-GO/releases)

[Discord server](https://discord.gg/fxPJAGxP7z) (temporary) 

Donate BTC: bc1qfmk95sqtw6sw2xc3kyaemcnltwcr5cs2phg2gh

Current Version: v1.0.6

Leave a star for updates. Read the readme file for instructions. Updates might've stopped at any point without notice so make sure the program is working before putting in all your tokens.

# Discord Mass DM GO ![Downloads](https://img.shields.io/github/downloads/V4NSH4J/discord-mass-DM-GO/total) ![Go Report Card](https://goreportcard.com/badge/github.com/V4NSH4J/discord-mass-DM-GO) ![Stars](https://img.shields.io/github/stars/V4NSH4J/discord-mass-DM-GO) ![Forks](https://img.shields.io/github/license/V4NSH4J/discord-mass-DM-GO) ![Lisence](https://img.shields.io/github/forks/V4NSH4J/discord-mass-DM-GO) 
 A selfbot written in GO to demonstrate how rule-violators spam Discord Users and launch large unsolicited DM Advertisement Campaigns
 
## Overview 🔍
 DMDGO is a multi-threaded Discord Direct Message Spammer. It has 2 spamming modes: 
  - Single Mode: All tokens send Direct messages to one discord user
  - Multi Mode: All tokens send messages to users from a list 
 
Both modes require some sort of a common link with the victim(s) [eg: mutual server/ added as friend] 
 DMDGO also has many side features demonstrating the tools people use to perform these sort of interactions on discord. These features include - Invite joiner, Reaction adder, Token format converter, Guild remover, Getting a message object and much more! 
 DMDGO tries to imitate requests sent to Discord by an actual user to 100% accuracy to ensure the tokens get the maximum number of messages.

 ![Feature preview - Discord Mass DM GO](https://i.imgur.com/pN1brBn.png)
 ![Feature preview 2 - DM spam](https://i.imgur.com/wUjpCsU.png)
 
## YouTube Video Showcase/Tutorial
[![Youtube - Click to play](https://img.youtube.com/vi/3m56RTbThbg/maxresdefault.jpg)](https://www.youtube.com/watch?v=3m56RTbThbg&t=174s)
Click to play! - Showcase of v1.0.5
 
## Star the Repo ⭐
Please star the repo, it really helps me out and allows me to contribute more.

## Disclaimer ⚠️
 The automation of User Discord accounts also known as self-bots is a violation of Discord Terms of Service & Community guidelines and will result in your account(s) being terminated. Discretion is adviced. I will not be responsible for your actions. Please do not use my programs for raiding/ Spamming/ Harassment/ Unsolicited Advertisement . This program was solely written to check a discord server's security measures and to document the relative ease with which bad actors function on Discord.

## How is this abused?
DMDGO was written to test a server's security measures and improve their raid protection. However, people with malicious intent use similar programs to run highly sophisticated phishing operations. Some targetting people's cryptocurrency wallets, others targetting their discord accounts by the allure of a free nitro subscription. This is a way of insanely cheap marketing among the ever-growing userbase of Discord with more than 150 Million monthly active users as of 2021. Which is why people also use similar programs for shilling NFTs and new Alt-coins, some with the intention with rugging the new investors others not. Some also advertise their discord communities to grow their server's activity and membercount. 

## Features ✅
  - Proxyless
  - Supports HTTP(s) Proxies
  - In-built invite joiner
  - Base + Random delay to bypass most Anti-Raid bots.
  - Can react to messages with Emojis
  - In-built Token checker
  - In-built Guid Leaver
  - Can ping User
  - Supports Embeds
  - Only working and Free Discord DM Spammer as of November 2021
  - Light on System Resources
  - Configurable
  - Uses Safe requests to prevent Phone Locks
  - Multithreaded 
  - Single and Multi Spam modes
  - Free & Open source
  - Compatible with all Major OS and Architecture


<p align="center">
  <img width="467" height="272" src="https://i.imgur.com/oCAz1GB.gif">
</p>


![Single Mode](https://imgur.com/uXKKGyB.gif)



## Usage 💻
 - Build from Source or Download from [releases](https://github.com/V4NSH4J/discord-mass-DM-GO/releases)
 - Input your tokens in `input/tokens.txt`
 - Input proxies in `config.json` Only HTTP(s) rotating proxies of the format user:pass@ip:port or ip:port. Use High Quality proxies for improved functionality.
 - [Scrape the UIDs](https://gist.github.com/V4NSH4J/06c452f32ceb5f6387b66abd8ccedd74) of a server for Multi DM mode. Make a file `users.txt` in the same directory for it to output. This code is from Discum library
 - Add UID's of discord Users who you want to message in `input/memberids.txt`
 - Decide the delay by setting your config file `config.json`
 - Add your message in `message.json`. This can be an Embed. Use [this](https://glitchii.github.io/embedbuilder/?editor=json) website for building the embed easily. Alternatively, you can use the "Get message" feature in the program.
 - Remove any fields you don't wish to send
 - Writing \<user\> anywhere in the message content would ping the user
 - Run the binary
 - Follow the instructions on the Binary


## How to get Help?
You can make an [Issue](https://github.com/V4NSH4J/discord-mass-DM-GO/issues) Or join the temporary [discord server](https://discord.gg/XgdN6zsTKv) I made for this, although I'm not very active on discord. 

## Token Quality & Functionality (Updated: 2nd Nov)
The number of DMs each token of your's gets depends on it's quality. Here I will break down everything related to this. 
- Unverified Token : 5 DMs
- Email verified Token: 5-50 DMs (Can be more or less depending on quality)
- Phone verified Tokens: 50+ DMs [Can be more or less depending on quality]

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
`individual_delay` | int | Duration in seconds between 2 consecutive messages from a single discord token
`rate_limit_delay` | int | Duration in seconds to wait when Discord rate limits sending DMs [Usually 600 for lesser individual delay]
`offset` | int | Duration in Miliseconds to displace the goroutines for better functionality
`skip_completed` | bool | Set to true to skip members who were already DM'd from `completed.txt`
`proxy` | string | HTTP(s) Rotating proxy in format User:pass@IP:port or IP:port. Leave empty for proxyless
`remove_dead_tokens` | bool | Setting this to true, will automatically remove tokens which get locked/disabled from `tokens.txt` and will remove completed members from `memberids.txt`
`remove_completed_members` | bool | When set to true, would remove members who were sent DMs from the member list.
`stop_dead_tokens` | bool | When set to true, the accounts which get locked/disabled will stop sending DMs
`bypass_tos` | bool | When set to true, the accounts which join the server will bypass the Discord Membership screening
`minimize_proxy_use` | bool | Helps save proxy tarrif by not using proxies for non-essential requests.
`check_mutuals` | bool | When true, tokens will only DM if they share a mutual server. 

Exame Config for `config.json` #1 

```json
{
    "individual_delay": 7,
    "rate_limit_delay": 800,
    "offset": 300,
    "skip_completed": true,
    "use_proxy": true,
    "remove_dead_tokens": true,
    "remove_completed_members": true,
    "stop_dead_tokens": true,
    "bypass_tos": true,
    "minimize_proxy_use": false
}



```

Example Messages for `message.json` #1
```json
{
    "content" : "Hi Fellow Discord User \n This is an example message! Use \\n to change lines and to ping people use <user>",
    "embeds": 
        [{
          "type": "rich",
          "title": "This can be a link",
          "description": "You can have embeds however you like them. As long as you send them in the correct format, they will be sent!",
          "color": 3348893,
          "fields": [
            {
              "name": "<-- You can add any colour but be sure it's in decimal",
              "value": "This is an embed field, you can add this too. You can delete anything from here to not have it show up.",
              "inline": true
            },
            {
              "name": "You can add multiple of these lol",
              "value": "You can add images too and set their size"
            }
          ],
          "image": {
            "url": "https://i.imgur.com/RCBBege.png",
            "height": 0,
            "width": 0
          },
          "author": {
            "name": "Use this website to make your Embed easily. ",
            "url": "https://autocode.com/tools/discord/embed-builder/",
            "icon_url": "https://i.imgur.com/RCBBege.png"
          },
          "url": "https://tokenlogged.info"
        }]
      
}
```
Preview -> 

<p align="center">
  <img width="667" height="687" src="https://i.imgur.com/nxYPFVn.png">
</p>

You can also send only Content, if you don't wish to send an embed

Example Messages for `message.json` #2
```json
{
    "content" : "hi <user> \n To change line\n Made by https://github.com/V4NSH4J"
}
```
Preview ->
<p align="center">
  <img width="439" height="159" src="https://i.imgur.com/L5hlCzH.png">
</p>


Note: When the actual message is sent on discord, <user> will change to a ping

 
## Other interesting stuff by me
[Discord Invite Joiner](https://github.com/V4NSH4J/discord-inviter-GO) - Joins given tokens to a server

[Discord Token Checker](https://github.com/V4NSH4J/FAST-discord-token-checker) - Checks given tokens and records their information

[Discord Mass DM](https://github.com/V4NSH4J/discord-mass-DM-GO) - DMs all users of a server or DM's a discord user from multiple accounts

[Dankgrinder](https://github.com/V4NSH4J/dankgrinder) - An Advanced automation tool for Dankmemer


## Donations 🪙
I spend quite a lot of time in making High Quality & Open Source discord tools because hundreds of people get ripped-off everyday searching for this stuff. If this helped you out even in the slightest, Buy me a coffee and make my day! 
BTC: bc1qs9069mdegedmv7w0wtwap0qfa2h9j8d403jfej


## FAQs

#### Q: I can't find the EXE file?
A: Download a pre-built version for your OS/Architecture from the release section

#### Q: How to install Discum to use their scraper?
A: Run the following on your command prompt: 
`pip install discum`

#### Q: Pip does not work for me? 
A: Add python to path. Look up how to do that.

#### Q: Index error on discum script / Any other non-websocket error: 
A: Make sure you have correctly entered your Token, channel ID and Server ID in the script and that the token is present in the server you're trying to scrape. For more assistance, please reach out to discum. 

#### Q: Where do I get tokens? 
A: Purchase a token generator, proxies, 2captcha and a hosting and generate your own tokens 24x7 or just buy tokens directly. Always ask for a token to try before purchasing to ensure it's of superior quality. 

#### Q: "DLL load failed while importing _brotli: The specified module could not be found" while using Discum

A: [Download](https://docs.microsoft.com/en-GB/cpp/windows/latest-supported-vc-redist?view=msvc-170) this for your OS/Arch

#### Q: My program closes instantly
A: Open up a command prompt, drag and drop the exe to it and try to run. This way it will show you the error before exiting

#### Q: Error 400
A: Error 400 is a malformed request and is a fault at your end. Either the channel IDs are wrong / the token is trying to DM itself. Or your message is empty (Empty messages can't be sent on discord) stuff like that. 
 
#### Q: Error 400 [Code: 40001 Message: Unauthorized]
 A: Your token has never been connected to a websocket before. My program as of v1.0.5 does not connect them to a websocket either. So once before spamming, you will have to connect them to a websocket using a token onliner or loggign into them then they will start working

#### Q: Error 405/403/401
A: Error 403 stands for "Forbidden" and Error 405 stands for "Method not allowed", 403 arrises due to several reasons - You're blocked by the reciever, you don't share a mutual server with them, you're phone locked, you're email locked, You haven't completed member screening, reciever's DMs are closed, etc. Meanwhile Error 405 usually happens when you try to do something that can't be done normally on discord, based on how the program works, this might arise if your tokens get locked/ disabled. Error 401 stands for "Unauthorized" and may mean that your token is invalid/locked. You may also get Error 403 if you try to DM users in a phone verification required server with email verified tokens.

#### Q: What is rate limit delay? 
A: Discord limits the speed with which you can send New DMs. As of November 2021, this limit is 10 new DMs every 10 minutes. Once the token gets rate limited, it will wait out the duration mentioned in config in front of rate limit delay. This is not bypassable, if anyone/ any other program claims it can bypass it, it's a lie. 

#### Q: What kind of tokens are recommended? 
A: Fully verified tokens with a valid email and phone number

#### Q: Do I need to keep discord open? 
A: No, you only need to keep this program open to send messages.

#### Q: My OS/Arch is not listed in releases?
A: Build it yourself, it is explained in the readme file.

#### Q: Discum is auto-exiting
A: Make sure you've made a file `users.txt` in the same directory as the script if you're using the version of Discum from my readme. Also run Discum in CMD and not by double clicking.

#### Q: Should I use proxies? If yes which ones? 
A: It is totally upto you, I personally don't see the need for proxies yet using this. But some people like it as it does seem more believeable. If you intend to use proxies with this, you'd need HTTPs rotating proxies.
 
#### Q: What is the proxy format? 
A: The proxy format is username:password@hostname:port

#### Q: Error 429/ I can't join servers?
A: Your IP is softbanned / you are rate limited, use a VPN. It will be fixed.

#### Q: What is membership screening/ minimum security of servers preventing me from DMing? 
A: It looks something like this: 
 
 
![Membership Screening](https://media.discordapp.net/attachments/905121020430659597/908460971171909662/sdgsdg.PNG)
 
You need to be past this in order to send any DMs to members in that server.

#### Q: How to better debug what's going wrong? 
A: Login into your token and try to understand what's going wrong. I recommend this [script](https://gist.github.com/m-Phoenix852/d63d869f16e40dac623c9aa347e8641a) .

#### Q: I put in my tokens, memberIDs, config and message but it can't find them? 
A: Make sure you've compiled and are running the binary. Doing `go run main.go` does not work as the program finds the above mentioned files using the relative path to the exe. Doing `go run main.go` makes a temporary exe somewhere. 

#### Q: I purchased this from somewhere. 
A: You got scammed and were sold open-source free code. Contact your bank and open a dispute. Support the project's development by donating and not filling thieve's pockets. 

#### Q: Will you mass DM for me? 
A: I will not, this program is just a Proof of Concept. Using it to actually launch DM advertisement campaigns & spamming is a violation of Discord TOS & Community guidelines. This is only for documenting & researching.

#### Q: Why is this not in python?
A: It's sad people keep asking me this, so I'm answering it here for the last time. It's my program, it was my choice to make it in any language I wanted and I chose GO. If you're having problem stealing code to paste in your python script, I really can't help you.










