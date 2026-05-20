package main

import (
	"fmt"
	"net/http"

	"ajirascan/internal/web"
)

func main() {

	http.HandleFunc("/", web.HomeHandler)

	http.HandleFunc("/analyze", web.HomeHandler)

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	fmt.Println("Server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}