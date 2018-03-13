package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [files]\n", os.Args[0])
		os.Exit(1)
	}

	var s scanner.Scanner
	fset := token.NewFileSet() // positions are relative to fset

	counts := make(map[string]int)

	for _, arg := range os.Args[1:] {
		src, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}

		file := fset.AddFile(arg, fset.Base(), len(src)) // register input "file"
		s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

		// Repeated calls to Scan yield the token sequence found in the input.
		for {
			_, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			if tok == token.IDENT {
				// fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
				counts[lit]++
			}
		}
	}

	type pair struct {
		s string
		n int
	}

	pairs := make([]pair, 0, len(counts))

	for s, n := range counts {
		// fmt.Printf("%q\t%d\n", s, n)
		pairs = append(pairs, pair{s, n})
	}

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].n > pairs[j].n })

	for i := 0; i < len(pairs) && i < 5; i++ {
		fmt.Printf("%6d %s\n", pairs[i].n, pairs[i].s)
	}
}
