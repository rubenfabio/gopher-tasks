package main

import (
	"fmt"
	"log"
	"net/http"
)

// handler de health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("gopher-tasks is running!"))
}

func main() {
    // rota raiz
    http.HandleFunc("/", healthHandler)

    addr := ":8080"
    fmt.Println("Starting server on", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}
