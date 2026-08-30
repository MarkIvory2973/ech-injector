package handlers

import (
	"ech-injector/internal/injectors"
	"ech-injector/pkg/logs"
	"ech-injector/pkg/resolvers"
	"encoding/json"
	"net/http"
)

// TODO
func HandleJSON() {
	http.HandleFunc("/resolve", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		queries := map[string]string{
			"name": request.URL.Query().Get("name"),
			"type": request.URL.Query().Get("type"),
		}
		dnsJSON, err := resolvers.Resolve(request.Context(), queries)
		if err != nil {
			logs.Warning("resolvers.Resolve", "couldn't exchange the DNS information in JSON format", err)
			writer.WriteHeader(http.StatusBadGateway)
			return
		}

		dnsJSON, err = injectors.InjectJSON(request.Context(), dnsJSON)
		if err != nil {
			logs.Fatal("injectors.InjectJSON", "couldn't inject ECH configuration into JSON", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		content, err := json.Marshal(dnsJSON)
		if err != nil {
			logs.Warning("json.Marshal", "couldn't marshal the DNS infomation to JSON format", err)
			writer.WriteHeader(http.StatusBadGateway)
			return
		}

		writer.Header().Set("Content-Type", "application/dns-message")
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(content)
		if err != nil {
			logs.Warning("writer.Write", "couldn't write the DNS message", err)
		}
	})
}
