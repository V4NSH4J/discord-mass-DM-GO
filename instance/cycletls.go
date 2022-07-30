// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package instance

import "github.com/Danny-Dasilva/CycleTLS/cycletls"

func (in *Instance) CycleOptions(body string, headers map[string]string) cycletls.Options {
	return cycletls.Options{
		Timeout:   in.Config.ProxySettings.Timeout,
		Proxy:     in.ProxyProt,
		Headers:   headers,
		Body:      body,
		Ja3:       in.JA3,
		UserAgent: in.UserAgent,
	}
}
