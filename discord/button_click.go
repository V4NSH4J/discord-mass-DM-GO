package discord

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/fatih/color"
	"github.com/zenthangplus/goccm"
)

func LaunchButtonClicker() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		color.Red("Error: %s", err)
		return 
	}
	color.Cyan("Enter a token which can see the message:")
	var token string
	fmt.Scanln(&token)
	color.White("Enter message ID: ")
	var id string
	fmt.Scanln(&id)
	color.White("Enter channel ID: ")
	var channel string
	fmt.Scanln(&channel)
	color.White("Enter Server ID: ")
	var server string
	fmt.Scanln(&server)
	msg, err := instance.FindMessage(channel, id, token)
	if err != nil {
		color.Red("Error while finding message: %v", err)
		return
	}
	color.Green("[%v] Message: %v", time.Now().Format("15:04:05"), msg)
	var Msg instance.Message 
	err = json.Unmarshal([]byte(msg), &Msg)
	if err != nil {
		color.Red("Error while marshalling message: %v", err)
		return
	}
	if len(Msg.Components) == 0 {
		color.Red("No buttons found in message")
		return
	}
	color.Cyan("Enter Row number: ")
	for i := 0; i < len(Msg.Components); i++ {
		fmt.Println(fmt.Sprintf("%v) Row %v", i, i))
	}
	var row int
	fmt.Scanln(&row)
	color.Cyan("Select Button:")
	var column int 
	for i := 0; i < len(Msg.Components[row].Buttons); i++ {
		if Msg.Components[row].Buttons[i].Label != "" {
			fmt.Println(fmt.Sprintf("%v) Button %v [%v]", i, i, Msg.Components[row].Buttons[i].Label))
		} else if Msg.Components[row].Buttons[i].Emoji.Name != "" {
			fmt.Println(fmt.Sprintf("%v) Button %v [%v]", i, i, Msg.Components[row].Buttons[i].Emoji))
		} else {
			fmt.Println(fmt.Sprintf("%v) Button %v [Name or Emoji not found]", i, i))
		}
	}
	fmt.Scanln(&column)
	color.Cyan("Enter number of threads")
	var threads int
	fmt.Scanln(&threads)
	if threads > len(instances) || threads == 0 {
		threads = len(instances)
	}
	c := goccm.New(threads)
	for i := 0; i < len(instances); i++ {
		c.Wait()
		go func(i int) {
			defer c.Done()
			err := instances[i].StartWS()
			if err != nil {
				color.Red("[%v] Error while opening websocket: %v", time.Now().Format("15:04:05"), err)
			} else {
				color.Green("[%v] Websocket opened %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
			}
			respCode, err := instances[i].PressButton(row, column, server,Msg)
			if err != nil {
				color.Red("Error while pressing button: %v", err)
				return
			}
			if respCode != 204 && respCode != 200 {
				color.Red("Error while pressing button: %v", respCode)
				return
			}
			color.Green("[%v] Button pressed on instance %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
			if instances[i].Ws != nil {
				if instances[i].Ws.Conn != nil {
					err = instances[i].Ws.Close()
					if err != nil {
						color.Red("[%v] Error while closing websocket: %v", time.Now().Format("15:04:05"), err)
					} else {
						color.Green("[%v] Websocket closed %v", time.Now().Format("15:04:05"), instances[i].CensorToken())
					}
				}
			}
		}(i)
	}
	c.WaitAllDone()


	
}