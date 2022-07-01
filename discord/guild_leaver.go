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
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/zenthangplus/goccm"
)

func LaunchGuildLeaver() {
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting instances or config %s", err)
		return
	}
	var LeftCount, TotalCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Successfully Left, %v Failed, %v Unprocessed]`, LeftCount, FailedCount, TotalCount-LeftCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	var tokenFile, successFile, failedFile string
	if cfg.OtherSettings.Logs {
		path := fmt.Sprintf(`logs/guild_leaver/DMDGO-GL-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
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
	threads := utilities.UserInputInteger("Enter number of threads (0 for unlimited):")
	if threads > len(instances) {
		threads = len(instances)
	}
	if threads == 0 {
		threads = len(instances)
	}
	delay := utilities.UserInputInteger("Enter delay between leaves on each thread (in seconds):")
	serverid := utilities.UserInput("Enter server ID:")
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
		c.Wait()
		TotalCount++
		go func(i int) {
			p := instances[i].Leave(serverid)
			if p == 0 {
				utilities.LogFailed("Error while leaving on token %v", instances[i].CensorToken())
				if cfg.OtherSettings.Logs {
					instances[i].WriteInstanceToFile(failedFile)
				}
				FailedCount++
			}
			if p == 200 || p == 204 {
				utilities.LogSuccess("Successfully left on token %v", instances[i].CensorToken())
				if cfg.OtherSettings.Logs {
					instances[i].WriteInstanceToFile(successFile)
				}
				LeftCount++
			} else {
				utilities.LogFailed("Invalid Status code %v while leaving on token %v", p, instances[i].CensorToken())
				if cfg.OtherSettings.Logs {
					instances[i].WriteInstanceToFile(failedFile)
				}
				FailedCount++
			}
			time.Sleep(time.Duration(delay) * time.Second)
			c.Done()
		}(i)
	}
	c.WaitAllDone()
	title <- true
	utilities.LogSuccess("All Threads Completed!")
}
