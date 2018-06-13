package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := []string{"."}
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	for _, arg := range args {
		err := tree(arg)
		if err != nil {
			log.Printf("tree %s: %v", arg, err)
		}
	}
}

func tree(root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name()[0] == '.' {
			return filepath.SkipDir
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("could not Rel(%s, %s): %v", root, path, err)
		}

		depth := len(strings.Split(rel, string(filepath.Separator)))
		fmt.Println(strings.Repeat("-", depth), info.Name())
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", root, err)
	}

	return nil
}
