#!/bin/env bash

addgroup -g 2000 medpot
adduser -S -s /bin/bash -u 2000 -D -g 2000 medpot
mkdir -p /var/log/medpot
