package uri

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildUri(t *testing.T) {
	b := NewBuildURI("http://localhost")

	assert.Equal(t, "http://localhost", b.String())
}

func TestBuildFullUri(t *testing.T) {
	b := NewBuildURI("http://localhost")
	b.Path("/v1/api/")
	b.Path("1")
	b.Path("/")
	b.Path(50)
	b.QueryParam("q", "text tex text")

	assert.Equal(t, "http://localhost/v1/api/1/50?q=text+tex+text", b.String())
}

func TestBuildUriForOnlyPath(t *testing.T) {
	b := NewBuildURI("http://localhost")
	b.Path("/v1/api/teste")
	b.Path("1")
	b.Path(50)

	assert.Equal(t, "http://localhost/v1/api/teste/1/50", b.String())
}

func TestBuildUriForOnlyQuery(t *testing.T) {
	b := NewBuildURI("http://localhost")
	b.QueryParam("q1", "text tex text")
	b.QueryParam("q2", "hi")

	assert.Equal(t, "http://localhost?q1=text+tex+text&q2=hi", b.String())
}

func TestBuildUriFluentApi(t *testing.T) {
	b := NewBuildURI("http://localhost").
		Path("1").
		QueryParam("q1", "text tex text").
		QueryParam("q2", "hi")

	assert.Equal(t, "http://localhost/1?q1=text+tex+text&q2=hi", b.String())
}

func TestPathBuild(t *testing.T) {
	paths := NewPaths().
		Path("1").
		Path("v1").
		Path("svc").
		Query("q1", "text tex text").
		Query("q2", "hi")

	assert.Equal(t, "/1/v1/svc?q1=text+tex+text&q2=hi", paths.ToBuild(""))
}

func TestBuildUriWithPaths(t *testing.T) {
	values := url.Values{}
	values.Add("n1", "1")
	values.Add("n2", "2")

	paths := NewPaths().
		Path("1").
		Path("v1").
		Path("svc").
		Query("q1", "text tex text").
		Query("q2", "hi").
		Query("q3", "1%0A 2%08 3%09 4%0B 5%0C 5%0D").
		QueryParams(values)

	result := NewBuildURI("http://localhost").WithPaths(paths)

	assert.Equal(t, "http://localhost/1/v1/svc?n1=1&n2=2&q1=text+tex+text&q2=hi&q3=1+2+3+4+5+5", result)
}

func TestBuildUriCustom(t *testing.T) {
	r := service().Path("v1").Path("1").Query("n", "1").ToBuild("")

	assert.Equal(t, "/v1/1?n=1", r)
}

func service() *URIPath {
	return NewPaths()
}
