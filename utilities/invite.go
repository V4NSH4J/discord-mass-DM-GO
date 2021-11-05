package utilities

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Discord tracks accounts on their website using a fingerprint, adding this is essential otherwise accounts would get phone locked
func GetFingerprint() string {
	log.SetOutput(ioutil.Discard)
	resp, err := http.Get("https://discordapp.com/api/v9/experiments")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	type Fingerprintx struct {
		Fingerprint string `json:"fingerprint"`
	}
	var fingerprinty Fingerprintx
	json.Unmarshal(body, &fingerprinty)
	color.Yellow("Obtained Fingerprint: " + fingerprinty.Fingerprint + "\n")
	return fingerprinty.Fingerprint

}

type cookie struct {
	Dcfduid  string
	Sdcfduid string
}

// Getting cookies for legit looking requests
func GetCookie() cookie {
	log.SetOutput(ioutil.Discard)
	time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	resp, err := http.Get("https://discord.com")
	if err != nil {
		fmt.Printf("[%v]Error while getting cookies %v", time.Now().Format("15:05:04"), err)
		CookieNil := cookie{}
		return CookieNil
	}
	defer resp.Body.Close()

	Cookie := cookie{}
	if resp.Cookies() != nil {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "__dcfduid" {
				Cookie.Dcfduid = cookie.Value
			}
			if cookie.Name == "__sdcfduid" {
				Cookie.Sdcfduid = cookie.Value
			}
		}
	}
	color.Yellow("Obtained Cookies: " + "__dcfduid= " + Cookie.Dcfduid + " " + "__sdcfduid= " + Cookie.Sdcfduid + "\n")
	return Cookie
}

func JoinGuild(inviteCode string, token string) {
	url := "https://discord.com/api/v9/invites/" + inviteCode
	fmt.Println(url)
	Cookie := GetCookie()
	if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
		fmt.Printf("[%v]Empty cookie", time.Now().Format("15:05:04"))
		return
	}

	Cookies := "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + "locale=us"
	fmt.Println(Cookies)
	//var headers struct {}
	var headers struct{}
	requestBytes, _ := json.Marshal(headers)

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
	if err != nil {
		color.Red("ERR: Error while creating request \n")
	}
	//req.Header.Set("", )
	req.Header.Set("cookie", Cookies)
	req.Header.Set("authorization", token)

	httpClient := http.Client{}
	resp, err := httpClient.Do(commonHeaders(req))
	if err != nil {
		color.Red("[%v]Error while sending request %v \n", time.Now().Format("15:05:04"), err)
	}
	if resp.StatusCode == 200 {
		color.Green("Succesfully joined guild")
	}
	if resp.StatusCode != 200 {
		fmt.Printf("ERR: Unexpected Status code %v while joining token %v \n", resp.StatusCode, token)
	}

}

