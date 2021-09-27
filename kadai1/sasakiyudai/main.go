package main

import (
	"flag"
	"log"
	"fmt"
	"imgconv/imgconv"
)

type Flags struct {
	src string
	dst string
	dir string
}

var flags Flags

func init() {
	flag.StringVar(&flags.src, "src", "jpg", "src extension [jpg]jpeg|png|gif]")
	flag.StringVar(&flags.dst, "dst", "png", "dst extension [jpg]jpeg|png|gif]")
}

func validateArg(ext string) error {
	switch ext {
		case "jpg", "jpeg", "png", "gif":
			return nil
		default:
			return fmt.Errorf("invalid ext: %s", ext)
	}
}

func parseArgs() error {
	flag.Parse()

	if err := validateArg(flags.src); err != nil {
		return err
	}

	if err := validateArg(flags.dst); err != nil {
		return err
	}

	if flag.Arg(0) == "" {
		return fmt.Errorf("directory is required")
	} else {
		flags.dir = flag.Arg(0)
	}

	return nil
}

func main() {
	if err := parseArgs(); err != nil {
		log.Fatal(err)
	}

	if err := imgconv.Converter(flags.dir, flags.src, flags.dst); err != nil {
		log.Fatal(err)
	}
}