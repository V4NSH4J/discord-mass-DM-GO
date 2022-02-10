<p align="center">
  <img src="https://i.imgur.com/z8ig6eN.png">
  <img src="https://img.shields.io/github/license/V4NSH4J/discord-mass-DM-GO?style=for-the-badge&logo=appveyor">
  <img src="https://img.shields.io/github/downloads/V4NSH4J/discord-mass-DM-GO/total?style=for-the-badge&logo=appveyor">
  <img src="https://goreportcard.com/badge/github.com/V4NSH4J/discord-mass-dm-GO?style=for-the-badge&logo=appveyor">
  <img src="https://img.shields.io/github/stars/V4NSH4J/discord-mass-DM-GO?style=for-the-badge&logo=appveyor">
  <img src="https://img.shields.io/github/forks/V4NSH4J/discord-mass-DM-GO?style=for-the-badge&logo=appveyor">
  </p>

# Discord Mass DM GO
**DMDGO** is a Multi-threaded Discord Self-Bot primarily used for mass messaging users on Discord. It has numerous other quality features to enhance the user experience and allowing the user to target the most users. 

## Community 

[Telegram Server for Support](https://t.me/tosviolators)

[Discord Community Server](https://discord.gg/mtW89Y7prK)


## **Features** : 
- Invite joiner [Single Invite/Multiple Invites] with delays to bypass Anti-Raid bots
- Mass DM Advertisement 
- Single DM Spam
- Reaction Adder
- Token format changer
- Token Checker
- Guild Leaver
- Token Onliner
- Scraping [3 Modes to get the maximum users including an Opcode 8 Scraper]
- Name Changer
- Avatar Changer
- Checks if tokens are in a server
- Multi-threaded and supports high number of simultaneous accounts
- Proxyless / HTTP(s) Proxies
- Can ping users
- ~Can send Embeds~ (can no longer send embeds after Discord update on 22nd January 2022 removing ability to send embeds completely from regular discord users)
- Supports multiple messages, sends one randomly
- Can receive messages from people it messaged
- Uses safe requests to prevent phone locks
- Free & Open source
- Compatible with all major OS and Architectures
- Highly documented to help the user run it without problems


<p align="center">
  <img src="https://i.imgur.com/go9h9Fn.png">
</p>

## Disclaimer 
 The automation of User Discord accounts also known as self-bots is a violation of Discord Terms of Service & Community guidelines and will result in your account(s) being terminated. Discretion is adviced. I will not be responsible for your actions. Read about Discord [Terms Of service](https://discord.com/terms) and [Community Guidelines](https://discord.com/guidelines)
 
## Tutorial / Showcase Video
### DMDGO v1.0.7
[![Youtube - Click to play](https://i.imgur.com/Jx4gk54.png)](https://youtu.be/9HX64DHJYWI)
Click to play

### DMDGO v1.0.5 
[![Youtube - Click to play](https://img.youtube.com/vi/3m56RTbThbg/maxresdefault.jpg)](https://www.youtube.com/watch?v=3m56RTbThbg&t=174s)
Click to play

 
## Basic Usage
1) [Build from source](https://github.com/V4NSH4J/discord-mass-DM-GO#building-from-source-) or download a pre-built version for your OS & Arch from [releases](https://github.com/V4NSH4J/discord-mass-DM-GO/releases)
2) Run the program via the binary. 
3) Set your [Config](https://github.com/V4NSH4J/discord-mass-DM-GO#configuration) by modifying the `config.json` file. 
4) If you already have memberids to DM, put them in `\input\memberids.txt` or obtain them from [Scraping]()
5) Put HTTP(s) Proxies in `\input\proxies.txt` if you enabled proxies in config. The format is IP:Port or User:pass@IP:Port if your proxies have a user-pass authentication. 
6) Enter your message(s) in `message.json` file. You can use [this](https://glitchii.github.io/embedbuilder/?editor=json) website to easily make JSON objects for embeds. However, if you do not want to/ are unable to format the file properly, you will have an option to input a simple message before mass DM-ing. Writing "\<user\>" without quotes anywhere in message content will ping the user to whom you're sending a message. 

## Building from source
1) Download and install [Golang](https://go.dev/) and verify your installation
2) Open a terminal window/ command prompt in the directory of the source code and type `go build`
3) A binary compatible with your OS/Arch should be made. If there are some problems on MacOS/Linux with executing the binary as a program. You can run this command `chmod +x ./discord-mass-dm-GO` or go to properties -> permissions -> Allow executing file as program. 

## How to get help? 
Read this documentation, try using `Ctrl + F` to find what you're looking for. Watch the tutorial video on YouTube. Other than that, feel free to make an [issue](https://github.com/V4NSH4J/discord-mass-DM-GO/issues) or try asking on our [Guilded Support Server](https://guilded.gg/tos)

## Configuration
Name | Type | Description
---- | ---- | ----
`individual_delay` | int | Duration in seconds that one account waits in between of sends 2 messages. Due to how discord rate limits are structured, you can only send 10 new messages every 10 minutes. To avoid hitting the rate limit, the recommended duration is 60.
`rate_limit_delay` | int | Duration in seconds that the account sleeps for when rate limited. If you send 10 messages very quickly, this rate limit is of 600 seconds (10 minutes). With the recommended settings, you'd never hit it, so it's recommended to set at 60.
`offset` | int | Duration in Milliseconds (1/1000th of a second) that the program waits in between of starting 2 instances. Perhaps one of the most important settings which is why it has it's [own section](). Recommended offset is (60/number of tokens) * 1000 but it does not matter with a few tokens and can be set to any small value like 50
`skip_completed` | bool | Program will avoid sending DMs to IDs it has already messaged. When someone is messaged, his ID is added to `input/completed.txt`. Adding IDs to `input/completed.txt` will in a way blacklist them. 
`skip_failed` | bool | Program will avoid sending DMs to IDs it has already attempted to message but failed. When someone is attempted to be messaged but fails, his ID is added to `input/failed.txt`.
`proxy` | string | Not included in the default config, but can be added. It is there to ensure backward compatibility, to use DMDGO like it was in Versions 1.0.6 and before. The format is IP:Port or User:Pass@IP:Port for User-Pass Authenticated proxies. See `proxy_from_file` to use proxies the new way. To use this option, `proxy_from_file` must be false. Adding proxy here uses it for gateway events even if `use_proxy_for_gateway` may be disabled.
`call` | bool | After a succesful DM, DMDGO will attempt to call the user. This is unnecessary as it would make no difference. The calls don't "ring" if the users are not friended which they're not in most cases. So in essence, it would be a call without a ring. Recommended is false.
`remove_dead_tokens` | bool | After completion of Mass DM, DMDGO will attempt to remove non-working tokens from `input/tokens.txt`
`remove_completed_members` | bool | After completion of Mass DM, DMDGO will attempt to remove completed members from `input/memberids.txt` leaving behind only the failed/unattempted IDs so the user can re-run the program to target them
`stop_dead_tokens` | bool | Once a token is locked/disabled, it will stop DMing if this is set to true. Recommended to be true. 
`check_mutual` | bool | DMDGO will check if the user has a mutual guild with the token and also gather the user's Name and Discrim. This has been buggy in the past. Recommended to be set to false. 
`friend_before_DM` | bool | DMDGO will send a friend request to everyone in the list before attempting to DM them. Requires `online_tokens` and `check_mutual`. It needs `online_tokens` as a precautionary measure. Incase a token has never connected to gateway before and it tries sending a friend request, it will get locked.
`online_tokens` | bool | DMDGO will online tokens before starting mass DM. This is a very important setting, certain functions like `receive_messages`, `friend_before_DM` and `call` cannot be used without it. But if you're not planning to use those features anyways, it can be set to false.
`online_scraper_delay` | int | Time in milliseconds DMDGO will sleep between 2 consecutive scrapes to avoid being websocket-rate-limited. 1000 by default.
`proxy_from_file` | bool | If set to true, DMDGO will use proxies from `input\proxies.txt`
`max_dms_per_token` | int | Number of maximum DMs you want your tokens to send. If set to 0 they will send DMs till they are locked or the list is completed. 
`receive_messages` | bool | If set to true, the messages from other users will be logged on console and in `input\received.txt` 
`use_proxy_for_gateway` | bool | If set to true, websocket connections will use proxy as well. Recommend keeping it at false unless you have very good proxies.
`Timeout` | int | Timeout for all requests in seconds. Increase if slow connection or proxies.
`captcha_api` | string | Domain of the Captcha API you wish to use. The current supported Captcha APIs are capmonster.cloud and anti-captcha.com
`captcha_api_key` | string | Your Captcha API Key with balance loaded 
`max_attempt_invite_rejoin` | int | Maximum attempts to rejoin token to server. Use minimum of 2. Introduced since Discord added Captchas to join servers on some tokens. 
`disable_keep_alives` | bool | Closes the underlying TCP connection after every request. Might be helpful to change IPs on every request with rotating proxies. Only used with proxies_from_file and not with the proxy field in config. Added because older versions of DMDGO which used only rotating proxies did this by default.

### Offset 
Offset is a duration in milliseconds. As the name suggests this offsets or displaces the goroutines (threads) by a short period of time to ensure that all accounts don't start at the exact same second. What is the recommended offset? If you have less than 100 tokens or are using short individual delays, it does not matter. You can put any offset like 50-300. But if you are running a large number of tokens, you should set your individual and rate limit delays to 60 each or higher. Your offset will come with this formula - (individual delay/number of tokens) * 1000 This ensures your tokens start evenly spread out throughout the individual delay period. 
You can do more interesting things with offset. Normally to bypass Anti-Raid bots like Beemo or Wick, you'd have to join your tokens with high delays then wait for all of them to join to start DMing. Now with Offset you can make it so that one account joins and starts DMing, 30 seconds or any duration of your choice later the second account joins and start DMing so you save A LOT of time. How to do this? Set your offset to the duration you want your accounts to join in, like 30,000 - 60,000 (Remember offset is in milliseconds) and don't join your accounts to the server. Before Mass DMing, you'd get an option for advanced settings. Enter the server invite and serverid there. Use multiple proxies/ rotating proxies to prevent Discord server IP bans by the Anti-Raid bots. This won't work while Proxyless. 

## Using Captcha APIs
Captcha Solving APIs were introduced to DMDGO on 8th February 2022 when Discord mandated Captchas for joining servers on some tokens they deemed untrustworthy. The supported Captcha APIs right now are capmonster.cloud and anti-captcha.com 
You can register an account there, load some balance and copy your Captcha API Key to config. Make sure to specify the service you're using as well. It is extremely inexpensive and can join thousands of accounts in a couple USD. If there is an error with the captcha APIs, You will get an error code. You can look it up on their documentation [here](https://anti-captcha.com/apidoc/errors)

### Example configuration
```json
{
    "individual_delay": 60,
    "rate_limit_delay": 60,
    "offset": 100,
    "skip_completed": true,
    "skip_failed": true,
    "remove_dead_tokens": true,
    "remove_completed_members": true,
    "stop_dead_tokens": true,
    "check_mutual": false,
    "friend_before_DM": false,
    "online_tokens": false,
    "online_scraper_delay": 2000,
    "call": false,
    "proxy_from_file": false,
    "max_dms_per_token": 0,
    "receive_messages": true,
    "use_proxy_for_gateway": false,
    "timeout": 60,
    "captcha_api": "capmonster.cloud",
    "captcha_api_key": "your_captcha_api_key",
    "max_attempt_invite_rejoin": 4,
    "disable_keep_alives": false
}
```
This is the config I'd use, with ofcourse the offset calculated accordingly. 

## Message in file
The `input/message.json` is an array of messages from which one is chosen at random to be sent before each DM. Message.json is an array of messages. Find the examples below to add multiple messages. You can use the "get message" option to get messages from discord as well. Be sure to have the [] around the whole message. The only way to change lines is adding `\n`. After discord update on 22nd January 2022; Embed support was removed from DMDGO V1.0.7.5 and higher as discord removed the capibility to send embeds completely from userbots

### Example message 1 : Single Message, No Embed
```json
[
  {
    "content": "Hi <user> join my guilded server https://guilded.gg/tos"
  }
]
```

### Example message 2: Multiple messages, No Embeds. 
```json
[
  {
    "content": "Hi <user> join my guilded server https://guilded.gg/tos"
  },
  {
    "content": "We had a discord but it got terminated"
  },
  {
    "content": "We might make one again but too lazy to do so"
  }
]
```

## How to Debug problems with Mass DMing / Message 
For problems with setting the JSON files. Read this document very carefully and try understanding a bit of JSON. You can use [JSON Lint](https://jsonlint.com/) or similar to validate your JSON files and fix errors. The structure for all files is clearly defined here. 
For problems with sending DMs/ Any other function, the best way to diagnose is logging into the token and see what's going on. I highly recommend not to use email:password to login as it might trigger the New Login Location prompt. It's better to login via tokens. [Click here](https://gist.github.com/m-Phoenix852/d63d869f16e40dac623c9aa347e8641a) for a simple and fast token login script by @m-Pheonix852
Once logged into the token, you can see if it's still in the server / diagnose other problems with channel veriifcation levels, etc. You may also read the FAQs. Always keep in mind, when using self-bots you can only do actions that normal users can do. If you try DMing someone with no mutuals or friends, you can't DM them. Similarly, this self-bot can't either. 

## Scraping [Experimental Menu as of DMDGO v1.0.7]
The Scraping menu is a new functionality introduced in DMDGO V1.0.7 Before that, DMDGO recommended the use of [Discum's Scraper](https://gist.github.com/V4NSH4J/06c452f32ceb5f6387b66abd8ccedd74) 
This menu is still unstable and needs a lot of improvement. For stability, you can use the Discum Scraper. But if you do decide to use the scraper from DMDGO, It's explained here. 
- *Online Scraper (Opcode 14)* : Scrapes members from the member list visible on the right hand side of a discord server. This is usually only online members in case of larger servers as the offline member list usually gets hidden when servers are larger than 1000 members. 

- *Scrape from Reactions*: Does not use websocket, incase you see a reaction with a large number of reacts from which you'd like to get users, you can use this option. The one downside to this is that it would also scrape reacts from users who may have already left the server. 

- *Offline Scraper* (Opcode 8) : Scrapes members using OP8 websocket requests. This is what goes on behind the screens when you search for members in the search bar or by using @ in chat. This is usually slow with 1 account as it's bruteforcing which is why this mode supports multiple tokens for faster scrapes. It can't get all the users because of limitations described in [Discum's Docs](https://github.com/Merubokkusu/Discord-S.C.U.M/blob/master/docs/using/fetchingGuildMembers.md) but it gets a substantial number more users than the Opcode 14 scraper. You'll have to press ENTER to start and stop this scraper as it has the potential to go on for a very long time. Recommended to use multiple tokens and they need to be in the server before you start scraping. This function was slowed down on purpose to avoid any sort of rate limits. It will save IDs automatically to memberids.txt and you'll have to stop it manually when it's no longer getting IDs.
<p align="center">
  <img src="https://i.imgur.com/cMscRo5.png">
</p>


## Token Functionality [Updated: 14th January 2022]
The number of DMs you get per token depends on your token's quality and it's verification status. Unverified tokens (No email, No Phone number) get around 5 DMs. Email Verified tokens get around 50 DMs and Phone verifieds have the potential to go more than 50. 
These are ofcourse the *maximum/ideal* number of DMs as they were tested on botted servers where every DM was open. In a real-DMing scenerio, this number can be lesser. 
Whenever you buy tokens from a new seller, buy less tokens first to test out their quality. You can do this by manually sending DMs. Discord often flags the domains on which tokens are made, when this happens your tokens will struggle to even get 1 DM. Sometimes the phone number on the accounts is flagged as well. So always check before buying. 
Aged tokens have the potential to do thousands of DMs without ever getting disabled. But the type of tokens you want to use would depend on the servers you want to target and the cost efficiency. The mean DM price from DM-services is $0.01. The cost of an email verified token is $0.03 on average right now. Which would make the cost of DMs if you use email verifieds around $0.0006 - $ 0.001 which is more than 10 times cheaper. 
If you're new to this and want to try out your hand, I recommend going for cheap email verifieds to test. Other than that, it is upto you to narrow your targets and find your cost efficiency. 

## Mass DMing
Before you start Mass DMing, You will have an option for Advanced Settings. You may set a Serverid and invite code there. If the token is not in the server, the token will stop sending messages or try to rejoin the server. 
For best efficiency, use the recommended setting for delays. 
If you get spammed with errors "Cannot send messages to this user", make sure you're using the right memberids and that your tokens are in the server. 
If you get the error "Channel verification too high", this could be because you're trying to use email verified tokens on a server which requires phone verification or because the server has a 10 minute timer which you'd need to wait out before sending messages (Can be confirmed by loggin in with a token)
Sometimes, servers have anti-raid bots which detect suspicious patterns in joining like a lot of accounts with similar recent dates of registeration, no profile picture and random names joining within a certain time period. They may kick/ban the accounts, in such an event, you will not be able to send messages. Check out the method described [here]() or use high delays while joining such servers

## Proxies, Tokens and the Discord Self-Bot market in general
DMDGO was tested using Proxiware's Static Proxies and Iproyal's Rotating proxies. It may or may not work properly with free proxies from proxiscrape. Those are the worst proxies you can find on the internet. Using a proxies with gateway functions is not recommended.
Tokens are Discord accounts, they will be sending DMs for you. There are few ways to get them, the simplest being to buy them. Whenever you buy tokens, please check the quality and only buy more if they're good. Or you can buy/make/find a token generator. 
The Discord-Self bot market is very risky. A word of advice, don't purchase from unreputed people and use middlemen on reputed forums for large transactions. You will get scammed most of the times otherwise. The market is full of highly elaborate scammers like [Exordium](https://www.youtube.com/watch?v=uw7wjBxNK-4&ab_channel=Exordium) targetting people with his purchased channel and botted impressions. He will take your money and block you. And the owners of Anonix who will hapilly sell you open source code. 
Exit scams happen here all the time. Take recent incident of one of the MassDN partners Certex who exit scammed $60,000+ by ratting their customers. This is not to scare you to make purchases, this is just to warn you that you are likely to get scammed especially as a newcomer so stay vigilant. I decided to include this in the readme because everyday I see several people getting scammed.

## Support my Journey!
Leave a star on the repository, helps out intensively! You can also buy me a cookie on these addresses if I helped you out in any way. DMDGO was made with <3 over a period of 3 months and 184 cans of Redbull which doesn't come cheap :)

- *ETH*: 0xE01118C55963fA92174802Dae87E1C6DE1dADC07

- *BTC*: bc1qs9069mdegedmv7w0wtwap0qfa2h9j8d403jfej

- *SOL*: 8QyA9dCetgVMxU2AjzfM3DrY1i3mXuE8nsgLkvAX1hTe

- *LTC*: LN5UPbL31TcPzpBKFbsNKZ5BxwUzKcyi1F

## Credits
DMDGO has not been a One man show! I would like to thank everyone for their contributions and my patreons. Special thanks to my lads -> 
- [The author of Dankgrinder](https://github.com/dankgrinder) and [OsOmE1](https://github.com/OsOmE1) for helping out numerous times since the time I started writing code! Also for the websocket code taken from [dankgrinder](https://github.com/dankgrinder/dankgrinder)
- [Sympthey](https://github.com/Zenoryne) for allowing use of his Websocket code from [DiscSpam](https://github.com/Zenoryne/DiscSpam) helping me to understand the protocol!
- Arandomnewaccount and [Dolfies](https://github.com/dolfies) Contributors on Python Libraries for Discord Selfbots like [Discum](https://github.com/Merubokkusu/Discord-S.C.U.M) and [Discord.py-Self](https://github.com/dolfies/discord.py-self) for helping with the op8 and op14 scrapers and their amazing [docs](https://arandomnewaccount.gitlab.io/discord-unofficial-docs/lazy_guilds.html)
- Woen for providing the configuration for the initial HTTP client which did not really solve the problem but gave me the right direction!
- My friend [Siegfried](https://github.com/dasbard) for helping out with the community servers and many functions of DMDGO since even before it existed.


## FAQs

#### Q: I can't find the EXE file?
A: Download a pre-built version for your OS/Architecture from the release section or build from source.

#### Q: How to install Discum to use their scraper?
A: Run the following on your command prompt: 
`pip install discum`

#### Q: Pip does not work for me? 
A: Add python to path, watch a tutorial on it. 

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
 A: Your token has never been connected to a websocket before. Use token onliner in DMDGO v1.0.7 and above.

#### Q: Error 405/403/401
A: Error 403 stands for "Forbidden" and Error 405 stands for "Method not allowed", 403 arrises due to several reasons - You're blocked by the Receiver, you don't share a mutual server with them, you're phone locked, you're email locked, You haven't completed member screening, Receiver's DMs are closed, etc. Meanwhile Error 405 usually happens when you try to do something that can't be done normally on discord, based on how the program works, this might arise if your tokens get locked/ disabled. Error 401 stands for "Unauthorized" and may mean that your token is invalid/locked. You may also get Error 403 if you try to DM users in a phone verification required server with email verified tokens.

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

#### Q: Will you mass DM for me? 
A: I will not, this program is just a Proof of Concept. Using it to actually launch DM advertisement campaigns & spamming is a violation of Discord TOS & Community guidelines. This is only for documenting & researching.

#### Q: Channel verification too high?
A: This happens in a few scenerios. You're trying to use unverified tokens to DM in a server which needs Phone/Email verification OR you're trying to use email verified tokens in a server which requires phone verification. This may also happen if the server has a 10 minute wait time before you can interact in it, to verify, login into a token and see. It may also happen if for some reason, DMDGO failed to bypass the token which it does automatically.

#### Q: Invalid character `e` looking for beginning of value error code: 1015 
A: Cloudflare Error 1015 is an IP Based Rate limit. You have to use proxies/ VPN to get around it


