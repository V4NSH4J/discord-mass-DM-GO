// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
)

func LogInfo(format string, a ...any) {
	color.Printf("<fg=white>[</><fg=cyan;op=bold>INFO   </><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> » %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
}

func LogErr(format string, a ...any) {
	color.Printf("<fg=white>[</><fg=red;op=bold>ERROR  </><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> » %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
}

func LogSuccess(format string, a ...any) {
	color.Printf("<fg=white>[</><fg=green;op=bold>SUCCESS</><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> » %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
}

func LogFailed(format string, a ...any) {
	color.Printf("<fg=white>[</><fg=red;op=bold>FAILED </><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> » %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
}

func LogWarn(format string, a ...any) {
	color.Printf("<fg=white>[</><fg=red;op=bold>WARNING</><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> » %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
}

func LogLocked(format string, a ...any) {
	color.Printf("<fg=white>[</><fg=red;op=bold>LOCKED </><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> » %s\n", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
}

func CaptchaDetected(token string, key string) {
	color.Printf("<fg=white>[</><fg=red;op=bold>CAPTCHA</><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> » Token %s Captcha Detected [%s]\n", time.Now().Format("15:04:05"), token, key)
}

func CaptchaSolved(token string, key string) {
	var k string
	if len(key) < 25 {
		k = key
	} else {
		k = key[:25]
		k += "..."
	}
	color.Printf("<fg=white>[</><fg=green;op=bold>CAPTCHA</><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> » Token %s Captcha Solved [%s]\n", time.Now().Format("15:04:05"), token, k)
}

func UserInput(format string, a ...any) string {
	reader := bufio.NewReader(os.Stdin)
	var out string
	color.Printf("<fg=white>[</><fg=cyan;op=bold>INPUT  </><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> %s » ", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
	out, err := reader.ReadString('\n')
	if err != nil {
		color.Printf("<fg=white>[</><fg=red;op=bold>FATAL  </><fg=white>]</> » Error %s\n", err)
		ExitSafely()
	}
	out = strings.TrimSuffix(out, "\r\n")
	out = strings.TrimSuffix(out, "\n")
	return out
}

func UserInputInteger(format string, a ...any) int {
	reader := bufio.NewReader(os.Stdin)
	var out string
	color.Printf("<fg=white>[</><fg=cyan;op=bold>INPUT  </><fg=white>]</><fg=white>[</><fg=white;op=bold>%s</><fg=white>]</> %s » ", time.Now().Format("15:04:05"), fmt.Sprintf(format, a...))
	out, err := reader.ReadString('\n')
	if err != nil {
		color.Printf("<fg=white>[</><fg=red;op=bold>FATAL  </><fg=white>]</> » Error %s\n", err)
		ExitSafely()
	}
	if out == "" || out == "\n" {
		return 0 // Return 0 if user didn't enter anything
	}
	out = strings.TrimSuffix(out, "\r\n")
	out = strings.TrimSuffix(out, "\n")
	i, err := strconv.Atoi(out)
	if err != nil {
		color.Printf("<fg=white>[</><fg=red;op=bold>FATAL  </><fg=white>]</> » Error %s\n", err)
		ExitSafely()
	}
	return i
}

func ExitSafely() {
	color.Printf("<fg=white>[</><fg=red;op=bold>FATAL  </><fg=white>]</> » Press ENTER to exit\n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	os.Exit(0)
}

func PrintMenu(a []string) {
	for i := 1; i < len(a)+1; i++ {
		if i < 10 {
			color.Printf("<fg=white>[</><fg=cyan;op=bold>0%d</><fg=white>]</> » %s\n", i, a[i-1])
		} else {
			color.Printf("<fg=white>[</><fg=cyan;op=bold>%d</><fg=white>]</> » %s\n", i, a[i-1])
		}

	}
}

func PrintMenu2(a []string) {
	for i := 0; i < len(a); i++ {
		if i < 10 {
			color.Printf("<fg=white>[</><fg=cyan;op=bold>00%d</><fg=white>]</> » %s\n", i, a[i])
		} else if i < 100 && i >= 10 {
			color.Printf("<fg=white>[</><fg=cyan;op=bold>0%d</><fg=white>]</> » %s\n", i, a[i])
		} else {
			color.Printf("<fg=white>[</><fg=cyan;op=bold>%d</><fg=white>]</> » %s\n", i, a[i])
		}

	}
}
