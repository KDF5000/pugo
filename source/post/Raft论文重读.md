```toml
title = "Raft论文重读"
date = "2017-07-13 17:16:40"
update_date = "2017-07-13 17:16:40"
author = "KDF5000"
thumb = ""
tags = ["分布式系统", "一致性算法", "Raft"]
draft = false
```
Raft是一个管理日志副本一致性的算法。相比Paxos结果一样，并且一样高效，但是理解起来更加的容易。Raft将一致性的主要元素分离开来，比如leader选举，log 复制，安全等。同时，也提供了一个新的机制实现cluster membership改变，其使用多数的原则来保证安全性。

<!--more-->

### 一致性算法
一致性算法的意义就是保证一致性的一组机器能够在其部分成员出现故障的时候依然能够存活下来(提供服务)。在Raft之前所有的一致性算法可以认为都是Paxos或者其变种的实现，但是Paxos难以理解，其工程实现往往都需要改变Paxos的架构。因此为了便于理解，以及工程上的实现，开发了Raft算法，该算法使用一些技术从而使其在保证正确性和性能的前提下，更加的容易理解和实现，这些技术主要有过程分解(领导选举，日志复制，安全)和状态空间减少。

相比其他一致性算法，Raft有几个明显的特点:
- strong leader: Raft使用很强的leadership，比所有的Log entries都是通过leader发送到所有其他的服务器。这样简化了管理也更加容易理解。
- leader election：采用随机化的定时器去选举leader，简单快速的解决了leader选举中的冲突。
- 成员变更: 使用一个joint consensus机制，保证在变更成员的时候依然能够正常的操作。

### Replicated state machines
**(Replicated State Machine)状态机**：一致性group的节点的某个时刻的状态(比如数据库里x=1,y=1是一个状态)转移可以看成自动机里的一个状态，所以叫状态机。
**Replicated Log**: 包含了来自客户端的用于执行在状态机上的命令(操作)。

replicate group里的每一个节点都包含一个**相同序列的**log,也就是相同序列的操作，因此每一个节点的state machine都是执行相同序列的操作，所以其结果也是相同。

一致性算法的主要目的就是保证每个节点上log的一致，其核心功能就是接收客户端的命令，并添加到log里，同时一致性模块与其他节点通信，保证他们的log都是**相同顺序的序列**，即使一些节点出现了故障，当log成功的复制到一定数量(通常是多数)的节点后，每个节点的状态机就执行该条命令并将结果返回给客户端。

总的来说一致性算法有下面一下特性:
- 安全。这里是指即使再机器出现故障(比如网络延迟，隔离，丢包，重复，重排等情况)的情况下绝不会返回错误的结果
- 可用性。在多数机器可用的情况下，整个集群依然总是可以对外服务。出现故障的机器后续依然可以通过持久化的数据恢复并且重新加入集群
- 不依赖时钟。一致性算法不依赖时钟就能保证log的一致性
- Single-round RPC. 通常情况下，一个命令只要集群中多数机器响应，就认为是成功了。即使一小部分机器出现了故障也不会影响整个系统的性能，而这个过程通常只需要一轮的RPC调用。

### Paxos的缺陷
Paxos几乎成为了分布式一致性算法的代名词，基本成为了一致性算法的教学典范，以及工业实现的参考。Paxos分为单命令的(single decree)*single-decree Paxos*和多命令的*multi-Paxos*. Paxos的正确性和性能都已经得到理论上的证明。

但是Paxos有两个主要的缺陷:
- 难以理解
- 工程实现困难

而且Paxos的架构实际系统的实现中并不实用，比如其收集日志的一个集合然后将其重排成有序的日志的做法并没有任何好处，相反直接append到一个有序的log更加的高效。还有其实用同等地位的点对点，没有leader的方式，对于一个确定一个决策没有问题，但是对于一些序列决策就存在一定的问题，虽然其最后建议使用一个弱的leader，但是还是显得很复杂。

