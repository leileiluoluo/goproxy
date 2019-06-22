package proxy

import "testing"

func TestPathValidation(t *testing.T) {
	paths := map[string]bool{
		"rsc.io/@latest":          false,
		"rsc.io/quote":            false,
		"github.com/olzhy/quote":  false,
		"rsc.io/quote/@v/@latest": true,
	}
	for path, valid := range paths {
		err := pathValidation(path)
		if valid && nil != err {
			t.Errorf("path valid, but got a err, path: %s", path)
			return
		}
		if !valid && nil == err {
			t.Errorf("path invalid, but has no err, path: %s", path)
			return
		}
	}
}

func TestPathRevision(t *testing.T) {
	paths := []struct {
		before string
		after  string
	}{
		{
			"gopkg.in/%21data%21dog/dd-trace-go.v1/@v/v1.10.0.info",
			"gopkg.in/DataDog/dd-trace-go.v1/@v/v1.10.0.info",
		},
		{
			"github.com/%21burnt%21sushi/toml/@v/v0.3.1.zip",
			"github.com/BurntSushi/toml/@v/v0.3.1.zip",
		},
	}
	for _, path := range paths {
		rev := pathRevision(path.before)
		if path.after != rev {
			t.Errorf("want: %s, got: %s", path.after, rev)
		}
	}
}
