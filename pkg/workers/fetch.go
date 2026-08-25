package workers

import (
	"bytes"
	"context"
	"io"
	"net/url"

	"github.com/syumai/workers/cloudflare/fetch"
)

type HTTPRequest struct {
	Method  string
	Scheme  string
	Host    string
	Path    string
	Queries map[string]string
	Headers map[string]string
	Content []byte
}

func (httpRequest HTTPRequest) URL() string {
	url := &url.URL{
		Scheme: httpRequest.Scheme,
		Host:   httpRequest.Host,
		Path:   httpRequest.Path,
	}

	queries := url.Query()
	for key, value := range httpRequest.Queries {
		queries.Set(key, value)
	}
	url.RawQuery = queries.Encode()

	content := url.String()

	return content
}

var client *fetch.Client

func init() {
	client = fetch.NewClient()
}

func Fetch(context context.Context, httpRequest HTTPRequest) ([]byte, error) {
	var body io.Reader
	if httpRequest.Content != nil {
		body = bytes.NewReader(httpRequest.Content)
	}

	request, err := fetch.NewRequest(
		context,
		httpRequest.Method,
		httpRequest.URL(),
		body,
	)
	if err != nil {
		return nil, err
	}

	for key, value := range httpRequest.Headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
