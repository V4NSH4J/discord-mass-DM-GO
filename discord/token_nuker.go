// Remove all friends and blocks
// Cancel all incoming friend requests

package discord

import (
	"fmt"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

func LaunchTokenNuker() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("Error: %s", err)
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
	color.White("Enter Delay between actions in milliseconds: ")
	var delay int
	fmt.Scanln(&delay)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			defer c.Done()
			// Getting all guilds
			func(i int) {
				respCode, _, guilds, err := instances[i].Guilds()
				if err != nil {
					color.Red("[%v] Instance %v Error while getting guilds: %s", time.Now().Format("15:04:05"), instances[i].Token, err)
					return
				}
				if respCode != 200 {
					color.Red("[%v] Instance %v Invalid Status Code while getting guilds: %s", time.Now().Format("15:04:05"), instances[i].Token, respCode)
					return
				}
				// Nuking all guilds
				for j := 0; j < len(guilds); j++ {
					p := instances[i].Leave(guilds[j])
					if p == 0 {
						color.Red("[%v] Instance %v Error while leaving %v", time.Now().Format("15:04:05"), instances[i].Token, guilds[j])
					}
					if p == 200 || p == 204 {
						color.Green("[%v] %v Left server %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), guilds[j])
						if delay > 0 {
							time.Sleep(time.Duration(delay) * time.Millisecond)
						}
					} else {
						color.Red("[%v] %v Error while leaving - Invalid status code %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), p)
					}
				}
				color.Green("[%v] Instance %v no servers left", time.Now().Format("15:04:05"), instances[i].CensorToken())
			}(i)

			// Getting all channels (DMs, GCs, etc)
			func(i int) {
				respCode, _, channels, err := instances[i].Channels()
				if err != nil {
					color.Red("[%v] Instance %v Error while getting channels: %s", time.Now().Format("15:04:05"), instances[i].Token, err)
					return
				}
				if respCode != 200 {
					color.Red("[%v] Instance %v Invalid Status Code while getting channels: %s", time.Now().Format("15:04:05"), instances[i].Token, respCode)
					return
				}
				// Nuking all channels
				for j := 0; j < len(channels); j++ {
					p, err := instances[i].CloseDMS(channels[j].ID)
					if p == -1 {
						color.Red("[%v] Instance %v Error %v while closing %v", time.Now().Format("15:04:05"), instances[i].Token, err, channels[j])
						continue
					}
					if p == 204 {
						if channels[j].Type == 1 {
							color.Green("[%v] %v Closed DM %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), channels[j])
						} else if channels[j].Type == 3 {
							color.Green("[%v] %v Closed Group DM %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), channels[j])
						} else {
							color.Green("[%v] %v Closed Channel %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), channels[j])
						}
						if delay > 0 {
							time.Sleep(time.Duration(delay) * time.Millisecond)
						}
					} else {
						color.Red("[%v] %v Error while closing - Invalid status code %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), p)
					}
					color.Green("[%v] Instance %v no channels left", time.Now().Format("15:04:05"), instances[i].CensorToken())
				}
			}(i)
			// Getting all Blocked, Pending, and Friends
			func(i int) {
				respCode, _, _, _, _, relations, err := instances[i].Relationships()
				if err != nil {
					color.Red("[%v] Instance %v Error while getting relationships: %s", time.Now().Format("15:04:05"), instances[i].Token, err)
					return
				}
				if respCode != 200 {
					color.Red("[%v] Instance %v Invalid Status Code while getting relationships: %s", time.Now().Format("15:04:05"), instances[i].Token, respCode)
					return
				}
				// Nuking all relationships
				for j := 0; j < len(relations); j++ {
					p, err := instances[i].EndRelation(relations[j].ID)
					if p == -1 {
						color.Red("[%v] Instance %v Error %v while ending %v", time.Now().Format("15:04:05"), instances[i].Token, err, relations[j])
						return
					}
					if p == 204 {
						color.Green("[%v] %v Ended relationship %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), relations[j])
						if delay > 0 {
							time.Sleep(time.Duration(delay) * time.Millisecond)
						}
					} else {
						color.Red("[%v] %v Error while ending relationship - Invalid status code %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), p)
					}

				}
				color.Green("[%v] Instance %v no relationships left", time.Now().Format("15:04:05"), instances[i].CensorToken())
			}(i)
			color.Green("[%v] Instance %v nuked", time.Now().Format("15:04:05"), instances[i].CensorToken())
		}(i)
	}
	c.WaitAllDone()
}
