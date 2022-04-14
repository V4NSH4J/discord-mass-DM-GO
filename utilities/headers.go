package utilities

import (
	"fmt"
	"net/http"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.61 Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36"
const XTrack = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2Ojk3LjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTcuMCIsImJyb3dzZXJfdmVyc2lvbiI6Ijk3LjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6OTk5OSwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbH0="
const XSuper = "eyJvcyI6Ik1hYyBPUyBYIiwiYnJvd3NlciI6IkRpc2NvcmQgQ2xpZW50IiwicmVsZWFzZV9jaGFubmVsIjoicHRiIiwiY2xpZW50X3ZlcnNpb24iOiIwLjAuNjEiLCJvc192ZXJzaW9uIjoiMjEuMy4wIiwib3NfYXJjaCI6ImFybTY0Iiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTIzODMzLCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ=="

func (in *Instance) cookieHeaders(req *http.Request) *http.Request {
	for k, v := range map[string]string{
		"Host":             "discord.com",
		"User-Agent":       UserAgent,
		"Accept":           "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Sec-Fetch-Site":   "none",
		"Sec-Fetch-Mode":   "navigate",
		"Sec-Fetch-User":   "?1",
		"Sec-Fetch-Dest":   "document",
		"sec-ch-ua-mobile": "?0",
		"Accept-Language":  "en-US,en;q=0.9",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) fingerprintHeaders(req *http.Request, cookie string) *http.Request {

	for k, v := range map[string]string{

		"Host":            "discord.com",
		"User-Agent":      UserAgent,
		"Accept":          "*/*",
		"Accept-Language": "en-US,en;q=0.5",
		"X-Track":         XTrack,
		"DNT":             "1",
		"Referer":         "https://discord.com/",
		"Cookie":          cookie,
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
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
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
	} {
		req.Header.Set(k, v)
	}
	return req
}

func (in *Instance) inviteHeaders(req *http.Request, cookie, fingerprint, xcontext string) *http.Request {
	for k, v := range map[string]string{
		"Host":                 "discord.com",
		"User-Agent":           UserAgent,
		"Accept":               "*/*",
		"Accept-Language":      "en-US,en;q=0.5",
		"Content-Type":         "application/json",
		"X-Context-Properties": xcontext,
		"Authorization":        in.Token,
		"X-Super-Properties":   XSuper,
		"X-Fingerprint":        fingerprint,
		"X-Discord-Locale":     "en-US",
		"X-Debug-Options":      "bugReporterEnabled",
		"Origin":               "https://discord.com",
		"Referer":              "https://discord.com/channels/@me",
		"Cookie":               cookie,
		"Sec-Fetch-Dest":       "empty",
		"Sec-Fetch-Mode":       "cors",
		"Sec-Fetch-Site":       "same-origin",
	} {
		req.Header.Set(k, v)
	}

	return req
}

func (in *Instance) xContextPropertiesHeaders(req *http.Request, cookie, fingerprint string) *http.Request {
	for k, v := range map[string]string{
		"Host":               `discord.com`,
		"User-Agent":         UserAgent,
		"Accept":             `*/*`,
		"Accept-Language":    `en-US,en;q=0.5`,
		"Authorization":      in.Token,
		"X-Super-Properties": XSuper,
		"X-Fingerprint":      fingerprint,
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
		"Origin":               "https://discord.com",
		"Referer":              "https://discord.com/channels/@me",
		"Cookie":               cookie,
		"Sec-Fetch-Dest":       "empty",
		"Sec-Fetch-Mode":       "cors",
		"Sec-Fetch-Site":       "same-origin",
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
		"Host": "discord.com",
		"User-Agent": UserAgent,
		"Accept": "*/*",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip, deflate",
		"Authorization": in.Token,
		"X-Super-Properties": XSuper,
		"X-Discord-Locale": "en-US",
		"X-Debug-Options": "bugReporterEnabled",
		"Origin": "https://discord.com",
		"Referer": fmt.Sprintf(`https://discord.com/channels/@me/%s`, snowflake),
		"Cookie": cookie,
		"Sec-Fetch-Dest": "empty",
		"Sec-Fetch-Mode": "cors",
		"Sec-Fetch-Site": "same-origin",
	} {
		req.Header.Set(k, v)
	}
	return req
}
