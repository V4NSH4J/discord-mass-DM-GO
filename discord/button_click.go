// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html
package discord

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/zenthangplus/goccm"
)

func LaunchButtonClicker() {
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting instances or config %s", err)
		return
	}
	var tokenFile, successFile, failedFile string
	if cfg.OtherSettings.Logs {
		path := fmt.Sprintf(`logs/button_clicker/DMDGO-BC-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
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
	token := utilities.UserInput("Enter a token which can see the message:")
	id := utilities.UserInput("Enter the ID of the message:")
	channel := utilities.UserInput("Enter the ID of the channel:")
	server := utilities.UserInput("Enter the ID of the server:")
	msg, err := instance.FindMessage(channel, id, token)
	if err != nil {
		utilities.LogErr("Error while finding message: %v", err)
		return
	}
	utilities.LogInfo("Message found!\n %s", msg)
	var Msg instance.Message
	err = json.Unmarshal([]byte(msg), &Msg)
	if err != nil {
		utilities.LogErr("Error while unmarshalling message: %v", err)
		return
	}
	if len(Msg.Components) == 0 {
		utilities.LogErr("Message has no components (Buttons or similar)")
		return
	}
	for i := 0; i < len(Msg.Components); i++ {
		fmt.Printf("%v) Row %v", i, i)
	}
	row := utilities.UserInputInteger("Enter Row number:")
	for i := 0; i < len(Msg.Components[row].Buttons); i++ {
		if Msg.Components[row].Buttons[i].Label != "" {
			fmt.Printf("%v) Button %v [%v]", i, i, Msg.Components[row].Buttons[i].Label)
		} else if Msg.Components[row].Buttons[i].Emoji.Name != "" {
			fmt.Printf("%v) Button %v [%v]", i, i, Msg.Components[row].Buttons[i].Emoji)
		} else {
			fmt.Printf("%v) Button %v [Name or Emoji not found]", i, i)
		}
	}
	column := utilities.UserInputInteger("Select Button:")
	threads := utilities.UserInputInteger("Enter number of threads:")
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			defer c.Done()
			err := instances[i].StartWS()
			if err != nil {
				utilities.LogFailed("Error while starting websocket: %v", err)
			} else {
				utilities.LogSuccess("Websocket opened %s", instances[i].CensorToken())
			}
			respCode, err := instances[i].PressButton(row, column, server, Msg)
			if err != nil {
				utilities.LogFailed("Error while pressing button: %v", err)
				if cfg.OtherSettings.Logs {
					instances[i].WriteInstanceToFile(failedFile)
				}
				return
			}
			if respCode != 204 && respCode != 200 {
				utilities.LogFailed("Error while pressing button: %v", respCode)
				if cfg.OtherSettings.Logs {
					instances[i].WriteInstanceToFile(failedFile)
				}

				return
			}
			utilities.LogSuccess("Button pressed on instance %v", instances[i].CensorToken())
			if cfg.OtherSettings.Logs {
				instances[i].WriteInstanceToFile(successFile)
			}
			if instances[i].Ws != nil {
				if instances[i].Ws.Conn != nil {
					err = instances[i].Ws.Close()
					if err != nil {
						utilities.LogFailed("Error while closing websocket: %v", err)
					} else {
						utilities.LogSuccess("Websocket closed %v", instances[i].CensorToken())
					}
				}
			}
		}(i)
	}
	c.WaitAllDone()

}
