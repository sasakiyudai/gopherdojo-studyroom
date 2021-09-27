package main

import (
	"fmt"
	"log"
	"context"
	"os"
	"bytes"
	"errors"
	"net/url"
	"runtime"
	"time"
	"path/filepath"
	"strconv"

	flags "github.com/jessevdk/go-flags"

	"github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/request"
	"github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/getheader"
	"github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/listen"
	"github.com/sasakiyudai/gopherdojo-studyroom/kadai3-2/sasakiyudai/download"
)

type Options struct {
	Help   bool   `short:"h" long:"help"`
	Procs  uint   `short:"P" long:"procs"`
	Output string `short:"o" long:"output" default:"./"`
	Tm     int    `short:"t" long:"timeout" default:"120"`
}

func (opts *Options) parse(argv []string) ([]string, error) {
	p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.ParseArgs(argv)
	if err != nil {
		_, err2 := os.Stderr.Write(opts.usage())
		if err2 != nil {
			return nil, fmt.Errorf("%s: invalid command line options: cannot print usage: %s", err, err2)
		}
		return nil, fmt.Errorf("%w: invalid command line options", err)
	}

	return args, nil
}

func (opts Options) usage() []byte {
	buf := bytes.Buffer{}

	fmt.Fprintln(&buf,
		`Usage: paraDW [options] URL (URL2, URL3, ...)
	Options:
	-h,   --help              print usage and exit
	-p,   --procs <num>       the number of split to download (default: the number of CPU cores)
	-o,   --output <filename> path of the file downloaded (default: current directory)
	-t,   --timeout <num>     Time limit of return of http response in seconds (default: 120)
	`,
	)

	return buf.Bytes()
}

func main() {
	var opts Options
	argv := os.Args[1:]
	if len(argv) == 0 {
		if _, err := os.Stdout.Write(opts.usage()); err != nil {
			log.Fatalf("err: %s: %s\n", errors.New("no options"), err)
		}
		log.Fatalf("err: %s\n", errors.New("no options"))
	}

	urlsStr, err := opts.parse(argv)
	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	var urls []*url.URL
	for _, u := range urlsStr {
		url, err := url.ParseRequestURI(u)
		if err != nil {
			log.Fatalf("err: url.ParseRequestURI: %s\n", err)
		}
		urls = append(urls, url)
	}

	fmt.Printf("timeout: %d\n", opts.Tm)

	if opts.Help {
		if _, err := os.Stdout.Write(opts.usage()); err != nil {
			log.Fatalf("err: cannot print usage: %s", err)
		}
		log.Fatal(errors.New("print usage"))
	}

	if opts.Procs == 0 {
		opts.Procs = uint(runtime.NumCPU())
	}

	if len(opts.Output) > 0 && opts.Output[len(opts.Output)-1] != '/' {
		opts.Output += "/"
	}

	for i, urlObj := range urls {
		downloadFromUrl(i, opts, urlObj)
	}
}

func downloadFromUrl(i int, opts Options, urlObj *url.URL) {
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), time.Duration(opts.Tm)*time.Second)
	defer cancelTimeout()

	resp, err := request.Request(ctxTimeout, "HEAD", urlObj.String(), "", "")
	if err != nil {
		log.Fatalf("err: &s\n", err)
	}

	fileSize, err := getheader.Getsize(resp)
	if err != nil {
		log.Fatalf("err: getheader.Getsize: %s\n", err)
	}
	if err = resp.Body.Close(); err != nil {
		log.Fatalf("err: %s", err)
	}

	partial := fileSize / opts.Procs

	outputPath := opts.Output + filepath.Base(urlObj.String())
	if isExists(outputPath) {
		err := os.Remove(outputPath)
		if err != nil {
			log.Fatalf("err: isExists: os.Remove: %s\n", err)
		}
	}

	out, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("err: os.Create: %s\n", err)
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Fatalf("err: %s", err)
		}
	}()
	
	tmpDirName := opts.Output + strconv.Itoa(i)
	err = os.Mkdir(tmpDirName, 0777)
	if err != nil {
		if err3 := out.Close(); err3 != nil {
			log.Fatalf("err: %s", err3)
		}
		if err2 := os.Remove(opts.Output + filepath.Base(urlObj.String())); err2 != nil {
			log.Fatalf("err: os.Mkdir: %s\nerr: os.Remove: %s\n", err, err2)
		}
		log.Fatalf("err: os.Mkdir: %s\n", err)
	}

	clean := func() {
		if err := out.Close(); err != nil {
			log.Fatalf("err: out.Close: %s\n", err)
		}
		if err := os.RemoveAll(tmpDirName); err != nil {
			log.Fatalf("err: RemoveAll: %s\n", err)
		}
		if err := os.Remove(opts.Output + filepath.Base(urlObj.String())); err != nil {
			log.Fatalf("err: os.Remove: %s\n", err)
		}
	}
	ctx, cancel := listen.Listen(ctxTimeout, os.Stdout, clean)
	defer cancel()

	var isPara bool = true
	_, err = getheader.ResHeader(os.Stdout, resp, "Accept-Ranges")
	if err != nil && err.Error() == "cannot find Accept-Ranges header" {
		isPara = false
	} else if err != nil {
		clean()
		log.Fatalf("err: getheader.ResHeader: %s\n", err)
	}

	// err = download.Downloader(urlObj, out, fileSize, partial, opts.Procs, isPara, tmpDirName, ctx)
	// if err != nil {
	// 	log.Fatalf("err: %s\n", err)
	// }

	fmt.Printf("download complete: %s\n", urlObj.String())

	
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}