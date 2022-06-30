// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

func LaunchNameChanger() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
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
			color.Red("[%v] %v No password set. It may be wrongly formatted. Only supported format is email:pass:token", time.Now().Format("15:04:05"), instances[i].CensorToken())
			continue
		}
	}
	color.Red("NOTE: Names are changed randomly from the file.")
	users, err := utilities.ReadLines("names.txt")
	if err != nil {
		color.Red("[%v] Error while reading names.txt: %v", time.Now().Format("15:04:05"), err)
		utilities.ExitSafely()
	}
	color.Green("[%v] Enter number of threads: (0 for unlimited)", time.Now().Format("15:04:05"))

	var threads int
	fmt.Scanln(&threads)
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
				color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
			} else {
				color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
			}
			r, err := instances[i].NameChanger(users[rand.Intn(len(users))])
			if err != nil {
				color.Red("[%v] %v Error while changing name: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), err)
				FailedCount++
				return
			}
			body, err := utilities.ReadBody(r)
			if err != nil {
				fmt.Println(err)
			}
			if r.StatusCode == 200 || r.StatusCode == 204 {
				color.Green("[%v] %v Changed name successfully", time.Now().Format("15:04:05"), instances[i].CensorToken())
				SuccessCount++
			} else {
				color.Red("[%v] %v Error while changing name: %v %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), r.Status, string(body))
				FailedCount++
			}
			if instances[i].Ws != nil {
				if instances[i].Ws.Conn != nil {
					err = instances[i].Ws.Close()
					if err != nil {
						color.Red("[%v] Error while closing websocket: %v", time.Now().Format("15:04:05"), err)
					} else {
						color.Green("[%v] Websocket closed %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
					}
					c.Done()
				}
			}
		}(i)
	}
	c.WaitAllDone()
	title <- true
	color.Green("[%v] All Done", time.Now().Format("15:04:05"))

}

func LaunchAvatarChanger() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
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
	color.Red("NOTE: Only PNG and JPEG/JPG supported. Profile Pictures are changed randomly from the folder. Use PNG format for faster results.")
	color.White("Loading Avatars..")
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
	var avatars []string

	for i := 0; i < len(images); i++ {
		av, err := instance.EncodeImg(images[i])
		if err != nil {
			color.Red("Couldn't encode image")
			continue
		}
		avatars = append(avatars, av)
	}
	color.Green("%v avatars loaded", len(avatars))
	color.Green("[%v] Enter number of threads: ", time.Now().Format("15:04:05"))
	var threads int
	fmt.Scanln(&threads)
	if threads > len(instances) {
		threads = len(instances)
	}
	TotalCount = len(instances)
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			err := instances[i].StartWS()
			if err != nil {
				color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
			} else {
				color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
			}
			r, err := instances[i].AvatarChanger(avatars[rand.Intn(len(avatars))])
			if err != nil {
				color.Red("[%v] %v Error while changing avatar: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), err)
				FailedCount++
			} else {
				if r.StatusCode == 204 || r.StatusCode == 200 {
					color.Green("[%v] %v Avatar changed successfully", time.Now().Format("15:04:05"), instances[i].CensorToken())
					SuccessCount++
				} else {
					color.Red("[%v] %v Error while changing avatar: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), r.StatusCode)
					FailedCount++
				}
			}
			if instances[i].Ws != nil {
				if instances[i].Ws.Conn != nil {
					err = instances[i].Ws.Close()
					if err != nil {
						color.Red("[%v] Error while closing websocket: %v", time.Now().Format("15:04:05"), err)
					} else {
						color.Green("[%v] Websocket closed %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
					}
					c.Done()
				}
			}

		}(i)
	}
	c.WaitAllDone()
	title <- true
	color.Green("[%v] All done", time.Now().Format("15:04:05"))
}

func LaunchBioChanger() {
	bios, err := utilities.ReadLines("bios.txt")
	if err != nil {
		color.Red("[%v] Error while reading bios.txt: %v", time.Now().Format("15:04:05"), err)
		utilities.ExitSafely()
	}
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
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
	color.Green("[%v] Loaded %v bios, %v instances", time.Now().Format("15:04:05"), len(bios), len(instances))
	color.Green("[%v] Enter number of threads: (0 for unlimited)", time.Now().Format("15:04:05"))
	var threads int
	fmt.Scanln(&threads)
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
				color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
			} else {
				color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
			}
			err = instances[i].BioChanger(bios)
			if err != nil {
				color.Red("[%v] %v Error while changing bio: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), err)
				FailedCount++
			} else {
				color.Green("[%v] %v Bio changed successfully", time.Now().Format("15:04:05"), instances[i].CensorToken())
				SuccessCount++
			}
			if instances[i].Ws != nil {
				if instances[i].Ws.Conn != nil {
					err = instances[i].Ws.Close()
					if err != nil {
						color.Red("[%v] Error while closing websocket: %v", time.Now().Format("15:04:05"), err)
					} else {
						color.Green("[%v] Websocket closed %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
					}
					c.Done()
				}
			}
		}(i)
	}
	title <- true
	c.WaitAllDone()
}

