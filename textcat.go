/*
A Go package for n-gram based text categorization, with support for utf-8 and raw text.
*/
package textcat

var (
	LanguagesUtf8 = make([]string, 0)
	LanguagesRaw  = make([]string, 0)
)

func init() {
	for _, lang := range data {
		LanguagesRaw = append(LanguagesRaw, lang.lang)
		if len(lang.patUtf8) > 0 {
			LanguagesUtf8 = append(LanguagesUtf8, lang.lang)
		}
	}
}