### Raft如何做到可理解性？
为了是一致性算法更加容易理解，Raft主要使用了一下两个技术:
- 问题分解(Decomposition)
将一致性问题分解成几个容易理解和实现的子问题，Raft讲算法分为：领导选举(Leader Election), 日志复制(Log Replication)，安全(Safety)和成员变更(Membership Change)四个部分
- 减少状态空间
尽可能的减少需要考虑的状态和不确定性，让系统看起来更加的清晰。和Paxos一个很不同的地方就是Raft不允许log存储漏洞(Holes)，并且限制了不同节点之间log不一致的情况。通过减少状态空间和一些限制，大大地增加了算法的可理解性。除此之外，Raft使用随机的方式简化Raft的领导选举，这样虽然增加了算法的不确定性，但是随机化通过对于所有的状态都是用同样的处理方式可以简化状态空间。

### Raft一致性算法
Raft算法的一个核心思路是在一组机器中间选出一个Leader, 然后又Leader去管理日志，通过选举一个leader大大地简化了管理的复杂度，客户端需要发送新的命令，通过leader将log复制到所有其他节点，当大多数节点都复制了log后，leader通知所有节点，可以将日志包含的命令应用在本地的状态机。基于leader机制，Raft讲算法分为三个部分:
- leader选举
当一个集群中没有leadr或者存在的leader出现故障时重新选举一个leader
- 日志复制
当前的leader接受客户端发送的日志(命令)，然后复制到集群的其他节点。
- 安全(Safety)
安全主要指状态机的安全属性。

#### Raft基本概念
**Role:** 在一个Raft集群中(通常五个节点，容忍两个节点故障), 每个节点或者说(Raft实例)都有三种角色，Follwer, Candidate, Leader，每个实例初始化都是Follower状态，设置一个定时器，当一定时间内没有发现leader，则会发起leader选举的请求，其状态变为candidate，当其他节点投票满足一定的条件后成为Leader, 其他所有的节点都转变或者维持Follower状态，具体的状态装换图如下:
![state transition.png](@media/archive/blog/images/state_transition.png?imageView/0/w/500/)

**Trem:** Raft算法将时间分成不同长度的term，每个一个term都是从选举一个leader开始，换句话说就是质只要出现新的leader选举那么就是进入了新的term，一旦选出leader后，该term后续的所有时间都有该leader负责与客户端交互，管理日志的复制。Term更像是一个逻辑时钟，每个server都会维护一个currentTerm，这么在server通信的过程中便能识别出每个server的term是高于还是低于自己的term，从而进行状态的转换和term的更新。

#### Leader选举
Raft使用心跳机制触发leader选举，所有服务器启动的时候都是follower，规定在一定时间间隔内如果能够收到leader的心跳或者candidate的投票消息，则一直保持follower，否则触发leader选举过程。

leader选举有两个阶段：首先增加currentTerm，转为candidate状态，然后**并行**的向其他服务器发送投票RPC(RequestVote)。candidate状态结束的条件有如下三个:
- **赢得了选举，成为leader**
如果candidate状态的节点收到了多数节点的投票那么就会赢的选举，成为leader, 然后向其他节点发送心跳，告诉他们自己是leader。每一个节点在投票的时候，至多投一次，并且按照先到先得的原则，同时请求投票的节点的term必须大于等于自身的term。**多数的规则**能够保证再一个选举周期保证最多一个节点赢得选举成为leader。
- **收到了自称leader的节点的心跳**
如果candidate状态的节点在等待投票的过程中收到了某个节点的心跳自称自己是leader，如果心跳里包含的term大于等于自己的currentTerm，那么就说明该leader是合法的，自己转为follower，否则拒绝RPC,并返回自己的term.
- **既没有赢得选举，也没有失败(没有其他leader产生)**
如果同一个时刻有多个candidate状态的机器，那么就会产生votes split，这种情况就不会满足多数的规则，所以不会产生leader，这个时候每个candidate增加当前term，重置election timeout，开始新一轮的选举。这个时候为了防止每个candidate的election timeout相同导致无休止的选举失败，Raft采用了一个简单但是非常有效的方法，**随机生成**election timeout，通常是一个范围比如150ms~300ms。这个**随机化**也是Raft的一个重要的技术点，很多地方都用到了随机化从而简化了系统同时，也保证了正确性。

