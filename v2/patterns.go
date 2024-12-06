package textcat

import (
	"regexp"
	"sort"
	"strings"
)

const (
	MaxPatterns = 400
)

var reInvalid = regexp.MustCompile("[ \t\f\r\n_0-9]+")

type countType struct {
	S string
	I int
}

type countsType []*countType

func (c countsType) Len() int {
	return len(c)
}

func (c countsType) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c countsType) Less(i, j int) bool {
	if c[i].I != c[j].I {
		return c[i].I > c[j].I
	}
	return c[i].S < c[j].S
}

func GetPatterns(s string, useRunes bool) []*countType {
	ngrams := make(map[string]int)
	if useRunes {
		s = reInvalid.ReplaceAllString(s, " ")
		for _, word := range strings.Fields(s) {
			b := []rune("_" + word + "____")
			n := len(b) - 4
			for i := 0; i < n; i++ {
				for j := 1; j < 6; j++ {
					s = string(b[i : i+j])
					if !strings.HasSuffix(s, "__") {
						ngrams[s] += 1
					}
				}
			}
		}
	} else {
		for _, word := range strings.Fields(s) {
			b := []byte("_" + word + "____")
			n := len(b) - 4
			for i := 0; i < n; i++ {
				for j := 1; j < 6; j++ {
					s = string(b[i : i+j])
					if !strings.HasSuffix(s, "__") {
						ngrams[s] += 1
					}
				}
			}
		}
	}
	size := len(ngrams)
	counts := make([]*countType, 0, size)
	for i := range ngrams {
		counts = append(counts, &countType{i, ngrams[i]})
	}
	sort.Sort(countsType(counts))
	if size > MaxPatterns {
		counts = counts[:MaxPatterns]
	}
	return counts
}
