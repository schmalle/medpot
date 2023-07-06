FROM alpine

# Based on Marco Ochse great Glutton / T-Pot template

# Setup apk
RUN apk -U --no-cache add \
                   build-base \
                   git \
                   go \
                   g++

# Setup go, medpot
RUN cd /tmp && \
    git clone https://github.com/s9rA16Bf4/medpot.git && \
    go get -d -v github.com/davecgh/go-spew/spew && \
    go get -d -v github.com/go-ini/ini && \
    go get -d -v github.com/mozillazg/request && \
    go get -d -v go.uber.org/zap && \
    cd /tmp/medpot && \
    go build -o medpot go/medpot.go go/logo.go && \
    mkdir -p /etc/medpot/ && \
	mkdir -p /var/log/medpot && \
	cp ./template/* /etc/medpot/ && \
	cp ./scripts/packet.txt /etc/medpot/ && \
	touch /var/log/medpot/medpot.log && \
	cp medpot /usr/bin/


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
WORKDIR /usr/bin/
USER medpot:medpot
CMD exec medpot
