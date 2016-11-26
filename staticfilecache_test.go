package staticfilecache

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"reflect"
	"testing"
)

var s struct {
	server *httptest.Server
}

func TestMain(m *testing.M) {
	flag.Parse()
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	mux := http.NewServeMux()
	s.server = httptest.NewServer(mux)

	mux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello world")
	}))
}

func teardown() {
	s.server.Close()
}

func TestStaticFileCache(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "staticfilecache")
	if err != nil {
		t.Fatalf("TempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cache := New(tempDir)

	key := s.server.URL + "/"
	client := &http.Client{}
	resp, _ := client.Get(key)
	respBytes, _ := httputil.DumpResponse(resp, true)

	_, ok := cache.Get(key)
	if ok {
		t.Fatal("retrieved key before adding it")
	}

	cache.Set(key, respBytes)

	retVal, ok := cache.Get(key)
	if !ok {
		t.Fatal("could not retrieve an element we just added")
	}
	if !bytes.Equal(retVal, respBytes) {
		t.Fatal("retrieved a different value than what we put in")
	}

	cache.Delete(key)

	_, ok = cache.Get(key)
	if ok {
		t.Fatal("deleted key still present")
	}
}

func TestSplitJoinResponse(t *testing.T) {
	key := s.server.URL + "/"

	client := &http.Client{}
	resp, _ := client.Get(key)

	header, body := splitHeaderBody(resp)
	resp2 := joinHeaderBody(header, body)

	if !reflect.DeepEqual(resp.Header, resp2.Header) {
		t.Errorf("compare header - expect %q, got %q", resp.Header, resp2.Header)
	}

	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)
	resp.Body.Close()

	var buf2 bytes.Buffer
	io.Copy(&buf2, resp.Body)
	resp2.Body.Close()

	if !bytes.Equal(buf.Bytes(), buf2.Bytes()) {
		t.Fatal("response body is not equals")
	}
}
