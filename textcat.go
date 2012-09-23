package textcat

import (
	"errors"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	thresholdValue = 1.03
	maxCandidates  = 5
	minDocSize     = 25
)

var (
	errShort   = errors.New("SHORT")
	errUnknown = errors.New("UNKNOWN")
	errAvail   = errors.New("NOPATTERNS")
)

type TextCat struct {
	utf8 bool
	raw  bool
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

func NewTextCat(enableUtf8, enableRaw bool) *TextCat {
	tc := &TextCat{utf8: enableUtf8, raw: enableRaw, lang: make(map[string]bool)}
	for d := range data {
		tc.lang[d] = true
	}
	return tc
}

func (tc *TextCat) ActiveLanguages() []string {
	a := make([]string, 0, len(tc.lang))
	for l := range tc.lang {
		if tc.lang[l] {
			if tc.utf8 && strings.HasSuffix(l, ".utf8") {
				a = append(a, l)
			} else if tc.raw && strings.HasSuffix(l, ".raw") {
				a = append(a, l)
			}
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

func (tc *TextCat) DisableRawLanguages() {
	tc.raw = false
}

func (tc *TextCat) DisableUtf8Languages() {
	tc.utf8 = false
}

func (tc *TextCat) EnableLanguages(language ...string) {
	for _, lang := range language {
		if _, exists := tc.lang[lang]; exists {
			tc.lang[lang] = true
		}
	}
}

func (tc *TextCat) EnableRawLanguages() {
	tc.raw = true
}

func (tc *TextCat) EnableUtf8Languages() {
	tc.utf8 = true
}

func (tc *TextCat) Classify(text string) (languages []string, err error) {
	languages = make([]string, 0, maxCandidates)

	if tc.raw && len(text) < minDocSize {
		err = errShort
		return
	}
	if tc.utf8 && utf8.RuneCountInString(text) < minDocSize {
		err = errShort
		return
	}

	scores := make([]*resultType, 0, len(tc.lang))
	pattypes := make([]bool, 0, 2)
	if tc.utf8 {
		pattypes = append(pattypes, true)
	}
	if tc.raw {
		pattypes = append(pattypes, false)
	}
	for _, utf8 := range pattypes {
		patt := getPatterns(text, utf8)
		suffix := ".raw"
		if utf8 {
			suffix = ".utf8"
		}
		for lang := range tc.lang {
			if !tc.lang[lang] || !strings.HasSuffix(lang, suffix) {
				continue
			}
			score := 0
			for n, p := range patt {
				i, ok := data[lang][p.s]
				if !ok {
					i = maxPatterns
				}
				if n > i {
					score += n - i
				} else {
					score += i - n
				}
			}
			scores = append(scores, &resultType{score, lang})
		}
	}
	if len(scores) == 0 {
		err = errAvail
		return
	}

	minScore := maxPatterns * maxPatterns
	for _, sco := range scores {
		if sco.score < minScore {
			minScore = sco.score
		}
	}
	threshold := float64(minScore) * thresholdValue
	nCandidates := 0
	for _, sco := range scores {
		if float64(sco.score) <= threshold {
			nCandidates += 1
		}
	}
	if nCandidates > maxCandidates {
		err = errUnknown
		return
	}

	lowScores := make([]*resultType, 0, nCandidates)
	for _, sco := range scores {
		if float64(sco.score) <= threshold {
			lowScores = append(lowScores, sco)
		}
	}
	sort.Sort(resultsType(lowScores))
	for _, sco := range lowScores {
		languages = append(languages, sco.lang)
	}

	return
}
