```toml
title = "Python Rq的使用"
date = "2015-08-31 12:17:11"
update_date = "2015-08-31 12:17:11"
author = "KDF5000"
thumb = ""
tags = ["Python", "分布式", "爬虫"]
draft = false
```
`Rq`的介绍和安装可以参考[Ubuntu 14.04 下安装使用Python rq模块](http://kdf5000.github.io/2015/08/23/Ubuntu-14-04-%E4%B8%8B%E5%AE%89%E8%A3%85%E4%BD%BF%E7%94%A8Python-rq%E6%A8%A1%E5%9D%97/)，此文章详细的介绍了安装`Rq`的全部过程，文章最后给出了`Rq`官方文档的地址。

 `Rq`官方文档详细的介绍了`Queues`, `Workers`，`Results Jobs`等每个组件的功能和使用，但是并没有指出这几个组件之间到底有什么关系，如果联合起来使用，对于像我这种菜鸟级的码农，一开始还真有点摸不着头脑，思考一番，懂了之后才发现原来是那么的简单...下面是我自己对`Rq`的整体认识图（windows 画图画的有点粗糙望见谅）

![](@media/archive/img_rq_model.png)

用白话解释就是：有一个大工厂，这个工厂有个大仓库（Redis），有一个或者几个管理员（CreateJob），一些工人(Worker)。管理员负责产生任务，把任务放到仓库里特定的位置（特定的Queue），然后工人自己去队列去任务，默默的完成任务，完成成功后按照要求放到特定的位置（可能是仓库也可能是其他地方）。如果任务比较多或者想在短期内完成任务，那么工厂就可以招聘更多的工人去完成这些任务。

使用`Rq`模块最简单的就是只需要一个管理员，一个仓库，一个工人（工人都是一样的，可以复制很多出来）
####建仓库
仓库就是我们安装的时候使用的redis，具体建造过程参见[Ubuntu 14.04 下安装使用Python rq模块](http://kdf5000.github.io/2015/08/23/Ubuntu-14-04-%E4%B8%8B%E5%AE%89%E8%A3%85%E4%BD%BF%E7%94%A8Python-rq%E6%A8%A1%E5%9D%97/)中安装`Redis`一节

<!--more-->

####招管理员
仓库建好后，我们就需要一个或者几个管理员来生成不同的任务（Job），放到仓库里即可，`Job`的调度不用担心，`Rq`模块会自己处理，其实就是不同的工人自发的从仓库里特定的区域拿任务即可。
一个管理员的主要任务就是产生任务，如下：
```
#连接redis
redis_conn = Redis(host='192.168.0.108', port=6379)
q = Queue(connection=redis_conn, async=True)  # 设置async为False则入队后会自己执行 不用调用perform

with open("companies.json", 'r') as f:
    i = 0
    for line in f:
        job = q.enqueue(parse_company, line.strip())
        i += 1
        print i, ":", job.id
```
这个管理员的工作就是从文件`Companies.json`读取每一行内容，将每一行内容放到仓库(Redis)默认的位置（Queue），并且指定那一类的`worker`去完成这些任务，也就是`parse_company `，其实只是指定工人所要做的工作流程，并不是一个工人实体。`parse_company `的流程如下：
```
def parse_company(json_data):
    try:
        obj = json.loads(json_data)
        company_data = obj['Company']
        new_company_id = insert_company_info(company_data)  #插入公司数据到数据库
        if new_company_id is None:
            conn.rollback()
            with open('error.txt', 'a') as ferr:
                ferr.write(json_data)
            return None
        # 股东结构
        partners = company_data['Partners']
        for val in partners:
            id = insert_parter_info(new_company_id, val)  #插入股东信息到数据库
            if id is None:
                cursor.close()
                conn.rollback()
                conn.close()
                with open('error.txt', 'a') as ferr:
                    ferr.write(json_data)
                return None
        cursor.close()
        conn.commit()
        conn.close()
        print 'success!'
        with open('success.txt', 'a') as fsu:
            fsu.write(json_data)
        return True
```

负责这个任务的工人所要做的工作也比较简单，就是解析每个任务内容（一行json文本），然后插入到数据库中。

其实管理员的工作还是比较繁重的，既要将大的任务分解成小的任务，又要指定那一类工人（那个工艺流程）去做这些事，这也正是管理员工资比工人工资高的地方吧，虽然不出体力，但是脑力劳动还是比较强的。

####招聘工人
招聘工人其实比较简单了，当然需要付出的，在计算机世界就是要么多买一些计算机或者多开几个线程，然后还是培训这些工人，告诉他们负责那个工艺流程。

将上一节中指定的工艺流程的文件拷贝一份，放到需要完成任务的计算机上，当然在该计算机也要安装`Rq`模块，到此时一个工人的培训已经结束了（有点填鸭式教育的感觉），让工人开始工作只要一个指令即可。
```
#在parse_company文件所在目录下执行
$rqworker -u "redis://192.168.0.108:6379"  #-u后面的地址是仓库（redis）的地址
```

这个指令需要告诉工人去哪个仓库取任务，就这么简单，想要招几个工人就招几个工人

####验收产品
每个工人都是按照特定的工艺流程进行的，每个工艺流程指定了产品的输出位置，到相应的位置验收产品即可。

整个工厂的工艺流程见[这里](https://github.com/KDF5000/ParseCompany)





