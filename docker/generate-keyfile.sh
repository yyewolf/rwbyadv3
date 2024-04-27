#!/bin/bash

openssl rand -base64 756 > mongo/dev-keyfile
chmod 400 mongo/dev-keyfile
chown 999:999 mongo/dev-keyfile
