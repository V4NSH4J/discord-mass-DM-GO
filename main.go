// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

func main() {
	// Credits
	color.Blue(logo)
	color.Green("Made by https://github.com/V4NSH4J\nStar repository on github for updates!")
	Options()
}

// Options menu
func Options() {
	reg := regexp.MustCompile(`(.+):(.+):(.+)`)
	color.White("Menu:\n |- 01) Invite Joiner [Token]\n |- 02) Mass DM advertiser [Token]\n |- 03) Single DM spam [Token]\n |- 04) Reaction Adder [Token]\n |- 05) Get message [Input]\n |- 06) Email:Pass:Token to Token [Email:Password:Token]\n |- 07) Token Checker [Token]\n |- 08) Guild Leaver [Token]\n |- 09) Token Onliner [Token]\n |- 10) Scraping Menu [Input]\n |- 11) Name Changer [Email:Password:Token]\n |- 12) Profile Picture Changer [Token]\n |- 13) Token Servers Check [Token]\n |- 14) Credits & Info\n |- 15) Exit")
	color.White("\nEnter your choice: ")
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	default:
		color.Red("Invalid choice!")
		Options()
	case 0:
		color.Cyan("Debug Mode")
	case 1:
		var invitechoice int
		color.White("Invite Menu:\n1) Single Invite\n2) Multiple Invites from file")
		fmt.Scanln(&invitechoice)
		if invitechoice != 1 && invitechoice != 2 {
			color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
			ExitSafely()
			return
		}
		switch invitechoice {
		case 1:
			color.Cyan("Single Invite Mode")
			color.White("This will join your tokens from tokens.txt to a server")
			cfg, instances, err := getEverything()
			if err != nil {
				color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
			}
			color.White("[%v] Enter your invite CODE (The part after discord.gg/): ", time.Now().Format("15:04:05"))
			var invite string
			fmt.Scanln(&invite)
			color.White("[%v] Enter number of Threads (0: Unlimited Threads. 1: For using proper delay): ", time.Now().Format("15:04:05"))
			var threads int
			fmt.Scanln(&threads)

			if threads > len(instances) {
				threads = len(instances)
			}
			if threads == 0 {
				threads = len(instances)
			}

			color.White("[%v] Enter base delay for joining in seconds (0 for none)", time.Now().Format("15:04:05"))
			var base int
			fmt.Scanln(&base)
			color.White("[%v] Enter random delay to be added upon base delay (0 for none)", time.Now().Format("15:04:05"))
			var random int
			fmt.Scanln(&random)
			var delay int
			if random > 0 {
				delay = base + rand.Intn(random)
			} else {
				delay = base
			}
			c := goccm.New(threads)
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				c.Wait()
				go func(i int) {
					err := instances[i].Invite(invite)
					if err != nil {
						color.Red("[%v] Error while joining: %v", time.Now().Format("15:04:05"), err)
					}
					time.Sleep(time.Duration(delay) * time.Second)
					c.Done()

				}(i)
			}
			c.WaitAllDone()
			color.Green("[%v] All threads finished", time.Now().Format("15:04:05"))

		case 2:
			color.Cyan("Multiple Invite Mode")
			color.White("This will join your tokens from tokens.txt to servers from invite.txt")
			cfg, instances, err := getEverything()
			if err != nil {
				color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
			}

			if len(instances) == 0 {
				color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:04:05"))
				ExitSafely()
			}
			invites, err := utilities.ReadLines("invite.txt")
			if err != nil {
				color.Red("Error while opening invite.txt: %v", err)
				ExitSafely()
				return
			}
			if len(invites) == 0 {
				color.Red("[%v] Enter your invites in invite.txt", time.Now().Format("15:04:05"))
				ExitSafely()
				return
			}
			color.White("Enter delay between 2 consecutive joins by 1 token in seconds: ")
			var delay int
			fmt.Scanln(&delay)
			color.White("Enter number of Threads (0 for unlimited): ")
			var threads int
			fmt.Scanln(&threads)
			if threads > len(instances) {
				threads = len(instances)
			}
			if threads == 0 {
				threads = len(instances)
			}
			c := goccm.New(threads)
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				c.Wait()
				go func(i int) {
					for j := 0; j < len(invites); j++ {
						err := instances[i].Invite(invites[j])
						if err != nil {
							color.Red("[%v] Error while joining: %v", time.Now().Format("15:04:05"), err)
						}
						time.Sleep(time.Duration(delay) * time.Second)
					}
					c.Done()
				}(i)
			}
			c.WaitAllDone()
			color.Green("[%v] All threads finished", time.Now().Format("15:04:05"))
		}
	case 2:

		color.Cyan("Mass DM Advertiser/Spammer")
		color.White("This will DM everyone in memberids.txt from your tokens")
		members, err := utilities.ReadLines("memberids.txt")
		if err != nil {
			color.Red("Error while opening memberids.txt: %v", err)
			ExitSafely()
		}
		cfg, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
		}
		var msg utilities.Message
		color.White("Press 1 to use messages from file or press 2 to enter a message: ")
		var messagechoice int
		fmt.Scanln(&messagechoice)
		if messagechoice != 1 && messagechoice != 2 {
			color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if messagechoice == 2 {
			color.White("Enter your message, use \\n for changing lines. To use an embed, put message in message.json: ")
			scanner := bufio.NewScanner(os.Stdin)
			var text string
			if scanner.Scan() {
				text = scanner.Text()
			}

			msg.Content = text
			msg.Content = strings.Replace(msg.Content, "\\n", "\n", -1)
			var msgs []utilities.Message
			msgs = append(msgs, msg)
			err := setMessages(instances, msgs)
			if err != nil {
				color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
		} else {
			var msgs []utilities.Message
			err := setMessages(instances, msgs)
			if err != nil {
				color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
		}
		color.White("[%v] Do you wish to use Advanced Settings? 0: No, 1: Yes: ", time.Now().Format("15:04:05"))
		var advancedchoice int
		var checkchoice int
		var serverid string
		var tryjoinchoice int
		var invite string
		var maxattempts int
		fmt.Scanln(&advancedchoice)
		if advancedchoice != 0 && advancedchoice != 1 {
			color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if advancedchoice == 1 {
			color.White("[%v] Do you wish to check if token is still in server before every DM? [0: No, 1: Yes]", time.Now().Format("15:04:05"))
			fmt.Scanln(&checkchoice)
			if checkchoice != 0 && checkchoice != 1 {
				color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
				ExitSafely()
			}
			if checkchoice == 1 {
				color.White("[%v] Enter Server ID", time.Now().Format("15:04:05"))
				fmt.Scanln(&serverid)
				color.White("[%v] Do you wish to try rejoining the server if token is not in server? [0: No, 1: Yes]", time.Now().Format("15:04:05"))
				fmt.Scanln(&tryjoinchoice)
				if tryjoinchoice != 0 && tryjoinchoice != 1 {
					color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
					ExitSafely()
				}
				if tryjoinchoice == 1 {
					color.White("[%v] Enter a permanent invite code", time.Now().Format("15:04:05"))
					fmt.Scanln(&invite)
					color.White("[%v] Enter max rejoin attempts", time.Now().Format("15:04:05"))
					fmt.Scanln(&maxattempts)
				}
			}
		}

		// Also initiate variables and slices for logging and counting
		var session []string
		var completed []string
		var failed []string
		var dead []string
		var failedCount = 0
		completed, err = utilities.ReadLines("completed.txt")
		if err != nil {
			color.Red("Error while opening completed.txt: %v", err)
			ExitSafely()
		}
		if cfg.Skip {
			members = utilities.RemoveSubset(members, completed)
		}
		if len(instances) == 0 {
			color.Red("[%v] Enter your tokens in tokens.txt ", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if len(members) == 0 {
			color.Red("[%v] Enter your member ids in memberids.txt or ensure that all of them are not in completed.txt", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if len(members) < len(instances) {
			instances = instances[:len(members)]
		}
		msgs := instances[0].Messages
		for i := 0; i < len(msgs); i++ {
			if msgs[i].Content == "" && msgs[i].Embeds == nil {
				color.Red("[%v] WARNING: Message %v is empty", time.Now().Format("15:04:05"), i)
			}
		}
		var wg sync.WaitGroup
		wg.Add(len(instances))
		if cfg.Websocket {
			for i := 0; i < len(instances); i++ {
				go func(i int) {
					err := instances[i].StartWS()
					if err != nil {
						color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
					} else {
						color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), instances[i].Token)
					}
					wg.Done()
				}(i)
			}
			wg.Wait()
			color.Green("[%v] Websocket started for all instances", time.Now().Format("15:04:05"))
		}

		wg.Add(len(instances))

		start := time.Now()
		for i := 0; i < len(instances); i++ {
			// Offset goroutines by a few milliseconds. Makes a big difference and allows for better concurrency
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)

			go func(i int) {
				//defer wg.Done()

				for j := i * (len(members) / len(instances)); j < (i+1)*(len(members)/len(instances)); j++ {
					// Breaking loop if max DMs are reached. If cfg.MaxDMS is at 0, means no limit.
					if cfg.MaxDMS != 0 && instances[i].Count >= cfg.MaxDMS {
						color.Yellow("[%v] Max DMs reached for %v", time.Now().Format("15:04:05"), instances[i].Token)

						break
					}
					if cfg.Websocket && cfg.Receive {
						go func() {
							for {
								mes := <-instances[i].Ws.Messages
								if !strings.Contains(string(mes), "guild_id") {
									var mar utilities.Event
									err := json.Unmarshal(mes, &mar)
									if err != nil {
										color.Red("[%v] Error while unmarshalling websocket message: %v", time.Now().Format("15:04:05"), err)
										continue
									}
									if instances[i].ID == "" {
										tokenPart := strings.Split(instances[i].Token, ".")[0]
										dec, err := base64.StdEncoding.DecodeString(tokenPart)
										if err != nil {
											color.Red("[%v] Error while decoding token: %v", time.Now().Format("15:04:05"), err)
											continue
										}
										instances[i].ID = string(dec)
									}
									if mar.Data.Author.ID == instances[i].ID {
										continue
									}
									color.Green("[%v] %v#%v sent a message to %v : %v", time.Now().Format("15:04:05"), mar.Data.Author.Username, mar.Data.Author.Discriminator, instances[i].Token, mar.Data.Content)
									newStr := "Username: " + mar.Data.Author.Username + "#" + mar.Data.Author.Discriminator + "\nID: " + mar.Data.Author.ID + "\n" + "Message: " + mar.Data.Content + "\n"
									err = utilities.WriteLines("received.txt", newStr)
									if err != nil {
										color.Red("[%v] Error while opening received.txt: %v", time.Now().Format("15:04:05"), err)
									}
								}
							}

						}()
					}
					// Check if token is still valid at start of loop. Close instance is non-functional.
					status := instances[i].CheckToken()
					if status != 200 && status != 204 && status != 429 && status != -1 {
						failedCount += (i+1)*(len(members)/len(instances)) - j
						color.Red("[%v] Token %v might be locked - Stopping instance and adding members to failed list. %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, status, failedCount)
						failed = append(failed, members[j:(i+1)*(len(members)/len(instances))]...)
						dead = append(dead, instances[i].Token)
						err := Append("input/failed.txt", members[j:(i+1)*(len(members)/len(instances))])
						if err != nil {
							fmt.Println(err)
						}
						if cfg.Stop {

							break
						}

					}
					if cfg.Websocket && instances[i].Ws == nil {
						err := instances[i].StartWS()
						if err != nil {
							color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
						} else {
							color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), instances[i].Token)
						}
					}
					if advancedchoice == 1 {
						if checkchoice == 1 {
							r, err := instances[i].ServerCheck(serverid)
							if err != nil {
								color.Red("[%v] Error while checking server: %v", time.Now().Format("15:04:05"), err)
								continue
							}
							if r != 200 && r != 204 && r != 429 {
								if tryjoinchoice == 0 {
									color.Red("[%v] Stopping token %v [Not in server]", time.Now().Format("15:04:05"), instances[i].Token)

									break
								} else {
									if instances[i].Rejoin >= maxattempts {
										color.Red("[%v] Stopping token %v [Max server rejoin attempts]", time.Now().Format("15:04:05"), instances[i].Token)

										break
									}
									err := instances[i].Invite(invite)
									if err != nil {
										color.Red("[%v] Error while joining server: %v", time.Now().Format("15:04:05"), err)
										instances[i].Rejoin++
										continue
									}
								}
							}
						}
					}

					var user string
					user = members[j]
					// Get user info and check for mutual servers with the victim. Continue loop if no mutual servers or error.
					if cfg.Mutual {
						info, err := instances[i].UserInfo(members[j])
						if err != nil {
							failedCount++
							color.Red("[%v] Error while getting user info: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
							err = WriteLine("input/failed.txt", members[j])
							if err != nil {
								fmt.Println(err)
							}
							failed = append(failed, members[j])

							continue
						}
						if len(info.Mutual) == 0 {
							failedCount++
							color.Red("[%v] Token %v failed to DM %v [No Mutual Server] [%v]", time.Now().Format("15:04:05"), instances[i].Token, info.User.Username+info.User.Discriminator, failedCount)
							err = WriteLine("input/failed.txt", members[j])
							if err != nil {
								fmt.Println(err)
							}
							failed = append(failed, members[j])
							continue
						}
						user = info.User.Username + "#" + info.User.Discriminator
						// Used only if Websocket is enabled as Unwebsocketed Tokens get locked if they attempt to send friend requests.
						if cfg.Friend && cfg.Websocket {
							x, err := strconv.Atoi(info.User.Discriminator)
							if err != nil {
								color.Red("[%v] Error while adding friend: %v", time.Now().Format("15:04:05"), err)
							}
							resp, err := instances[i].Friend(info.User.Username, x)
							if err != nil {
								color.Red("[%v] Error while adding friend: %v", time.Now().Format("15:04:05"), err)
							}
							if resp.StatusCode != 204 && err != nil {
								body, _ := utilities.ReadBody(*resp)
								color.Red("[%v] Error while adding friend: %v", time.Now().Format("15:04:05"), string(body))
							} else {
								color.Green("[%v] Added friend %v", time.Now().Format("15:04:05"), info.User.Username+"#"+info.User.Discriminator)
							}
						}
					}

					// Send DM to victim. Continue loop if error.
					snowflake, err := instances[i].OpenChannel(members[j])
					if err != nil {
						failedCount++
						color.Red("[%v] Error while opening DM channel: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, members[j])
						continue
					}

					resp, err := instances[i].SendMessage(snowflake, members[j])
					if err != nil {
						failedCount++
						color.Red("[%v] Error while sending message: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, members[j])
						continue
					}
					body, err := utilities.ReadBody(resp)
					if err != nil {
						failedCount++
						color.Red("[%v] Error while reading body: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, members[j])
						continue
					}
					var response jsonResponse
					errx := json.Unmarshal(body, &response)
					if errx != nil {
						failedCount++
						color.Red("[%v] Error while unmarshalling body: %v [%v]", time.Now().Format("15:04:05"), errx, failedCount)
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, members[j])
						continue
					}
					// Everything is fine, continue as usual
					if resp.StatusCode == 200 {
						err = WriteLine("input/completed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						completed = append(completed, members[j])
						session = append(session, members[j])
						color.Green("[%v] Token %v sent DM to %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, len(session))
						if cfg.Websocket && cfg.Call {
							err := instances[i].Call(snowflake)
							if err != nil {
								color.Red("[%v] %v Error while calling %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, err)
							}
							// resp, err := utilities.Ring(instances[i].Client, instances[i].Token, snowflake)
							// if err != nil {
							//      color.Red("[%v] %v Error while ringing %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, err)
							// }
							// if resp == 200 || resp == 204 {
							//      color.Green("[%v] %v Ringed %v", time.Now().Format("15:04:05"), instances[i].Token, user)
							// } else {
							//      color.Red("[%v] %v Error while ringing %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, resp)
							// }

						}
						// Forbidden - Token is being rate limited
					} else if resp.StatusCode == 403 && response.Code == 40003 {

						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Yellow("[%v] Token %v sleeping for %v minutes!", time.Now().Format("15:04:05"), instances[i].Token, int(cfg.LongDelay/60))
						time.Sleep(time.Duration(cfg.LongDelay) * time.Second)
						color.Yellow("[%v] Token %v continuing!", time.Now().Format("15:04:05"), instances[i].Token)
						// Forbidden - DM's are closed
					} else if resp.StatusCode == 403 && response.Code == 50007 {
						failedCount++
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v failed to DM %v User has DMs closed or not present in server %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, string(body), failedCount)
						// Forbidden - Locked or Disabled
					} else if (resp.StatusCode == 403 && response.Code == 40002) || resp.StatusCode == 401 || resp.StatusCode == 405 {
						failedCount++
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v is locked or disabled. Stopping instance. %v %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, resp.StatusCode, string(body), failedCount)
						dead = append(dead, instances[i].Token)
						// Stop token if locked or disabled
						if cfg.Stop {

							break
						}
						// Forbidden - Invalid token
					} else if resp.StatusCode == 403 && response.Code == 50009 {
						failedCount++
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v can't DM %v. It may not have bypassed membership screening or it's verification level is too low or the server requires new members to wait 10 minutes before they can interact in the server. %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, string(body), failedCount)
						// General case - Continue loop. If problem with instance, it will be stopped at start of loop.
					} else if resp.StatusCode == 429 {
						color.Red("[%v] Token %v is being rate limited. Sleeping for 10 seconds", time.Now().Format("15:04:05"), instances[i].Token)
						time.Sleep(10 * time.Second)
					} else {
						failedCount++
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v couldn't DM %v Error Code: %v; Status: %v; Message: %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, response.Code, resp.Status, response.Message, failedCount)
					}
					time.Sleep(time.Duration(cfg.Delay) * time.Second)
					instances[i].Count++
				}
				wg.Done()
			}(i)
		}
		wg.Wait()

		color.Green("[%v] Threads have finished! Writing to file", time.Now().Format("15:04:05"))

		elapsed := time.Since(start)
		color.Green("[%v] DM advertisement took %v. Successfully sent DMs to %v IDs. Failed to send DMs to %v IDs. %v tokens are dis-functional & %v tokens are functioning", time.Now().Format("15:04:05"), elapsed.Seconds(), len(completed), len(failed), len(dead), len(instances)-len(dead))
		if cfg.Remove {
			var tokens []string
			for i := 0; i < len(instances); i++ {
				tokens = append(tokens, instances[i].Token)
			}
			m := utilities.RemoveSubset(tokens, dead)
			err := Truncate("input/tokens.txt", m)
			if err != nil {
				fmt.Println(err)
			}
			color.Green("Updated tokens.txt")
		}
		if cfg.RemoveM {
			m := utilities.RemoveSubset(members, completed)
			err := Truncate("input/memberids.txt", m)
			if err != nil {
				fmt.Println(err)
			}
			color.Green("Updated memberids.txt")

		}
		if cfg.Websocket {
			for i := 0; i < len(instances); i++ {
				instances[i].Ws.Close()
			}
		}

	case 3:
		color.Cyan("Single DM Spammer")
		color.White("Enter 0 for one message; Enter 1 for continuous spam")
		var choice int
		fmt.Scanln(&choice)
		cfg, instances, err := getEverything()
		if err != nil {
			fmt.Println(err)
			ExitSafely()
		}
		var msg utilities.Message
		color.White("Press 1 to use message from file or press 2 to enter a message: ")
		var messagechoice int
		fmt.Scanln(&messagechoice)
		if messagechoice != 1 && messagechoice != 2 {
			color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
			ExitSafely()
		}
		if messagechoice == 2 {
			color.White("Enter your message, use \\n for changing lines. To use an embed, put message in message.json: ")
			scanner := bufio.NewScanner(os.Stdin)
			var text string
			if scanner.Scan() {
				text = scanner.Text()
			}

			msg.Content = text
			msg.Content = strings.Replace(msg.Content, "\\n", "\n", -1)
			var msgs []utilities.Message
			msgs = append(msgs, msg)
			err := setMessages(instances, msgs)
			if err != nil {
				color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
		} else {
			var msgs []utilities.Message
			err := setMessages(instances, msgs)
			if err != nil {
				color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
		}

		color.White("Ensure a common link and enter victim's ID: ")
		var victim string
		fmt.Scanln(&victim)
		var wg sync.WaitGroup
		wg.Add(len(instances))
		if choice == 0 {
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)

				go func(i int) {
					defer wg.Done()
					snowflake, err := instances[i].OpenChannel(victim)
					if err != nil {
						fmt.Println(err)
					}
					resp, err := instances[i].SendMessage(snowflake, victim)
					if err != nil {
						fmt.Println(err)
					}
					body, err := utilities.ReadBody(resp)
					if err != nil {
						fmt.Println(err)
					}
					if resp.StatusCode == 200 {
						color.Green("[%v] Token %v DM'd %v", time.Now().Format("15:04:05"), instances[i].Token, victim)
					} else {
						color.Red("[%v] Token %v failed to DM %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, victim, string(body))
					}
				}(i)
			}
			wg.Wait()
		}
		if choice == 1 {
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()

					var c int
					for {
						snowflake, err := instances[i].OpenChannel(victim)
						if err != nil {
							fmt.Println(err)
						}
						resp, err := instances[i].SendMessage(snowflake, victim)
						if err != nil {
							fmt.Println(err)
						}
						if resp.StatusCode == 200 {
							color.Green("[%v] Token %v DM'd %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, victim, c)
						} else {
							color.Red("[%v] Token %v failed to DM %v", time.Now().Format("15:04:05"), instances[i].Token, victim)
						}
						c++
					}
				}(i)
				wg.Wait()
			}
		}
		color.Green("[%v] Threads have finished!", time.Now().Format("15:04:05"))

	case 4:
		color.Cyan("Reaction Adder")
		color.White("Note: You don't need to do this to send DMs in servers.")
		color.White("Menu:\n1) From message\n2) Manually")
		var choice int
		fmt.Scanln(&choice)
		cfg, instances, err := getEverything()
		if err != nil {
			fmt.Println(err)
			ExitSafely()
		}
		var wg sync.WaitGroup
		wg.Add(len(instances))
		if choice == 1 {
			color.Cyan("Enter a token which can see the message:")
			var token string
			fmt.Scanln(&token)
			color.White("Enter message ID: ")
			var id string
			fmt.Scanln(&id)
			color.White("Enter channel ID: ")
			var channel string
			fmt.Scanln(&channel)
			msg, err := utilities.GetRxn(channel, id, token)
			if err != nil {
				fmt.Println(err)
			}
			color.White("Select Emoji")
			for i := 0; i < len(msg.Reactions); i++ {
				color.White("%v) %v %v", i, msg.Reactions[i].Emojis.Name, msg.Reactions[i].Count)
			}
			var emoji int
			var send string
			fmt.Scanln(&emoji)
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					if msg.Reactions[emoji].Emojis.ID == "" {
						send = msg.Reactions[emoji].Emojis.Name

					} else if msg.Reactions[emoji].Emojis.ID != "" {
						send = msg.Reactions[emoji].Emojis.Name + ":" + msg.Reactions[emoji].Emojis.ID
					}
					err := instances[i].React(channel, id, send)
					if err != nil {
						fmt.Println(err)
						color.Red("[%v] %v failed to react", time.Now().Format("15:04:05"), instances[i].Token)
					} else {
						color.Green("[%v] %v reacted to the emoji", time.Now().Format("15:04:05"), instances[i].Token)
					}

				}(i)
			}
			wg.Wait()
			color.Green("[%v] Completed all threads.", time.Now().Format("15:04:05"))
		}
		if choice == 2 {
			color.Cyan("Enter channel ID")
			var channel string
			fmt.Scanln(&channel)
			color.White("Enter message ID")
			var id string
			fmt.Scanln(&id)
			color.Red("If you have a message, please use choice 1. If you want to add a custom emoji. Follow these instructions, if you don't, it won't work.\n If it's a default emoji which appears on the emoji keyboard, just copy it as TEXT not how it appears on Discord with the colons. Type it as text, it might look like 2 question marks on console but ignore.\n If it's a custom emoji (Nitro emoji) type it like this -> name:emojiID To get the emoji ID, copy the emoji link and copy the emoji ID from the URL.\nIf you do not follow this, it will not work. Don't try to do impossible things like trying to START a nitro reaction with a non-nitro account.")
			color.White("Enter emoji")
			var emoji string
			fmt.Scanln(&emoji)
			for i := 0; i < len(instances); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					err := instances[i].React(channel, id, emoji)
					if err != nil {
						fmt.Println(err)
						color.Red("[%v] %v failed to react", time.Now().Format("15:04:05"), instances[i].Token)
					}
					color.Green("[%v] %v reacted to the emoji", time.Now().Format("15:04:05"), instances[i].Token)
				}(i)
			}
			wg.Wait()
			color.Green("[%v] Completed all threads.", time.Now().Format("15:04:05"))
		}

	case 5:
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
		message, err := utilities.FindMessage(channelID, messageID, token)
		if err != nil {
			color.Red("Error while finding message: %v", err)
			ExitSafely()
			return
		}
		color.Green("[%v] Message: %v", time.Now().Format("15:04:05"), message)

	case 6:
		// Quick way to interconvert tokens from a popular format to the one this program supports.
		color.Cyan("Email:Password:Token to Token")
		Tokens, err := utilities.ReadLines("tokens.txt")
		if err != nil {
			color.Red("Error while opening tokens.txt: %v", err)
			ExitSafely()
			return
		}
		if len(Tokens) == 0 {
			color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:04:05"))
			ExitSafely()
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
			ExitSafely()
			return
		}

	case 7:
		// Basic token checker
		color.Cyan("Token checker")
		cfg, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
			ExitSafely()
		}
		color.White("Enter the number of threads: (0 for Unlimited)\n")
		var threads int
		fmt.Scanln(&threads)
		if threads > len(instances) {
			threads = len(instances)
		}
		if threads == 0 {
			threads = len(instances)
		}
		c := goccm.New(threads)
		var working []string
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
			c.Wait()
			go func(i int) {
				err := instances[i].CheckToken()
				if err != 200 {
					color.Red("[%v] Token Invalid %v", time.Now().Format("15:04:05"), instances[i].Token)
				} else {
					color.Green("[%v] Token Valid %v", time.Now().Format("15:04:05"), instances[i].Token)
					working = append(working, instances[i].Token)
				}
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		t := utilities.TruncateLines("tokens.txt", working)
		if t != nil {
			color.Red("[%v] Error while truncating tokens.txt: %v", time.Now().Format("15:04:05"), t)
			ExitSafely()
			return
		}

		color.Green("[%v] All threads finished", time.Now().Format("15:04:05"))

	case 8:
		// Leavs tokens from a server
		color.Cyan("Guild Leaver")
		cfg, instances, err := getEverything()
		if err != nil {
			color.Red("Error while getting necessary data %v", err)
			ExitSafely()

		}
		color.White("Enter the number of threads (0 for unlimited): ")
		var threads int
		fmt.Scanln(&threads)
		if threads > len(instances) {
			threads = len(instances)
		}
		if threads == 0 {
			threads = len(instances)
		}
		color.White("Enter delay between leaves: ")
		var delay int
		fmt.Scanln(&delay)
		color.White("Enter serverid: ")
		var serverid string
		fmt.Scanln(&serverid)
		c := goccm.New(threads)
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
			c.Wait()
			go func(i int) {
				p := instances[i].Leave(serverid)
				if p == 0 {
					color.Red("[%v] Error while leaving", time.Now().Format("15:04:05"))
				}
				if p == 200 || p == 204 {
					color.Green("[%v] Left server", time.Now().Format("15:04:05"))
				} else {
					color.Red("[%v] Error while leaving", time.Now().Format("15:04:05"))
				}
				time.Sleep(time.Duration(delay) * time.Second)
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] All threads finished", time.Now().Format("15:04:05"))
	case 9:

		color.Blue("Token Onliner")
		_, instances, err := getEverything()
		if err != nil {
			color.Red("Error while getting necessary data %v", err)
			ExitSafely()
		}
		var wg sync.WaitGroup
		wg.Add(len(instances))
		for i := 0; i < len(instances); i++ {
			go func(i int) {
				err := instances[i].StartWS()
				if err != nil {
					color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
				} else {
					color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), instances[i].Token)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		color.Green("[%v] All Token online. Press ENTER to disconnect and continue the program", time.Now().Format("15:04:05"))
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		wg.Add(len(instances))
		for i := 0; i < len(instances); i++ {
			go func(i int) {
				instances[i].Ws.Close()
				wg.Done()
			}(i)
		}
		wg.Wait()
		color.Green("[%v] All Token offline", time.Now().Format("15:04:05"))

	case 10:
		color.Blue("Scraping Menu")
		cfg, _, err := getEverything()
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

			Is := utilities.Instance{Token: token}
			t := 0
			for {
				if t >= 5 {
					color.Red("[%v] Couldn't connect to websocke after retrying.", time.Now().Format("15:04:05"))
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

			i := 0
			for {
				err := utilities.Scrape(Is.Ws, serverid, channelid, i)
				if err != nil {
					color.Red("[%v] Error while scraping: %v", time.Now().Format("15:04:05"), err)
				}
				color.Green("[%v] Token %v Scrape Count: %v", time.Now().Format("15:04:05"), Is.Token, len(Is.Ws.Members))
				if Is.Ws.Complete {
					break
				}
				i++
				time.Sleep(time.Duration(cfg.SleepSc) * time.Millisecond)
			}
			Is.Ws.Close()
			color.Green("[%v] Scraping finished. Scraped %v members", time.Now().Format("15:04:05"), len(Is.Ws.Members))
			clean := utilities.RemoveDuplicateStr(Is.Ws.Members)
			color.Green("[%v] Removed Duplicates. Scraped %v members", time.Now().Format("15:04:05"), len(clean))
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
				err := WriteFile("scraped/"+serverid+".txt", clean)
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
				msg, err := utilities.GetRxn(channelid, messageid, token)
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
			for {
				if len(allUIDS) == 0 {
					m = ""
				} else {
					m = allUIDS[len(allUIDS)-1]
				}
				rxn, err := utilities.GetReactions(channelid, messageid, token, send, m)
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
				for k := 0; k < len(allUIDS); k++ {
					err := utilities.WriteLines("memberids.txt", clean[k])
					if err != nil {
						color.Red("[%v] Error while writing to memberids.txt: %v", time.Now().Format("15:04:05"), err)
					}
				}
				color.Green("[%v] Wrote to memberids.txt", time.Now().Format("15:04:05"))
				err := WriteFile("scraped/"+messageid+".txt", allUIDS)
				if err != nil {
					color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
				}
			}
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
			cfg, instances, err := getEverything()
			if err != nil {
				color.Red("[%v] Error while getting config: %v", time.Now().Format("15:04:05"), err)
				ExitSafely()
			}
			var scraped []string
			// Input the number of tokens to be used
			color.Green("[%v] How many tokens do you wish to use? ", time.Now().Format("15:04:05"))
			var numTokens int
			quit := make(chan bool)
			var allQueries []string
			fmt.Scanln(&numTokens)

			chars := " 0123456789abcdefghijklmnopqrstuvwxyz"
			queriesLeft := make(chan string)
			var queriesCompleted []string

			for i := 0; i < len(chars); i++ {
				go func(i int) {
					queriesLeft <- string(chars[i])
				}(i)
			}

			if numTokens > len(instances) {
				color.Red("[%v] You only have %v tokens in your tokens.txt Using the maximum number of tokens possible", time.Now().Format("15:04:05"), len(instances))
			} else if numTokens <= 0 {
				color.Red("[%v] You must atleast use 1 token", time.Now().Format("15:04:05"))
				ExitSafely()
			} else if numTokens < len(instances) {
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
			// Starting the instances as GOroutines
			for i := 0; i < len(instances); i++ {
				go func(i int) {
					instances[i].ScrapeCount = 0
					for {
						instances[i].ScrapeCount++
						// Start websocket, reconnect if disconnected.
						if instances[i].ScrapeCount%5 == 0 {
							instances[i].Ws.Close()
							instances[i].Ws = nil
						}
						if instances[i].Ws == nil {
							err := instances[i].StartWS()
							if err != nil {
								color.Red("[%v] Error while starting websocket: %v", time.Now().Format("15:04:05"), err)
								continue
							}
							time.Sleep(3000 * time.Millisecond)
						}
						// Get a query from the channel / Await for close response
						select {
						case <-quit:
							return
						default:
							query := <-queriesLeft
							allQueries = append(allQueries, query)
							err := utilities.ScrapeOffline(instances[i].Ws, serverid, query)
							if err != nil {
								color.Red("[%v] %v Error while scraping: %v", time.Now().Format("15:04:05"), instances[i].Token, err)
								go func() {
									queriesLeft <- query
								}()
								continue
							}

							memInfo := <-instances[i].Ws.OfflineScrape
							queriesCompleted = append(queriesCompleted, query)
							var MemberInfo utilities.Event
							err = json.Unmarshal(memInfo, &MemberInfo)
							if err != nil {
								color.Red("[%v] Error while unmarshalling: %v", time.Now().Format("15:04:05"), err)
								continue
							}

							if len(MemberInfo.Data.Members) == 0 {
								continue
							}
							for _, member := range MemberInfo.Data.Members {
								// Avoiding Duplicates
								if !utilities.Contains(scraped, member.User.ID) {
									scraped = append(scraped, member.User.ID)
								}
							}
							color.Green("[%v] Token %v Query %v Scraped %v [+%v]", time.Now().Format("15:04:05"), instances[i].Token, query, len(scraped), len(MemberInfo.Data.Members))

							for i := 0; i < len(MemberInfo.Data.Members); i++ {
								id := MemberInfo.Data.Members[i].User.ID
								err := utilities.WriteLines("memberids.txt", id)
								if err != nil {
									color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
									continue
								}
							}
							if len(MemberInfo.Data.Members) < 100 {
								time.Sleep(time.Duration(cfg.SleepSc) * time.Millisecond)
								continue
							}
							lastName := MemberInfo.Data.Members[len(MemberInfo.Data.Members)-1].User.Username

							nextQueries := findNextQueries(query, lastName, queriesCompleted, chars)
							for i := 0; i < len(nextQueries); i++ {
								go func(i int) {
									queriesLeft <- nextQueries[i]
								}(i)
							}
							if len(MemberInfo.Data.Members) == 100 {
								time.Sleep(500 * time.Millisecond)
								instances[i].Ws.Close()
								instances[i].Ws = nil
								time.Sleep(1500 * time.Millisecond)
								err := instances[i].StartWS()
								if err != nil {
									color.Red("[%v] Error while starting websocket: %v", time.Now().Format("15:04:05"), err)
									continue
								}

							}

						}

					}
				}(i)
			}

			bufio.NewReader(os.Stdin).ReadBytes('\n')
			color.Green("[%v] Stopping All instances", time.Now().Format("15:04:05"))
			for i := 0; i < len(instances); i++ {
				go func() {
					quit <- true
				}()
			}

			color.Green("[%v] Scraping Complete. %v members scraped.", time.Now().Format("15:04:05"), len(scraped))

		}
	case 11:
		color.Blue("Name Changer")
		_, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
		}
		for i := 0; i < len(instances); i++ {
			if !reg.MatchString(instances[i].Token) {
				color.Red("[%v] Name changer requires tokens in email:pass:token format, there might be wrongly formatted tokens", time.Now().Format("15:04:05"))
				continue
			}
			fullz := instances[i].Token
			instances[i].Token = strings.Split(fullz, ":")[2]
			instances[i].Password = strings.Split(fullz, ":")[1]
		}
		color.Red("NOTE: Profile pictures are changed randomly from the file.")
		users, err := utilities.ReadLines("names.txt")
		if err != nil {
			color.Red("[%v] Error while reading names.txt: %v", time.Now().Format("15:04:05"), err)
			ExitSafely()
		}
		color.Green("[%v] Enter number of threads: ", time.Now().Format("15:04:05"))

		var threads int
		fmt.Scanln(&threads)
		if threads > len(instances) {
			threads = len(instances)
		}

		c := goccm.New(threads)
		for i := 0; i < len(instances); i++ {
			c.Wait()
			go func(i int) {
				r, err := instances[i].NameChanger(users[rand.Intn(len(users))])
				if err != nil {
					color.Red("[%v] %v Error while changing name: %v", time.Now().Format("15:04:05"), instances[i].Token, err)
				} else {
					if r.StatusCode == 200 || r.StatusCode == 204 {
						color.Green("[%v] %v Changed name successfully", time.Now().Format("15:04:05"), instances[i].Token)
					} else {
						color.Red("[%v] %v Error while changing name: %v", time.Now().Format("15:04:05"), instances[i].Token, r.Status)
					}
				}
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] All Done", time.Now().Format("15:04:05"))

	case 12:
		color.Blue("Profile Picture Changer")
		_, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
		}
		color.Red("NOTE: Only PNG and JPEG/JPG supported. Profile Pictures are changed randomly from the folder. Use PNG format for faster results.")
		color.White("Loading Avatars..")
		ex, err := os.Executable()
		if err != nil {
			color.Red("Couldn't find Exe")
			ExitSafely()
		}
		ex = filepath.ToSlash(ex)
		path := path.Join(path.Dir(ex) + "/input/pfps")

		images, err := utilities.GetFiles(path)
		if err != nil {
			color.Red("Couldn't find images in PFPs folder")
			ExitSafely()
		}
		color.Green("%v files found", len(images))
		var avatars []string

		for i := 0; i < len(images); i++ {
			av, err := utilities.EncodeImg(images[i])
			if err != nil {
				color.Red("Couldn't encode image")
				continue
			}
			avatars = append(avatars, av)
		}
		color.Green("%v avatars loaded", len(avatars))
		color.Green("[%v] Enter number of threads: ", time.Now().Format("15:04:05"))
		var threads int
		fmt.Scanln(&threads)
		if threads > len(instances) {
			threads = len(instances)
		}

		c := goccm.New(threads)
		for i := 0; i < len(instances); i++ {
			c.Wait()

			go func(i int) {
				r, err := instances[i].AvatarChanger(avatars[rand.Intn(len(avatars))])
				if err != nil {
					color.Red("[%v] %v Error while changing avatar: %v", time.Now().Format("15:04:05"), instances[i].Token, err)
				} else {
					if r.StatusCode == 204 || r.StatusCode == 200 {
						color.Green("[%v] %v Avatar changed successfully", time.Now().Format("15:04:05"), instances[i].Token)
					} else {
						color.Red("[%v] %v Error while changing avatar: %v", time.Now().Format("15:04:05"), instances[i].Token, r.StatusCode)
					}
				}

				c.Done()
			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] All done", time.Now().Format("15:04:05"))
	case 13:
		color.White("Check if your tokens are still in the server")
		_, instances, err := getEverything()
		if err != nil {
			color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
			ExitSafely()
		}
		var serverid string
		var inServer []string
		color.Green("[%v] Enter server ID: ", time.Now().Format("15:04:05"))
		fmt.Scanln(&serverid)
		var wg sync.WaitGroup
		wg.Add(len(instances))
		for i := 0; i < len(instances); i++ {
			go func(i int) {
				defer wg.Done()
				r, err := instances[i].ServerCheck(serverid)
				if err != nil {
					color.Red("[%v] %v Error while checking server: %v", time.Now().Format("15:04:05"), instances[i].Token, err)
				} else {
					if r == 200 || r == 204 {
						color.Green("[%v] %v is in server %v ", time.Now().Format("15:04:05"), instances[i].Token, serverid)
						inServer = append(inServer, instances[i].Token)
					} else if r == 429 {
						color.Green("[%v] %v is rate limited", time.Now().Format("15:04:05"), instances[i].Token)
					} else if r == 400 {
						color.Red("[%v] Bad request - Invalid Server ID", time.Now().Format("15:04:05"))
					} else {
						color.Red("[%v] %v is not in server [%v] [%v]", time.Now().Format("15:04:05"), instances[i].Token, serverid, r)
					}
				}
			}(i)
		}
		wg.Wait()
		color.Green("[%v] All done. Do you wish to save only tokens in the server to tokens.txt ? (y/n)", time.Now().Format("15:04:05"))
		var save string
		fmt.Scanln(&save)
		if save == "y" || save == "Y" {
			err := utilities.TruncateLines("tokens.txt", inServer)
			if err != nil {
				color.Red("[%v] Error while saving tokens: %v", time.Now().Format("15:04:05"), err)
			} else {
				color.Green("[%v] Tokens saved to tokens.txt", time.Now().Format("15:04:05"))
			}
		}

	case 14:
		color.Blue("Made with <3 by github.com/V4NSH4J for free. If you were sold this program, you got scammed. Full length documentation for this is available on the github readme.")
	case 15:
		// Exit without error
		os.Exit(0)

	}
	time.Sleep(1 * time.Second)
	Options()

}

type jsonResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func getEverything() (utilities.Config, []utilities.Instance, error) {
	var cfg utilities.Config
	var instances []utilities.Instance
	var err error
	var tokens []string
	var proxies []string
	var proxy string

	// Load config
	cfg, err = utilities.GetConfig()
	if err != nil {
		return cfg, instances, err
	}
	if cfg.Proxy != "" && os.Getenv("HTTPS_PROXY") != "" {
		os.Setenv("HTTPS_PROXY", "http://"+cfg.Proxy)
	}

	// Load instances
	tokens, err = utilities.ReadLines("tokens.txt")
	if err != nil {
		return cfg, instances, err
	}
	if len(tokens) == 0 {
		return cfg, instances, fmt.Errorf("no tokens found in tokens.txt")
	}
	if cfg.ProxyFromFile {
		proxies, err = utilities.ReadLines("proxies.txt")
		if err != nil {
			return cfg, instances, err
		}
		if len(proxies) == 0 {
			return cfg, instances, fmt.Errorf("no proxies found in proxies.txt")
		}
	}
	for i := 0; i < len(tokens); i++ {
		if cfg.ProxyFromFile {
			proxy = proxies[rand.Intn(len(proxies))]
		} else {
			proxy = ""
		}
		client, err := initClient(proxy)
		if err != nil {
			return cfg, instances, fmt.Errorf("couldn't initialize client: %v", err)
		}
		if !cfg.GatewayProxy {
			proxy = ""
		}
		instances = append(instances, utilities.Instance{Client: client, Token: tokens[i], Proxy: proxy})
	}
	return cfg, instances, nil

}

func setMessages(instances []utilities.Instance, messages []utilities.Message) error {
	var err error
	if len(messages) == 0 {
		messages, err = utilities.GetMessage()
		if err != nil {
			return err
		}
		if len(messages) == 0 {
			return fmt.Errorf("no messages found in messages.txt")
		}
		for i := 0; i < len(instances); i++ {
			instances[i].Messages = messages
		}
	} else {
		for i := 0; i < len(instances); i++ {
			instances[i].Messages = messages
		}
	}

	return nil
}

// Append items from slice to file
func Append(filename string, items []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range items {
		if _, err = file.WriteString(item + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Truncate items from slice to file
func Truncate(filename string, items []string) error {
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range items {
		if _, err = file.WriteString(item + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// Write line to file
func WriteLine(filename string, line string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(line + "\n"); err != nil {
		return err
	}

	return nil
}

// Create a New file and add items from a slice or append to it if it already exists
func WriteFile(filename string, items []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range items {
		if _, err = file.WriteString(item + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func initClient(proxy string) (*http.Client, error) {
	// If proxy is empty, return a default client (if proxy from file is false)
	if proxy == "" {
		return http.DefaultClient, nil
	}
	if !strings.Contains(proxy, "http://") {
		proxy = "http://" + proxy
	}
	// Error while converting proxy string to url.url would result in default client being returned
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return http.DefaultClient, err
	}
	// Creating a client and modifying the transport.
	Client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion:         tls.VersionTLS12,
				CipherSuites:       []uint16{0x1301, 0x1303, 0x1302, 0xc02b, 0xc02f, 0xcca9, 0xcca8, 0xc02c, 0xc030, 0xc00a, 0xc009, 0xc013, 0xc014, 0x009c, 0x009d, 0x002f, 0x0035, 0x000a},
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	return Client, nil

}

func ExitSafely() {
	color.Red("\nPress ENTER to EXIT")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)
}

const logo = "\r\n\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2557 \u2588\u2588\u2588\u2557   \u2588\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2588\u2588\u2557  \u2588\u2588\u2588\u2588\u2588\u2588\u2557  \u2588\u2588\u2588\u2588\u2588\u2588\u2557 \r\n\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2557 \u2588\u2588\u2588\u2588\u2551\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2550\u2550\u255D \u2588\u2588\u2554\u2550\u2550\u2550\u2588\u2588\u2557\r\n\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2554\u2588\u2588\u2588\u2588\u2554\u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2588\u2557\u2588\u2588\u2551   \u2588\u2588\u2551\r\n\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551\u255A\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2588\u2588\u2551  \u2588\u2588\u2551\u2588\u2588\u2551   \u2588\u2588\u2551\u2588\u2588\u2551   \u2588\u2588\u2551\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551 \u255A\u2550\u255D \u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\r\n\u255A\u2550\u2550\u2550\u2550\u2550\u255D \u255A\u2550\u255D     \u255A\u2550\u255D\u255A\u2550\u2550\u2550\u2550\u2550\u255D  \u255A\u2550\u2550\u2550\u2550\u2550\u255D  \u255A\u2550\u2550\u2550\u2550\u2550\u255D \r\nDISCORD MASS DM GO V1.0.7\n"

func findNextQueries(query string, lastName string, completedQueries []string, chars string) []string {
	if query == "" {
		color.Red("[%v] Query is empty", time.Now().Format("15:04:05"))
		return nil
	}
	lastName = strings.ToLower(lastName)
	indexQuery := strings.Index(lastName, query)
	if indexQuery == -1 {
		return nil
	}
	wantedCharIndex := indexQuery + len(query)
	if wantedCharIndex >= len(lastName) {

		return nil
	}
	wantedChar := lastName[wantedCharIndex]
	queryIndexDone := strings.Index(chars, string(wantedChar))
	if queryIndexDone == -1 {

		return nil
	}

	var nextQueries []string
	for j := queryIndexDone; j < len(chars); j++ {
		newQuery := query + string(chars[j])
		if !utilities.Contains(completedQueries, newQuery) && !strings.Contains(newQuery, "  ") && string(newQuery[0]) != "" {
			nextQueries = append(nextQueries, newQuery)
		}
	}
	return nextQueries
}
