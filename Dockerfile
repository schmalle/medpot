FROM alpine

# Based on Marco Ochse great Glutton / T-Pot template

# Setup apk
RUN apk -U --no-cache add \
                   build-base \
                   git \
                   go \
                   g++

RUN wget -O dep.zip https://github.com/golang/dep/releases/download/v0.3.0/dep-linux-amd64.zip && \
    echo '96c191251164b1404332793fb7d1e5d8de2641706b128bf8d65772363758f364  dep.zip' | sha256sum -c - && \
    unzip -d /usr/bin dep.zip && rm dep.zip


RUN dep ensure -vendor-only

# Setup go, medpot
RUN    export GOPATH=/opt/go/ && \
    mkdir /opt && \
    mkdir /opt/go && \
    cd /opt/go && \
    mkdir ./src && cd src && \
    git clone https://github.com/schmalle/medpot.git && \
    cd medpot

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
WORKDIR /opt/go/medpot
USER medpot:medpot
CMD dep ensure
# CMD exec bin/server -i $(/sbin/ip address | grep '^2: ' | awk '{ print $2 }' | tr -d [:punct:]) -l /var/log/glutton/glutton.log
