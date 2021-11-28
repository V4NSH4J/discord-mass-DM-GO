// Copyright (C) 2021 github.com/V4NSH4J
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package utilities

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Inputs 2 slices of strings and returns a slice of strings which does not contain elements from the second slice
func RemoveSubset(s []string, r []string) []string {
	var n []string
	for _, v := range s {
		if !Contains(r, v) {
			n = append(n, v)
		}
	}
	return n
}
