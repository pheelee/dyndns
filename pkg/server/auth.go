package server

import (
	"net/http"
)

func basicAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok == true && user == appConfig.HTTPReqAuth.Username && pass == appConfig.HTTPReqAuth.Password {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	}

	return http.HandlerFunc(fn)
}
