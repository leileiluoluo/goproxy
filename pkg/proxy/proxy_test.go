package proxy

import "testing"

func TestPathRevision(t *testing.T) {
	paths := []struct {
		path  string
		after string
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
		rev := pathRevision(path.path)
		if path.after != rev {
			t.Errorf("want: %s, got: %s", path.after, rev)
		}
	}
}
