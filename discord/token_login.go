// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func LaunchTokenLogin() {
	var token string
	utilities.LogWarn("You NEED Google Chrome installed to use this functionalit")
	token = utilities.UserInput("Enter a token which you want to login into: ")
	// Navigate to discord.com/login
	// We have to place this token into local storage
	// Refresh the page
	byt, err := os.ReadFile("tokenLogin.js")
	if err != nil {
		utilities.LogErr("Error while opening tokenLogin.js %v", err)
		utilities.ExitSafely()
	}
	x := strings.ReplaceAll(string(byt), "asdfgh", token)
	opts := append(
		chromedp.DefaultExecAllocatorOptions[3:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	)
	parentCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(parentCtx)
	defer cancel()
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://discord.com/login"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(x).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		})); err != nil {
		utilities.LogErr("Error while running chromedp. Error while evaluating %v", err)
		utilities.ExitSafely()
	}
	utilities.LogInfo("Press ENTER to close browser and continue program")
	fmt.Scanln()

}
