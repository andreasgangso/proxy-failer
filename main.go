package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	requestCounts = make(map[string]int)
	mu            sync.Mutex
)

func main() {
	http.HandleFunc("/", handleProxy)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	// Extract target URL from the path
	targetURLStr := r.URL.Query().Get("proxy")
	println(targetURLStr)
	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Check request count for this URL
	mu.Lock()
	if count, ok := requestCounts[targetURLStr]; ok && count == 0 {
		requestCounts[targetURLStr]++
		mu.Unlock()
		http.Error(w, "Simulated failure", http.StatusInternalServerError)
		return
	}
	requestCounts[targetURLStr] = 0
	mu.Unlock()

	// Proxy the request
	r.Host = targetURL.Host
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}
