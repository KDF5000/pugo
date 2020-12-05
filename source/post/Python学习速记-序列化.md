```toml
title = "Python学习速记--序列化"
date = "2015-07-21 12:24:38"
update_date = "2015-07-21 12:24:38"
author = "KDF5000"
thumb = ""
tags = ["Python", "学习笔记"]
draft = false
```
####使用`pickle`模块序列化
```
#pickle序列化对象
import pickle
d = {'name':'tom','age':22}
print(pickle.dumps(d))
f = open('seria.txt','wb')
pickle.dump(d,f)  #序列化到文件
f.close()
#反序列化
f = open('seria.txt','rb')
d2 = pickle.load(f)
f.close()
print(d2)
```

<!--more-->

####使用`JSON`序列化
#####序列化`dict`
```
import json
#使用`json`序列化
d = {'name':'Tom','age':23}
f = open('json.txt','w')
json.dump(d,f) #序列化到文件中
f.close()
# 从文件中加载
f = open('json.txt','r')
d2 = json.load(f)
print(d2)

# 序列化到字符串 
str = json.dumps(d) #序列化为字符串
print(str)
# 反序列化
d3 = json.loads(str)
print(d3)
```
#####序列化类
```
# 序列化类
class Student(object):
	def __init__(self,name,age):
		self.name = name
		self.age = age

#将student转化为dict
def student2dict(std):
	return {
		'name':std.name,
		'age':std.age
	}

#将dict转化为student
def dict2student(d):
	return Student(d['name'],d['age'])

stu = Student('kdf',23)
str = json.dumps(stu,default=student2dict) #指定序列化的函数
#str = json.dumps(stu,default=lambda obj: obj.__dict__) #通用的序列化，但是有的class没有__dict__变量
print(str)
#反序列化
stu2 = json.loads(str,object_hook=dict2student)
print(stu2)
```
