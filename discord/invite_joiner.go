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
	"strings"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/zenthangplus/goccm"
)

func LaunchinviteJoiner() {
	utilities.PrintMenu([]string{"Single Invite", "Multiple Invites from File"})
	invitechoice := utilities.UserInputInteger("Enter your choice:")
	if invitechoice != 1 && invitechoice != 2 {
		utilities.LogErr("Invalid choice")
		return
	}
	switch invitechoice {
	case 1:
		cfg, instances, err := instance.GetEverything()
		if err != nil {
			utilities.LogErr("Error while getting config or instances %s", err)
		}
		var tokenFile, jointFile, failedFile, reactedFile string
		if cfg.OtherSettings.Logs {
			path := fmt.Sprintf(`logs/invite_joiner/DMDGO-IJ-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
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
			jointFileX, err := os.Create(fmt.Sprintf(`%s/joint.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating joint file: %s", err)
				utilities.ExitSafely()
			}
			jointFileX.Close()
			failedFileX, err := os.Create(fmt.Sprintf(`%s/failed.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating failed file: %s", err)
				utilities.ExitSafely()
			}
			failedFileX.Close()
			reactedFileX, err := os.Create(fmt.Sprintf(`%s/reacted.txt`, path))
			if err != nil {
				utilities.LogErr("Error creating reacted file: %s", err)
				utilities.ExitSafely()
			}
			tokenFile, jointFile, failedFile, reactedFile = tokenFileX.Name(), jointFileX.Name(), failedFileX.Name(), reactedFileX.Name()
			for i := 0; i < len(instances); i++ {
				instances[i].WriteInstanceToFile(tokenFile)
			}
		}

		invite := utilities.UserInput("Enter your Invite Code or Link:")
		invite = processInvite(invite)
		threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")

		if threads > len(instances) {
			threads = len(instances)
		}
		if threads == 0 {
			threads = len(instances)
		}
		verif := utilities.UserInputInteger("Use additional adding reaction verification passing. 0) No 1) Yes")
		var channelid string
		var msgid string
		var emoji string

		if verif == 1 {
			channelid = utilities.UserInput("ID of the channel with verification message")

			msgid = utilities.UserInput("ID of the message with verification reaction")
			emoji = utilities.UserInput("Enter emoji")

		}
		base := utilities.UserInputInteger("Enter base delay per thread for joining in seconds: ")
		random := utilities.UserInputInteger("Enter random delay per thread for joining in seconds: ")
		var delay int
		if random > 0 {
			delay = base + rand.Intn(random)
		} else {
			delay = base
		}
		c := goccm.New(threads)
		for i := 0; i < len(instances); i++ {
			c.Wait()
			go func(i int) {
				err := instances[i].Invite(invite)
				if err != nil {
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedFile, instances[i].Token)
					}
				} else {
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(jointFile, instances[i].Token)
					}
				}
				if verif == 1 {
					err := instances[i].React(channelid, msgid, emoji)
					if err != nil {
						utilities.LogFailed("%v failed to react to %v", instances[i].CensorToken(), emoji)
					} else {
						utilities.LogSuccess("%v reacted to the emoji %v", instances[i].CensorToken(), emoji)
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(reactedFile, instances[i].Token)
						}
					}

				}
				time.Sleep(time.Duration(delay) * time.Second)
				c.Done()

			}(i)
		}
		c.WaitAllDone()
		utilities.LogSuccess("All Threads Completed!")

	case 2:
		cfg, instances, err := instance.GetEverything()
		if err != nil {
			utilities.LogErr("Error while getting config or instances %s", err)
		}
		var tokenFile string
		path := fmt.Sprintf(`logs/multi_joiner/DMDGO-MJ-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
		if cfg.OtherSettings.Logs {
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
			tokenFile = tokenFileX.Name()
			for i := 0; i < len(instances); i++ {
				instances[i].WriteInstanceToFile(tokenFile)
			}
		}
		invites, err := utilities.ReadLines("invite.txt")
		if err != nil {
			utilities.LogErr("Error while opening invite.txt file %s", err)
			return
		}
		var inviteFiles []string
		if cfg.OtherSettings.Logs {
			for i := 0; i < len(invites); i++ {
				f, err := os.Create(fmt.Sprintf(`%s/%s.txt`, path, processInvite(invites[i])))
				if err != nil {
					utilities.LogErr("Error creating invite file %v: %s", invites[i], err)
				}
				inviteFiles = append(inviteFiles, f.Name())
			}
		}

		if len(invites) == 0 {
			utilities.LogErr("No invites found in invite.txt")
			return
		}
		delay := utilities.UserInputInteger("Enter delay between 2 consecutive joins by 1 token in seconds: ")
		threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
		if threads > len(instances) {
			threads = len(instances)
		}
		if threads == 0 {
			threads = len(instances)
		}
		c := goccm.New(threads)
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
			c.Wait()
			go func(i int) {
				for j := 0; j < len(invites); j++ {
					err := instances[i].Invite(processInvite(invites[j]))
					if err == nil {
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(inviteFiles[j], instances[i].Token)
						}
					}
					time.Sleep(time.Duration(delay) * time.Second)
				}
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		utilities.LogSuccess("All Threads Completed!")
	}
}

func processInvite(rawInvite string) string {
	if !strings.Contains(rawInvite, "/") {
		return rawInvite
	} else {
		return strings.Split(rawInvite, "/")[len(strings.Split(rawInvite, "/"))-1]
	}
}
