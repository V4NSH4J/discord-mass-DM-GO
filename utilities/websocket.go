// Copyright (C) 2021 github.com/dankgrinder & github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Define WebSocket connection struct
type Connection struct {
	Members       []string
	OfflineScrape chan []byte
	AllMembers    []string
	Messages      chan []byte
	Complete      bool
	ws            *websocket.Conn
	sessionID     string
	in            chan string
	out           chan []byte
	fatalHandler  func(err error)
	seq           int
	closeChan     chan struct{}
}

// Input Discord token and start a new websocket connection
func NewConnection(token string, fatalHandler func(err error), proxy string) (*Connection, error) {
	var dialer websocket.Dialer
	if proxy == "" {
		dialer = *websocket.DefaultDialer
	} else {
		if !strings.Contains(proxy, "http://") {
			proxy = "http://" + proxy
		}
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		dialer = websocket.Dialer{Proxy: http.ProxyURL(proxyURL)}
	}
	// Dial Connection to Discord
	ws, _, err := dialer.Dial("wss://gateway.discord.gg/?v=9&encoding=json", nil)
	if err != nil {
		return nil, err
	}

	c := Connection{
		ws:            ws,
		in:            make(chan string),
		out:           make(chan []byte),
		OfflineScrape: make(chan []byte),
		Messages:      make(chan []byte),
		fatalHandler:  fatalHandler,
		closeChan:     make(chan struct{}),
	}
	// Recieve Hello message
	interval, err := c.ReadHello()
	if err != nil {
		c.ws.Close()
		return nil, err
	}
	// Authenticate with Discord
	err = c.ws.WriteJSON(&Event{
		Op: OpcodeIdentify,
		Data: Data{
			ClientState: ClientState{
				HighestLastMessageID:     "0",
				ReadStateVersion:         0,
				UserGuildSettingsVersion: -1,
			},
			Identify: Identify{
				Token: token,
				Properties: Properties{
					OS:                "Linux",
					Browser:           "Chrome",
					BrowserUserAgent:  "Chrome/86.0.4240.75",
					BrowserVersion:    "86.0.4240.75",
					Referrer:          "https://discord.com/new",
					ReferringDomain:   "discord.com",
					ReleaseChannel:    "stable",
					ClientBuildNumber: 73683,
				},
				Capabilities: 61,
				Presence: Presence{
					Status: "online",
					Since:  0,
					AFK:    false,
				},
				Compress: false,
			},
		}})
	if err != nil {
		c.ws.Close()
		return nil, fmt.Errorf("error while sending authentication message: %v", err)
	}

	if err = c.awaitEvent(EventNameReady); err != nil {
		c.ws.Close()
		return nil, fmt.Errorf("error while waiting for ready event: %v", err)
	}
	go c.Ping(time.Duration(interval) * time.Millisecond)
	go c.listen()

	return &c, nil

}

// Read Hello function to read hello message from websocket return 0 if next message is not a hello message or return the heartbeat interval
func (c *Connection) ReadHello() (int, error) {
	_, message, err := c.ws.ReadMessage()
	if err != nil {
		return 0, err
	}
	var body Event
	if err := json.Unmarshal(message, &body); err != nil {
		return 0, fmt.Errorf("Error while Unmarshalling incoming websocket message: %v", err)
	}
	if body.Op != OpcodeHello {
		return 0, fmt.Errorf("Expected OpcodeHello but got %v", body.Op)
	}

	if body.Data.HeartbeatInterval <= 0 {
		return 0, fmt.Errorf("Heartbeat interval is not valid")
	}

	return body.Data.HeartbeatInterval, nil

}

// Ping Heartbeat interval

func (c *Connection) Ping(interval time.Duration) {
	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-c.closeChan:
				return
			case <-t.C:

			}
			_ = c.ws.WriteJSON(&Event{
				Op: OpcodeHeartbeat,
			})
		}
	}()
}

func (c *Connection) awaitEvent(e string) error {
	_, b, err := c.ws.ReadMessage()
	if err != nil {
		return fmt.Errorf("error while reading message from websocket: %v", err)
	}
	var body Event
	if err = json.Unmarshal(b, &body); err != nil {
		return fmt.Errorf("error while unmarshalling incoming websocket message: %v", err)
	}
	if body.EventName != e {
		return fmt.Errorf("unexpected event name for received websocket message: %v, expected %v", body.EventName, e)
	}
	return nil
}

func (c *Connection) listen() {
	for {
		_, b, err := c.ws.ReadMessage()

		if err != nil {
			c.closeChan <- struct{}{}
			c.ws.Close()
			fmt.Println(err)
			c.fatalHandler(err)
			break
		}

		var body Event
		if err := json.Unmarshal(b, &body); err != nil {
			// All messages which don't decode properly are likely caused by the
			// data object and are ignored for now.
			continue
		}

		if body.EventName == "GUILD_MEMBERS_CHUNK" {
			go func() {
				c.OfflineScrape <- b
			}()

		}
		if body.EventName == "GUILD_MEMBER_LIST_UPDATE" {
			for i := 0; i < len(body.Data.Ops); i++ {
				if len(body.Data.Ops[i].Items) == 0 && body.Data.Ops[i].Op == "SYNC" {
					c.Complete = true
				}
			}

			for i := 0; i < len(body.Data.Ops); i++ {
				if body.Data.Ops[i].Op == "SYNC" {
					for j := 0; j < len(body.Data.Ops[i].Items); j++ {
						fmt.Println(body.Data.Ops[i].Items[j].Member.User.ID)
						c.Members = append(c.Members, body.Data.Ops[i].Items[j].Member.User.ID)
					}
				}
			}
		}

		switch body.Op {
		default:
			c.seq = body.Sequence
			if body.Data.SessionID != "" {
				c.sessionID = body.Data.SessionID
			}
			if body.EventName == EventNameMessageCreate || body.EventName == EventNameMessageUpdate {
				go func() {
					c.Messages <- b
				}()

			}
		case OpcodeInvalidSession:
			c.fatalHandler(fmt.Errorf("session invalidated"))
			c.Close()
		case OpcodeReconnect:
			c.fatalHandler(fmt.Errorf("reconnecting"))
			c.Close()

		}
	}
}

func (c *Connection) Close() error {
	c.fatalHandler = func(err error) {}
	c.closeChan <- struct{}{}
	err := c.ws.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseGoingAway, "going away"),
		time.Now().Add(time.Second*10),
	)
	if err != nil {
		c.ws.Close()
	}
	return nil
}

// Send interface to websocket
func (c *Connection) WriteRaw(e interface{}) error {
	return c.ws.WriteJSON(e)
}

// Function to write event
func (c *Connection) WriteJSONe(e *Event) error {
	return c.ws.WriteJSON(e)
}
