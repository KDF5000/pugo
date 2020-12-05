```toml
title = "Ubuntu下Scrapy的安装"
date = "2015-08-11 23:06:56"
update_date = "2015-08-11 23:06:56"
author = "KDF5000"
thumb = ""
tags = ["Python", "Scrapy", "Linux"]
draft = false
```
最近在学习爬虫，早就听说`Python`写爬虫极爽（貌似pythoner说python都爽，不过也确实，python的类库非常丰富，不用重复造轮子），还有一个强大的框架`Scrapy`，于是决定尝试一下。

要想使用`Scrapy`第一件事，当然是安装`Scrapy`，尝试了`Windows`和`Ubuntu`的安装，本文先讲一下 `Ubuntu`的安装，比`Windows`的安装简单太多了。。。抽时间也会详细介绍一下怎么在`Windows`下进行安装。

[官方介绍](http://scrapy-chs.readthedocs.org/zh_CN/latest/intro/install.html#ubuntu-9-10)，在安装`Scrapy`前需要安装一系列的依赖.
*  `Python 2.7`： `Scrapy`是`Python`框架，当然要先安装`Python` ，不过由于`Scrapy`暂时只支持 `Python2.7`，因此首先确保你安装的是`Python 2.7`
*  `lxml`：大多数`Linux`发行版自带了`lxml`
*  `OpenSSL`：除了`windows`之外的系统都已经提供
*  `Python Package`: pip and setuptools. 由于现在`pip`依赖`setuptools `,所以安装`pip`会自动安装`setuptools `

有上面的依赖可知，在非windows的环境下安装 Scrapy的相关依赖是比较简单的，只用安装`pip`即可。`Scrapy`使用`pip`完成安装。

<!--more-->

####检查`Scrapy`依赖是否安装
你可能会不放心自己的电脑是否已经安装了，上面说的已经存在的依赖，那么你可以使用下面的方法检查一下，本文使用的是`Ubuntu 14.04`。
#####检查`Python`的版本
```
$ python --version
```
如果看到下面的输出，说明`Python`的环境已经安装，我这里显示的是`Python 2.7.6`，版本也是`2.7`的满足要求。如果没有出现下面的信息，那么请读者自行百度安装`Python`，本文不介绍`Python`的安装（网上一搜一堆）。
![](@media/archive/img_scrapy_python.png)

#####检查`lxml`和`OpenSSL `是否安装
假设已经安装了`Python`，在控制台输入`python`，进入`Python`的交互环境。
![](@media/archive/img_scrapy_lxml_ssl.png)

然后分别输入`import lxml`和`import OpenSSL`如果没有报错，说明两个依赖都已经安装。
![](@media/archive/img_scrapy_lxml_openssl.png)

####安装`python-dev`和`libevent`
`python-dev`是`linux`上开发`python`比较重要的工具，以下的情况你需要安装
* 你需要自己安装一个源外的python类库, 而这个类库内含需要编译的调用python api的c/c++文件
* 你自己写的一个程序编译需要链接libpythonXX.(a|so)

`libevent`是一个时间出发的高性能的网络库，很多框架的底层都使用了`libevent`

上面两个库是需要安装的，不然后面后报错。使用下面的指令安装
```
$sudo apt-get install python-dev
$sudo apt-get install libevent-dev
```
####安装`pip`
因为`Scrapy`可以使用`pip`方便的安装，因此我们需要先安装`pip`，可以使用下面的指令安装`pip`
```
$ sudo apt-get install python-pip
```

####使用`pip`安装`Scrapy`
使用下面的指令安装`Scrapy`。
```
$ sudo pip install scrapy
```
记住一定要获得`root`权限，否则会出现下面的错误。
![](@media/archive/img_scrapy_exception.png)

至此`scrapy`安装完成，使用下面的命令检查`Scrapy`是否安装成功。
```
$ scrapy version
```
显示如下结果说明安装成功，此处的安装版本是`1.02`
![](@media/archive/img_scrapy_version.png)


