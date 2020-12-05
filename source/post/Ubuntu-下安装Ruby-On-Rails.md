```toml
title = "Ubuntu 下安装Ruby On Rails"
date = "2015-09-25 23:39:46"
update_date = "2015-09-25 23:39:46"
author = "KDF5000"
thumb = ""
tags = ["Ruby", "Rails"]
draft = false
```
###Ruby On Rails的配置
`Rails`是一个优秀的`Ruby`框架.结合了`PHP`快速开发和`JAVA`程序规整的特点,使用`MVC`的设计模式,是一个用于开发数据库驱动的网络应用程序的完整框架.

`Ruby On Rails`的官网: [https://www.ruby-lang.org](https://www.ruby-lang.org)  

安装`Ruby On Rails`需要安装`Ruby`,`gem`,`Sqlite3`,和`Rails`. 本教程的平台使用的使用的`Ubuntu`,因此只介绍`Ubuntu`下的安装.

###安装`Ruby`
安装`Ruby`可以参考`Ruby`的[官方教程](https://www.ruby-lang.org/en/documentation/installation/)

再`Ubuntu`下可以使用`sudo apt-get install ruby-fully`进行快速的安装,但是使用这种方式安装的最新版本是`1.9.3`,而最新版本的`Ruby`是`2.2.3`,因此本教程使用编译源码进行安装.

到`Ruby`下载最新的[源码](https://cache.ruby-lang.org/pub/ruby/2.2/ruby-2.2.3.tar.gz),解压到指定的位置,进入解压后的目录,执行下面的命令:
```
$./configuren
$ make 
$ sudo make install
```
耐心的等待一段时间,执行下面的命令:
```
$ ruby -v
```
如果出现下图中的结果,说明安装成功了.  

![](@media/archive/img_ruby-v.png)

<!--more-->

###安装`Gem`
`Gem`是`Ruby`的包管理工具,相当于`Python`的`pip`,`PHP`的`Composer`和`PEAR`.

可以去`RubyGem`的[官网](https://rubygems.org/)查看更多的信息.

安装`Gem`也非常简单,下载最新的[RubyGem](https://rubygems.org/rubygems/rubygems-2.4.8.tgz),然后解压到指定目录,进入目录,执行下面的命令即可
```
$ sudo ruby setup.rb 
```
如果出现下面的结果说明安装成功

![](@media/archive/img_gem-v.png)

###安装`sqlite3`
要想使用`Rails`,需要安装`Sqlite3`,大多数的`Unix`系统已经安装的有了,可以使用**sqlite3 -verion**检查是否安装了`sqlite3`，如果已经安装，则直接进入下一步，如果没有安装则，运行下面的命令进行安装
```
$ sudo apt-get install sqlite3
```

###安装`Rails`
终于到了安装`Rails`的时候,当然这也是最后一步,安装`Rainls`使用`Gem`安装即可,非常方便,再控制台下执行下面的命令即可.
```
$sudo gem install rails
```
安装完成后执行下面的命令检查`Rails`是否安装成功.
```
$rails -v
```
如果出现下图,恭喜你安装成功了.  

![](@media/archive/img_rails-v.png)

###测试`Rails`
使用下面的命令创建一个新的`Rails`项目
```
$sudo rails new course   # course是项目的名称
```
然后启动`Rails Server`,使用下面的命令
```
$sudo rails server
```
如果非常不幸出现了下面的错误信息,则说明没有`js`的运行环境,有两种解决方案.

* 在Gemfile里添加下面两行,然后执行`bundle install`即可.
> 
gem 'execjs'  
gem 'therubyracer'

* 安装`Node.js`,
> $sudo apt-get install nodejs

![](@media/archive/img_rails-server-failed.png)

再次启动`Rails server`,如果出现下面的信息则说明启动成功.

![](@media/archive/img_server_start.png)

然后访问**http://localhost:3000**,出现下面的页面则说明安装`Rails`成功.

![](@media/archive/img_server.png)

**PS:** 在启动`Rails server`的过程中可能还会出现一个缺少`libreadlinedev`(确切的名字不太记得了)的错误,这时只需按照提示安装缺失库,然后重新编译安装`Ruby`即可.

