// If emoji info doesn't come with Emoji, then explore option for reactWith otherwise Explore both options if enabled. 
package discord

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func LaunchAutoReact() {
	color.Cyan("Make sure your tokens are present in the server you will be reacting!")
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("Error while getting necessary data %v", err)
		utilities.ExitSafely()
	}
	servers, channels, users, messages, emojis, subscribe := cfg.AutoReact.Servers, cfg.AutoReact.Channels, cfg.AutoReact.Users, cfg.AutoReact.Messages, cfg.AutoReact.Emojis, cfg.AutoReact.Subscribe
	if cfg.AutoReact.Observer == "" {
		color.Red("Set an Observer token to use Auto React")
		utilities.ExitSafely()
	}
	if cfg.AutoReact.Randomness > 100 {
		color.Red("Randomness cannot be more than 100")
		cfg.AutoReact.Randomness = 100
	} else if cfg.AutoReact.Randomness < 0 {
		color.Red("Randomness cannot be less than 0")
		cfg.AutoReact.Randomness = 0
	}

	color.Green("[%v][O] Initializing Observer token [%v]", time.Now().Format("15:04:05"), cfg.AutoReact.Observer)
	var observerInstance instance.Instance
	observerInstance.Token = cfg.AutoReact.Observer
	var proxy string
	var proxies []string
	if cfg.ProxySettings.ProxyFromFile {
		proxies, err = utilities.ReadLines("proxies.txt")
		if err != nil {
			color.Red("Error while reading proxies.txt: %v", err)
			utilities.ExitSafely()
		}
		if len(proxies) == 0 {
			cfg.ProxySettings.ProxyFromFile = false
			color.Red("[%v][!] No proxies found in proxies.txt - Disabling proxy usage", time.Now().Format("15:04:05"))
		}
	}
	if cfg.ProxySettings.ProxyFromFile {
		proxy = proxies[rand.Intn(len(proxies))]
		observerInstance.Proxy = proxy
		if !cfg.ProxySettings.GatewayProxy {
			proxy = ""
		}
		observerInstance.GatewayProxy = proxy
	}
	client, err := instance.InitClient(proxy, cfg)
	if err != nil {
		color.Red("[%v][!] Error while initializing client: %s Using Default client for listener token", time.Now().Format("15:04:05"), err)
		client = http.DefaultClient
	}
	observerInstance.Client = client
	observerInstance.Config = cfg
	color.Green("[%v][O] Successfully initialized Observer token [%v]", time.Now().Format("15:04:05"), observerInstance.CensorToken())
	// To make sure the observer doesn't observe reactions caused by our own instances.
	var blacklisted []string
	for i := 0; i < len(instances); i++ {
		if !strings.Contains(instances[i].Token, ".") {
			continue
		} else {
			b64id := strings.Split(instances[i].Token, ".")[0]
			// Decode this base64 string
			id, err := base64.StdEncoding.DecodeString(b64id)
			if err != nil {
				color.Red("[%v][!] Error while decoding base64 string: %v", time.Now().Format("15:04:05"), err)
				continue
			}
			blacklisted = append(blacklisted, string(id))
		}
	}
	ticker := make(chan bool)
	kill := make(chan bool)
	reactQueue := make(chan instance.ReactInfo, 9999)
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			ticker <- true
		}
	}()
	go func() {
	Listener:
		for {
			if observerInstance.Ws != nil {
				if observerInstance.Ws.Conn != nil {
					select {
					case <-ticker:
						// Closing websocket and re-opening it.
						if observerInstance.Ws.Conn != nil {
							observerInstance.Ws.Close()
						}
						if observerInstance.Ws != nil {
							observerInstance.Ws = nil
						}
						color.Yellow("[%v][O] Disconnected Observer token to reconnect", time.Now().Format("15:04:05"))
					case x := <-observerInstance.Ws.Reactions:
						if cfg.AutoReact.ReactAll {
							var event instance.Event
							err := json.Unmarshal(x, &event)
							if err != nil {
								color.Red("[%v][!] Error while unmarshalling event: %s", time.Now().Format("15:04:05"), err)
								continue Listener
							}
							color.Cyan("[%v][O] Event received: %v reacted [Guild:%v|Message:%v|Channel:%v|Emoji:%v]", time.Now().Format("15:04:05"), event.Data.UserID, event.Data.GuildId, event.Data.MessageID, event.Data.ChannelID, event.Data.Emoji.Name)
							if !utilities.Contains(messages, event.Data.MessageID) && !utilities.Contains(emojis, event.Data.Emoji.Name) && !utilities.Contains(emojis, fmt.Sprintf(`%s:%s`, event.Data.Emoji.Name, event.Data.Emoji.ID)) && !utilities.Contains(blacklisted, event.Data.UserID) && !utilities.Contains(users, event.Data.UserID) && !utilities.Contains(channels, event.Data.ChannelID) && !utilities.Contains(servers, fmt.Sprintf("%v", event.Data.GuildId)) {
								continue Listener
							}
							if event.Data.Emoji.ID != "" {
								event.Data.Emoji.Name += ":" + event.Data.Emoji.ID
							}
							//Approved
							go func() {
								reactQueue <- instance.ReactInfo{MessageID: event.Data.MessageID, ChannelID: event.Data.ChannelID, Emoji: event.Data.Emoji.Name}
							}()
							continue Listener
						} else {
							continue Listener
						}
					case y := <-observerInstance.Ws.Messages:
						if len(cfg.AutoReact.ReactWith) != 0 {
							var event instance.Event
							err := json.Unmarshal(y, &event)
							if err != nil {
								color.Red("[%v][!] Error while unmarshalling event: %s", time.Now().Format("15:04:05"), err)
								continue Listener
							}
							color.Cyan("[%v][O] Event received: %v messaged [Guild:%v|Channel:%v]", time.Now().Format("15:04:05"), event.Data.Author.ID, event.Data.GuildId, event.Data.ChannelID)
							if !utilities.Contains(messages, event.Data.MessageId) && !utilities.Contains(emojis, event.Data.Emoji.Name) && !utilities.Contains(emojis, fmt.Sprintf(`%s:%s`, event.Data.Emoji.Name, event.Data.Emoji.ID)) && !utilities.Contains(blacklisted, event.Data.UserID) && !utilities.Contains(users, event.Data.UserID) && !utilities.Contains(channels, event.Data.ChannelID) && !utilities.Contains(servers, event.Data.GuildID) {
								continue Listener
							}
							if event.Data.Emoji.ID != "" {
								event.Data.Emoji.Name += ":" + event.Data.Emoji.ID
							}
							go func() {
								reactQueue <- instance.ReactInfo{MessageID: event.Data.MessageId, ChannelID: event.Data.ChannelID, Emoji: event.Data.Emoji.Name}
							}()
						} else {
							continue Listener
						}
					case <-kill:
						break Listener
					}

				} else {
					err := observerInstance.StartWS()
					if err != nil {
						color.Red("[%v][!] Error while starting observer websocket: %s", time.Now().Format("15:04:05"), err)
						time.Sleep(10 * time.Second)
						continue Listener
					}
					if len(subscribe) != 0 && observerInstance.Ws != nil {

						for i := 0; i < len(subscribe); i++ {
							err := instance.Subscribe(observerInstance.Ws, subscribe[i])
							if err != nil {
								color.Red("[%v][!] Error while subscribing to server: %s", time.Now().Format("15:04:05"), err)
								time.Sleep(10 * time.Second)
								continue Listener
							}
						}

					}
					color.Yellow("[%v][O] Reconnected Observer token", time.Now().Format("15:04:05"))
					continue Listener
				}
			} else {
				// Opening Websocket
				err := observerInstance.StartWS()
				if err != nil {
					color.Red("[%v][!] Error while starting observer websocket: %s", time.Now().Format("15:04:05"), err)
					time.Sleep(10 * time.Second)
					continue Listener
				}
				if len(subscribe) != 0 && observerInstance.Ws != nil {

					for i := 0; i < len(subscribe); i++ {
						err := instance.Subscribe(observerInstance.Ws, subscribe[i])
						if err != nil {
							color.Red("[%v][!] Error while subscribing to server: %s", time.Now().Format("15:04:05"), err)
							time.Sleep(10 * time.Second)
							continue Listener
						}
					}

				}
				color.Yellow("[%v][O] Reconnected Observer token", time.Now().Format("15:04:05"))
				continue Listener
			}
		}
		color.Red("[%v][O] Permanently disconnected Observer token", time.Now().Format("15:04:05"))
	}()
	for i := 0; i < len(instances); i++ {
		instances[i].ReactChannel = make(chan instance.ReactInfo)
	}
	go func() {
		for {
			react := <-reactQueue
			for i := 0; i < len(instances); i++ {
				if instances[i].ReactChannel != nil {
					instances[i].ReactChannel <- react
				}
			}
			time.Sleep(time.Duration(cfg.AutoReact.Delay) * time.Millisecond)
		}
	}()
	var wg sync.WaitGroup
	for i := 0; i < len(instances); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				x := <-instances[i].ReactChannel
				func(i int) {
					if cfg.AutoReact.ReactAll {
						for i := 0; i < len(instances[i].Reacted); i++ {
							if instances[i].Reacted[i].MessageID == x.MessageID && instances[i].Reacted[i].Emoji == x.Emoji && instances[i].Reacted[i].ChannelID == x.ChannelID {
								return
							}
						}
						randomNum := rand.Intn(100) + 1
						if randomNum > cfg.AutoReact.Randomness {
							return
						}
						if x.Emoji == "" {
							return
						}
						err := instances[i].React(x.ChannelID, x.MessageID, x.Emoji)
						if err != nil {
							color.Red("[%v][!] Instance %v Error while reacting: %s", time.Now().Format("15:04:05"), instances[i].Token, err)
						} else {
							color.Green("[%v][X] Instance %v Reacted: %v", time.Now().Format("15:04:05"), instances[i].Token, x.Emoji)
							instances[i].Reacted = append(instances[i].Reacted, x)
							if cfg.AutoReact.IndividualDelay != 0 {
								time.Sleep(time.Duration(cfg.AutoReact.IndividualDelay) * time.Millisecond)
							}
						}
						return 
					}
				}(i)
				func(i int) {
					if len(cfg.AutoReact.ReactWith) != 0 {
						for i := 0; i < len(instances[i].Reacted); i++ {
							if instances[i].Reacted[i].MessageID == x.MessageID && instances[i].Reacted[i].Emoji == x.Emoji && instances[i].Reacted[i].ChannelID == x.ChannelID {
								return
							}
						}
					Emoji:
						for i := 0; i < len(cfg.AutoReact.ReactWith); i++ {
							randomNum := rand.Intn(100) + 1
							if randomNum > cfg.AutoReact.Randomness {
								continue Emoji
							}
							if cfg.AutoReact.ReactWith[i] == "" {
								continue Emoji
							}
							err := instances[i].React(x.ChannelID, x.MessageID, cfg.AutoReact.ReactWith[i])
							if err != nil {
								color.Red("[%v][!] Instance %v Error while reacting: %s", time.Now().Format("15:04:05"), instances[i].Token, err)
							} else {
								color.Green("[%v][X] Instance %v Reacted: %v", time.Now().Format("15:04:05"), instances[i].Token, cfg.AutoReact.ReactWith[i])
								instances[i].Reacted = append(instances[i].Reacted, x)
								if cfg.AutoReact.IndividualDelay != 0 {
									time.Sleep(time.Duration(cfg.AutoReact.IndividualDelay) * time.Millisecond)
								}
							}
							return 
						}
					}
				}(i)
				continue 

			}
		}(i)
	}
	wg.Wait()

}
