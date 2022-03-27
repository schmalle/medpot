FROM alpine

# Based on Marco Ochse great Glutton / T-Pot template

# Setup apk
RUN apk -U --no-cache add build-base git go g++

# Setup go, medpot
RUN cd /tmp && \
    git clone https://github.com/s9rA16Bf4/medpot.git && \
    cd /tmp/medpot && \
    git checkout unstable && \
    go get -v github.com/davecgh/go-spew/spew && \
    go get -v github.com/go-ini/ini && \
    go get -v github.com/mozillazg/request && \
    go get -v go.uber.org/zap && \
    cd ../ && \
    go build -o medpot go/*.go && \
    cp medpot /usr/bin/ && \
    mkdir /etc/medpot/ && \
    cp template/* /etc/medpot/ && \
    mkdir /var/log/medpot && \
    touch /var/log/medpot/medpot.log

# Setup user, groups and configs
RUN addgroup -g 2000 medpot && \
    adduser -S -s /bin/bash -u 2000 -D -g 2000 medpot

# Clean up
RUN apk del --purge build-base git go g++ && \
    rm -rf /var/cache/apk/* /opt/go /root/dist && \
    rm -rf /tmp/medpot


# Start medpot
WORKDIR /opt/go/src/medpot
USER medpot:medpot
CMD exec medpot
