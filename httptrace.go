package httptrace

import (
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"
)

var once = sync.Once{}
var idCh = make(chan uint64)

type flusher interface {
	Flush() error
}

type interceptor struct {
	http.ResponseWriter
	id  uint64
	to  string
	dst io.Writer
}

func (i interceptor) WriteHeader(code int) {
	dst := i.dst
	io.WriteString(dst, "--- response "+strconv.FormatUint(i.id, 10)+" to "+i.to+"---\n")
	io.WriteString(dst, strconv.Itoa(code)+" "+http.StatusText(code)+"\n")
	i.Header().Write(dst)

	if f, ok := dst.(flusher); ok {
		f.Flush()
	}
	i.ResponseWriter.WriteHeader(code)
}

func Wrap(h http.Handler, dst io.Writer) http.Handler {
	once.Do(func() {
		go func() {
			var i uint64 = 1
			for ; ; i++ {
				idCh <- i
			}
		}()
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := <-idCh

		io.WriteString(dst, "--- request "+strconv.FormatUint(id, 10)+" from "+r.RemoteAddr+"---\n")
		dump, err := httputil.DumpRequest(r, false)
		if err == nil {
			dst.Write(dump)
			if f, ok := dst.(flusher); ok {
				f.Flush()
			}
		}

		h.ServeHTTP(interceptor{w, id, r.RemoteAddr, dst}, r)
	})
}