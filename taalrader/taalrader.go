package taalrader

import (
	"github.com/pebbe/textcat"

	"fmt"
	"strings"
)

func raadtaal(tc *textcat.TextCat, text string) string {
	l, e := tc.Classify(text)
	if e != nil {
		return fmt.Sprintln(e)
	}
	return fmt.Sprintln(strings.Join(l, "\n"))
}

func Raadtaal(s string) string {
	tc := textcat.NewTextCat()
	tc.EnableAllUtf8Languages()
	return raadtaal(tc, s)
}

func Raadtaalboth(s string) string {
	tc := textcat.NewTextCat()
	tc.EnableAllUtf8Languages()
	tc.EnableAllRawLanguages()
	return raadtaal(tc, s)
}

func Raadtaalraw(s string) string {
	tc := textcat.NewTextCat()
	tc.EnableAllRawLanguages()
	return raadtaal(tc, s)
}
