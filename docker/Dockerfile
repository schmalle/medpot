FROM alpine

# Based on Marco Ochse great Glutton / T-Pot template

# Setup apk
RUN apk -U --no-cache add build-base git go g++

RUN cd /tmp && git clone https://github.com/s9rA16Bf4/medpot.git && cd /tmp/medpot && / 
    git checkout unstable

# Setup go, medpot
RUN "/bin/env bash /tmp/medpot/docker/step_2.bash"

RUN "/bin/env bash /tmp/medpot/docker/step_3.bash"

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
