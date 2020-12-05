```toml
title = "使用Python和BitTorrentSync定期给Kindle推送电子书"
date = "2016-11-12 18:45:59"
update_date = "2016-11-12 18:45:59"
author = "KDF5000"
thumb = ""
tags = ["Python", "Kindle"]
draft = false
```
#### Kindle伴侣

最近发现一个很好用的Kindle电子书分享网站[Kindle伴侣](http://kindlefere.com/)，资源丰富，除了提供各种热门图书免费下载，还提供了很多关于Kindle的使用技巧，里面个人最喜欢的功能就是每周一书，使用BitTorrent每周自动同步热门电子书。使用的是著名的同步软件Resilio Sync， 最开始是在自己电脑上安装了软件同步，但是要一直开着，感觉好麻烦，有时候电脑还要关了，影响同步。刚好手里有个闲置的树莓派，于是就想着用它来同步，然后通过邮件将电子书发送到kindle，这样就可以每周都看到热门的电子书了。

<!--more-->

![kindle-0](@media/archive/blog/image/kindle-0.png)

#### Resilio Sync服务搭建

Linux下的安装方法可以参考官网：[https://help.getsync.com/hc/en-us/articles/206178924](https://help.getsync.com/hc/en-us/articles/206178924)

下面的方法适合Raspbian，你可以按照官网的说明安装在任何一台可以上网的服务器上。

**Step1:** 下载deb安装包

```shell
$ wget https://download-cdn.getsync.com/2.0.128/PiWD/bittorrent-sync-pi-server_2.0.128_armhf.deb
```

**Step2:** 安装Sync包

```SHELL
$ sudo dpkg -i bittorrent-sync-pi-server_2.0.128_armhf.deb
```

**Step3:**访问Sync Web UI

通过浏览器访问http://yourrespbianip:8888 ,其中是你respian的ip地址或者你托管的服务器地址，你应该看到下面的界面 

 ![kindle-1](@media/archive/blog/image/kindle-1.png)

然后到[Kindle伴侣](http://kindlefere.com/)随便点一个[每周一书]可以找到使用bittorrentSync同步的同步key，通过web页面添加到你自己的sync服务器即可。记得你选择的文件夹要有写入权限，我是直接给的777。

#### 监控进程

这里使用Python作为监控程序监控上面Sync服务器的同步文件夹，一旦出现新的指定格式（如MOBI)的文件，则将会向指定的kindle邮箱发送新增加的电子书，邮箱的设置需要到amazon自己设置，想必用过kindle的应该不陌生。

代码可以从github下载：[https://github.com/KDF5000/WeeklyBook](https://github.com/KDF5000/WeeklyBook) , 然后修改`server.py`里面的kindle邮箱，以及授权的邮箱

> EXT_LIST = ['mobi'] #想要发送的邮件格式
> HOST_NAME = 'smtp.cstnet.cn' #授权邮箱的smtp服务器
> HOST_PORT = 25 #smtp服务器端口
> USER_NAME = 'kongdefei@ict.ac.cn' #授权邮箱的用户名
> USER_PASS = '*********************'
> KINDLE_MAILS = ['kdf5000@kindle.cn'] #接收电子书的kindle邮箱，可以在亚马逊查看
> FROM_NAME = 'kongdefei@ict.ac.cn' #发送邮件的from名字，建议使用发送的邮箱地址

##### 后台启动

进入下载的文件夹，执行下面的命令

```shell
$nohup python server.py /home/ubuntu/kongdefei/SyncBook > sync.out 2>&1 &
```

其中`/home/ubuntu/kongdefei/SyncBook`是sync同步的文件夹，sync.out是server.py的输出日志，该程序将会在后台持续运行监控`/home/ubuntu/kongdefei/SyncBook`下面是否有新的电子书的出现。


至此，你就可以享受你心爱的kindle每周定时收到一本优秀的电子书了。
顺便晒一下树莓派[憨笑]
![](@media/archive/blog/image/raspberry.jpeg)
