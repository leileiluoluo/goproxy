package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/olzhy/goproxy/internal/modfetch"
	"github.com/olzhy/goproxy/internal/modfetch/codehost"
	"github.com/olzhy/goproxy/internal/module"
)

const (
	ListSuffix   = "/@v/list"
	LatestSuffix = "/@latest"
	InfoSuffix   = ".info"
	ModSuffix    = ".mod"
	ZipSuffix    = ".zip"
	VInfix       = "/@v/"

	UpperCaseSign = "%21"
)

func init() {
	goPath := os.Getenv("GOPATH")
	if "" == goPath {
		goPath = "/tmp"
	}
	// mod dir is $GOPATH/pkg/mod
	modfetch.PkgMod = filepath.Join(goPath, "pkg", "mod")
	// work dir is $GOPATH/pkg/mod/cache/vcs
	codehost.WorkRoot = filepath.Join(modfetch.PkgMod, "cache", "vcs")
}

func Proxy() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.RequestURI, "/")
		log.Printf("%s req path: %s", r.RemoteAddr, path)

		// req path validation
		if err := pathValidation(path); nil != err {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		// path revision
		// such as gopkg.in/%21data%21dog/dd-trace-go.v1/@v/v1.10.0.info
		// revise to gopkg.in/DataDog/dd-trace-go.v1/@v/v1.10.0.info
		path = pathRevision(path)

		switch {
		// suffix is /@v/list
		case strings.HasSuffix(path, ListSuffix):
			mod := strings.TrimSuffix(path, ListSuffix)
			vers := lookupVersions(mod)
			fmt.Fprintln(w, strings.Join(vers, VInfix))
			return
		// suffix is /@latest
		case strings.HasSuffix(path, LatestSuffix):
			mod := strings.TrimSuffix(path, LatestSuffix)
			rev := lookupLatestRev(mod)
			j, _ := json.Marshal(rev)
			fmt.Fprintln(w, string(j))
			return
		// suffix is .info
		case strings.HasSuffix(path, InfoSuffix):
			mod, ver := parseModAndVersion(path, VInfix, InfoSuffix)
			rev := loadRev(mod, ver)
			j, _ := json.Marshal(rev)
			fmt.Fprintln(w, string(j))
			return
		// suffix is .mod
		case strings.HasSuffix(path, ModSuffix):
			mod, ver := parseModAndVersion(path, VInfix, ModSuffix)
			rev := loadModContent(mod, ver)
			fmt.Fprintln(w, string(rev))
			return
		// suffix is .zip
		case strings.HasSuffix(path, ZipSuffix):
			mod, ver := parseModAndVersion(path, VInfix, ZipSuffix)
			zipFile := loadZip(mod, ver)
			if "" == zipFile {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			http.ServeFile(w, r, zipFile)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "please give me a correct module query")
		}
	})
}

func pathValidation(path string) error {
	msg := `give me a correct module query,
	such as github.com/olzhy/quote/@latest`
	if "" == path {
		return errors.New("empty path")
	}

	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		return errors.New(msg)
	}

	suffixes := []string{ListSuffix, LatestSuffix,
		InfoSuffix, ModSuffix, ZipSuffix}
	ok := false
	for _, suffix := range suffixes {
		if strings.HasSuffix(path, suffix) {
			ok = true
			break
		}
	}
	if !ok {
		log.Printf("suffix invalid, path: %s", path)
		return errors.New("suffix should match /@v/list|/@latest|.info|.zip")
	}
	return nil
}

func pathRevision(path string) string {
	after := ""
	i := 0
	for i < len(path) {
		c := path[i]
		if '%' == c {
			if i+3 < len(path) &&
				UpperCaseSign == path[i:i+3] {
				after += strings.ToUpper(string(path[i+3]))
				i += 4
				continue
			}
		}
		after += string(c)
		i++
	}
	return after
}

func parseModAndVersion(path, infix, suffix string) (mod string, ver string) {
	verIndex := strings.LastIndex(path, infix)
	mod = path[:verIndex]
	ver = path[verIndex+len(infix) : len(path)-len(suffix)]
	return
}

func lookupVersions(mod string) []string {
	log.Printf("lookup latest versions, mod: %s", mod)
	repo, err := modfetch.Lookup(mod)
	if nil != err {
		log.Printf("lookup module error, err: %s", err)
		return []string{}
	}
	vers, _ := repo.Versions("")
	return vers
}

func lookupLatestRev(mod string) *modfetch.RevInfo {
	log.Printf("lookup latest rev, mod: %s", mod)
	repo, err := modfetch.Lookup(mod)
	if nil != err {
		log.Printf("lookup module error, err: %s", err)
		return nil
	}
	rev, _ := repo.Latest()
	return rev
}

func loadRev(mod, rev string) *modfetch.RevInfo {
	log.Printf("load rev, mod: %s, rev: %s", mod, rev)
	revInfo, err := modfetch.Stat(mod, rev)
	if nil != err {
		log.Printf("fetch stat error, err: %s", err)
		return nil
	}
	return revInfo
}

func loadModContent(mod, rev string) []byte {
	log.Printf("load mod content, mod: %s, rev: %s", mod, rev)
	modContent, err := modfetch.GoMod(mod, rev)
	if nil != err {
		log.Printf("fetch mod content error, err: %s", err)
		return nil
	}
	return modContent
}

func loadZip(mod, rev string) string {
	log.Printf("load zip file, mod: %s, rev: %s", mod, rev)
	zipfile, err := modfetch.DownloadZip(module.Version{Path: mod, Version: rev})
	if nil != err {
		log.Printf("fetch zip file error, err: %s", err)
		return ""
	}
	log.Printf("fetch zip file success, file: %s", zipfile)
	return zipfile
}
