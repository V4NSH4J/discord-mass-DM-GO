// Copyright (C) github.com/dankgrinder & github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import "time"

const (
	OpcodeDispatch = iota
	OpcodeHeartbeat
	OpcodeIdentify
	OpcodePresenceUpdate
	OpcodeVoiceStateUpdate
	OpcodeResume = iota + 1
	OpcodeReconnect
	OpcodeRequestGuildMembers
	OpcodeInvalidSession
	OpcodeHello
	OpcodeHeartbeatACK
)

const (
	EventNameMessageCreate = "MESSAGE_CREATE"
	EventNameMessageUpdate = "MESSAGE_UPDATE"
	EventNameReady         = "READY"
	EventNameResumed       = "RESUMED"
)

type Intent int

// Constants for the different bit offsets of intents
const (
	IntentsGuilds                 Intent = 1 << 0
	IntentsGuildMembers           Intent = 1 << 1
	IntentsGuildBans              Intent = 1 << 2
	IntentsGuildEmojis            Intent = 1 << 3
	IntentsGuildIntegrations      Intent = 1 << 4
	IntentsGuildWebhooks          Intent = 1 << 5
	IntentsGuildInvites           Intent = 1 << 6
	IntentsGuildVoiceStates       Intent = 1 << 7
	IntentsGuildPresences         Intent = 1 << 8
	IntentsGuildMessages          Intent = 1 << 9
	IntentsGuildMessageReactions  Intent = 1 << 10
	IntentsGuildMessageTyping     Intent = 1 << 11
	IntentsDirectMessages         Intent = 1 << 12
	IntentsDirectMessageReactions Intent = 1 << 13
	IntentsDirectMessageTyping    Intent = 1 << 14

	IntentsAllWithoutPrivileged = IntentsGuilds |
		IntentsGuildBans |
		IntentsGuildEmojis |
		IntentsGuildIntegrations |
		IntentsGuildWebhooks |
		IntentsGuildInvites |
		IntentsGuildVoiceStates |
		IntentsGuildMessages |
		IntentsGuildMessageReactions |
		IntentsGuildMessageTyping |
		IntentsDirectMessages |
		IntentsDirectMessageReactions |
		IntentsDirectMessageTyping
	IntentsAll = IntentsAllWithoutPrivileged |
		IntentsGuildMembers |
		IntentsGuildPresences
	IntentsNone Intent = 0
)

type Event struct {
	Op        int    `json:"op"`
	Data      Data   `json:"d,omitempty"`
	Sequence  int    `json:"s,omitempty"`
	EventName string `json:"t,omitempty"`
}

type Data struct {
	Message
	Identify
	PresenceChange
	ClientState       ClientState            `json:"client_state,omitempty"`
	HeartbeatInterval int                    `json:"heartbeat_interval,omitempty"`
	SessionID         string                 `json:"session_id,omitempty"`
	Sequence          int                    `json:"seq,omitempty"` // For sending only
	GuildId           interface{}            `json:"guild_id,omitempty"`
	Channels          map[string]interface{} `json:"channels,omitempty"`
	Ops               []Ops                  `json:"ops,omitempty"`
	ChannelID         string                 `json:"channel_id,omitempty"`
	Members           []Member               `json:"members,omitempty"`
	Typing            bool                   `json:"typing,omitempty"`
	Threads           bool                   `json:"threads,omitempty"`
	Activities        bool                   `json:"activities,omitempty"`
	ThreadMemberLists interface{}            `json:"thread_member_lists,omitempty"`
	// Emoji React
	UserID    string `json:"user_id,omitempty"`
	Emoji     Emoji  `json:"emoji,omitempty"`
	MessageID string `json:"message_id,omitempty"`
}

type Ops struct {
	Items []Userinfo  `json:"items,omitempty"`
	Range interface{} `json:"range,omitempty"`
	Op    string      `json:"op,omitempty"`
}
type Userinfo struct {
	Member Member `json:"member,omitempty"`
}
type Member struct {
	User                       User        `json:"user,omitempty"`
	Roles                      []string    `json:"roles"`
	PremiumSince               interface{} `json:"premium_since"`
	Pending                    bool        `json:"pending"`
	Nick                       string      `json:"nick"`
	Mute                       bool        `json:"mute"`
	JoinedAt                   time.Time   `json:"joined_at"`
	Flags                      int         `json:"flags"`
	Deaf                       bool        `json:"deaf"`
	CommunicationDisabledUntil interface{} `json:"communication_disabled_until"`
	Avatar                     interface{} `json:"avatar"`
}
type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Bot           bool   `json:"bot"`
}

type Identify struct {
	Token        string     `json:"token,omitempty"`
	Properties   Properties `json:"properties,omitempty"`
	Capabilities int        `json:"capabilities,omitempty"`
	Compress     bool       `json:"compress,omitempty"`
	Presence     Presence   `json:"presence,omitempty"`
}

type Presence struct {
	Status     string   `json:"status,omitempty"`
	Since      int      `json:"since,omitempty"`
	Activities []string `json:"activities,omitempty"`
	AFK        bool     `json:"afk,omitempty"`
}

type ClientState struct {
	HighestLastMessageID     string `json:"highest_last_message_id,omitempty"`
	ReadStateVersion         int    `json:"read_state_version,omitempty"`
	UserGuildSettingsVersion int    `json:"user_guild_settings_version,omitempty"`
}

type Properties struct {
	OS                     string `json:"os,omitempty"`
	Browser                string `json:"browser,omitempty"`
	Device                 string `json:"device,omitempty"`
	BrowserUserAgent       string `json:"browser_user_agent,omitempty"`
	BrowserVersion         string `json:"browser_version,omitempty"`
	OSVersion              string `json:"os_version,omitempty"`
	Referrer               string `json:"referrer,omitempty"`
	ReferringDomain        string `json:"referring_domain,omitempty"`
	ReferrerCurrent        string `json:"referrer_current,omitempty"`
	ReferringDomainCurrent string `json:"referring_domain_current,omitempty"`
	ReleaseChannel         string `json:"release_channel,omitempty"`
	ClientBuildNumber      int    `json:"client_build_number,omitempty"`
}
