// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

type twoCaptchaSubmitResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

type CapmonsterPayload struct {
	ClientKey string `json:"clientKey,omitempty"`
	Task      Task   `json:"task,omitempty"`
	TaskId    int    `json:"taskId,omitempty"`
	SoftID    int    `json:"softId,omitempty"`
}

type Task struct {
	CaptchaType   string     `json:"type,omitempty"`
	WebsiteURL    string     `json:"websiteURL,omitempty"`
	WebsiteKey    string     `json:"websiteKey,omitempty"`
	IsInvisible   bool       `json:"isInvisible,omitempty"`
	Data          string     `json:"data,omitempty"`
	ProxyType     string     `json:"proxyType,omitempty"`
	ProxyAddress  string     `json:"proxyAddress,omitempty"`
	ProxyPort     int        `json:"proxyPort,omitempty"`
	ProxyLogin    string     `json:"proxyLogin,omitempty"`
	ProxyPassword string     `json:"proxyPassword,omitempty"`
	UserAgent     string     `json:"userAgent,omitempty"`
	Cookies       string     `json:"cookies,omitempty"`
	Enterprise    Enterprise `json:"enterprisePayload,omitempty"`
}

type Enterprise struct {
	RqData      string `json:"rqdata,omitempty"`
	Sentry      bool   `json:"sentry,omitempty"`
	ApiEndpoint string `json:"apiEndpoint,omitempty"`
	Endpoint    string `json:"endpoint,omitempty"`
	ReportAPI   string `json:"reportapi,omitempty"`
	AssetHost   string `json:"assethost,omitempty"`
	ImageHost   string `json:"imghost,omitempty"`
}

type CapmonsterSubmitResponse struct {
	ErrorID int `json:"errorId,omitempty"`
	TaskID  int `json:"taskId,omitempty"`
}

type CapmonsterOutResponse struct {
	ErrorID   int      `json:"errorId,omitempty"`
	ErrorCode string   `json:"errorCode,omitempty"`
	Status    string   `json:"status,omitempty"`
	Solution  Solution `json:"solution"`
}

type Solution struct {
	CaptchaResponse string `json:"gRecaptchaResponse,omitempty"`
}

type UserInf struct {
	User   User     `json:"user"`
	Mutual []Guilds `json:"mutual_guilds"`
}

type Guilds struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
}

type captchaDetected struct {
	CaptchaKey []string `json:"captcha_key"`
	Sitekey    string   `json:"captcha_sitekey"`
	Service    string   `json:"captcha_service"`
	RqData     string   `json:"captcha_rqdata"`
	RqToken    string   `json:"captcha_rqtoken"`
}

type Reactionx struct {
	ID string `json:"id"`
}

type guild struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type joinresponse struct {
	VerificationForm bool  `json:"show_verification_form"`
	GuildObj         guild `json:"guild"`
}

type bypassInformation struct {
	Version    string      `json:"version"`
	FormFields []FormField `json:"form_fields"`
}

