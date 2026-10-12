package main

import (
	httpapi "backend-challange/internal/httpAPI"
	"log"
	"net/http"
)

//cors midleware
func cors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/wallets", httpapi.HandlerWallet)

	log.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", cors(mux)))

}
