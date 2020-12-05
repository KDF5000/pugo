```toml
title = "Python学习速记--IO操作"
date = "2015-07-21 12:22:18"
update_date = "2015-07-21 12:22:18"
author = "KDF5000"
thumb = ""
tags = ["Python", "学习笔记"]
draft = false
```
#### 文件读写
##### 读取文本文件
```
#IO操作--文件读取
f = open('test.txt','r') # 和c一样打开文件
#print(f.read())  # read()读取所有内容

#print(f.readline()) # readline 读取一行

#print(f.readlines()) # 读取所有内容到list,每一行一个单位

#异常处理
try:
	f = open('test1.txt','r')
except IOError:
	print('打开文件失败')
finally:
	if f:
		f.close()

#python提供了更简洁的方式
with open('test.txt','r') as f:
	print(f.read())
```

<!--more-->

#####读取二进制文件
读取文本文件是以UTF-8进行编码，读取二进制文件其实和读取文本文件是一样的，还是打开的模式不一样使用`rb`打卡文件，读出来的内容是二进制数据
```
#读取二进制文件
f = open('Penguins.jpg','rb') # 以rb打开
print(f.read())
```

#####字符编码
读取文本文件默认是utf-8进行编码，也可以在打开文件的时候指定读取的编码，甚至可以自定遇到编码错误时如何处理错误，如下所示
```
#字符编码
f = open('test.txt', 'r', encoding='gbk', errors='ignore')
print(f.read())
```

#####写文件
写文件与打开文件也比较类似，主要区别在于打开文件时的模式不同，写文件使用`w`或者`wb`，也可以指定编码，由于写文件不是立即写到磁盘，而是先写到缓冲区，随后写入，因此在写完后一定有节的调动`close()`关闭文件，这样才真正写到磁盘，为了避免忘记写`close`，最好使用`with`打开文件。示例如下：
```
#写文件
with open('demo.txt','w')  as f:
	f.write('nihaoma')
```

####文件读写
#####读取文本文件
```
#IO操作--文件读取
f = open('test.txt','r') # 和c一样打开文件
#print(f.read())  # read()读取所有内容

#print(f.readline()) # readline 读取一行

#print(f.readlines()) # 读取所有内容到list,每一行一个单位

#异常处理
try:
	f = open('test1.txt','r')
except IOError:
	print('打开文件失败')
finally:
	if f:
		f.close()

#python提供了更简洁的方式
with open('test.txt','r') as f:
	print(f.read())

```
#####读取二进制文件
读取文本文件是以UTF-8进行编码，读取二进制文件其实和读取文本文件是一样的，还是打开的模式不一样使用`rb`打卡文件，读出来的内容是二进制数据
```
#读取二进制文件
f = open('Penguins.jpg','rb') # 以rb打开
print(f.read())
```

#####字符编码
读取文本文件默认是utf-8进行编码，也可以在打开文件的时候指定读取的编码，甚至可以自定遇到编码错误时如何处理错误，如下所示
```
#字符编码
f = open('test.txt', 'r', encoding='gbk', errors='ignore')
print(f.read())
```

#####写文件
写文件与打开文件也比较类似，主要区别在于打开文件时的模式不同，写文件使用`w`或者`wb`，也可以指定编码，由于写文件不是立即写到磁盘，而是先写到缓冲区，随后写入，因此在写完后一定有节的调动`close()`关闭文件，这样才真正写到磁盘，为了避免忘记写`close`，最好使用`with`打开文件。示例如下：
```
#写文件
with open('demo.txt','w')  as f:
	f.write('nihaoma')
```
