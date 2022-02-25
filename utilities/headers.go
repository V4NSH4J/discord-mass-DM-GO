package utilities

import (
	"net/http"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0"
const XTrack = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2Ojk3LjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTcuMCIsImJyb3dzZXJfdmVyc2lvbiI6Ijk3LjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6OTk5OSwiY2xpZW50X2V2ZW50X3NvdXJjZSI6bnVsbH0="
const XSuper = "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2OjgzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvODMuMCIsImJyb3dzZXJfdmVyc2lvbiI6IjgzLjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTE2MzExLCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ=="

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

func (in *Instance) inviteHeaders(req *http.Request, cookie string, fingerprint string) *http.Request {

	for k, v := range map[string]string{
		"Host":                 "discord.com",
		"User-Agent":           UserAgent,
		"Accept":               "*/*",
		"Accept-Language":      "en-US,en;q=0.5",
		"Authorization":        in.Token,
		"X-Super-Properties":   XSuper,
		"X-Fingerprint":        fingerprint,
		"X-Context-Properties": "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6IjkyNTA2NjE3MjM0MzM5ODQyMCIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI5MjUwNjYxNzIzNDMzOTg0MjMiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9",
		"X-Discord-Locale":     "en-US",
		"X-Debug-Options":      "bugReporterEnabled",
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
