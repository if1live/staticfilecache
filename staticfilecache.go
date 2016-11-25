package staticfilecache

import "os"

type Cache struct {
	basePath string
}

func (c *Cache) Get(key string) (resp []byte, ok bool) {
	seg, ok := ParseUrl(key)
	if !ok {
		return nil, false
	}
	filepath := seg.ToCacheFilePath(c.basePath)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	// normal state
	fileInfo, _ := file.Stat()
	data := make([]byte, fileInfo.Size())
	file.Read(data)

	return data, true
}

func (c *Cache) Set(key string, resp []byte) {
	seg, ok := ParseUrl(key)
	if !ok {
		return
	}
	cacheDir := seg.ToCacheDir(c.basePath)
	os.MkdirAll(cacheDir, 0755)

	filepath := seg.ToCacheFilePath(c.basePath)
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(resp)
}

func (c *Cache) Delete(key string) {
	seg, ok := ParseUrl(key)
	if !ok {
		return
	}
	filepath := seg.ToCacheFilePath(c.basePath)
	os.Remove(filepath)
}

func New(basePath string) *Cache {
	return &Cache{basePath}
}
