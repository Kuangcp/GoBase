package main

import (
	"embed"
	"log"
	"net/http"
	"os"
)

//go:embed frontend
var frontendFS embed.FS

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(frontendFS))

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
