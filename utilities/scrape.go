// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities


func Scrape(ws *Connection, Guild string, Channel string, index int) error {
	var x []interface{}
	if index == 0 {
		x = []interface{}{[2]int{0, 99}}
	} else if index == 1 {
		x = []interface{}{[2]int{0, 99}, [2]int{100, 199}}
	} else if index == 2 {
		x = []interface{}{[2]int{0, 99}, [2]int{100, 199}, [2]int{200, 299}}
	} else {
		x = []interface{}{[2]int{0, 99}, [2]int{100, 199}, [2]int{index * 100, (index * 100) + 99}}
	}

	payload := Data{
		GuildId: Guild,
		Channels: map[string]interface{}{
			Channel: x,
		},
	}

	err := ws.WriteJSONe(&Event{
		Op:   14,
		Data: payload,
	})

	if err != nil {
		return err
	}

	return nil

}

type CustomEvent struct {
	Op   int    `json:"op,omitempty"`
	Data Custom `json:"d,omitempty"`
}
type Custom struct {
	GuildID  interface{} `json:"guild_id"`
	Limit    int         `json:"limit"`
	Query    string      `json:"query"`
	Presence bool        `json:"presence"`
}

// Write a function which would input the connection, guildid, query, limit and presence.
// The function would then make an Event struct and send it to the websocket.
// The guild ID is to be put as a list of one item

func ScrapeOffline(c *Connection, guild string, query string) error {

	custom := Custom{
		GuildID:  []string{guild},
		Limit:    100,
		Query:    query,
		Presence: true,
	}
	eventx := CustomEvent{
		Op:   8,
		Data: custom,
	}


	err := c.ws.WriteJSON(eventx)
	if err != nil {
		return err
	}
	err = c.WriteRaw(eventx)
	if err != nil {
		return err
	}

	return nil
}
