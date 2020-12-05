```toml
title = "auditd的使用"
date = "2016-11-12 19:37:32"
update_date = "2016-11-12 19:37:32"
author = "KDF5000"
thumb = ""
tags = ["Linux", "安全"]
draft = false
```
linux系统中用户空间的组件，负责系统安全审计，将审计的记录写入磁盘

#### 安装

ubuntu可以直接使用apt-get安装，centos使用yum安装

安装后有下面相关的工具：

- **auditctl :** 即时控制审计守护进程的行为的工具，比如如添加规则等等。
- **/etc/audit/audit.rules :** 记录审计规则的文件。
- **aureport :** 查看和生成审计报告的工具。
- **ausearch :** 查找审计事件的工具
- **auditspd :** 转发事件通知给其他应用程序，而不是写入到审计日志文件中。
- **autrace :** 一个用于跟踪进程的命令。
- **/etc/audit/auditd.conf :** auditd工具的配置文件。

首次安装auditd后，审计规则是空的，使用下面的命令查看

```shell
$auditctl -l
```

<!-- more -->

#### 添加auditd规则

监控文件和目录的更改

```shell
$sudo auditctl -w /etc/passwd -p rwxa
```

选项：

* -w path: 指定监控的路径，上面指定监控的文件路径/etc/passwd
* -p: 指定出发审计的文件/目录的访问权限
* rwxa: 指定的触发条件，r读，w写，x执行权限，a属性

如果不加-p参数，-w参数后面指定的是目录，则将会对指定目录所有访问进行监控。默认的权限是rwax，可以使用`auditctl -l`查看已经添加的规则

![](@media/archive/blog/image/auditd-0.png)



#### 查看审计日志

使用ausearch可以查看auditd的审计日志。比如上面已经对/etc/passwd文件添加了审计，那么使用下面的命令可以查看审计日志。

```shell
$sudo ausearch - f/etc/passwd
```

`-f`需要查找的审计目标的日志

输出可能如下

> ----
> time->Wed Sep 21 15:07:51 2016
> type=PATH msg=audit(1474441671.907:881): item=1 name="/root/kongdefei/.test.swp" inode=44826697 dev=08:03 mode=0100600 ouid=0 ogid=0 rdev=00:00
> type=PATH msg=audit(1474441671.907:881): item=0 name="/root/kongdefei/" inode=80856866 dev=08:03 mode=040755 ouid=0 ogid=0 rdev=00:00
> type=CWD msg=audit(1474441671.907:881):  cwd="/root/kongdefei"
> type=SYSCALL msg=audit(1474441671.907:881): arch=c000003e syscall=87 success=yes exit=0 a0=699f130 a1=1 a2=1 a3=2 items=2 ppid=18763 pid=18798 auid=4294967295 uid=0 gid=0 euid=0 s
> uid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=pts1 ses=4294967295 comm="vi" exe="/bin/vi" key=(null)

* time: 审计时间
* name: 审计对象
* cwd: 当前路径
* sys call：系统调用
* auid: 设计用户id
* uid和gid：访问文件的用户id和用户组id,uid=0 gid=0表明是root用户
* comm: 用户访问文件的命令
* exe: 上面命令的执行文件路径

上面内容表示/root/kongdefei/.test.swp文件被root用户使用vi命令编辑过



#### 查看审计报告

一旦定义规则后，他会自动运行，过一段时间后，可以使用auditd的工具aureport生成简要的日志报告。直接使用下面命令即可

```shell
$sudo aureport
```

 ![](@media/archive/blog/image/auditd-1.png)

可以看出有51次授权失败，102次登录。可以使用下面命令查看授权失败的详细信息

```shell
$sudo aureport -au
```

 ![](@media/archive/blog/image/auditd-4.png)

凡是no（如2，5）的都是授权失败。可以看出是root用户ssh登录的时候失败。

`-m`可以查看所有账户修改相关的事件

```shell
$sudo aureport -m
```

 ![](@media/archive/blog/image/auditd-2.png)



#### auditd配置文件

上面使用-w添加审计规则只是暂时的，系统重启后就会消失，可以修改/etc/audit/audit.rules文件是规则持久有效。上面添加的规则可以直接写入/etc/audit/audit.rules文件中

如下：

 ![](@media/archive/blog/image/auditd-3.png)

然后重启auditd即可

```shell
$ service auditd restart #或者/etc/init.d/auditd restart
```



#### 总结

Auditd是Linux上的一个审计工具。你可以阅读auidtd文档获取更多使用auditd和工具的细节。例如，输入 **man auditd** 去看auditd的详细说明，或者键入 **man ausearch** 去看有关 ausearch 工具的详细说明。

**请谨慎创建规则**。太多规则会使得日志文件急剧增大！
