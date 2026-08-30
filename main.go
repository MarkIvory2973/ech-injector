package main

import (
	"ech-injector/internal/handlers"

	"github.com/syumai/workers"
)

func main() {
	handlers.HandleRFC8484()
	handlers.HandleJSON()
	handlers.HandleFallback()
	workers.Serve(nil)
}
