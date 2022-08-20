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
)

func LaunchScraperMenu() {
	cfg, _, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting neccessary information %v", err)
		utilities.ExitSafely()
	}
	utilities.PrintMenu([]string{"Online Scraper (Opcode 14)", "Scrape from Reactions (REST API)", "Offline Scraper (Opcode 8)"})
	options := utilities.UserInputInteger("Select an option: ")
	if options == 1 {
		token := utilities.UserInput("Enter the token: ")
		serverid := utilities.UserInput("Enter the server ID: ")
		channelid := utilities.UserInput("Enter the channel ID: ")
		var botsFile, avatarFile, nameFile, path, rolePath, scrapedFile, userDataFile string
		if cfg.OtherSettings.Logs {
			path = fmt.Sprintf(`logs/online_scraper/DMDGO-OS-%s-%s-%s-%s`, serverid, channelid, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
			rolePath = fmt.Sprintf(`%s/roles`, path)
			err := os.MkdirAll(path, 0755)
			if err != nil && !os.IsExist(err) {
				utilities.LogErr("Error creating logs directory: %s", err)
				utilities.ExitSafely()
			}
			err = os.MkdirAll(rolePath, 0755)
			if err != nil && !os.IsExist(err) {
				utilities.LogErr("Error creating roles directory: %s", err)
				utilities.ExitSafely()
			}
			botsFileX, err := os.Create(fmt.Sprintf(`%s/bots.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating bots file: %s", err)
				utilities.ExitSafely()
			}
			botsFileX.Close()
			AvatarFileX, err := os.Create(fmt.Sprintf(`%s/avatars.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating avatars file: %s", err)
				utilities.ExitSafely()
			}
			AvatarFileX.Close()
			NameFileX, err := os.Create(fmt.Sprintf(`%s/names.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating names file: %s", err)
				utilities.ExitSafely()
			}
			NameFileX.Close()
			ScrapedFileX, err := os.Create(fmt.Sprintf(`%s/scraped.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating scraped file: %s", err)
				utilities.ExitSafely()
			}
			UserDataFileX, err := os.Create(fmt.Sprintf(`%s/user_data.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating user data file: %s", err)
				utilities.ExitSafely()
			}
			botsFile, avatarFile, nameFile, scrapedFile, userDataFile = botsFileX.Name(), AvatarFileX.Name(), NameFileX.Name(), ScrapedFileX.Name(), UserDataFileX.Name()
		}
		Is := instance.Instance{Token: token}
		title := make(chan bool)
		go func() {
		Out:
			for {
				select {
				case <-title:
					break Out
				default:
					if Is.Ws != nil {
						if Is.Ws.Conn != nil {
							cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Scraped]`, len(Is.Ws.Members)))
							_ = cmd.Run()
						}
					}
				}
			}
		}()
		t := 0
		for {
			if t >= 5 {
				utilities.LogErr("Couldn't connect to websocket after retrying.")
				break
			}
			err := Is.StartWS()
			if err != nil {
				utilities.LogFailed("Error while opening websocket: %v", err)
			} else {
				break
			}
			t++
		}

		utilities.LogErr("Websocket opened %v", Is.Token)
		i := 0
		for {
			err := instance.Scrape(Is.Ws, serverid, channelid, i)
			if err != nil {
				utilities.LogErr("Error while scraping: %v", err)
			}
			utilities.LogSuccess("Token %v Scrape Count: %v", Is.Token, len(Is.Ws.Members))
			if Is.Ws.Complete {
				break
			}
			i++
			time.Sleep(time.Duration(cfg.ScraperSettings.SleepSc) * time.Millisecond)
		}
		if Is.Ws != nil {
			Is.Ws.Close()
		}
		utilities.LogSuccess("Scraping finished. Scraped %v members", len(Is.Ws.Members))
		if cfg.OtherSettings.Logs {
			for i := 0; i < len(Is.Ws.Members); i++ {
				if Is.Ws.Members[i].User.Bot {
					utilities.WriteLinesPath(botsFile, fmt.Sprintf("%v %v %v", Is.Ws.Members[i].User.ID, Is.Ws.Members[i].User.Username, Is.Ws.Members[i].User.Discriminator))
				}
				if Is.Ws.Members[i].User.Avatar != "" {
					utilities.WriteLinesPath(avatarFile, fmt.Sprintf("%v:%v", Is.Ws.Members[i].User.ID, Is.Ws.Members[i].User.Avatar))
				}
				if Is.Ws.Members[i].User.Username != "" {
					utilities.WriteLinesPath(nameFile, fmt.Sprintf("%v", Is.Ws.Members[i].User.Username))
				}
				for x := 0; x < len(Is.Ws.Members[i].Roles); x++ {
					utilities.WriteRoleFile(Is.Ws.Members[i].User.ID, rolePath, Is.Ws.Members[i].Roles[x])
				}
				if Is.Ws.Members[i].User.Discriminator != "" && Is.Ws.Members[i].User.Username != "" {
					utilities.WriteLinesPath(userDataFile, fmt.Sprintf("%v#%v", Is.Ws.Members[i].User.Username, Is.Ws.Members[i].User.Discriminator))
				}
				utilities.WriteLinesPath(scrapedFile, Is.Ws.Members[i].User.ID)
			}
		}
		var memberids []string
		for _, member := range Is.Ws.Members {
			memberids = append(memberids, member.User.ID)
		}
		clean := utilities.RemoveDuplicateStr(memberids)
		utilities.LogSuccess("Removed Duplicates. Scraped %v members", len(clean))
		write := utilities.UserInput("Write to memberids.txt? (y/n)")
		title <- true
		if write == "y" {
			for k := 0; k < len(clean); k++ {
				err := utilities.WriteLines("memberids.txt", clean[k])
				if err != nil {
					utilities.LogErr("Error while writing to file: %v", err)
				}
			}
			utilities.LogSuccess("Wrote %v members to memberids.txt", len(clean))
		}

	}
	if options == 2 {
		token := utilities.UserInput("Enter the token: ")
		channelid := utilities.UserInput("Enter the channel ID: ")
		messageid := utilities.UserInput("Enter the message ID: ")
		utilities.PrintMenu([]string{"Get Emoji from Message", "Enter Emoji Manually"})
		option := utilities.UserInputInteger("Select an option: ")
		var send string
		if option == 2 {
			send = utilities.UserInput("Enter the emoji [Format emojiName or emojiName:emojiID for nitro emojis]: ")
		} else {
			msg, err := instance.GetRxn(channelid, messageid, token)
			if err != nil {
				utilities.LogErr("Error while getting message: %v", err)
			}
			var selection []string
			for i := 0; i < len(msg.Reactions); i++ {
				selection = append(selection, fmt.Sprintf("Emoji: %v | Count: %v", msg.Reactions[i].Emojis.Name, msg.Reactions[i].Count))
			}
			utilities.PrintMenu2(selection)
			index := utilities.UserInputInteger("Select an option: ")
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
				utilities.LogErr("Error while getting reactions: %v", err)
				continue
			}
			if len(rxn) == 0 {
				break
			}
			utilities.LogInfo("Scraped %v members", len(rxn))
			allUIDS = append(allUIDS, rxn...)

		}
		utilities.LogInfo("Scraping finished. Scraped %v lines - Removing Duplicates", len(allUIDS))
		clean := utilities.RemoveDuplicateStr(allUIDS)
		path := fmt.Sprintf(`logs/reaction_scraper/DMDGO-RS-%s-%s-%s-%s`, channelid, messageid, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
		err := os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			utilities.LogErr("Error creating logs directory: %s", err)
			utilities.ExitSafely()
		}
		scrapedFileX, err := os.Create(fmt.Sprintf(`%s/scraped.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating scraped file: %s", err)
			utilities.ExitSafely()
		}
		defer scrapedFileX.Close()
		scrapedFile := scrapedFileX.Name()
		for i := 0; i < len(clean); i++ {
			utilities.WriteLinesPath(scrapedFile, clean[i])
		}
		write := utilities.UserInput("Write to memberids.txt? (y/n)")
		if write == "y" {
			for k := 0; k < len(clean); k++ {
				err := utilities.WriteLines("memberids.txt", clean[k])
				if err != nil {
					utilities.LogErr("Error while writing to file: %v", err)
				}
			}
			utilities.LogSuccess("Wrote %v members to memberids.txt", len(clean))
		}
		title <- true
		utilities.LogSuccess("Finished")
	}
	if options == 3 {
		cfg, instances, err := instance.GetEverything()
		if err != nil {
			utilities.LogErr("Error while getting instances: %v", err)
			return
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
		numTokens := utilities.UserInputInteger("How many tokens do you wish to use? You have %v ", len(instances))
		quit := make(chan bool)
		var allQueries []string
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
			utilities.LogWarn("You only have %v tokens in your tokens.txt Using the maximum number of tokens possible", len(instances))
		} else if numTokens <= 0 {
			utilities.LogErr("You must atleast use 1 token")
			utilities.ExitSafely()
		} else if numTokens <= len(instances) {
			utilities.LogInfo("You have %v tokens in your tokens.txt Using %v tokens", len(instances), numTokens)
			instances = instances[:numTokens]
		} else {
			utilities.LogErr("Invalid input")
		}

		serverid := utilities.UserInput("Enter the server ID: ")
		var tokenFile, botsFile, avatarFile, nameFile, path, rolePath, scrapedFile, userDataFile string
		if cfg.OtherSettings.Logs {
			path = fmt.Sprintf(`logs/offline_scraper/DMDGO-OS-%s-%s-%s`, serverid, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
			rolePath = fmt.Sprintf(`%s/roles`, path)
			err := os.MkdirAll(path, 0755)
			if err != nil && !os.IsExist(err) {
				utilities.LogErr("Error creating logs directory: %s", err)
				utilities.ExitSafely()
			}
			err = os.MkdirAll(rolePath, 0755)
			if err != nil && !os.IsExist(err) {
				utilities.LogErr("Error creating roles directory: %s", err)
				utilities.ExitSafely()
			}
			tokenFileX, err := os.Create(fmt.Sprintf(`%s/token.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating token file: %s", err)
				utilities.ExitSafely()
			}
			tokenFileX.Close()
			botsFileX, err := os.Create(fmt.Sprintf(`%s/bots.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating bots file: %s", err)
				utilities.ExitSafely()
			}
			botsFileX.Close()
			AvatarFileX, err := os.Create(fmt.Sprintf(`%s/avatars.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating avatars file: %s", err)
				utilities.ExitSafely()
			}
			AvatarFileX.Close()
			NameFileX, err := os.Create(fmt.Sprintf(`%s/names.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating names file: %s", err)
				utilities.ExitSafely()
			}
			NameFileX.Close()
			ScrapedFileX, err := os.Create(fmt.Sprintf(`%s/scraped.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating scraped file: %s", err)
				utilities.ExitSafely()
			}
			defer ScrapedFileX.Close()
			UserDataFileX, err := os.Create(fmt.Sprintf(`%s/userdata.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating userdata file: %s", err)
				utilities.ExitSafely()
			}
			defer UserDataFileX.Close()
			tokenFile, botsFile, avatarFile, nameFile, scrapedFile, userDataFile = tokenFileX.Name(), botsFileX.Name(), AvatarFileX.Name(), NameFileX.Name(), ScrapedFileX.Name(), UserDataFileX.Name()
			for i := 0; i < len(instances); i++ {
				instances[i].WriteInstanceToFile(tokenFile)
			}
		}
		utilities.LogInfo("Press ENTER to START and STOP scraping")
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
						if instances[i].Ws != nil {
							instances[i].Ws.Close()
						}
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
							utilities.LogErr("%v Error while scraping: %v", instances[i].CensorToken(), err)
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
							utilities.LogErr("Error while unmarshalling: %v", err)
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
						utilities.LogSuccess("Token %v Query %v Scraped %v [+%v]", instances[i].CensorToken(), query, len(scraped), len(MemberInfo.Data.Members))

						for i := 0; i < len(MemberInfo.Data.Members); i++ {
							id := MemberInfo.Data.Members[i].User.ID
							err := utilities.WriteLines("memberids.txt", id)
							if err != nil {
								utilities.LogErr("Error while writing to file: %v", err)
								continue
							}
							if cfg.OtherSettings.Logs {
								utilities.WriteLinesPath(scrapedFile, id)
								if MemberInfo.Data.Members[i].User.Bot {
									utilities.WriteLinesPath(botsFile, fmt.Sprintf("%v %v %v", id, MemberInfo.Data.Members[i].User.Username, MemberInfo.Data.Members[i].User.Discriminator))
								}
								if MemberInfo.Data.Members[i].User.Avatar != "" {
									utilities.WriteLinesPath(avatarFile, fmt.Sprintf("%v:%v", id, MemberInfo.Data.Members[i].User.Avatar))
								}
								if MemberInfo.Data.Members[i].User.Username != "" {
									utilities.WriteLinesPath(nameFile, fmt.Sprintf("%v", MemberInfo.Data.Members[i].User.Username))
								}
								for x := 0; x < len(MemberInfo.Data.Members[i].Roles); x++ {
									utilities.WriteRoleFile(id, rolePath, MemberInfo.Data.Members[i].Roles[x])
								}
								if MemberInfo.Data.Members[i].User.Username != "" && MemberInfo.Data.Members[i].User.Discriminator != "" {
									utilities.WriteLinesPath(userDataFile, fmt.Sprintf("%v#%v", MemberInfo.Data.Members[i].User.Username, MemberInfo.Data.Members[i].User.Discriminator))
								}
							}

							if cfg.ScraperSettings.ScrapeUsernames {
								nom := MemberInfo.Data.Members[i].User.Username
								if !utilities.Contains(namesScraped, nom) {
									err := utilities.WriteLines("names.txt", nom)
									if err != nil {
										utilities.LogErr("Error while writing to file: %v", err)
										continue
									}
								}
							}
							if cfg.ScraperSettings.ScrapeAvatars {
								av := MemberInfo.Data.Members[i].User.Avatar
								if !utilities.Contains(avatarsScraped, av) {
									err := utilities.ProcessAvatar(av, id)
									if err != nil {
										utilities.LogErr("Error while processing avatar: %v", err)
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
		utilities.LogInfo("Stopping All instances")
		title <- true
		for i := 0; i < len(instances); i++ {
			go func() {
				quit <- true
			}()
		}

		utilities.LogInfo("Scraping Complete. %v members scraped.", len(scraped))
		choice := utilities.UserInput("Do you wish to write to file again? (y/n) [This will remove pre-existing IDs from memberids.txt]")
		if choice == "y" || choice == "Y" {
			clean := utilities.RemoveDuplicateStr(scraped)
			err := utilities.TruncateLines("memberids.txt", clean)
			if err != nil {
				utilities.LogErr("Error while truncating file: %v", err)
			}
		}

	}
}
