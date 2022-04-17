// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"bufio"
	"os"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func LaunchTokenOnliner() {
	color.Blue("Token Onliner")
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("Error while getting necessary data %v", err)
		utilities.ExitSafely()
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
}
