```toml
title = "Python学习速记--多线程"
date = "2015-07-21 12:26:25"
update_date = "2015-07-21 12:26:25"
author = "KDF5000"
thumb = ""
tags = ["Python", "学习笔记", "多线程"]
draft = false
```
####多线程
`Python`支持多线程，使用多线程也比较简单，直接使用`threading`模块即可
```
import time ,threading
 
def my_thread():
	print('新的线程%s'%threading.current_thread().name)
	n = 0
	while n<5:
		n = n + 1
		print('线程%s -->  %s'%(threading.current_thread().name,n))
		time.sleep(1)
	print('线程%s结束'% threading.current_thread().name)

print('线程 %s 正在运行' % threading.current_thread().name)
t = threading.Thread(target=my_thread, name='my_thread')
t.start()
t.join()
print('线程 %s 结束' % threading.current_thread().name)
```

<!--more-->

####锁
进程的资源线程是共享的，因此对于进程的变量，如果再线程中需要修改就需要加入锁的机制
```
import time,threading

#全局变量
balance = 0
lock = threading.Lock()
def change_balance(n):
	global balance
	balance = balance + n
	balance = balance - n

def run_thread(n):
	for i in range(100000):
		lock.acquire()
		try:
			change_balance(n)
		finally:
			lock.release()
		

if __name__ == '__main__':
	t1 = threading.Thread(target=run_thread, args=(5,))
	t2 = threading.Thread(target=run_thread, args=(8,))
	t1.start()
	t2.start()
	t1.join()
	t2.join()
	print(balance)
```

#####ThreadLocal
```
import threading

local_school = threading.local()

def process_student():
	std = local_school.student
	print('%s线程的学生是%s'%(threading.current_thread().name,std))

def process_thread(name):
	local_school.student = name
	process_student()

if __name__ == '__main__':
	t1 = threading.Thread(target=process_thread,args=('Tom',))
	t2 = threading.Thread(target=process_thread,args=('Alice',))
	t1.start()
	t2.start()
	t1.join()
	t2.join()
```
