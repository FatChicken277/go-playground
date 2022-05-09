package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	methodNotAllowedError = "method is not allowed"
	internalServerError   = "internal error"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, methodNotAllowedError, http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Hello, World")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowedError, http.StatusMethodNotAllowed)
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, internalServerError, http.StatusInternalServerError)
	}

	name := r.FormValue("name")
	address := r.FormValue("address")

	fmt.Fprintln(w, "POST request successful")
	fmt.Fprintf(w, "%s, %s", name, address)
}

func main() {
	fileHandler := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileHandler)

	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Serving at port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
