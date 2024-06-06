package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"os/exec"

	"github.com/sourcegraph/conc/pool"
)

func main() {
	var saveMode, listMode, dumpMode, remMode, hidden bool

	flag.BoolVar(&saveMode, "s", false, "")
	flag.BoolVar(&saveMode, "save", false, "")

	flag.BoolVar(&listMode, "l", false, "")
	flag.BoolVar(&listMode, "list", false, "")

	flag.BoolVar(&dumpMode, "d", false, "")
	flag.BoolVar(&dumpMode, "dump", false, "")

	flag.BoolVar(&remMode, "rm", false, "")

	flag.BoolVar(&hidden, "i", false, "")
	flag.BoolVar(&hidden, "hidden", false, "")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	switch {
	case listMode:
		pats, err := listPatterns()
		if err != nil {
			log.Fatalf(errGetPattern, err)
		}

		if len(pats) > 0 {
			println(strings.Join(pats, "\n"))
		}

		return
	case saveMode:
		if flag.NArg() < 3 {
			flag.Usage()
			log.Fatal(errNoPatternInput)
		}

		name := flag.Arg(0)
		flags := flag.Arg(1)
		pattern := flag.Arg(2)

		if err := savePattern(name, flags, pattern); err != nil {
			log.Fatalf(errSavePattern, err)
		}
		return
	}

	patName := flag.Arg(0)
	files := "."
	if flag.NArg() > 1 {
		files = flag.Arg(1)
	}

	if patName == "" {
		flag.Usage()
		log.Fatal(errNoPatternInput)
	}

	pats, err := getPatterns(patName)
	if err != nil {
		log.Fatal(err)
	}

	p := pool.New().WithMaxGoroutines(10)
	for _, pat := range pats {
		pat := pat
		p.Go(func() {
			operator := "grep"
			if pat.Engine != "" {
				operator = pat.Engine
			}

			_, err = exec.LookPath(operator)
			if err != nil {
				log.Fatalf(errOperatorCmdNotFound, operator)
			}

			var patternFlags = "--exclude='.*' " + pat.Flags
			if hidden {
				patternFlags = pat.Flags
			}

			switch {
			case dumpMode:
				fmt.Printf("[%s] %s %s %q %s\n", pat.Filename, operator, patternFlags, pat.Pattern, files)
			case remMode:
				_ = os.Remove(pat.Filepath)
			default:
				var cmd *exec.Cmd

				if isStdin() {
					cmd = exec.Command(operator, pat.Flags, patternFlags)
				} else {
					cmd = exec.Command(operator, pat.Flags, patternFlags, files)
				}
				doSearch(cmd)
			}
		})
	}

	p.Wait()
}
