```toml
title = "Laravel入门教程之一--安装"
date = "2015-07-27 12:11:58"
update_date = "2015-07-27 12:11:58"
author = "KDF5000"
thumb = ""
tags = ["PHP"]
draft = false
```
`Laravel`是一套简洁、优雅的`PHP Web`开发框架(`PHP Web Framework`)。它可以让你从面条一样杂乱的代码中解脱出来；它可以帮你构建一个完美的网络`APP`，而且每行代码都可以简洁、富于表达力。
在安装`Laravel`之前需要先安装`PHP`环境，    `Laravel`的安装是通过[`Composer`](https://getcomposer.org/)安装的，所以必须先安装`PHP`和`Composer`。
####安装`PHP`    
本教程使用最常用的`PHP`集成环境`XAMPP`，`XAMPP`已经集成了`Apache`，`PHP`，`MySQL`等开发环境，非常方便，而且搭建简单。
从`XAMPP`的[官网](https://www.apachefriends.org/zh_cn/index.html)下载安装包即可，本教程是在`Windows`下进行的，所以下载`Windows`版即可。
下载完成后，进行安装，基本上是一路`Next`。安装完成后在开始里会出现一个`XAMPP`的文件夹，点击里面的`XAMPP Control Panel`既可启动`Xampp`的控制面板，`Apache`,`MySQL`等启动停止都可以在这里控制。如下图所示
![](@media/archive/img_PHP_ENV_1.png)
启动`Apache`，在浏览器里输入`http://localhost`，如果出现如下页面则说明配置成功
![](@media/archive/img_PHP_ENV_2.png)
本教程安装`PHP`的介绍到此为止，网上有很多关于`XAMPP`安装以及配置的教程，可以自行搜索。如果安装过程中出现问题，可以评论留言。
####安装`Composer`
`Composer`是`PEAR`之后的一个很好用的包管理工具，也是安装`Laravel`的必备工具，详细的安装可以参考鄙人的博客[PHP依赖管理Composer](http://kdf5000.github.io/2015/07/26/PHP%E4%BE%9D%E8%B5%96%E7%AE%A1%E7%90%86Ciomposer%E7%9A%84%E5%AE%89%E8%A3%85/)。
####安装`Laravel`
`Laravel`的安装有两种方法，一种是安装`Laravel`到本地环境，将`Laravel`加入系统变量，使用`Laravel`内置的命令创建项目。另一种是使用`Composer`的`create-project`命令创建项目。当然两种方式各有各的优点，第一种方法使用更加简单，第二种方法不需要安装`Laravel`环境
。应该有其他比较大的差一点，暂时还不知道，等日后研究一下分享给大家，如果知道的大家分享一下。

<!--more-->

#####使用`Composer`命令创建`Laravel`项目
使用`Composer`创建比较简单，代开命令行，输入下面指令即可
```
composer create-project laravel/laravel learnlaravel5 5.0.22
```
其中`learnlaravel5`是创建的项目的名字，`5.0.22`是使用的`Laravel`版本号，也可以不指定，默认应该是使用最新版本。这个过程会下载很多组件，时间可能会比较久，留意一下命令行下载过程如果出现错误，立刻停止(`Ctrl + C`)，看看错误原因是什么。
创建完成后，可以在`Apache`创建一个虚拟主机，将目录只想新建的`Laravel`项目的`Public`目录下。
* 修改`http-vhosts.conf`，该文件在`xampp安装目录\apache\conf\extra`
```
Listen 8008  #监听的端口
<VirtualHost *:8008>
    ##ServerAdmin webmaster@dummy-host.example.com
    DocumentRoot "F:/CodingDemo/Laravel/LearnLaravel/learnlaravel5/public"  #laravel项目的public目录
    ##ServerName dummy-host.example.com
    ##ServerAlias www.dummy-host.example.com
    ErrorLog "logs/laravel.com-error.log"
    CustomLog "logs/laravel.com-access.log" common
</VirtualHost>
```
该配置指定了`apache`服务器监听`8008`是访问`Laravel`项目。此时启动`Apache`，访问`localhost:8008`，可能出现下面的`Access forbidden`错误页面。
![](@media/archive/img_laravel_2.png)

* 出现上面的访问禁止的原因主要是`Apache`配置拒绝了所有的请求，最简单的就是注释掉，打开`xampp安装目录\apache\conf\httpd.conf`，修改下面的地方，注释掉`Require all denied`，将`AloolwOverride`改为`All`。
```
<Directory />
    AllowOverride All
    #Require all denied
</Directory>
```
再次访问`localhost:8008`，出现下面的页面，说明安装成功。
![](@media/archive/img_laravel_3.png)


