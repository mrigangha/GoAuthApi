package services

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const origins = "http://localhost:5000"

func getIP(r *http.Request) string {
	// If behind a trusted proxy, use X-Forwarded-For
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP in the list (the original client)
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If SplitHostPort fails, return the full RemoteAddr
		return r.RemoteAddr
	}
	return ip
}

var (
	rateLimiter = make(map[string]float64)
	mu          sync.Mutex // Protect the map from concurrent access
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		currentTime := float64(time.Now().Unix())

		mu.Lock()
		lastRequest, exists := rateLimiter[ip]
		if exists && currentTime-lastRequest < 1 {
			mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rateLimiter[ip] = currentTime
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers BEFORE calling next.ServeHTTP
		w.Header().Set("Access-Control-Allow-Origin", origins)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println(getIP(r))
		next.ServeHTTP(w, r)
	})
}
