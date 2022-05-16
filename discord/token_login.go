package discord

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
)

func LaunchTokenLogin() {
	var token string
	color.Red("You NEED Google Chrome installed to run this function.")
	color.White("[%v] Enter your token: ", time.Now().Format("15:04:05"))
	fmt.Scanln(&token)
	// Navigate to discord.com/login
	// We have to place this token into local storage
	// Refresh the page
	byt, err := os.ReadFile("instance/tokenLogin.js")
	if err != nil {
		color.Red("Error while reading tokenLogin Javascript file %v", err)
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
		color.Red("Error while Evaluating %v", err)
		utilities.ExitSafely()
	}
	color.Green("Press ENTER to close window and continue program")
	fmt.Scanln()

}
