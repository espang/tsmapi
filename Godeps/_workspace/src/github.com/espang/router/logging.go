package router

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Logging() Decorator {
	return func(s Server) Server {
		return ServerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func(start time.Time) {
				msg := fmt.Sprintf("request [%s] took %s", r.URL.String(), time.Since(start))
				log.Println(msg)
			}(time.Now())
			s.ServeHTTP(w, r)
		})
	}
}
