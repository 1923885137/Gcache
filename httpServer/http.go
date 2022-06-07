package httpServer

import (
	"Gcache/cache"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	cache.Cache
}

func (S *Server) Listen() {
	http.Handle("/cache/", S.cacheHandler())
	http.Handle("/status/", S.statusHandler())
	http.ListenAndServe(":12345", nil)
}
func New(c cache.Cache) *Server {
	return &Server{c}
}
func (S *Server) cacheHandler() http.Handler {
	return &cacheHandler{S}
}
func (S *Server) statusHandler() http.Handler {
	return &cacheHandler{S}
}

type cacheHandler struct {
	*Server
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//EscapedPath(),返回的是一个转义前的url
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m := r.Method
	if m == http.MethodPut {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			e := h.Set(key, b)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}
	if m == http.MethodGet {
		b, e := h.Get(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(b) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write(b)
		return
	}
	if m == http.MethodDelete {
		e := h.Del(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
