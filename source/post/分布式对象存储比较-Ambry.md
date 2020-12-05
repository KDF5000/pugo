```toml
title = "分布式对象存储比较-Ambry"
date = "2017-07-02 21:07:36"
update_date = "2017-07-02 21:07:36"
author = "KDF5000"
thumb = ""
tags = ["Distributed System", "Object Storage", "Ambry", "Haystack"]
draft = false
```
在前面的文章[对象存储比较-Haystack](http://kdf5000.com/2017/06/25/对象存储比较-Haystack/)中详细介绍了对象存储的经典之作Haystack，该系统合并小文件，元数据分开管理的思想到现在依然是很多系统的基石，想TFS, FastDFS多少都借用了这个思想。去年(2016)Linkedid在数据管理会议SIGMOD发表了自己内部使用的对象存储系统[Ambry](http://dprg.cs.uiuc.edu/docs/SIGMOD2016-a/ambry.pdf)，能发到 SIGMOD，想必也不是等闲之物，因此就拜读了一下。整体感觉相比Haystack考虑的比较全面，诸如数据的负载均衡，partition的异步同步，甚至数据的冷热全都考虑了进来，但是依然能看到Haystack的影子，只能说是在基础上更加的完善，可用性，维护性也更强了。

## 动机
在介绍Ambry之前，先看一下为什么要重新自己设计一套对象存储系统，Linkedin的工程师想必也都是有经验的工程师，他们在分析了一下自己的业务场景，已经现在的解决方案，最后得出了一下的结论：
* 存储各种各样大小的对象(图片，文档，音视频)
* 各种不同的负载
* 现有系统大多针对小文件(100k一下)的文件存储进行了优化，但是不能适应各种大小的对象
* 现有系统不能很好的解决，数据访问倾斜，新加入节点导致数据分配不均衡等问题

因此他们决定在前人的肩膀上，重新设计一套适合自己的解决方案，他们认为他们的系统应该有如下的特点：
* 低延迟，高吞吐
* 跨多个数据中心
* 扩展性要强
* 能够很好的解决负载均衡的问题

## 架构
![ambry-arch.png](@media/archive/blog/images/ambry-arch.png)
Ambry主要有Cluster Manager，Datanode，Frontend三个组件。在介绍每个组件负责的职责前，先介绍Ambry里的一个重要概念: partition.

**Partition:** Partition是Ambry里的最小的存储单元，是一个虚拟的概念，相当于haystack里的Logical Volume,  Twitter BlobStore里的Virtual Bucket，Partition在实现的时候实际上是一个大的文件，存放着成千上万的小文件，并且作为replication的最小单位。因为存在replication，所以每个parition对应着多个位于不同节点，甚至不同机架和数据中心的物理文件，每个文件的大小通常100G，尽可能大，但是不能太大，因为移动数据的时候可能会消耗太多的时间。

**Cluster Manager:** Cluster Manager负责管理整个集群的节点，磁盘等布局信息，跟踪partition的状态，以及映射的位置。不同的数据中心都有自己的cluster Manager，通过zookeeper保持状态的同步。

**Datanode:** Datanaode是一个实际的物理节点，管理多个磁盘，存储真实的数据，每个datanode有多个partition，并且每个parition会在不同的节点有相应的副本。同时Datanode为每个partition都维护一个常驻内存的索引，Journals和布隆过滤，索引后面会详细介绍。

**Frontend:** Frontend是直接客户端交互的组件，所有的请求都是发送给Frontend然后再有其转发到相应的Datanode节点，然后接受数据返回给客户端。Frontend主要有下面三个功能：
- 请求处理转发
- 安全检查。Frontend在接收到客户端的请求的时候会根据设置进行一点的权限之类的安全检查。
- 操作拦截，push到Kafaka。这个功能主要用于监控集群的请求信息，用于分析请求的模式等信息。

Frontend有一个重要的模块叫Router Library，主要用于请求的转发，以及执行相关的逻辑。该模块的的请求操作行为是基于规则的。用户可以设置，每个写入或者get特定partition的请求是写入或者读取几个datanode才算是成功。除了转发请求，该模块还负责将大对象分成固定大小的chunck，然后作为新的blob分发到不同的partition，同时使用一个额外的保存blob的编号等信息，保存在某个paitition，并将partition返回给客户端。Router Library同时还兼顾故障检测，请求代理的功能。

## 索引、Journal和Bloom Filter
### 索引
前面提到Datanode的主要作用就是存储不同的paitition，每个partition内保存了成千上万的小blob，每个blob在Partition的布局如下图：
![ambry-blob_layout.png](@media/archive/blog/images/ambry-blob_layout.png)
当增加一个blob的时候直接append到特定的partition即可。那么当获取一个blob的时候怎么才能够快速地找到相应的数据呢？为此Ambry为每一个Partition创建了一个索引，保存了每个Blob在paitition文件中的偏移位置，大小等信息，不过建索引不是什么新鲜的东西，Haystack也有。但是Linkedin的工程师并没有就此止步，他们为了节约成本，减少索引占用的内存，Ambry讲索引分成不同的segment，然后只有经常访问的blob，也就是热索引才会常驻内存，然后那些冷数据通过在内存维护一个布隆过滤器，当访问冷数据的时候现在布隆过滤器里查找是否存在其对应的索引段，然后再决定加载那段索引 ，从而在保证用少量内存的情况下依然能够有较高的访问速度。
![ambry-index.png](@media/archive/blog/images/ambry-index.png)
### replication和Journal
replication是保证数据可靠性的一个有效方法，通常在保证数据的一致性上最经典的做法就是使用Paxos等协议，但是在对象存储的场景下，直接向所有或者多数的副本节点发送同步或者异步的请求就可以了，大大的简化了实现难度，并且保证了延迟，如果使用一致性协议则会大大的增加延迟(个人的猜测，也许不是这个原因)。Ambry的replication发生在put操作时候，客户端向Frontend发送put请求，然后根据贪心的策略选择replication的位置，选择的时候根据下面两个原则：
- 不能再同一个Datanode
- 处于多个不同的数据中心

副本的数量有管理员自己配置，选择好replicatio的位置后，异步的像这些节点发送请求，但是会根据管理员的配置等待指定的k(1，多数，all)个节点成功返回才向客户端返回成功的响应。

因为replication的操作是异步的，因此失败的情况是可能发生的。为了保证不同副本的一致性，Ambry使用了一个replication算法，周期性的在后台执行，每个replica都作为一个master从其他节点拉取丢失的blob，为了加快检测的速度，使用journal维护每个节点上每个parititon的lastOffset，然后分下面两个阶段同步副本的blob
- 请求lastOffset后的所有blob id，然后过滤掉本地已经有的blob
- 请求missing的blob

![ambry-replication.png](@media/archive/blog/images/ambry-replication.png)

## 负载均衡
负载均衡是Ambry相对于Haystack的一个重要的不同点，针对业务请求负载可能的倾斜，大量大的blobs以及集群节点的变化导致的负载不均衡，Ambry提出了分别针对静态集群，以及动态添加节点时候两种情况提出了自己的均衡算法。
### 静态集群
针对静态集群，导致其负载均衡主要的原因是大量的大Blob导致数据分布的均衡，已经数据冷热不同从而导致访问的不均衡。对于大的blob，ambry将其分成固定大小(一般4M-8M)的chunk，分发到不同的partition；对于热门数据通过CDN缓存，从而可以减少冷热数据访问不均的问题。

### 动态集群
对于新增节点的情况，因为新增节点，那么其上面的partition的状态都是可写的，因为新写入的数据往往有更大可能性被访问，因此新节点的partition往往承担大量的读和写，所以导致访问严重不均。Ambry使用了一个rebanlance的算法，该算法定义了一个理想下的读写，只读，磁盘使用率，计算方法都是根据整个集群当前读写状态的，只读磁盘，以及所有使用的磁盘除以整个集群所有的磁盘数得到；然后将每个磁盘上RW, RO状态的partition放入一个pool，最后将pool里的partition分发给低于理想状态的磁盘。整个算法的伪代码如下:
![ambry-reblance.png](@media/archive/blog/images/ambry-reblance.png)

那么如何移动一个partition到另一个节点呢？Ambry分为三个阶段去安全的移动Partition
- 在新的datanode创键一个replica
- 同步旧的replica所有的Blob到新的replica，同步的时候如果有新的blob写入旧的replica，那么两个replica都要执行append操作
- 同步完成后，删除旧的replica，然后更新映射信息

## 操作
Ambry主要提供了Put, Get, Delete三种操作，这三种操作都是客户端直接和Frontend交互，然后Frontend分发请求到相应的partition对应的Datanode，最后将数据或者操作状态返回给客户端。

Delete和Put操作相同，同步写本地数据中心的replica，然后异步写不同数据中心的replica, 只有返回管理员设置的k个成功响应即可。后续通过前面描述的replication算法达到副本的一致性。Put操作在Frontend层的时候会随机选择Partition希尔如，并且声称一个PartitionId + UUID的fid作为blob的唯一id，返回给客户端。对于Delete操作会置delete flag位，并且在后续周期性的compact操作中将其删除。

Get操作请求在Router模块，会提取请求的fid找到对应partitionId，向对应的replica节点发送请求，然后根据设置的成功请求数判断是否成功。如果get时候replica没有完成，则proxy到其他已经存在的数据中心。

## 总结
前面介绍了Ambry的整体架构以及核心组件的实现，中间也多多少少和Haystack进行了对比。总的来说Ambry的设计是在Haystack基础上又前进了一步，主要是为了能够适应不同大小的对象，以及解决负载均衡的问题。整个系统印象最深的有三点：
- 将大对象分成多个固定大小的chunck，并用一个新的Blob记录了不同chunck的信息，虽然多消耗了一些空间，但是提高了get数据的并行度，从而能够减低延迟。
- 对Haystack存在的多个副本同步失败的问题，通过引入异步复制，按需设置一致性，后台定期执行replication算法实现最终一致性巧妙的解决了。
- 对于Haystack没有解决的负载均衡的问题，创新性的提出了负载 均衡算法，后台动态调整不同状态partition的分布，实现动态的reblance。
