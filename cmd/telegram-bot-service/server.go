package main

import (
	"fmt"
	"net/http"
	"os"
)

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("it works!"))
	})

	port := os.Getenv("PORT")
	fmt.Println("listening on port " + port)

	go http.ListenAndServe(":"+port, mux)
}
