```toml
title = "Ansible: 自动化管理服务器(1)"
date = "2017-01-12 00:12:33"
update_date = "2017-01-12 00:12:33"
author = "KDF5000"
thumb = ""
tags = ["Linux", "DevOps", "Ansible"]
draft = false
```
## Ansible配置

Ansible 提供一种最简单的方式用于发布、管理和编排计算机系统的工具，你可在数分钟内搞定

Ansible 是一个模型驱动的配置管理器，支持多节点发布、远程任务执行。默认使用 SSH 进行远程连接。无需在被管理节点上安装附加软件，可使用各种编程语言进行扩展

![Ansible](http://7xtgvp.com2.z0.glb.clouddn.com/blog/images/Ansible.png)

<!--more-->

### 安装

#### 安装节点

| 节点 | ip          | 操作系统     | 角色            |
| ---- | ----------- | ------------ | --------------- |
| gd67 | 172.16.1.67 | Ubuntu 14.04 | Controller Node |
| gd86 | 172.16.1.86 | Ubuntu 14.04 |                 |
| gd87 | 172.16.1.87 | Ubuntu 14.04 |                 |
| gd88 | 172.16.1.88 | Ubuntu 14.04 |                 |
| gd89 | 172.16.1.89 | Ubuntu 14.04 |                 |
| gd92 | 172.16.1.92 | Ubuntu 14.04 |                 |
| gd93 | 172.16.1.93 | Ubuntu 14.04 |                 |
| gd94 | 172.16.1.94 | Ubuntu 14.04 |                 |
| gd95 | 172.16.1.95 | Ubuntu 14.04 |                 |
| gd96 | 172.16.1.96 | Ubuntu 14.04 |                 |

#### 安装Ansible

```shell
$ sudo apt-get install software-properties-common
$ sudo apt-add-repository ppa:ansible/ansible
$ sudo apt-get update
$ sudo apt-get install ansible
```

#### 免密登录

设置ControllerNode能够免密登录被管理的远程主机，使用系统自带的`ssh-copy-id`复制ssh的公钥到被控主机，下面是自动运行的脚本

```shell
#!/usr/bin/env bash
hosts=(
"172.16.1.86"
"172.16.1.87"
"172.16.1.88"
"172.16.1.89"
"172.16.1.92"
"172.16.1.93"
"172.16.1.94"
"172.16.1.95"
"172.16.1.96"
)
for host in ${hosts[@]};
do
    expect <<EOF
    spawn ./ssh-copy-id root@$host
    expect {
        "(yes/no)?" { send "yes\n"; exp_continue  }
        "password:" { send "123456\n"  }
    }
    expect eof
 EOF
done
```
脚本里的被控主机列表，用户名和密码根据自己的情况修改

### 使用

现在就可以开始你的第一条命令了，不过首先需要将所有被控主机添加到ansible的host里，即/etc/ansible/hosts中

```
172.16.1.86
172.16.1.87
172.16.1.88
172.16.1.89
172.16.1.92
172.16.1.93
172.16.1.94
172.16.1.95
```

现在就可以执行第一条命令了

```shell
$ ansible all -m ping
```

Ansible会像SSH那样试图用你的当前用户名来连接你的远程机器.要覆写远程用户名,只需使用’-u’参数. 如果你想访问 sudo模式,这里也有标识(flags)来实现:

```shell
# as devin
$ ansible all -m ping -u root
# as devin, sudoing to root
$ ansible all -m ping -u root --sudo
# as devin, sudoing to batman
$ ansible all -m ping -u root --sudo --sudo-user kdf
```

**注意：**上面免密登录的步骤使用的用户名要和现在一样

如果得到下面的结果，那么恭喜你已经安装配置成功了！

```shell
root@gd67:/home/kongdefei/scripts# ansible all -m ping
172.16.1.86 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
172.16.1.88 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
172.16.1.87 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
172.16.1.92 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
172.16.1.89 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
172.16.1.93 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
172.16.1.94 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
172.16.1.95 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
172.16.1.96 | SUCCESS => {
    "changed": false, 
    "ping": "pong"
}
```