type FormField struct {
	FieldType   string   `json:"field_type"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Values      []string `json:"values"`
	Response    bool     `json:"response"`
}

type XContext struct {
	Location            string  `json:"location"`
	LocationGuildID     string  `json:"location_guild_id"`
	LocationChannelID   string  `json:"location_channel_id"`
	LocationChannelType float64 `json:"location_channel_type"`
}

type RingData struct {
	Recipients interface{} `json:"recipients"`
}

type invitePayload struct {
	CaptchaKey string `json:"captcha_key,omitempty"`
	RqToken    string `json:"captcha_rqtoken,omitempty"`
}

type friendRequest struct {
	Username string `json:"username"`
	Discrim  int    `json:"discriminator"`
}

type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}
type Embed struct {
	Title string `json:"title,omitempty"`

	// The type of embed. Always EmbedTypeRich for webhook embeds.
	Type        string             `json:"type,omitempty"`
	Description string             `json:"description,omitempty"`
	URL         string             `json:"url,omitempty"`
	Image       *MessageEmbedImage `json:"image,omitempty"`

	// The color code of the embed.
	Color     int                    `json:"color,omitempty"`
	Footer    EmbedFooter            `json:"footer,omitempty"`
	Thumbnail *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Provider  EmbedProvider          `json:"provider,omitempty"`
	Author    EmbedAuthor            `json:"author,omitempty"`
	Fields    []EmbedField           `json:"fields,omitempty"`
}
type Emoji struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Animated bool   `json:"animated,omitempty"`
}
type Reaction struct {
	Emojis Emoji `json:"emoji,omitempty"`
	Count  int   `json:"count,omitempty"`
}

type Message struct {
	Content    string             `json:"content,omitempty"`
	ChannelID  string             `json:"channel_id,omitempty"`
	Embeds     []Embed            `json:"embeds,omitempty"`
	Reactions  []Reaction         `json:"reactions,omitempty"`
	Author     User               `json:"author,omitempty"`
	GuildID    string             `json:"guild_id,omitempty"`
	MessageId  string             `json:"id,omitempty"`
	Components []MessageComponent `json:"components,omitempty"`
	Flags      int                `json:"flags,omitempty"`
}

type MessageComponent struct {
	Type    int       `json:"type"`
	Buttons []Buttons `json:"components"`
}

type Buttons struct {
	Type     int         `json:"type,omitempty"`
	Style    int         `json:"style,omitempty"`
	Label    string      `json:"label,omitempty"`
	CustomID string      `json:"custom_id,omitempty"`
	Hash     string      `json:"hash,omitempty"`
	Emoji    ButtonEmoji `json:"emoji,omitempty"`
	Disabled bool        `json:"disabled,omitempty"`
}

type ButtonEmoji struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Animated bool   `json:"animated,omitempty"`
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

type NameChange struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AvatarChange struct {
	Avatar string `json:"avatar"`
}

type TokenInfo struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Avatar            string `json:"avatar"`
	AvatarDecoration  string `json:"avatar_decoration"`
	Discriminator     string `json:"discriminator"`
	PublicFlags       int    `json:"public_flags"`
	PremiumUsageFlags int    `json:"premium_usage_flags"`
	PremiumType       int    `json:"premium_type"`
	Flags             int    `json:"flags"`
	Bio               string `json:"bio"`
	Pronouns          string `json:"pronouns"`
	Locale            string `json:"locale"`
	NSFWAllowed       bool   `json:"nsfw_allowed"`
	MFAEnabled        bool   `json:"mfa_enabled"`
	Email             string `json:"email"`
	Verified          bool   `json:"verified"`
	Phone             string `json:"phone"`
}

type ReactInfo struct {
	ChannelID string
	MessageID string
	Emoji     string
}

type NickNameChange struct {
	Nickname string `json:"nick"`
}
type Activity struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type PresenceChange struct {
	Since      int        `json:"since,omitempty"`
	Activities []Activity `json:"activities"`
	Status     string     `json:"status"`
	Afk        bool       `json:"afk"`
}

type CfBm struct {
	M       string   `json:"m"`
	Results []string `json:"results"`
	Timing  int      `json:"timing"` // Time taken
	Fp      struct { // Fingerprint
		ID int      `json:"id"` // ID
		E  struct { // Engine
			R  []int   `json:"r"`  // Screen Width, Screen Height (Total)
			Ar []int   `json:"ar"` // Available screen Width, Available screen Height
			Pr float64 `json:"pr"` // Pixel ratio
			Cd int     `json:"cd"` // Color depth
			Wb bool    `json:"wb"`
			Wp bool    `json:"wp"`
			Wn bool    `json:"wn"`
			Ch bool    `json:"ch"` // Chrome browser
			Ws bool    `json:"ws"`
			Wd bool    `json:"wd"`
		} `json:"e"`
	} `json:"fp"`
}

type Fingerprints struct {
	JA3              string `json:"ja3"`
	XSuperProperties string `json:"x-super-properties"`
	Useragent        string `json:"useragent"`
}
