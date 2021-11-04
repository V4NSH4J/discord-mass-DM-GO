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
	color.Green("\n\nMade by https://github.com/V4NSH4J ")

	color.Blue("Enter 0 for Single mode, Enter 1 for Multi mode (DM Advertising)")
	var option int
	fmt.Scanln(&option)
	if option != 0 && option != 1 {
		log.Panicf("Invalid Mode")
		return
	}

	type Config struct {
		Message   string `json:"message"`
		Delay     int    `json:"individual_delay"`
		LongDelay int    `json:"rate_limit_delay"`
		Offset    int    `json:"offset"`
		Skip      bool   `json:"skip_completed"`
	}
	var config Config
	ex, err := os.Executable()
	if err != nil {
		return
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "config.json"))
	if err != nil {
		return
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	errr := json.Unmarshal(bytes, &config)
	if errr != nil {
		fmt.Println(err)
		return
	}

	tokens, err := utilities.ReadLines("tokens.txt")
	if err != nil {
		fmt.Printf("Error while opening tokens.txt, %v \n", err)
		return
	}
	if option == 1 {
		color.Blue("\nMake sure everything is configured and press ENTER to begin SPAM")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		var memberids []string
		var completed []string
		var failed []string
		type jsonResponse struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}

		memberids, err = utilities.ReadLines("memberids.txt")
		if err != nil {
			fmt.Printf("[%v]Error while opening Memberids: %v \n", time.Now().Format("15:05:04"), err)
			return
		}
		completed, err = utilities.ReadLines("completed.txt")
		if err != nil {
			fmt.Printf("[%v]Error while opening Completed member list: %v \n", time.Now().Format("15:05:04"), err)
			return
		}

		if len(tokens) == 0 {
			fmt.Printf("[%v]No tokens loaded", time.Now().Format("15:05:04"))
			return
		}

		if len(memberids) == 0 {
			fmt.Printf("[%v]No Member ID's loaded", time.Now().Format("15:05:04"))
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
						a := directmessage.OpenChannel(tokens[i], memberids[j])
						b := directmessage.SendMessage(tokens[i], a, config.Message)
						defer b.Body.Close()

						body, err := ioutil.ReadAll(b.Body)
						if err != nil {
							log.Fatal(err)
						}
						var JsonB jsonResponse
						json.Unmarshal(body, &JsonB)
						if b.StatusCode == 200 {
							completed = append(completed, memberids[j])
							color.Green("[%v]Succesfully sent DM to %v", time.Now().Format("15:05:04"), memberids[j])
							w := utilities.WriteLines("completed.txt", memberids[j])
							if w != nil {
								fmt.Println(w)
							}

						} else if b.StatusCode == 403 && JsonB.Code == 40003 {
							color.Cyan("[%v] Token sleeping for %v minutes! Consider setting this delay to an appropriate amount (10-20 Minutes) to ensure your tokens last long!", tokens[i], int(config.LongDelay/60))
							time.Sleep(time.Duration(config.LongDelay) * time.Second)

						} else {
							failed = append(failed, memberids[j])
							color.Red("[%v]Failed to send DM to %v (Error %v)", time.Now().Format("15:05:04"), memberids[j], b)
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
					a := directmessage.OpenChannel(tokens[i], memberids[i])
					b := directmessage.SendMessage(tokens[i], a, config.Message)
					var JsonB jsonResponse
					if b.StatusCode == 200 {
						completed = append(completed, memberids[i])
						color.Green("[%v]Succesfully sent DM to %v", time.Now().Format("15:05:04"), memberids[i])

					} else if b.StatusCode == 403 && JsonB.Code == 40003 {
						time.Sleep(10 * time.Minute)
						color.Cyan("[%v] Token sleeping for 10 minutes!", tokens[i])
						time.Sleep(time.Duration(config.LongDelay) * time.Second)
					} else {
						failed = append(failed, memberids[i])
						color.Red("[%v]Failed to send DM to %v", time.Now().Format("15:05:04"), memberids[i])
					}
				}(i)
			}
			wg.Wait()
		}
		elapsed := time.Since(start)
		color.Blue("[%v]DM advertisement took %s. DM'd %v users and failed to DM %v users", time.Now().Format("15:05:04"), elapsed, len(completed), len(failed))
		fmt.Println("Writing to file, please wait!")

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

					a := directmessage.OpenChannel(tokens[i], UUID)
					b := directmessage.SendMessage(tokens[i], a, config.Message)
					if b.StatusCode == 200 {
						color.Green("[%v]Succesfully sent message from %v\n", time.Now().Format("15:05:04"), tokens[i])
					} else {
						color.Red("[%v]Failed to send message from %v\n", time.Now().Format("15:05:04"), tokens[i])
					}
				}(i)

			}
			wg.Wait()
		} else {
			var wg sync.WaitGroup
			wg.Add(len(tokens))
			for i := 0; i < len(tokens); i++ {
				go func(i int) {
					for {
						a := directmessage.OpenChannel(tokens[i], UUID)
						b := directmessage.SendMessage(tokens[i], a, config.Message)
						if b.StatusCode == 200 {
							color.Green("[%v]Succesfully sent message from %v\n", time.Now().Format("15:05:04"), tokens[i])

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
