#!/bin/sh


while getopts "a:" opt
do
	case "$opt" in
	n) CNAME=$OPTARG;;
	a) ALG=$OPTARG;;
	*) echo "Unknown option: $opt";;
	esac
done

if [ $# -ne 2 ]
then
    echo "Usage: $0 -a [rsa|ec]"
    exit 1;
fi


HOME=$ALG-sub-ca
CERT_PATH=$HOME/certs
DB_PATH=$HOME/db
PRIVATE_PATH=$HOME/private

CONF=$HOME/$ALG-sub-ca.conf

CSR=$CERT_PATH/$ALG-sub-ca.csr
CRT=$CERT_PATH/$ALG-sub-ca.crt
KEY=$PRIVATE_PATH/$ALG-sub-ca.key


mkdir -p $CERT_PATH $DB_PATH $PRIVATE_PATH
chmod 700 $PRIVATE_PATH
touch $DB_PATH/index
openssl rand -hex 16 > $DB_PATH/serial
echo 1001 > $DB_PATH/crlnumber

cat >> $CONF <<EOF
[default]
name                    = $ALG-sub-ca
domain_suffix           = orzrz.com
aia_url                 = http://\$name.\$domain_suffix/\$name.crt
crl_url                 = http://\$name.\$domain_suffix/\$name.crl
ocsp_url                = http://ocsp.\$name.\$domain_suffix:9081
default_ca              = ca_default
name_opt                = utf8,esc_ctrl,multiline,lname,align

[ca_dn]
countryName             = "CN"
organizationName        = "orzrz"
commonName              = "Sub CA"


[ca_default]
home                    = $ALG-root-ca
database                = \$home/db/index
serial                  = \$home/db/serial
crlnumber               = \$home/db/crlnumber
certificate             = \$home/certs/\$home.crt
private_key             = \$home/private/\$home.key
RANDFILE                = \$home/private/random
new_certs_dir           = \$home/certs
unique_subject          = no
copy_extensions         = none
default_days            = 3650
default_crl_days        = 730
default_md              = sha256
policy                  = policy_c_o_match

[policy_c_o_match]
countryName             = match
stateOrProvinceName     = optional
organizationName        = match
organizationalUnitName  = optional
commonName              = supplied
emailAddress            = optional


[req]
default_bits            = 4096
encrypt_key             = yes
default_md              = sha256
utf8                    = yes
string_mask             = utf8only
prompt                  = no
distinguished_name      = ca_dn
req_extensions          = ca_ext

[ca_ext]
basicConstraints        = critical,CA:true
keyUsage                = critical,keyCertSign,cRLSign
subjectKeyIdentifier    = hash


[sub_ca_ext]
authorityInfoAccess     = @issuer_info
authorityKeyIdentifier  = keyid:always
basicConstraints        = critical,CA:true,pathlen:0
crlDistributionPoints   = @crl_info
extendedKeyUsage        = clientAuth,serverAuth
keyUsage                = critical,keyCertSign,cRLSign
nameConstraints         = @name_constraints
subjectKeyIdentifier    = hash

[crl_info]
URI.0                   = \$crl_url

[issuer_info]
caIssuers;URI.0         = \$aia_url
OCSP;URI.0              = \$ocsp_url

[name_constraints]
permitted;DNS.0=orzrz.com
permitted;DNS.1=orzrz.org
excluded;IP.0=0.0.0.0/0.0.0.0
excluded;IP.1=0:0:0:0:0:0:0:0/0:0:0:0:0:0:0:0


[ocsp_ext]
authorityKeyIdentifier  = keyid:always
basicConstraints        = critical,CA:false
extendedKeyUsage        = OCSPSigning
noCheck                 = yes
keyUsage                = critical,digitalSignature
subjectKeyIdentifier    = hash

[server_ext]
authorityInfoAccess     = @issuer_info
authorityKeyIdentifier  = keyid:always
basicConstraints        = critical,CA:false
crlDistributionPoints   = @crl_info
extendedKeyUsage        = clientAuth,serverAuth
keyUsage                = critical,digitalSignature,keyEncipherment
subjectKeyIdentifier    = hash

[client_ext]
authorityInfoAccess     = @issuer_info
authorityKeyIdentifier  = keyid:always
basicConstraints        = critical,CA:false
crlDistributionPoints   = @crl_info
extendedKeyUsage        = clientAuth
keyUsage                = critical,digitalSignature
subjectKeyIdentifier    = hash
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

# RootCA 签发：
openssl ca -config $CONF -in $CSR -out $CRT -extensions sub_ca_ext