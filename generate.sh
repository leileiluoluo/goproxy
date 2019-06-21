#!/bin/bash

mkdir internal

# copy dependencies
cp -r $GOROOT/src/cmd/go/internal/modfetch ./internal/
cp -r $GOROOT/src/cmd/go/internal/modfile ./internal/
cp -r $GOROOT/src/cmd/go/internal/modinfo ./internal/
cp -r $GOROOT/src/cmd/go/internal/base ./internal/
cp -r $GOROOT/src/cmd/go/internal/cache ./internal/
cp -r $GOROOT/src/cmd/go/internal/lockedfile ./internal/
cp -r $GOROOT/src/cmd/go/internal/module ./internal/
cp -r $GOROOT/src/cmd/go/internal/par ./internal/
cp -r $GOROOT/src/cmd/go/internal/renameio ./internal/
cp -r $GOROOT/src/cmd/go/internal/semver ./internal/
cp -r $GOROOT/src/cmd/go/internal/cfg ./internal/
cp -r $GOROOT/src/cmd/go/internal/str ./internal/
cp -r $GOROOT/src/cmd/go/internal/dirhash ./internal/
cp -r $GOROOT/src/cmd/go/internal/get ./internal/
cp -r $GOROOT/src/cmd/go/internal/web ./internal/
cp -r $GOROOT/src/cmd/go/internal/web2 ./internal/
cp -r $GOROOT/src/cmd/go/internal/load ./internal/
cp -r $GOROOT/src/cmd/go/internal/search ./internal/
cp -r $GOROOT/src/cmd/go/internal/work ./internal/

cp -r $GOROOT/src/cmd/internal/sys ./internal/
cp -r $GOROOT/src/cmd/internal/objabi ./internal/
cp -r $GOROOT/src/cmd/internal/buildid ./internal/
cp -r $GOROOT/src/cmd/internal/browser ./internal/
cp -r $GOROOT/src/internal/testenv ./internal/
cp -r $GOROOT/src/internal/singleflight ./internal/
cp -r $GOROOT/src/internal/xcoff ./internal/

# replace import paths
find . -type f -name "*.go" -exec sed -i '' 's#cmd/go/internal/#github.com/olzhy/goproxy/internal/#g' {} \; 
find . -type f -name "*.go" -exec sed -i '' 's#cmd/internal/#github.com/olzhy/goproxy/internal/#g' {} \; 
find . -type f -name "*.go" -exec sed -i '' 's#internal/testenv#github.com/olzhy/goproxy/internal/testenv#g' {} \; 
find . -type f -name "*.go" -exec sed -i '' 's#internal/singleflight#github.com/olzhy/goproxy/internal/singleflight#g' {} \; 
find . -type f -name "*.go" -exec sed -i '' 's#internal/xcoff#github.com/olzhy/goproxy/internal/xcoff#g' {} \; 
