```toml
title = "面向HBase的内存key-value缓存的实现"
date = "2016-06-19 09:54:46"
update_date = "2016-06-19 09:54:46"
author = "KDF5000"
thumb = ""
tags = ["HBase", "Cache", "Redis"]
draft = false
```
### 0x01 背景
之所以要实现这个缓存主要原因如下(但是由于不是实际业务场景需求，所以可能不太准确，也可能不存在这个需求):
* 非结构化数据的爆炸式增长
* 处理速度的要求越来越高
* HBase是面向硬盘的
* 内存容量越来越大
* 热点数据可以在内存放下

### 0x02 设计方案
通常的要实现缓存，主要是在有两个大方向实现，一个是在客户端实现，另一个时在服务端实现
* 客户端实现
 - 修改Hbase Client的源码,在Put, Get等关键操作的地方加入缓存机制
 - 在client端设计一种缓存服务层,并实现一个分布式的Key-value缓存系统, 对Hbase Client进行重新封装
* 服务端实现
 - 修改HBase Server端的源码,在Put, Get等关键操作的地方加入缓存机制
 - 在服务端添加一层代理服务,解析所有client的请求,对Put, Get等关键

两种大方向的第一种方案，都是直接修改hbase的源代码，直接修改源码性能可能会更好一些，但是修改源码，会过于依赖Hbase的版本，对于每一个base的版本更新都可能要重新查看源码，重新修改，除此之外1)中的方案一，由于是在本地客户端进行的缓存，所以没有实现分布式的缓存，因此可能存在缓存命中率低，缓存数据不一致的情况。

客户端实现中的方案二，进行重新地封装，不对源码进行修改，同时使用分布式的缓存，既可以提高缓存的命中率，同时又解决的对hbase过度依赖的问题，但是可能会降低性能。

服务端实现中的方案二，通过设计一个缓存代理服务，同样可以解决对hbase的过度依赖，降低了整个系统的耦合性，但是实现一个代理服务并非那么简单，而且需要对hbase的通信机制以及相关协议比较了解。

通过上面的分析比较，考虑到可行性和技术背景，客户端实现的方案二是最合适的方案，而且其性能和可维护性，扩展性相对来说都是比较好的。分布式缓存系统我们使用Redis实现，并且在client本地实现一个LocalCache，尽可能的提高缓存的命中率，减少通信造成的时间延迟（局部性假设，同一个客户端最有可能访问它put的数据）。

<!--more-->

### 0x03 实现
系统的整体架构图如下
![Alt text](@media/archive/blog/image/hydracache/0.png)
按照0x02的方案，要实现的是图中的client，对原生的hbase client进行封装，包括redis主要实现了如下功能
![Alt text](@media/archive/blog/image/hydracache/Functions.png)
主要包括redis集群的实现，以及redis client的封装和localcache的实现。
* Redis集群部署：可以快速的搭建一个redis集群，用于HydraCache缓存系统的分布式缓存。
* Redis缓存：使用Redis缓存关键数据，提高系统的读取速度。
* LocalCache: 一个本地的缓存系统，此处有一个局部性假设，认为同一个客户端最有可能访问它put的数据，实现常用的LRU和LFU等缓存策略。

#### HBase集群的搭建
本系统由于处于实验阶段，并且没有真是的分布式环境，所以使用Docker在本机大家一个分布式HBase环境。Docker 是一个开源的应用容器引擎，让开发者可以打包他们的应用以及依赖包到一个可移植的容器中，然后发布到任何流行的 Linux 机器上，也可以实现虚拟化。容器是完全使用沙箱机制，相互之间不会有任何接口。相比传统的虚拟机，Docker更加的节省资源，在一台普通的机器上启动多个容器，基本上没有压力，因此在单机上使用Docker搭建HBase分布式集群可以比较真实地模拟真实的分布式集群。
本系统实现的基于Docker的分布式集群，具有下面的特性：
* 使用serf和dnsmasq 作为集群节点管理和dns解析
* 可以自定义集群hadoop和hbase的配置，配置完后只需重新build镜像即可
* ssh远程登录集群节点容器

