```toml
title = "Ansible: 自动化管理服务器(2)"
date = "2017-01-12 21:07:04"
update_date = "2017-01-12 21:07:04"
author = "KDF5000"
thumb = ""
tags = ["Linux", "DevOps", "Ansible"]
draft = false
```
### Ansible 命令参数

| 参数             | 含义                                |
| -------------- | --------------------------------- |
| -v             | 详细模式，如果命令执行成功，输出详细的结果             |
| -i Path        | 指定host文件的路径，默认是/etc/ansible/hosts |
| -f Num         | 指定一个整数，默认是5，指定fork开启同步进程的个数       |
| -m Name        | 指定使用的module名称，默认是command          |
| -M path        | 指定module的目录。默认是/usr/share/ansible |
| -a Module_Args | 指定module的参数                       |
| -k             | 提示输入ssh密码，而不是基于ssh的密钥认证           |
| -sudo          | 指定使用sudo获得root权限                  |
| -K             | 提示输入sudo密码，与-sudo一起使用             |
| -u Username    | 指定被控机器的执行用户                       |
| -C             | 测试此命令会改变什么内容，不会真正的去执行             |

<!--more-->

### Ansible的常用模块

#### shell模块

使用shell来执行命令，如

```shell
$ ansible all -m shell -a "mkdir kongdefei"
```

将会在每一个客户机上创建kongdefei的目录

#### copy模块

实现主控机向客户机复制文件，类似scp的功能，比如复制主控机器的/etc/hosts到客户端机

```shell
$ansible all -m copy -a "src=/etc/hosts dest=/tmp/hosts"
```

#### file模块

该模块称之为文件属性模块，可以创建，删除，修改文件属性等操作

```shell
#创建文件
$ansible all -m file -a "dest=/tmp/a.txt state=touch"
#更改文件的用户和权限
$ansible all -m file -a "dest=/tmp/a.txt mode=600 owener=www-data group=root"
#创建目录
$ansible all -m file -a "dest=/tmp/kongdefei owner=kongdefei group=root state=directory"
#删除文件或者目录
$ansible all -m file -a "dest=/tmp/kongdefei state=absent"
```

**注：**state的其他选项：link(链接)、hard(硬链接)

#### stat模块

获取指定文件的状态信息，比如atime, ctime, mtime, md5, uid, gid等

```shell
$ ansible all -m stat -a "path=/home/kongdefei/tidb/tidb-latest-linux-amd64.tar.gz"
```

返回

```shell
172.16.1.96 | SUCCESS => {
    "changed": false, 
    "stat": {
        "atime": 1484225125.1479475, 
        "checksum": "b3afe00a995630184d98194c1a61d6f8fced1f05", 
        "ctime": 1484225125.5839477, 
        "dev": 2049, 
        "executable": false, 
        "exists": true, 
        "gid": 0, 
        "gr_name": "root", 
        "inode": 27787277, 
        "isblk": false, 
        "ischr": false, 
        "isdir": false, 
        "isfifo": false, 
        "isgid": false, 
        "islnk": false, 
        "isreg": true, 
        "issock": false, 
        "isuid": false, 
        "md5": "1de3012618bc0c1d1e42534d365a8ca5", 
        "mode": "0644", 
        "mtime": 1484225124.8959477, 
        "nlink": 1, 
        "path": "/home/kongdefei/tidb/tidb-latest-linux-amd64.tar.gz", 
        "pw_name": "root", 
        "readable": true, 
        "rgrp": true, 
        "roth": true, 
        "rusr": true, 
        "size": 51645042, 
        "uid": 0, 
        "wgrp": false, 
        "woth": false, 
        "writeable": true, 
        "wusr": true, 
        "xgrp": false, 
        "xoth": false, 
        "xusr": false
    }
...
```

#### 管理软件模块

指像yum, apt等管理包的软件

```shell
#安装nginx软件包
$ansible web -m yum -a "name=nginx state=present"
$ansible web -m apt -a "name=nginx state=present"

#安装包到一个特定的版本
$ansible web -m yum -a "name=nginx-1.6.2 state=present"
$ansible web -m apt -a "name=nginx-1.6.2 state=present"

#指定某个源仓库安装某软件包
$ansible web -m yum -a "name=php55w enablerepo=remi state=present"

#更新一个软件包是最新版本
$ansible web -m yum -a "name=nginx state=latest"
$ansible web -m apt -a "name=nginx state=latest"

#卸载一个软件
$ansible web -m yum -a "name=nginx state=absent"
$ansible web –m apt -a "name=nginx state=absent"
```

Ansible 支持很多操作系统的软件包管理，使用时 -m 指定相应的软件包管理工具模块，如果没有这样的模块，可以自己定义类似的模块或者使用 command 模块来安装软件包

#### User模块

用户模块主要用于创建，更改，删除用户

```shell
#创建一个用户
$ansible all -m user -a "name=hadoop password=123456"

#删除用户
ansible all -m user -a "name=hadoop state=absent"
```

#### service模块

主要是控制被控机的各种service。

```shell
#启动,重启，停止httpd服务
$ansible all -m service -a "name=httpd state=started(restarted|stopped)"
```
