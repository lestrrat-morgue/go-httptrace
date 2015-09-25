package httptrace

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrace(t *testing.T) {
	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hello, World!")
	})

	buf := &bytes.Buffer{}
	s := httptest.NewServer(Wrap(app, buf))
	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		t.Logf("Error while requesting %s: %s", s.URL, err)
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Logf("Response did not succeed")
		return
	}

	t.Logf("buf = %s", buf.String())
}