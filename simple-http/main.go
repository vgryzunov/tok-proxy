package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <port>", os.Args[0])
	}

	port := os.Args[1]

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("Invalid port: %s (%s)\n", os.Args[1], err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		println("--->", port, req.URL.String())

		fmt.Fprintf(w, "%s %s %s\n", req.Method, req.URL, req.Proto)

		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
		fmt.Fprintf(w, "Host = %q\n", req.Host)
		fmt.Fprintf(w, "RemoteAddr = %q\n", req.RemoteAddr)
	})

	log.Printf("Listening port %s", port)
	http.ListenAndServe(":"+port, nil)
}
