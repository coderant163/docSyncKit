#!/bin/bash

# 本脚本废弃
# 推荐使用：./build/bin/docSyncKit rsa create

# 生成非对称密钥
TS=`date +"%Y%M%d_%H%m%S"`
KEY_DIR=../rsa_keys
PRIVATE_KEY=${KEY_DIR}/privatekey_${TS}.pem
PUBLIC_KEY=${KEY_DIR}/publickey_${TS}.pem

# 生成私钥：openssl genrsa -out rsaprivatekey.pem 4096
openssl genrsa -out ${PRIVATE_KEY} 4096

# 生成公钥：openssl rsa -in rsaprivatekey.pem -out rsapublickey.pem -pubout
openssl rsa -in ${PRIVATE_KEY} -out ${PUBLIC_KEY} -pubout


