package urlshortener

import (
	"github.com/go-chi/chi"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type URLShortener struct {
	urls map[string]string
	addr string
}

func NewShortener(addr string) *URLShortener {
	rand.Seed(time.Now().UnixNano())
	var m map[string]string
	m = make(map[string]string)
	return &URLShortener{
		urls: m,
		addr: addr,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *URLShortener) HandleSave(rw http.ResponseWriter, req *http.Request) {
	str := req.URL.Query().Get("u")
	_, err := url.Parse(str)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	h := randSeq(10)
	if _, emp := s.urls[h]; emp {
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	s.urls[h] = str
	_, _ = rw.Write([]byte(s.addr + "/" + h))
}

func (s *URLShortener) HandleExpand(rw http.ResponseWriter, req *http.Request) {
	str := chi.URLParam(req, "key")
	if s.urls[str] == "" {
		http.Error(rw, "", http.StatusNotFound)
	}

	http.Redirect(rw, req, s.urls[str], http.StatusMovedPermanently)
}
