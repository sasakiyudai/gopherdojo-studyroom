package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "unicode"
	"bufio"
)

var p [1]byte

func readbyte(r io.Reader) (rune, error) {
    n, err := r.Read(p[:])
    if n > 0 {
        return rune(p[0]), nil
    }
    return 0, err
}

func main() {
    filename := os.Args[1]
    f, err := os.Open(filename)
    if err != nil {
        log.Fatalf("cannot open file %q: %v", filename, err)
    }
    defer f.Close()

	b := bufio.NewReader(f)

    words := 0
    inword := false

    for {
        r, err := readbyte(b)
        if unicode.IsSpace(r) {
            if inword {
                words++
            }
            inword = false
        } else {
            inword = true
        }

        if err == io.EOF {
            break
        }

        if err != nil {
            log.Fatalf("read failed: %v", err)
        }
    }

    fmt.Printf("%q: %d words\n", filename, words)
}
