/*
The program `textpat` creates language patterns that can be loaded into
running Go programs that are using the textcat library.

See function `AddLanguage` in library `textcat`.

Usage:

    textpat [-r|-u] < sample data

Reads text samples from standard input, write to standard output
text patterns for package github.com/pebbe/textcat

Options:

    -r : raw patterns only
    -u : utf8 patterns only
*/
package main

import (
	"fmt"
	"github.com/pebbe/textcat"
	"github.com/pebbe/util"
	"io/ioutil"
	"os"
)

func main() {
	doUtf8 := true
	doRaw := true

	if util.IsTerminal(os.Stdin) {
		syntax()
		return
	}

	for _, arg := range os.Args[1:] {
		switch arg {
		case "-r":
			doUtf8 = false
		case "-u":
			doRaw = false
		default:
			syntax()
			return
		}
	}
	if !doUtf8 && !doRaw {
		syntax()
		return
	}

	data, err := ioutil.ReadAll(os.Stdin)
	util.CheckErr(err)
	str := string(data)

	if doRaw {
		fmt.Println("[[[RAW]]]")
		n := 0
		for i, p := range textcat.GetPatterns(str, false) {
			if i == textcat.MaxPatterns {
				break
			}
			n += 1
			fmt.Printf("%s\t%d\n", p.S, p.I)
		}
		if n < textcat.MaxPatterns {
			fmt.Fprintf(os.Stderr, "Warning: there are less than %d raw patterns\n", textcat.MaxPatterns)
		}
	}

	if doUtf8 {
		fmt.Println("[[[UTF8]]]")
		n := 0
		for i, p := range textcat.GetPatterns(str, true) {
			if i == textcat.MaxPatterns {
				break
			}
			n += 1
			fmt.Printf("%s\t%d\n", p.S, p.I)
		}
		if n < textcat.MaxPatterns {
			fmt.Fprintf(os.Stderr, "Warning: there are less than %d utf8 patterns\n", textcat.MaxPatterns)
		}
	}

}

func syntax() {
	fmt.Fprintf(os.Stderr, `
Usage: %s [-r|-u] < sample data

Reads text samples from standard input, write to standard output
text patterns for package github.com/pebbe/textcat

Options:

    -r : raw patterns only
    -u : utf8 patterns only

`, os.Args[0])
}
