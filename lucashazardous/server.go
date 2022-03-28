package main

import (
	"library_server/lucashazardous/book_handler"
	"net/http"
)

func main() {
	adminPanel := book_handler.NewAdminPanel()
	bookHandler := book_handler.NewBookHandlers(adminPanel)
	fileServer := http.FileServer(http.Dir("./static"))

	go http.Handle("/", http.StripPrefix("/", fileServer))

	go http.HandleFunc("/books", bookHandler.Books)
	go http.HandleFunc("/books/", bookHandler.GetBook)

	go http.HandleFunc("/admin", adminPanel.AdminHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