为了防止不包含之前leader已经commit的log entry，在投票的时候有一个限制，candidate在发送RequestVote RPC的时候会附带其最后一个log的index和term，只有candidate的term大于等于(等于时候比较index)自身的term才投票给candidate，否则拒绝投票。这样能够保证按照这种规则投票 选出的leader包含了所有之前leader已经committed的log entry。

#### Log Replication
当一个raft实例当选为leader后，该实例就开始负责与客户端交互，并将客户端的请求作为一个新的log entry保存到本地logs内，然后并行的发给其他的follower，如果收到多数follower保存该entry成功的相应，那么就commit该log entry，并apply到本地的state machine, 其他follower也会在~~收到commit的请求后或者在~~后续的一致性检查的时候apply改entry到本地的state machine，最终实现所有实例的一致性。

正常的情况下，按照上述的流程是不会出错的。但是往往follower, candidate, leader都有可能出现宕机或者其他故障导致log的不一致，如下图所示，leader和follower都有可能缺少log或者有多余的Log.
![Log inconsistent.png](@media/archive/blog/images/Log_inconsistent.png?imageView/0/w/500/)
最上面的log是当前的leader，a和b两个follower的log属于丢失log，c，d, e和f都含有没有提交的log，f相对于现在的leader少了index 4之后所有的log entry，但是多了term 2和term3的log, 这种情况可能就是由于f在term2的时候是leader，生成了三个log，但是在没有提交的的时候就挂了，迅速重启后又称为了term 3的leader,同样生成了5个log entry，但是也是没提交就挂了直到term 8的时候才启动起来。

由于leader选举时候的限制，当前的leader包含了之前leader已经commit的所有log, 对于follower缺失log的情况，leader在当选为leader的时候初始化每个follower已经匹配到的index为其log的长度，也就是最有一个log entry的下一个index，然后在appendEntries的时候，在follower查看其log对应的该index时候和leader对应的index对应的log相同，如果相同则删除其后所有的log，强制用leader的log覆盖，否则index减一向前查找，知道找到第一个匹配的log，然后删除后面的用leader的log覆盖，这样就达到了follower的log和leader的一致。对于leader缺失的log，由于选举时候的限制，说明该确实的log是没有commit的，所以可以忽略，直接覆盖即可。

那么如果leader在replicate log entry的时候，如果leader挂了，然后重新选出的leader可能不是之前的实例，那么对于之前的leader已经commit的log entry，或者可能没有commit的log entry怎么处理呢？leader出现故障可能发生在下面三个阶段：
- log entry没有replicate在多数的follower，然后crash
- log entry已经replicate在多数的follower，没有commit，然后crash
- log entry已经replicate在多数的follower，并且已经commit，然后crash

还是以论文中的例子看一下怎么处理之前term的log entry:
![commit previous log.png](@media/archive/blog/images/commit_previous_log.png?imageView/0/w/500/)
我们用(termID,indexID)表示图中的term为termID，index是indexID的log entry.图中(a)对应第一种情况，此时S1是leader，(2,2)对应的log entry只在S1,S2完成了复制，然后S1挂了，这个时候根据leader选举的约束条件，S5可以当选为leader，然后接收了一个新的log entry，对应(3,2)，这个时候还没有在其他节点完成复制，就挂了，接着S1又成为了term 4阶段的新的leader, 此时只有接收新的log entry，复制的到其他follower时候，完成之前没有完成的复制，也就是图中S3对应(2,2)的entry,这个时候完成多数的复制后，即图中的(c)。这个时候(2,2)处于复制了多数但是没有commit的状态，(4,2)只有S1的logs里有，这个时候(2,2)虽然是多数但是是不会提交的，原因后面会解释。所以这个时候就有两种情况，分别对应上面的第二种和第三种情况，对于第二种情况对应(d)，S1中的(2,2)和(4,3)都没有提交，然后挂了，这个时候S5当选为leader(term 5, S2, S3, S4，都可能投票给S5)，然后强制用S5的(3,2)覆盖S1，S2, S3, S4未提交的日志。第三种情况对应(e), (2,2)和(4,3)已经提交，然后S1挂了，这个时候因为S1，S2，S3最后一个log的term是4，index是3，而S5最后一个log的term是3，不满足leader选举的限制条件，所以不能当选leader，新的leader只能从S1, S2,S3中选举，所以之前提交的log都是没有丢失。

