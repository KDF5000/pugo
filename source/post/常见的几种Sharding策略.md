```toml
title = "常见的几种Sharding策略"
date = "2017-04-17 00:00:38"
update_date = "2017-04-17 00:00:38"
author = "KDF5000"
thumb = ""
tags = ["Sharding", "分布式系统"]
draft = false
```
### Range
range假设key有序，好处是临近的key经常在一起，比如共同前缀的key,可以很好的支持scan操作，hbase的region就是range策略。缺点是对压力较大的顺序写不太友好，比如日志类型的写入，一般日志的key都是和时间相关的，时间是单调递增的，因此写入的热点永远在最后一个region。

但是对于关系型的数据库，因为经常性的需要**表扫描或者索引扫描**，基本上都会使用range的shard策略。
![](http://static.zybuluo.com/zyytop/5h4vfs0g6t7y3609lbuslozw/屏幕快照_2016-10-13_下午4.40.22.png)
*图片来自PingCAP博客*

### Hash
Hash的策略是将key先做一个hash,然后得到的hash值作为sharding ID, 这样每一个Key的分布都是随机的，因此可以任务是均匀分布的，这样对于写压力较大的系统非常友好，同事随机的读因为会将读的压力均匀分布在不同的节点，因此也是非常友好的。但是对于需要scan的操作几乎是不可能的。
![](http://static.zybuluo.com/zyytop/8kaltq5ww337kgxq63asdbnz/屏幕快照_2016-10-19_下午6.14.37.png)
*图片来自PingCAP博客*
#### Round Robin
这种hash是最简单的其实也是上面描述的Hash的最朴素的实现方式，Hash首先一般是讲数据进行分片，如果分片后每一片直接对应一个物理节点，一般的实现是直接在hash函数设计的时候将物理节点个数考虑进去，比如对物理节点的个数取余。
$$H(key) = hash(key) \%  k$$
这样H(key)的值即为物理节点的编号，这种方式因为讲数据分片和分片与物理节点映射的功能合二为一，因此新增加或者减少一个物理节点re-hash的代价非常高。

对于这个问题有多种解决方案，一种就是想办法解除合二为一的哈希功能，虚拟桶就是这种思路，另一种是一致性哈希所采用的方法。

#### 虚拟通(Virtual Buckets)
虚拟桶的方法就是为了解决上面提到的Round Robin的哈希方法的扩展性问题，其实也是很简单的，就是讲hash的映射分为两部分，首先 根据key做个hash映射到一个虚拟桶(数据分片)，然后维护一个从虚拟桶到物理节点的映射表。当新增加一个节点的时候，将某些虚拟桶从原先的物理节点移动到新的物理节点，然后更新映射表即可，不需要重新对所有的key进行hash.
![sharding.png](@media/archive/blog/images/sharding.png)
#### 一致性哈希
Cassandra和Twemproxy都是使用了一致性哈希的方法，这个讲起来有点复杂，详情看[待定]()

#### Pre-sharding
Pre-sharding其实就是前面提到的虚拟桶的方法，只要经过两次hash都可以认为是pre-sharding的方式。

### 动态扩展
range在动态扩展方面可以通过分裂的方式将一个大的shard分裂成不同大小的shard然后做shard的迁移，但是对于hash来说就必须来进行re-hash，这样的代价是非常大的，比如添加一个物理节点，此时hash的模如果是3，则添加后变成了4, 对于已有系统的抖动非常大。 虽然一致性hash可以一定程度降低系统的抖动，但是并不能彻底避免。


----
**参考资料**
1. [基于 Raft 构建弹性伸缩的存储系统的一些实践](https://pingcap.com/blog-building-distributed-db-with-raft-zh)
2. [Redis集群的数据划分与扩容探讨](http://engineering.xueqiu.com/blog/2014/12/26/redis-capacity-planning/)
3. 《大数据日知录》第一章：数据分片与路由
