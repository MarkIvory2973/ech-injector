package handlers

import (
	"encoding/base64"
	"io"
	"net/http"
	"ech-injector/internal/injectors"
)

func HandleRFC8484() {
	http.HandleFunc("/dns-query", func(writer http.ResponseWriter, request *http.Request) {
		var content []byte
		var err error

		switch request.Method {
		case "GET":
			dns := request.URL.Query().Get("dns")
			if dns == "" {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			content, err = base64.RawURLEncoding.DecodeString(dns)
		case "POST":
			content, err = io.ReadAll(request.Body)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
			return
		}

		content, err = injectors.InjectRFC8484(request.Context(), content)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}

		writer.Header().Set("Content-Type", "application/dns-message")
		writer.WriteHeader(http.StatusOK)
		writer.Write(content)
	})
}
