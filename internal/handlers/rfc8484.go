package handlers

import (
	"ech-injector/internal/injectors"
	"ech-injector/pkg/logs"
	"ech-injector/pkg/resolvers"
	"encoding/base64"
	"io"
	"net/http"
)

func HandleRFC8484() {
	http.HandleFunc("/dns-query", func(writer http.ResponseWriter, request *http.Request) {
		var content []byte
		var err error

		switch request.Method {
		case http.MethodGet:
			dns := request.URL.Query().Get("dns")
			content, err = base64.RawURLEncoding.DecodeString(dns)
		case http.MethodPost:
			body := http.MaxBytesReader(writer, request.Body, 65536)
			content, err = io.ReadAll(body)
		default:
			writer.Header().Set("Allow", "GET, POST")
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err != nil {
			logs.Warning("handlers.HandleRFC8484", "couldn't parse HTTP request", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		dnsQuestion, err := resolvers.UnpackMessage(content)
		if err != nil {
			logs.Warning("resolvers.UnpackMessage", "couldn't unpack the DNS message", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		dnsAnswer, err := resolvers.ExchangeMessage(request.Context(), dnsQuestion)
		if err != nil {
			logs.Warning("resolvers.ExchangeMessage", "couldn't exchange the DNS message", err)
			writer.WriteHeader(http.StatusBadGateway)
			return
		}

		err = injectors.InjectRFC8484(request.Context(), dnsQuestion, dnsAnswer)
		if err != nil {
			logs.Warning("injectors.InjectDNSMessage", "couldn't inject ECH configuration into the DNS message", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		content, err = resolvers.PackMessage(dnsAnswer)
		if err != nil {
			logs.Warning("resolvers.PackMessage", "couldn't pack the DNS message", err)
			writer.WriteHeader(http.StatusInternalServerError)
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