func commonHeaders(req *http.Request) *http.Request {
	req.Header.Set("accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-GB")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("sec-ch-ua", "'Chromium';v='92', ' Not A;Brand';v='99', 'Google Chrome';v='92'")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-context-properties", "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6Ijg4NTkwNzE3MjMwNTgwOTUxOSIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI4ODU5MDcxNzIzMDU4MDk1MjUiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9")
	req.Header.Set("x-fingerprint", GetFingerprint())
	req.Header.Set("x-super-properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2OjkzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTMuMCIsImJyb3dzZXJfdmVyc2lvbiI6IjkzLjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTAwODA0LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("origin", "https://discord.com")
	req.Header.Set("referer", "https://discord.com/channels/@me")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.16 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36")
	req.Header.Set("te", "trailers")
	return req
}

func readLines(filename string) ([]string, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	ex = filepath.ToSlash(ex)
	fmt.Println(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/" + filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func LaunchInviteJoiner() {
	color.Blue("\u2593\u2588\u2588\u2588\u2588\u2588\u2584\u0020\u2588\u2588\u2593\u0020\u2588\u2588\u2588\u2588\u2588\u2588\u0020\u2584\u2588\u2588\u2588\u2588\u2584\u0020\u0020\u2592\u2588\u2588\u2588\u2588\u2588\u0020\u0020\u2588\u2588\u2580\u2588\u2588\u2588\u0020\u2593\u2588\u2588\u2588\u2588\u2588\u2584\u0020\u0020\u0020\u0020\u0020\u2584\u2584\u2584\u2588\u2588\u2580\u2580\u2592\u2588\u2588\u2588\u2588\u2588\u0020\u0020\u2588\u2588\u2593\u2588\u2588\u2588\u2584\u0020\u0020\u0020\u0020\u2588\u2593\u2588\u2588\u2588\u2588\u2588\u0020\u2588\u2588\u2580\u2588\u2588\u2588\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2584\u2588\u2588\u2588\u2588\u0020\u2592\u2588\u2588\u2588\u2588\u2588\u0020\u0020\u000d\u000a\u2592\u2588\u2588\u2580\u0020\u2588\u2588\u2593\u2588\u2588\u2592\u2588\u2588\u0020\u0020\u0020\u0020\u2592\u2592\u2588\u2588\u2580\u0020\u2580\u2588\u0020\u2592\u2588\u2588\u2592\u0020\u0020\u2588\u2588\u2593\u2588\u2588\u0020\u2592\u0020\u2588\u2588\u2592\u2588\u2588\u2580\u0020\u2588\u2588\u258c\u0020\u0020\u0020\u0020\u0020\u0020\u2592\u2588\u2588\u0020\u2592\u2588\u2588\u2592\u0020\u0020\u2588\u2588\u2593\u2588\u2588\u2592\u2588\u2588\u0020\u2580\u2588\u0020\u0020\u0020\u2588\u2593\u2588\u0020\u0020\u0020\u2580\u2593\u2588\u2588\u0020\u2592\u0020\u2588\u2588\u2592\u0020\u0020\u0020\u0020\u2588\u2588\u2592\u0020\u2580\u2588\u2592\u2588\u2588\u2592\u0020\u0020\u2588\u2588\u2592\u000d\u000a\u2591\u2588\u2588\u0020\u0020\u0020\u2588\u2592\u2588\u2588\u2591\u0020\u2593\u2588\u2588\u2584\u0020\u0020\u2592\u2593\u2588\u0020\u0020\u0020\u0020\u2584\u2592\u2588\u2588\u2591\u0020\u0020\u2588\u2588\u2593\u2588\u2588\u0020\u2591\u2584\u2588\u0020\u2591\u2588\u2588\u0020\u0020\u0020\u2588\u258c\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u2588\u2588\u0020\u2592\u2588\u2588\u2591\u0020\u0020\u2588\u2588\u2592\u2588\u2588\u2593\u2588\u2588\u0020\u0020\u2580\u2588\u0020\u2588\u2588\u2592\u2588\u2588\u2588\u0020\u0020\u2593\u2588\u2588\u0020\u2591\u2584\u2588\u0020\u2592\u0020\u0020\u0020\u2592\u2588\u2588\u2591\u2584\u2584\u2584\u2592\u2588\u2588\u2591\u0020\u0020\u2588\u2588\u2592\u000d\u000a\u2591\u2593\u2588\u2584\u0020\u0020\u0020\u2591\u2588\u2588\u2591\u0020\u2592\u0020\u0020\u0020\u2588\u2588\u2592\u2593\u2593\u2584\u0020\u2584\u2588\u2588\u2592\u2588\u2588\u0020\u0020\u0020\u2588\u2588\u2592\u2588\u2588\u2580\u2580\u2588\u2584\u0020\u2591\u2593\u2588\u2584\u0020\u0020\u0020\u258c\u0020\u0020\u0020\u2593\u2588\u2588\u2584\u2588\u2588\u2593\u2592\u2588\u2588\u0020\u0020\u0020\u2588\u2588\u2591\u2588\u2588\u2593\u2588\u2588\u2592\u0020\u0020\u2590\u258c\u2588\u2588\u2592\u2593\u2588\u0020\u0020\u2584\u2592\u2588\u2588\u2580\u2580\u2588\u2584\u0020\u0020\u0020\u0020\u0020\u2591\u2593\u2588\u0020\u0020\u2588\u2588\u2592\u2588\u2588\u0020\u0020\u0020\u2588\u2588\u2591\u000d\u000a\u2591\u2592\u2588\u2588\u2588\u2588\u2593\u2591\u2588\u2588\u2592\u2588\u2588\u2588\u2588\u2588\u2588\u2592\u2592\u0020\u2593\u2588\u2588\u2588\u2580\u0020\u2591\u0020\u2588\u2588\u2588\u2588\u2593\u2592\u2591\u2588\u2588\u2593\u0020\u2592\u2588\u2588\u2591\u2592\u2588\u2588\u2588\u2588\u2593\u0020\u0020\u0020\u0020\u0020\u2593\u2588\u2588\u2588\u2592\u0020\u2591\u0020\u2588\u2588\u2588\u2588\u2593\u2592\u2591\u2588\u2588\u2592\u2588\u2588\u2591\u0020\u0020\u0020\u2593\u2588\u2588\u2591\u2592\u2588\u2588\u2588\u2588\u2591\u2588\u2588\u2593\u0020\u2592\u2588\u2588\u2592\u0020\u0020\u0020\u2591\u2592\u2593\u2588\u2588\u2588\u2580\u2591\u0020\u2588\u2588\u2588\u2588\u2593\u2592\u2591\u000d\u000a\u0020\u2592\u2592\u2593\u0020\u0020\u2592\u2591\u2593\u0020\u2592\u0020\u2592\u2593\u2592\u0020\u2592\u0020\u2591\u0020\u2591\u2592\u0020\u2592\u0020\u0020\u2591\u0020\u2592\u2591\u2592\u2591\u2592\u2591\u2591\u0020\u2592\u2593\u0020\u2591\u2592\u2593\u2591\u2592\u2592\u2593\u0020\u0020\u2592\u0020\u0020\u0020\u0020\u0020\u2592\u2593\u2592\u2592\u2591\u0020\u2591\u0020\u2592\u2591\u2592\u2591\u2592\u2591\u2591\u2593\u0020\u2591\u0020\u2592\u2591\u0020\u0020\u0020\u2592\u0020\u2592\u2591\u2591\u0020\u2592\u2591\u0020\u2591\u0020\u2592\u2593\u0020\u2591\u2592\u2593\u2591\u0020\u0020\u0020\u0020\u2591\u2592\u0020\u0020\u0020\u2592\u2591\u0020\u2592\u2591\u2592\u2591\u2592\u2591\u0020\u000d\u000a\u0020\u2591\u0020\u2592\u0020\u0020\u2592\u0020\u2592\u0020\u2591\u0020\u2591\u2592\u0020\u0020\u2591\u0020\u2591\u0020\u2591\u0020\u0020\u2592\u0020\u0020\u0020\u0020\u2591\u0020\u2592\u0020\u2592\u2591\u0020\u0020\u2591\u2592\u0020\u2591\u0020\u2592\u2591\u2591\u0020\u2592\u0020\u0020\u2592\u0020\u0020\u0020\u0020\u0020\u2592\u0020\u2591\u2592\u2591\u0020\u0020\u0020\u2591\u0020\u2592\u0020\u2592\u2591\u0020\u2592\u0020\u2591\u0020\u2591\u2591\u0020\u0020\u0020\u2591\u0020\u2592\u2591\u2591\u0020\u2591\u0020\u0020\u2591\u0020\u2591\u2592\u0020\u2591\u0020\u2592\u2591\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u2591\u0020\u0020\u2591\u0020\u2592\u0020\u2592\u2591\u0020\u000d\u000a\u0020\u2591\u0020\u2591\u0020\u0020\u2591\u0020\u2592\u0020\u2591\u0020\u0020\u2591\u0020\u0020\u2591\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u2591\u0020\u2592\u0020\u0020\u0020\u2591\u2591\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u2591\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u2591\u0020\u2591\u0020\u2591\u0020\u2591\u0020\u2592\u0020\u0020\u2592\u0020\u2591\u0020\u0020\u2591\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u2591\u2591\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u0020\u0020\u2591\u2591\u0020\u2591\u0020\u2591\u0020\u2592\u0020\u0020\u000d\u000a\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u2591\u0020\u0020\u2591\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u2591\u0020\u2591\u0020\u0020\u000d\u000a\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u2591\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\n Made by https://www.github.com/V4NSH4J\n")
	var code string
	color.Green("Enter Server Invite code (Not the invite link, just the code): ")
	fmt.Scanln(&code)
	lines, err := readLines("tokens.txt")
	red := color.New(color.FgRed).SprintFunc()
	if err != nil {
		fmt.Printf("%s Error while reading tokens.txt: %v", red("ERR"), err)
		return
	}
	start := time.Now()
	color.Red("Starting joining guilds with tokens!")
	var wg sync.WaitGroup
	wg.Add(len(lines))
	for i := 0; i < len(lines); i++ {
		time.Sleep(5 * time.Millisecond)
		go func(i int) {
			defer wg.Done()
			JoinGuild(code, lines[i])
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	color.Blue("Consider Starring this Repo on github for further updates! Happy Malicious Activity!")
	fmt.Printf("Joining took only %s", elapsed)

}
