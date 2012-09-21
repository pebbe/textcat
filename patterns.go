package textcat

import (
	"regexp"
	"sort"
	"strings"
)

const (
	maxPatterns = 400
)

var (
	re = regexp.MustCompile("[^\\p{L}]")
)

type countType struct {
	s string
	i int
}

type countsType []*countType

func (c countsType) Len() int {
	return len(c)
}

func (c countsType) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c countsType) Less(i, j int) bool {
	if c[i].i != c[j].i {
		return c[i].i > c[j].i
	}
	return c[i].s < c[j].s
}

func getPatterns(s string, useRunes bool) ([]*countType) {
	ngrams := make(map[string]int)
	if useRunes {
		s = re.ReplaceAllString(s, " ")
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
	if size > maxPatterns {
		counts = counts[:maxPatterns]
	}
	return counts
}