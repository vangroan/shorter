#!/bin/bash

sudo docker run -it \
  --name shorter \
  --publish 4003:8000 \
  --network nginx-proxy \
  -e SHORTER_BASEURL='https://u.vangroan.com/' \
  -e VIRTUAL_HOST='u.vangroan.com' \
  shorter:latest

