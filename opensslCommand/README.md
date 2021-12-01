> 该文件夹下的脚本根据 [openssl cookbook - Creating a Private Certification Authority](https://www.feistyduck.com/library/openssl-cookbook/online/ch-openssl.html#openssl-private-ca) 一节整理形成


#### 使用
##### 签发证书
该文件夹根据CA和密钥算法已经提前生成了对应的 CA 证书，可以使用脚本直接生成新的私钥并对证书进行签发：
- rsa 证书：
  ```shell
  ./genCert.sh -n example2 -a rsa
  ```
  生成的证书和私钥文件位于 rsa-example2 文件夹下
- ec 证书：
  ```shell
  ./genCert.sh -n example -a ec
  ```
##### 从 RootCA 开始
建议先删除该文件夹下的其他 CA 证书目录
- rsa
  ```shell
  ./genRootCA.sh -a rsa
  ./genSubCA.sh -a rsa
  ./genCert.sh -n example2 -a rsa
  ```
- ec
  ```shell
  ./genRootCA.sh -a ec
  ./genSubCA.sh -a ec
  ./genCert.sh -n example2 -a ec
  ```