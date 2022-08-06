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
	"os/exec"
	"path"
	"path/filepath"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/zenthangplus/goccm"
)

func LaunchNameChanger() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %v", err)
	}
	var TotalCount, SuccessCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Success, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount-SuccessCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	for i := 0; i < len(instances); i++ {
		if instances[i].Password == "" {
			utilities.LogWarn("Token %v does not have password set. Name changer requires token in format email:password:token", instances[i].CensorToken())
			continue
		}
	}
	utilities.LogWarn("Usernames are changed randomly from file.")
	users, err := utilities.ReadLines("names.txt")
	if err != nil {
		utilities.LogErr("Error while reading names.txt: %v", err)
		return
	}
	if len(users) == 0 {
		utilities.LogErr("names.txt is empty")
		return
	}
	threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	TotalCount = len(instances)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			err := instances[i].StartWS()
			if err != nil {
				utilities.LogErr("Token %v Error while opening websocket %v", instances[i].CensorToken(), err)
			} else {
				utilities.LogSuccess("Token %v websocket open", instances[i].CensorToken())
			}
			r, err := instances[i].NameChanger(users[rand.Intn(len(users))])
			if err != nil {
				utilities.LogErr("Token %v Error while changing name: %v", instances[i].CensorToken(), err)
				FailedCount++
				return
			}
			body, err := utilities.ReadBody(r)
			if err != nil {
				utilities.LogErr("Token %v Error while reading body: %v", instances[i].CensorToken(), err)
				FailedCount++
				return
			}
			if r.StatusCode == 200 || r.StatusCode == 204 {
				utilities.LogSuccess("Token %v Name changed successfully", instances[i].CensorToken())
				SuccessCount++
			} else {
				utilities.LogFailed("Token %v Error while changing name: %v %v", instances[i].CensorToken(), r.Status, string(body))
				FailedCount++
			}
			if instances[i].Ws != nil {
				if instances[i].Ws.Conn != nil {
					err = instances[i].Ws.Close()
					if err != nil {
						utilities.LogFailed("Token %v Error while closing websocket: %v", instances[i].CensorToken(), err)
					} else {
						utilities.LogSuccess("Token %v websocket closed", instances[i].CensorToken())
					}
					c.Done()
				}
			}
		}(i)
	}
	c.WaitAllDone()
	title <- true
	utilities.LogSuccess("Name changer finished")

}

func LaunchAvatarChanger() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %v", err)
	}
	var TotalCount, SuccessCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Success, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount-SuccessCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	utilities.LogWarn("NOTE: Only PNG and JPEG/JPG supported. Profile Pictures are changed randomly from the folder. Use PNG format for faster results.")
	utilities.LogInfo("Loading Avatars..")
	ex, err := os.Executable()
	if err != nil {
		utilities.LogErr("Error while getting executable path: %v", err)
		utilities.ExitSafely()
	}
	ex = filepath.ToSlash(ex)
	path := path.Join(path.Dir(ex) + "/input/pfps")

	images, err := instance.GetFiles(path)
	if err != nil {
		utilities.LogErr("Error while getting files from %v: %v", path, err)
		utilities.ExitSafely()
	}
	utilities.LogInfo("%v files loaded", len(images))
	var avatars []string

	for i := 0; i < len(images); i++ {
		av, err := instance.EncodeImg(images[i])
		if err != nil {
			utilities.LogErr("Error while encoding image %v: %v", images[i], err)
			continue
		}
		avatars = append(avatars, av)
	}
	utilities.LogInfo("%v avatars loaded", len(avatars))
	threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	TotalCount = len(instances)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			err := instances[i].StartWS()
			if err != nil {
				utilities.LogFailed("Token %v Error while opening websocket", instances[i].CensorToken())
			} else {
				utilities.LogSuccess("Websocket opened %v", instances[i].CensorToken())
			}
			r, err := instances[i].AvatarChanger(avatars[rand.Intn(len(avatars))])
			if err != nil {
				utilities.LogFailed("Token %v Error while changing avatar: %v", instances[i].CensorToken(), err)
				FailedCount++
			} else {
				if r.StatusCode == 204 || r.StatusCode == 200 {
					utilities.LogSuccess("Token %v Avatar changed successfully", instances[i].CensorToken())
					SuccessCount++
				} else {
					utilities.LogFailed("Token %v Error while changing avatar: %v", instances[i].CensorToken(), r.StatusCode)
					FailedCount++
				}
			}
			if instances[i].Ws != nil {
				if instances[i].Ws.Conn != nil {
					err = instances[i].Ws.Close()
					if err != nil {
						utilities.LogFailed("Token %v Error while closing websocket: %v", instances[i].CensorToken(), err)
					} else {
						utilities.LogSuccess("Token %v websocket closed", instances[i].CensorToken())
					}
					c.Done()
				}
			}

		}(i)
	}
	c.WaitAllDone()
	title <- true
	utilities.LogSuccess("Avatar changer finished")
}

