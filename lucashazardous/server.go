package main

import (
	"net/http"
	"library_server/lucashazardous/book_handler"
)

func main() {
	bookHandler := book_handler.NewBookHandlers()
	go http.HandleFunc("/books", bookHandler.Books)
	go http.HandleFunc("/books/", bookHandler.GetBook)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
