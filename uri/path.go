package uri

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

var invalidCaracters = []string{"%0A", "%08", "%09", "%0B", "%0C", "%0D"}

type URIPath struct {
	paths  *strings.Builder
	values *url.Values
}

func NewPaths() *URIPath {
	b := &strings.Builder{}
	v := &url.Values{}
	return &URIPath{paths: b, values: v}
}

func (u *URIPath) Path(p interface{}) *URIPath {
	newPath := fmt.Sprintf("%v", p)

	if !strings.HasSuffix(u.paths.String(), "/") &&
		!strings.HasPrefix(newPath, "/") {
		newPath = "/" + newPath
	}

	u.paths.WriteString(newPath)
	return u
}

func (u *URIPath) Query(key string, value interface{}) *URIPath {
	safeQuery := fmt.Sprintf("%v", value)

	for _, safe := range invalidCaracters {
		safeQuery = strings.ReplaceAll(safeQuery, safe, "")
	}

	u.values.Add(key, fmt.Sprintf("%v", safeQuery))
	return u
}

func (u *URIPath) QueryParams(from url.Values) *URIPath {
	for k, v := range from {
		u.Query(k, v[0])
	}
	return u
}

func (u *URIPath) ToBuild(base string) string {
	paths, err := url.Parse(base + u.paths.String())
	if err != nil {
		log.Fatal(err)
	}

	paths.RawQuery = u.values.Encode()
	return paths.String()
}
