// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
)

func LaunchDMReact() {
	color.Cyan("DM On React")
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("Error while obtaining config and instances: %s", err)
	}
	// Checking config for observer token
	if cfg.DMonReact.Observer == "" {
		color.Red("Set an Observer token to use DM on react")
		utilities.ExitSafely()
	}
	if cfg.DMonReact.ServerID == "" {
		color.Red("Set a Server ID to use DM on react")
		utilities.ExitSafely()
	}
	if cfg.DMonReact.Invite == "" {
		color.Red("Set an Invite to use DM on react")
		utilities.ExitSafely()
	}
	var msg instance.Message
	color.Green("[%v] Press 1 to use messages from file or press 2 to enter a message: ", time.Now().Format("15:04:05"))
	var messagechoice int
	fmt.Scanln(&messagechoice)
	if messagechoice != 1 && messagechoice != 2 {
		color.Red("[%v] Invalid choice", time.Now().Format("15:04:05"))
		utilities.ExitSafely()
	}
	if messagechoice == 2 {
		color.Green("[%v] Enter your message, use \\n for changing lines. You can also set a constant message in message.json", time.Now().Format("15:04:05"))
		scanner := bufio.NewScanner(os.Stdin)
		var text string
		if scanner.Scan() {
			text = scanner.Text()
		}

		msg.Content = text
		msg.Content = strings.Replace(msg.Content, "\\n", "\n", -1)
		var msgs []instance.Message
		msgs = append(msgs, msg)
		err := instance.SetMessages(instances, msgs)
		if err != nil {
			color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
			utilities.ExitSafely()
		}
	} else {
		var msgs []instance.Message
		err := instance.SetMessages(instances, msgs)
		if err != nil {
			color.Red("[%v] Error while setting messages: %v", time.Now().Format("15:04:05"), err)
			utilities.ExitSafely()
		}
	}
	// Initializing Files & Variables
	var completed []string
	var failed []string
	var names []string
	var avatars []string
	var proxies []string
	if cfg.DMonReact.SkipCompleted {
		completed, err = utilities.ReadLines("completed.txt")
		if err != nil {
			color.Red("Error while reading completed.txt: %s", err)
			utilities.ExitSafely()
		}
	}
	if cfg.DMonReact.SkipFailed {
		failed, err = utilities.ReadLines("failed.txt")
		if err != nil {
			color.Red("Error while reading failed.txt: %s", err)
			utilities.ExitSafely()
		}
	}
	if cfg.DMonReact.ChangeAvatar {
		color.Green("[%v] Loading Avatars..", time.Now().Format("15:04:05"))
		ex, err := os.Executable()
		if err != nil {
			color.Red("Couldn't find Exe")
			utilities.ExitSafely()
		}
		ex = filepath.ToSlash(ex)
		path := path.Join(path.Dir(ex) + "/input/pfps")

		images, err := instance.GetFiles(path)
		if err != nil {
			color.Red("Couldn't find images in PFPs folder")
			utilities.ExitSafely()
		}
		color.Green("%v files found", len(images))
		if len(images) == 0 {
			cfg.DMonReact.ChangeAvatar = false
			color.Red("[%v][!] No images found in PFPs folder - Disabling avatar change", time.Now().Format("15:04:05"))
		} else {
			for i := 0; i < len(images); i++ {
				av, err := instance.EncodeImg(images[i])
				if err != nil {
					color.Red("Couldn't encode image")
					continue
				}
				avatars = append(avatars, av)
			}
			color.Green("%v avatars loaded", len(avatars))
		}
	}
	if cfg.DMonReact.ChangeName {
		names, err = utilities.ReadLines("names.txt")
		if err != nil {
			color.Red("Error while reading names.txt: %s", err)
			utilities.ExitSafely()
		}
		if len(names) == 0 {
			cfg.DMonReact.ChangeName = false
			color.Red("[%v][!] No names found in names.txt - Disabling name change", time.Now().Format("15:04:05"))
		}
	}
	if cfg.DMonReact.ChangeName {
		for i := 0; i < len(instances); i++ {
			if instances[i].Password == "" {
				cfg.DMonReact.ChangeName = false
				color.Red("[%v][!] Token %v has no password, perhaps using the wrong format. Need to use email:password:token to use the name changer", time.Now().Format("15:04:05"), instances[i].Token)
				break
			}
		}
	}
	if cfg.ProxySettings.ProxyFromFile {
		proxies, err = utilities.ReadLines("proxies.txt")
		if len(proxies) == 0 {
			cfg.ProxySettings.ProxyFromFile = false
			color.Red("[%v][!] No proxies found in proxies.txt - Disabling proxy usage", time.Now().Format("15:04:05"))
		}
	}
	if cfg.CaptchaSettings.ClientKey == "" {
		color.Red("[%v][!] You're not using a Captcha key, if anytime the token is met with a captcha, it will be switched.")
	}
	tokenPool := make(chan instance.Instance, len(instances))
	for i := 0; i < len(instances); i++ {
		go func(i int) {
			if instances[i].Token != cfg.DMonReact.Observer {
				tokenPool <- instances[i]
			} else {
				color.Red("[%v][!] Skipping Observer token %v", time.Now().Format("15:04:05"), instances[i].Token)
			}
		}(i)
	}

	// All files and variables loaded and errors handled.
	color.Green("[%v][O] Initializing Observer token [%v]", time.Now().Format("15:04:05"), cfg.DMonReact.Observer)
	var observerInstance instance.Instance
	observerInstance.Token = cfg.DMonReact.Observer
	var proxy string
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
	if cfg.DMonReact.ServerID != "" {
		r, err := observerInstance.ServerCheck(cfg.DMonReact.ServerID)
		if err != nil {
			color.Red("[%v][!] Error while checking server: %s", time.Now().Format("15:04:05"), err)
		} else {
			if r != 200 && r != 204 {
				// Token not in server or some other issue like rate limit
				err := observerInstance.Invite(cfg.DMonReact.Invite)
				if err != nil {
					color.Red("[%v][!] Error while inviting to server: %s", time.Now().Format("15:04:05"), err)
				} else {
					color.Green("[%v][O] Successfully invited to server", time.Now().Format("15:04:05"))
				}

			}
		}
		r, err = observerInstance.ServerCheck(cfg.DMonReact.ServerID)
		if err != nil {
			color.Red("[%v][!] Error while checking server: %s", time.Now().Format("15:04:05"), err)
			utilities.ExitSafely()
		} else {
			if r != 200 && r != 204 {
				color.Red("[%v][!] Error while joining server: %v - If you meant to listen to reactions globally, make serverid an empty string", time.Now().Format("15:04:05"), r)
				// Token was tried to be invited after it was found that it might not be in the server. And after checking again the listener token may still not be present
				utilities.ExitSafely()
			}
		}

	}
	color.Green("[%v][O] Successfully initialized Observer token [%v]", time.Now().Format("15:04:05"), cfg.DMonReact.Observer)
	// Start Listening for reactions.
	ticker := make(chan bool)
	kill := make(chan bool)
	filteredReacts := make(chan string, 10000)

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
						var event instance.Event
						err := json.Unmarshal(x, &event)
						if err != nil {
							color.Red("[%v][!] Error while unmarshalling event: %s", time.Now().Format("15:04:05"), err)
							continue Listener
						}
						color.Cyan("[%v][O] Event received: %v reacted [%v|%v|%v|%v]", time.Now().Format("15:04:05"), event.Data.UserID, event.Data.GuildId, event.Data.MessageID, event.Data.ChannelID, event.Data.Emoji.Name)
						if cfg.DMonReact.MaxAntiRaidQueue > 0 {
							if len(filteredReacts) >= cfg.DMonReact.MaxAntiRaidQueue {
								color.Red("[%v][!] Anti-Raid queue is full, skipping this reaction [%v]", time.Now().Format("15:04:05"), event.Data.UserID)
								continue Listener
							}
						}
						if cfg.DMonReact.SkipCompleted {
							if utilities.Contains(completed, event.Data.UserID) {
								color.Cyan("[%v][O] Skipping completed user [%v]", time.Now().Format("15:04:05"), event.Data.UserID)
								continue Listener
							}
						}
						if cfg.DMonReact.SkipFailed {
							if utilities.Contains(failed, event.Data.UserID) {
								color.Cyan("[%v][O] Skipping failed user [%v]", time.Now().Format("15:04:05"), event.Data.UserID)
								continue Listener
							}
						}
						if cfg.DMonReact.ServerID != "" {
							if event.Data.GuildId != cfg.DMonReact.ServerID {
								color.Cyan("[%v][O] Skipping event from other server [%v]", time.Now().Format("15:04:05"), event.Data.GuildId)
								continue Listener
							}
						}
						if cfg.DMonReact.ChannelID != "" {
							if event.Data.ChannelID != cfg.DMonReact.ChannelID {
								color.Cyan("[%v][O] Skipping event from other channel [%v]", time.Now().Format("15:04:05"), event.Data.ChannelID)
								continue Listener
							}
						}
						if cfg.DMonReact.MessageID != "" {
							if event.Data.MessageID != cfg.DMonReact.MessageID {
								color.Cyan("[%v][O] Skipping event from other message [%v]", time.Now().Format("15:04:05"), event.Data.MessageID)
								continue Listener
							}
						}
						if cfg.DMonReact.Emoji != "" {
							if event.Data.Emoji.Name != cfg.DMonReact.Emoji && fmt.Sprintf(`%v:%v`, event.Data.Emoji.Name, event.Data.Emoji.ID) != cfg.DMonReact.Emoji {
								color.Cyan("[%v][O] Skipping event from other emoji [%v]", time.Now().Format("15:04:05"), event.Data.Emoji.Name)
								continue Listener
							}
						}
						// React is approved
						go func() {
							filteredReacts <- event.Data.UserID
						}()
						continue Listener
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
					if cfg.DMonReact.ServerID != "" && observerInstance.Ws != nil {
						err := instance.Subscribe(observerInstance.Ws, cfg.DMonReact.ServerID)
						if err != nil {
							color.Red("[%v][!] Error while subscribing to server: %s", time.Now().Format("15:04:05"), err)
							time.Sleep(10 * time.Second)
							continue Listener
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
				if cfg.DMonReact.ServerID != "" && observerInstance.Ws != nil {
					err := instance.Subscribe(observerInstance.Ws, cfg.DMonReact.ServerID)
					if err != nil {
						color.Red("[%v][!] Error while subscribing to server: %s", time.Now().Format("15:04:05"), err)
						time.Sleep(10 * time.Second)
						continue Listener
					}

				}
				color.Yellow("[%v][O] Reconnected Observer token", time.Now().Format("15:04:05"))
				continue Listener
			}
		}
		color.Red("[%v][O] Permanently disconnected Observer token", time.Now().Format("15:04:05"))
	}()
	// Starting token
