```toml
title = "Python学习速记--多进程"
date = "2015-07-21 12:25:31"
update_date = "2015-07-21 12:25:31"
author = "KDF5000"
thumb = ""
tags = ["Python", "学习笔记", "多进程"]
draft = false
```
####多进程
`python`封装了操作系统底层的进程函数，`Linux`下使用`fork`开辟多进程，`windows`下使用`multiprocessing`模块使用多进程
#####`windows`下使用多进程
需要使用模块`multiproccessing`中的Process
```
from multiprocessing import Process
import os

def proc_run(name):
	print('我是子进程%s(%s)' %(name,os.getpid()))

if __name__ == '__main__':
	print('即将启动一个子进程...')
	child = Process(target=proc_run,args=('child',))
	print("开始启动子进程")
	child.start()
	child.join() #等待子进程结束后再继续往下运行
	print("子进程结束")
```

<!--more-->

#####进程池
`Python`也可以使用`multiprocessing`中的线程池来管理子进程
```
from multiprocessing import Pool
import os,time,random

def long_time_task(name):
	print('任务%d(%s) 开始执行' %(name,os.getpid()))
	start = time.time()
	time.sleep(random.random() * 3)
	end = time.time()
	print('任务%d(%s)执行了%.2f秒'%(name,os.getpid(),end-start))

if __name__ == '__main__':
	print('开始执行批量任务')
	p = Pool(4) #进程池只允许四个进程 默认是CPU的核数
	for i in range(5):# 五个进程
		p.apply_async(long_time_task,args=(i,))
	print('等待任务执行完毕')
	p.close() #关闭后不能再添加进程
	p.join() #等待所有子进程执行完毕
	print('任务执行完毕')
```

#####进程间的通信
`python`进程的通信使用`multiprocessing`模块提供的队列`(Queue)`和管道`(Pipes)`
```
from multiprocessing import Process,Queue
import os,time,random

#写进程
def write(q):
	for val in ['A','B','C']:
		print('写字母%s' %val)
		q.put(val)
		time.sleep(random.random())

#读进程
def read(q):
	while True:
		val = q.get(True)
		print('读字母%s'%val)

if __name__ == '__main__':
	#创建一个队列
	q = Queue()
	#开辟两个子进程
	pw = Process(target=write,args=(q,)) #写进程
	pr = Process(target=read,args=(q,))  #读进程

	pw.start()
	pr.start()

	pw.join()
	pr.terminate() #读进程是死循环 只能通过terminate结束
```
