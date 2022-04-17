// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func LaunchMassDM() {

	color.Cyan("Mass DM Advertiser/Spammer")
	color.White("This will DM everyone in memberids.txt from your tokens")
	members, err := utilities.ReadLines("memberids.txt")
	if err != nil {
		color.Red("Error while opening memberids.txt: %v", err)
		utilities.ExitSafely()
	}
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
	}
	var msg instance.Message
	color.White("Press 1 to use messages from file or press 2 to enter a message: ")
	var messagechoice int
	fmt.Scanln(&messagechoice)
	if messagechoice != 1 && messagechoice != 2 {
		color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
		utilities.ExitSafely()
	}
	if messagechoice == 2 {
		color.White("Enter your message, use \\n for changing lines. You can also set a constant message in message.json")
		scanner := bufio.NewScanner(os.Stdin)
		var text string
		if scanner.Scan() {
			text = scanner.Text()
		}

		msg.Content = text
		msg.Content = strings.Replace(msg.Content, "\\n", "\n", -1)
		var msgs []instance.Message
		msgs = append(msgs, msg)
		err := instance.SetMessages(instances, msgs)
		if err != nil {
			color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
			utilities.ExitSafely()
		}
	} else {
		var msgs []instance.Message
		err := instance.SetMessages(instances, msgs)
		if err != nil {
			color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
			utilities.ExitSafely()
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
		utilities.ExitSafely()
	}
	if advancedchoice == 1 {
		color.White("[%v] Do you wish to check if token is still in server before every DM? [0: No, 1: Yes]", time.Now().Format("15:04:05"))
		fmt.Scanln(&checkchoice)
		if checkchoice != 0 && checkchoice != 1 {
			color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
			utilities.ExitSafely()
		}
		if checkchoice == 1 {
			color.White("[%v] Enter Server ID", time.Now().Format("15:04:05"))
			fmt.Scanln(&serverid)
			color.White("[%v] Do you wish to try rejoining the server if token is not in server? [0: No, 1: Yes]", time.Now().Format("15:04:05"))
			fmt.Scanln(&tryjoinchoice)
			if tryjoinchoice != 0 && tryjoinchoice != 1 {
				color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
				utilities.ExitSafely()
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
		utilities.ExitSafely()
	}
	if cfg.DirectMessage.Skip {
		members = utilities.RemoveSubset(members, completed)
	}
	if cfg.DirectMessage.SkipFailed {
		failedSkip, err := utilities.ReadLines("failed.txt")
		if err != nil {
			color.Red("Error while opening failed.txt: %v", err)
			utilities.ExitSafely()
		}
		members = utilities.RemoveSubset(members, failedSkip)
	}
	if len(instances) == 0 {
		color.Red("[%v] Enter your tokens in tokens.txt ", time.Now().Format("15:04:05"))
		utilities.ExitSafely()
	}
	if len(members) == 0 {
		color.Red("[%v] Enter your member ids in memberids.txt or ensure that all of them are not in completed.txt", time.Now().Format("15:04:05"))
		utilities.ExitSafely()
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
	// Send members to a channel
	mem := make(chan string, len(members))
	go func() {
		for i := 0; i < len(members); i++ {
			mem <- members[i]
		}
	}()
	// Setting information to windows titlebar by github.com/foxzsz
	go func() {
		for {
			cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%d sent, %v failed, %d locked, %v avg. dms, %d tokens left]`, len(session), len(failed), len(dead), len(session)/len(instances), len(instances)-len(dead)))
			_ = cmd.Run()
		}
	}()
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < len(instances); i++ {
		// Offset goroutines by a few milliseconds. Makes a big difference and allows for better concurrency
		time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				// Get a member from the channel
				if len(mem) == 0 {
					break
				}
				member := <-mem

				// Breaking loop if maximum DMs reached
				if cfg.DirectMessage.MaxDMS != 0 && instances[i].Count >= cfg.DirectMessage.MaxDMS {
					color.Yellow("[%v] Maximum DMs reached for %v", time.Now().Format("15:04:05"), instances[i].Token)
					break
				}
				// Start websocket connection if not already connected and reconnect if dead
				if cfg.DirectMessage.Websocket && instances[i].Ws == nil {
					err := instances[i].StartWS()
					if err != nil {
						color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
					} else {
						color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), instances[i].Token)
					}
				}
				if cfg.DirectMessage.Websocket && cfg.DirectMessage.Receive && instances[i].Ws != nil && !instances[i].Receiver {
					instances[i].Receiver = true
					go func() {
						for {
							if !instances[i].Receiver {
								break
							}
							mes := <-instances[i].Ws.Messages
							if !strings.Contains(string(mes), "guild_id") {
								var mar instance.Event
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
				// Check if token is valid
				status := instances[i].CheckToken()
				if status != 200 && status != 204 && status != 429 && status != -1 {
					failedCount++
					color.Red("[%v] Token %v might be locked - Stopping instance and adding members to failed list. %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, status, failedCount)
					failed = append(failed, member)
					dead = append(dead, instances[i].Token)
					err := utilities.WriteLines("failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					if cfg.DirectMessage.Stop {
						break
					}
				}
				// Advanced Options
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
								if instances[i].Retry >= maxattempts {
									color.Red("[%v] Stopping token %v [Max server rejoin attempts]", time.Now().Format("15:04:05"), instances[i].Token)
									break
								}
								err := instances[i].Invite(invite)
								if err != nil {
									color.Red("[%v] Error while joining server: %v", time.Now().Format("15:04:05"), err)
									instances[i].Retry++
									continue
								}
							}
						}
					}
				}
				var user string
				user = member
				// Check Mutual
				if cfg.DirectMessage.Mutual {
					info, err := instances[i].UserInfo(member)
					if err != nil {
						failedCount++
						color.Red("[%v] Error while getting user info: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
						err = utilities.WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, member)

						continue
					}
					if len(info.Mutual) == 0 {
						failedCount++
						color.Red("[%v] Token %v failed to DM %v [No Mutual Server] [%v]", time.Now().Format("15:04:05"), instances[i].Token, info.User.Username+info.User.Discriminator, failedCount)
						err = utilities.WriteLine("input/failed.txt", member)
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, member)
						continue
					}
					user = info.User.Username + "#" + info.User.Discriminator
					// Used only if Websocket is enabled as Unwebsocketed Tokens get locked if they attempt to send friend requests.
					if cfg.DirectMessage.Friend && cfg.DirectMessage.Websocket {
						x, err := strconv.Atoi(info.User.Discriminator)
						if err != nil {
							color.Red("[%v] Error while adding friend: %v", time.Now().Format("15:04:05"), err)
							continue
						}
						resp, err := instances[i].Friend(info.User.Username, x)
						if err != nil {
							color.Red("[%v] Error while adding friend: %v", time.Now().Format("15:04:05"), err)
							continue
						}
						if resp.StatusCode != 204 && err != nil {
							if !errors.Is(err, io.ErrUnexpectedEOF) {
								body, err := utilities.ReadBody(*resp)
								if err != nil {
									color.Red("[%v] Error while adding friend: %v", time.Now().Format("15:04:05"), fmt.Sprintf("error reading body: %v", err))
									continue
								}
								color.Red("[%v] Error while adding friend: %v", time.Now().Format("15:04:05"), string(body))
								continue
							}
							color.Red("[%v] Error while adding friend: %v", time.Now().Format("15:04:05"), err)
							continue
						} else {
							color.Green("[%v] Added friend %v", time.Now().Format("15:04:05"), info.User.Username+"#"+info.User.Discriminator)
						}
					}
				}
				// Open channel to get snowflake
				snowflake, err := instances[i].OpenChannel(member)
				if err != nil {
					failedCount++
					color.Red("[%v] Error while opening DM channel: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					failed = append(failed, member)
					continue
				}
				if cfg.SuspicionAvoidance.RandomDelayOpenChannel != 0 {
					time.Sleep(time.Duration(rand.Intn(cfg.SuspicionAvoidance.RandomDelayOpenChannel)) * time.Second)
				}
				resp, err := instances[i].SendMessage(snowflake, member)
				if err != nil {
					failedCount++
					color.Red("[%v] Error while sending message: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					failed = append(failed, member)
					continue
				}
				body, err := utilities.ReadBody(resp)
				if err != nil {
					failedCount++
					color.Red("[%v] Error while reading body: %v [%v]", time.Now().Format("15:04:05"), err, failedCount)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					failed = append(failed, member)
					continue
				}
				var response jsonResponse
				errx := json.Unmarshal(body, &response)
				if errx != nil {
					failedCount++
					color.Red("[%v] Error while unmarshalling body: %v [%v]", time.Now().Format("15:04:05"), errx, failedCount)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					failed = append(failed, member)
					continue
				}
				// Everything is fine, continue as usual
				if resp.StatusCode == 200 {
					err = utilities.WriteLine("input/completed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					completed = append(completed, member)
					session = append(session, member)
					color.Green("[%v] Token %v sent DM to %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, len(session))
					if cfg.DirectMessage.Websocket && cfg.DirectMessage.Call && instances[i].Ws != nil {
						err := instances[i].Call(snowflake)
						if err != nil {
							color.Red("[%v] %v Error while calling %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, err)
						}
						// Unfriended people can't ring.
						//
						// resp, err := instance.Ring(instances[i].Client, instances[i].Token, snowflake)
						// if err != nil {
						//      color.Red("[%v] %v Error while ringing %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, err)
						// }
						// if resp == 200 || resp == 204 {
						//      color.Green("[%v] %v Ringed %v", time.Now().Format("15:04:05"), instances[i].Token, user)
						// } else {
						//      color.Red("[%v] %v Error while ringing %v: %v", time.Now().Format("15:04:05"), instances[i].Token, user, resp)
						// }

					}
					if cfg.DirectMessage.Block {
						r, err := instances[i].BlockUser(member)
						if err != nil {
							color.Red("[%v] Error while blocking user: %v", time.Now().Format("15:04:05"), err)
						} else {
							if r == 204 {
								color.Green("[%v] Blocked %v", time.Now().Format("15:04:05"), user)
							} else {
								color.Red("[%v] Error while blocking user: %v", time.Now().Format("15:04:05"), r)
							}
						}
					}
					if cfg.DirectMessage.Close {
						r, err := instances[i].CloseDMS(snowflake)
						if err != nil {
							color.Red("[%v] Error while closing DM: %v", time.Now().Format("15:04:05"), err)
						} else {
							if r == 200 {
								color.Green("[%v] Succesfully closed DM %v", time.Now().Format("15:04:05"), user)
							} else {
								color.Red("[%v] Failed to close DM %v", time.Now().Format("15:04:05"), user)
							}
						}
					}
					// Forbidden - Token is being rate limited
				} else if resp.StatusCode == 403 && response.Code == 40003 {

					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					mem <- member
					color.Yellow("[%v] Token %v sleeping for %v minutes!", time.Now().Format("15:04:05"), instances[i].Token, int(cfg.DirectMessage.LongDelay/60))
					time.Sleep(time.Duration(cfg.DirectMessage.LongDelay) * time.Second)
					if cfg.SuspicionAvoidance.RandomRateLimitDelay != 0 {
						time.Sleep(time.Duration(rand.Intn(cfg.SuspicionAvoidance.RandomRateLimitDelay)) * time.Second)
					}
					color.Yellow("[%v] Token %v continuing!", time.Now().Format("15:04:05"), instances[i].Token)
					// Forbidden - DM's are closed
				} else if resp.StatusCode == 403 && response.Code == 50007 {
					failedCount++
					failed = append(failed, member)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					color.Red("[%v] Token %v failed to DM %v User has DMs closed or not present in server %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, string(body), failedCount)
					// Forbidden - Locked or Disabled
				} else if (resp.StatusCode == 403 && response.Code == 40002) || resp.StatusCode == 401 || resp.StatusCode == 405 {
					failedCount++
					failed = append(failed, member)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					color.Red("[%v] Token %v is locked or disabled. Stopping instance. %v %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, resp.StatusCode, string(body), failedCount)
					dead = append(dead, instances[i].Token)
					// Stop token if locked or disabled
					if cfg.DirectMessage.Stop {
						break
					}
					// Forbidden - Invalid token
				} else if resp.StatusCode == 403 && response.Code == 50009 {
					failedCount++
					failed = append(failed, member)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					color.Red("[%v] Token %v can't DM %v. It may not have bypassed membership screening or it's verification level is too low or the server requires new members to wait 10 minutes before they can interact in the server. %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, string(body), failedCount)
					// General case - Continue loop. If problem with instance, it will be stopped at start of loop.
				} else if resp.StatusCode == 429 {
					failed = append(failed, member)
					color.Red("[%v] Token %v is being rate limited. Sleeping for 10 seconds", time.Now().Format("15:04:05"), instances[i].Token)
					time.Sleep(10 * time.Second)
				} else if resp.StatusCode == 400 && strings.Contains(string(body), "captcha") {
					color.Red("[%v] Token %v Captcha was attempted to solve but appeared again", time.Now().Format("15:04:05"), instances[i].Token)
					instances[i].Retry++
					if instances[i].Retry >= cfg.CaptchaSettings.MaxCaptcha {
						color.Red("[%v] Stopping token %v max captcha solves reached", time.Now().Format("15:04:05"), instances[i].Token)
						break
					}
				} else {
					failedCount++
					failed = append(failed, member)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						fmt.Println(err)
					}
					color.Red("[%v] Token %v couldn't DM %v Error Code: %v; Status: %v; Message: %v [%v]", time.Now().Format("15:04:05"), instances[i].Token, user, response.Code, resp.Status, response.Message, failedCount)
				}
				time.Sleep(time.Duration(cfg.DirectMessage.Delay) * time.Second)
				if cfg.SuspicionAvoidance.RandomIndividualDelay != 0 {
					time.Sleep(time.Duration(rand.Intn(cfg.SuspicionAvoidance.RandomIndividualDelay)) * time.Second)
				}
			}
		}(i)
	}
	wg.Wait()

	color.Green("[%v] Threads have finished! Writing to file", time.Now().Format("15:04:05"))

	elapsed := time.Since(start)
	color.Green("[%v] DM advertisement took %v. Successfully sent DMs to %v IDs. Failed to send DMs to %v IDs. %v tokens are dis-functional & %v tokens are functioning", time.Now().Format("15:04:05"), elapsed.Seconds(), len(completed), len(failed), len(dead), len(instances)-len(dead))
	var left []string
	if cfg.DirectMessage.Remove {
		for i := 0; i < len(instances); i++ {
			if !utilities.Contains(dead, instances[i].Token) {
				if instances[i].Password == "" {
					left = append(left, instances[i].Token)
				} else {
					left = append(left, fmt.Sprintf(`%v:%v:%v`, instances[i].Email, instances[i].Password, instances[i].Token))
				}
			}
		}
		err := utilities.Truncate("input/tokens.txt", left)
		if err != nil {
			fmt.Println(err)
		}
		color.Green("Updated tokens.txt")
	}
	if cfg.DirectMessage.RemoveM {
		m := utilities.RemoveSubset(members, completed)
		err := utilities.Truncate("input/memberids.txt", m)
		if err != nil {
			fmt.Println(err)
		}
		color.Green("Updated memberids.txt")

	}
	if cfg.DirectMessage.Websocket {
		for i := 0; i < len(instances); i++ {
			if instances[i].Ws != nil {
				instances[i].Ws.Close()
			}
		}
	}

}

type jsonResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func LaunchSingleDM() {
	color.Cyan("Single DM Spammer")
	color.White("Enter 0 for one message; Enter 1 for continuous spam")
	var choice int
	fmt.Scanln(&choice)
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		fmt.Println(err)
		utilities.ExitSafely()
	}
	var msg instance.Message
	color.White("Press 1 to use message from file or press 2 to enter a message: ")
	var messagechoice int
	fmt.Scanln(&messagechoice)
	if messagechoice != 1 && messagechoice != 2 {
		color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
		utilities.ExitSafely()
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
		var msgs []instance.Message
		msgs = append(msgs, msg)
		err := instance.SetMessages(instances, msgs)
		if err != nil {
			color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
			utilities.ExitSafely()
		}
	} else {
		var msgs []instance.Message
		err := instance.SetMessages(instances, msgs)
		if err != nil {
			color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
			utilities.ExitSafely()
		}
	}

	color.White("Ensure a common link and enter victim's ID: ")
	var victim string
	fmt.Scanln(&victim)
	var wg sync.WaitGroup
	wg.Add(len(instances))
	if choice == 0 {
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)

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
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
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
}
