package discord

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchFriendSpammer() {
	utilities.PrintMenu([]string{"Send Friend Requests", "Send DMs to friends", "Check Status of Friends"})
	choice := utilities.UserInputInteger("Enter your choice!")
	switch choice {
	default:
		utilities.LogErr("Invalid choice!")
		LaunchFriendSpammer()
	case 1:
		utilities.LogInfo("Enter people to be friended in memberids.txt in format username#discriminator (All scrapers make such a file in logs folder if enabled)")
		_, instances, err := instance.GetEverything()
		if err != nil {
			utilities.LogErr(err.Error())
			LaunchFriendSpammer()
		}
		friends, err := utilities.ReadLines("memberids.txt")
		if err != nil {
			utilities.LogErr(err.Error())
			LaunchFriendSpammer()
		}
		var toBeFriended []FriendsInfo
		r := regexp.MustCompile(`%s#%s`)
		for _, friend := range friends {
			if r.MatchString(friend) {
				x := strings.Split(friend, "#")
				toBeFriended = append(toBeFriended, FriendsInfo{
					Username: x[0],
					Discrim:  x[1],
				})
			}
		}
		utilities.LogInfo("%d friends to be friended", len(toBeFriended))
		var friendChan chan FriendsInfo = make(chan FriendsInfo, len(toBeFriended))
		for i := 0; i < len(toBeFriended); i++ {
			go func(i int) {
				friendChan <- toBeFriended[i]
			}(i)
		}
		var success int
		delay := utilities.UserInputInteger("Enter delay between requests (in seconds)")
		var wg sync.WaitGroup
		for i := 0; i < len(instances); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for {
					if len(friendChan) == 0 {
						break
					}
					friend := <-friendChan
					discrimInt, err := strconv.Atoi(friend.Discrim)
					if err != nil {
						utilities.LogErr(" Error %s while converting Discrim %s to an integer", err, friend.Discrim)
						continue
					}
					resp, err := instances[i].Friend(friend.Username, discrimInt)
					if err != nil {
						utilities.LogErr("Token %s Error %s while friending %s#%s", instances[i].CensorToken(), err, friend.Username, friend.Discrim)
						time.Sleep(time.Duration(delay) * time.Second)
						continue
					}
					if resp.StatusCode == 200 || resp.StatusCode == 204 {
						utilities.LogInfo("Token %s Successfully friended %s#%s", instances[i].CensorToken(), friend.Username, friend.Discrim)
						success++
					} else {
						utilities.LogFailed("Token %s Invalid Status Code %s while friending %s#%s", instances[i].CensorToken(), resp.Status, friend.Username, friend.Discrim)
					}
					resp.Body.Close()
					time.Sleep(time.Duration(delay) * time.Second)
				}

			}(i)
		}
		wg.Wait()
		utilities.LogInfo("%d friends friended!", success)
	case 2:
		utilities.LogInfo("This will check tokens which have friends and send them DMs")
		_, instances, err := instance.GetEverything()
		if err != nil {
			utilities.LogErr(err.Error())
			LaunchFriendSpammer()
		}
		delay := utilities.UserInputInteger("Enter delay between messages (in seconds)")
		var wg sync.WaitGroup
		for i := 0; i < len(instances); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				r, friendCount, blockedCount, incomingCount, outgoingCount, ids, err := instances[i].Relationships()
				if err != nil {
					utilities.LogErr("Token %s Error %s while getting relationships", instances[i].CensorToken(), err)
					return
				}
				if r != 200 && r != 204 {
					utilities.LogFailed("Token %s Invalid Status Code %s while getting relationships", instances[i].CensorToken(), r)
					return
				}
				utilities.LogSuccess("Token %s [%d friends, %d blocked, %d incoming requests, %d pending friends]", instances[i].CensorToken(), friendCount, blockedCount, incomingCount, outgoingCount)
				for x := 0; x < len(ids); x++ {
					if x != 1 {
						// Not a friend
						continue
					}
					snowflake, err := instances[i].OpenChannel(ids[x].ID)
					if err != nil {
						utilities.LogErr("Token %s Error %s while opening channel %s", instances[i].CensorToken(), err, ids[x].ID)
						continue
					}
					r, bytes, err := instances[i].SendMessage(snowflake, ids[x].ID)
					if err != nil {
						utilities.LogErr("Token %s Error %s while sending message to %s", instances[i].CensorToken(), err, ids[x].ID)
						continue
					}
					if r != 200 && r != 204 {
						utilities.LogFailed("Token %s Invalid Status Code %s while sending message to %s [%s]", instances[i].CensorToken(), r, ids[x].ID, string(bytes))
						continue
					}
					utilities.LogSuccess("Token %s Successfully sent message to %s", instances[i].CensorToken(), ids[x].ID)
					time.Sleep(time.Duration(delay) * time.Second)
				}

			}(i)
		}
		wg.Wait()
		utilities.LogInfo("All Completed!")

	case 3:
		utilities.LogInfo("This will check tokens which have friends")
		_, instances, err := instance.GetEverything()
		if err != nil {
			utilities.LogErr(err.Error())
			LaunchFriendSpammer()
		}
		var wg sync.WaitGroup
		for i := 0; i < len(instances); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				r, friendCount, blockedCount, incomingCount, outgoingCount, _, err := instances[i].Relationships()
				if err != nil {
					utilities.LogErr("Token %s Error %s while getting relationships", instances[i].CensorToken(), err)
					return
				}
				if r != 200 && r != 204 {
					utilities.LogFailed("Token %s Invalid Status Code %s while getting relationships", instances[i].CensorToken(), r)
					return
				}
				utilities.LogSuccess("Token %s [%d friends, %d blocked, %d incoming requests, %d pending friends]", instances[i].CensorToken(), friendCount, blockedCount, incomingCount, outgoingCount)

			}(i)
		}
		wg.Wait()
		utilities.LogInfo("All Completed!")

	}

}

type FriendsInfo struct {
	Username string
	Discrim  string
}
