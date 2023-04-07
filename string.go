package zutils

import "sort"

func Private(str string, p bool) string {
	if !p {
		return str
	}
	l := len([]rune(str))
	if l < 5 {
		return "*****"
	}
	if l < 10 {
		return "**********"
	}
	return "*****…*****"
}

func UnPrivate(real, view string) string {
	if view == "*****" || view == "**********" || view == "*****…*****" {
		return real
	}
	return view
}

func IsHas(arr []string, str string) bool {
	sort.Strings(arr)
	index := sort.SearchStrings(arr, str)
	return index < len(arr) && arr[index] == str
}
