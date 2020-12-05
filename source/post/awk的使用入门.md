```toml
title = "awk的使用入门"
date = "2016-04-30 16:05:46"
update_date = "2016-04-30 16:05:46"
author = "KDF5000"
thumb = ""
tags = ["Linux"]
draft = false
```
#### 简介
`awk`是一个强大的文本分析工具，在对数据分析生成报告时，特别枪弹。`awk`的原理是把文件以行读入，以指定的分隔符对每行分片，提供一些函数或者自己写程序逻辑，可以对每片进行出来。

#### 使用方法
```
awk [options] 'script' var=value file(s) 
awk [options] -f scriptfile var=value file(s)
```
##### 常用命令选项
* **-F field_separator** 制定每行文本的域分隔符，可以是字符串也可以是正则表达式
* **-v var=value** 赋值一个用户变量，将外部变量值传入awk
* **-f** 制定awk脚本文件

#### awk的模式和操作
awk 脚本，也就是上一节中的script部分，格式遵循下面的形式
> awk "pattern action"  files
#####模式(Pattern)
模式可以是下面的任意一个:
* /正则表达式/: 使用//包围的正则表达式
* 关系表达式：使用运算符进行操作，可以是字符串或数字的比较测试
* 模式匹配表达式： 用运算符`~`(匹配)和`~!`(不匹配)
*  BEGIN语句块，pattern匹配块，END语句块

<!-- more -->

##### 操作
操作有一个或者多个命令，函数，表达式组成，有换行符或者分毫分开，必须位于大括号内部，主要有下面操作：
* 变量或数组赋值
* 输出命令
* 内置函数
* 控制流程语句

awk操作中的变量或者数组可以直接使用，不用进行初始化，其他语法和c风格比较像

#### 入门实例
假设last -n 5的输出如下：
```
hadoop@kdf5000:~$ last -n 5
hadoop   pts/9        :1               Sat Apr 30 15:18   still logged in   
hadoop   pts/22       :1               Sat Apr 30 14:49 - 14:49  (00:00)    
hadoop   :1           :1               Sat Apr 30 14:48   still logged in   
hust-lh  :0           :0               Sat Apr 30 14:46   still logged in   
hadoop   pts/1        :0               Sat Apr 30 14:03 - 14:46  (00:43)    
```
如果只显示最近登陆的5个账号
```
hadoop@hust-lh:~$ last -n 5 | awk '{print $1}'
hadoop
hadoop
hadoop
hust-lh
hadoop
```
awk工作流程是这样的：读入有`\n`换行符分割的一条记录，然后将记录按指定的域分隔符划分域，填充域，`$0` 则表示所有域,`$1`表示第一个域,`$n`表示第n个域。默认域分隔符是"空白键" 或 "[tab]键",所以`$1`表示登录用户，`$3`表示登录用户ip,以此类推。

如果只是显示/etc/passwd的账户
```
hadoop@hust-lh:~$ cat /etc/passwd |awk  -F ':'  '{print $1}'  
root
daemon
bin
sys
sync
games
man
```
这种是`awk+action`的示例，每行都会执行`action{print $1}`。`-F`指定域分隔符为`':'`。
如果只是显示`/etc/passwd`的账户和账户对应的`shell`,而账户与`shell`之间以逗号分割,而且在所有行添加列名`name`,`shell`,在最后一行添加`"blue,/bin/nosh"`。
```
hadoop@hust-lh:~$ cat /etc/passwd |awk  -F ':'  'BEGIN {print "name,shell"}  {print $1","$7} END {print "blue,/bin/nosh"}'
name,shell
root,/bin/bash
daemon,/usr/sbin/nologin
bin,/usr/sbin/nologin
sys,/usr/sbin/nologin
sync,/bin/sync
blue,/bin/nosh
```
awk工作流程是这样的：先执行`BEGING`，然后读取文件，读入有`\n`换行符分割的一条记录，然后将记录按指定的域分隔符划分域，填充域，`$0`则表示所有域,`$1`表示第一个域,`$n`表示第`n`个域,随后开始执行模式所对应的动作`action`。接着开始读入第二条记录······直到所有的记录都读完，最后执行`END`操作。

搜索/etc/passwd有root关键字的所有行
```
hadoop@hust-lh:~$ awk -F: '/root/' /etc/passwd
root:x:0:0:root:/root:/bin/bash
```
这种是`pattern`的使用示例，匹配了`pattern`(这里是`root`)的行才会执行`action`(没有指定`action`，默认输出每行的内容)。

搜索支持正则，例如找`root`开头的: `awk -F: '/^root/' /etc/passwd`
搜索/etc/passwd有root关键字的所有行，并显示对应的`shell`
```
hadoop@hust-lh:~$ awk -F: '/root/{print $7}' /etc/passwd          
/bin/bash
```
这里指定了action{print $7}

#### 脚本编程
awk可是将脚本写到文件中，然后使用-f调用脚本文件
下面是一个统计一个文件中对应的单词出现的次数，文件的每一行有两个字段，一个是单词，一个是次数，每个单词可能出现多次，求每个单词后面每个数组的和
下面是输入文件示例：
```
hello 1
world 2
hello 2
world 1
haha  4
demo  3
```
脚本如下：
```
#!/bin/awk
BEGIN{
    print "word", "count"
    count=0
}
{
    word[$1] = word[$1] + $2;
    count = count + $2
}
END{
   for (v in word){
      print v, word[v]
   }
   print "Total Count:", count
}
```
执行脚本
```
hadoop@hust-lh:~/BigData/Homework/Hw2/GraphLite/GraphLite-0.20/res$ awk -f count.awk word
word count
demo 3
haha 4
hello 3
world 3
Total Count: 13
```
其实和c语言的风格时一样的，程序控制语句支持if, for, while, do while。

##### awk内置变量
下面是awk的内置变量，这些变量时可以直接在脚本里使用的
> 
ARGC ：命令行参数个数
ARGV  ：命令行参数排列
ENVIRON   ：支持队列中系统环境变量的使用
FILENAME   ：awk浏览的文件名
FNR     ： 浏览文件的记录数
FS    ：设置输入域分隔符，等价于命令行 -F选项
NF  ：浏览记录的域的个数
NR：已读的记录数
OFS ：输出域分隔符
ORS ：输出记录分隔符
RS ：控制记录分隔符


#### 结语
`awk`是一个非常强大好用的文本处理工具，`awk`的更多使用方法参考[http://www.gnu.org/software/gawk/manual/gawk.html](http://www.gnu.org/software/gawk/manual/gawk.html)

