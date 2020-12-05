```toml
title = "Ubuntu 14.04 下安装使用Python rq模块"
date = "2015-08-23 14:00:16"
update_date = "2015-08-23 14:00:16"
author = "KDF5000"
thumb = ""
tags = ["Python", "Ubuntu", "分布式", "Python队列"]
draft = false
```
`rq` 是`Python`的一个第三方模块，使用`rq`可以方便快速的实现`Python`的队列操作，实现多态电脑的分布式架构。其中 *R*是`Redis`的意思，*Q*是`Queue`的首字母，`rq`使用`Redis`和`Queue`实现分布式，分别实现了`Master`和`Worker`，通过`Redis`存储任务队列。
####Ubuntu14.04 安装rq
假设已经安装了`Python`和`pip`，本文通过``pip`来安装`rq`
```
$sudo pip install rq 
```

####安装`Redis`
`rq`模块使用`redis`保存队列信息，因此可以保证多台机器同时读取同一个队列，也就是多个``worker`同时工作，这也就达到了我们的目的。在`Ubuntu` 下安装`Redis`比较简单，使用下面的命令即可，该命令除了安装 `Redis`外，也会好心地帮你安装了`redis-cli`。 
```
$sudo apt-get install redis-server
```

<!--more-->

安装完成后可以尝试启动一下`Reids`，检查是否安装成功。
```
$ redis-server
```
上面的命令会使用默认的设置启动`Redis`服务，如果你看到下面漂亮启动界面说明安装成功了。

不过还没完额，使用下面命令看看我们可以看到什么
```
$ netstat -an | grep 6379
```
结果：

![](@media/archive/img_rq-redis-bind.png)

因为`Redis`默认使用的端口是6379，该命令可以查看6379端口监听的ip ，可以看到 `Redis`默认绑定的是`127.0.0.1`，可以在`/etc/redis/redis.conf`中看到该设置。

![](@media/archive/img_rq-redis-redis-conf.png)

`Redis`的默认配置绑定了`127.0.0.1`，注释掉**bind 127.0.0.1**即可。然后重启`Redis`。
```
$ sudo /etc/init.d/redis-server restart
```
再次执行`netstat -an | grep 6379`

![](@media/archive/img_rq-redis-redis-netstat.png)

可以看到改变了 ，`Redis`已经可以接受同一个局域网内的`redis cli`连接了
####安装`rq-dashboard`
`rq-dashboard`是一个监控`rq`执行状况的`python`库，它可以显示当前有哪些`Queue`，每个`Queue`有多少`Job`，以及有多少`Worker`处于工作状态，还显示了失败的`Job`。可以使用`pip`方便的安装`Dashboard`.
```
$sudo pip install rq-dashboard
```

安装成功后，使用下面的命令启动`rq-dashboard`
```
$rq-dashboard -u "redus://192.168.0.107:6379"
```
其中`-u`参数是需要使用的`Redis`连接地址，启动成功后可以看到下面的信息

![](@media/archive/img_rq-rqdashboard-start.png)

可以看出`Rq dashboard`的版本信息，以及运行的地址端口，也就是我们可以通过浏览器访问，默认的端口是*9181*，`IP`地址是启动`rq-dashboard`的机器`ip`，在同一局域网的电脑访问`http://192.168.0.107:9181`，其中`192.168.0.107`是启动`rq-dashboard`的电脑`ip`。

![](@media/archive/img_rq-rqdashboard-web.png)

`Rq-dashboard`是一个很有用的工具，可以图形化的监控`rq`的工作状态，但是美中不足，不能控制`worker`的工作，不过相信应该很快就会支持这些功能了。

####`rq`的使用
**参考[官方文档](http://python-rq.org/)**

