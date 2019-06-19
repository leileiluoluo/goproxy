package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/olzhy/goproxy/internal/modfetch"
	"github.com/olzhy/goproxy/internal/modfetch/codehost"
	"github.com/olzhy/goproxy/internal/module"
)

const (
	List   = "/@v/list"
	Latest = "/@latest"
	Info   = ".info"
	Mod    = ".mod"
	Zip    = ".zip"
)

func init() {
	codehost.WorkRoot = "/tmp"
	modfetch.PkgMod = "/tmp"
}

func Proxy() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimLeft(r.RequestURI, "/")
		if len(path) <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("request path: %s", path)
		switch {
		// suffix is /@v/list
		case strings.HasSuffix(path, List):
			mod := strings.TrimSuffix(path, List)
			vers := lookupVersions(mod)
			fmt.Fprintln(w, strings.Join(vers, "\n"))
			return
		// suffix is /@latest
		case strings.HasSuffix(path, Latest):
			mod := strings.TrimSuffix(path, Latest)
			rev := lookupLatestRev(mod)
			j, _ := json.Marshal(rev)
			fmt.Fprintln(w, string(j))
			return
		// suffix is .info
		case strings.HasSuffix(path, Info):
			verIndex := strings.LastIndex(path, "/@v/")
			mod := path[:verIndex]
			ver := path[verIndex+len("/@v/") : len(path)-len(Info)]
			rev := loadRev(mod, ver)
			j, _ := json.Marshal(rev)
			fmt.Fprintln(w, string(j))
			return
		// suffix is .mod
		case strings.HasSuffix(path, Mod):
			verIndex := strings.LastIndex(path, "/@v/")
			mod := path[:verIndex]
			ver := path[verIndex+len("/@v/") : len(path)-len(Mod)]
			rev := loadModContent(mod, ver)
			fmt.Fprintln(w, string(rev))
			return
			// suffix is .mod
		case strings.HasSuffix(path, Zip):
			verIndex := strings.LastIndex(path, "/@v/")
			mod := path[:verIndex]
			ver := path[verIndex+len("/@v/") : len(path)-len(Zip)]
			zipFile := loadZip(mod, ver)
			if len(zipFile) <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			http.ServeFile(w, r, zipFile)
			return
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

func lookupVersions(mod string) []string {
	log.Printf("lookup module %s", mod)
	repo, err := modfetch.Lookup(mod)
	if nil != err {
		log.Printf("lookup module error, err: %s", err)
		return []string{}
	}
	vers, _ := repo.Versions("")
	return vers
}

func lookupLatestRev(mod string) *modfetch.RevInfo {
	log.Printf("lookup module %s", mod)
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
	modContext, err := modfetch.GoMod(mod, rev)
	if nil != err {
		log.Printf("fetch mod content error, err: %s", err)
		return nil
	}
	return modContext
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
