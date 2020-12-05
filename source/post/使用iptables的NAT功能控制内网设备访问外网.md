```toml
title = "使用iptables的NAT功能控制内网设备访问外网"
date = "2016-12-08 19:59:19"
update_date = "2016-12-08 19:59:19"
author = "KDF5000"
thumb = ""
tags = ["Linux"]
draft = false
```
通过iptables的NAT功能控制内网上网。

##### 前提条件

一台能够上网的主机，并且和其他需要控制上网的主机在同一个内网。

hw98: 172.18.11.98 可以上外网

hw网段: 172.18.11.0/24

目的：通过hw98装发内网内的其他机器数据包，从而实现上网控制的目的

#### 步骤

hw98设置路由转发的功能

* 开启内核的ip转发功能

  ```shell
  $ echo 1 > /proc/sys/net/ipv4/ip_forward
  或者
  $ vim /etc/sysctl.conf
  # Controls IP packet forwarding
  net.ipv4.ip_forward = 1
  ```

* 添加确认包和关联包的通过

  ```shell
  $ iptables -A FORWARD -m state --state ESTABLISHED,RELATED -j ACCEPT
  ```

<!--more -->

* 设置iptables的nat表，添加FORWARD规则

  ```shell
  $ iptables -t nat -A POSTROUTING -s 172.18.11.0/24 -j SNAT --to 172.18.11.98
  ```

内网内需要进行上网的机器

* 删除默认的路由，如果没有就不用删除

  ```shell
  $route del default gw *.*.*.*  #通过route 查看
  ```

* 添加hw98为默认路由

  ```shell
  $route add default gw 172.18.11.98
  ```

上面内网的设置在重启网络的时候会被修改，所以最好的方法就是在网络配置里就行设置，在centos就是`/etc/sysconfig/network-scripts/ifcfg-eth0`，只用设置GATEWAY即可，如果有必要也可以设置一下DNS

```
GATEWAY=172.18.11.98
DNS1=159.226.39.1
```

**说明：**没有深入理解整个运行过程，其实可以通过hw98的iptables设置进行更加细致的访问控制，比如控制内网内特定的主机访问外网，这里是全部允许，也可以控制特定的端口等等。
