package book_handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Book struct {
	Title  string  `json:"title"`
	Price  float32 `json:"price"`
	Year   int     `json:"year"`
	Author string  `json:"author"`
	Id     string  `json:"id"`
}

type bookHandler struct {
	sync.Mutex
	library map[string]Book
	panel   adminPanel
}

func NewLibraryHandler(panel *adminPanel) *bookHandler {
	return &bookHandler{
		library: map[string]Book{
			"1": {
				Title:  "The Remains of the Day",
				Price:  100.00,
				Year:   1989,
				Author: "Kazuo Ishiguro",
				Id:     "1",
			},
		},
		panel: *panel,
	}
}

func (b *bookHandler) BooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		b.getBooks(w, r)
		return
	case "POST":
		b.createBook(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (b *bookHandler) SpecificBookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		b.deleteSpecificBookById(w, r)
	case "GET":
		b.getSpecificBookById(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (b *bookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if username != "admin" || password != b.panel.password || !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("need content application/json"))
		return
	}

	var book Book
	err = json.Unmarshal(bodyBytes, &book)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book.Id = fmt.Sprintf("%d", time.Now().UnixNano())

	b.Lock()
	b.library[book.Id] = book
	b.Unlock()
}

func (b *bookHandler) getBooks(w http.ResponseWriter, r *http.Request) {

	books := make([]Book, len(b.library))

	b.Lock()
	i := 0
	for _, book := range b.library {
		books[i] = book
		i++
	}
	b.Unlock()

	jsonBytes, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (b *bookHandler) deleteSpecificBookById(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	if len(urlParts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b.Lock()
	delete(b.library, urlParts[2])
	b.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (b *bookHandler) getSpecificBookById(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	if len(urlParts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b.Lock()
	book, ok := b.library[urlParts[2]]
	b.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

type adminPanel struct {
	password string
}

func NewAdminPanel() *adminPanel {
	password, err := ioutil.ReadFile("./password.txt")
	if err != nil || string(password) == "" {
		panic("error reading password")
	}
	return &adminPanel{password: string(password)}
}

func (a adminPanel) AdminHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.getAdminWebsite(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (a adminPanel) getAdminWebsite(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if username != "admin" || password != a.password || !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	adminPage, err := ioutil.ReadFile("./admin_static/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to read admin panel"))
		return
	}

	w.Write(adminPage)
}
