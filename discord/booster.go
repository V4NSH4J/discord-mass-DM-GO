package discord

import (
	"fmt"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchBoosterInfo() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %s", err)
		return
	}

	boostAmount := 0
	for _, instance := range instances {
		boosts := instance.GetBoostSlots()
		for _, boost := range boosts {
			fmt.Printf("id: %s", boost.SubscriptionID)
			if boost.Cooldown != "" {
				fmt.Printf(" cooldown: %s", boost.Cooldown)
			} else {
				boostAmount += 1
			}
			fmt.Println()
		}
	}

	utilities.LogInfo("%d boosts available in %d accounts\n", boostAmount, len(instances))
}

func LaunchBooster() {
	_, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting config or instances %s", err)
		return
	}

	serverid := utilities.UserInput("Enter server ID:")
	boost_amount := utilities.UserInputInteger("Enter boost amount: ")

	current_boost_amount := 0
	for _, instance := range instances {
		boosts := instance.GetBoostSlots()
		for _, boost := range boosts {
			if current_boost_amount >= boost_amount {
				break
			}

			if boost.Cooldown != "" {
				continue
			}

			instance.BoostServer(serverid, boost.ID)
			current_boost_amount++
		}
	}

}
