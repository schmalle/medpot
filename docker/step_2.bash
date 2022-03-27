#!/bin/env bash

export GOPATH=/opt/go/

# Install all the requirements
go get -d -v github.com/davecgh/go-spew/spew
go get -d -v github.com/go-ini/ini
go get -d -v github.com/mozillazg/request
go get -d -v go.uber.org/zap


bash ../scripts/compile_medpot.sh
go build medpot
cp ./medpot /usr/bin/medpot
