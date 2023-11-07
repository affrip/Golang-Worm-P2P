package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var secret int = randInt(1, 10000000)

func flushHttp(rw http.ResponseWriter, req *http.Request) {
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}
}

func downloadBot(w http.ResponseWriter, r *http.Request) {
	var botcode []byte

	rd, err := os.Open(os.Args[0])
	if err != nil {
		fmt.Fprintf(w, "Error opening file")
		return
	}

	botcode, err = io.ReadAll(rd)
	if err != nil {
		fmt.Fprintf(w, "Error reading code")
		return
	}

	w.Write(botcode)
	flushHttp(w, r)
}

func downloadBook(w http.ResponseWriter, r *http.Request) {
	var book []byte

	rd, err := os.Open(".book.dat")
	if err != nil {
		fmt.Fprintf(w, "Error opening file")
		return
	}

	book, err = io.ReadAll(rd)
	if err != nil {
		fmt.Fprintf(w, "Error reading code")
		return
	}

	w.Write(book)
	flushHttp(w, r)
}

func setupRoutes() {
	_ = os.Mkdir(".static", os.ModePerm)
	http.HandleFunc("/api/v1/get", downloadBot)
	http.HandleFunc("/api/v1/book", downloadBook)
	fs := http.FileServer(http.Dir("./.static"))
	http.Handle("/", fs)
}

func setupServ() {
	http.ListenAndServe("0.0.0.0:80", nil)
}

func setupServAlt() {
	http.ListenAndServe("0.0.0.0:8888", nil)
}

func init_web() {
	setupRoutes()
	go setupServ()
	go setupServAlt()
}
