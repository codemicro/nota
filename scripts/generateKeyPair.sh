#!/usr/bin/env bash
mkdir -p ../run
ssh-keygen -t rsa -b 4096 -m PEM -f ../run/private.key
openssl rsa -in ../run/private.key -pubout -outform PEM -out ../run/public.key