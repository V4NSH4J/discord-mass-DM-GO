// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"

	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/directmessage"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

type jsonResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ExitSafely() {
	color.Red("\nPress ENTER to EXIT")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func main() {

	// Credits
	color.Blue("\r\n\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2584    \u2584\u2584\u2584\u2584\u2588\u2588\u2588\u2584\u2584\u2584\u2584   \u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2584          \u2584\u2588\u2588\u2588\u2588\u2588\u2588\u2584   \u2584\u2588\u2588\u2588\u2588\u2588\u2588\u2584  \r\n\u2588\u2588\u2588   \u2580\u2588\u2588\u2588 \u2584\u2588\u2588\u2580\u2580\u2580\u2588\u2588\u2588\u2580\u2580\u2580\u2588\u2588\u2584 \u2588\u2588\u2588   \u2580\u2588\u2588\u2588        \u2588\u2588\u2588    \u2588\u2588\u2588 \u2588\u2588\u2588    \u2588\u2588\u2588 \r\n\u2588\u2588\u2588    \u2588\u2588\u2588 \u2588\u2588\u2588   \u2588\u2588\u2588   \u2588\u2588\u2588 \u2588\u2588\u2588    \u2588\u2588\u2588        \u2588\u2588\u2588    \u2588\u2580  \u2588\u2588\u2588    \u2588\u2588\u2588 \r\n\u2588\u2588\u2588    \u2588\u2588\u2588 \u2588\u2588\u2588   \u2588\u2588\u2588   \u2588\u2588\u2588 \u2588\u2588\u2588    \u2588\u2588\u2588       \u2584\u2588\u2588\u2588        \u2588\u2588\u2588    \u2588\u2588\u2588 \r\n\u2588\u2588\u2588    \u2588\u2588\u2588 \u2588\u2588\u2588   \u2588\u2588\u2588   \u2588\u2588\u2588 \u2588\u2588\u2588    \u2588\u2588\u2588      \u2580\u2580\u2588\u2588\u2588 \u2588\u2588\u2588\u2588\u2584  \u2588\u2588\u2588    \u2588\u2588\u2588 \r\n\u2588\u2588\u2588    \u2588\u2588\u2588 \u2588\u2588\u2588   \u2588\u2588\u2588   \u2588\u2588\u2588 \u2588\u2588\u2588    \u2588\u2588\u2588        \u2588\u2588\u2588    \u2588\u2588\u2588 \u2588\u2588\u2588    \u2588\u2588\u2588 \r\n\u2588\u2588\u2588   \u2584\u2588\u2588\u2588 \u2588\u2588\u2588   \u2588\u2588\u2588   \u2588\u2588\u2588 \u2588\u2588\u2588   \u2584\u2588\u2588\u2588        \u2588\u2588\u2588    \u2588\u2588\u2588 \u2588\u2588\u2588    \u2588\u2588\u2588 \r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2580   \u2580\u2588   \u2588\u2588\u2588   \u2588\u2580  \u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2580         \u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2580   \u2580\u2588\u2588\u2588\u2588\u2588\u2588\u2580  \r\n                                                                   \r\n\rDISCORD MASS DM GO")
	color.Green("\nV1.0.6 - Made by https://github.com/V4NSH4J ")
	color.Red("For educational purposes only. Read the full disclaimer and terms of use on GitHub readme file.")
	time.Sleep(2 * time.Second)
	// Check all files
	color.Green("[%v] Checking all files", time.Now().Format("15:05:04"))
	cfg, err := utilities.GetConfig()
	if err != nil {
		color.Red("Error while opening config.json: %v", err)
		ExitSafely()
		return
	}
	color.Green("[%v] Config Validated!", time.Now().Format("15:05:04"))

	msg, err := utilities.GetMessage()
	if err != nil {
		fmt.Printf("Error while opening message.json: %v", err)
		ExitSafely()
		return
	}
	color.Green("[%v] Message validated: %v\n", time.Now().Format("15:05:04"), msg)

	tkns, err := utilities.ReadLines("tokens.txt")
	if err != nil {
		color.Red("Error while opening tokens.txt: %v", err)
		ExitSafely()
		return
	}

	if len(tkns) == 0 {
		color.Red("[%v] Enter your tokens in tokens.txt")
		ExitSafely()
		return
	}
	color.Green("[%v] Tokens validated: %v tokens loaded \n", time.Now().Format("15:05:04"), len(tkns))

	members, err := utilities.ReadLines("memberids.txt")
	if err != nil {
		color.Red("Error while opening memberids.txt: %v", err)
		ExitSafely()
		return
	}

	if len(members) == 0 {
		color.Red("[%v] Enter your member ids in memberids.txt")
		ExitSafely()

		return
	}

	color.Green("[%v] Member ids validated: %v member ids loaded \n", time.Now().Format("15:05:04"), len(members))

	if cfg.Proxy != "" {
		color.Green("[%v] Now setting proxy as %v", time.Now().Format("15:05:04"), cfg.Proxy)
		os.Setenv("http_proxy", "http://"+cfg.Proxy)
		os.Setenv("https_proxy", "http://"+cfg.Proxy)
	} else {
		color.Green("[%v] Proxyless mode", time.Now().Format("15:05:04"))
	}
	// All Files validated.
	Options()
}

// Options menu
func Options() {
	color.Yellow("Leave a star on https://github.com/V4NSH4J/discord-mass-DM-GO for updates!")
	color.White("Menu:\n1) Invite Joiner\n2) Mass DM advertiser\n3) Single DM spam\n4) Reaction Adder\n5) Get message\n6) Email:Pass:Token to Token\n7) Token Checker\n8) Guild Leaver\n9) Credits & Info\n10) Exit")
	color.White("\nEnter your choice: ")
	var choice int
	fmt.Scanln(&choice)
	if choice != 1 && choice != 2 && choice != 3 && choice != 4 && choice != 5 && choice != 6 && choice != 7 && choice != 8 && choice != 9 && choice != 0 {
		color.Red("[%v] Invalid choice", time.Now().Format("15:05:04"))

		return
	}
	switch choice {
	case 1:
		var invitechoice int
		color.White("Invite Menu:\n1) Single Invite\n2) Multiple Invites from file")
		fmt.Scanln(&invitechoice)
		if invitechoice != 1 && invitechoice != 2 {
			color.Red("[%v] Invalid choice", time.Now().Format("15:05:04"))
			ExitSafely()
			return
		}
		switch invitechoice {
		case 1:
			color.Cyan("Single Invite Mode")
			color.White("Enter your invite CODE: ")
			var invite string
			fmt.Scanln(&invite)
			color.White("Enter number of Threads (0 for unlimited): ")
			var threads int
			fmt.Scanln(&threads)
			cfg, err := utilities.GetConfig()
			if err != nil {
				color.Red("Error while opening config.json: %v", err)
				ExitSafely()
				return
			}
			tokens, err := utilities.ReadLines("tokens.txt")
			if err != nil {
				color.Red("Error while opening tokens.txt: %v", err)
				ExitSafely()
				return
			}
			if len(tokens) == 0 {
				color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:05:04"))
				ExitSafely()
				return
			}
			if threads > len(tokens) {
				threads = len(tokens)
			}
			if threads == 0 {
				threads = len(tokens)
			}
			color.White("Enter base delay for joining in seconds (0 for none)")
			var base int
			fmt.Scanln(&base)
			color.White("Enter random delay to be added upon base delay (0 for none)")
			var random int
			fmt.Scanln(&random)
			var delay int
			if random > 0 {
				delay = base + rand.Intn(random)
			} else {
				delay = base
			}
			c := goccm.New(threads)
			for i := 0; i < len(tokens); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				c.Wait()
				go func(i int) {
					err := utilities.Invite(invite, tokens[i])
					if err != nil {
						color.Red("[%v] Error while joining: %v", time.Now().Format("15:05:04"), err)
					}
					time.Sleep(time.Duration(delay) * time.Second)
					c.Done()

				}(i)
			}
			c.WaitAllDone()
			color.Green("[%v] All threads finished", time.Now().Format("15:05:04"))
		case 2:
			color.Cyan("Multiple Invite Mode")
			cfg, err := utilities.GetConfig()
			if err != nil {
				color.Red("Error while opening config.json: %v", err)
				ExitSafely()
				return
			}
			tokens, err := utilities.ReadLines("tokens.txt")
			if err != nil {
				color.Red("Error while opening tokens.txt: %v", err)
				ExitSafely()
				return
			}
			if len(tokens) == 0 {
				color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:05:04"))
				ExitSafely()
				return
			}
			invites, err := utilities.ReadLines("invite.txt")
			if err != nil {
				color.Red("Error while opening invite.txt: %v", err)
				ExitSafely()
				return
			}
			if len(invites) == 0 {
				color.Red("[%v] Enter your invites in invite.txt", time.Now().Format("15:05:04"))
				ExitSafely()
				return
			}
			color.White("Enter delay between 2 consecutive joins by 1 token in seconds: ")
			var delay int
			fmt.Scanln(&delay)
			color.White("Enter number of Threads (0 for unlimited): ")
			var threads int
			fmt.Scanln(&threads)
			if threads > len(tokens) {
				threads = len(tokens)
			}
			if threads == 0 {
				threads = len(tokens)
			}
			c := goccm.New(threads)
			for i := 0; i < len(tokens); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				c.Wait()
				go func(i int) {
					for j := 0; j < len(invites); j++ {
						err := utilities.Invite(invites[j], tokens[i])
						if err != nil {
							color.Red("[%v] Error while joining: %v", time.Now().Format("15:05:04"), err)
						}
						time.Sleep(time.Duration(delay) * time.Second)
					}
					c.Done()
				}(i)
			}
			c.WaitAllDone()
			color.Green("[%v] All threads finished", time.Now().Format("15:05:04"))
		}
	case 2:
		// DM Advertiser - Blueprint
		// 1. Load all files and manage errors
		// 2. Manage the threading
		// 3. Check token {if working, continue; if not working - end instance & add IDs to failed}
		// 4. Check mutuals (if error/not mutual, continue loop & add to failed; if mutual, continue)
		// 5. Open the channel (if error, continue loop & add to failed; if no error, continue)
		// 6. Send DM (if sent, add to completed slice, print to file. If not sent; continue loop if require checking, close instance if locked)
		// 7. Truncate members with members left
		// 8. If all DMs gone, truncate tokens with tokens left
		// 9. Exit out to menu
		color.Cyan("Mass DM Advertiser/Spammer")
		color.Red("Please ensure you have used the invite joiner to join your tokens to the server and that they haven't been kicked/banned by an anti-raid bot")
		// Load files & Check for sources of error to prevent sudden crashes.
		// Also initiate variables and slices for logging and counting

		var completed []string
		var failed []string
		var dead []string
		completed, err := utilities.ReadLines("completed.txt")
		if err != nil {
			color.Red("Error while opening completed.txt: %v", err)
			ExitSafely()
			return
		}
		tokens, err := utilities.ReadLines("tokens.txt")
		if err != nil {
			color.Red("Error while opening tokens.txt: %v", err)
			ExitSafely()
			return
		}
		members, err := utilities.ReadLines("memberids.txt")
		if err != nil {
			color.Red("Error while opening members.txt: %v", err)
			ExitSafely()
			return
		}

		cfg, err := utilities.GetConfig()
		if err != nil {
			color.Red("Error while opening config.json: %v", err)
			ExitSafely()
			return
		}
		if cfg.Skip {
			members = utilities.RemoveSubset(members, completed)
		}
		msg, err := utilities.GetMessage()
		if err != nil {
			color.Red("Error while opening message.txt: %v", err)
			ExitSafely()
			return
		}
		if len(tokens) == 0 || len(members) == 0 {
			color.Red("[%v] Enter your tokens in tokens.txt and members in members.txt", time.Now().Format("15:05:04"))
			ExitSafely()
			return
		}
		if len(members) < len(tokens) {
			tokens = tokens[:len(members)]
		}
		// Threading system to prevent large-scale concurrency. Might be detectable.

		var wg sync.WaitGroup
		wg.Add(len(tokens))

		start := time.Now()
		for i := 0; i < len(tokens); i++ {
			// Offset goroutines by a few milliseconds. Made a big difference in my testing.
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)

			go func(i int) {
				for j := i * (len(members) / len(tokens)); j < (i+1)*(len(members)/len(tokens)); j++ {
					// Check if token is still valid at start of loop. Close instance is non-functional.
					status := utilities.CheckToken(tokens[i])
					if status != 200 && status != 204 && status != 429 && status != -1 {
						color.Red("[%v] Token %v might be locked - Stopping instance and adding members to failed list. %v", time.Now().Format("15:05:04"), tokens[i], status)
						failed = append(failed, members[j:(i+1)*(len(members)/len(tokens))]...)
						dead = append(dead, tokens[i])
						err := Append("input/failed.txt", members[j:(i+1)*(len(members)/len(tokens))])
						if err != nil {
							fmt.Println(err)
						}
						if cfg.Stop {
							break
						}

					}
					var user string
					user = members[j]
					// Get user info and check for mutual servers with the victim. Continue loop if no mutual servers or error.
					if cfg.Mutual {
						info, err := directmessage.UserInfo(tokens[i], members[j])
						if err != nil {
							color.Red("[%v] Error while getting user info: %v", time.Now().Format("15:05:04"), err)
							err = WriteLine("input/failed.txt", members[j])
							if err != nil {
								fmt.Println(err)
							}
							failed = append(failed, members[j])

							continue
						}
						if len(info.Mutual) == 0 {
							color.Red("[%v] Token %v failed to DM %v [No Mutual Server]", time.Now().Format("15:05:04"), tokens[i], info.User.Username+info.User.Discriminator)
							err = WriteLine("input/failed.txt", members[j])
							if err != nil {
								fmt.Println(err)
							}
							failed = append(failed, members[j])
							continue
						}
						user = info.User.Username + "#" + info.User.Discriminator
					}

					// Send DM to victim. Continue loop if error.
					snowflake, err := directmessage.OpenChannel(tokens[i], members[j])
					if err != nil {
						color.Red("[%v] Error while opening DM channel: %v", time.Now().Format("15:05:04"), err)
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, members[j])
						continue
					}

					resp, err := directmessage.SendMessage(tokens[i], snowflake, &msg, members[j])
					if err != nil {
						color.Red("[%v] Error while sending message: %v", time.Now().Format("15:05:04"), err)
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						failed = append(failed, members[j])
						continue
					}
					body, err := utilities.ReadBody(*resp)
					if err != nil {
						color.Red("[%v] Error while reading body: %v", time.Now().Format("15:05:04"), err)
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
						color.Red("[%v] Error while unmarshalling body: %v", time.Now().Format("15:05:04"), errx)
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
						color.Green("[%v] Token %v sent DM to %v [%v]", time.Now().Format("15:05:04"), tokens[i], user, len(completed))
						// Case-based error, something unusual with data enterred
					} else if resp.StatusCode == 400 {
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v failed to DM %v Check wether the token tried to DM itself or tried sending an empty message!", time.Now().Format("15:05:04"), tokens[i], user)
						// Forbidden - Token is being rate limited
					} else if resp.StatusCode == 403 && response.Code == 40003 {
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Cyan("[%v] Token %v sleeping for %v minutes!", time.Now().Format("15:05:04"), tokens[i], int(cfg.LongDelay/60))
						time.Sleep(time.Duration(cfg.LongDelay) * time.Second)
						color.Cyan("[%v] Token %v continuing!", time.Now().Format("15:05:04"), tokens[i])
						// Forbidden - DM's are closed
					} else if resp.StatusCode == 403 && response.Code == 50007 {
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v failed to DM %v User has DMs closed or not present in server", time.Now().Format("15:05:04"), tokens[i], user)
						// Forbidden - Locked or Disabled
					} else if (resp.StatusCode == 403 && response.Code == 40002) || resp.StatusCode == 401 || resp.StatusCode == 405 {
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v is locked or disabled. Stopping instance. %v %v", time.Now().Format("15:05:04"), tokens[i], resp.StatusCode, response.Message)
						dead = append(dead, tokens[i])
						// Stop token if locked or disabled
						if cfg.Stop {
							break
						}
						// Forbidden - Invalid token
					} else if resp.StatusCode == 403 && response.Code == 50009 {
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v can't DM %v. It might not have bypassed community screening.", time.Now().Format("15:05:04"), tokens[i], user)
						// General case - Continue loop. If problem with instance, it will be stopped at start of loop.
					} else {
						err = WriteLine("input/failed.txt", members[j])
						if err != nil {
							fmt.Println(err)
						}
						color.Red("[%v] Token %v couldn't DM %v Error Code: %v; Status: %v; Message: %v", time.Now().Format("15:05:04"), tokens[i], user, response.Code, resp.Status, response.Message)
					}
					time.Sleep(time.Duration(cfg.Delay) * time.Second)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()

		color.Green("[%v] Threads have finished! Writing to file", time.Now().Format("15:05:04"))
		elapsed := time.Since(start)
		color.Green("[%v] DM advertisement took %v. Successfully sent DMs to %v IDs. Failed to send DMs to %v IDs. %v tokens are dis-functional & %v tokens are functioning", time.Now().Format("15:04:05"), elapsed.Seconds(), len(completed), len(failed), len(dead), len(tokens)-len(dead))
		if cfg.Remove {
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
	case 3:
		color.Cyan("Single DM Spammer")
		color.White("Enter 0 for one message; Enter 1 for continuous spam")
		var choice int
		fmt.Scanln(&choice)
		tokens, err := utilities.ReadLines("tokens.txt")
		if err != nil {
			fmt.Println(err)
		}
		cfg, err := utilities.GetConfig()
		if err != nil {
			fmt.Println(err)
		}
		msg, err := utilities.GetMessage()
		if err != nil {
			fmt.Println(err)
		}
		color.White("Ensure a common link and enter victim's ID: ")
		var victim string
		fmt.Scanln(&victim)
		var wg sync.WaitGroup
		wg.Add(len(tokens))
		if choice == 0 {
			for i := 0; i < len(tokens); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					snowflake, err := directmessage.OpenChannel(tokens[i], victim)
					if err != nil {
						fmt.Println(err)
					}
					resp, err := directmessage.SendMessage(tokens[i], snowflake, &msg, victim)
					if err != nil {
						fmt.Println(err)
					}
					if resp.StatusCode == 200 {
						color.Green("[%v] Token %v DM'd %v", time.Now().Format("15:05:04"), tokens[i], victim)
					} else {
						color.Red("[%v] Token %v failed to DM %v", time.Now().Format("15:05:04"), tokens[i], victim)
					}
				}(i)
			}
			wg.Wait()
		}
		if choice == 1 {
			for i := 0; i < len(tokens); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					var c int
					for {
						snowflake, err := directmessage.OpenChannel(tokens[i], victim)
						if err != nil {
							fmt.Println(err)
						}
						resp, err := directmessage.SendMessage(tokens[i], snowflake, &msg, victim)
						if err != nil {
							fmt.Println(err)
						}
						if resp.StatusCode == 200 {
							color.Green("[%v] Token %v DM'd %v [%v]", time.Now().Format("15:05:04"), tokens[i], victim, c)
						} else {
							color.Red("[%v] Token %v failed to DM %v", time.Now().Format("15:05:04"), tokens[i], victim)
						}
						c++
					}
				}(i)
				wg.Wait()
			}
		}
		color.Green("[%v] Threads have finished!", time.Now().Format("15:05:04"))

	case 4:
		color.Cyan("Reaction Adder")
		color.White("Menu:\n1) From message\n2) Manually")
		var choice int
		fmt.Scanln(&choice)
		tokens, err := utilities.ReadLines("tokens.txt")
		if err != nil {
			fmt.Println(err)
		}
		cfg, err := utilities.GetConfig()
		if err != nil {
			fmt.Println(err)
		}
		var wg sync.WaitGroup
		wg.Add(len(tokens))
		if choice == 1 {
			color.White("Enter a token which can see the message:")
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
			for i := 0; i < len(tokens); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					if msg.Reactions[emoji].Emojis.ID == "" {
						send = msg.Reactions[emoji].Emojis.Name

					} else if msg.Reactions[emoji].Emojis.ID != "" {
						send = msg.Reactions[emoji].Emojis.Name + "" + msg.Reactions[emoji].Emojis.ID
					}
					err := utilities.React(tokens[i], channel, id, send)
					if err != nil {
						fmt.Println(err)
						color.Red("[%v] %v failed to react", time.Now().Format("15:05:04"), tokens[i])
					}
					color.Green("[%v] %v reacted to the emoji", time.Now().Format("15:05:04"), tokens[i])

				}(i)
			}
			wg.Wait()
			color.Green("[%v] Completed all threads.", time.Now().Format("15:05:04"))
		}
		if choice == 2 {
			color.White("Enter channel ID")
			var channel string
			fmt.Scanln(&channel)
			color.White("Enter message ID")
			var id string
			fmt.Scanln(&id)
			color.Red("If you have a message, please use choice 1. If you want to add a custom emoji. Follow these instructions, if you don't, it won't work.\n If it's a default emoji which appears on the emoji keyboard, just copy it as TEXT not how it appears on Discord with the colons. Type it as text, it might look like 2 question marks on console but ignore.\n If it's a custom emoji (Nitro emoji) type it like this -> name:emojiID To get the emoji ID, copy the emoji link and copy the emoji ID from the URL.\nIf you do not follow this, it will not work. Don't try to do impossible things like trying to START a nitro reaction with a non-nitro account.")
			color.White("Enter emoji")
			var emoji string
			fmt.Scanln(&emoji)
			for i := 0; i < len(tokens); i++ {
				time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					err := utilities.React(tokens[i], channel, id, emoji)
					if err != nil {
						fmt.Println(err)
						color.Red("[%v] %v failed to react", time.Now().Format("15:05:04"), tokens[i])
					}
					color.Green("[%v] %v reacted to the emoji", time.Now().Format("15:05:04"), tokens[i])
				}(i)
			}
			wg.Wait()
			color.Green("[%v] Completed all threads.", time.Now().Format("15:05:04"))
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
		color.Green("[%v] Message: %v", time.Now().Format("15:05:04"), message)

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
			color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:05:04"))
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
			color.Red("[%v] Error while truncating tokens.txt: %v", time.Now().Format("15:05:04"), t)
			ExitSafely()
			return
		}
	case 7:
		// Basic token checker
		color.Cyan("Token checker")
		tokens, err := utilities.ReadLines("tokens.txt")
		if err != nil {
			color.Red("Error while opening tokens.txt: %v", err)
			ExitSafely()
			return
		}
		if len(tokens) == 0 {
			color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:05:04"))
			ExitSafely()
			return
		}
		cfg, err := utilities.GetConfig()
		if err != nil {
			color.Red("Error while opening config.json: %v", err)
			ExitSafely()
			return
		}
		color.White("Enter the number of threads: \n")
		var threads int
		fmt.Scanln(&threads)
		if threads > len(tokens) {
			threads = len(tokens)
		}
		if threads == 0 {
			threads = len(tokens)
		}
		c := goccm.New(threads)
		var working []string
		for i := 0; i < len(tokens); i++ {
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
			c.Wait()
			go func(i int) {
				err := utilities.CheckToken(tokens[i])
				if err != 200 {
					color.Red("[%v] Token Invalid %v", time.Now().Format("15:05:04"), tokens[i])
				} else {
					color.Green("[%v] Token Valid %v", time.Now().Format("15:05:04"), tokens[i])
					working = append(working, tokens[i])
				}
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		t := utilities.TruncateLines("tokens.txt", working)
		if t != nil {
			color.Red("[%v] Error while truncating tokens.txt: %v", time.Now().Format("15:05:04"), t)
			ExitSafely()
			return
		}
		color.Green("[%v] All threads finished", time.Now().Format("15:05:04"))

	case 8:
		// Leavs tokens from a server
		color.Cyan("Guild Leaver")
		cfg, err := utilities.GetConfig()
		if err != nil {
			color.Red("Error while opening config.json: %v", err)
			ExitSafely()
			return
		}
		tokens, err := utilities.ReadLines("tokens.txt")
		if err != nil {
			color.Red("Error while opening tokens.txt: %v", err)
			ExitSafely()
			return
		}
		if len(tokens) == 0 {
			color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:05:04"))
			ExitSafely()
			return
		}
		color.White("Enter the number of threads (0 for unlimited): ")
		var threads int
		fmt.Scanln(&threads)
		if threads > len(tokens) {
			threads = len(tokens)
		}
		if threads == 0 {
			threads = len(tokens)
		}
		color.White("Enter delay between leaves: ")
		var delay int
		fmt.Scanln(&delay)
		color.White("Enter serverid: ")
		var serverid string
		fmt.Scanln(&serverid)
		c := goccm.New(threads)
		for i := 0; i < len(tokens); i++ {
			time.Sleep(time.Duration(cfg.Offset) * time.Millisecond)
			c.Wait()
			go func(i int) {
				p := utilities.Leave(serverid, tokens[i])
				if p == 0 {
					color.Red("[%v] Error while leaving", time.Now().Format("15:05:04"))
				}
				if p == 200 || p == 204 {
					color.Green("[%v] Left server", time.Now().Format("15:05:04"))
				} else {
					color.Red("[%v] Error while leaving", time.Now().Format("15:05:04"))
				}
				time.Sleep(time.Duration(delay) * time.Second)
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] All threads finished", time.Now().Format("15:05:04"))

	case 9:
		color.Blue("Made with <3 by github.com/V4NSH4J for free. If you were sold this program, you got scammed.")
	case 10:
		// Exit without error
		os.Exit(0)

	}
	time.Sleep(1 * time.Second)
	Options()

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
