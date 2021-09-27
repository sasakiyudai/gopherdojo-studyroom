package getheader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func ResHeader(w io.Writer, r *http.Response, header string) ([]string, error) {
	h, is := r.Header[header]
	if !is {
		return nil, fmt.Errorf("cannot find %s header", header)
	}
	if _, err := fmt.Fprintf(w, "Header[%s] = %s\n", header, h); err != nil {
		return nil, err
	}
	return h, nil
}

func GetSize(resp *http.Response) (uint, error) {
	contLen, err := ResHeader(os.Stdout, resp, "Content-Length")
	if err != nil {
		return 0, err
	}
	ret, err := strconv.ParseUint(contLen[0], 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(ret), nil
}