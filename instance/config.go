// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DirectMessage      DirectMessage      `yaml:"direct_message_settings"`
	ProxySettings      ProxySettings      `yaml:"proxy_settings"`
	ScraperSettings    ScraperSettings    `yaml:"scraper_settings"`
	CaptchaSettings    CaptchaSettings    `yaml:"captcha_settings"`
	OtherSettings      OtherSettings      `yaml:"other_settings"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"suspicion_avoidance"`
	DMonReact          DMonReact          `yaml:"dm_on_react"`
}

type DirectMessage struct {
	Delay                        int  `yaml:"individual_delay"`
	LongDelay                    int  `yaml:"rate_limit_delay"`
	Offset                       int  `yaml:"offset"`
	Skip                         bool `yaml:"skip_completed"`
	Call                         bool `yaml:"call"`
	Remove                       bool `yaml:"remove_dead_tokens"`
	RemoveM                      bool `yaml:"remove_completed_members"`
	Stop                         bool `yaml:"stop_dead_tokens"`
	Mutual                       bool `yaml:"check_mutual"`
	Friend                       bool `yaml:"friend_before_DM"`
	Websocket                    bool `yaml:"online_tokens"`
	MaxDMS                       int  `yaml:"max_dms_per_token"`
	Receive                      bool `yaml:"receive_messages"`
	SkipFailed                   bool `yaml:"skip_failed"`
	Block                        bool `yaml:"block_after_dm"`
	Close                        bool `yaml:"close_dm_after_message"`
	MultipleMessages             bool `yaml:"multiple_messages"`
	DelayBetweenMultipleMessages int  `yaml:"delay_between_multiple_messages"`
}
type ProxySettings struct {
	Proxy           string `yaml:"proxy"`
	ProxyFromFile   bool   `yaml:"proxy_from_file"`
	ProxyForCaptcha bool   `yaml:"proxy_for_captcha"`
	GatewayProxy    bool   `yaml:"use_proxy_for_gateway"`
	Timeout         int    `yaml:"timeout"`
}

type ScraperSettings struct {
	SleepSc         int    `yaml:"scraper_delay"`
	ScrapeUsernames bool   `yaml:"scrape_usernames"`
	ScrapeAvatars   bool   `yaml:"scrape_avatars"`
	ExtendedChars   string `yaml:"query_brute_extra_chars"`
}

type CaptchaSettings struct {
	ClientKey     string `yaml:"captcha_api_key"`
	CaptchaAPI    string `yaml:"captcha_api"`
	Timeout       int    `yaml:"max_captcha_wait"`
	MaxCaptchaDM  int    `yaml:"max_captcha_retry_dm"`
	MaxCaptchaInv int    `yaml:"max_captcha_retry_invite"`
	Self          string `yaml:"self"`
	SelfUsername  string `yaml:"self_username"`
	SelfPassword  string `yaml:"self_password"`
}

type OtherSettings struct {
	DisableKL        bool   `yaml:"disable_keep_alives"`
	Mode             int    `yaml:"mode"`
	ConstantCookies  bool   `yaml:"constant_cookies"`
	CensorToken      bool   `yaml:"censor_token"`
	Logs             bool   `yaml:"logs"`
	GatewayStatus    int    `yaml:"gateway_status"`
	DolfiesHeaders   bool   `yaml:"dolfies_headers"`
	XSuperProperties string `yaml:"x_super_properties"`
	Useragent        string `yaml:"useragent"`
	JA3              string `yaml:"ja3"`
}

type SuspicionAvoidance struct {
	RandomIndividualDelay  int  `yaml:"random_individual_delay"`
	RandomRateLimitDelay   int  `yaml:"random_rate_limit_delay"`
	RandomDelayOpenChannel int  `yaml:"random_delay_before_dm"`
	Typing                 bool `yaml:"typing"`
	TypingVariation        int  `yaml:"typing_variation"`
	TypingSpeed            int  `yaml:"typing_speed"`
	TypingBase             int  `yaml:"typing_base"`
}

type DMonReact struct {
	Observer              string `yaml:"observer_token"`
	ChangeName            bool   `yaml:"change_name"`
	ChangeAvatar          bool   `yaml:"change_avatar"`
	Invite                string `yaml:"invite"`
	ServerID              string `yaml:"server_id"`
	ChannelID             string `yaml:"channel_id"`
	MessageID             string `yaml:"message_id"`
	SkipCompleted         bool   `yaml:"skip_completed"`
	SkipFailed            bool   `yaml:"skip_failed"`
	LeaveTokenOnRateLimit bool   `yaml:"leave_token_on_rate_limit"`
	Emoji                 string `yaml:"emoji"`
	RotateTokens          bool   `yaml:"rotate_tokens"`
	MaxAntiRaidQueue      int    `yaml:"max_anti_raid_queue"`
	MaxDMsPerToken        int    `yaml:"max_dms_per_token"`
}

type AutoReact struct {
	Observer        string   `yaml:"observer_token"`
	Servers         []string `yaml:"servers"`
	Channels        []string `yaml:"channels"`
	Messages        []string `yaml:"messages"`
	Users           []string `yaml:"users"`
	Emojis          []string `yaml:"emojis"`
	ReactWith       []string `yaml:"react_with"`
	ReactAll        bool     `yaml:"react_all"`
	Delay           int      `yaml:"delay_between_reacts"`
	Subscribe       []string `yaml:"subscribe_to_servers"`
	Randomness      int      `yaml:"minimum_percent_react"`
	IndividualDelay int      `yaml:"individual_delay"`
}

func GetConfig() (Config, error) {
	ex, err := os.Executable()
	if err != nil {
		utilities.LogErr("Error getting executable path %v", err)
		return Config{}, err
	}
	ex = filepath.ToSlash(ex)
	var file *os.File
	file, err = os.Open(path.Join(path.Dir(ex) + "/" + "config.yml"))
	if err != nil {
		utilities.LogErr("Error while opening config file %v", err)
		return Config{}, err
	} else {
		defer file.Close()
		var config Config
		bytes, _ := io.ReadAll(file)
		err = yaml.Unmarshal(bytes, &config)
		if err != nil {
			fmt.Println(err)
			return Config{}, err
		}
		return config, nil
	}
}

func GetMessage() ([]Message, error) {
	var messages []Message
	ex, err := os.Executable()
	if err != nil {
		utilities.LogErr("Error while finding executable %v", err)
		return []Message{}, err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + "message.json"))
	if err != nil {
		utilities.LogErr("Error while Opening Message %v", err)
		fmt.Println(err)
		return []Message{}, err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	errr := json.Unmarshal(bytes, &messages)
	if errr != nil {
		fmt.Println(errr)

		return []Message{}, errr
	}

	return messages, nil
}
