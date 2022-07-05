#!/bin/bash
openssl genrsa -out server.pem 2048

# 本地测试 Common Name 填写 127.0.0.1
openssl req -new -x509 -key server.pem -out public.crt -days 99999