**<font color=red>对于(c)中的(2,2)为什么虽然已经是多数，但是还是不会提交呢?</font>**

> 因为Raft为了简化提交旧term的log entry的过程，Raft不根据旧的log entry已经完成的副本的数量进行commit，只有当前term的log entry根据多数原则提交，但是因为在完成当前term Log entry复制的过程中，会强制复制leader的log 到确实相应log entry的follower，所以当commit当前某个index对应的log的时候，会将leader该index之前的log都进行提交，而且存在多数的follower的log在改index之前的log entry是一致的，然后在appendEntrye(心跳也是这个方法)的时候会将leader已经commit的最大index发送给follower, follower会根据该index提交本地的log。

**<font color=red>会不会出现(b)中S1的(2,2)已经是多数，也就是S3也已经存在(2,2)，然后已经提交然后S1挂了，这个时候S5当选为leader，强制覆盖了S1，S2,S3里的(2,2)？</font>**
> 这种情况是不会出现的，因为(2,2)是term 2阶段生成的日志，如果(2,2)已经提交，然后leader挂了，这个时候S5是不会当选为leader的，也就是说不会出现图(b)中S5里的(3,2)这个log entry.

#### Cluster Membership Changes
当集群成员发生改变的时候，通常最直接的方法就是每个节点的配置直接更新成新的配置，但是这种方法会导致在切换的过程中，由于每个节点切换成功的时刻不一致，所以导致新旧配置产生两个Leader.还有一种方法是，将切换过程分两阶段，首先让集群暂时不可用，然后切换成新的配置，最后使用新的配置启用集群，这样显然会造成集群短暂的不可用。

为了确保在变更配置的时候，集群依然能够提供该服务，Raft提供了一种Joint Consesus的机制。其核心想法是讲配置作为一个特殊的Log Entry，使用上面的replication算法分发到每个节点。Joint Consesus将新旧配置组合到一起：
- Log Entries会复制到新旧配置中的所有节点
- 任何一个server，可能是新配置的也可能是旧配置的，都能够当选为leader
- 选举或者提交entry需要新配置节点的多数节点同意，也需要旧配置多数节点同意

下图是配置变更的一个过程:
![membership change.png](@media/archive/blog/images/membership_change.png?imageView/0/w/500/)
当当前集群leader收到请求需要将配置从$C_{old}$变更到$C_{new}$时，leader创建一个新的用于Joint Consesus的配置$C_{old,new}$，改配置作为一个log entry使用Raft的replicate算法复制到新旧配置里的所有节点，<font color=red>一旦节点接收到新的配置Log, 不管其是否已经提交都会使用该新的配置</font>，也就是说当前集群的leader会使用$C_{old,new}$来决策什么时候提交该配置，如果这个时候leader宕机，那么新选举出来的leader可能是使用old配置的节点，也可能使用$C_{old,new}$的节点。一旦$C_{old,new}$提交，这个时候有多数的节点已经有$C_{old,new}$配置，这个时候如果leader宕机，只有有$C_{old,new}$的节点才会成为leader。所以这个时候leader就会创建一个$C_{new}$log进行replicate，一旦$C_{new}$已经commit，那么旧配置不相关的节点就可以关掉了。从图中可以看出不存在一个时刻$C_{old}$和$C_{new}$同时作用的时刻，这就能保证系统配置变更的正确性。


**问:** <font color=red>新加入的节点没有任何的log entry，将会导致其短期能不能commit 新的log entries，怎么解决?</font>
> 因为新加入的节点，其log可能需要一段时间才能接收完leader之前提交的log，所以会导致需要一段时间才能接受提交的新的log，Raft新增加了一个阶段，对于新加入的节点，在其log复制完也就是跟上集群内其他节点的日志前，不参与多数节点的行列，只有接受完旧的日志，才能够正常参与表决(投票和commit表决)。

