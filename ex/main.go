// ## Objectif
// Créer une API en Go combinant **middleware**, **context** et **timeout**.

// L’API devra :

// 1. Ajouter un identifiant de requête (`requestID`) à chaque requête via un middleware.
// 2. Simuler un traitement long dans un endpoint `/process_order`.
// 3. Respecter un timeout fixé par `context.WithTimeout`.
// 4. Retourner le `requestID` dans la réponse pour vérifier le middleware.

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := "123abc"
		ctx := context.WithValue(r.Context(), "reqID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ProcessOrder(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Second)
	defer cancel()

	fmt.Printf(req.Context().Value("reqID").(string))

	select {
	case <-time.After(4 * time.Second):
		io.WriteString(w, "Order has been processed")
	case <-ctx.Done():
		http.Error(w, "TimeOut, processing was too long", http.StatusGatewayTimeout)
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /process_order", ProcessOrder)

	wrappedMux := middleware(mux)

	server := &http.Server{
		Addr:    ":3000",
		Handler: wrappedMux,
	}

	server.ListenAndServe()

}
