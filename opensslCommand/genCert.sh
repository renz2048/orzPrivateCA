#!/bin/sh


while getopts "n:a:" opt
do
	case "$opt" in
	n) NAME=$OPTARG;;
	a) ALG=$OPTARG;;
	*) echo "Unknown option: $opt";;
	esac
done

if [ $# -ne 4 ]
then
    echo "Usage: $0 -n [cn] -a [rsa|ec]"
    exit 1;
fi

if [ "$ALG"x != "rsa"x || "$ALG"x != "ec"x ]; then
    echo " -a [rsa|ec]"
    exit 1;
fi


HOME=www.$ALG-$NAME.com
CNAME=www.$ALG-$NAME.com
CONF=$HOME/$CNAME.conf
CSR=$HOME/$CNAME.csr
CRT=$HOME/$CNAME.crt
KEY=$HOME/$CNAME.key
SUBCACONF=$ALG-sub-ca/$ALG-sub-ca.conf

mkdir -p $HOME


cat >> $CONF <<EOF
[req]
prompt = no
distinguished_name = dn
req_extensions = ext
input_password = 12345678

[dn]
CN = $CNAME
emailAddress = test@orzrz.com
O = orzrz
L = Beijing
C = CN

[ext]
subjectAltName = DNS:www.$NAME.com,DNS:$NAME.com
EOF


# 生成 key ：
if [ "$ALG"x = "rsa"x ]; then
openssl genpkey -out $KEY.secure -algorithm RSA -pkeyopt rsa_keygen_bits:2048
else
openssl genpkey -out $KEY.secure -algorithm EC -pkeyopt ec_paramgen_curve:P-256
fi
# 去除 key 的 pass phrase
openssl $ALG -in $KEY.secure -out $KEY
# 创建 CSR：
openssl req -new -config $CONF -key $KEY -out $CSR
# Sub CA 签发：
openssl ca -config $SUBCACONF -in $CSR -out $CRT -extensions server_ext

exit 0;
