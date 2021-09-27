package request

import (
	"fmt"
	"context"
	"net/http"
)

func Request(ctx context.Context, method string, urlStr string, setH string, setV stinger) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, urlStr, nil)
	if err != nil {
		return nil, err
	}

	if len(setH) != 0 {
		req.Header.Set(setH, setV)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request.Request err: %s", err)
	}
	return resp, nil
}