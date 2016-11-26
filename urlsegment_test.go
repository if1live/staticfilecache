package staticfilecache

import (
	"reflect"
	"testing"
)

func TestParseUrl(t *testing.T) {
	cases := []struct {
		url string
		seg UrlSegment
	}{
		// starts with slash
		{"/", UrlSegment{"", []string{}, "index.html"}},
		{"/file.html", UrlSegment{"", []string{}, "file.html"}},

		{"/foo/bar", UrlSegment{"", []string{"foo"}, "bar"}},
		{"/foo/bar/", UrlSegment{"", []string{"foo", "bar"}, "index.html"}},

		// starts with host with double slash
		{"//a.com", UrlSegment{"a.com", []string{}, "index.html"}},
		{"//a.com/", UrlSegment{"a.com", []string{}, "index.html"}},

		{"//a.com/foo/bar", UrlSegment{"a.com", []string{"foo"}, "bar"}},
		{"//a.com/foo/bar/", UrlSegment{"a.com", []string{"foo", "bar"}, "index.html"}},

		// starts with schema
		{"http://a.com", UrlSegment{"a.com", []string{}, "index.html"}},
		{"http://a.com/", UrlSegment{"a.com", []string{}, "index.html"}},

		{"http://a.com/foo/bar", UrlSegment{"a.com", []string{"foo"}, "bar"}},
		{"http://a.com/foo/bar/", UrlSegment{"a.com", []string{"foo", "bar"}, "index.html"}},

		// single dot
		{"./foo", UrlSegment{"", []string{}, "foo"}},
		{"http://a.com/foo/./bar", UrlSegment{"a.com", []string{"foo"}, "bar"}},
		{"http://a.com/./foo/bar", UrlSegment{"a.com", []string{"foo"}, "bar"}},

		// double dot
		{"../foo", UrlSegment{"", []string{}, "foo"}},
		{"https://a.com/foo/../", UrlSegment{"a.com", []string{}, "index.html"}},
		{"https://a.com/foo/../../", UrlSegment{"a.com", []string{}, "index.html"}},

		// port num
		{"http://127.0.0.1:1234", UrlSegment{"127.0.0.1_1234", []string{}, "index.html"}},
	}
	for _, c := range cases {
		got, ok := ParseUrl(c.url)
		if !ok {
			t.Errorf("ParseUrl - cannot parer url %q", c.url)
		}
		if !reflect.DeepEqual(got, c.seg) {
			t.Errorf("ParseUrl, %s - expected %q, got %q", c.url, c.seg, got)
		}
	}
}
