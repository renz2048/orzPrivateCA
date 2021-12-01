#!/bin/sh

if [ $# -ne 1 ]
then
    echo "Usage: $0 [cn]"
    exit 1;
fi

HOME=rsa-$1
CONF=$HOME/rsa-$1.conf
CSR=$HOME/rsa-$1.csr
CRT=$HOME/rsa-$1.crt
KEY=$HOME/rsa-$1.key
SUBCACONF=rsa-sub-ca/rsa-sub-ca.conf

mkdir -p $HOME

cat >> $CONF <<EOF
[req]
prompt = no
distinguished_name = dn
req_extensions = ext
input_password = 12345678

[dn]
CN = $1
emailAddress = test@chinatelecom.cn
O = Example
L = London
C = GB

[ext]
subjectAltName = DNS:www.example.com,DNS:example.com
EOF


# 生成 key ：
openssl genpkey -out $KEY.secure -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -aes-128-cbc
# 去除 key 的 pass phrase
openssl rsa -in $KEY.secure -out $KEY
# 创建 CSR：
openssl req -new -key $KEY -out $CSR
# Sub CA 签发：
openssl ca -config $SUBCACONF -in $CSR -out $CRT -extensions server_ext

exit 0;