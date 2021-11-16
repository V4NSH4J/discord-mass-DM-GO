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
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/directmessage"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func main() {

	color.Blue("\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2591\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2591\r\n\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2551\u2588\u2588\u2554\u2550\u2550\u2550\u2550\u255D\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\r\n\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\u2588\u2588\u2551\u255A\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2551\u2591\u2591\u255A\u2550\u255D\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\r\n\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\u2588\u2588\u2551\u2591\u255A\u2550\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2557\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u255A\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\r\n\u255A\u2550\u2550\u2550\u2550\u2550\u255D\u2591\u255A\u2550\u255D\u255A\u2550\u2550\u2550\u2550\u2550\u255D\u2591\u2591\u255A\u2550\u2550\u2550\u2550\u255D\u2591\u2591\u255A\u2550\u2550\u2550\u2550\u255D\u2591\u255A\u2550\u255D\u2591\u2591\u255A\u2550\u255D\u255A\u2550\u2550\u2550\u2550\u2550\u255D\u2591\r\n\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2557\u2591\u2591\u2591\u2588\u2588\u2588\u2557\r\n\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2588\u2551\r\n\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\u2588\u2588\u2554\u2588\u2588\u2588\u2588\u2554\u2588\u2588\u2551\r\n\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\u2588\u2588\u2551\u255A\u2588\u2588\u2554\u255D\u2588\u2588\u2551\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2591\u255A\u2550\u255D\u2591\u2588\u2588\u2551\r\n\u255A\u2550\u2550\u2550\u2550\u2550\u255D\u2591\u255A\u2550\u255D\u2591\u2591\u2591\u2591\u2591\u255A\u2550\u255D\r\n\r\n\u2591\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2591\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2557\u2591\u2591\u2591\u2588\u2588\u2588\u2557\u2588\u2588\u2588\u2557\u2591\u2591\u2591\u2588\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2591\r\n\u2588\u2588\u2554\u2550\u2550\u2550\u2550\u255D\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2588\u2551\u2588\u2588\u2554\u2550\u2550\u2550\u2550\u255D\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\r\n\u255A\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2551\u2588\u2588\u2554\u2588\u2588\u2588\u2588\u2554\u2588\u2588\u2551\u2588\u2588\u2554\u2588\u2588\u2588\u2588\u2554\u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2557\u2591\u2591\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\r\n\u2591\u255A\u2550\u2550\u2550\u2588\u2588\u2557\u2588\u2588\u2554\u2550\u2550\u2550\u255D\u2591\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2551\u2588\u2588\u2551\u255A\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2588\u2588\u2551\u255A\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2588\u2588\u2554\u2550\u2550\u255D\u2591\u2591\u2588\u2588\u2554\u2550\u2550\u2588\u2588\u2557\r\n\u2588\u2588\u2588\u2588\u2588\u2588\u2554\u255D\u2588\u2588\u2551\u2591\u2591\u2591\u2591\u2591\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\u2588\u2588\u2551\u2591\u255A\u2550\u255D\u2591\u2588\u2588\u2551\u2588\u2588\u2551\u2591\u255A\u2550\u255D\u2591\u2588\u2588\u2551\u2588\u2588\u2588\u2588\u2588\u2588\u2588\u2557\u2588\u2588\u2551\u2591\u2591\u2588\u2588\u2551\r\n\u255A\u2550\u2550\u2550\u2550\u2550\u255D\u2591\u255A\u2550\u255D\u2591\u2591\u2591\u2591\u2591\u255A\u2550\u255D\u2591\u2591\u255A\u2550\u255D\u255A\u2550\u255D\u2591\u2591\u2591\u2591\u2591\u255A\u2550\u255D\u255A\u2550\u255D\u2591\u2591\u2591\u2591\u2591\u255A\u2550\u255D\u255A\u2550\u2550\u2550\u2550\u2550\u2550\u255D\u255A\u2550\u255D\u2591\u2591\u255A\u2550\u255D")
	color.Green("\nV1.0.5\nMade by https://github.com/V4NSH4J ")
	type Config struct {
		Delay     int    `json:"individual_delay"`
		LongDelay int    `json:"rate_limit_delay"`
		Offset    int    `json:"offset"`
		Skip      bool   `json:"skip_completed"`
		Proxy     string `json:"proxy"`
		Call      bool   `json:"call"`
		Remove    bool   `json:"remove_dead_tokens"`
	}
	var config Config
	ex, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "config.json"))
	if err != nil {
		color.Red("Error while Opening config.json")
		fmt.Println(err)
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	errr := json.Unmarshal(bytes, &config)
	if errr != nil {
		fmt.Println(err)
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	if config.Proxy != "" {
		proxy := "http://" + config.Proxy + "/"
		proxys := "http://" + config.Proxy + "/"
		os.Setenv("HTTP_PROXY", proxy)
		os.Setenv("HTTPS_PROXY", proxys)

		color.Green("Now using Proxy %v", os.Getenv("HTTP_PROXY"))
	}
	color.Blue("Do you wish to join tokens to a server? Enter 0 for Yes, 1 for No")
	var invite int
	fmt.Scanln(&invite)
	if invite != 0 && invite != 1 {
		color.Red("Invalid option")
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	if invite == 0 {
		utilities.LaunchInviteJoiner()
	}
	color.Blue("Enter 0 for Single mode, Enter 1 for Multi mode (DM Advertising)")
	var option int
	fmt.Scanln(&option)
	if option != 0 && option != 1 {
		color.Red("Invalid Mode")
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	var message directmessage.Message
	exx, err := os.Executable()
	if err != nil {
		color.Red("Error while finding executable")
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	exx = filepath.ToSlash(exx)
	filee, errx := os.Open(path.Join(path.Dir(exx) + "/" + "message.json"))
	if err != nil {
		color.Red("Error while Opening message.json")
		fmt.Println(errx)
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	defer filee.Close()
	bytess, _ := io.ReadAll(filee)
	errrr := json.Unmarshal(bytess, &message)
	if errrr != nil {
		color.Red("Error while Unmarshalling Message - Please make sure your Embed colour is a decimal number and not a Hex. Also please make sure you're using the right kind of quotes \"  and not ' (You can add it to notepad and use Ctrl + H to replace all quotes with the right ones) and use escape characters (backslash n) to add new lines instead of EOF")
		fmt.Println(errrr)
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}



	tokens, err := utilities.ReadLines("tokens.txt")
	if err != nil {
		fmt.Printf("Error while opening tokens.txt, %v \n", err)
		color.Red("\nPress ENTER to EXIT")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	fingerprintSlice := make([]string, len(tokens))
	cookieSlice := make([]string, len(tokens))

	if option == 1 {
		color.Blue("\nMake sure everything is configured and press ENTER to begin SPAM")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		var memberids []string
		var completed []string
		var failed []string
		var deadtoken []string
		var left []string
		type jsonResponse struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}

		memberids, err = utilities.ReadLines("memberids.txt")
		if err != nil {
			fmt.Printf("[%v]Error while opening Memberids: %v \n", time.Now().Format("15:05:04"), err)
			color.Red("\nPress ENTER to EXIT")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		completed, err = utilities.ReadLines("completed.txt")
		if err != nil {
			fmt.Printf("[%v]Error while opening Completed member list: %v \n", time.Now().Format("15:05:04"), err)
			color.Red("\nPress ENTER to EXIT")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}

		if len(tokens) == 0 {
			fmt.Printf("[%v]No tokens loaded", time.Now().Format("15:05:04"))
			color.Red("\nPress ENTER to EXIT")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}

		if len(memberids) == 0 {
			fmt.Printf("[%v]No Member ID's loaded", time.Now().Format("15:05:04"))
			color.Red("\nPress ENTER to EXIT")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}
		var mode int
		if len(memberids) >= len(tokens) {
			mode = 1
		}
		if len(memberids) < len(tokens) {
			mode = 2
		}
		var wg sync.WaitGroup
		color.Green("[%v]Starting now", time.Now().Format("15:05:04"))
		start := time.Now()

		if mode == 1 {
			wg.Add(len(tokens))
			ChannelPerToken := len(memberids) / len(tokens)
			for i := 0; i < len(tokens); i++ {
				go func(i int) {
					defer wg.Done()
					time.Sleep(time.Duration(rand.Intn(config.Offset)) * time.Millisecond)
					x := ChannelPerToken * i
					y := x + ChannelPerToken

					for j := x; j < y; j++ {
						if len(completed) > 0 && config.Skip && utilities.Contains(completed, memberids[j]) {
							color.Green("[%v] Skipping Member %v [Already DM'd]", time.Now().Format("15:05:04"), memberids[j])
							continue
						}
						if fingerprintSlice[i] == "" {
							fingerprintSlice[i] = utilities.GetFingerprint()
							if fingerprintSlice[i] == "" {
								color.Red("[%v] No Fingerprint for <%v>", time.Now().Format("15:05:04"), tokens[i])
								continue
							}
						}
						if cookieSlice[i] == "" {
							Cookie := utilities.GetCookie()
							if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
								color.Red("[%v] No Cookie for <%v>", time.Now().Format("15:05:04"), tokens[i])
								continue
							}
							cookieSlice[i] = "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + " locale=us" + "; __cfruid=d2f75b0a2c63c38e6b3ab5226909e5184b1acb3e-1634536904"
						}
						a := directmessage.OpenChannel(tokens[i], memberids[j], cookieSlice[i], fingerprintSlice[i])
						b := directmessage.SendMessage(tokens[i], a, &message, memberids[j], cookieSlice[i], fingerprintSlice[i])
						defer b.Body.Close()

						body, err := ioutil.ReadAll(b.Body)
						if err != nil {
							log.Fatal(err)
						}
						var JsonB jsonResponse
						json.Unmarshal(body, &JsonB)
						if b.StatusCode == 200 {
							completed = append(completed, memberids[j])
							color.Green("[%v]Successfully sent DM to %v from <%v> [%v]", time.Now().Format("15:05:04"), memberids[j], tokens[i], len(completed))
							w := utilities.WriteLines("completed.txt", memberids[j])
							if w != nil {
								fmt.Println(w)
							}

						} else if b.StatusCode == 403 && JsonB.Code == 40003 {
							color.Cyan("[%v]Token <%v> sleeping for %v minutes! Consider setting this delay to an appropriate amount (10-20 Minutes) to ensure your tokens last long!", time.Now().Format("15:05:04"), tokens[i], int(config.LongDelay/60))
							time.Sleep(time.Duration(config.LongDelay) * time.Second)
							color.Cyan("[%v]Token <%v> waking up, starting DMs again", time.Now().Format("15:05:04"), tokens[i])
							q := utilities.WriteLines("failed.txt", memberids[j])
							if q != nil {
								fmt.Println(q)
							}
							failed = append(failed, memberids[j])

						} else if b.StatusCode == 403 && JsonB.Code == 50007 {
							color.Red("[%v] User %v has either closed DMs or is not in a mutual server or has blocked the token <%v>", time.Now().Format("15:05:04"), memberids[j], tokens[i])
							q := utilities.WriteLines("failed.txt", memberids[j])
							if q != nil {
								fmt.Println(q)
							}
							failed = append(failed, memberids[j])

						} else if b.StatusCode == 403 && JsonB.Code == 40002 {
							color.Red("[%v] Token <%v> is LOCKED. Adding it's memberIDs to the failed list & Stopping the instance", time.Now().Format("15:05:04"), tokens[i])
							deadtoken = append(deadtoken, tokens[i])
							for m := j; m < y; m++ {
								q := utilities.WriteLines("failed.txt", memberids[m])
								if q != nil {
									color.Red("[%v] Failed to write to failed.txt %v", time.Now().Format("15:05:04"), q)
								}
								failed = append(failed, memberids[m])
								left = append(left, memberids[m])

							}
							break
						} else if b.StatusCode == 401 {
							color.Red("[%v] Token <%v> is wrong or Disabled. Adding it's memberIDs to the failed list & Stopping instance", time.Now().Format("15:05:04"), tokens[i])
							deadtoken = append(deadtoken, tokens[i])
							for m := j; m < y; m++ {
								q := utilities.WriteLines("failed.txt", memberids[m])
								if q != nil {
									color.Red("[%v] Failed to write to failed.txt %v", time.Now().Format("15:05:04"), q)
								}
								failed = append(failed, memberids[m])
								left = append(left, memberids[m])

							}
							break

						} else if b.StatusCode == 405 && a == "" {
							color.Red("[%v] Token <%v> is wrong or Disabled. Adding it's memberIDs to the failed list & Stopping instance", time.Now().Format("15:05:04"), tokens[i])
							deadtoken = append(deadtoken, tokens[i])
							for m := j; m < y; m++ {
								q := utilities.WriteLines("failed.txt", memberids[m])
								if q != nil {
									color.Red("[%v] Failed to write to failed.txt %v", time.Now().Format("15:05:04"), q)
								}
								failed = append(failed, memberids[m])
								left = append(left, memberids[m])

							}
							break

						} else if b.StatusCode == 403 && JsonB.Code == 50009 {
							color.Red("[%v] Token <%v> can't DM %v - It might not have completed discord's community server member screening or the User is only accepting DMs from friends", time.Now().Format("15:05:04"), tokens[i], memberids[j])
							q := utilities.WriteLines("failed.txt", memberids[j])
							if q != nil {
								fmt.Println(q)
							}
							failed = append(failed, memberids[j])

						} else if b.StatusCode == 405 {
							color.Red("[%v] Token <%v> might be phone locked or disabled or may not have a mutual server  %v %v", time.Now().Format("15:05:04"), tokens[i], JsonB.Code, JsonB.Message)
							q := utilities.WriteLines("failed.txt", memberids[j])
							if q != nil {
								fmt.Println(q)
							}
							failed = append(failed, memberids[j])
						} else {
							failed = append(failed, memberids[j])
							color.Red("[%v]Failed to send DM to %v (Error %v) token <%v> - %v (%v)", time.Now().Format("15:05:04"), memberids[j], b.StatusCode, tokens[i], JsonB.Code, JsonB.Message)
							q := utilities.WriteLines("failed.txt", memberids[j])
							if q != nil {
								fmt.Println(q)
							}
						}
						time.Sleep(time.Duration(config.Delay) * time.Second)
					}
				}(i)

			}
			wg.Wait()

		}

		if mode == 2 {
			wg.Add(len(memberids))
			for i := 0; i < len(memberids); i++ {
				if len(completed) > 0 && config.Skip && utilities.Contains(completed, memberids[i]) {
					color.Green("[%v] Skipping Member %v [Already DM'd]", time.Now().Format("15:05:04"), memberids[i])
					continue
				}
				go func(i int) {
					defer wg.Done()
					if fingerprintSlice[i] == "" {
						fingerprintSlice[i] = utilities.GetFingerprint()
						if fingerprintSlice[i] == "" {
							color.Red("[%v] No Fingerprint for <%v>", time.Now().Format("15:05:04"), tokens[i])
							return
						}
					}
					if cookieSlice[i] == "" {
						Cookie := utilities.GetCookie()
						if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
							color.Red("[%v] No Cookie for <%v>", time.Now().Format("15:05:04"), tokens[i])
							return
						}
						cookieSlice[i] = "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + " locale=us" + "; __cfruid=d2f75b0a2c63c38e6b3ab5226909e5184b1acb3e-1634536904"

					}

					a := directmessage.OpenChannel(tokens[i], memberids[i], cookieSlice[i], fingerprintSlice[i])
					b := directmessage.SendMessage(tokens[i], a, &message, memberids[i], cookieSlice[i], fingerprintSlice[i])
					var JsonB jsonResponse
					if b.StatusCode == 200 {
						completed = append(completed, memberids[i])
						color.Green("[%v]Successfully sent DM to %v from <%v>", time.Now().Format("15:05:04"), memberids[i], tokens[i])

					} else if b.StatusCode == 403 && JsonB.Code == 40003 {
						time.Sleep(10 * time.Minute)
						color.Cyan("[%v] Token sleeping for 10 minutes!", tokens[i])
						time.Sleep(time.Duration(config.LongDelay) * time.Second)
					} else if b.StatusCode == 403 && JsonB.Code == 40002 {
						color.Red("[%v] Token <%v> is LOCKED. Adding it's memberIDs to the failed list & Stopping the instance", time.Now().Format("15:05:04"), tokens[i])
						deadtoken = append(deadtoken, tokens[i])
						q := utilities.WriteLines("failed.txt", memberids[i])
						if q != nil {
							color.Red("[%v] Failed to write to failed.txt %v", time.Now().Format("15:05:04"), q)
						}
						failed = append(failed, memberids[i])
						left = append(left, memberids[i])
					} else if b.StatusCode == 405 && a == "" {
						color.Red("[%v] Token <%v> is LOCKED. Adding it's memberIDs to the failed list & Stopping the instance", time.Now().Format("15:05:04"), tokens[i])
						deadtoken = append(deadtoken, tokens[i])
						q := utilities.WriteLines("failed.txt", memberids[i])
						if q != nil {
							color.Red("[%v] Failed to write to failed.txt %v", time.Now().Format("15:05:04"), q)
						}
						failed = append(failed, memberids[i])
						left = append(left, memberids[i])
					} else if b.StatusCode == 401 {
						color.Red("[%v] Token <%v> is wrong or Disabled. Adding it's memberIDs to the failed list & Stopping instance", time.Now().Format("15:05:04"), tokens[i])
						deadtoken = append(deadtoken, tokens[i])

						q := utilities.WriteLines("failed.txt", memberids[i])
						if q != nil {
							color.Red("[%v] Failed to write to failed.txt %v", time.Now().Format("15:05:04"), q)
						}
						failed = append(failed, memberids[i])
						left = append(left, memberids[i])

					} else {
						failed = append(failed, memberids[i])
						color.Red("[%v]Failed to send DM to %v from <%v> Code %v Err %v", time.Now().Format("15:05:04"), memberids[i], tokens[i], JsonB.Code, JsonB.Message)
					}
				}(i)
			}
			wg.Wait()
		}
		elapsed := time.Since(start)
		color.Blue("[%v]DM advertisement took %s. DM'd %v users and failed to DM %v users", time.Now().Format("15:05:04"), elapsed, len(completed), len(failed))
		fmt.Println("Writing to file, please wait!")
		if config.Remove {
			e := utilities.RemoveSubset(tokens, deadtoken)
			if err != nil {
				color.Red("[%v] Failed to remove dead tokens", time.Now().Format("15:05:04"))
			}
			utilities.TruncateLines("tokens.txt", e)
			utilities.TruncateLines("memberids.txt", left)

		}
		color.Blue("Do you wish to leave tokens from the server? 0 for yes, 1 for no")
		var think int
		fmt.Scanln(&think)
		if think != 0 && think != 1 {
			color.Red("Invalid Option")
			return
		}
		if think == 1 {
			color.Green("All completed")
			return
		}
		if think == 0 {
			var serverID string
			color.Blue("Enter Server ID")
			fmt.Scanln(&serverID)

			var wg sync.WaitGroup
			wg.Add(len(tokens))
			for i := 0; i < len(tokens); i++ {
				time.Sleep(5 * time.Millisecond)
				go func(i int) {
					defer wg.Done()
					if fingerprintSlice[i] == "" {
						fingerprintSlice[i] = utilities.GetFingerprint()
						if fingerprintSlice[i] == "" {
							color.Red("[%v] No Fingerprint for <%v> Couldn't leave server", time.Now().Format("15:05:04"), tokens[i])
							return
						}
					}
					if cookieSlice[i] == "" {
						Cookie := utilities.GetCookie()
						if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
							color.Red("[%v] No Cookie for <%v> Couldn't leave server", time.Now().Format("15:05:04"), tokens[i])
							return
						}
						cookieSlice[i] = "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + " locale=us" + "; __cfruid=d2f75b0a2c63c38e6b3ab5226909e5184b1acb3e-1634536904"

					}
					xa := utilities.Leave(serverID, tokens[i], fingerprintSlice[i])
					if xa == 204 || xa == 200 {
						color.Green("Successfully left from <%v>", tokens[i])
					} else {
						color.Red("Failed to leave from <%v>", tokens[i])
					}
				}(i)
			}
			wg.Wait()
			color.Red("\nPress ENTER to EXIT")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			return
		}

	}
	if option == 0 {
		color.Blue("Please make sure the tokens are in a mutual server as the victim and enter the victim's discord UID here: ")
		var UUID string
		fmt.Scanln(&UUID)
		color.Blue("Press 0 for single message, press 1 for continuous spam: ")
		var mode int
		fmt.Scanln(&mode)
		if mode != 0 && mode != 1 {
			log.Panic("Invalid mode")
			return
		}
		if mode == 0 {
			var wg sync.WaitGroup
			wg.Add(len(tokens))
			for i := 0; i < len(tokens); i++ {
				go func(i int) {
					color.Yellow("Loading")
					time.Sleep(500 * time.Millisecond)
					if fingerprintSlice[i] == "" {
						fingerprintSlice[i] = utilities.GetFingerprint()
						if fingerprintSlice[i] == "" {
							color.Red("[%v] No Fingerprint for %v", time.Now().Format("15:05:04"), tokens[i])
							return
						}
					}
					if cookieSlice[i] == "" {
						Cookie := utilities.GetCookie()
						if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
							color.Red("[%v] No Cookie for %v", time.Now().Format("15:05:04"), tokens[i])
							return
						}
						cookieSlice[i] = "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + " locale=us" + "; __cfruid=d2f75b0a2c63c38e6b3ab5226909e5184b1acb3e-1634536904"
					}
					a := directmessage.OpenChannel(tokens[i], UUID, cookieSlice[i], fingerprintSlice[i])
					b := directmessage.SendMessage(tokens[i], a, &message, UUID, cookieSlice[i], fingerprintSlice[i])
					if b.StatusCode == 200 {
						color.Green("[%v]Successfully sent message from %v\n", time.Now().Format("15:05:04"), tokens[i])
					} else {
						color.Red("[%v]Failed to send message from %v\n", time.Now().Format("15:05:04"), tokens[i])
					}
				}(i)

			}
			color.Blue("Do you wish to leave tokens from the server? 0 for yes, 1 for no")
			var think int
			fmt.Scanln(&think)
			if think != 0 && think != 1 {
				color.Red("Invalid Option")
				return
			}
			if think == 0 {
				color.Green("All completed")
				return
			}
			if think == 1 {
				var serverID string
				color.Blue("Enter Server ID")
				fmt.Scanln(&serverID)

				var wg sync.WaitGroup
				wg.Add(len(tokens))
				for i := 0; i < len(tokens); i++ {
					time.Sleep(5 * time.Millisecond)
					go func(i int) {
						defer wg.Done()
						if fingerprintSlice[i] == "" {
							fingerprintSlice[i] = utilities.GetFingerprint()
							if fingerprintSlice[i] == "" {
								color.Red("[%v] No Fingerprint for <%v> Couldn't leave server", time.Now().Format("15:05:04"), tokens[i])
								return
							}
						}
						if cookieSlice[i] == "" {
							Cookie := utilities.GetCookie()
							if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
								color.Red("[%v] No Cookie for <%v> Couldn't leave server", time.Now().Format("15:05:04"), tokens[i])
								return
							}
							cookieSlice[i] = "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + " locale=us" + "; __cfruid=d2f75b0a2c63c38e6b3ab5226909e5184b1acb3e-1634536904"

						}
						xa := utilities.Leave(serverID, tokens[i], cookieSlice[i])
						if xa == 204 || xa == 200 {
							color.Green("Successfully left from <%v>", tokens[i])
						} else {
							color.Red("Failed to leave from <%v>", tokens[i])
						}
					}(i)
				}
			}
			wg.Wait()
			color.Red("\nPress ENTER to EXIT")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		} else {
			var wg sync.WaitGroup
			wg.Add(len(tokens))
			for i := 0; i < len(tokens); i++ {
				go func(i int) {
					for {
						if fingerprintSlice[i] == "" {
							fingerprintSlice[i] = utilities.GetFingerprint()
							if fingerprintSlice[i] == "" {
								color.Red("[%v] No Fingerprint for %v", time.Now().Format("15:05:04"), tokens[i])
								return
							}
						}
						if cookieSlice[i] == "" {
							Cookie := utilities.GetCookie()
							if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
								color.Red("[%v] No Cookie for %v", time.Now().Format("15:05:04"), tokens[i])
								return
							}
							cookieSlice[i] = "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + " locale=us" + "; __cfruid=d2f75b0a2c63c38e6b3ab5226909e5184b1acb3e-1634536904"

						}
						a := directmessage.OpenChannel(tokens[i], UUID, cookieSlice[i], fingerprintSlice[i])
						b := directmessage.SendMessage(tokens[i], a, &message, UUID, cookieSlice[i], fingerprintSlice[i])
						if b.StatusCode == 200 {
							color.Green("[%v]Successfully sent message from %v\n", time.Now().Format("15:05:04"), tokens[i])

						} else {
							color.Red("[%v]Failed to send message from %v\n", time.Now().Format("15:05:04"), tokens[i])
							break
						}
						time.Sleep(time.Duration(config.Delay) * time.Second)
					}
				}(i)
			}
			wg.Wait()
		}
	}
}
