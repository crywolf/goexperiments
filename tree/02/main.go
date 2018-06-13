package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := []string{"."}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	for _, arg := range args {
		err := tree(arg, "")
		if err != nil {
			log.Printf("tree %s: %v", arg, err)
		}
	}
}

func tree(root, ident string) error {
	fi, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("could not Stat(%q): %v", root, err)
	}

	fmt.Printf("%s\n", fi.Name())
	if !fi.IsDir() {
		return nil
	}

	fis, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("could not read dir %s: %v", root, err)
	}

	var names []string

	for _, fi := range fis {
		if fi.Name()[0] != '.' {
			names = append(names, fi.Name())
		}
	}

	for i, name := range names {
		add := "│  "
		if i == len(names)-1 { // last one
			fmt.Print(ident + "└──")
			add = "   "
		} else {
			fmt.Print(ident + "├──")
		}

		if err := tree(filepath.Join(root, name), ident+add); err != nil {
			return err
		}
	}

	return nil
}
