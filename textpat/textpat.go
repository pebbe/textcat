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
		for i, p := range textcat.GetPatterns(str, false) {
			if i == 400 {
				break
			}
			fmt.Printf("%s\t%d\n", p.S, p.I)
		}
	}

	if doUtf8 {
		fmt.Println("[[[UTF8]]]")
		for i, p := range textcat.GetPatterns(str, true) {
			if i == 400 {
				break
			}
			fmt.Printf("%s\t%d\n", p.S, p.I)
		}
	}

}

func syntax() {
	fmt.Fprintf(os.Stderr, `
Usage: %s [-r|-u] < sample-data

Reads text samples from standard input, write to standard output
text patterns for package github.com/pebbe/textcat

Options:

    -r : raw patterns only
    -u : utf8 patterns only

`, os.Args[0])
}
