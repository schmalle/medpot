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
    bash scripts/compile_medpot.sh && \
    cp medpot /usr/bin

# Setup user, groups and configs
RUN "/bin/env bash /tmp/medpot/docker/step_4.bash"

# Clean up
RUN apk del --purge build-base git go g++ && \
    rm -rf /var/cache/apk/* /opt/go /root/dist && \
    rm -rf /tmp/medpot


# Start medpot
WORKDIR /opt/go/src/medpot
USER medpot:medpot
CMD exec medpot
