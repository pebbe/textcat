/*
A Go package for n-gram based text categorization, with support for utf-8 and raw text.
*/
package textcat

import (
	"sort"
)

type TextCat struct {
	utf8 bool
	lang map[string]bool
}

func NewTextCat(utf8 bool) *TextCat {
	tc := &TextCat{utf8: utf8, lang: make(map[string]bool)}
	data := dataRaw
	if utf8 {
		data = dataUtf8
	}
	for d := range data {
		tc.lang[d] = true
	}
	return tc
}

func (tc *TextCat) ActiveLanguages() []string {
	a := make([]string, 0, len(tc.lang))
	for l := range tc.lang {
		if tc.lang[l] {
			a = append(a, l)
		}
	}
	sort.Strings(a)
	return a
}

func (tc *TextCat) AvailableLanguages() []string {
	a := make([]string, 0, len(tc.lang))
	for l := range tc.lang {
		a = append(a, l)
	}
	sort.Strings(a)
	return a
}

func (tc *TextCat) DisableLanguages(language ...string) {
	for _, lang := range language {
		if _, exists := tc.lang[lang]; exists {
			tc.lang[lang] = false
		}
	}
}

func (tc *TextCat) EnableLanguages(language ...string) {
	for _, lang := range language {
		if _, exists := tc.lang[lang]; exists {
			tc.lang[lang] = true
		}
	}
}
