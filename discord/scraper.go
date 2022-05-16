// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func LaunchScraperMenu() {
	cfg, _, err := instance.GetEverything()
	if err != nil {
		color.Red("Error while getting necessary data %v", err)
	}
	color.White("1) Online Scraper (Opcode 14)\n2) Scrape from Reactions\n3) Offline Scraper (Opcode 8)")
	var options int
	fmt.Scanln(&options)
	if options == 1 {
		var token string
		color.White("Enter token: ")
		fmt.Scanln(&token)
		var serverid string
		color.White("Enter serverid: ")
		fmt.Scanln(&serverid)
		var channelid string
		color.White("Enter channelid: ")
		fmt.Scanln(&channelid)
		Is := instance.Instance{Token: token}
		title := make(chan bool)
		go func() {
		Out:
			for {
				select {
				case <-title:
					break Out
				default:
					cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Scraped]`, len(Is.Ws.Members)))
					_ = cmd.Run()
				}

			}
		}()
		t := 0
		for {
			if t >= 5 {
				color.Red("[%v] Couldn't connect to websocket after retrying.", time.Now().Format("15:04:05"))
				break
			}
			err := Is.StartWS()
			if err != nil {
				color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
			} else {
				break
			}
			t++
		}

		color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), Is.Token)
		if Is.Ws != nil {
			err := instance.Subscribe(Is.Ws, serverid)
			if err != nil {
				color.Red("[%v][!] Error while subscribing to server: %s", time.Now().Format("15:04:05"), err)
			}
		}
		i := 0
		for {
			err := instance.Scrape(Is.Ws, serverid, channelid, i)
			if err != nil {
				color.Red("[%v] Error while scraping: %v", time.Now().Format("15:04:05"), err)
			}
			color.Green("[%v] Token %v Scrape Count: %v", time.Now().Format("15:04:05"), Is.Token, len(Is.Ws.Members))
			if Is.Ws.Complete {
				break
			}
			i++
			time.Sleep(time.Duration(cfg.ScraperSettings.SleepSc) * time.Millisecond)
		}
		if Is.Ws != nil {
			Is.Ws.Close()
		}
		color.Green("[%v] Scraping finished. Scraped %v members", time.Now().Format("15:04:05"), len(Is.Ws.Members))
		clean := utilities.RemoveDuplicateStr(Is.Ws.Members)
		color.Green("[%v] Removed Duplicates. Scraped %v members", time.Now().Format("15:04:05"), len(clean))
		color.Green("[%v] Write to memberids.txt? (y/n)", time.Now().Format("15:04:05"))
		title <- true
		var write string
		fmt.Scanln(&write)
		if write == "y" {
			for k := 0; k < len(clean); k++ {
				err := utilities.WriteLines("memberids.txt", clean[k])
				if err != nil {
					color.Red("[%v] Error while writing to memberids.txt: %v", time.Now().Format("15:04:05"), err)
				}
			}
			color.Green("[%v] Wrote to memberids.txt", time.Now().Format("15:04:05"))
			err := utilities.WriteFile("scraped/"+serverid+".txt", clean)
			if err != nil {
				color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
			}
		}

	}
	if options == 2 {
		var token string
		color.White("Enter token: ")
		fmt.Scanln(&token)
		var messageid string
		color.White("Enter messageid: ")
		fmt.Scanln(&messageid)
		var channelid string
		color.White("Enter channelid: ")
		fmt.Scanln(&channelid)
		color.White("1) Get Emoji from Message\n2) Enter Emoji manually")
		var option int
		var send string
		fmt.Scanln(&option)
		var emoji string
		if option == 2 {
			color.White("Enter emoji [For Native Discord Emojis, just copy and paste emoji as unicode. For Custom/Nitro Emojis enter Name:EmojiID exactly in this format]: ")
			fmt.Scanln(&emoji)
			send = emoji
		} else {
			msg, err := instance.GetRxn(channelid, messageid, token)
			if err != nil {
				fmt.Println(err)
			}
			color.White("Select Emoji")
			for i := 0; i < len(msg.Reactions); i++ {
				color.White("%v) %v %v", i, msg.Reactions[i].Emojis.Name, msg.Reactions[i].Count)
			}
			var index int
			fmt.Scanln(&index)
			if msg.Reactions[index].Emojis.ID == "" {
				send = msg.Reactions[index].Emojis.Name

			} else if msg.Reactions[index].Emojis.ID != "" {
				send = msg.Reactions[index].Emojis.Name + ":" + msg.Reactions[index].Emojis.ID
			}
		}

		var allUIDS []string
		var m string
		title := make(chan bool)
		go func() {
		Out:
			for {
				select {
				case <-title:
					break Out
				default:
					cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Scraped]`, len(allUIDS)))
					_ = cmd.Run()
				}

			}
		}()
		for {
			if len(allUIDS) == 0 {
				m = ""
			} else {
				m = allUIDS[len(allUIDS)-1]
			}
			rxn, err := instance.GetReactions(channelid, messageid, token, send, m)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(rxn) == 0 {
				break
			}
			fmt.Println(rxn)
			allUIDS = append(allUIDS, rxn...)

		}
		color.Green("[%v] Scraping finished. Scraped %v lines - Removing Duplicates", time.Now().Format("15:04:05"), len(allUIDS))
		clean := utilities.RemoveDuplicateStr(allUIDS)
		color.Green("[%v] Write to memberids.txt? (y/n)", time.Now().Format("15:04:05"))
		var write string
		fmt.Scanln(&write)
		if write == "y" {
			for k := 0; k < len(clean); k++ {
				err := utilities.WriteLines("memberids.txt", clean[k])
				if err != nil {
					color.Red("[%v] Error while writing to memberids.txt: %v", time.Now().Format("15:04:05"), err)
				}
			}
			color.Green("[%v] Wrote to memberids.txt", time.Now().Format("15:04:05"))
			err := utilities.WriteFile("scraped/"+messageid+".txt", allUIDS)
			if err != nil {
				color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
			}
		}
		title <- true
		fmt.Println("Done")
	}
	if options == 3 {
		// Query Brute. This is a test function. Try using the compressed stream to appear legit.
		// Make a list of possible characters - Space can only come once, double spaces are counted as single ones and Name can't start from space. Queries are NOT case-sensitive.
		// Start from a character, check the returns. If it's less than 100, that query is complete and no need to go further down the rabbit hole.
		// If it's more than 100 or 100 and the last name starts from the query, pick the letter after our query and go down the rabbit hole.
		// Wait 0.5s (Or better, needs testing) Between scrapes and systematically connect and disconnect from websocket to avoid rate limiting.
		// Global var where members get appended (even repeats, will be cleared later) list of queries completed, list of queries left to complete and last query the instance searched to be in struct
		// Scan line for user input to stop at any point and proceed with the memberids scraped at hand.
		// Multiple instance support. Division of queries and hence completes in lesser time.
		// Might not need to worry about spaces at all as @ uses no spaces.
		// Starting Websocket(s) Appending to a slice. 1 for now, add more later.
		color.Cyan("Opcode 8 Scraper (Offline Scraper)")
		color.White("This feature is intentionally slowed down with high delays. Please use multiple tokens and ensure they are in the server before starting to complete it quick.")
		cfg, instances, err := instance.GetEverything()
		if err != nil {
			color.Red("[%v] Error while getting config: %v", time.Now().Format("15:04:05"), err)
			utilities.ExitSafely()
		}
		var scraped []string
		var queriesCompleted []string
		// Input the number of tokens to be used
		title := make(chan bool)
		go func() {
		Out:
			for {
				select {
				case <-title:
					break Out
				default:
					cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Scraped %v Queries Completed]`, len(scraped), len(queriesCompleted)))
					_ = cmd.Run()
				}

			}
		}()
		color.Green("[%v] How many tokens do you wish to use? You have %v ", time.Now().Format("15:04:05"), len(instances))
		var numTokens int
		quit := make(chan bool)
		var allQueries []string
		fmt.Scanln(&numTokens)
		var chars string
		rawChars := " !\"#$%&'()*+,-./0123456789:;<=>?@[]^_`abcdefghijklmnopqrstuvwxyz{|}~" + cfg.ScraperSettings.ExtendedChars
		// Removing duplicates
		for i := 0; i < len(rawChars); i++ {
			if !strings.Contains(rawChars[0:i], string(rawChars[i])) {
				chars += string(rawChars[i])
			}
		}

		queriesLeft := make(chan string)

		for i := 0; i < len(chars); i++ {
			go func(i int) {
				queriesLeft <- string(chars[i])
			}(i)
		}

		if numTokens > len(instances) {
			color.Red("[%v] You only have %v tokens in your tokens.txt Using the maximum number of tokens possible", time.Now().Format("15:04:05"), len(instances))
		} else if numTokens <= 0 {
			color.Red("[%v] You must atleast use 1 token", time.Now().Format("15:04:05"))
			utilities.ExitSafely()
		} else if numTokens <= len(instances) {
			color.Green("[%v] You have %v tokens in your tokens.txt Using %v tokens", time.Now().Format("15:04:05"), len(instances), numTokens)
			instances = instances[:numTokens]
		} else {
			color.Red("[%v] Invalid input", time.Now().Format("15:04:05"))
		}

		color.Green("[%v] Enter the ServerID", time.Now().Format("15:04:05"))
		var serverid string
		fmt.Scanln(&serverid)
		color.Green("[%v] Press ENTER to START and STOP scraping", time.Now().Format("15:04:05"))
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		var namesScraped []string
		var avatarsScraped []string
		// Starting the instances as GOroutines
		for i := 0; i < len(instances); i++ {
			go func(i int) {
				instances[i].ScrapeCount = 0
				for {

					// Start websocket, reconnect if disconnected.
					if instances[i].ScrapeCount%5 == 0 || instances[i].LastCount%100 == 0 {
						if instances[i].Ws != nil {
							instances[i].Ws.Close()
						}
						time.Sleep(2 * time.Second)
						err := instances[i].StartWS()
						if err != nil {
							fmt.Println(err)
							continue
						}
						time.Sleep(2 * time.Second)

					}
					instances[i].ScrapeCount++

					// Get a query from the channel / Await for close response
					select {
					case <-quit:
						return
					default:
						query := <-queriesLeft
						allQueries = append(allQueries, query)
						if instances[i].Ws == nil {
							continue
						}
						if instances[i].Ws.Conn == nil {
							continue
						}
						err := instance.ScrapeOffline(instances[i].Ws, serverid, query)
						if err != nil {
							color.Red("[%v] %v Error while scraping: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), err)
							go func() {
								queriesLeft <- query
							}()
							continue
						}

						memInfo := <-instances[i].Ws.OfflineScrape
						queriesCompleted = append(queriesCompleted, query)
						var MemberInfo instance.Event
						err = json.Unmarshal(memInfo, &MemberInfo)
						if err != nil {
							color.Red("[%v] Error while unmarshalling: %v", time.Now().Format("15:04:05"), err)
							queriesLeft <- query
							continue
						}

						if len(MemberInfo.Data.Members) == 0 {
							instances[i].LastCount = -1
							continue
						}
						instances[i].LastCount = len(MemberInfo.Data.Members)
						for _, member := range MemberInfo.Data.Members {
							// Avoiding Duplicates
							if !utilities.Contains(scraped, member.User.ID) {
								scraped = append(scraped, member.User.ID)
							}
						}
						color.Green("[%v] Token %v Query %v Scraped %v [+%v]", time.Now().Format("15:04:05"), instances[i].CensorToken(), query, len(scraped), len(MemberInfo.Data.Members))

						for i := 0; i < len(MemberInfo.Data.Members); i++ {
							id := MemberInfo.Data.Members[i].User.ID
							err := utilities.WriteLines("memberids.txt", id)
							if err != nil {
								color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
								continue
							}
							if cfg.ScraperSettings.ScrapeUsernames {
								nom := MemberInfo.Data.Members[i].User.Username
								if !utilities.Contains(namesScraped, nom) {
									err := utilities.WriteLines("names.txt", nom)
									if err != nil {
										color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
										continue
									}
								}
							}
							if cfg.ScraperSettings.ScrapeAvatars {
								av := MemberInfo.Data.Members[i].User.Avatar
								if !utilities.Contains(avatarsScraped, av) {
									err := utilities.ProcessAvatar(av, id)
									if err != nil {
										color.Red("[%v] Error while processing avatar: %v", time.Now().Format("15:04:05"), err)
										continue
									}
								}
							}
						}
						if len(MemberInfo.Data.Members) < 100 {
							time.Sleep(time.Duration(cfg.ScraperSettings.SleepSc) * time.Millisecond)
							continue
						}
						lastName := MemberInfo.Data.Members[len(MemberInfo.Data.Members)-1].User.Username

						nextQueries := instance.FindNextQueries(query, lastName, queriesCompleted, chars)
						for i := 0; i < len(nextQueries); i++ {
							go func(i int) {
								queriesLeft <- nextQueries[i]
							}(i)
						}

					}

				}
			}(i)
		}

		bufio.NewReader(os.Stdin).ReadBytes('\n')
		color.Green("[%v] Stopping All instances", time.Now().Format("15:04:05"))
		title <- true
		for i := 0; i < len(instances); i++ {
			go func() {
				quit <- true
			}()
		}

		color.Green("[%v] Scraping Complete. %v members scraped.", time.Now().Format("15:04:05"), len(scraped))
		color.Green("Do you wish to write to file again? (y/n) [This will remove pre-existing IDs from memberids.txt]")
		var choice string
		fmt.Scanln(&choice)
		if choice == "y" || choice == "Y" {
			clean := utilities.RemoveDuplicateStr(scraped)
			err := utilities.TruncateLines("memberids.txt", clean)
			if err != nil {
				color.Red("[%v] Error while truncating file: %v", time.Now().Format("15:04:05"), err)
			}
			err = utilities.WriteFile("scraped/"+serverid, clean)
			if err != nil {
				color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
			}
		}

	}
}
