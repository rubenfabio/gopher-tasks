package main

import (
	"net/http"
	"os"

	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/logger"
)

func main() {
	log := logger.New("debug", "text", os.Stdout)

    log.Info("Starting gopher-tasks server")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.WithField("path", r.URL.Path).Info("health check")
        w.Write([]byte("gopher-tasks is running!"))
    })

    addr := ":8080"
    log.Info("Listening on", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}