func LaunchBioChanger() {
	bios, err := utilities.ReadLines("bios.txt")
	if err != nil {
		utilities.LogErr("Error while reading bios.txt: %v", err)
		utilities.ExitSafely()
	}
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %v", err)
		utilities.ExitSafely()
	}
	var TotalCount, SuccessCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Success, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount-SuccessCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	bios = instance.ValidateBios(bios)
	utilities.LogInfo("Loaded %v bios, %v instances", len(bios), len(instances))
	threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	TotalCount = len(instances)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			err := instances[i].StartWS()
			if err != nil {
				utilities.LogFailed("Token %v Error while opening websocket", instances[i].CensorToken())
			} else {
				utilities.LogSuccess("Token %v Websocket opened", instances[i].CensorToken())
			}
			err = instances[i].BioChanger(bios)
			if err != nil {
				utilities.LogFailed("%v Error while changing bio: %v", instances[i].CensorToken(), err)
				FailedCount++
			} else {
				utilities.LogSuccess("%v Bio changed successfully", instances[i].CensorToken())
				SuccessCount++
			}
			if instances[i].Ws != nil {
				if instances[i].Ws.Conn != nil {
					err = instances[i].Ws.Close()
					if err != nil {
						utilities.LogFailed("Token %v Error while closing websocket: %v", instances[i].CensorToken(), err)
					} else {
						utilities.LogSuccess("Token %v Websocket closed", instances[i].CensorToken())
					}
					c.Done()
				}
			}
		}(i)
	}
	title <- true
	c.WaitAllDone()
	utilities.LogSuccess("Bio changer finished")
}

func LaunchHypeSquadChanger() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %v", err)
		utilities.ExitSafely()
	}
	var TotalCount, SuccessCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Success, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount-SuccessCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	TotalCount = len(instances)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			err := instances[i].RandomHypeSquadChanger()
			if err != nil {
				utilities.LogFailed("Token %v Error while changing hype squad: %v", instances[i].CensorToken(), err)
				FailedCount++
			} else {
				utilities.LogSuccess("Token %v Hype squad changed successfully", instances[i].CensorToken())
				SuccessCount++
			}
			c.Done()
		}(i)
	}
	title <- true
	c.WaitAllDone()
	utilities.LogSuccess("Hype squad changer finished")
}

