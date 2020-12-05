```toml
title = "Python学习速记--面向对象"
date = "2015-07-21 12:18:18"
update_date = "2015-07-21 12:18:18"
author = "KDF5000"
thumb = ""
tags = ["Python", "学习笔记"]
draft = false
```
python支持面向对象编程，定义一个类的方式如下:
```
class Student(object):
	pass
```
其中class是类的声明，Student是类名，后面括号里的object是该类继承的类，即父类
```
#创建一个类
stu_instance = Student('jack','male')
print(stu_instance.name)
stu_instance.print_info()

#添加私有变量
class Student2(object):
	def __init__(self,name,sex):
		self.__name = name
		self.__sex = sex
	def print_info(self):#self必须作为参数
		print('name:',self.__name,'sex:',self.__sex)
	
	def set_name(self,name):
		self.__name = name
	def get_name(self):
		return self__name

	def set_sex(self,sex):
		self.__sex = sex
	def get_sex(self):
		return self.__sex

stu2 = Student2('jack','male')
print(stu2._Student2__name) #可以这样访问私有变量
stu2.print_info()
print(stu2.get_sex())
```

<!--more-->
