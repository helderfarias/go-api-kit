package uri

import (
	"log"
	"net/url"
)

type URI struct {
	base  string
	paths *URIPath
}

type URIWithPaths struct {
	delegate *URI
}

func NewBuildURI(u string) *URI {
	return &URI{base: u, paths: NewPaths()}
}

func (u *URI) Path(p interface{}) *URI {
	u.paths.Path(p)
	return u
}

func (u *URI) QueryParam(key string, value interface{}) *URI {
	u.paths.Query(key, value)
	return u
}

func (u *URI) WithPaths(paths *URIPath) string {
	u.paths = paths
	return u.String()
}

func (u *URI) String() string {
	base, err := url.Parse(u.paths.ToBuild(u.base))
	if err != nil {
		log.Fatal(err)
	}
	return base.String()
}