func LaunchHypeSquadChanger() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
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
	color.Green("[%v] Enter number of threads: (0 for unlimited)", time.Now().Format("15:04:05"))
	var threads int
	fmt.Scanln(&threads)
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
				color.Red("[%v] %v Error while changing hype squad: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), err)
				FailedCount++
			} else {
				color.Green("[%v] %v Hype squad changed successfully", time.Now().Format("15:04:05"), instances[i].CensorToken())
				SuccessCount++
			}
			c.Done()
		}(i)
	}
	title <- true
	c.WaitAllDone()

}

func LaunchTokenChanger() {

	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
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
			color.Red("[%v] %v No password set. It may be wrongly formatted. Only supported format is email:pass:token", time.Now().Format("15:04:05"), instances[i].CensorToken())
			continue
		}
	}
	color.Green("[%v] Enter 0 to change passwords randomly and 1 to change them to a constant input", time.Now().Format("15:04:05"))
	var mode int
	fmt.Scanln(&mode)
	if mode != 0 && mode != 1 {
		color.Red("[%v] Invalid mode", time.Now().Format("15:04:05"))
		utilities.ExitSafely()
	}
	var password string
	if mode == 1 {
		color.Green("[%v] Enter Password:", time.Now().Format("15:04:05"))
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			password = scanner.Text()
		}
	}
	color.Green("[%v] Enter number of threads: (0 for unlimited)", time.Now().Format("15:04:05"))

	var threads int
	fmt.Scanln(&threads)
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
				color.Red("[%v] %v Error while changing token: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), err)
				FailedCount++
				err := utilities.WriteLine("input/changed_tokens.txt", fmt.Sprintf(`%s:%s:%s`, instances[i].Email, instances[i].Password, instances[i].CensorToken()))
				if err != nil {
					color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
				}
			} else {
				color.Green("[%v] %v Token changed successfully", time.Now().Format("15:04:05"), instances[i].CensorToken())
				SuccessCount++
				err := utilities.WriteLine("input/changed_tokens.txt", fmt.Sprintf(`%s:%s:%s`, instances[i].Email, password, newToken))
				if err != nil {
					color.Red("[%v] Error while writing to file: %v", time.Now().Format("15:04:05"), err)
				}
			}
			c.Done()
		}(i)
	}
	c.WaitAllDone()
	title <- true
	color.Green("[%v] All Done", time.Now().Format("15:04:05"))

}

func LaunchServerNicknameChanger() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("[%v] Error while getting necessary data: %v", time.Now().Format("15:04:05"), err)
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
	color.Red("NOTE: Nicknames are changed randomly from the file.")
	nicknames, err := utilities.ReadLines("nicknames.txt")
	if err != nil {
		color.Red("[%v] Error while reading nicknames.txt: %v", time.Now().Format("15:04:05"), err)
		utilities.ExitSafely()
	}

	var guildid int
	color.Green("[%v] Enter guild id in which nicknames should be changed", time.Now().Format("15:04:05"))
	fmt.Scanln(&guildid)

	color.Green("[%v] Enter number of threads: (0 for unlimited)", time.Now().Format("15:04:05"))
	var threads int
	fmt.Scanln(&threads)
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
				color.Red("[%v] %v Error while changing nickname: %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), err)
				FailedCount++
				return
			}
			body, err := utilities.ReadBody(r)
			if err != nil {
				fmt.Println(err)
			}
			if r.StatusCode == 200 || r.StatusCode == 204 {
				color.Green("[%v] %v Changed nickname successfully", time.Now().Format("15:04:05"), instances[i].CensorToken())
				SuccessCount++
			} else {
				color.Red("[%v] %v Error while changing nickname: %v %v", time.Now().Format("15:04:05"), instances[i].CensorToken(), r.Status, string(body))
				FailedCount++
			}
			c.Done()
		}(i)
	}
	c.WaitAllDone()
	title <- true
	color.Green("[%v] All Done", time.Now().Format("15:04:05"))

}