func LaunchTokenChanger() {

	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %v", err)
	}
	var TotalCount, SuccessCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Changed, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount-SuccessCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	for i := 0; i < len(instances); i++ {
		if instances[i].Password == "" {
			utilities.LogWarn("%v No password set. It may be wrongly formatted. Only supported format is email:pass:token", instances[i].CensorToken())
			continue
		}
	}
	mode := utilities.UserInputInteger("Enter 0 to change passwords randomly and 1 to change them to a constant input")

	if mode != 0 && mode != 1 {
		utilities.LogErr("Invalid mode")
		utilities.ExitSafely()
	}
	var password string
	if mode == 1 {
		password = utilities.UserInput("Enter password to change tokens to:")
	}
	threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	TotalCount = len(instances)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			if password == "" {
				password = utilities.RandStringBytes(12)
			}
			newToken, err := instances[i].ChangeToken(password)
			if err != nil {
				utilities.LogFailed("Token %v Error while changing token: %v", instances[i].CensorToken(), err)
				FailedCount++
				err := utilities.WriteLine("input/changed_tokens.txt", fmt.Sprintf(`%s:%s:%s`, instances[i].Email, instances[i].Password, instances[i].Token))
				if err != nil {
					utilities.LogErr("Error while writing to file: %v", err)
				}
			} else {
				utilities.LogSuccess("%v Token changed successfully", instances[i].CensorToken())
				SuccessCount++
				err := utilities.WriteLine("input/changed_tokens.txt", fmt.Sprintf(`%s:%s:%s`, instances[i].Email, password, newToken))
				if err != nil {
					utilities.LogErr("Error while writing to file: %v", err)
				}
			}
			c.Done()
		}(i)
	}
	c.WaitAllDone()
	title <- true
	utilities.LogSuccess("Token changer finished")

}

func LaunchServerNicknameChanger() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %v", err)
	}
	var TotalCount, SuccessCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Success, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount-SuccessCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	utilities.LogWarn("NOTE: Nicknames are changed randomly from the file.")
	nicknames, err := utilities.ReadLines("nicknames.txt")
	if err != nil {
		utilities.LogErr("Error while reading nicknames.txt: %v", err)
		utilities.ExitSafely()
	}

	guildid := utilities.UserInput("Enter guild ID:")

	threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	TotalCount = len(instances)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			r, err := instances[i].NickNameChanger(nicknames[rand.Intn(len(nicknames))], guildid)
			if err != nil {
				utilities.LogFailed("Token %v Error while changing nickname: %v", instances[i].CensorToken(), err)
				FailedCount++
				return
			}
			body, err := utilities.ReadBody(r)
			if err != nil {
				fmt.Println(err)
			}
			if r.StatusCode == 200 || r.StatusCode == 204 {
				utilities.LogSuccess("Token %v Changed nickname successfully", instances[i].CensorToken())
				SuccessCount++
			} else {
				utilities.LogFailed("Token %v Error while changing nickname: %v %v", instances[i].CensorToken(), r.Status, string(body))
				FailedCount++
			}
			c.Done()
		}(i)
	}
	c.WaitAllDone()
	title <- true
	utilities.LogSuccess("All Done")

}

func LaunchFriendRequestSpammer() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %v", err)
		return
	}
	var TotalCount, SuccessCount, FailedCount int
	title := make(chan bool)
	go func() {
	Out:
		for {
			select {
			case <-title:
				break Out
			default:
				cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf(`DMDGO [%v Success, %v Failed, %v Unprocessed]`, SuccessCount, FailedCount, TotalCount-SuccessCount-FailedCount))
				_ = cmd.Run()
			}

		}
	}()
	threads := utilities.UserInputInteger("Enter number of threads (0 for maximum):")
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	username := utilities.UserInput("Enter username to spam (Only Username, not Discrim):")
	discrim := utilities.UserInputInteger("Enter discriminator to spam (Only Discrim, not Username):")
	TotalCount = len(instances)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			defer c.Done()
			r, err := instances[i].Friend(username, discrim)
			if err != nil {
				utilities.LogFailed("Token %v Error while sending friend request: %v", instances[i].CensorToken(), err)
				FailedCount++
				return
			}
			body, err := utilities.ReadBody(*r)
			if err != nil {
				utilities.LogErr("Error while reading body: %v", err)
				FailedCount++
				return
			}
			if r.StatusCode == 200 || r.StatusCode == 204 {
				utilities.LogSuccess("Token %v Sent friend request successfully", instances[i].CensorToken())
				SuccessCount++
			} else {
				utilities.LogFailed("Token %v Error while sending friend request: %v %v", instances[i].CensorToken(), r.Status, string(body))
				FailedCount++
			}
		}(i)
	}
	c.WaitAllDone()
	title <- true
	utilities.LogSuccess("All Done")
}
