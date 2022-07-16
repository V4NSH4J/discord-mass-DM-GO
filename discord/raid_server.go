package discord

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchRaidServer() {
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances%s", err)
		return
	}
	server := utilities.GetConfigOrInputString(
		cfg.RaidSettings.AttackedServer,
		"Enter server ID to perform the raid on: ",
	)
	channel := utilities.GetConfigOrInputString(
		cfg.RaidSettings.AttackedChannel,
		"Enter the channel ID to attack: ",
	)
	// USE MESSAGES JSON
	var msg instance.Message
	messagechoice := utilities.UserInputInteger("Enter 1 to use message from file (message.json), 2 to use message from console: ")
	if messagechoice != 1 && messagechoice != 2 {
		utilities.LogErr("Invalid choice")
		return
	}
	if messagechoice == 2 {
		text := utilities.UserInput("Enter your message, use \\n for changing lines. You can also set a constant message in message.json")
		msg.Content = text
		msg.Content = strings.Replace(msg.Content, "\\n", "\n", -1)
		var msgs []instance.Message
		msgs = append(msgs, msg)
		err := instance.SetMessages(instances, msgs)
		if err != nil {
			utilities.LogErr("Error while setting messages: %s", err)
			return
		}
	} else {
		var msgs []instance.Message
		err := instance.SetMessages(instances, msgs)
		if err != nil {
			utilities.LogErr("Error while setting messages: %s", err)
			return
		}
	}

	var sentMessages []int

	// START CONCURRENCY
	var wg sync.WaitGroup
	wg.Add(len(instances))

	for i := 0; i < len(instances); i++ {
		// Don't sleep the first one
		if i != 0 {
			time.Sleep(time.Duration(cfg.RaidSettings.Offset) * time.Millisecond)
		}
		go func(i int) {
			defer wg.Done()

			index := rand.Intn(len(instances[i].Messages))

			if cfg.RaidSettings.DuplicateMessage == false {
				if len(sentMessages) < len(instances[i].Messages) {
					for utilities.ContainsInt(sentMessages, index) {
						index = rand.Intn(len(instances[i].Messages))
					}
				} else {
					utilities.LogWarn("Token %v will send a duplicate message.", instances[i].Token)
				}
			}

			if !utilities.ContainsInt(sentMessages, index) {
				sentMessages = append(sentMessages, index)
			}

			if cfg.RaidSettings.OnlineTokens && instances[i].Ws == nil {
				err := instances[i].StartWS()
				if err != nil {
					utilities.LogFailed("Token %v error while going online: %v", instances[i].Token, err)
				}
			}

			respCode, body, err := instances[i].SendMessageToChannel(index, utilities.SnowflakeParams{ChannelId: channel, ServerId: server})
			if err != nil {
				utilities.LogErr("Token %v error while sending message %s", instances[i].Token, err)
			}
			if respCode == 200 {
				utilities.LogSuccess("Token %v raided channel: %v", instances[i].Token, channel)
			} else {
				utilities.LogFailed("Token %v failed to raid channel: %v [%v]", instances[i].Token, channel, string(body))
			}
			if cfg.RaidSettings.RandomIndividualDelay != 0 {
				var sleepTime = rand.Intn(cfg.RaidSettings.RandomIndividualDelay)
				//utilities.LogInfo("Next token will type in %v milliseconds", sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Millisecond)
			}
		}(i)
	}
	wg.Wait()

	if cfg.RaidSettings.OnlineTokens {
		for i := 0; i < len(instances); i++ {
			if instances[i].Ws != nil {
				instances[i].Ws.Close()
			}
		}
	}
	utilities.LogSuccess("All threads finished")
}
