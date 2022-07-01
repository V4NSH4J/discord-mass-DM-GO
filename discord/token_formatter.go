// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package discord

import (
	"fmt"
	"os"
	"time"

	"github.com/V4NSH4J/discord-mass-dm-GO/instance"
	"github.com/V4NSH4J/discord-mass-dm-GO/utilities"
)

func LaunchTokenFormatter() {
	cfg, instances, err := instance.GetEverything()
	if err != nil {
		utilities.LogErr("Error while getting neccessary information %v", err)
	}
	var tokenFile, changedFile string
	if cfg.OtherSettings.Logs {
		path := fmt.Sprintf(`logs/token_formatter/DMDGO-TF-%s-%s`, time.Now().Format(`2006-01-02 15-04-05`), utilities.RandStringBytes(5))
		err := os.MkdirAll(path, 0755)
		if err != nil && !os.IsExist(err) {
			utilities.LogErr("Error creating logs directory: %s", err)
			utilities.ExitSafely()
		}
		tokenFileX, err := os.Create(fmt.Sprintf(`%s/token.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating token file: %s", err)
			utilities.ExitSafely()
		}
		tokenFileX.Close()
		ChangedFileX, err := os.Create(fmt.Sprintf(`%s/changed.txt`, path))
		if err != nil {
			utilities.LogErr("Error creating success file: %s", err)
			utilities.ExitSafely()
		}
		ChangedFileX.Close()
		tokenFile, changedFile = tokenFileX.Name(), ChangedFileX.Name()
		for i := 0; i < len(instances); i++ {
			instances[i].WriteInstanceToFile(tokenFile)
		}
	}
	var tokens []string

	for i := 0; i < len(instances); i++ {
		if cfg.OtherSettings.Logs {
			instances[i].Email = ""
			instances[i].Password = ""
			instances[i].WriteInstanceToFile(changedFile)
		}
		tokens = append(tokens, instances[i].Token)
	}
	_ = utilities.TruncateLines("tokens.txt", tokens)
	utilities.LogSuccess("Token formatter has finished")
}
