// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func LaunchReactionAdder() {
	color.White("Use Offset (milliseconds) in config to use delay while reaction adding.")
	color.White("Menu:\n1) From message\n2) Manually")
	var choice int
	fmt.Scanln(&choice)
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		fmt.Println(err)
		utilities.ExitSafely()
	}
	var TotalCount, SuccessCount, FailedCount int 
	title := make(chan bool)
	go func() {
		Out:
		for {
			select {
			case<- title: 
				break Out
			default: 
			cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Success, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount - SuccessCount - FailedCount))
			_ = cmd.Run()
			}

		}
	}()
	TotalCount = len(instances)
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
		msg, err := instance.GetRxn(channel, id, token)
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
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
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
					color.Red("[%v] %v failed to react", time.Now().Format("15:04:05"), instances[i].CensorToken)
					FailedCount++
				} else {
					color.Green("[%v] %v reacted to the emoji", time.Now().Format("15:04:05"), instances[i].CensorToken)
					SuccessCount++
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
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
			go func(i int) {
				defer wg.Done()
				err := instances[i].React(channel, id, emoji)
				if err != nil {
					fmt.Println(err)
					color.Red("[%v] %v failed to react", time.Now().Format("15:04:05"), instances[i].CensorToken)
					FailedCount++
				}
				color.Green("[%v] %v reacted to the emoji", time.Now().Format("15:04:05"), instances[i].CensorToken)
				SuccessCount++
			}(i)
		}
		wg.Wait()
		title <- true 
		color.Green("[%v] Completed all threads.", time.Now().Format("15:04:05"))
	}

}
