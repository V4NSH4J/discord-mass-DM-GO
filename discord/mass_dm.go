// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchMassDM() {
	members, err := utilities.ReadLines("memberids.txt")
	if err != nil {
		utilities.LogErr("Error while opening MemberIDs file %s", err)
		return
	}
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %s", err)
		return
	}
	var tokenFile, completedUsersFile, failedUsersFile, lockedFile, quarantinedFile, logsFile string
	if cfg.OtherSettings.Logs {
		path := fmt.Sprintf(`logs/mass_dm/DMDGO-MASSDM-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
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
		LogsX, err := os.Create(fmt.Sprintf(`%s/logs.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating failed file: %s", err)
			utilities.ExitSafely()
		}
		LogsX.Close()
		tokenFile, completedUsersFile, failedUsersFile, lockedFile, quarantinedFile, logsFile = tokenFileX.Name(), completedUsersFileX.Name(), failedUsersFileX.Name(), lockedFileX.Name(), quarantinedFileX.Name(), LogsX.Name()
		for i := 0; i < len(instances); i++ {
			instances[i].WriteInstanceToFile(tokenFile)
		}
	}
	if cfg.OtherSettings.Logs {
		utilities.WriteLinesPath(logsFile, fmt.Sprintf("Start Time: %v", time.Now()))
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
	if cfg.OtherSettings.Logs {
		if len(instances) > 0 {
			utilities.WriteLinesPath(logsFile, fmt.Sprintf("Messages Loaded: %v", instances[0].Messages))
		}
	}
	advancedchoice := utilities.UserInputInteger("Do you wish to use Advanced Settings? 0: No, 1: Yes: ")

	var checkchoice int
	var serverid string
	var tryjoinchoice int
	var invite string
	var maxattempts int
	if advancedchoice != 0 && advancedchoice != 1 {
		utilities.LogErr("Invalid choice")
		return
	}
	if advancedchoice == 1 {
		checkchoice := utilities.UserInputInteger("Do you wish to check if token is still in server before every DM? [0: No, 1: Yes]")
		if checkchoice != 0 && checkchoice != 1 {
			utilities.LogErr("Invalid choice")
			return
		}
		if checkchoice == 1 {
			serverid = utilities.UserInput("Enter the server ID: ")
			tryjoinchoice := utilities.UserInputInteger("Do you wish to try rejoining the server if token is not in server? [0: No, 1: Yes]")
			if tryjoinchoice != 0 && tryjoinchoice != 1 {
				utilities.LogErr("Invalid choice")
				return
			}
			if tryjoinchoice == 1 {
				invite = utilities.UserInput("Enter a permanent invite code")
				maxattempts = utilities.UserInputInteger("Enter max rejoin attempts")
			}
		}
	}
	// Also initiate variables and slices for logging and counting
	var session []string
	var completed []string
	var failed []string
	var dead []string
	var failedCount = 0
	var openedChannels = 0
	completed, err = utilities.ReadLines("completed.txt")
	if err != nil {
		utilities.LogErr("Error while opening completed.txt %s", err)
		return
	}
	if cfg.DirectMessage.Skip {
		members = utilities.RemoveSubset(members, completed)
		if cfg.OtherSettings.Logs {
			utilities.WriteLinesPath(logsFile, fmt.Sprintf("Users blacklisted from completed.txt: %v", len(completed)))
		}
	}
	if cfg.DirectMessage.SkipFailed {
		failedSkip, err := utilities.ReadLines("failed.txt")
		if err != nil {
			utilities.LogErr("Error while opening failed.txt %s", err)
			return
		}
		if cfg.OtherSettings.Logs {
			utilities.WriteLinesPath(logsFile, fmt.Sprintf("Users blacklisted from failed.txt: %v", len(failedSkip)))
		}
		members = utilities.RemoveSubset(members, failedSkip)
	}
	if len(instances) == 0 {
		utilities.LogErr("Enter your tokens in tokens.txt")
		if cfg.OtherSettings.Logs {
			utilities.WriteLinesPath(logsFile, fmt.Sprintf("Tokens loaded: %v", len(instances)))
		}
		return
	}
	if len(members) == 0 {
		utilities.LogErr("Enter your memberids and ensure they're not all in completed.txt or failed.txt")
		return
	}
	if len(members) < len(instances) {
		instances = instances[:len(members)]
	}
	if cfg.OtherSettings.Logs {
		utilities.WriteLinesPath(logsFile, fmt.Sprintf("Unique members loaded: %v", len(members)))
	}
	msgs := instances[0].Messages
	for i := 0; i < len(msgs); i++ {
		if msgs[i].Content == "" && msgs[i].Embeds == nil {
			utilities.LogWarn("Message %v is empty", i)
		}
	}
	// Send members to a channel
	mem := make(chan string, len(members))
	go func() {
		for i := 0; i < len(members); i++ {
			mem <- members[i]
		}
	}()
	ticker := make(chan bool)
	// Setting information to windows titlebar by github.com/foxzsz
	go func() {
	Out:
		for {
			select {
			case <-ticker:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%d sent, %v failed, %d locked, %v avg. dms, %v avg. channels, %d tokens left]`, len(session), len(failed), len(dead), len(session)/len(instances), openedChannels/len(instances), len(instances)-len(dead)))
				_ = cmd.Run()
			}

		}
	}()
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < len(instances); i++ {
		// Offset goroutines by a few milliseconds. Makes a big difference and allows for better concurrency
		time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				// Get a member from the channel
				if len(mem) == 0 {
					break
				}
				member := <-mem
				instances[i].LastIDstr = ""
				// Breaking loop if maximum DMs reached
				if cfg.DirectMessage.MaxDMS != 0 && instances[i].Count >= cfg.DirectMessage.MaxDMS {
					utilities.LogInfo("Maximum DMs reached for %v", instances[i].CensorToken())
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token max DMs reached", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
					}
					break
				}
				// Start websocket connection if not already connected and reconnect if dead
				if cfg.DirectMessage.Websocket && instances[i].Ws == nil {
					err := instances[i].StartWS()
					if err != nil {
						utilities.LogFailed("Error while opening websocket: %v", err)
					} else {
						utilities.LogSuccess("Websocket opened %v", instances[i].CensorToken())
					}
				}
				// Check if token is valid
				status := instances[i].CheckToken()
				if status != 200 && status != 204 && status != 429 && status != -1 {
					failedCount++
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					utilities.LogLocked("Token %v might be locked - Stopping instance and adding members to failed list. %v [%v]", instances[i].CensorToken(), status, failedCount)
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token locked", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
					}
					failed = append(failed, member)
					dead = append(dead, instances[i].Token)
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(lockedFile)
					}
					err := utilities.WriteLines("failed.txt", member)
					if err != nil {
						utilities.LogErr("Error while writing to failed.txt %s", err)
					}
					if cfg.DirectMessage.Stop {
						break
					}
				}
				// Advanced Options
				if advancedchoice == 1 {
					if checkchoice == 1 {
						r, err := instances[i].ServerCheck(serverid)
						if err != nil {
							utilities.LogErr("Error while checking server %s", err)
							continue
						}
						if r != 200 && r != 204 && r != 429 {
							if tryjoinchoice == 0 {
								utilities.LogFailed("Token %s is not present in server %s", instances[i].CensorToken(), serverid)
								if cfg.OtherSettings.Logs {
									utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token not present in server %v", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), serverid))
								}
								break
							} else {
								if instances[i].Retry >= maxattempts {
									if cfg.OtherSettings.Logs {
										utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token max rejoin attempts", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
									}
									utilities.LogFailed("Stopping token %v [Max server rejoin attempts]", instances[i].CensorToken())
									break
								}
								err := instances[i].Invite(invite)
								if err != nil {
									utilities.LogFailed("Error while joining server: %v", err)
									if cfg.OtherSettings.Logs {
										utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token error while joining server %v", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), err))
									}
									instances[i].Retry++
									continue
								}
							}
						}
					}
				}
				var user string
				user = member
				// Check Mutual
				if cfg.DirectMessage.Mutual {
					info, err := instances[i].UserInfo(member)
					if err != nil {
						failedCount++
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(failedUsersFile, member)
						}
						utilities.LogErr("Error while getting user info: %v [%v]", err, failedCount)
						err = utilities.WriteLine("input/failed.txt", member)
						if err != nil {
							utilities.LogErr("Error while writing to failed.txt %s", err)
						}
						failed = append(failed, member)

						continue
					}
					if len(info.Mutual) == 0 {
						failedCount++
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(failedUsersFile, member)
						}
						utilities.LogFailed("Token %v failed to DM %v [No Mutual Server] [%v]", instances[i].CensorToken(), info.User.Username+info.User.Discriminator, failedCount)
						if cfg.OtherSettings.Logs {
							utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token failed to DM %v [No mutuals]", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), member))
						}
						err = utilities.WriteLine("input/failed.txt", member)
						if err != nil {
							utilities.LogErr("Error while writing to failed.txt %s", err)
						}
						failed = append(failed, member)
						continue
					}
					user = info.User.Username + "#" + info.User.Discriminator
					// Used only if Websocket is enabled as Unwebsocketed Tokens get locked if they attempt to send friend requests.
					if cfg.DirectMessage.Friend && cfg.DirectMessage.Websocket {
						x, err := strconv.Atoi(info.User.Discriminator)
						if err != nil {
							utilities.LogErr("Error while converting discriminator to int: %v", err)
							continue
						}
						resp, err := instances[i].Friend(info.User.Username, x)
						if err != nil {
							utilities.LogErr("Error while sending friend request: %v", err)
							continue
						}
						defer resp.Body.Close()
						if resp.StatusCode != 204 && err != nil {
							if !errors.Is(err, io.ErrUnexpectedEOF) {
								body, err := utilities.ReadBody(*resp)
								if err != nil {
									utilities.LogErr("Error while reading body: %v", err)
									continue
								}
								utilities.LogFailed("Error while sending friend request: %v", body)
								continue
							}
							utilities.LogErr("Error while sending friend request: %v", err)
							continue
						} else {
							utilities.LogSuccess("Friend request sent to %v", info.User.Username+info.User.Discriminator)
							if cfg.OtherSettings.Logs {
								utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token friended %v", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), member))
							}
						}
					}
				}
				// Open channel to get snowflake
				snowflake, err := instances[i].OpenChannel(member)
				if err != nil {
					failedCount++
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					utilities.LogErr("Error while opening DM channel: %v [%v]", err, failedCount)
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token error %v while opening channel", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), err))
					}
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						utilities.LogErr("Error while writing to failed.txt %s", err)
					}
					failed = append(failed, member)
					if instances[i].Quarantined {
						break
					}
					continue
				}
				if cfg.SuspicionAvoidance.RandomDelayOpenChannel != 0 {
					time.Sleep(time.Duration(rand.Intn(cfg.SuspicionAvoidance.RandomDelayOpenChannel)) * time.Second)
				}
				respCode, body, err := instances[i].SendMessage(snowflake, member)
				openedChannels++
				if err != nil {
					failedCount++
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					utilities.LogErr("Error while sending message: %v [%v]", err, failedCount)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						utilities.LogErr("Error while writing to failed.txt %s", err)
					}
					failed = append(failed, member)
					continue
				}
				var response jsonResponse
				errx := json.Unmarshal(body, &response)
				if errx != nil {
					failedCount++
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					utilities.LogErr("Error while unmarshalling body: %v [%v]", errx, failedCount)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						utilities.LogErr("Error while writing to failed.txt %s", err)
					}
					failed = append(failed, member)
					continue
				}
				// Everything is fine, continue as usual
				if respCode == 200 {
					err = utilities.WriteLine("input/completed.txt", member)
					if err != nil {
						utilities.LogErr("Error while writing to completed.txt %s", err)
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(completedUsersFile, member)
					}
					completed = append(completed, member)
					session = append(session, member)
					utilities.LogSuccess("[DM-%v] Token %v sent DM to %v", len(session), instances[i].CensorToken(), user)
					if cfg.DirectMessage.Websocket && cfg.DirectMessage.Call && instances[i].Ws != nil {
						err := instances[i].Call(snowflake)
						if err != nil {
							utilities.LogErr("Token %v Error while calling %v: %v", instances[i].CensorToken(), user, err)
						}
						// Unfriended people can't ring.
						//
						// resp, err := instance.Ring(instances[i].Client, instances[i].Token, snowflake)
						// if err != nil {
						//      color.Red("%v Error while ringing %v: %v", instances[i].Token, user, err)
						// }
						// if resp == 200 || resp == 204 {
						//      color.Green("%v Ringed %v", instances[i].Token, user)
						// } else {
						//      color.Red("%v Error while ringing %v: %v", instances[i].Token, user, resp)
						// }

					}
					if cfg.DirectMessage.Block {
						r, err := instances[i].BlockUser(member)
						if err != nil {
							utilities.LogErr("Error while blocking user: %v", err)
						} else {
							if r == 204 {
								utilities.LogSuccess("Blocked %v", user)
								if cfg.OtherSettings.Logs {
									utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token blocked user %v", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), member))
								}
							} else {
								utilities.LogErr("Error while blocking user: %v", r)
							}
						}
					}
					if cfg.DirectMessage.Close {
						r, err := instances[i].CloseDMS(snowflake)
						if err != nil {
							utilities.LogErr("Error while closing DM: %v", err)
						} else {
							if r == 200 {
								utilities.LogSuccess("Closed %v", user)
								if cfg.OtherSettings.Logs {
									utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token closed DM of user %v", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), member))
								}
							} else {
								utilities.LogErr("Error while closing DM: %v", r)
							}
						}
					}
					// Forbidden - Token is being rate limited
				} else if response.Code == 20026 {
					utilities.LogLocked("Token %v is Quarantined, considering it locked", instances[i].CensorToken())
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token quarantined", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
					}
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					dead = append(dead, instances[i].Token)
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(lockedFile)
					}
					// Stop token if locked or disabled
					if cfg.DirectMessage.Stop {
						break
					}
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(quarantinedFile)
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token quarantined", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
					}
					mem <- member

				} else if respCode == 403 && response.Code == 40003 {
					mem <- member
					utilities.LogInfo("Token %v sleeping for %v minutes!", instances[i].CensorToken(), int(cfg.DirectMessage.LongDelay/60))
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token rate limited, sleeping for %v seconds", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), cfg.DirectMessage.LongDelay))
					}
					time.Sleep(time.Duration(cfg.DirectMessage.LongDelay) * time.Second)
					if cfg.SuspicionAvoidance.RandomRateLimitDelay != 0 {
						time.Sleep(time.Duration(rand.Intn(cfg.SuspicionAvoidance.RandomRateLimitDelay)) * time.Second)
					}
					utilities.LogInfo("Token %v continuing!", instances[i].CensorToken())
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token continuing", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
					}
					// Forbidden - DM's are closed
				} else if respCode == 403 && response.Code == 50007 {
					failedCount++
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					failed = append(failed, member)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						utilities.LogErr("Error while writing to failed.txt %s", err)
					}
					utilities.LogFailed("Token %v failed to DM %v User has DMs closed or not present in server ", instances[i].CensorToken(), user)
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token failed to DM %v [DMs closed or no mutual servers]", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), member))
					}
					// Forbidden - Locked or Disabled
				} else if (respCode == 403 && response.Code == 40002) || respCode == 401 || respCode == 405 {
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					utilities.LogFailed("Token %v is locked or disabled. Stopping instance. %v %v", instances[i].CensorToken(), respCode, string(body))
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token locked or disabled", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
					}
					dead = append(dead, instances[i].Token)
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(lockedFile)
					}
					// Stop token if locked or disabled
					if cfg.DirectMessage.Stop {
						break
					}
					mem <- member
					// Forbidden - Invalid token
				} else if respCode == 403 && response.Code == 50009 {
					failedCount++
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					failed = append(failed, member)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						utilities.LogErr("Error while writing to failed.txt %s", err)
					}
					utilities.LogFailed("Token %v can't DM %v. It may not have bypassed membership screening or it's verification level is too low or the server requires new members to wait 10 minutes before they can interact in the server. %v", instances[i].CensorToken(), user, string(body))
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token channel verification level too high", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
					}
					// General case - Continue loop. If problem with instance, it will be stopped at start of loop.
				} else if respCode == 429 {
					utilities.LogFailed("Token %v is being rate limited. Sleeping for 10 seconds", instances[i].CensorToken())
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token rate limited", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken()))
					}
					mem <- member
					time.Sleep(10 * time.Second)
				} else if respCode == 400 && strings.Contains(string(body), "captcha") {
					mem <- member
					utilities.LogFailed("Token %v Captcha was solved incorrectly", instances[i].CensorToken())
					if instances[i].Config.CaptchaSettings.CaptchaAPI == "anti-captcha.com" {
						err := instances[i].ReportIncorrectRecaptcha()
						if err != nil {
							utilities.LogFailed("Error while reporting incorrect hcaptcha: %v", err)
						} else {
							utilities.LogSuccess("Succesfully reported incorrect hcaptcha [%v]", instances[i].LastID)
						}
					}
					instances[i].Retry++
					if instances[i].Retry >= cfg.CaptchaSettings.MaxCaptchaDM && cfg.CaptchaSettings.MaxCaptchaDM != 0 {
						utilities.LogFailed("Stopping token %v max captcha solves reached", instances[i].CensorToken())
						break
					}
				} else {
					failedCount++
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(failedUsersFile, member)
					}
					failed = append(failed, member)
					err = utilities.WriteLine("input/failed.txt", member)
					if err != nil {
						utilities.LogErr("Error while writing to failed.txt %s", err)
					}
					utilities.LogFailed("Token %v couldn't DM %v Error Code: %v; Status: %v; Message: %v", instances[i].CensorToken(), user, response.Code, respCode, response.Message)
					if cfg.OtherSettings.Logs {
						utilities.WriteLinesPath(logsFile, fmt.Sprintf("[%v][Success:%v][Failed:%v] %v token failed to DM %v error %v", time.Now().Format("15:04:05"), len(session), len(failed), instances[i].CensorToken(), user, response.Message))
					}
				}
				time.Sleep(time.Duration(cfg.DirectMessage.Delay) * time.Second)
				if cfg.SuspicionAvoidance.RandomIndividualDelay != 0 {
					time.Sleep(time.Duration(rand.Intn(cfg.SuspicionAvoidance.RandomIndividualDelay)) * time.Second)
				}
			}
		}(i)
	}
	wg.Wait()

	utilities.LogSuccess("Threads have finished! Writing to file")
	ticker <- true
	elapsed := time.Since(start)
	utilities.LogSuccess("DM advertisement took %v. Successfully sent DMs to %v IDs. Failed to send DMs to %v IDs. %v tokens are dis-functional & %v tokens are functioning", elapsed.Seconds(), len(session), len(failed), len(dead), len(instances)-len(dead))
	if cfg.OtherSettings.Logs {
		utilities.WriteLinesPath(logsFile, fmt.Sprintf("DM advertisement took %v. Successfully sent DMs to %v IDs. Failed to send DMs to %v IDs. %v tokens are dis-functional & %v tokens are functioning", elapsed.Seconds(), len(session), len(failed), len(dead), len(instances)-len(dead)))
	}
	var left []string
	if cfg.DirectMessage.Remove {
		for i := 0; i < len(instances); i++ {
			if !utilities.Contains(dead, instances[i].Token) {
				if instances[i].Password == "" {
					left = append(left, instances[i].Token)
				} else {
					left = append(left, fmt.Sprintf(`%v:%v:%v`, instances[i].Email, instances[i].Password, instances[i].Token))
				}
			}
		}
		err := utilities.Truncate("input/tokens.txt", left)
		if err != nil {
			utilities.LogErr("Error while writing to failed.txt %s", err)
		}
		utilities.LogSuccess("Updated tokens.txt")
	}
	if cfg.DirectMessage.RemoveM {
		m := utilities.RemoveSubset(members, completed)
		err := utilities.Truncate("input/memberids.txt", m)
		if err != nil {
			utilities.LogErr("Error while writing to failed.txt %s", err)
		}
		utilities.LogSuccess("Updated memberids.txt")

	}
	if cfg.DirectMessage.Websocket {
		for i := 0; i < len(instances); i++ {
			if instances[i].Ws != nil {
				instances[i].Ws.Close()
			}
		}
	}

}

type jsonResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func LaunchSingleDM() {
	choice := utilities.UserInputInteger("Enter 0 for one message; Enter 1 for continuous spam")
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances%s", err)
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

	victim := utilities.UserInput("Ensure a common link and enter victim's ID: ")
	var wg sync.WaitGroup
	wg.Add(len(instances))
	if choice == 0 {
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
			go func(i int) {
				defer wg.Done()
				snowflake, err := instances[i].OpenChannel(victim)
				if err != nil {
					utilities.LogErr("Error while opening channel %s", err)
				}
				respCode, body, err := instances[i].SendMessage(snowflake, victim)
				if err != nil {
					utilities.LogErr("Error while sending message%s", err)
				}
				if respCode == 200 {
					utilities.LogSuccess("Token %v DM'd %v", instances[i].Token, victim)
				} else {
					utilities.LogFailed("Token %v failed to DM %v [%v]", instances[i].Token, victim, string(body))
				}
			}(i)
		}
		wg.Wait()
	}
	if choice == 1 {
		for i := 0; i < len(instances); i++ {
			time.Sleep(time.Duration(cfg.DirectMessage.Offset) * time.Millisecond)
			go func(i int) {
				defer wg.Done()

				var c int
				for {
					snowflake, err := instances[i].OpenChannel(victim)
					if err != nil {
						utilities.LogErr("Error while opening channel %s", err)
					}
					respCode, _, err := instances[i].SendMessage(snowflake, victim)
					if err != nil {
						utilities.LogErr("Error while sending message %s", err)
					}
					if respCode == 200 {
						utilities.LogSuccess("Token %v DM'd %v [%v]", instances[i].CensorToken(), victim, c)
					} else {
						utilities.LogFailed("Token %v failed to DM %v", instances[i].CensorToken(), victim)
					}
					c++
				}
			}(i)
			wg.Wait()
		}
	}
	utilities.LogSuccess("All threads finished")
}
