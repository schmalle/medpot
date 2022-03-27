#!/bin/env bash

apk del --purge build-base git go g++
rm -rf /var/cache/apk/* /opt/go /root/dist
