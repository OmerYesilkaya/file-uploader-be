package middleware

import (
	"net/http"
	"strings"

	"github.com/OmerYesilkaya/fileuploader/internal/utils"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			utils.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		const bearerPrefix = "Bearer "
		if strings.HasPrefix(token, bearerPrefix) {
			token = strings.TrimSpace(token[len(bearerPrefix):])
		}

		_, err := utils.ParseJWT(token)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})

}
