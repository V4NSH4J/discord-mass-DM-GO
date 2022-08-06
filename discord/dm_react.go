// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	http "github.com/Danny-Dasilva/fhttp"

	goclient "github.com/V4NSH4J/discord-mass-dm-GO/client"
	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchDMReact() {
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %s", err)
		return
	}
	// Setting the titlebar on windows
	var ReactCount, ApproveCount, SuccessCount, LockedCount int
	var LastDM time.Time
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Reacts, %v Approved, %v Success, %v Failed, %v Locked, Last DM %v ago]`, ReactCount, ApproveCount, SuccessCount, ApproveCount-SuccessCount, LockedCount, time.Since(LastDM).Round(time.Second)))
				_ = cmd.Run()
			}

		}
	}()
	var tokenFile, completedUsersFile, failedUsersFile, lockedFile, quarantinedFile, eventsFile, logsFile string
	if cfg.OtherSettings.Logs {
		path := fmt.Sprintf(`logs/dm_react/DMDGO-DMR-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
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
		completedUsersFileX, err := os.Create(fmt.Sprintf(`%s/success.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating success file: %s", err)
			utilities.ExitSafely()
		}
		completedUsersFileX.Close()
		failedUsersFileX, err := os.Create(fmt.Sprintf(`%s/failed.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating failed file: %s", err)
			utilities.ExitSafely()
		}
		failedUsersFileX.Close()
		lockedFileX, err := os.Create(fmt.Sprintf(`%s/locked.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating failed file: %s", err)
			utilities.ExitSafely()
		}
		lockedFileX.Close()
		quarantinedFileX, err := os.Create(fmt.Sprintf(`%s/quarantined.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating failed file: %s", err)
			utilities.ExitSafely()
		}
		quarantinedFileX.Close()
		eventsFileX, err := os.Create(fmt.Sprintf(`%s/events.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating failed file: %s", err)
			utilities.ExitSafely()
		}
		eventsFileX.Close()
		LogsX, err := os.Create(fmt.Sprintf(`%s/logs.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating failed file: %s", err)
			utilities.ExitSafely()
		}
		LogsX.Close()
		tokenFile, completedUsersFile, failedUsersFile, lockedFile, quarantinedFile, eventsFile, logsFile = tokenFileX.Name(), completedUsersFileX.Name(), failedUsersFileX.Name(), lockedFileX.Name(), quarantinedFileX.Name(), eventsFileX.Name(), LogsX.Name()
		for i := 0; i < len(instances); i++ {
			instances[i].WriteInstanceToFile(tokenFile)
		}
	}
	// Checking config for observer token
	if cfg.DMonReact.Observer == "" {
		utilities.LogErr("Observer token not set in config")
		return
	}
	if cfg.DMonReact.ServerID == "" {
		utilities.LogErr("Set a Server ID to use DM on react")
		return
	}
	if cfg.DMonReact.Invite == "" {
		utilities.LogErr("Set an Invite to use DM on react")
		return
	}
	var msg instance.Message
	messagechoice := utilities.UserInputInteger("Enter 1 to use message from file, 2 to use message from console: ")
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
	// Initializing Files & Variables
	var completed []string
	var failed []string
	var names []string
	var avatars []string
	var proxies []string
	if cfg.DMonReact.SkipCompleted {
		completed, err = utilities.ReadLines("completed.txt")
		if err != nil {
			utilities.LogErr("Error while reading completed.txt: %s", err)
			return
		}
	}
	if cfg.DMonReact.SkipFailed {
		failed, err = utilities.ReadLines("failed.txt")
		if err != nil {
			utilities.LogErr("Error while reading failed.txt: %s", err)
			return
		}
	}
	if cfg.DMonReact.ChangeAvatar {
		utilities.LogInfo("Loading Avatars..")
		ex, err := os.Executable()
		if err != nil {
			utilities.LogErr("Error while getting executable path: %s", err)
			return
		}
		ex = filepath.ToSlash(ex)
		path := path.Join(path.Dir(ex) + "/input/pfps")

		images, err := instance.GetFiles(path)
		if err != nil {
			utilities.LogErr("Couldn't find Images in pfps folder %s", err)
			return
		}
		utilities.LogInfo("Found %d images", len(images))
		if len(images) == 0 {
			cfg.DMonReact.ChangeAvatar = false
			utilities.LogWarn("No images found in pfps folder, disabling avatar change")
		} else {
			for i := 0; i < len(images); i++ {
				av, err := instance.EncodeImg(images[i])
				if err != nil {
					utilities.LogErr("Error while encoding image: %s", err)
					continue
				}
				avatars = append(avatars, av)
			}
			utilities.LogInfo("%v avatars loaded", len(avatars))
		}
	}
	if cfg.DMonReact.ChangeName {
		names, err = utilities.ReadLines("names.txt")
		if err != nil {
			utilities.LogErr("Error while reading names.txt: %s", err)
			return
		}
		if len(names) == 0 {
			cfg.DMonReact.ChangeName = false
			utilities.LogWarn("No names found in names.txt, disabling name change")
		}
	}
	if cfg.DMonReact.ChangeName {
		for i := 0; i < len(instances); i++ {
			if instances[i].Password == "" {
				cfg.DMonReact.ChangeName = false
				utilities.LogWarn("Token %s has no password, require tokens in format email:password:token to use the name changer", instances[i].CensorToken())
				break
			}
		}
	}
	if cfg.ProxySettings.ProxyFromFile {
		proxies, err = utilities.ReadLines("proxies.txt")
		if err != nil {
			utilities.LogErr("Error while reading proxies.txt: %s", err)
			return
		}
		if len(proxies) == 0 {
			cfg.ProxySettings.ProxyFromFile = false
			utilities.LogWarn("No proxies found in proxies.txt, disabling proxy change")
		}
	}
	if cfg.CaptchaSettings.ClientKey == "" {
		utilities.LogWarn("You're not using a captcha key, if you're met with a captcha the token will be cycled.")
	}
	tokenPool := make(chan instance.Instance, len(instances))
	for i := 0; i < len(instances); i++ {
		go func(i int) {
			if instances[i].Token != cfg.DMonReact.Observer {
				tokenPool <- instances[i]
			} else {
				utilities.LogInfo("Skipping observer token %s", instances[i].CensorToken())
			}
		}(i)
	}

	// All files and variables loaded and errors handled.
	utilities.LogInfo("Initializing Observer Token from config")
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

	httpclient, err := goclient.NewClient(goclient.Browser{JA3: "771,4865-4867-4866-49195-49199-52393-52392-49196-49200-49162-49161-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-51-43-13-45-28-21,29-23-24-25-256-257,0", UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:103.0) Gecko/20100101 Firefox/103.0", Cookies: nil}, cfg.ProxySettings.Timeout, false, "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:103.0) Gecko/20100101 Firefox/103.0", "")
	if err != nil {
		utilities.LogWarn("Error while initializing client: %s using default client for observer", err)
		httpclient = http.DefaultClient
	}
	observerInstance.Client = httpclient
	observerInstance.Config = cfg
	if cfg.DMonReact.ServerID != "" {
		r, err := observerInstance.ServerCheck(cfg.DMonReact.ServerID)
		if err != nil {
			utilities.LogErr("Error while checking if observer token is present in server: %s", err)
		} else {
			if r != 200 && r != 204 {
				// Token not in server or some other issue like rate limit
				err := observerInstance.Invite(cfg.DMonReact.Invite)
				if err != nil {
					utilities.LogErr("Error while inviting observer token: %s", err)
				} else {
					utilities.LogSuccess("Observer token invited to server")
				}

			}
		}
		r, err = observerInstance.ServerCheck(cfg.DMonReact.ServerID)
		if err != nil {
			utilities.LogErr("Error while checking if observer token is present in server: %s", err)
			return
		} else {
			if r != 200 && r != 204 {
				utilities.LogErr("Invalid response code %s while checking if observer token is present in server", r)
				return
			}
		}

	}
	utilities.LogSuccess("Observer Token initialized")
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
						utilities.LogInfo("Disconnected observer token to reconnect")
					case x := <-observerInstance.Ws.Reactions:
						var event instance.Event
						err := json.Unmarshal(x, &event)
						if err != nil {
							utilities.LogErr("Error while unmarshalling event: %s", err)
							continue Listener
						}
						utilities.LogInfo("Event received: %v reacted [%v|%v|%v|%v]", event.Data.UserID, event.Data.GuildId, event.Data.MessageID, event.Data.ChannelID, event.Data.Emoji.Name)
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(eventsFile, fmt.Sprintf(`[%v] User ID: %s | Guild ID: %s | Message ID: %s | Channel ID: %s | Emoji: %s (%s)`, time.Now().Format("2006-01-02 15:04:05"), event.Data.UserID, event.Data.GuildId, event.Data.MessageID, event.Data.ChannelID, event.Data.Emoji.Name, event.Data.Emoji.ID))
						}
						ReactCount++
						if cfg.DMonReact.MaxAntiRaidQueue > 0 {
							if len(filteredReacts) >= cfg.DMonReact.MaxAntiRaidQueue {
								utilities.LogErr("Anti-Raid queue is full, dropping reaction. Queue length: %s", len(filteredReacts))
								continue Listener
							}
						}
						if cfg.DMonReact.SkipCompleted {
							if utilities.Contains(completed, event.Data.UserID) {
								utilities.LogInfo("Skipping completed user %s", event.Data.UserID)
								continue Listener
							}
						}
						if cfg.DMonReact.SkipFailed {
							if utilities.Contains(failed, event.Data.UserID) {
								utilities.LogInfo("Skipping failed user %s", event.Data.UserID)
								continue Listener
							}
						}
						if cfg.DMonReact.ServerID != "" {
							if event.Data.GuildId != cfg.DMonReact.ServerID {
								utilities.LogInfo("Skipping reaction from other server %s", event.Data.GuildId)
								continue Listener
							}
						}
						if cfg.DMonReact.ChannelID != "" {
							if event.Data.ChannelID != cfg.DMonReact.ChannelID {
								utilities.LogInfo("Skipping reaction from other channel %s", event.Data.ChannelID)
								continue Listener
							}
						}
						if cfg.DMonReact.MessageID != "" {
							if event.Data.MessageID != cfg.DMonReact.MessageID {
								utilities.LogInfo("Skipping reaction from other message %s", event.Data.MessageID)
								continue Listener
							}
						}
						if cfg.DMonReact.Emoji != "" {
							if event.Data.Emoji.Name != cfg.DMonReact.Emoji && fmt.Sprintf(`%v:%v`, event.Data.Emoji.Name, event.Data.Emoji.ID) != cfg.DMonReact.Emoji {
								utilities.LogInfo("Skipping reaction from other emoji %s", event.Data.Emoji.Name)
								continue Listener
							}
						}
						// React is approved
						ApproveCount++
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
						utilities.LogErr("Error while starting observer websocket: %s", err)
						time.Sleep(10 * time.Second)
						continue Listener
					}
					if cfg.DMonReact.ServerID != "" && observerInstance.Ws != nil {
						err := instance.Subscribe(observerInstance.Ws, cfg.DMonReact.ServerID, cfg.DMonReact.ChannelID)
						if err != nil {
							utilities.LogErr("Error while subscribing to server: %s", err)
							time.Sleep(10 * time.Second)
							continue Listener
						}

					}
					utilities.LogInfo("Reconnected observer websocket")
					continue Listener
				}
			} else {
				// Opening Websocket
				err := observerInstance.StartWS()
				if err != nil {
					utilities.LogErr("Error while starting observer websocket: %s", err)
					time.Sleep(10 * time.Second)
					continue Listener
				}
				if cfg.DMonReact.ServerID != "" && observerInstance.Ws != nil {
					err := instance.Subscribe(observerInstance.Ws, cfg.DMonReact.ServerID, cfg.DMonReact.ChannelID)
					if err != nil {
						utilities.LogErr("Error while subscribing to server: %s", err)
						time.Sleep(10 * time.Second)
						continue Listener
					}

				}
				utilities.LogInfo("Reconnected observer websocket")
				continue Listener
			}
		}
		utilities.LogInfo("Observer stopped permanently")
		title <- true
	}()
	// Starting token
