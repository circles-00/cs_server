package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func runHttpServer(wg *sync.WaitGroup, mux *http.ServeMux) {
	defer wg.Done()

	fs := http.FileServer(http.Dir("./src"))
	mux.Handle("/", fs)
}

func runProxyServer(wg *sync.WaitGroup, mux *http.ServeMux) {
	defer wg.Done()

	// TODO: Get from env variable
	proxyURL, err := url.Parse("http://127.0.0.1:5000")
	if err != nil {
		log.Fatal("Error parsing proxy URL: ", err)
	}

	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		targetURL := proxyURL.ResolveReference(r.URL)
		targetURL.Path = strings.Replace(targetURL.Path, "/api", "", 1)

		proxyRequest, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
			return
		}

		proxyRequest.Header = r.Header

		resp, err := http.DefaultClient.Do(proxyRequest)
		if err != nil {
			http.Error(w, "Failed to contact the target server", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for key, value := range resp.Header {
			w.Header()[key] = value
		}
		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Println("Failed to copy response body:", err)
		}
	})
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)

	mux := http.NewServeMux()

	go runHttpServer(&wg, mux)
	go runProxyServer(&wg, mux)

	log.Println("WWW server listening on 3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
	wg.Wait()
}
