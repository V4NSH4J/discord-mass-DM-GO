// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

func LaunchTokenChecker() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
		utilities.ExitSafely()
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
	title := make(chan bool)
	var valid, invalid int 
	go func() {
		Out:
		for {
			select {
			case<- title: 
				break Out
			default: 
			cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Unchecked %v Valid %v Invalid]`, len(instances)-valid-invalid, valid, invalid))
			_ = cmd.Run()
			}

		}
	}()
	c := goccm.New(threads)
	var working []instance.Instance
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			err := instances[i].CheckToken()
			if err != 200 {
				color.Red("[%v] Token Invalid %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
				invalid++
			} else {
				color.Green("[%v] Token Valid %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
				working = append(working, instances[i])
				valid++
			}
			c.Done()
		}(i)
	}
	c.WaitAllDone()
	var workingTokens []string
	for i := 0; i < len(working); i++ {
		if working[i].Password != "" && working[i].Email != "" {
			workingTokens = append(workingTokens, fmt.Sprintf(`%v:%v:%v`, working[i].Email, working[i].Password, working[i].Token))
		} else {
			workingTokens = append(workingTokens, fmt.Sprintf(`%v`, working[i].Token))
		}
	}
	t := utilities.TruncateLines("tokens.txt", workingTokens)
	if t != nil {
		color.Red("[%v] Error while truncating tokens.txt: %v", time.Now().Format("15:04:05"), t)
		utilities.ExitSafely()
		return
	}
	title <- true 
	color.Green("[%v] All threads finished", time.Now().Format("15:04:05"))
}


