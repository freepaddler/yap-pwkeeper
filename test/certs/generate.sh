#!/bin/sh

# we don't need libreSSL
for e in $(which -a openssl); do
    if "$e" version | grep -qi '^openssl'; then
        alias openssl='$e'
        break
    fi
done

# root CA
openssl req \
    -subj "/CN=root CA" -out ca.crt \
    -newkey rsa:4096 -nodes -keyout ca.key\
    -x509 -days 365 \
    -addext basicConstraints=critical,CA:TRUE \
    -addext subjectKeyIdentifier=hash \
    -addext authorityKeyIdentifier=keyid:always \
    -addext keyUsage=critical,digitalSignature,cRLSign,keyCertSign

# intermediate CA
openssl req \
    -subj "/CN=intermediate CA" -out ca-int.crt \
    -newkey rsa:4096 -nodes -keyout ca-int.key \
    -x509 -days 365 -CA ca.crt -CAkey ca.key \
    -addext basicConstraints=critical,CA:TRUE,pathlen:0 \
    -addext subjectKeyIdentifier=hash \
    -addext authorityKeyIdentifier=keyid:always \
    -addext keyUsage=critical,digitalSignature,cRLSign,keyCertSign

# server certificate
openssl req \
    -subj "/CN=server certificate" -out server.crt \
    -newkey rsa:4096 -nodes -keyout server.key \
    -x509 -days 180 -CA ca-int.crt -CAkey ca-int.key \
    -addext subjectAltName=DNS:localhost,IP:127.0.0.1 \
    -addext basicConstraints=critical,CA:FALSE \
    -addext subjectKeyIdentifier=hash \
    -addext authorityKeyIdentifier=keyid:always \
    -addext keyUsage=critical,nonRepudiation,digitalSignature,keyEncipherment \
    -addext extendedKeyUsage=serverAuth

# client certificate
openssl req \
    -subj "/CN=client certificate" -out client.crt \
    -newkey rsa:4096 -nodes -keyout client.key \
    -x509 -days 180 -CA ca-int.crt -CAkey ca-int.key \
    -addext basicConstraints=critical,CA:FALSE \
    -addext subjectKeyIdentifier=hash \
    -addext authorityKeyIdentifier=keyid:always \
    -addext keyUsage=critical,nonRepudiation,digitalSignature,keyEncipherment \
    -addext extendedKeyUsage=clientAuth,emailProtection
