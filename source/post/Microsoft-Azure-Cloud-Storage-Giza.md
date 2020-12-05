```toml
title = "Microsoft Azure Cloud Storage: Giza"
date = "2020-02-17 23:45:14"
update_date = "2020-02-17 23:45:14"
author = "KDF5000"
thumb = ""
tags = ["Distributed System", "Distributed Storage"]
draft = false
```

Microsoft Azure Cloud Storage Series:

[Windows Azure Storage: A Highly Available Cloud Storage Service with Strong](https://webcourse.cs.technion.ac.il/236802/Spring2018/ho/WCFiles/Azure_Cloud_Storage.pdf)

[Erasure Coding in Windows](https://www.cs.princeton.edu/courses/archive/spring13/cos598C/atc12-final181.pdf) 

[Giza: Erasure Coding Objects across Global Data Centers](http://true/)

# **Overview**

Giza是构建在Azure Blob Storage之上的一层服务，在此之前Azure BLOB Storage只支持单Datacenter(虽然可以通过异步复制在不同的Stamp之间同步数据，但只能是一个最终一致性的系统)，并且通过LRC的优化在保证可靠性和可用性的前提下将存储成本降到了足够低(1.3)，但是随着业务的发展，对存储的需求要求越来越高，对DC级别的故障容灾的需求越来越强烈，Giza正是为了解决这个问题，为了保证尽量利用现有的存储基础设施，微软选择了在原有存储系统之上构建多DC的存储系统，这就诞生了Giza：**一个强一致的，支持多版本，使用EC编码并且跨全球数据中心的对象存储系统**。

当前Azure在全球有38个不同的Region，数EB的数据存储量，为了保护用户的数据在磁盘、机器、机架故障的情况下数据依然完好无损，微软设计LRC[1]在保证数据高可用和高可靠的情况下极大的降低了存储成本。但是随着业务的发展，用户对数据安全性的要求越来越高，希望能够在地震、洪水等可能导致Region或者机房级别故障的情况下，数据依然是可以访问的，因此这就要求Azure需要将数据在其他Region也保存一份，这就导致了存储成功成倍的增加。使用EC在单机房非常有效的将存储成本将了下来，自然而然的如果能够将EC数据分布在多个Region或者Datacenter就能够在保证数据可靠性的情况下将低存储成本。但是将EC的多个分片数据分布在不同的Datacenter就会导致每次请求都要从不同的Datacenter获取数据，导致占用cross-dc的带宽占用并增加请求延迟，而且在Datacenter出现故障恢复后会产生大量的跨DC的数据修复。在Datacenter专线带宽有限且费用昂贵的情况下，这种跨Datacenter的EC方案并非最好的解决方案。

<!--more-->

**Design Goals**

1. 最下化请求延迟的前提下保证数据的强一致性
2. 尽量充分利用现有的基础设施去实现和部署

**Challenge**

Giza的一个目标之一就是最小化请求延迟，如果做Cross-DC的EC，那么必然会带来跨DC的读写请求，所以如果想要对跨DC的延迟进行优化，最好的方法就是尽量减少跨DC的请求。微软分析了他们内部负载特点，发现确实存在很大的优化空间，比如OneDrive的负载，对象很少更新，并且存在的更新往往都是发生在同一个数据中心，很少出现不同数据中心的更新冲突。同一个数据中心的更新就比较容易了，但是如果出现不同数据中心的更新，就可能导致数据不一致性，那么就需要一种方法去保证多个DC之间数据强一致性。所以要想在这种负载的情况下实现Giza有两个重要的挑战：

- 对于没有更新冲突的请求在一次跨DC网络下完成
- 确保冲突更新的强一致性

要想保证强一致性，传统的方案是指定一个主，所有的写入请求relay到主，这样就不存在写入冲突的问题，但是也就也为这所有来自非主DC的请求都会多一次跨DC的网络请求，这与系统的设计目标是相悖的。Giza则使用了FastPaxos和经典Paxos两种一致性算法，在没有冲突的时候使用FastPaxos尽量减少跨DC的请求，当出现更新冲突的情况下则回退到经典Paxos算法，确保数据的强一致。

**Summary**

Giza在现有云存储Azure的基础上实现了一套基于纠删码多版本的、全球强一致的分布式存储系统。当一个对象写入的时候，Giza将该对象EC后的不同Fragment（通常是2+1）分布在不通过的DC，并且将元数据全量同步到不同的DC。当前Giza已经在全球跨三个洲11个数据中心部署并稳定运行。

# **Giza overview**

**OneDrive Characteristics**

1. Large Objects Dominate
- 0.9%的存储空间是有小于4MB的对象占用，如果想要降低存储成本，**只用考虑大于等于4MB的对象，小于4MB的对象直接使用多副本的方案**，这样可以减少小对象进行EC带来的开销（如元数据的存储）
2. 对象降温很快
- 大多数据的读都是发生在对象创建的时候
- 47%的读在对象创建一天内发生，87%的读在对象创建一周内发生在，还有不到2%的读发生在对象创建一个月
3. **缓存可以减少跨dc请求**
- 没有缓存的情况下，读写请求比是2.3，其中跨dc读写比是1.15
- 数据写后缓存到本DC，缓存一天跨dc读写比可以降低到0.61，一周可以降到0.18，一个月降低到0.05
4. **并发很少，但是需要多版本**
- 57.96%的对象写入后三个月内不会有更新，40.88%会更新一次，1.16%更新多余2次，0.5的对象会出现并发更新
5. **删除并不少见，一年前的创建的对象被删除的空间占总空间的26.7%，因此回收被删除的空间很有必要**

**Giza tradeoff**

![Giza Tradeoffs](@media/azure/tradeoffs.png)

Giza采用k+1的方式对数据进行编码，单个DC使用论文[3]描述的LRC方法可以将存储成本降到1.3(12,2,2)，相比Geo-Replication的方法，Giza的k+1 EC在多机房的情况下既降低了存储成本同时也停工了更高的数据可靠性。当然也有存在一些不可避免的成本，Geo-Repl的方案下，每个对象的写也要多一次跨dc的写入，因此写入放大1X，Giza同样需要在本地写入1个分片后，跨DC写入k个数据，写入放大一样是1X。读取的话，Geo-Repl不需要跨DC的请求，在本dc就可以读到完整的数据，但是Giza的方法需要跨机房读取(k-1)个分片，跨DC读取占比(k-1)/k（实际上如果读请求发生在Parity所在机房，那么需要跨dc读取k个数据分片，也就是跨DC需要产生k个流量，当然也可以读取k-1个走EC恢复过程）。

实际上为了减少跨DC的流量，Facebook F4也提出了一种可行的方案，两个DC以本机房的Volume为单位，进行EC或者XOR，然后将Parity或者XOR结果放到第三个机房，这样能够保证每次读取对象都是在单个机房就可以完成，不过这种方案对删除不是很友好，F4也因此采用了外部DB的方案，记录数据的删除，底层存储空间不进行回收。前面分析过OneDrive的删除比例还是比较高的，因此回收删除数据可以节省很大的成本，因此Giza需要提供支持删除的方案，并且微软认为F4的方案如果想要回收删除空间，需要引入复杂的处理方式，使得系统变得很复杂，这违背了他们的设计原则。

# **设计实现**

![Giza Architecture](@media/azure/giza_arch.png)

上图是Giza的整体架构，可以看出Giza充分利用了现有的系统：对象存储和表格存储。这两个系统都是在单DC稳定运行多年，提供高可靠、高可用的服务。Giza将每个对象的写入，采用k+1的方式进行EC编码，每个对象被分成k个分片，经过编码后生成一个Parity分片，分别放在k+1的数据中心，如图中的{a，b,p}。每个对象的写入又分为数据操作和元数据操作，这两个操作是并行完成的，其中元数据包含了每个分片的唯一ID，以及存储位置，在每个dc全量复制一份。k+1个分片，每个dc放置一份，利用单个dc内的容灾能力保证数据的可靠性，也方便单个dc的存储进行单独的优化，Giza只作为一个无状态的服务，可以说这种设计充分利用了当前已有的基础服务，也减少系统内各个模块的耦合度，保证每个模块都可以按照自己的特点进行极致优化。

**技术挑战：**

1. 利用现有的单DC的Cloud table构建一个强一致的geo-replicated元数据存储
2. 联合优化数据访问和元数据访问，实现单次跨DC请求的读写操作
3. 高效适时的实现垃圾回收，从而回收删除的对象和元数据

**强一致的元数据**

前面提到Giza一个关键的挑战就是实现跨地域的元数据的强一致性，Giza选择了Paxos和FastPaxos一致性算法达到此目的。但是和传统的Paxos实现又有很大的不同，传统的Paxos的Acceptor是一个有状态的进程，可以通过自身的状态决定是否投票，Giza则是基于Azure Table利用其原子条件更新的特性实现Acceptor的逻辑。

Giza的一个很重要的优化目标是减少跨DC的请求，传统的Paxos的流程分为两个阶段：提议和提交，也就是说一次元数据的操作，需要两次跨dc的网络请求。为了减少跨DC的请求，Giza使用FastPaxos [3]将两个过程合并成一次网络请求，但是需要更多的Acceptor同意提议。

**元数据的布局**

![Metadata Layout](@media/azure/meta_layout.png)

Giza将每个对象的元数据在表格存储中存储一行，比较特殊的是这一行是一个可变列长的一行，其中每个版本的对象占用三列，每个版本包括是哪个记录: 当前看到的最大投票号、当前接受的最大投票号和当前接受的最大值（EC的schema、每个分片ID和所有分片所在的DC）。除此之外，Giza还维护了一个​Known committed versions​ 的集合，记录了当前已经被提交的版本。

**元数据写**

元数据在写之前会先查询当前DC的​known committed version​表，找出当前DC认为已经提交的最新版本，然后加一作为新的版本号。因为异常情况下，改DC所知道的已经提交的版本可能不是最新的版本，所以在拿着新版本进行提交的时候可能会失败，失败的时候Giza就可以知道当前应该使用的最新版本是多少。按照FastPaxos的流程，Giza会发送一个PreAccept的请求给所有的DC，每个DC收到请求后会进行原则更新操作，如果Giza收到一个Fast Quorum的成功返回后，则认为该元数据的更新成功，此时异步发送请求更新各个DC的​known committed version​表。

**元数据并发写**

Fast Paxos在失败的情况下，可能是因为存在并发写导致不能返回Fast Quorum的成功响应，此时会转入经典的Paxos流程。Giza首先选择一个可以区分的Ballot号，然后发送给各个DC，这个阶段称为Prepare阶段。每个DC接收到Prepare请求后，只有当Prepare中的Ballot号比当前表中的Highest ballot seen大时才返回成功，并且会把整个一行的数据返回。Giza在收到多数成功返回后，会选择一个值进行提交，选择哪个值进行提交的规则如下：

1. 在所有响应中选择最大的已经被Accepted的值，如果有的话
2. 选择最大的Pre-accepted值，
3. 1和2如果存在任何一个就选择这个值使用一个新的版本使用fastpaxos进行提交
4. 如果1和2都不存在，那么说明不存在冲突，直接使用当前的metadata进行提交

选择一个值进行提交后，每个dc的table只有在当前highest ballot seen和highest accepted ballot比要写入的小的时候才能写入成功。当多数写入成功的时候，Giza返回客户端写入成功，然后异步复制Commit版本到各个dc的table。

**元数据的读**

元数据的读主要是如何找到改对象已经被提交的版本，因为本地的​known committed version​可能没有包含已经被提交的最大版本，所以Giza需要读取多个dc的元数据的。通常Giza会读取多个DC的对应的元数据行和​known committed version​，当找到多数认可已经提交的版本，并且没有比该版本更高的版本的accpeted value，那么就返回该版本的元数据。但是如果发现存在高于该版本的accepted value，则会走Paxos流程确认这个更高的版本是否已经被提交。

**联合优化数据和元数据的操作**

Giza的读写操作包含两个路径，数据操作和元数据操作。正常的逻辑是先写数据成功后写元数据或者先读元数据然后读数据，完全串行的一个读写流程，并且还是跨DC进行的。为了降低读写延迟，Giza将读写数据和元数据的操作并行进行，这必然会带来一致性的问题，因此需要一些策略来保证读写的正确性。

1. 写操作

并行去写数据和元数据，Giza等待两者都返回的时候再返回客户端并且提交更新​known committed vesrion​，这个流程可以确保​known committed version​所有的版本的元数据和数据都是提交成功的。但是存在两种异常case：

- 数据写入成功，元数据写入失败
- 数据写入成功，元数据写入失败，相当于写入成功的数据成为了孤儿，这个时候对外部来说不会影响读取，因为对客户端来说这个数据是写入失败的，永远不会读取到。但是写入成功的数据需要进行清理，Giza会启动一个清除线程，标记当前版本为​no-op​，意味着孤儿的分片在元数据表里是不存在的，然后从各个数据中心删除孤儿分片
- 元数据写入成功，数据写入失败
- 元数据写入成功，以为着客户端来读取的时候是有可能读取到改版本的元数据，然后尝试去读取数据，但是实际上数据是不存在的，就会导致读取失败，所以在读取的时候需要考虑到这种情况

2. 读操作

Giza首先从本地读取​known committed version​找到该对象最新的版本，然后去读取数据，同时启动一个异步操作去验证改版本确实是最新的版本，如果验证失败那么就说明存在最新的版本，然后拿着最新的版本重新去读取数据。这样确实会可能存在多次读取数据的情况，但是这种情况发生的概率比较低，只有在大量并发的时候才会发生。 针对上面写操作的第二种情况，Giza发现数据读取失败的时候，会马上去查找之前一个​known committed version​的版本读取数据。

**删除和GC**

Giza的删除操作会作为一个更新操作处理，针对特定版本或者整个对象的删除，会修改元数据或者增加一个新的版本标识该对象被删除，这个过程 只涉及到元数据的更新，一旦更新完成便返回客户端删除成功。另外会有个GC服务回收删除对象的空间，回收的过程分为三步：

1. 从元数据表读取需要回收的对象
2. 从对象存储系统删除对应的分片
3. 从元数据行内删除对应的版本（特定列）

如果是整个对象的删除，则需要从元数据表中删除整行数据。因为是Giza是允许多点写入的，所以删除正常数据存在冲突的风险，比如一个DC对应的对象已经删除了，但是其他DC还没有删除，这个时候如果有个put请求，那么该DC会认为不存在这个对象，会生成一个最小版本的元数据进行写入，这个时候就会和其他没有删除元数据的DC的数据存在不一致。所以针对删除整个元数据的情况，Giza采用两阶段提交的方法去保证，首先标记所有DC要删除的元数据行为​confined​状态，不接受任何的读写请求，然后第二阶段进行删除。如果存在DC不可用，那么这个操作是不能成功，只有故障的DC恢复后才能继续进行。

# **故障恢复**

Giza是构建在Azure Blob storage和Azure Table Storage之上的，所以只用关心DC级别的容灾，单个DC的数据容灾交给Blob和Table服务，这让容灾变得会简单许多。

**短暂的DC故障**

对于短暂的DC界别故障，Giza运行降级的读写，会发生多余一次的跨DC网络请求。不过论文里主要讨论的都是元数据的读写，元数据在多个DC之间是多副本的，但是数据分片的读写并没有太多介绍，所以这里多个EC分片的多数写入读取成功可能也是走的FastPaxos或者Paxos的方法，支持多数写入成功。等故障DC故障恢复后，利用Paxos Learning的机制，将缺失的元数据和数据分片补回来，当然对于EC分片会产生EC重建，计算出缺失的分片数据。

**永久性故障**

极端情况下，DC可能存在永久性故障，意味着这个DC的数据直接放弃，以及这个DC不在接受任何的读写请求，Giza采用逻辑DC的概念，用户选择的DC只是一个逻辑上的DC，使用一个额外的映射表记录逻辑DC和物理DC的关系，当一个物理DC永久性故障的时候，只需要改变映射指向新的物理DC，然后启动恢复任务从其他DC补全缺失的元数据和数据分片即可。

# **性能评估**

![Performance Comparision](@media/azure/performance_cmp.png)

上表是Giza的配置，对比的对象是CockroachDB，因为CockroachDB不支持全球范围的复制，所以只部署了US-2-1的情况进行对比，并且每个DC部署三个实例。

Giza机器配置：16c，56G，1Gbps的虚拟机

1. Metadata延迟：Fast Paxos和Classic Paxos对比

![metadata latency](@media/azure/meta_latency.png)

2. Giza延迟

![Overall Latency](@media/azure/overall_latency.png)

对象大小为4MB的情况下，对于put操作，经典的Paxos和FastPaxos的中位数延迟分别是852ms和598ms，Giza的延迟是374ms左右，只比只传输数据的延迟多了30ms。get操作，Giza的延迟只有223ms。

3. 和CockroachDB的对比

在CockroachDB上实现Giza，因为CockroachDB不支持大对象，所以只测试了128KB的情况，put延迟的中位数CockroachDB有333ms，Giza不到100ms。get操作，cockroachDB延迟比Giza低了20%，因为CockroachDB是读本地盘，但是Giza是读Azure，所以要快一些。

4. 并发更新

![Concurrent Update](@media/azure/contention_update.png)

**参考文献**

1. [Giza: Erasure Coding Objects across Global Data Centers](https://www.usenix.org/sites/default/files/conference/protected-files/atc17_slides_chen.pdf)
2. [Fast Paxos](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/tr-2005-112.pdf)
3. [Erasure Coding in Windows](https://www.cs.princeton.edu/courses/archive/spring13/cos598C/atc12-final181.pdf) 
4. [Windows Azure Storage: A Highly Available Cloud Storage Service with Strong](https://webcourse.cs.technion.ac.il/236802/Spring2018/ho/WCFiles/Azure_Cloud_Storage.pdf)
