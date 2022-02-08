package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Function to use a captcha solving service and return a solved captcha key
func (in *Instance) SolveCaptcha(sitekey string) (string, error) {
	jsonx := Pload{
		ClientKey: in.Config.ClientKey,
		Task: Task{
			Type:       "HCaptchaTaskProxyless",
			WebsiteURL: "https://discord.com/channels/@me",
			WebsiteKey: sitekey,
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9003 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36",
		},
	}

	bytes, err := json.Marshal(jsonx)
	if err != nil {
		return "", fmt.Errorf("error marshalling json [%v]", err)

	}
	// Almost all solving services have similar API, so we can use the same function and replace the domain.
	resp, err := http.Post("https://api."+in.Config.CaptchaAPI+"/createTask", "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		return "", fmt.Errorf("error creating the request for captcha [%v]", err)
	}
	defer resp.Body.Close()
	// Get taskID from response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response body [%v]", err)
	}
	var response Resp
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling the response body [%v]", err)
	}
	switch response.ErrorID {
	case 0:
		// Poling server for the solved captcha
		jsonx = Pload{
			ClientKey: in.Config.ClientKey,
			TaskID:    response.TaskID,
		}
		y, err := json.Marshal(jsonx)
		if err != nil {
			return "", fmt.Errorf("error marshalling json [%v]", err)
		}
		// Anti Captcha documentation prescribes to use a delay of 5 seconds before requesting the captcha and 3 seconds delays after that.
		time.Sleep(5 * time.Second)
		p := 0
		for {
			if p > 100 {
				// Max retries
				break
			}
			resp, err := http.Post("https://api."+in.Config.CaptchaAPI+"/getTaskResult", "application/json", strings.NewReader(string(y)))
			if err != nil {
				return "", fmt.Errorf("error creating the request for captcha [%v]", err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("error reading the response body [%v]", err)
			}
			var response Resp
			err = json.Unmarshal(body, &response)
			if err != nil {
				return "", fmt.Errorf("error unmarshalling the response body [%v]", err)
			}
			if response.ErrorID == 16 {
				// The captcha being requested no longer exists in the active captchas
				return "", fmt.Errorf("error getting captcha [%v]", response.ErrorID)
			}
			if response.Status == "ready" {
				return response.Solution.Ans, nil

			} else if response.Status == "processing" {
				p++ // Incrementing the counter
				time.Sleep(3 * time.Second)
			}

		}
		return "", fmt.Errorf("max captcha retries reached [%v]", err)
	case 2:
		color.Red("No Available Captcha Workers - Increase your Maximum Bid in your Captcha API settings. Also try reducing the number of threads. Sleeping 10 seconds")
		time.Sleep(10 * time.Second)
		return "", fmt.Errorf("no captcha workers were available, retrying")
	case 3:
		return "", fmt.Errorf("the size of the captcha you are uploading is less than 100 bytes.")
	case 4:
		return "", fmt.Errorf("the size of the captcha you are uploading is greater than 500,000 bytes.")

	case 10:
		return "", fmt.Errorf("you have zero or negative captcha API balance")
	case 11:
		return "", fmt.Errorf("captcha was unsolvable.")
	default:
		return "", fmt.Errorf("unknown error [%v]", response.ErrorID)
	}
}

type Pload struct {
	ClientKey string `json:"clientKey"`
	Task      Task   `json:"task"`
	ErrorID   int    `json:"ErrorId"`
	TaskID    int    `json:"taskId"`
}

type Task struct {
	Type       string `json:"type"`
	WebsiteURL string `json:"websiteURL"`
	WebsiteKey string `json:"websiteKey"`
	UserAgent  string `json:"userAgent"`
}

type Resp struct {
	TaskID   int    `json:"taskID"`
	ErrorID  int    `json:"ErrorId"`
	Status   string `json:"status"`
	Solution Sol    `json:"solution"`
}

type Sol struct {
	Ans string `json:"gRecaptchaResponse"`
}
