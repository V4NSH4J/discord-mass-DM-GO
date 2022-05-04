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

const UserAgent = "Discord/32114 CFNetwork/1331.0.7 Darwin/21.4.0"
const XTrack = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2Ojk3LjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTcuMCIsImJyb3dzZXJfdmVyc2lvbiI6Ijk3LjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6OTk5OSwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbH0="
const XSuper = "eyJvcyI6ImlPUyIsImJyb3dzZXIiOiJEaXNjb3JkIGlPUyIsImRldmljZSI6ImlQYWQxMywxNiIsInN5c3RlbV9sb2NhbGUiOiJlbi1JTiIsImNsaWVudF92ZXJzaW9uIjoiMTI0LjAiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJkZXZpY2VfYWR2ZXJ0aXNlcl9pZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImRldmljZV92ZW5kb3JfaWQiOiJBMTgzNkNFRC1BRDI5LTRGRTAtQjVDNC0zODQ0NDU0MEFFQTciLCJicm93c2VyX3VzZXJfYWdlbnQiOiIiLCJicm93c2VyX3ZlcnNpb24iOiIiLCJvc192ZXJzaW9uIjoiMTUuNC4xIiwiY2xpZW50X2J1aWxkX251bWJlciI6MzIyNDcsImNsaWVudF9ldmVudF9zb3VyY2UiOm51bGx9"


