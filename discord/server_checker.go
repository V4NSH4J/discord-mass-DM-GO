// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchServerChecker() {
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting neccessary information %v", err)
		return
	}
	var tokenFile, presentFile, notPresentFile string
	if cfg.OtherSettings.Logs {
		path := fmt.Sprintf(`logs/server_checker/DMDGO-SC-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
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
		presentFileX, err := os.Create(fmt.Sprintf(`%s/present.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating present file: %s", err)
			utilities.ExitSafely()
		}
		presentFileX.Close()
		notPresentFileX, err := os.Create(fmt.Sprintf(`%s/not_present.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating failed file: %s", err)
			utilities.ExitSafely()
		}
		notPresentFileX.Close()
		tokenFile, presentFile, notPresentFile = tokenFileX.Name(), presentFileX.Name(), notPresentFileX.Name()
		for i := 0; i < len(instances); i++ {
			instances[i].WriteInstanceToFile(tokenFile)
		}
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
	serverid = utilities.UserInput("Enter the server ID: ")
	var wg sync.WaitGroup
	wg.Add(len(instances))
	for i := 0; i < len(instances); i++ {
		go func(i int) {
			defer wg.Done()
			r, err := instances[i].ServerCheck(serverid)
			if err != nil {
				utilities.LogErr("%v Error while checking server: %v", instances[i].CensorToken(), err)
				if cfg.OtherSettings.Logs {
					instances[i].WriteInstanceToFile(notPresentFile)
				}
			} else {
				if r == 200 || r == 204 {
					utilities.LogSuccess("%v is in server %v ", instances[i].CensorToken(), serverid)
					inServer = append(inServer, instances[i].Token)
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(presentFile)
					}
				} else if r == 429 {
					utilities.LogFailed("%v is rate limited", instances[i].CensorToken())
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(notPresentFile)
					}
				} else if r == 400 {
					utilities.LogFailed("Bad request - Invalid Server ID")
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(notPresentFile)
					}
				} else {
					utilities.LogFailed("%v is not in server [%v]", instances[i].CensorToken(), serverid, r)
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(notPresentFile)
					}
				}
			}
		}(i)
	}
	wg.Wait()
	title <- true
	save := utilities.UserInput("Do you want to save the results? (y/n)")
	if save == "y" || save == "Y" {
		err := utilities.TruncateLines("tokens.txt", inServer)
		if err != nil {
			utilities.LogErr("Error while truncating file: %v", err)
		} else {
			utilities.LogSuccess("Successfully truncated file")
		}
	}
}
