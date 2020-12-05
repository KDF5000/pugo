#!/bin/bash

set -e

QINIU_LOCAL_PATH="/Users/kongdefei/Documents/Personal/qiuniu-personal-blog"
# 从七牛云存储同步文件到本地
for key in `cat $1`
do
    path=`dirname $key`
    if [ ! -d "$path" ];then
        mkdir -p ./$path
    fi
    cp ${QINIU_LOCAL_PATH}$key ./$key
done
