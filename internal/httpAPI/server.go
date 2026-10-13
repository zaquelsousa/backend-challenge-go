package httpapi

import (
	"net/http"
)

//cors midleware
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewServer(handler *WalletHandler) http.Handler {
	mux := http.NewServeMux()
	//POST endpoints
	mux.HandleFunc("/wallets", handler.HandlerCreateWallet)

	//GET endpoints
	mux.HandleFunc("/wallets/{walletId}", handler.HandlerGetWalletById)

	return cors(mux)
}

