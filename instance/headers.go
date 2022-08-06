// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"

	http "github.com/Danny-Dasilva/fhttp"
)

func (in *Instance) cookieHeaders(req *http.Request) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":            "discord.com",
			"User-Agent":      "Discord/125.0 (iPad; iOS 15.4.1; Scale/2.00)",
			"Accept-Language": "en-US;q=1",
			"Accept":          "*/*",
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"accept-language":           "en-US,en;q=0.9",
			"sec-ch-ua-mobile":          "?0",
			"sec-fetch-dest":            "document",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-site":            "none",
			"sec-fetch-user":            "?1",
			"upgrade-insecure-requests": "1",
			"user-agent":                in.UserAgent,
		} {
			req.Header.Set(k, v)
		}
	}

	return req
}

func (in *Instance) cfBmHeaders(req *http.Request, cookie string) *http.Request {
	for k, v := range map[string]string{
		"accept":          "*/*",
		"accept-language": "en-US,en;q=0.9",
		"content-type":    "application/json",
		"cookie":          cookie,
		"origin":          "https://discord.com",
		"referer":         "https://discord.com/",
		"sec-fetch-mode":  "cors",
		"sec-fetch-site":  "same-origin",
		"user-agent":      in.UserAgent,
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) inviteHeaders(req *http.Request, cookie, xcontext string) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":                 "discord.com",
			"Cookie":               cookie,
			"Content-Type":         "application/json",
			"X-Debug-Options":      "bugReporterEnabled",
			"Accept":               "*/*",
			"Authorization":        in.Token,
			"X-Discord-Locale":     "en-US",
			"Accept-Language":      "en-US",
			"X-Context-Properties": xcontext,
			"User-Agent":           in.UserAgent,
			"X-Super-Properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"accept":               "*/*",
			"accept-language":      "en-US,en;q=0.9",
			"authorization":        in.Token,
			"content-type":         "application/json",
			"cookie":               cookie,
			"origin":               "https://discord.com",
			"referer":              "https://discord.com/channels/@me",
			"sec-fetch-dest":       "empty",
			"sec-fetch-mode":       "cors",
			"sec-fetch-site":       "same-origin",
			"user-agent":           in.UserAgent,
			"x-context-properties": xcontext,
			"x-debug-options":      "bugReporterEnabled",
			"x-discord-locale":     "en-US",
			"x-super-properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	}

	return req
}

func (in *Instance) xContextPropertiesHeaders(req *http.Request, cookie string) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":               "discord.com",
			"Cookie":             cookie,
			"X-Debug-Options":    "bugReporterEnabled",
			"Accept":             "*/*",
			"X-Discord-Locale":   "en-US",
			"X-Super-Properties": in.XSuper,
			"Authorization":      in.Token,
			"User-Agent":         in.UserAgent,
			"Accept-Language":    "en-US",
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"accept":             "*/*",
			"accept-language":    "en-US,en;q=0.9",
			"authorization":      in.Token,
			"cookie":             cookie,
			"referer":            "https://discord.com/channels/@me",
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-origin",
			"user-agent":         in.UserAgent,
			"x-debug-options":    "bugReporterEnabled",
			"x-discord-locale":   "en-US",
			"x-super-properties": in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	}

	return req
}

func (in *Instance) OpenChannelHeaders(req *http.Request, cookie string) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":                 "discord.com",
			"Cookie":               cookie,
			"Content-Type":         "application/json",
			"X-Debug-Options":      "bugReporterEnabled",
			"Accept":               "*/*",
			"Authorization":        in.Token,
			"X-Discord-Locale":     "en-US",
			"Accept-Language":      "en-US",
			"X-Context-Properties": "e30=",
			"User-Agent":           in.UserAgent,
			"X-Super-Properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"accept": "*/*",

			"accept-language":      "en-US,en;q=0.9",
			"authorization":        in.Token,
			"content-type":         "application/json",
			"cookie":               cookie,
			"origin":               "https://discord.com",
			"referer":              "https://discord.com/channels/@me",
			"sec-fetch-dest":       "empty",
			"sec-fetch-mode":       "cors",
			"sec-fetch-site":       "same-origin",
			"user-agent":           in.UserAgent,
			"x-context-properties": "e30=",
			"x-debug-options":      "bugReporterEnabled",
			"x-discord-locale":     "en-US",
			"x-super-properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	}
	return req
}

func (in *Instance) SendMessageHeaders(req *http.Request, cookie, recipient string) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":               "discord.com",
			"Cookie":             cookie,
			"Content-Type":       "application/json",
			"X-Debug-Options":    "bugReporterEnabled",
			"Accept":             "*/*",
			"Authorization":      in.Token,
			"X-Discord-Locale":   "en-US",
			"Accept-Language":    "en-US",
			"User-Agent":         in.UserAgent,
			"X-Super-Properties": in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"accept": "*/*",

			"accept-language":      "en-US,en;q=0.9",
			"authorization":        in.Token,
			"content-type":         "application/json",
			"cookie":               cookie,
			"origin":               "https://discord.com",
			"referer":              fmt.Sprintf("https://discord.com/channels/@me/%s", recipient),
			"sec-fetch-dest":       "empty",
			"sec-fetch-mode":       "cors",
			"sec-fetch-site":       "same-origin",
			"user-agent":           in.UserAgent,
			"x-context-properties": "e30=",
			"x-debug-options":      "bugReporterEnabled",
			"x-discord-locale":     "en-US",
			"x-super-properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	}
	return req
}

