// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchReactionAdder() {
	utilities.PrintMenu([]string{"From Message", "Manually"})
	choice := utilities.UserInputInteger("Select an option: ")
	if choice != 1 && choice != 2 {
		utilities.LogErr("Invalid option")
		return
	}
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting instances or config %s", err)
		return
	}
	var tokenFile, successFile, failedFile string
	if cfg.OtherSettings.Logs {
		path := fmt.Sprintf(`logs/reaction_adder/DMDGO-RA-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
		err := os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			utilities.LogErr("Error creating logs directory: %s", err)
			utilities.ExitSafely()
		}
		tokenFileX, err := os.Create(fmt.Sprintf(`%s/token.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating token file: %s", err)
			utilities.ExitSafely()
		}
		tokenFileX.Close()
		successFileX, err := os.Create(fmt.Sprintf(`%s/success.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating success file: %s", err)
			utilities.ExitSafely()
		}
		successFileX.Close()
		failedFileX, err := os.Create(fmt.Sprintf(`%s/failed.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating failed file: %s", err)
			utilities.ExitSafely()
		}
		failedFileX.Close()
		tokenFile, successFile, failedFile = tokenFileX.Name(), successFileX.Name(), failedFileX.Name()
		for i := 0; i < len(instances); i++ {
			instances[i].WriteInstanceToFile(tokenFile)
		}
	}
	var TotalCount, SuccessCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Success, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount-SuccessCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	var wg sync.WaitGroup
	wg.Add(len(instances))
	if choice == 1 {
		token := utilities.UserInput("Enter a token which can see the message: ")
		id := utilities.UserInput("Enter the message ID: ")
		channel := utilities.UserInput("Enter the channel ID: ")
		msg, err := instance.GetRxn(channel, id, token)
		if err != nil {
			utilities.LogErr("Error while getting message %s", err)
			return
		}
		var selection []string
		if len(msg.Reactions) == 0 {
			utilities.LogErr("Message has no reactions. React to the message to select it or use the manual option.")
			return
		}
		for i := 0; i < len(msg.Reactions); i++ {
			var j string
			if msg.Reactions[i].Emojis.ID == "" {
				j = msg.Reactions[i].Emojis.Name
			} else {
				j = fmt.Sprintf("<:%s:%s>", msg.Reactions[i].Emojis.Name, msg.Reactions[i].Emojis.ID)
			}
			selection = append(selection, fmt.Sprintf(`%v [%v Reacts]`, j, msg.Reactions[i].Count))
		}
		utilities.PrintMenu2(selection)

		var emojis []int
		x := utilities.UserInput("Select Emoji, seperate them by commas to select multiple:")
		if !strings.Contains(x, ",") {
			index, err := strconv.Atoi(x)
			if err != nil {
				utilities.LogErr("Error while converting %s to int %s", x, err)
				return
			}
			emojis = append(emojis, index)
		} else {
			for _, v := range strings.Split(x, ",") {
				index, err := strconv.Atoi(v)
				if err != nil {
					utilities.LogErr("Error while converting %s to int %s", v, err)
				} else {
					emojis = append(emojis, index)
				}

			}
		}
		var DelayBetweenReacts int
		var RandomDelayBetweenReacts int
		if len(emojis) > 1 {
			DelayBetweenReacts = utilities.UserInputInteger("Enter the delay between reacts for each token in seconds: ")
			RandomDelayBetweenReacts = utilities.UserInputInteger("Enter the random delay between reacts for each token in seconds: ")
		} else {
			DelayBetweenReacts = 0
		}
		TotalCount = len(instances) * len(emojis)
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < len(emojis); j++ {
					var send string
					if msg.Reactions[emojis[j]].Emojis.ID == "" {
						send = msg.Reactions[emojis[j]].Emojis.Name

					} else if msg.Reactions[emojis[j]].Emojis.ID != "" {
						send = msg.Reactions[emojis[j]].Emojis.Name + ":" + msg.Reactions[emojis[j]].Emojis.ID
					}
					err := instances[i].React(channel, id, send)
					if err != nil {
						utilities.LogFailed("Token %v failed to react to emoji %v %v", instances[i].CensorToken(), send, err)
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(failedFile, fmt.Sprintf(`Token %v Emoji %v Channel %v Message %v Error %v`, instances[i].CensorToken(), send, channel, id, err))
						}
						FailedCount++
					} else {
						utilities.LogSuccess("Token %v successfully reacted to emoji %v", instances[i].CensorToken(), send)
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(successFile, fmt.Sprintf(`Token %v Emoji %v Channel %v Message %v`, instances[i].CensorToken(), send, channel, id))
						}
						SuccessCount++
					}
					if DelayBetweenReacts != 0 {
						utilities.LogInfo("Token %v sleeping for %v seconds (Base Delay)", instances[i].CensorToken(), DelayBetweenReacts)
						time.Sleep(time.Duration(DelayBetweenReacts) * time.Second)
					}
					if RandomDelayBetweenReacts != 0 {
						x := rand.Intn(RandomDelayBetweenReacts)
						utilities.LogInfo("Token %v sleeping for %v seconds (Random Delay)", instances[i].CensorToken(), x)
						time.Sleep(time.Second * time.Duration(x))
					}
				}
			}(i)
		}
		wg.Wait()
		utilities.LogSuccess("Finished All threads")
	}
	if choice == 2 {
		id := utilities.UserInput("Enter the message ID: ")
		channel := utilities.UserInput("Enter the channel ID: ")
		var emojis []string
		x := utilities.UserInput("Enter the emoji, seperate them by commas to select multiple. Format is emojiName or emojiName:emojiID for nitro emojis: ")
		if !strings.Contains(x, ",") {
			emojis = append(emojis, x)
		} else {
			emojis = append(emojis, strings.Split(x, ",")...)
		}
		var DelayBetweenReacts int
		var RandomDelayBetweenReacts int
		if len(emojis) > 1 {
			DelayBetweenReacts = utilities.UserInputInteger("Enter the delay between reacts for each token in seconds: ")
			RandomDelayBetweenReacts = utilities.UserInputInteger("Enter the random delay between reacts for each token in seconds: ")
		} else {
			DelayBetweenReacts = 0
		}
		TotalCount = len(instances) * len(emojis)
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < len(emojis); j++ {
					err := instances[i].React(channel, id, emojis[j])
					if err != nil {
						utilities.LogFailed("Token %v failed to react to emoji %v %v", instances[i].CensorToken(), emojis[j], err)
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(failedFile, fmt.Sprintf(`Token %v Emoji %v Channel %v Message %v Error %v`, instances[i].CensorToken(), emojis[j], channel, id, err))
						}
						FailedCount++
					} else {
						utilities.LogSuccess("Token %v successfully reacted to emoji %v", instances[i].CensorToken(), emojis[j])
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(successFile, fmt.Sprintf(`Token %v Emoji %v Channel %v Message %v`, instances[i].CensorToken(), emojis[j], channel, id))
						}
						SuccessCount++
					}
					if DelayBetweenReacts != 0 {
						utilities.LogInfo("Token %v sleeping for %v seconds (Base Delay)", instances[i].CensorToken(), DelayBetweenReacts)
						time.Sleep(time.Duration(DelayBetweenReacts) * time.Second)
					}
					if RandomDelayBetweenReacts != 0 {
						x := rand.Intn(RandomDelayBetweenReacts)
						utilities.LogInfo("Token %v sleeping for %v seconds (Random Delay)", instances[i].CensorToken(), x)
						time.Sleep(time.Second * time.Duration(x))
					}
				}
			}(i)
		}
		wg.Wait()
		title <- true
		utilities.LogSuccess("Finished All threads")
	}

}