**问:** <font color=red>如果leader不在新的配置里怎么办？</font>
> 这种情况下，leader在提交$C_{new}$后再关掉，也就是说在提交新配置前，leader管理了一个不包含自己的集群，这个时候replicate log的时候，计算多数的时候不算他自己在内。

**问:** <font color=red>移除的节点在没有关掉的情况下，因为这些节点收不到leader的心跳，所以会重新发起选举，这个时候会像新集群的节点发送rpc，这个时候可能影响新节点的状态</font>
> Raft规定，如果新配置的集群能够收到leader的心跳，即使收到了选举的RPC，也会拒绝掉，不会给他投票。


#### Log Compaction
随着请求的增多，Raft每个实例的log不断的增加，现实中是不可能任其无限增长的，因此Raft也提供了快照的方法来压缩日志。每个Raft实例都单独执行快照算法，当日志大小达到一定大小的时候触发快照操作，保存最后提交的日志时刻状态机的状态，以及最后提交的一个日志的index和term。

当一个节点的日志远落后于leader节点，比如新加入节点，这个时候就需要将leader的快照发送给该节点，因此Raft也提供了一个InstallSnapshot的RPC接口用来发送快照。当follower收到leader的快照时候，根据快照里包含的LastIncludeIndex和LastIncludeTerm以及自身log的index和term来确定如果处理，如果当前节点包含一些没有提交但是和快照冲突的日志，那么清楚有冲突的日志，保留快照后面的日志(正常情况不会出现，有可能快照是重传或者产生了错误，所以收到一个自身日志之前某个index前的快照，这个时候覆盖该index之前的日志保留后续日志即可)。

#### Client Interaction
客户端与Raft集群进行交互主要存在两个问题：第一，客户端如何确定leader的位置？第二，如何保证Raft的线性一致性(Linearizable)?

客户端如何确定leader的位置？当客户端启动的时候会随机的选择一个节点请求，如果这个节点不是leader，那么会返回该节点最近知道的leader的位置(可能leader已经切换，该节点还没感知)。如果leader宕机了，那么客户端随机选择一个节点请求。

第二个问题，Raft提供线性一致性，也就是说对于一个操作，Raft只执行一次，并且立即执行，后面的读一定读到最新的数据。但是如果不加任何的限制，目前描述的Raft算法是可能出现同一个命令执行多次的情况，比如如果leader在提交了log entry后但是没有返回给客户端的时候宕机了，那么客户端超时后会重新请求改命令。这样就可能提交两次操作。针对这种情况，客户端给每个操作都分配一个递增的序列号，Raft集群每个状态机都记录一下当前每个客户端已经执行的操作的最新的序列号，如果客户端请求了同一个序列号的操作，一点状态机发现已经执行过了，就直接返回上次的结果给客户端。

### 总结
Raft在保证正确性的前提下实现一个了容易理解和实现的一致性算法，相比Paxos简化了很多，与Paxos最大的不同就是Raft是一个强leader的协议，所以的操作都依赖于leader，操作流的方向也只能从Leader发往其他节点，所以Raft整个协议分为两阶段，每个term首先进行选主，然后后续leader掌权接受所有客户端的请求，这样大大的简化了协议的复杂度，但是也存在leader负载过高的问题，不过通常实现的时候用的都是multi-raft，在接受请求前已经进行了负载的均衡，有时候还会使用旁路的监控，动态调整leader的位置达到更好的性能。


----
**附录**
Raft正确性原则：
- Election Safety
一个term只能选出一个leader
- Leader Append-Only
leader绝不覆盖或者删除日志，只会增加新的日志
- Log Matching
如果两个logs的包含一个entry其term和index都相同，那么该entry之前所有的log entry都相同
- Leader Completeness
一个选举出来的leader包含了当前term之前所有term已经提交的log entry
- State Machine Safety
如果一个节点的状态机已经应用了某个index的操作，那么其他节点的状态机在改index也是执行同样的操作，不会是其他操作