func (in *Instance) TypingHeaders(req *http.Request, cookie, snowflake string) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":               "discord.com",
			"User-Agent":         in.UserAgent,
			"Accept":             "*/*",
			"Accept-Language":    "en-US,en;q=0.5",
			"Authorization":      in.Token,
			"X-Super-Properties": in.XSuper,
			"X-Discord-Locale":   "en-US",
			"Cookie":             cookie,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"accept":               "*/*",
			"accept-language":      "en-US,en;q=0.9",
			"authorization":        in.Token,
			"content-type":         "application/json",
			"cookie":               cookie,
			"origin":               "https://discord.com",
			"referer":              fmt.Sprintf("https://discord.com/channels/@me/%s", snowflake),
			"sec-fetch-dest":       "empty",
			"sec-fetch-mode":       "cors",
			"sec-fetch-site":       "same-origin",
			"user-agent":           in.UserAgent,
			"x-context-properties": "e30=",
			"x-debug-options":      "bugReporterEnabled",
			"x-discord-locale":     "en-US",
			"x-super-properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	}

	return req
}

func (in *Instance) AtMeHeaders(req *http.Request, cookie string) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":               "discord.com",
			"Cookie":             cookie,
			"Content-Type":       "application/json",
			"X-Debug-Options":    "bugReporterEnabled",
			"Accept":             "*/*",
			"Authorization":      in.Token,
			"X-Discord-Locale":   "en-US",
			"Accept-Language":    "en-US",
			"User-Agent":         in.UserAgent,
			"X-Super-Properties": in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"accept":               "*/*",
			"accept-language":      "en-US,en;q=0.9",
			"authorization":        in.Token,
			"content-type":         "application/json",
			"cookie":               cookie,
			"origin":               "https://discord.com",
			"sec-fetch-dest":       "empty",
			"sec-fetch-mode":       "cors",
			"sec-fetch-site":       "same-origin",
			"user-agent":           in.UserAgent,
			"x-context-properties": "e30=",
			"x-debug-options":      "bugReporterEnabled",
			"x-discord-locale":     "en-US",
			"x-super-properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	}
	return req
}

func CommonHeaders(req *http.Request) *http.Request {

	for k, v := range map[string]string{
		"Host":               "discord.com",
		"Content-Type":       "application/json",
		"X-Debug-Options":    "bugReporterEnabled",
		"Accept":             "*/*",
		"X-Discord-Locale":   "en-US",
		"Accept-Language":    "en-US",
		"User-Agent":         "Discord/32114 CFNetwork/1331.0.7 Darwin/21.4.0",
		"X-Super-Properties": "eyJvcyI6ImlPUyIsImJyb3dzZXIiOiJEaXNjb3JkIGlPUyIsImRldmljZSI6ImlQYWQxMywxNiIsInN5c3RlbV9sb2NhbGUiOiJlbi1JTiIsImNsaWVudF92ZXJzaW9uIjoiMTI0LjAiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJkZXZpY2VfYWR2ZXJ0aXNlcl9pZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImRldmljZV92ZW5kb3JfaWQiOiJBMTgzNkNFRC1BRDI5LTRGRTAtQjVDNC0zODQ0NDU0MEFFQTciLCJicm93c2VyX3VzZXJfYWdlbnQiOiIiLCJicm93c2VyX3ZlcnNpb24iOiIiLCJvc192ZXJzaW9uIjoiMTUuNC4xIiwiY2xpZW50X2J1aWxkX251bWJlciI6MzIyNDcsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) UserInfoHeaders(req *http.Request, cookie string) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":               "discord.com",
			"Cookie":             cookie,
			"Content-Type":       "application/json",
			"X-Debug-Options":    "bugReporterEnabled",
			"Accept":             "*/*",
			"Authorization":      in.Token,
			"X-Discord-Locale":   "en-US",
			"Accept-Language":    "en-US",
			"User-Agent":         in.UserAgent,
			"X-Super-Properties": in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"accept":               "*/*",
			"accept-language":      "en-US,en;q=0.9",
			"authorization":        in.Token,
			"content-type":         "application/json",
			"cookie":               cookie,
			"origin":               "https://discord.com",
			"sec-fetch-dest":       "empty",
			"sec-fetch-mode":       "cors",
			"sec-fetch-site":       "same-origin",
			"user-agent":           in.UserAgent,
			"x-context-properties": "e30=",
			"x-debug-options":      "bugReporterEnabled",
			"x-discord-locale":     "en-US",
			"x-super-properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	}
	return req
}
