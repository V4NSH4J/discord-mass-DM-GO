package utilities

import "fmt"

func Contains(s []string, e string) bool {
	defer HandleOutOfBounds()
	if len(s) == 0 {
		return false
	}
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveSubset(s []string, r []string) []string {
	var n []string
	for _, v := range s {
		if !Contains(r, v) {
			n = append(n, v)
		}
	}
	return n
}

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func HandleOutOfBounds() {
	if r := recover(); r != nil {
		fmt.Printf("Recovered from Panic %v", r)
	}
}
