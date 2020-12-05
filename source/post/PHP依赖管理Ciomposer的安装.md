```toml
title = "PHP依赖管理Ciomposer的安装"
date = "2015-07-26 22:22:28"
update_date = "2015-07-26 22:22:28"
author = "KDF5000"
thumb = ""
tags = ["PHP", "Laravel"]
draft = false
```
对于现代语言，包管理基本上市标配，通过包管理可以大大提高项目的开发效率，也使得项目的结构更加的清晰。`Java`有`Maven`，`Python`有`PIP`，`Ruby`有`gem`，`NodeJs`有`PEAR`，不过`PEAR`坑不少：
* 依赖处理容易出问题
* 配置复杂
* 难用的命令行接口

`Composer`的出现，是的`PHP`的依赖管理更加容易，相比`PEAR`简单易用，而且是开源的，提交自己的包夜很容易

####安装`Composer`
`Composer`的需要`PHP 5.3.2+`才能运行
#####`Windows`下安装
进入`Composer`的[官网](https://getcomposer.org/download/)下载`Windows Installer`进行安装，安装的前提是你已经安装了`PHP`，安装过程中会让你选择`PHP`的可执行程序路径，也可以选择将`Composer`组件加入到右键快捷键，如下图所示
![](@media/archive/img_composer_1.png)
安装成功后，打开命令行，执行`composer --version`，如果出现下图所示，则说明安装成功。
![](@media/archive/img_composer_2.png)

修改`Composer`的配置文件，指定获取包的镜像服务地址。

<!--more -->

打开命令行，输入下面的命令即可。此方法是修改全局配置，也可以针对单个项目进行配置。参考[http://pkg.phpcomposer.com/](http://pkg.phpcomposer.com/)
```
composer config -g repositories.packagist composer http://packagist.phpcomposer.com
```
#####类`Unix`的安装
这里介绍一种全局安装的方法，其他放发可以参考[http://docs.phpcomposer.com/00-intro.html](http://docs.phpcomposer.com/00-intro.html)。
此方法将`Composer`放在系统的`PATH`里，这样就可一再任何地方使用`Composer`命令。
```
curl -sS https://getcomposer.org/installer | php
mv composer.phar /usr/local/bin/composer
```