Token:
	for {
		if len(tokenPool) == 0 {
			utilities.LogWarn("No more tokens left, Exiting")
			kill <- true
			break Token
		}
		instance := <-tokenPool
		utilities.LogInfo("Initializing token %s", instance.CensorToken())
	React:
		for {
			status := instance.CheckToken()
			if status != 200 {
				utilities.LogFailed("Token %v may be invalid [%v], skipping!", instance.CensorToken(), status)
				if cfg.OtherSettings.Logs {
					instance.WriteInstanceToFile(lockedFile)
				}
				if cfg.OtherSettings.Logs {
					utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v may be invalid. Response code %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), status))
				}
				LockedCount++
				continue Token
			}
			if cfg.DMonReact.Invite != "" && !instance.Invited {
				err := instance.Invite(cfg.DMonReact.Invite)
				if err != nil {
					utilities.LogFailed("Error while inviting Token %v: %v Switching!", instance.CensorToken(), err)
					continue Token
				}
				instance.Invited = true
			}
			if instance.Cookie == "" {
				cookie, err := instance.GetCookieString()
				if err != nil {
					utilities.LogErr("Error while getting cookie for Token %v: %v Switching!", instance.CensorToken(), err)
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
					utilities.LogErr("Error while opening websocket %v: %v", instance.CensorToken(), err)
					if cfg.DMonReact.RotateTokens {
						go func() {
							tokenPool <- instance
						}()
					}
					continue Token
				} else {
					utilities.LogSuccess("Websocket opened %v", instance.CensorToken())
				}
				if cfg.DMonReact.ChangeAvatar && !instance.ChangedAvatar {
					p := avatars[rand.Intn(len(avatars))]
					r, err := instance.AvatarChanger(p)
					if err != nil {
						utilities.LogErr("%v Error while changing avatar: %v", instance.CensorToken(), err)
						if cfg.DMonReact.RotateTokens {
							go func() {
								tokenPool <- instance
							}()
						}
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v failed to change avatar. Error %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), err))
						}
						continue Token
					} else {
						if r.StatusCode == 204 || r.StatusCode == 200 {
							utilities.LogSuccess("%v Avatar changed successfully", instance.CensorToken())
							if cfg.OtherSettings.Logs {
								utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v changed Avatar to %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), p))
							}
							instance.ChangedAvatar = true
						} else {
							utilities.LogErr("%v Error while changing avatar: %v", instance.CensorToken(), r.StatusCode)
							if cfg.DMonReact.RotateTokens {
								go func() {
									tokenPool <- instance
								}()
							}
							if cfg.OtherSettings.Logs {
								utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v failed to change avatar. %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), r.StatusCode))
							}
							continue Token
						}
					}

				}
				if cfg.DMonReact.ChangeName && !instance.ChangedName {
					p := names[rand.Intn(len(names))]
					r, err := instance.NameChanger(p)
					if err != nil {
						utilities.LogErr("%v Error while changing name: %v", instance.CensorToken(), err)
						if cfg.DMonReact.RotateTokens {
							go func() {
								tokenPool <- instance
							}()
						}
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v failed to change name. Error %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), err))
						}
						continue Token
					}
					body, err := utilities.ReadBody(r)
					if err != nil {
						utilities.LogErr("%v Error while reading body: %v", instance.CensorToken(), err)
						if cfg.DMonReact.RotateTokens {
							go func() {
								tokenPool <- instance
							}()
						}
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v failed to change name. Error %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), err))
						}
						continue Token
					}
					if r.StatusCode == 200 || r.StatusCode == 204 {
						utilities.LogSuccess("%v Changed name successfully", instance.CensorToken())
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v name changed to %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), p))
						}
						instance.ChangedName = true
					} else {
						utilities.LogErr("%v Error while changing name: %v %v", instance.CensorToken(), r.StatusCode, string(body))
						if cfg.DMonReact.RotateTokens {
							go func() {
								tokenPool <- instance
							}()
						}
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v failed to change name. %v %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), r.StatusCode, string(body)))
						}
						continue Token
					}
				}
				// Closing websocket
				if instance.Ws != nil {
					err = instance.Ws.Close()
					if err != nil {
						utilities.LogErr("Error while closing websocket: %v", err)
					} else {
						utilities.LogSuccess("Websocket closed %v", instance.CensorToken())
					}
				}
			}
			if cfg.DMonReact.ServerID != "" && (instance.TimeServerCheck.Second() >= 120 || instance.TimeServerCheck.IsZero()) {
				r, err := instance.ServerCheck(cfg.DMonReact.ServerID)
				if err != nil {
					utilities.LogErr("Error while checking if token %v is present in server %v: %v Switching!", instance.CensorToken(), cfg.DMonReact.ServerID, err)
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Error while checking if token %v is present in server %v: %v Switching!`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), cfg.DMonReact.ServerID, err))
					}
					continue Token
				} else {
					if r != 200 && r != 204 {
						utilities.LogFailed("Token %v is not present in server %v: Response %v Switching!", instance.CensorToken(), cfg.DMonReact.ServerID, r)
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v is not present in server %v: Response %v Switching!`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), cfg.DMonReact.ServerID, r))
						}
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
					utilities.LogInfo("%v Max DMs per token reached, switching", instance.CensorToken())
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v Max DMs reached %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), cfg.DMonReact.MaxDMsPerToken))
					}
					continue Token
				}
				t := time.Now()
				snowflake, err := instance.OpenChannel(uuid)
				if err != nil {
					utilities.LogErr("Error while opening channel: %v", err)
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v Error %v while opening channel`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), err))
					}
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						utilities.LogErr("Error while writing to failed file %v: %v", uuid, err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, uuid)
					}
					failed = append(failed, uuid)
					continue React
				}
				respCode, body, err := instance.SendMessage(snowflake, uuid)
				if err != nil {
					utilities.LogErr("Error while sending message: %v", err)
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v Error %v while sending message`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), err))
					}
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						utilities.LogErr("Error while writing to failed file %v: %v", uuid, err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, uuid)
					}
					failed = append(failed, uuid)
					continue React
				}
				var response jsonResponse
				err = json.Unmarshal(body, &response)
				if err != nil {
					utilities.LogErr("Error while unmarshalling body: %v", err)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						utilities.LogErr("Error while writing to failed file %v: %v", uuid, err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, uuid)
					}
					failed = append(failed, uuid)
					continue React
				}
				if respCode == 200 {
					utilities.LogSuccess("Token %v messaged %v [%v milliseconds]", instance.CensorToken(), uuid, time.Since(t).Milliseconds())
					SuccessCount++
					LastDM = time.Now()
					completed = append(completed, uuid)
					err = utilities.WriteLine("input/completed.txt", uuid)
					if err != nil {
						utilities.LogErr("Error while writing to completed file %v: %v", uuid, err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(completedUsersFile, uuid)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v messaged %v [%v milliseconds]`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), uuid, time.Since(t).Milliseconds()))
					}
					continue React
				} else if response.Code == 20026 {
					utilities.LogLocked("Token %s is Quarantined. It is being stopped.", instance.CensorToken())
					if cfg.OtherSettings.Logs {
						instance.WriteInstanceToFile(quarantinedFile)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v] Token %v is Quarantined`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken()))
					}
					LockedCount++
					continue Token
				} else if respCode == 403 && response.Code == 40003 {
					// Token is rate limited
					go func() {
						filteredReacts <- uuid
					}()
					if cfg.DMonReact.LeaveTokenOnRateLimit && cfg.DMonReact.ServerID != "" {
						re := instance.Leave(cfg.DMonReact.ServerID)
						if re == 200 || re == 204 {
							utilities.LogSuccess("Token %v left server %v", instance.CensorToken(), cfg.DMonReact.ServerID)
							if cfg.OtherSettings.Logs {
								utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v left server %v`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), cfg.DMonReact.ServerID))
							}
						} else {
							utilities.LogErr("Error while leaving server %v: %v", cfg.DMonReact.ServerID, re)
							if cfg.OtherSettings.Logs {
								utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Error while leaving server %v: %v`, time.Now().Format("2006-01-02 15:04:05"), cfg.DMonReact.ServerID, re))
							}

						}
					}
					if cfg.DMonReact.RotateTokens {
						go func() {
							tokenPool <- instance
						}()
					}
					continue Token
				} else if respCode == 403 && response.Code == 50007 {
					failed = append(failed, uuid)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						utilities.LogErr("Error while writing to failed file %v: %v", uuid, err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, uuid)
					}
					utilities.LogFailed("Token %v failed to message %v [DMs closed or no Mutual servers] [%v milliseconds]", instance.CensorToken(), uuid, time.Since(t).Milliseconds())
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v failed to messaged %v [%v milliseconds] [DMs closed or No Mutuals]`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), uuid, time.Since(t).Milliseconds()))
					}
					continue React
				} else if respCode == 403 && response.Code == 40002 || respCode == 401 || respCode == 405 {
					failed = append(failed, uuid)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						utilities.LogErr("Error while writing to failed file %v: %v", uuid, err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, uuid)
					}
					utilities.LogFailed("Token %v failed to message %v [Locked/Disabled]", instance.CensorToken(), uuid)
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v failed to messaged %v [%v milliseconds] [Locked or Disabled]`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), uuid, time.Since(t).Milliseconds()))
					}
					continue React
				} else if respCode == 429 {
					failed = append(failed, uuid)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						utilities.LogErr("Error while writing to failed file %v: %v", uuid, err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, uuid)
					}
					utilities.LogFailed("Token %v failed to DM %v [Rate Limited][%vms]", instance.CensorToken(), uuid, time.Since(t).Milliseconds())
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v failed to messaged %v [%v milliseconds] [Rate Limited]`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), uuid, time.Since(t).Milliseconds()))
					}
					time.Sleep(5 * time.Second)
					continue React
				} else if respCode == 400 && strings.Contains(string(body), "captcha") {
					utilities.LogFailed("Token %v Captcha was solved incorrectly", instance.CensorToken())
					if instance.Config.CaptchaSettings.CaptchaAPI == "anti-captcha.com" {
						err := instance.ReportIncorrectRecaptcha()
						if err != nil {
							utilities.LogErr("Error while reporting incorrect hcaptcha: %v", err)
						} else {
							utilities.LogSuccess("Reported incorrect hcaptcha %v", instance.LastID)
						}
					}

				} else {
					failed = append(failed, uuid)
					err = utilities.WriteLine("input/failed.txt", uuid)
					if err != nil {
						utilities.LogErr("Error while writing to failed file %v: %v", uuid, err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, uuid)
					}
					utilities.LogFailed("Token %v failed to DM %v [%v][%vms]", instance.CensorToken(), uuid, string(body), time.Since(t).Milliseconds())
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf(`[%v]Token %v failed to messaged %v [%v milliseconds] [%v]`, time.Now().Format("2006-01-02 15:04:05"), instance.CensorToken(), uuid, time.Since(t).Milliseconds(), string(body)))
					}
					continue React
				}
			case <-ticker:
				utilities.LogInfo("Token %v is being re-checked", instance.CensorToken())
				continue React
			}
		}
	}
}
