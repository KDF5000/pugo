```toml
title = "Python学习速记--python函数相关"
date = "2015-07-21 09:01:02"
update_date = "2015-07-21 09:01:02"
author = "KDF5000"
thumb = ""
tags = ["Python", "学习笔记"]
draft = false
```
####python数据结构
* list = [1,2,3] 方法：list.append(ele),list.insert(pos,ele),list.pop(),list.pop(pos)
* tuple = (1,2)或(1,2,[3,4])  方法：tuble元素不可变
* dict = {'Tom':12,'Michael':23}  方法：dict.get(ele),dict.pop(ele)
* set = set([1,2,3])  方法：add, remobe

####函数
函数名其实就是指向一个函数对象的引用，完全可以把函数名赋给一个变量，相当于给这个函数起了一个“别名”
```
>>> a = abs #变量a指向abs函数
>>> a(-1) # 所以也可以通过a调用abs函数
1
```
python中，自己定义一个函数使用def语句，依次写出函数名、括号、括号中的参数和冒号:，然后，在缩进块中编写函数体，函数的返回值用`return`语句返回。

<!--more-->

```
#自己定义函数
def my_abs(x):
	if not isinstance(x,(int,float)):
		raise TypeError('参数类型不匹配')
	if x>0:
		return x
	else:
		return -x;
```
python 中的函数可以返回多个值，返回的return语句中直接用`,`隔开即可，接收变量同样`,`隔开，如下：
```
#返回多个值
import math # 导入math包 

def move(x, y, step, angle=0):
    nx = x + step * math.cos(angle)
    ny = y - step * math.sin(angle)
    return nx, ny
nx,ny = move(0,0,10)
print(nx,ny)
#其实返回的还是一个值，是一个tuple
r1 = move(0,0,10)
print(r1)
```
其实返回值是一个tuple
####函数参数
python支持默认参数，放在后面，而且默认参数必须指向不可对象，否则该默认参数的值在多次调用改变参数变量的值是默认参数的值也会放生变化。如下
```
def add_end(L=[]):
	L.append('end')
	return L;
print(add_end()) # L的默认值是[]
print(add_end()) # L的默认值是['end'],所以结果是['end','end']
#修改如下，将默认参数设置为不可变对象，None
def add_end_2(L=None):
	if L is None:
		L = []
	L.append('end')
	return L
print(add_end_2()) #L的默认值是None
print(add_end_2()) #L的默认值是None
```

#####可变参数
pytho如何实现可变参数呢？
* 可以将list或tuple类型作为参数，这样就实现了输入可变参数，但是调用的时候必须先组装一个list或者tuple
* 以`*param`的形式作为参数，这样调用的时候有几个参数就写几个参数
对比如下：
```
#list形式的可变参数
def fun_params(list):
	sum = 0
	for num in list:
		sum = sum + num
	return sum
print(fun_params([1,2,3])) #参数必须是个list
# “指针”形式的可变参数
def fun_params_2(*numbers):
	sum = 0
	for num in numbers:
		sum = sum + num
	return sum
print(fun_params_2(1,2,3)) #参数不用组装list,直接有个参数就写几个参数,其实内部参数numbers接收的是一个tuple
```
对于已经有的list，如何调用第二种像是的可变参数函数呢？如下代码
```
list = [1,2,3]
print(fun_params_2(*list))#只用在list或者tuple变量前加一个*即可
```

*总结*：对于可变参数，允许输入0个或者任意多个参数，内部都是组装成一个tuple


#####关键字参数
可变参数允许你传入0个或任意个参数，这些可变参数在函数调用时自动组装为一个tuple。而关键字参数允许你传入0个或任意个含参数名的参数，这些关键字参数在函数内部自动组装为一个dict
```
#关键字参数
def dict_fun(name,age,**kw):
	print('name:',name,'age:',age,'other:',kw)
dict_fun('jack',23,city='beijing',sex='male')#以键值对的方式调用
#dict形式的调用
dict_demo = {'city':'bijing','sex':'male'}
dict_fun('jack',23,**dict_demo)#dict形式的调用,dict变量前加两个**
```

#####明名关键字参数
和关键字参数不同的是，命名关键字参数制定了关键字参数的个数和名字，如下
```
#命名关键字参数
def kw_fun(name,age,*,city='Beijing',job):
	print(name,age,city,job)
kw_fun('tom',23,job="engineer")
```
还不知道什么时候会使用，，，

#####组合参数
在Python中定义函数，可以用必选参数、默认参数、可变参数、关键字参数和命名关键字参数，这5种参数都可以组合使用，除了可变参数无法和命名关键字参数混合。但是请注意，参数定义的顺序必须是：必选参数、默认参数、可变参数/命名关键字参数和关键字参数。

```

def f1(a, b, c=0, *args, **kw):
    print('a =', a, 'b =', b, 'c =', c, 'args =', args, 'kw =', kw)
f1(1,2,3,4,5,name='jack')

def f2(a, b, c=0, *, d,e, **kw):
    print(a,  b,c, d, e,kw)

f2(1,2,3,d=4,e=5,name='jack')


args = [1,2,3,4]
ext = {'name':'jack'}
f1(*args,**ext)

```
###总结

`Python`的函数具有非常灵活的参数形态，既可以实现简单的调用，又可以传入非常复杂的参数。
默认参数一定要用不可变对象，如果是可变对象，程序运行时会有逻辑错误！
要注意定义可变参数和关键字参数的语法：
`*args`是可变参数，`args`接收的是一个`tuple`；
`**kw`是关键字参数，`kw`接收的是一个`dict`。
以及调用函数时如何传入可变参数和关键字参数的语法：
可变参数既可以直接传入：`func(1, 2, 3)`，又可以先组装`list`或`tuple`，再通过`*args`传入：`func(*(1, 2, 3))`；
关键字参数既可以直接传入：`func(a=1, b=2)`，又可以先组装`dict`，再通过`**kw`传入：`func(**{'a': 1, 'b': 2})`。
使用`*args`和`**kw`是`Python`的习惯写法，当然也可以用其他参数名，但最好使用习惯用法。
命名的关键字参数是为了限制调用者可以传入的参数名，同时可以提供默认值。
定义命名的关键字参数不要忘了写分隔符*，否则定义的将是位置参数。


