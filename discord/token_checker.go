// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/gookit/color"
	"github.com/zenthangplus/goccm"
)

func LaunchTokenChecker() {
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting neccessary information %v", err)
		return
	}
	var tokenFile, workingFile, lockedFile, uncheckedFile, emailVerifiedFile, phoneVerifiedFile, unverifiedFile, spammerFile, quarantinedFile string
	if cfg.OtherSettings.Logs {
		path := fmt.Sprintf(`logs/token_checker/DMDGO-TC-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
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
		workingFileX, err := os.Create(fmt.Sprintf(`%s/working.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating working file: %s", err)
			utilities.ExitSafely()
		}
		workingFileX.Close()
		lockedFileX, err := os.Create(fmt.Sprintf(`%s/locked.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating locked file: %s", err)
			utilities.ExitSafely()
		}
		lockedFileX.Close()
		uncheckedFileX, err := os.Create(fmt.Sprintf(`%s/unchecked.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating unchecked file: %s", err)
			utilities.ExitSafely()
		}
		uncheckedFileX.Close()
		emailVerifiedFileX, err := os.Create(fmt.Sprintf(`%s/email_verified.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating email verified file: %s", err)
			utilities.ExitSafely()
		}
		emailVerifiedFileX.Close()
		phoneVerifiedFileX, err := os.Create(fmt.Sprintf(`%s/phone_verified.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating phone verified file: %s", err)
			utilities.ExitSafely()
		}
		phoneVerifiedFileX.Close()
		unverifiedFileX, err := os.Create(fmt.Sprintf(`%s/unverified.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating unverified file: %s", err)
			utilities.ExitSafely()
		}
		unverifiedFileX.Close()
		spammerFileX, err := os.Create(fmt.Sprintf(`%s/spammer.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating spammer file: %s", err)
			utilities.ExitSafely()
		}
		spammerFileX.Close()
		quarantinedFileX, err := os.Create(fmt.Sprintf(`%s/quarantined.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating quarantined file: %s", err)
			utilities.ExitSafely()
		}
		quarantinedFileX.Close()
		tokenFile, workingFile, lockedFile, uncheckedFile, emailVerifiedFile, phoneVerifiedFile, unverifiedFile, spammerFile, quarantinedFile = tokenFileX.Name(), workingFileX.Name(), lockedFileX.Name(), uncheckedFileX.Name(), emailVerifiedFileX.Name(), phoneVerifiedFileX.Name(), unverifiedFileX.Name(), spammerFileX.Name(), quarantinedFileX.Name()
		for i := 0; i < len(instances); i++ {
			instances[i].WriteInstanceToFile(tokenFile)
		}
	}
	threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
	if threads > len(instances) {
		threads = len(instances)
	}
	if threads == 0 {
		threads = len(instances)
	}
	title := make(chan bool)
	var validTokens []instance.Instance
	var valid, invalid, errors int
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Unchecked %v Valid %v Invalid %v Errors]`, len(instances)-valid-invalid-errors, valid, invalid, errors))
				_ = cmd.Run()
			}

		}
	}()
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			defer c.Done()
			var printStrings []string
			// Checking Validity of token & Information
			r, err := instances[i].CheckTokenNew()
			if err != nil {
				if cfg.OtherSettings.Logs {
					instances[i].WriteInstanceToFile(uncheckedFile)
				}
				printStrings = append(printStrings, MakeColoredString("red", "FAILED", " Token %v: %v", i, instances[i].CensorToken()))
				printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Error %v while checking token", err))
				errors++
			} else {
				if r == 200 || r == 204 {
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(workingFile)
					}
					validTokens = append(validTokens, instances[i])
					printStrings = append(printStrings, MakeColoredString("green", "WORKING", " Token %v: %v", i, instances[i].CensorToken()))
					r, info, err := instances[i].AtMe()
					if err != nil {
						if cfg.OtherSettings.Logs {
							instances[i].WriteInstanceToFile(uncheckedFile)
						}
						printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Error %v while checking token", err))
						errors++
					}
					if r != 200 && r != 204 {
						printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Invalid status code %v while checking token", r))
					} else {
						if info.ID != "" {
							printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " ID: %v", info.ID))
							printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Age: %v", utilities.TimeDifference(utilities.ReverseSnowflake(info.ID), time.Now())))
							if info.Avatar != "" {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Avatar: true"))
							}
							if info.Username != "" && info.Discriminator != "" {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Username: %v#%v", info.Username, info.Discriminator))
							}
							if info.Flags != 0 {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Flags: %v", info.Flags))
							}
							if info.PublicFlags != 0 {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Public Flags: %v", info.PublicFlags))
							}
							if info.Flags != 0 {
								f := info.Flags - info.PublicFlags
								if f == 17592186044416 {
									printStrings = append(printStrings, MakeColoredString("red", "QUARANTINED", " Token is QUARANTINED"))
									if cfg.OtherSettings.Logs {
										instances[i].WriteInstanceToFile(quarantinedFile)
									}
								} else if f == 1048576 {
									printStrings = append(printStrings, MakeColoredString("red", "SPAMMER", " Token is flagged as spammer"))
									if cfg.OtherSettings.Logs {
										instances[i].WriteInstanceToFile(spammerFile)
									}
								} else if f == 17592186044416+1048576 {
									if cfg.OtherSettings.Logs {
										instances[i].WriteInstanceToFile(spammerFile)
									}
									if cfg.OtherSettings.Logs {
										instances[i].WriteInstanceToFile(quarantinedFile)
									}
									printStrings = append(printStrings, MakeColoredString("red", "SPAMMER & QUARANTINED", " Token is flagged as spammer and QUARANTINED"))
								}
							}
							if info.Bio != "" {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Bio: %v", info.Bio))
							}
							if info.MFAEnabled {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " 2FA: true"))
							}
							if info.Email != "" {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Email: %v", info.Email))
							}
							if info.Verified {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Email Verified: true"))
								if cfg.OtherSettings.Logs {
									instances[i].WriteInstanceToFile(emailVerifiedFile)
								}
							}
							if !info.Verified && info.Phone == "" {
								if cfg.OtherSettings.Logs {
									instances[i].WriteInstanceToFile(unverifiedFile)
								}
							}
							if info.Phone != "" {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Phone Verified: true [%v]", info.Phone))
								if cfg.OtherSettings.Logs {
									instances[i].WriteInstanceToFile(phoneVerifiedFile)
								}
							} else {
								printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Phone Verified: false"))
							}
							// Check Guilds
							r, guilds, _, err := instances[i].Guilds()
							if err == nil {
								if r == 200 || r == 204 {
									printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Guilds: %v", guilds))
								} else {
									printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Unexpected Response Status Code %v while checking guilds", r))
								}
							} else {
								printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Error %v while checking guilds", err))
							}
							r, channels, _, err := instances[i].Channels()
							if err == nil {
								if r == 200 || r == 204 {
									printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Open DMs: %v", channels))
								} else {
									printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Unexpected Response Status Code %v while checking channels", r))
								}
							} else {
								printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Error %v while checking channels", err))
							}
							r, friends, blocked, incoming, outgoing, _, err := instances[i].Relationships()
							if err == nil {
								if r == 200 || r == 204 {
									printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Friends: %v", friends))
									printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Blocked: %v", blocked))
									printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Incoming Requests: %v", incoming))
									printStrings = append(printStrings, MakeColoredString("cyan", "INFO", " Outgoing Requests: %v", outgoing))
								} else {
									printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Unexpected Response Status Code %v while checking relationships", r))
								}
							} else {
								printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Error %v while checking relationships", err))
							}
							printStrings = append(printStrings, "===================\n\n")
						}
					}

					valid++
				} else if r == 401 || r == 403 {
					// Locked/Invalid
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(lockedFile)
					}
					printStrings = append(printStrings, MakeColoredString("red", "LOCKED", " Token %v: %v", i, instances[i].CensorToken()))
					invalid++
				} else {
					// Invalid StatusCode
					if cfg.OtherSettings.Logs {
						instances[i].WriteInstanceToFile(uncheckedFile)
					}

					errors++
					printStrings = append(printStrings, MakeColoredString("red", "FAILED", " Token %v: %v", i, instances[i].CensorToken()))
					printStrings = append(printStrings, MakeColoredString("red", "ERROR", " Unexpected Response Status Code %v while checking token", r))
				}

			}
			var e string
			for x := 0; x < len(printStrings); x++ {
				e += printStrings[x]
			}
			color.Printf(e)

		}(i)
	}
	c.WaitAllDone()

	title <- true
	var validTokenStrings []string
	for i := 0; i < len(validTokens); i++ {
		if validTokens[i].Password != "" && validTokens[i].Email != "" {
			validTokenStrings = append(validTokenStrings, fmt.Sprintf("%v:%v:%v", validTokens[i].Email, validTokens[i].Password, validTokens[i].Token))
		} else {
			validTokenStrings = append(validTokenStrings, validTokens[i].Token)
		}
	}
	utilities.TruncateLines("tokens.txt", validTokenStrings)
	utilities.LogSuccess("All Done!")
}

func MakeColoredString(color, title, format string, args ...interface{}) string {
	return fmt.Sprintf("<fg=white>[</><fg=%v;op=bold>%v</><fg=white>]</>%s\n", color, title, fmt.Sprintf(format, args...))
}
