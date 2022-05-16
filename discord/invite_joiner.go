// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

func LaunchinviteJoiner() {
	var invitechoice int
	color.White("Invite Menu:\n1) Single Invite\n2) Multiple Invites from file")
	fmt.Scanln(&invitechoice)
	if invitechoice != 1 && invitechoice != 2 {
		color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
		utilities.ExitSafely()
		return
	}
	switch invitechoice {
	case 1:
		_, instances, err := instance.GetEverything()
		if err != nil {
			color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
		}

		color.White("[%v] Enter your invite code [Only the CODE or the Link: ", time.Now().Format("15:04:05"))
		var invite string
		fmt.Scanln(&invite)
		invite = processInvite(invite)
		color.White("[%v] Enter number of Threads (0: Unlimited Threads. 1: For using proper delay. It may be a good idea to use less threads if you're looking to solve captchas): ", time.Now().Format("15:04:05"))
		var threads int
		fmt.Scanln(&threads)

		if threads > len(instances) {
			threads = len(instances)
		}
		if threads == 0 {
			threads = len(instances)
		}

		color.White("[%v] Enter base delay for joining in seconds (0 for none)", time.Now().Format("15:04:05"))
		var base int
		fmt.Scanln(&base)
		color.White("[%v] Enter random delay to be added upon base delay (0 for none)", time.Now().Format("15:04:05"))
		var random int
		fmt.Scanln(&random)
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
					color.Red("[%v] Error while joining: %v", time.Now().Format("15:04:05"), err)
				}
				time.Sleep(time.Duration(delay) * time.Second)
				c.Done()

			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] All threads finished", time.Now().Format("15:04:05"))

	case 2:
		color.Cyan("Multiple Invite Mode")
		color.White("This will join your tokens from tokens.txt to servers from invite.txt")
		cfg, instances, err := instance.GetEverything()
		if err != nil {
			color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
		}

		if len(instances) == 0 {
			color.Red("[%v] Enter your tokens in tokens.txt", time.Now().Format("15:04:05"))
			utilities.ExitSafely()
		}
		invites, err := utilities.ReadLines("invite.txt")
		if err != nil {
			color.Red("Error while opening invite.txt: %v", err)
			utilities.ExitSafely()
			return
		}
		if len(invites) == 0 {
			color.Red("[%v] Enter your invites in invite.txt", time.Now().Format("15:04:05"))
			utilities.ExitSafely()
			return
		}
		color.White("Enter delay between 2 consecutive joins by 1 token in seconds: ")
		var delay int
		fmt.Scanln(&delay)
		color.White("Enter number of Threads (0 for unlimited): ")
		var threads int
		fmt.Scanln(&threads)
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
					if err != nil {
						color.Red("[%v] Error while joining: %v", time.Now().Format("15:04:05"), err)
					}
					time.Sleep(time.Duration(delay) * time.Second)
				}
				c.Done()
			}(i)
		}
		c.WaitAllDone()
		color.Green("[%v] All threads finished", time.Now().Format("15:04:05"))
	}
}

func processInvite(rawInvite string) string {
	if !strings.Contains(rawInvite, "/") {
		return rawInvite
	} else {
		return strings.Split(rawInvite, "/")[len(strings.Split(rawInvite, "/"))-1]
	}
}
