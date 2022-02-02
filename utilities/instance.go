// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type Instance struct {
	Token       string
	Password    string
	Proxy       string
	Cookie      string
	Fingerprint string
	Messages    []Message
	Count       int
	LastQuery   string
	LastCount   int
	Members     []User
	AllMembers  []User
	Rejoin      int
	ScrapeCount int
	ID          string
	Receiver    bool

	Client *http.Client
	WG     *sync.WaitGroup
	Ws     *Connection
	fatal  chan error
}

func (in *Instance) StartWS() error {
	ws, err := NewConnection(in.Token, in.wsFatalHandler, in.Proxy)
	if err != nil {
		return fmt.Errorf("failed to create websocket connection: %s", err)
	}
	in.Ws = ws
	return nil
}

func (in *Instance) wsFatalHandler(err error) {
	if closeErr, ok := err.(*websocket.CloseError); ok && closeErr.Code == 4004 {
		in.fatal <- fmt.Errorf("websocket closed: authentication failed, try using a new token")
		return
	}
	color.Red("Websocket closed %v %v", err, in.Token)
	in.Receiver = false
	in.Ws, err = NewConnection(in.Token, in.wsFatalHandler, in.Proxy)
	if err != nil {
		in.fatal <- fmt.Errorf("failed to create websocket connection: %s", err)
		return
	}
	color.Green("Reconnected To Websocket")

}

type CallEvent struct {
	Op   int      `json:"op"`
	Data CallData `json:"d"`
}

type CallData struct {
	ChannelId string      `json:"channel_id"`
	GuildId   interface{} `json:"guild_id"`
	SelfDeaf  bool        `json:"self_deaf"`
	SelfMute  bool        `json:"self_mute"`
	SelfVideo bool        `json:"self_video"`
}

func (in *Instance) Call(snowflake string) error {
	if in.Ws == nil {
		return fmt.Errorf("websocket is not initialized")
	}
	e := CallEvent{
		Op: 4,
		Data: CallData{
			ChannelId: snowflake,
			GuildId:   nil,
			SelfDeaf:  false,
			SelfMute:  false,
			SelfVideo: false,
		},
	}
	err := in.Ws.WriteRaw(e)
	if err != nil {
		return fmt.Errorf("failed to write to websocket: %s", err)
	}

	return nil
}
