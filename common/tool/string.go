package tool

import (
	"unicode/utf8"
)

func CutText(text string, limit int, endText string) string {
	tll := utf8.RuneCountInString(text)
	ell := utf8.RuneCountInString(endText)
	if limit < 0 || limit >= tll || tll < ell || ell >= limit {
		return text
	}
	if tll > limit {
		if ell > 0 {
			return string([]rune(text)[:limit-ell]) + endText
		}
		return string([]rune(text)[:limit])
	}
	return text
}
