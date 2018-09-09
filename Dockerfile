FROM alpine

# Based on Marco Ochse great Glutton / T-Pot template

# Setup apk
RUN apk -U --no-cache add \
                   build-base \
                   git \
                   go \
                   g++

# Setup go, medpot
RUN    export GOPATH=/opt/go/ && \
    mkdir /opt && \
    mkdir /opt/go && \
    cd /opt/go && \
    mkdir ./src && cd src && \
    git clone https://github.com/schmalle/medpot.git && \
    go get -d -v github.com/davecgh/go-spew/spew && \
    go get -d -v github.com/go-ini/ini && \
    go get -d -v github.com/mozillazg/request && \
    go get -d -v go.uber.org/zap && \
    cd /opt/go/src/medpot && \
    go build medpot && \
    cp ./medpot /usr/bin/medpot && \
    chmod 777 /var/log && \
    cp ./template/medpot.log > /var/log/medpot.log && \
    chmod 777 /var/log/medpot.log && \
    mkdir /data && \
    mkdir /data/medpot && \
    chmod 700 /data/medpot && \
    cp ./template/ews.xml /data/medpot/ && \
    cp ./template/dummyerror.xml /data/medpot/

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
