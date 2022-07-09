// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/zenthangplus/goccm"
)

func LaunchTokenNuker() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting necessary data: %v", err)
		utilities.ExitSafely()
	}
	threads := utilities.UserInputInteger("Enter the number of threads: (0 for Maximum)")
	if threads > len(instances) {
		threads = len(instances)
	}
	if threads == 0 {
		threads = len(instances)
	}
	delay := utilities.UserInputInteger("Enter the delay between each request: (0 for No Delay)")
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			defer c.Done()
			// Getting all guilds
			func(i int) {
				respCode, _, guilds, err := instances[i].Guilds()
				if err != nil {
					utilities.LogFailed("Instance %v Error while getting guilds: %s", instances[i].Token, err)
					return
				}
				if respCode != 200 {
					utilities.LogFailed("Instance %v Invalid Status Code while getting guilds: %s", instances[i].Token, respCode)
					return
				}
				// Nuking all guilds
				for j := 0; j < len(guilds); j++ {
					p := instances[i].Leave(guilds[j])
					if p == 0 {
						utilities.LogFailed("Instance %v Error while leaving %v", instances[i].Token, guilds[j])
					}
					if p == 200 || p == 204 {
						utilities.LogSuccess("%v Left server %v", instances[i].CensorToken(), guilds[j])
						if delay > 0 {
							time.Sleep(time.Duration(delay) * time.Millisecond)
						}
					} else {
						utilities.LogFailed("%v Error while leaving - Invalid status code %v", instances[i].CensorToken(), p)
					}
				}
				utilities.LogSuccess("Instance %v no servers left", instances[i].CensorToken())
			}(i)

			// Getting all channels (DMs, GCs, etc)
			func(i int) {
				respCode, _, channels, err := instances[i].Channels()
				if err != nil {
					utilities.LogFailed("Instance %v Error while getting channels: %s", instances[i].Token, err)
					return
				}
				if respCode != 200 {
					utilities.LogFailed("Instance %v Invalid Status Code while getting channels: %s", instances[i].Token, respCode)
					return
				}
				// Nuking all channels
				for j := 0; j < len(channels); j++ {
					p, err := instances[i].CloseDMS(channels[j].ID)
					if p == -1 {
						utilities.LogFailed("Instance %v Error %v while closing %v", instances[i].Token, err, channels[j])
						continue
					}
					if p == 204 {
						if channels[j].Type == 1 {
							utilities.LogSuccess("%v Closed DM %v", instances[i].CensorToken(), channels[j])
						} else if channels[j].Type == 3 {
							utilities.LogSuccess("%v Closed Group DM %v", instances[i].CensorToken(), channels[j])
						} else {
							utilities.LogSuccess("%v Closed Channel %v", instances[i].CensorToken(), channels[j])
						}
						if delay > 0 {
							time.Sleep(time.Duration(delay) * time.Millisecond)
						}
					} else {
						utilities.LogFailed("%v Error while closing - Invalid status code %v", instances[i].CensorToken(), p)
					}
					utilities.LogSuccess("Instance %v no channels left", instances[i].CensorToken())
				}
			}(i)
			// Getting all Blocked, Pending, and Friends
			func(i int) {
				respCode, _, _, _, _, relations, err := instances[i].Relationships()
				if err != nil {
					utilities.LogFailed("Instance %v Error while getting relationships: %s", instances[i].Token, err)
					return
				}
				if respCode != 200 {
					utilities.LogFailed("Instance %v Invalid Status Code while getting relationships: %s", instances[i].Token, respCode)
					return
				}
				// Nuking all relationships
				for j := 0; j < len(relations); j++ {
					p, err := instances[i].EndRelation(relations[j].ID)
					if p == -1 {
						utilities.LogFailed("Instance %v Error %v while ending %v", instances[i].Token, err, relations[j])
						return
					}
					if p == 204 {
						utilities.LogSuccess("%v Ended relationship %v", instances[i].CensorToken(), relations[j])
						if delay > 0 {
							time.Sleep(time.Duration(delay) * time.Millisecond)
						}
					} else {
						utilities.LogFailed("%v Error while ending relationship - Invalid status code %v", instances[i].CensorToken(), p)
					}

				}
				utilities.LogSuccess("Instance %v no relationships left", instances[i].CensorToken())
			}(i)
			utilities.LogSuccess("Instance %v nuked", instances[i].CensorToken())
		}(i)
	}
	c.WaitAllDone()
}
