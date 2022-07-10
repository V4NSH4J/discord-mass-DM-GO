// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import (
	"fmt"
	"net/http"
)

func (in *Instance) cookieHeaders(req *http.Request) *http.Request {
	if in.Config.OtherSettings.Mode == 2 {
		for k, v := range map[string]string{
			"Host":            "discord.com",
			"User-Agent":      "Discord/125.0 (iPad; iOS 15.4.1; Scale/2.00)",
			"Accept-Language": "en-IN;q=1",
			"Accept":          "*/*",
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"Host": "discord.com",
			"User-Agent": in.UserAgent,
			"Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
			"Accept-Language": "en-US,en;q=0.5",
			"Accept-Encoding": "gzip, deflate",
			"sec-gpc": "1",
			"Upgrade-Insecure-Requests": "1",
			"Sec-Fetch-Dest": "document",
			"Sec-Fetch-Mode": "navigate",
			"Sec-Fetch-Site": "none",
			"Sec-Fetch-User": "?1",
			
		} {
			req.Header.Set(k, v)
		}
	}

	return req
}

func (in *Instance) cfBmHeaders(req *http.Request, cookie string) *http.Request {
	for k, v := range map[string]string{
		"Host": "discord.com",
		"Cookie": cookie,
		"User-Agent": in.UserAgent,
		"Accept": "*/*",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip, deflate",
		"Content-Type": "application/json",
		"Origin": "https://discord.com",
		"sec-gpc": "1",
		"Referer": "https://discord.com/",
		"Sec-Fetch-Dest": "empty",
		"Sec-Fetch-Mode": "cors",
		"Sec-Fetch-Site": "same-origin",
		
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
			"Accept-Language":      "en-IN",
			"X-Context-Properties": xcontext,
			"User-Agent":           in.UserAgent,
			"X-Super-Properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"Host": "discord.com",
			"Cookie": in.Cookie,
			"User-Agent": in.UserAgent,
			"Accept": "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			"Content-Type": "application/json",
			"X-Context-Properties": xcontext,
			"Authorization": in.Token,
			"X-Super-Properties": in.XSuper,
			"X-Discord-Locale": "en-US",
			"X-Debug-Options": "bugReporterEnabled",
			"Origin": "https://discord.com",
			"sec-gpc": "1",
			"Referer": "https://discord.com/channels/@me",
			"Sec-Fetch-Dest": "empty",
			"Sec-Fetch-Mode": "cors",
			"Sec-Fetch-Site": "same-origin",
			
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
			"Accept-Language":    "en-IN",
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"Cookie": cookie,
			"User-Agent": in.UserAgent,
			"Accept": "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			"Authorization": in.Token,
			"X-Super-Properties": in.XSuper,
			"X-Discord-Locale": "en-US",
			"X-Debug-Options": "bugReporterEnabled",
			"sec-gpc": "1",
			"Referer": "https://discord.com/channels/@me",
			"Sec-Fetch-Dest": "empty",
			"Sec-Fetch-Mode": "cors",
			"Sec-Fetch-Site": "same-origin",
			
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
			"Accept-Language":      "en-IN",
			"X-Context-Properties": "e30=",
			"User-Agent":           in.UserAgent,
			"X-Super-Properties":   in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"Host": "discord.com",
			"Cookie": cookie,
			"User-Agent": in.UserAgent,
			"Accept": "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			"Content-Type": "application/json",
			"X-Context-Properties": "e30=",
			"Authorization": in.Token,
			"X-Super-Properties": in.XSuper,
			"X-Discord-Locale": "en-US",
			"X-Debug-Options": "bugReporterEnabled",
			"Origin": "https://discord.com",
			"sec-gpc": "1",
			"Referer": "https://discord.com/channels/@me",
			"Sec-Fetch-Dest": "empty",
			"Sec-Fetch-Mode": "cors",
			"Sec-Fetch-Site": "same-origin",
			
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
			"Accept-Language":    "en-IN",
			"User-Agent":         in.UserAgent,
			"X-Super-Properties": in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"Host": "discord.com",
			"Cookie": cookie,
			"User-Agent": in.UserAgent,
			"Accept": "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			"Content-Type": "application/json",
			"Authorization": in.Token,
			"X-Super-Properties": in.XSuper,
			"X-Discord-Locale": "en-US",
			"X-Debug-Options": "bugReporterEnabled",
			"Origin": "https://discord.com",
			"sec-gpc": "1",
			"Referer": fmt.Sprintf("https://discord.com/channels/@me/%s", recipient),
			"Sec-Fetch-Dest": "empty",
			"Sec-Fetch-Mode": "cors",
			"Sec-Fetch-Site": "same-origin",
			
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
			"Host":               "discord.com",
			"User-Agent":         in.UserAgent,
			"Accept":             "*/*",
			"Accept-Language":    "en-US,en;q=0.5",
			"Authorization":      in.Token,
			"X-Super-Properties": in.XSuper,
			"X-Discord-Locale":   "en-US",
			"X-Debug-Options":    "bugReporterEnabled",
			"Origin":             "https://discord.com",
			"Referer":            fmt.Sprintf(`https://discord.com/channels/@me/%s`, snowflake),
			"Cookie":             cookie,
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "same-origin",
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
			"Accept-Language":    "en-IN",
			"User-Agent":         in.UserAgent,
			"X-Super-Properties": in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"Host": "discord.com",
			"Cookie": cookie,
			"User-Agent": in.UserAgent,
			"Accept": "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			"Authorization": in.Token,
			"X-Super-Properties": in.XSuper,
			"X-Discord-Locale": "en-US",
			"X-Debug-Options": "bugReporterEnabled",
			"sec-gpc": "1",
			"Referer": "https://discord.com/login",
			"Sec-Fetch-Dest": "empty",
			"Sec-Fetch-Mode": "cors",
			"Sec-Fetch-Site": "same-origin",
			
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
		"Accept-Language":    "en-IN",
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
			"Accept-Language":    "en-IN",
			"User-Agent":         in.UserAgent,
			"X-Super-Properties": in.XSuper,
		} {
			req.Header.Set(k, v)
		}
	} else {
		for k, v := range map[string]string{
			"Host": "discord.com",
			"Cookie": cookie,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:102.0) Gecko/20100101 Firefox/102.0",
			"Accept": "*/*",
			"Accept-Language": "en-US,en;q=0.5",
			"Accept-Encoding": "gzip, deflate",
			"Authorization": in.Token,
			"X-Super-Properties": in.XSuper,
			"X-Discord-Locale": "en-US",
			"X-Debug-Options": "bugReporterEnabled",
			"sec-gpc": "1",
			"Referer": "https://discord.com/channels/@me",
			"Sec-Fetch-Dest": "empty",
			"Sec-Fetch-Mode": "cors",
			"Sec-Fetch-Site": "same-origin",
			
		} {
			req.Header.Set(k, v)
		}
	}
	return req
}
