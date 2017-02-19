// https://en.wikipedia.org/wiki/ROT13

package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rotReader rot13Reader) Read(b []byte) (int, error) {
	reader := rotReader.r
	origData := make([]byte, len(b))

	n, err := reader.Read(origData)
	if err != nil {
		return n, err
	}

	for i, v := range origData {
		switch {
		case v >= 'A' && v <= 'Z':
			b[i] = 'A' + (v-'A'+13)%26
		case v >= 'a' && v <= 'z':
			b[i] = 'a' + (v-'a'+13)%26
		default:
			b[i] = v
		}
	}

	return len(b), nil
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
