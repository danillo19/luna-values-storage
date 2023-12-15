package middleware

import (
	"net/http"
)

func CreateCorsAllowOriginFunc() func(r *http.Request, origin string) bool {
	return func(r *http.Request, origin string) bool {
		return true
	}
}