Token:
	for {
		if len(tokenPool) == 0 {
			color.Red("[%v][!] No more tokens available, exiting", time.Now().Format("15:04:05"))
			kill <- true
			break Token
		}
		instance := <-tokenPool
		color.Yellow("[%v][X] Starting token %v", time.Now().Format("15:04:05"), instance.Token)
	React:
		for {
			status := instance.CheckToken()
			if status != 200 {
				color.Red("[%v][!] Token %v may be invalid [%v], skipping!", time.Now().Format("15:04:05"), instance.Token, status)
				continue Token
			}
			if cfg.DMonReact.Invite != "" && !instance.Invited {
				err := instance.Invite(cfg.DMonReact.Invite)
				if err != nil {
					color.Red("[%v][!] Error while inviting Token %v: %v Switching!", time.Now().Format("15:04:05"), instance.Token, err)
					continue Token
				}
				instance.Invited = true
			}
			if instance.Cookie == "" {
				cookie, err := instance.GetCookieString()
				if err != nil {
					color.Red("[%v][!] Error while getting cookie for Token %v: %v Switching!", time.Now().Format("15:04:05"), instance.Token, err)
					if cfg.DMonReact.RotateTokens {
						go func() {
							tokenPool <- instance
						}()
					}
					continue Token
				}
				instance.Cookie = cookie
			}
			if (cfg.DMonReact.ChangeAvatar && !instance.ChangedAvatar) || (cfg.DMonReact.ChangeName && !instance.ChangedName) {
				// Opening Websocket to change name/avatar
				err := instance.StartWS()
				if err != nil {
					color.Red("[%v][X] Error while opening websocket %v: %v", time.Now().Format("15:04:05"), instance.Token, err)
					if cfg.DMonReact.RotateTokens {
						go func() {
							tokenPool <- instance
						}()
					}
					continue Token
				} else {
					color.Green("[%v][X] Websocket opened %v", time.Now().Format("15:04:05"), instance.Token)
				}
				if cfg.DMonReact.ChangeAvatar && !instance.ChangedAvatar {

					r, err := instance.AvatarChanger(avatars[rand.Intn(len(avatars))])
					if err != nil {
						color.Red("[%v][X] %v Error while changing avatar: %v", time.Now().Format("15:04:05"), instance.Token, err)
						if cfg.DMonReact.RotateTokens {
							go func() {
								tokenPool <- instance
							}()
						}
						continue Token
					} else {
						if r.StatusCode == 204 || r.StatusCode == 200 {
							color.Green("[%v][X] %v Avatar changed successfully", time.Now().Format("15:04:05"), instance.Token)
							instance.ChangedAvatar = true
						} else {
							color.Red("[%v][X] %v Error while changing avatar: %v", time.Now().Format("15:04:05"), instance.Token, r.StatusCode)
							if cfg.DMonReact.RotateTokens {
								go func() {
									tokenPool <- instance
								}()
							}
							continue Token
						}
					}

				}
				if cfg.DMonReact.ChangeName && !instance.ChangedName {
					r, err := instance.NameChanger(names[rand.Intn(len(names))])
					if err != nil {
						color.Red("[%v]][X] %v Error while changing name: %v", time.Now().Format("15:04:05"), instance.Token, err)
						if cfg.DMonReact.RotateTokens {
							go func() {
								tokenPool <- instance
							}()
						}
						continue Token
					}
					body, err := utilities.ReadBody(r)
					if err != nil {
						fmt.Println(err)
						if cfg.DMonReact.RotateTokens {
							go func() {
								tokenPool <- instance
							}()
						}
						continue Token
					}
					if r.StatusCode == 200 || r.StatusCode == 204 {
						color.Green("[%v][X] %v Changed name successfully", time.Now().Format("15:04:05"), instance.Token)
						instance.ChangedName = true
					} else {
						color.Red("[%v][X] %v Error while changing name: %v %v", time.Now().Format("15:04:05"), instance.Token, r.Status, string(body))
						if cfg.DMonReact.RotateTokens {
							go func() {
								tokenPool <- instance
							}()
						}
						continue Token
					}
				}
				// Closing websocket
				if instance.Ws != nil {
					err = instance.Ws.Close()
					if err != nil {
						color.Red("[%v][X] Error while closing websocket: %v", time.Now().Format("15:04:05"), err)
					} else {
						color.Green("[%v][X] Websocket closed %v", time.Now().Format("15:04:05"), instance.Token)
					}
				}
			}
			if cfg.DMonReact.ServerID != "" && (instance.TimeServerCheck.Second() >= 120 || instance.TimeServerCheck.IsZero()) {
				r, err := instance.ServerCheck(cfg.DMonReact.ServerID)
				if err != nil {
					color.Red("[%v][!] Error while checking if token %v is present in server %v: %v Switching!", time.Now().Format("15:04:05"), instance.Token, cfg.DMonReact.ServerID, err)
					continue Token
				} else {
					if r != 200 && r != 204 {
						color.Red("[%v][!] Token %v is not present in server %v: %v Switching!", time.Now().Format("15:04:05"), instance.Token, cfg.DMonReact.ServerID, err)
						continue Token
					}
				}
				instance.TimeServerCheck = time.Now()
			}
			ticker := make(chan bool)
			go func() {
				for {
					time.Sleep(180 * time.Second)
					ticker <- true
				}
			}()
			select {
			case uuid := <-filteredReacts:
				instance.Count++
				if cfg.DMonReact.MaxDMsPerToken != 0 && instance.Count >= cfg.DMonReact.MaxDMsPerToken {
					color.Red("[%v] %v Max DMs reached, switching token", time.Now().Format("15:04:05"), instance.Token)
					continue Token
				}
				t := time.Now()
				snowflake, err := instance.OpenChannel(uuid)
				if err != nil {
					color.Red("[%v] Error while opening channel: %v", time.Now().Format("15:04:05"), err)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					failed = append(failed, uuid)
					continue React
				}
				resp, err := instance.SendMessage(snowflake, uuid)
				if err != nil {
					color.Red("[%v] Error while sending message: %v", time.Now().Format("15:04:05"), err)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					failed = append(failed, uuid)
					continue React
				}
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					color.Red("[%v] Error while reading body: %v", time.Now().Format("15:04:05"), err)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					failed = append(failed, uuid)
					continue React
				}
				var response jsonResponse
				err = json.Unmarshal(body, &response)
				if err != nil {
					color.Red("[%v] Error while unmarshalling body: %v", time.Now().Format("15:04:05"), err)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					failed = append(failed, uuid)
					continue React
				}
				if resp.StatusCode == 200 {
					color.Green("[%v][X] Token %v DM'd %v [%vms]", time.Now().Format("15:04:05"), instance.Token, uuid, time.Since(t).Milliseconds())
					completed = append(completed, uuid)
					err = utilities.WriteLine("input/completed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					continue React
				} else if resp.StatusCode == 403 && response.Code == 40003 {
					// Token is rate limited
					go func() {
						filteredReacts <- uuid
					}()
					if cfg.DMonReact.LeaveTokenOnRateLimit && cfg.DMonReact.ServerID != "" {
						re := instance.Leave(cfg.DMonReact.ServerID)
						if re == 200 || re == 204 {
							color.Green("[%v][!] Token %v left server %v", time.Now().Format("15:04:05"), instance.Token, cfg.DMonReact.ServerID)
						} else {
							color.Red("[%v][!] Error while leaving server %v: %v", time.Now().Format("15:04:05"), cfg.DMonReact.ServerID, re)
						}
					}
					if cfg.DMonReact.RotateTokens {
						go func() {
							tokenPool <- instance
						}()
					}
					continue Token
				} else if resp.StatusCode == 403 && response.Code == 50007 {
					failed = append(failed, uuid)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					color.Red("[%v][X] Token %v failed to DM %v [DMs Closed or No mutual servers] [%vms]", time.Now().Format("15:04:05"), instance.Token, uuid, time.Since(t).Milliseconds())
					continue React
				} else if resp.StatusCode == 403 && response.Code == 40002 || resp.StatusCode == 401 || resp.StatusCode == 405 {
					failed = append(failed, uuid)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					color.Red("[%v][X] Token %v failed to DM %v [Locked/Disabled][%vms]", time.Now().Format("15:04:05"), instance.Token, uuid, time.Since(t).Milliseconds())
					continue React
				} else if resp.StatusCode == 429 {
					failed = append(failed, uuid)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					color.Red("[%v][X] Token %v failed to DM %v [Rate Limited][%vms]", time.Now().Format("15:04:05"), instance.Token, uuid, time.Since(t).Milliseconds())
					time.Sleep(5 * time.Second)
					continue React
				} else if resp.StatusCode == 400 && strings.Contains(string(body), "captcha") {
					color.Red("[%v] Token %v Captcha was solved incorrectly", time.Now().Format("15:04:05"), instance.Token)
					if instance.Config.CaptchaSettings.CaptchaAPI == "anti-captcha.com" {
						err := instance.ReportIncorrectRecaptcha()
						if err != nil {
							color.Red("[%v] Error while reporting incorrect hcaptcha: %v", time.Now().Format("15:04:05"), err)
						} else {
							color.Green("[%v] Succesfully reported incorrect hcaptcha [%v]", time.Now().Format("15:04:05"), instance.LastID)
						}
					}
				
				} else {
					failed = append(failed, uuid)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						fmt.Println(err)
					}
					color.Red("[%v][X] Token %v failed to DM %v [%v][%vms]", time.Now().Format("15:04:05"), instance.Token, uuid, string(body), time.Since(t).Milliseconds())
					continue React
				}
			case <-ticker:
				color.Yellow("[%v][X] %v Refreshing token", time.Now().Format("15:04:05"), instance.Token)
				continue React
			}
		}
	}
}
