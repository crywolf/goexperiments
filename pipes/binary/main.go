package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/crywolf/goexperiments/pipes/imgcat"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "missing paths of images to cat")
		os.Exit(2)
	}

	for _, path := range os.Args[1:] {
		if err := cat(path); err != nil {
			fmt.Fprintf(os.Stderr, "could not cat %s: %v\n", path, err)
		}
	}
}

func cat(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "could not open image")
	}
	defer file.Close()
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		wc := imgcat.NewWriter(pw, false)
		defer wc.Close()
		if _, err = io.Copy(wc, file); err != nil {
			pw.CloseWithError(errors.Wrap(err, "could not copy the image"))
			return
		}
	}()

	newImg, err := os.Create("newImage.png")
	decodedImg := base64.NewDecoder(base64.StdEncoding, pr)
	_, err = io.Copy(newImg, decodedImg)
	if err != nil {
		return errors.Wrap(err, "could not copy to new image")
	}
	return nil
}
