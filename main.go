package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {

	// middleware pour ajouter une valeur dans un contexte
	requestMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			reqID := 155
			ctx := context.WithValue(req.Context(), "requestID", reqID)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
		if req.Context().Value("requestID") != nil {
			fmt.Printf("Request ID: %v\n", req.Context().Value("requestID"))
			io.WriteString(w, "Request ID: "+fmt.Sprint(req.Context().Value("requestID"))+"\n")
		}
	}

	longLoadingHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		defer cancel()

		select {
		case <-time.After(3 * time.Second):
			io.WriteString(w, "Réponse calculée\n")
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("Contexte annulé car timeout atteint")
				http.Error(w, "Timeout atteint", http.StatusGatewayTimeout)
			} else if ctx.Err() == context.Canceled {
				fmt.Println("Contexte annulé:", ctx.Err())
				http.Error(w, "Contexte annulé", http.StatusRequestTimeout)
			}
		}

	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", helloHandler)
	mux.HandleFunc("GET /long_loading", longLoadingHandler)

	middleware := requestMiddleware(mux)

	server := &http.Server{
		Addr:         ":3000",
		Handler:      middleware,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
