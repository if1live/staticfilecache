package staticfilecache

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

type Cache struct {
	basePath string
}

const headerExtension = ".header"

func (c *Cache) Get(key string) (resp []byte, ok bool) {
	seg, ok := ParseUrl(key)
	if !ok {
		return nil, false
	}
	bodyfilepath := seg.ToCacheFilePath(c.basePath)
	headerfilepath := seg.ToCacheFilePath(c.basePath) + headerExtension

	header, err := readFile(headerfilepath)
	if err != nil {
		return nil, false
	}
	body, err := readFile(bodyfilepath)
	if err != nil {
		return nil, false
	}

	response := joinHeaderBody(header, body)
	respBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}

	return respBytes, true
}

func (c *Cache) Set(key string, resp []byte) {
	seg, ok := ParseUrl(key)
	if !ok {
		return
	}
	cacheDir := seg.ToCacheDir(c.basePath)
	os.MkdirAll(cacheDir, 0755)

	// byte response -> real response
	b := bytes.NewBuffer(resp)
	req := &http.Request{}
	response, err := http.ReadResponse(bufio.NewReader(b), req)
	if err != nil {
		return
	}
	header, body := splitHeaderBody(response)

	headerfilepath := seg.ToCacheFilePath(c.basePath) + headerExtension
	writeFile(headerfilepath, header)

	bodyfilepath := seg.ToCacheFilePath(c.basePath)
	writeFile(bodyfilepath, body)
}

func (c *Cache) Delete(key string) {
	seg, ok := ParseUrl(key)
	if !ok {
		return
	}
	filepath := seg.ToCacheFilePath(c.basePath)
	os.Remove(filepath)
}

func splitHeaderBody(resp *http.Response) (header []byte, body []byte) {
	header, _ = httputil.DumpResponse(resp, false)

	var bodyBuf bytes.Buffer
	_, err := io.Copy(&bodyBuf, resp.Body)
	if err != nil {
		panic(err)
	}
	err = resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body = bodyBuf.Bytes()
	return
}

func joinHeaderBody(header, body []byte) *http.Response {
	data := append(header, body...)
	b := bytes.NewBuffer(data)
	req := &http.Request{}
	resp, _ := http.ReadResponse(bufio.NewReader(b), req)
	return resp
}

func writeFile(filepath string, data []byte) {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(data)
}

func readFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	data := make([]byte, fileInfo.Size())
	file.Read(data)
	return data, nil
}

func New(basePath string) *Cache {
	return &Cache{basePath}
}
