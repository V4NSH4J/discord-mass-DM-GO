// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

import (
	"strconv"
	"time"
)

func Snowflake() int64 {
	snowflake := strconv.FormatInt((time.Now().UTC().UnixNano()/1000000)-1420070400000, 2) + "0000000000000000000000"
	nonce, _ := strconv.ParseInt(snowflake, 2, 64)
	return nonce
}