具体的安装部署使用方法，[使用Docker单机搭建Hadoop完全分布式环境](http://kdf5000.github.io/2016/05/16/%E4%BD%BF%E7%94%A8Docker%E5%8D%95%E6%9C%BA%E6%90%AD%E5%BB%BAHadoop%E5%AE%8C%E5%85%A8%E5%88%86%E5%B8%83%E5%BC%8F%E7%8E%AF%E5%A2%83/)

#### Redis集群的搭建
Redis是一个高性能的key-value数据库，由于其数据是放在内存中的，因此经常被用做缓存系统，提供系统的响应速度，不过与像memcached这样的内存缓存系统不同的是，redis会周期性地把数据写入磁盘。
Redis 3.0后开始支持集群的部署，整个集群的架构如图
![Alt text](@media/archive/blog/image/hydracache/redis-cluster.jpg)
Redis 集群中内置了 16384 个哈希槽，当需要在 Redis 集群中放置一个 key-value 时，redis 先对 key 使用 CRC16 算法算出一个结果，然后把结果对 16384 求余数，这样每个 key 都会对应一个编号在 0-16383 之间的哈希槽，redis 会根据节点数量大致均等的将哈希槽映射到不同的节点。
    使用哈希槽的好处就在于可以方便的添加或移除节点。
* 当需要增加节点时，只需要把其他节点的某些哈希槽挪到新节点就可以了；
* 当需要移除节点时，只需要把移除节点上的哈希槽挪到其他节点就行了；

详细的使用说明可以参考我实现的一个快速搭建分布式redis的解决方案[redis集群的搭建](https://github.com/KDF5000/redis-cluster)

#### Client的实现
　　　Client主要对HBase Client进行封装，结合Redis实现缓存机制，并且实现一个LocalCache的功能。通过Redis缓存，可以实现多个客户端共享缓存的数据，缩短响应时间，LocalCache提高了同一个客户端读取近期读取的数据的响应速度，对于有些场景下的应用，可以减少通信时间从而减少了响应时间。
　　　Redis和HBase封装的类如图
　　　![Alt text](@media/archive/blog/image/hydracache/class.png)
由图可知，Client大致分为三个部分，HydraCacheClientImpl，Cache(LocalCache)和RedisCluster.
HydraCacheImpl负责是Client对外的核心接口，调用Cache和RedisCluster控制整个缓存策略。一共有四中模式，不使用缓存，只使用Redis, 只使用LocalCache, Redis和LocalCache都使用。
* 不使用缓存模式下，内部其实就是HBase Client的一个简单调用；
* 只使用Redis模式，在读取指定的key时，会先查看redis是否已经缓存了该数据，如果存在则直接读出返回客户端，如果不存在则去HBase读取数据并加入到Redis缓存里；
* 只使用LocalCache模式，用户可以在初始化的时候选择使用LRU(Least Recently Used)和LFU(Least Frequently Used)中的任何一种淘汰策略，当读取指定的key时，会先判断本地缓存中是否存在这个数据，如果存在则直接返回，否则读取HBase并存入缓存中；
* 使用Redis和LocalCache模式，在读取数据的时候会先判断本地缓存中是否存在对应的数据，如果存在则直接返回，否则读取redis判断是否存在相应数据，如果存在则直接返回，否则再从HBase读取然后分别存入本地缓存和redis里.

缓存策略流程图如下
![Alt text](@media/archive/blog/image/hydracache/cacheprocess.png)
HydraCacheImpl实现数据的缓存主要是在进行get操作中进行，如果缓存中没有命中，则读取hbase，然后对数据进行缓存。核心代码(省去异常判断)如下:
```
public String get(String tableName, String rowKey, String family, String columnName, int expireTime) {
		//先判断缓存里有没有
		String key = tableName+"_"+rowKey+"_"+family+"_"+columnName;
		String valString = getDataFromCache(key);
		if(valString != null){
			return valString;
		}
		Table table = null;
		Connection connection = null;
		connection = ConnectionFactory.createConnection(HydraCacheClientImpl.conf);
		table = connection.getTable(TableName.valueOf(tableName));
		Get g = new Get(rowKey.getBytes());
		g.addColumn(family.getBytes(), columnName.getBytes());
		Result result = table.get(g);
		byte[] bytes = result.getValue(family.getBytes(), columnName.getBytes());
		String valueStr = new String(bytes);
		//set cache
		if(valueStr != null){
			this.setData2Cache(key, valueStr, expireTime);
		}
		return valueStr;
}
private String getDataFromCache(String key){
   		String val = null;
		if(this.localCacheOn && this.localCache != null){
			val = this.localCache.get(key);
			if(val !=null){
				return val;
			}
		}
		if(this.cacheOn == true && this.redisCluster != null){
			val = this.redisCluster.get(key);
			if(val !=null){
　　　　		return val;
			}
		}
		return val;
}
private void setData2Cache(String key, String val, int expire){
		if(this.localCacheOn == true){
			this.localCache.set(key, val, expire*1000);//将s转换为ms
		}
		if(this.cacheOn == true && this.redisCluster != null){
　　this.redisCluster.set(key, val, expire);
		}
}
```
另外一个在缓存系统中重要的问题就是淘汰策略的指定，什么时候进行缓存过期的清除以及缓存达到限制大小是淘汰哪些缓存。

Redis提供了以下几种策略拱用户选择: 
**noenviction**：不清除数据，只是返回错误，这样会导致浪费掉更多的内存，对大多数写命令（DEL 命令和其他的少数命令例外）
**allkeys-lru**：从所有的数据集（server.db[i].dict）中挑选最近最少使用的数据淘汰，以供新数据使用
**volatile-lru**：从已设置过期时间的数据集（server.db[i].expires）中挑选最近最少使用的数据淘汰，以供新数据使用
**allkeys-random**：从所有数据集（server.db[i].dict）中任意选择数据淘汰，以供新数据使用
**volatile-random**：从已设置过期时间的数据集（server.db[i].expires）中任意选择数据淘汰，以供新数据使用
**volatile-ttl**：从已设置过期时间的数据集（server.db[i].expires）中挑选将要过期的数据淘汰，以供新数据使用

LocalCache实现了LRU和LFU两种缓存策略，用户可以根据自己的业务场景进行选择。
缓存失效的时候怎么办？通常主要两种方法，一种是消极的方法，在主键被访问时如果发现它已经失效，那么就删除它；另一种时积极的方法，周期性地从设置了失效时间的主键中选择一部分失效的主键删除。Redis对两种方法都有实现，LocaCache则采用消极的策略，在get请求缓存数据的时候，判断数据是否过期，如果过期则将其删除.

### 0x04 测试
#### 测试环境
　　笔记本型号：lenovo y470
　　操作系统：Ubuntu 14.04
　　内存：DDR3８G 
　　网卡：1000Mbps
　　Hbase分布式环境(Docker)：三个节点，一个master，２个slave
　　Redis集群：三个master和每个master一个slave
#### 测试结果
Hbase集群使用docker搭建，运行在本地机器上，启动三个节点，一个master，两个slave。Redis集群同样运行在本地机器上，使用不同的端口代表一个redis实例，一共启动六个，三个作为master，三个分别作为master的slave。

为了测试hbase原声client和hydraCache的性能，我们开发了一个HBench用于产生数据集和负载，写入hbase的数据集是一个顺序序列，负载读取用的数据集随机生成。主要测试了get操作的性能，测试的模式涵盖了第五节中的四种模式。在数据集规模的选择上，分别进行了四组规模的数据集，分别是300, 600,900和1200条数据。测试的结果如图
![Alt text](@media/archive/blog/image/hydracache/Size-Time_2.png)

图中每一组第第一列代表使用原生hbase客户端读取相应条数数据的时间，第二列是单独使用redis缓存所需的时间，第三列是使用LocalCache所需时间，第四列是同时使用redis和LocalCache所消耗的时间。从图中可以看出，使用缓存的响应时间明显要小于不适用缓存，使用缓存的情况下，仅仅使用LocalCache的时间是最短的，次之是同时使用redis和localCache，然后是单独使用redis。这与预期的结果还是比较一致的，使用缓存的情况下，由于在第二次读取某些key值的时候会从缓存中读取，这比从hbase中读取要快的多，而从localcache中读取数据因为减少了网络通信通信所以时间要更少一些。而同时使用redis和localcache的时间介于两者之间的原因应该是，读取数据的时候大部分情况下是可以从localcache中读取的，但是在设置缓存的时候要保存数据到redis，因此增加了操作的时间，故总的响应时间要多于单独使用localcache但是要比单独使用redis稍微好那么一点，不过并不是很明显。
 
对于缓存模式的选择，要根据实际的业务场景进行分析择取，比如有些应用，同一个客户端存取的数据很少在进行读取，但是可能读取其他客户端添加的数据，这个时候使用redis既可以满足，就能达到很好的效果；而另一些客户端则读取的数据大部分都是前面自己添加的数据，这样的业务场景就很适合使用localcache。

----
**相关项目:**
HydraCache: [https://github.com/KDF5000/HydraHbaseCache]()
Docker分布式Hbase: [https://github.com/KDF5000/hydra-hadoop](https://github.com/KDF5000/hydra-hadoop)
分布式Redis集群：[https://github.com/KDF5000/redis-cluster](https://github.com/KDF5000/redis-cluster)
性能测试：[https://github.com/KDF5000/HBench](https://github.com/KDF5000/HBench)

