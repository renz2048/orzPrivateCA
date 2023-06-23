# 介绍
支持创建一个二级CA与用此二级CA签发的站点证书；
支持使用certs/目录下的二级CA创建一个测试用的站点证书
certs 目录下保存了一份二级CA（root CA: orzrootCA.pem、orzrootCA.key; sub CA: orzsubCA.pem、orzsubCA.key）和一份站点证书

# 参考
[openssl cookbook - Creating a Private Certification Authority](https://www.feistyduck.com/library/openssl-cookbook/online/ch-openssl.html#openssl-private-ca)
