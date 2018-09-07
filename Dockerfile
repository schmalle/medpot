FROM alpine

# Based on Marco Ochse great Glutton / T-Pot template

# Setup apk
RUN apk -U --no-cache add \
                   build-base \
                   git \
                   go \
                   g++
RUN which go

# Setup go, medpot
RUN    export GOPATH=/opt/go/ && \
    mkdir /opt && \
    mkdir /opt/go && \
    cd /opt/go && \
    mkdir ./src && cd src && \
    git clone https://github.com/schmalle/medpot.git



RUN export GOPATH=/opt/go/ && go get -d -v github.com/davecgh/go-spew/spew
RUN export GOPATH=/opt/go/ && go get -d -v github.com/go-ini/ini
RUN export GOPATH=/opt/go/ && go get -d -v github.com/mozillazg/request

RUN export GOPATH=/opt/go/ && cd /opt/go/src/medpot && go build medpot && cp ./medpot /usr/bin

# Setup user, groups and configs
RUN    addgroup -g 2000 medpot && \
    adduser -S -s /bin/ash -u 2000 -D -g 2000 medpot && \
    mkdir -p /var/log/medpot

# Clean up
RUN    apk del --purge build-base \
                    git \
                    go \
                    g++ && \
    rm -rf /var/cache/apk/* \
           /opt/go \
           /root/dist

# Start medpot
WORKDIR /opt/go/src/medpot
USER medpot:medpot
CMD exec medpot
