package router

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type cachedResponseWriter struct {
	http.ResponseWriter
	buf []byte
}
type loggedCachedResponseWriter struct {
	http.ResponseWriter
	buf    []byte
	cached bool
}

func (w *cachedResponseWriter) Write(buf []byte) (int, error) {
	w.buf = buf
	return len(buf), nil
}

func (w *loggedCachedResponseWriter) Write(buf []byte) (int, error) {
	w.buf = buf
	return len(buf), nil
}

func Caching(c Cacher) Decorator {
	return func(s Server) Server {
		return ServerFunc(func(w http.ResponseWriter, r *http.Request) {
			v, err := c.Get(r.URL.String())
			if err == nil {
				w.Write(v)
				return
			}
			cw := &cachedResponseWriter{w, nil}
			s.ServeHTTP(cw, r)
			if len(cw.buf) > 0 {
				c.Add(r.URL.String(), cw.buf)
			}
			cw.ResponseWriter.Write(cw.buf)
		})
	}
}

func LoggedCaching(c Cacher) Decorator {
	return func(s Server) Server {
		return ServerFunc(func(w http.ResponseWriter, r *http.Request) {
			cw := &loggedCachedResponseWriter{w, nil, false}
			defer func(start time.Time) {
				msg := fmt.Sprintf(
					"reponse on request [%s] took %s (cached=%s)",
					r.URL.String(),
					time.Since(start),
					cw.cached,
				)
				log.Println(msg)
			}(time.Now())
			v, err := c.Get(r.URL.String())
			if err == nil {
				w.Write(v)
				cw.cached = true
				return
			}
			cw.cached = false
			s.ServeHTTP(cw, r)
			if len(cw.buf) > 0 {
				c.Add(r.URL.String(), cw.buf)
			}
			cw.ResponseWriter.Write(cw.buf)
		})
	}
}
