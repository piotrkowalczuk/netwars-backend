package user

import (
	"net/http"
)

func AuthenticationMiddleware(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-API-KEY") == "" {
		res.WriteHeader(http.StatusUnauthorized)
	} else {
		res.WriteHeader(200)
	}
}