func (in *Instance) cookieHeaders(req *http.Request) *http.Request {
	for k, v := range map[string]string{
		"Host":             "discord.com",
		"User-Agent":       UserAgent,
		"Accept":           "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		// "Sec-Fetch-Site":   "none",
		// "Sec-Fetch-Mode":   "navigate",
		// "Sec-Fetch-User":   "?1",
		// "Sec-Fetch-Dest":   "document",
		// "sec-ch-ua-mobile": "?0",
		"Accept-Language":  "en-US,en;q=0.9",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) cfBmHeaders(req *http.Request, cookie string) *http.Request {
	for k, v := range map[string]string{
		"Host":            "discord.com",
		"User-Agent":      UserAgent,
		"Accept":          "*/*",
		"Accept-Language": "en-US,en;q=0.5",
		"Content-Type":    "application/json",
		"Origin":          "https://discord.com",
		"DNT":             "1",
		"Referer":         "https://discord.com/",
		"Cookie":          cookie,
		// "Sec-Fetch-Dest":  "empty",
		// "Sec-Fetch-Mode":  "cors",
		// "Sec-Fetch-Site":  "same-origin",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) inviteHeaders(req *http.Request, cookie, xcontext string) *http.Request {
	for k, v := range map[string]string{
		"Host":                 "discord.com",
		"User-Agent":           UserAgent,
		"Accept":               "*/*",
		"Accept-Language":      "en-US,en;q=0.5",
		"Content-Type":         "application/json",
		"X-Context-Properties": xcontext,
		"Authorization":        in.Token,
		"X-Super-Properties":   XSuper,
		"X-Discord-Locale":     "en-US",
		"X-Debug-Options":      "bugReporterEnabled",
		// "Origin":               "https://discord.com",
		// "Referer":              "https://discord.com/channels/@me",
		"Cookie":               cookie,
		// "Sec-Fetch-Dest":       "empty",
		// "Sec-Fetch-Mode":       "cors",
		// "Sec-Fetch-Site":       "same-origin",
	} {
		req.Header.Set(k, v)
	}

	return req
}

func (in *Instance) xContextPropertiesHeaders(req *http.Request, cookie string) *http.Request {
	for k, v := range map[string]string{
		"Host":               `discord.com`,
		"User-Agent":         UserAgent,
		"Accept":             `*/*`,
		"Accept-Language":    `en-US,en;q=0.5`,
		"Authorization":      in.Token,
		"X-Super-Properties": XSuper,
		"X-Discord-Locale":   `en-US`,
		"X-Debug-Options":    `bugReporterEnabled`,
		"Referer":            `https://discord.com/channels/@me`,
		"Cookie":             cookie,
		"Sec-Fetch-Dest":     `empty`,
		"Sec-Fetch-Mode":     `cors`,
		"Sec-Fetch-Site":     `same-origin`,
	} {
		req.Header.Set(k, v)
	}

	return req
}

func (in *Instance) OpenChannelHeaders(req *http.Request, cookie string) *http.Request {
	for k, v := range map[string]string{
		"Host":                 "discord.com",
		"User-Agent":           UserAgent,
		"Accept":               "*/*",
		"Accept-Language":      "en-US,en;q=0.5",
		"Content-Type":         "application/json",
		"X-Context-Properties": "e30=",
		"Authorization":        in.Token,
		"X-Super-Properties":   XSuper,
		"X-Discord-Locale":     "en-US",
		"X-Debug-Options":      "bugReporterEnabled",
		// "Origin":               "https://discord.com",
		// "Referer":              "https://discord.com/channels/@me",
		"Cookie":               cookie,
		// "Sec-Fetch-Dest":       "empty",
		// "Sec-Fetch-Mode":       "cors",
		// "Sec-Fetch-Site":       "same-origin",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) SendMessageHeaders(req *http.Request, cookie, recipient string) *http.Request {
	for k, v := range map[string]string{
		"Host":               "discord.com",
		"User-Agent":         UserAgent,
		"Accept":             "*/*",
		"Accept-Language":    "en-US,en;q=0.5",
		"Content-Type":       "application/json",
		"Authorization":      in.Token,
		"X-Super-Properties": XSuper,
		"X-Discord-Locale":   "en-US",
		"X-Debug-Options":    "bugReporterEnabled",
		"Origin":             "https://discord.com",
		"Referer":            fmt.Sprintf(`https://discord.com/channels/@me/%s`, recipient),
		"Cookie":             cookie,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) TypingHeaders(req *http.Request, cookie, snowflake string) *http.Request {
	for k, v := range map[string]string{
		"Host":               "discord.com",
		"User-Agent":         UserAgent,
		"Accept":             "*/*",
		"Accept-Language":    "en-US,en;q=0.5",
		"Accept-Encoding":    "gzip, deflate",
		"Authorization":      in.Token,
		"X-Super-Properties": XSuper,
		"X-Discord-Locale":   "en-US",
		"X-Debug-Options":    "bugReporterEnabled",
		// "Origin":             "https://discord.com",
		// "Referer":            fmt.Sprintf(`https://discord.com/channels/@me/%s`, snowflake),
		"Cookie":             cookie,
		// "Sec-Fetch-Dest":     "empty",
		// "Sec-Fetch-Mode":     "cors",
		// "Sec-Fetch-Site":     "same-origin",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) AtMeHeaders(req *http.Request, cookie string) *http.Request {
	for k, v := range map[string]string{

		"User-Agent":         UserAgent,
		"Accept":             "*/*",
		"Accept-Language":    "en-US,en;q=0.5",
		"Authorization":      in.Token,
		"X-Super-Properties": XSuper,
		"X-Discord-Locale":   "en-US",
		"X-Debug-Options":    "bugReporterEnabled",
		// "Origin":             "https://discord.com",
		// "Referer":            `https://discord.com/channels/@me/`,
		"Content-Type":       "application/json",
		"Cookie":             cookie,
		// "Sec-Fetch-Dest":     "empty",
		// "Sec-Fetch-Mode":     "cors",
		// "Sec-Fetch-Site":     "same-origin",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func CommonHeaders(req *http.Request) *http.Request {

	req.Header.Set("X-Super-Properties", XSuper)
	// req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("x-debug-options", "bugReporterEnabled")
	// req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("X-Discord-Locale", "en-US")
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	// req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("accept-language", "en-US")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", UserAgent)
	req.Header.Set("TE", "trailers")
	return req
}

func (in *Instance) UserInfoHeaders(req *http.Request, cookie string) *http.Request {
	for k, v := range map[string]string{
		"Host": "discord.com",
		"Cookie": cookie,
		"User-Agent": UserAgent,
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip, deflate",
		"Authorization": in.Token,
		"X-Super-Properties": XSuper,
		"X-Discord-Locale": "en-US",
		"X-Debug-Options": "bugReporterEnabled",
		// "Dnt": "1",
		//"Referer": "https://discord.com/channels/942940056157560862/942940056157560865",
		// "Sec-Fetch-Dest": "empty",
		// "Sec-Fetch-Mode": "cors",
		// "Sec-Fetch-Site": "same-origin",
		// "Te": "trailers",
	} {
		req.Header.Set(k, v)
	}

	return req 
}