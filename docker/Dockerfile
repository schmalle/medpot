FROM alpine

# Based on Marco Ochse great Glutton / T-Pot template

# Setup apk
RUN bash step_1.bash

# Setup go, medpot
RUN bash step_2.bash

RUN bash step_3.bash

# Setup user, groups and configs
RUN bash step_4.bash

# Clean up
RUN bash step_5.bash

# Start medpot
WORKDIR /opt/go/src/medpot
USER medpot:medpot
CMD exec medpot
