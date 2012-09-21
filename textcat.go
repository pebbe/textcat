/*
A Go package for n-gram based text categorization, with support for utf-8 and raw text.
*/
package textcat

var Languages []string

func init() {
	Languages = make([]string, 0)
	for _, lang := range data {
		Languages = append(Languages, lang.lang)
	}
}

