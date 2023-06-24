#!/usr/bin/python3

import caArgParser as cap
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography import x509
from cryptography.x509.oid import NameOID
from cryptography.hazmat.primitives import hashes
import datetime
import os

class certProducer():
    def __init__(self, common_name) -> None:
        self.key = None
        self.key_size = 2048
        self.key_pass_phrase = b"12345678"
        self.cert = None
        self.country_name = u"CN"
        self.state_or_province_name = u"Sichuan"
        self.locality_name = u"Chengdu"
        self.organization_name = u"ctyun elb"
        self.common_name = common_name
        self.file_name = common_name.replace(' ', '')
        self.validity_day = 365

        self.tempcerts_dir = 'tempcerts/'

        # 若不存在 tempcerts, 则创建; 若已存在 tempcerts, 则忽略.
        os.makedirs(self.tempcerts_dir, mode=622, exist_ok=True)

    def setKeySize(self, key_size) -> None:
        self.key_size = key_size

    def setKeyPassPhrase(self, key_pass_phrase) -> None:
        self.key_pass_phrase = key_pass_phrase

    def setCountryName(self, country_name) -> None:
        self.country_name = country_name

    def setStateOrProvinceName(self, state_or_province_name) -> None:
        self.state_or_province_name = state_or_province_name

    def setLocalityName(self, locality_name) -> None:
        self.locality_name = locality_name

    def setOrganizationName(self, organization_name) -> None:
        self.organization_name = organization_name

    def setCommonName(self, common_name) -> None:
        self.common_name = common_name

    def setValidityDay(self, validity_day) -> None:
        self.validity_day = validity_day

    def genPrivateKey(self) -> None:
        # 生成 2048 位 rsa 私钥
        self.key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=2048,
        )

        # 保存私钥到磁盘，并使用 key pass 进行加密
        with open(self.tempcerts_dir + self.file_name + ".key.secure", "wb") as f:
            f.write(self.key.private_bytes(
                encoding=serialization.Encoding.PEM,
                format=serialization.PrivateFormat.TraditionalOpenSSL,
                encryption_algorithm=serialization.BestAvailableEncryption(self.key_pass_phrase)
            ))

        # 保存私钥到磁盘, 不加密
        with open(self.tempcerts_dir + self.file_name + ".key.unsecure", "wb") as f:
            f.write(self.key.private_bytes(
                encoding=serialization.Encoding.PEM,
                format=serialization.PrivateFormat.TraditionalOpenSSL,
                encryption_algorithm=serialization.NoEncryption()
            ))

    def genSelfSignCert(self):
        # Various details about who we are. For a self-signed certificate the
        # subject and issuer are always the same.
        subject = issuer = x509.Name([
            x509.NameAttribute(NameOID.COUNTRY_NAME, self.country_name),
            x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, self.state_or_province_name),
            x509.NameAttribute(NameOID.LOCALITY_NAME, self.locality_name),
            x509.NameAttribute(NameOID.ORGANIZATION_NAME, self.organization_name),
            x509.NameAttribute(NameOID.COMMON_NAME, self.common_name),
        ])
        self.cert = x509.CertificateBuilder().subject_name(
            subject
        ).issuer_name(
            issuer
        ).public_key(
            self.key.public_key()
        ).serial_number(
            x509.random_serial_number()
        ).not_valid_before(
            datetime.datetime.utcnow()
        ).not_valid_after(
            # Our certificate will be valid for 10 years default
            datetime.datetime.utcnow() + datetime.timedelta(days=self.validity_day)
        ).add_extension(
            x509.SubjectAlternativeName([x509.DNSName(u"localhost")]),
            critical=False,
        # 使用root CA的私钥签发root CA证书
        ).sign(self.key, hashes.SHA256())

        # 若不存在 tempcerts, 则创建; 若已存在 tempcerts, 则忽略.
        os.makedirs(self.tempcerts_dir, mode=622, exist_ok=True)
        # 保存证书到磁盘
        with open(self.tempcerts_dir + self.file_name + ".pem", "wb") as f:
            f.write(self.cert.public_bytes(serialization.Encoding.PEM))

    def genCert(self, cert, key) -> None:
        # Various details about who we are. For a self-signed certificate the
        # subject and issuer are always the same.
        subject = x509.Name([
            x509.NameAttribute(NameOID.COUNTRY_NAME, self.country_name),
            x509.NameAttribute(NameOID.STATE_OR_PROVINCE_NAME, self.state_or_province_name),
            x509.NameAttribute(NameOID.LOCALITY_NAME, self.locality_name),
            x509.NameAttribute(NameOID.ORGANIZATION_NAME, self.organization_name),
            x509.NameAttribute(NameOID.COMMON_NAME, self.common_name),
        ])
        self.cert = x509.CertificateBuilder().subject_name(
            subject
        ).issuer_name(
            # issuer 指代签发者
            cert.subject
        ).public_key(
            self.key.public_key()
        ).serial_number(
            x509.random_serial_number()
        ).not_valid_before(
            datetime.datetime.utcnow()
        ).not_valid_after(
            # 默认证书有效期设置为 1 年
            datetime.datetime.utcnow() + datetime.timedelta(days=self.validity_day)
        ).add_extension(
            x509.SubjectAlternativeName([x509.DNSName(u"localhost")]),
            critical=False,
        # Sign our certificate with our private key
        ).sign(key, hashes.SHA256())

        # 若不存在 tempcerts, 则创建; 若已存在 tempcerts, 则忽略.
        os.makedirs(self.tempcerts_dir, mode=622, exist_ok=True)
        with open(self.tempcerts_dir + self.file_name + ".pem", "wb") as f:
            f.write(self.cert.public_bytes(serialization.Encoding.PEM))

class CAStore(cap.caArgParser):
    def __init__(self) -> None:
        self.rootca = certProducer("orz root CA")
        self.subca = certProducer("orz sub CA")
        self.user = certProducer('www.example.com')

    def createCA(self) -> None:
        self.rootca.setValidityDay(3650)
        self.rootca.genPrivateKey()
        self.rootca.genSelfSignCert()

        self.subca.setValidityDay(3650)
        self.subca.genPrivateKey()
        self.subca.genCert(self.rootca.cert, self.rootca.key)

        self.user.genPrivateKey()
        self.user.genCert(self.subca.cert, self.subca.key)
    def loadCA(self, cert_path, key_path) -> None:
        with open(cert_path, 'rb') as f:
            datas = f.read()
            self.subca.cert = x509.load_pem_x509_certificate(datas)
        with open(key_path, 'rb') as f:
            datas = f.read()
            self.subca.key = serialization.load_pem_private_key(datas,b'12345678')
    def createUserCert(self) -> None:
        self.user.setCommonName(self.opt.common_name)
        self.user.genPrivateKey()
        self.user.genCert(self.subca.cert, self.subca.key)

    def _ca_handler(self) -> None:
        self.createCA()

    def _user_handler(self) -> None:
        print(self.opt.common_name)
        self.loadCA('certs/orzsubCA.pem', 'certs/orzsubCA.key')
        self.createUserCert()

def main():
    store = CAStore()
    store.parseArgs()
    store.opt.func()

if __name__ == "__main__":
    main()