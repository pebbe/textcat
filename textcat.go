/*
A Go package for n-gram based text categorization, with support for utf-8 and raw text.
*/
package textcat

import (
	"errors"
	"sort"
	"unicode/utf8"
)

const (
	thresholdValue = 1.03
	maxCandidates  = 5
	minDocSize     = 25
)

type TextCat struct {
	utf8 bool
	lang map[string]bool
}

type resultType struct {
	score int
	lang  string
}

type resultsType []*resultType

func (r resultsType) Len() int {
	return len(r)
}

func (r resultsType) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r resultsType) Less(i, j int) bool {
	if r[i].score != r[j].score {
		return r[i].score < r[j].score
	}
	return r[i].lang < r[j].lang
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

func (tc *TextCat) Classify(text string) (languages []string, err error) {
	languages = make([]string, 0)

	l := len(text)
	if tc.utf8 {
		l = utf8.RuneCountInString(text)
	}
	if l < minDocSize {
		err = errors.New("SHORT")
		return
	}

	scores := make([]*resultType, 0, len(tc.lang))
	data := dataRaw
	if tc.utf8 {
		data = dataUtf8
	}
	for lang := range tc.lang {
		if !tc.lang[lang] {
			continue
		}
		patt := getPatterns(text, tc.utf8)
		score := 0
		for n, p := range patt {
			i, ok := data[lang][p.s]
			if !ok {
				i = maxPatterns + 1
			}
			score += abs(n + 1 - i)
		}
		scores = append(scores, &resultType{score, lang})
	}

	sort.Sort(resultsType(scores)) // TO DO: dit kan sneller, niet alles sorteren, maar alleen de scores onder de threshold

	threshold := float64(scores[0].score) * thresholdValue

	for _, sco := range scores {
		if float64(sco.score) > threshold {
			break
		}
		languages = append(languages, sco.lang)
	}
	if len(languages) > maxCandidates {
		languages = languages[0:0]
		err = errors.New("UNKNOWN")
	}
	return
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
