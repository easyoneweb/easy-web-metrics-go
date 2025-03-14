package main

import "net/http"

func (apiCfg *ApiConfig) middlewareCheckApiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api-key")
		if apiCfg.ApiKey != apiKey {
			w.WriteHeader(401)
			w.Write([]byte("not allowed"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
