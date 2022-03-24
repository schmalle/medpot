FROM alpine

# Based on Marco Ochse great Glutton / T-Pot template

# Setup apk
RUN /bin/env bash ./docker/step_1.bash

# Setup go, medpot
RUN /bin/env bash ./docker/step_2.bash

RUN /bin/env bash ./docker/step_3.bash

# Setup user, groups and configs
RUN /bin/env bash ./docker/step_4.bash

# Clean up
RUN /bin/env bash ./docker/step_5.bash

# Start medpot
WORKDIR /opt/go/src/medpot
USER medpot:medpot
CMD exec medpot
