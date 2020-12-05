```toml
title = "Python学习速记--操作文件和目录"
date = "2015-07-21 12:23:33"
update_date = "2015-07-21 12:23:33"
author = "KDF5000"
thumb = ""
tags = ["Python", "学习笔记"]
draft = false
```
####操作文件和目录
python有一个os模块封装了操作与底层相关的函数，比如查看环境变量，系统名字，操作目录等
#####查看系统参数
```
import os
#查看系统参数
print(os.name) #系统名字
print(os.environ) #查看环境变量
print(os.environ.get('PATH')) #获取指定的环境变量值
```
#####操作文件和目录
利用`os`模块提供的功能可以操控目录，包括查看，添加，删除，遍历

<!--more-->

```
#操作文件和目录
cur_path = os.path.abspath('.')# 获取当前目录
print(cur_path)
new_path = os.path.join(cur_path,'test') # 通过join获取新的目录，不要自己组装，不同系统下的分隔符是不一样的
print(new_path)
os.mkdir(new_path) #创建一个目录
#os.rmdir(new_path) # 移除一个目录
#遍历目录
for x in os.listdir('.'):
	if os.path.isdir(x):
		print(x)
	elif os.path.isfile(x)  and os.path.splitext(x) == '.py':
		print('python file-->',x)
	else:
		print('file-->',x)
```
