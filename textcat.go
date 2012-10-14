package textcat

import (
	"errors"
	"github.com/pebbe/util"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	defaultThresholdValue = 1.03
	defaultMaxCandidates  = 5
	defaultMinDocSize     = 25
)

var (
	errShort   = errors.New("SHORT")
	errUnknown = errors.New("UNKNOWN")
	errAvail   = errors.New("NOPATTERNS")
)

type TextCat struct {
	utf8           bool
	raw            bool
	lang           map[string]bool
	thresholdValue float64
	maxCandidates  int
	minDocSize     int
	extra          map[string]map[string]int
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

func NewTextCat() *TextCat {
	tc := &TextCat{
		lang:           make(map[string]bool),
		thresholdValue: defaultThresholdValue,
		maxCandidates:  defaultMaxCandidates,
		minDocSize:     defaultMinDocSize}
	for d := range data {
		tc.lang[d] = false
	}
	tc.extra = make(map[string]map[string]int)
	return tc
}

func (tc *TextCat) SetThresholdValue(thresholdValue float64) {
	tc.thresholdValue = thresholdValue
}

func (tc *TextCat) GetThresholdValue() float64 {
	return tc.thresholdValue
}

func (tc *TextCat) SetMaxCandidates(maxCandidates int) {
	tc.maxCandidates = maxCandidates
}

func (tc *TextCat) GetMaxCandidates() int {
	return tc.maxCandidates
}

func (tc *TextCat) SetMinDocSize(minDocSize int) {
	tc.minDocSize = minDocSize
}

func (tc *TextCat) GetMinDocSize() int {
	return tc.minDocSize
}

func (tc *TextCat) ActiveLanguages() []string {
	a := make([]string, 0, len(tc.lang))
	for lang := range tc.lang {
		if tc.lang[lang] {
			a = append(a, lang)
		}
	}
	sort.Strings(a)
	return a
}

func (tc *TextCat) AvailableLanguages() []string {
	a := make([]string, 0, len(tc.lang))
	for lang := range tc.lang {
		a = append(a, lang)
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
	tc.raw = false
	tc.utf8 = false
	for lang := range tc.lang {
		if tc.lang[lang] {
			if !tc.raw && strings.HasSuffix(lang, ".raw") {
				tc.raw = true
			} else if !tc.utf8 && strings.HasSuffix(lang, ".utf8") {
				tc.utf8 = true
			}
			if tc.raw && tc.utf8 {
				break
			}
		}
	}
}

func (tc *TextCat) DisableAllRawLanguages() {
	tc.raw = false
	for lang := range tc.lang {
		if strings.HasSuffix(lang, ".raw") {
			tc.lang[lang] = false
		}
	}
}

func (tc *TextCat) DisableAllUtf8Languages() {
	tc.utf8 = false
	for lang := range tc.lang {
		if strings.HasSuffix(lang, ".utf8") {
			tc.lang[lang] = false
		}
	}
}

func (tc *TextCat) EnableLanguages(language ...string) {
	for _, lang := range language {
		if _, exists := tc.lang[lang]; exists {
			tc.lang[lang] = true
			if strings.HasSuffix(lang, ".raw") {
				tc.raw = true
			} else if strings.HasSuffix(lang, ".utf8") {
				tc.utf8 = true
			}
		}
	}
}

func (tc *TextCat) EnableAllRawLanguages() {
	for lang := range tc.lang {
		if strings.HasSuffix(lang, ".raw") {
			tc.lang[lang] = true
			tc.raw = true
		}
	}
}

func (tc *TextCat) EnableAllUtf8Languages() {
	for lang := range tc.lang {
		if strings.HasSuffix(lang, ".utf8") {
			tc.lang[lang] = true
			tc.utf8 = true
		}
	}
}

func (tc *TextCat) Classify(text string) (languages []string, err error) {
	var mydata map[string]int

	languages = make([]string, 0, tc.maxCandidates)

	if tc.raw && len(text) < tc.minDocSize {
		err = errShort
		return
	}
	if tc.utf8 && utf8.RuneCountInString(strings.TrimSpace(reInvalid.ReplaceAllString(text, " "))) < tc.minDocSize {
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
		patt := GetPatterns(text, utf8)
		suffix := ".raw"
		if utf8 {
			suffix = ".utf8"
		}
		for lang := range tc.lang {
			if !tc.lang[lang] || !strings.HasSuffix(lang, suffix) {
				continue
			}
			if _, ok := tc.extra[lang]; ok {
				mydata = tc.extra[lang]
			} else {
				mydata = data[lang]
			}
			score := 0
			for n, p := range patt {
				i, ok := mydata[p.S]
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
	threshold := float64(minScore) * tc.thresholdValue
	nCandidates := 0
	for _, sco := range scores {
		if float64(sco.score) <= threshold {
			nCandidates += 1
		}
	}
	if nCandidates > tc.maxCandidates {
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

// Add language created by 'textpat' in package github.com/pebbe/textcat/textpat
func (tc *TextCat) AddLanguage(language, filename string) error {
	rawlines := make([]string, 0, 400)
	utf8lines := make([]string, 0, 400)
	target := 0
	r, err := util.NewLinesReaderFromFile(filename)
	if err != nil {
		return err
	}
	for line := range r.ReadLines() {
		if line == "[[[RAW]]]" {
			target = 1
		} else if line == "[[[UTF8]]]" {
			target = 2
		} else {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				if target == 1 {
					rawlines = append(rawlines, fields[0])
				} else if target == 2 {
					utf8lines = append(utf8lines, fields[0])
				}
			}
		}
	}
	if len(rawlines) == 0 && len(utf8lines) == 0 {
		return errors.New("No patterns found in file \"" + filename + "\"")
	}

	if len(rawlines) > 0 {
		a := make(map[string]int)
		for i, p := range rawlines {
			if i == 400 {
				break
			}
			a[p] = i
		}
		l := language + ".raw"
		tc.extra[l] = a
		if _, ok := tc.lang[l]; !ok {
			tc.lang[l] = false
		}
	}

	if len(utf8lines) > 0 {
		a := make(map[string]int)
		for i, p := range utf8lines {
			if i == 400 {
				break
			}
			a[p] = i
		}
		l := language + ".utf8"
		tc.extra[l] = a
		if _, ok := tc.lang[l]; !ok {
			tc.lang[l] = false
		}
	}

	return nil
}
