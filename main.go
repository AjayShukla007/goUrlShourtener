package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UrlData struct {
	ID        string    `json:"id"`
	OgUrl     string    `json:"ogUrl"`
	ShortUrl  string    `json:"shortUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

var urlData = make(map[string]UrlData)

func generateShortUrl(ogUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(ogUrl))

	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:6]
}
func createShortUrl(ogUrl string) string {
	shortUrl := generateShortUrl(ogUrl)
	id := shortUrl
	urlData[id] = UrlData{
		ID:        id,
		OgUrl:     ogUrl,
		ShortUrl:  shortUrl,
		CreatedAt: time.Now(),
	}
	return shortUrl
}
func getFullUrl(id string) (UrlData, error) {
	url, ok := urlData[id]
	// fmt.Println(urlData)
	if !ok {
		return UrlData{}, errors.New("URL not found")
	}
	return url, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to URL shortener!")
}
func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Url string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl := createShortUrl(request.Url)
	// fmt.Fprintf(w, shortUrl)
	response := struct {
		ShortUrl string `json:"shortUrl"`
	}{ShortUrl: shortUrl}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getFullUrl(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OgUrl, http.StatusFound)
}

func main() {
	fmt.Println("Hello, World!")
	// shortUrl := createShortUrl("https://www.google.com")
	// fmt.Println(shortUrl)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/shorten", shortUrlHandler)
	http.HandleFunc("/redirect/{id}", redirectHandler)
	// startign the server
	fmt.Println("Server is running on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
