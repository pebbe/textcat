package main

import (
	"flag"
	"fmt"
	"github.com/pebbe/textcat"
	"github.com/pebbe/util"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	opt_b = flag.Bool("b", false, "both raw and utf8 patterns")
	opt_r = flag.Bool("r", false, "raw patterns, instead of utf8")
	opt_f = flag.String("f", "", "file name")
	opt_p = flag.String("p", "", "pattern file names, separated by comma's (no spaces)")
)

func main() {
	flag.Parse()

	if *opt_f == "" && flag.NArg() == 0 && util.IsTerminal(os.Stdin) {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [args] [text]\n\nargs with default values are:\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nIf both -f and text are missing, read from stdin\n\n")
		return
	}

	tc := textcat.NewTextCat()
	if *opt_p != "" {
		for _, i := range strings.Split(*opt_p, ",") {
			e := tc.AddLanguage(path.Base(i), i)
			util.CheckErr(e)
		}
	}
	if *opt_r || *opt_b {
		tc.EnableAllRawLanguages()
	}
	if *opt_b || !*opt_r {
		tc.EnableAllUtf8Languages()
	}

	var text string
	if *opt_f != "" {
		fp, err := os.Open(*opt_f)
		util.CheckErr(err)
		t, err := ioutil.ReadAll(fp)
		util.CheckErr(err)
		fp.Close()
		text = string(t)
	} else if flag.NArg() > 0 {
		text = strings.Join(flag.Args(), " ")
	} else {
		t, err := ioutil.ReadAll(os.Stdin)
		util.CheckErr(err)
		text = string(t)
	}

	l, e := tc.Classify(text)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(strings.Join(l, "\n"))
	}
}
