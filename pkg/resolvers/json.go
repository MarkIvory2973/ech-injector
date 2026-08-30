package resolvers

import (
	"context"
	"ech-injector/pkg/workers"
	"encoding/json"
	"net/http"
)

const hostJSON = "dns.google"
const pathJSON = "/resolve"

// TODO
func Resolve(context context.Context, queries map[string]string) (map[string]any, error) {
	httpRequest := workers.HTTPRequest{
		Method:  http.MethodGet,
		Scheme:  "https",
		Host:    hostJSON,
		Path:    pathJSON,
		Queries: queries,
		Headers: map[string]string{
			"Accept": "application/json",
		},
	}
	content, err := workers.Fetch(context, httpRequest)
	if err != nil {
		return nil, err
	}

	var dnsJSON map[string]any
	err = json.Unmarshal(content, &dnsJSON)
	if err != nil {
		return nil, err
	}

	return dnsJSON, err
}
