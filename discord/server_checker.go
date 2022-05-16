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

func LaunchServerChecker() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
		utilities.ExitSafely()
	}
	var serverid string
	var inServer []string
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Present in Server]`, len(inServer)))
				_ = cmd.Run()
			}

		}
	}()
	color.Green("[%v] Enter server ID: ", time.Now().Format("15:04:05"))
	fmt.Scanln(&serverid)
	var wg sync.WaitGroup
	wg.Add(len(instances))
	for i := 0; i < len(instances); i++ {
		go func(i int) {
			defer wg.Done()
			r, err := instances[i].ServerCheck(serverid)
			if err != nil {
				color.Red("[%v] %v Error while checking server: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), err)
			} else {
				if r == 200 || r == 204 {
					color.Green("[%v] %v is in server %v ", time.Now().Format("15:04:05"), instances[i].CensorToken(), serverid)
					inServer = append(inServer, instances[i].Token)
				} else if r == 429 {
					color.Green("[%v] %v is rate limited", time.Now().Format("15:04:05"), instances[i].CensorToken())
				} else if r == 400 {
					color.Red("[%v] Bad request - Invalid Server ID", time.Now().Format("15:04:05"))
				} else {
					color.Red("[%v] %v is not in server [%v] [%v]", time.Now().Format("15:04:05"), instances[i].CensorToken(), serverid, r)
				}
			}
		}(i)
	}
	wg.Wait()
	title <- true
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
}
