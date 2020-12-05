```toml
title = "Ubuntu 14.04 安装PHP环境和反向代理"
date = "2015-08-24 20:47:55"
update_date = "2015-08-24 20:47:55"
author = "KDF5000"
thumb = ""
tags = ["PHP", "反向代理"]
draft = false
```
####安装Apache
```
$ sudo apt-get install apache2
```

####安装php
```
$ sudo apt-get install php5 libapache2-mod-php5
```

####安装mysql
```
$ sudo apt-get install mysql-server
```

#### 反向代理
经过测试最小的配置。。。。

<!--more-->

* 启用apache的mod_proxy 模块
```
$ sudo a2enmod mod_proxy
```
* 修改配置文件 /etc/apache2/sites-enabled/000-default.conf （此文件是默认的80端口的配置文件，也可以添加在自己想添加的虚拟主机配置文件），在<VirtualHost></VirtualHost>内添加下面的代码
```
ProxyPassReverse    /      http://192.168.0.3:8006/
ProxyPass           /      http://192.168.0.3:8006/
```
如下所示

![](@media/archive/img_apache-proxy.png)

然后重启apache服务器
```
$ sudo service apache2 restart
```

如果不成功，则尝试进行下面的操作
* 重新load apache
```
$ sudo service apache2 reload
```

* 在最开始配置代理的地方添加下面两句
```
ProxyPreserveHost On
ProxyRequests On
```

* 添加代理的外部访问权限，在配置虚拟主机的地方添加下面几句
```
<Proxy *>
       Order deny,allow
       Allow from all
</Proxy>
```

**记得每次修改都要重启`Apache`**


