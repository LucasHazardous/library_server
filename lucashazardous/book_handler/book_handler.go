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
}

func (b *bookHandler) Books(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		b.Get(w, r)
		return
	case "POST":
		b.Post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("not allowed"))
		return
	}
}

func (b *bookHandler) Post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
		w.Write([]byte(err.Error()))
	}

	book.Id = fmt.Sprintf("%d", time.Now().UnixNano())

	b.Lock()
	b.library[book.Id] = book
	b.Unlock()
}

func (b *bookHandler) Get(w http.ResponseWriter, r *http.Request) {

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
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func NewBookHandlers() *bookHandler {
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
	}
}

func (b *bookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.String(), "/")
	if len(urlParts) != 3 {
		w.WriteHeader(http.StatusNotFound)
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
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
