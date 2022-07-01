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

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchTokenOnliner() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting neccessary information %v", err)
		utilities.ExitSafely()
	}
	var wg sync.WaitGroup
	wg.Add(len(instances))
	for i := 0; i < len(instances); i++ {
		go func(i int) {
			err := instances[i].StartWS()
			if err != nil {
				utilities.LogErr("Token %v Error while starting websocket %v", instances[i].CensorToken(), err)
			} else {
				utilities.LogSuccess("Websocket opened %v", instances[i].CensorToken())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	utilities.LogInfo("All Token online. Press ENTER to disconnect and continue the program")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	wg.Add(len(instances))
	for i := 0; i < len(instances); i++ {
		go func(i int) {
			instances[i].Ws.Close()
			wg.Done()
		}(i)
	}
	wg.Wait()
	utilities.LogInfo("All Token offline")
}
