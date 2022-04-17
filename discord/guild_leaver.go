// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

func LaunchGuildLeaver() {
	color.Cyan("Guild Leaver")
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("Error while getting necessary data %v", err)
		utilities.ExitSafely()

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
		time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
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
}
