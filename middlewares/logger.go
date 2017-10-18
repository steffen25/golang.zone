package middlewares

import (
	"log"
	"net/http"
	"time"
)

// Logger logs the current request to the console printing the date, HTTP method, path and elapsed time
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next(res, req)
		log.Printf("[%s] %q %v\n", req.Method, req.URL.String(), time.Since(start))
	}
}
