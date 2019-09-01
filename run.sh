#!/bin/bash

sudo docker build . -t shorter:latest

sudo docker run -d \
  --rm \
  --name shorter \
  --publish 4003:8000 \
  --network nginx-proxy \
  -e SHORTER_BASEURL='https://u.vangroan.com/' \
  -e VIRTUAL_HOST='u.vangroan.com' \
  -e LETSENCRYPT_HOST='u.vangroan.com' \
  -e LETSENCRYPT_EMAIL='wimpievictor+letsencrypt@gmail.com' \
  shorter:latest

