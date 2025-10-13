#!/bin/bash
export CGO_CFLAGS="-I/opt/homebrew/opt/leptonica/include -I/opt/homebrew/opt/tesseract/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/leptonica/lib -L/opt/homebrew/opt/tesseract/lib"
export PKG_CONFIG_PATH="/opt/homebrew/opt/leptonica/lib/pkgconfig:/opt/homebrew/opt/tesseract/lib/pkgconfig"
export DYLD_LIBRARY_PATH="/opt/homebrew/opt/leptonica/lib:/opt/homebrew/opt/tesseract/lib"
$(go env GOPATH)/bin/air
