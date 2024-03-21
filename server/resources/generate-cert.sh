#!/bin/bash

openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -subj '/CN=h41stur.com/O=H41stur/C=US' -new -x509 -sha256 -key server.key -out server.crt -days 3650