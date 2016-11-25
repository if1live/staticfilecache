package staticfilecache

import (
	"fmt"
	"net/url"
	"strings"
)

type UrlSegment struct {
	Domain   string
	Path     []string
	FileName string
}

func (seg *UrlSegment) ToCacheFilePath(root string) string {
	tokens := []string{root}
	tokens = append(tokens, seg.Domain)
	tokens = append(tokens, seg.Path...)
	tokens = append(tokens, seg.FileName)
	return joinTokens(tokens)
}

func (seg *UrlSegment) ToCacheDir(root string) string {
	tokens := []string{root}
	tokens = append(tokens, seg.Domain)
	tokens = append(tokens, seg.Path...)
	return joinTokens(tokens)
}

func joinTokens(tokens []string) string {
	path := strings.Join(tokens, "/")
	path = strings.Replace(path, "//", "/", -1)
	return path
}

func ParseUrl(rawurl string) (seg UrlSegment, ok bool) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return UrlSegment{}, false
	}

	scheme := u.Scheme
	if len(scheme) == 0 {
		scheme = "http"
	}

	host := u.Host
	baseurl := fmt.Sprintf("%s://%s", scheme, host)
	base, err := url.Parse(baseurl)
	if err != nil {
		return UrlSegment{}, false
	}
	u = base.ResolveReference(u)

	tokens := strings.Split(u.Path, "/")
	fileName := tokens[len(tokens)-1]
	if len(fileName) == 0 {
		fileName = "index.html"
	}
	if rawurl == host {
		fileName = "index.html"
	}

	if len(tokens) == 1 {
		ok = true
		seg = UrlSegment{
			host,
			[]string{},
			fileName,
		}
		return
	}

	ok = true
	seg = UrlSegment{
		host,
		tokens[1 : len(tokens)-1],
		fileName,
	}
	return
}
