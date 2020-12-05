```toml
title = "SeaweedFs索引-CompactMap"
date = "2017-06-18 23:39:23"
update_date = "2017-06-18 23:39:23"
author = "KDF5000"
thumb = ""
tags = ["Distributed System", "Object Storage", "SeaweedFS"]
draft = false
```
SeaweedFS提供了几种不同的needle索引策略，包括memory, btree, blotdb, leveldb四种，其中默认的是memory，也是其内部唯一自己实现的一种索引，btree使用google的btree开源实现，boltdb和leveldb都依赖一个db.

memory的索引实现，使用了一个叫CompactMap(这是作者自己起的名)的数据结构。本文后面将会重点介绍compactMap是一个怎样的数据结构，以及如何使用这个数据结构简历needle的索引。

<!--more -->

## CompactMap
compactMap本质是一个数组(Golang中可以动态扩展的数组)，其元素类型是CompactSection，定义如下：
```go
//This map assumes mostly inserting increasing keys
type CompactMap struct {
    list []*CompactSection
}
```
CompactMap结构体内部只有一个CompactSection的结构体的指针数组，那么为什么作者假设插入的key是增序的呢？我们先看CompactSection是一个什么样的结构，答案后面就会很清楚。
下面是CompactSection的结构：
```go
type CompactSection struct {
    sync.RWMutex                 //读写锁
    values   []NeedleValue       //NeedleValue数组，保存了Needle的key, offset, size
    overflow map[Key]NeedleValue //用于保存超出batch，但是在(start, end)范围的needleValue
    start    Key                 //该section保存的key最小值
    end      Key                 //该section保存的key最大值
    counter  int                 //当前该section的needleValue的数量
}
```
CompactSection保存了key值在[start, end]范围内的所有needleValue，每个section都会记录当前处于该段的所有needleValue数量，并且有一个固定的容量限制，目前版本的限制是100000，如果该段已经存储的needleValue已经达到了最大容量，但是其key那么就会将当前[start,end]范围，那么needleValue保存到map里，也就是说values数组最大长度为100000， 主要是因为当前使用二分搜索，查找key对应的needleValue, 所以数组还是不能太大。使用map也能够快速的定位。

上面简单介绍了compactSection的结构，先不管怎么向某个section插入一个needleValue, 我们只用知道一个compactSection保存从[start,end]范围内的所有needleValue。所以CompactMap的list就是将key分成一段一段管理，至于每一段的key是怎么管理的可以不用管，只用知道每个段的范围即可，当要插入一个key的时候，使用二分查找的方法找到该key所处的section, 然后插入到该section。

下面是二分查找的代码：
```go
func (cm *CompactMap) binarySearchCompactSection(key Key) int {
    l, h := 0, len(cm.list)-1
    if h < 0 {
        return -5
    }
    if cm.list[h].start <= key {
        if cm.list[h].counter < batch || key <= cm.list[h].end {
            return h
        }
        return -4
    }
    for l <= h {
        m := (l + h) / 2
        if key < cm.list[m].start {
            h = m - 1
        } else { // cm.list[m].start <= key
            if cm.list[m+1].start <= key {
                l = m + 1
            } else {
                //貌似这里有bug, 如果处于(end, $$start_{m+1}$$), 则同样要插入一个新的section
                /*if(key > cm.list[m].end){
                    return -5
                }*/
                return m
            }
        }
    }
    return -3
}
```
下面以上面的代码分析怎么找到一个NeedleValue应该插入的section:
* 如果list为空则直接返回一个负值
* 先判断是否能够插入到最后一个section, 能够插入的条件是该key大于section的最小key值，如果当前section没有超过容量，或者超过了但是处于[start, end]范围，则返回当前section, 否则返回一个负值。
* 前两个条件都不满足的话，就二分搜索前面所有section能够插入的section, 如果到第一个section都没有找到，即key< list[0].start， 说明需要插入一个新的section, 返回一个负值

总的来说，搜索有两个结果，返回一个可以插入的section和返回一个负值，当返回一个负值的时候，就需要重新生成一个compactSection，使用插入排序的方法，将其插入到list
代码如下：
```go
func (cm *CompactMap) Set(key Key, offset, size uint32) (oldOffset, oldSize uint32) {
    x := cm.binarySearchCompactSection(key)
    if x < 0 {
        //println(x, "creating", len(cm.list), "section, starting", key)
        cm.list = append(cm.list, NewCompactSection(key))
        x = len(cm.list) - 1
        //keep compact section sorted by start
        for x > 0 {
            //插入排序，保证start有序
            if cm.list[x-1].start > cm.list[x].start {
                cm.list[x-1], cm.list[x] = cm.list[x], cm.list[x-1]
                x = x - 1
            } else {
                break
            }
        }
    }
    return cm.list[x].Set(key, offset, size)
}
```
> **Note:** ~~感觉这里存在一个bug, 因为二分搜索的时候是根据start来查找能够插入的section， 但是存在一种情况，比如假设list[i].start <= key, 且list[i+1].start > key， 那么上面的程序就会返回插入到第i个section, 但是可能list[i].end < key， 也就是说实际上这个时候应该生成一个新的setction插入到除以i和i+1之间。当然作者假设key是递增的（并不一定要连续），如果key是递增的，这样就能能够保证不会出现出现插入一个中间断节的list.但是如果假设是递增的，何必使用二分查找合适的section呢？肯定是最后一个或者需要重新在后面追加一个，所以目前感觉这个逻辑还是有问题的。考虑的非递增的情况，但是非递增的情况下，在二分查找section又没有考虑断节的情况，所以这里应该有问题，二分查找的代码里21-23行是我添加的，应该可以解决这个bug。后分析发现，因为插入除了最后一个section时候，同样也更新了end所以就不存在这个问题~~

### CompactSectionMap
前面已经简单的介绍了CompactSectionMap的结构，主要有一个数组和一个map组成，并且每个section的数组都有一个容量，下面主要分析怎么set和get一个needleValue。
下面是set的核心代码：
```go
//return old entry size
func (cs *CompactSection) Set(key Key, offset, size uint32) (oldOffset, oldSize uint32) {
    cs.Lock()
    if key > cs.end {
        cs.end = key
    }
    if i := cs.binarySearchValues(key); i >= 0 {
        oldOffset, oldSize = cs.values[i].Offset, cs.values[i].Size
        //println("key", key, "old size", ret)
        cs.values[i].Offset, cs.values[i].Size = offset, size
    } else {
        needOverflow := cs.counter >= batch
        needOverflow = needOverflow || cs.counter > 0 && cs.values[cs.counter-1].Key > key
        if needOverflow {
            //println("start", cs.start, "counter", cs.counter, "key", key)
            if oldValue, found := cs.overflow[key]; found {
                oldOffset, oldSize = oldValue.Offset, oldValue.Size
            }
            cs.overflow[key] = NeedleValue{Key: key, Offset: offset, Size: size}
        } else {
            p := &cs.values[cs.counter]
            p.Key, p.Offset, p.Size = key, offset, size
            //println("added index", cs.counter, "key", key, cs.values[cs.counter].Key)
            cs.counter++
        }
    }
    cs.Unlock()
    return
}
```
下面是set的过程：
* 在插入前，首先会进行加锁的操作
* 然后如果当前的key大于end, 那么更新end(这里优于，在执行set之前，compactMap的set里有个判断，如果key大于最后一个section时，如果小于当前容量，则直接返回最后一个section，如果大于当前容量，那么要求要小于最后一个sectiond的end, 这样的结果就是最后一个section的end只有在小于当前容量的时候才会更新。低于其他的section就没有这个限制，这样可以防止中间出现断节从而需要插入新的section)
* 使用二分查找在values数组中查找是否含有该key, 如果含有那么就更新对应的needleValue，返回旧的needleValue
* 如果没有找到，则判断是否需要插入到overflow map中，这里判断是否需要插入到溢出map中并不单单，根据overflow来判断，还有一种情况是，没有大于数组的容量，但是该key小于values数组中最后一个needleValue的key, 这样能够保证values数组的key是有序的
* 如果不需要溢出，那么久直接追加在values数组的最后面

从compactSectionMap中获取一个needleValue的操作比较简单，下面是核心代码：
```go
func (cs *CompactSection) Get(key Key) (*NeedleValue, bool) {
    cs.RLock()
    if v, ok := cs.overflow[key]; ok {
        cs.RUnlock()
        return &v, true
    }
    if i := cs.binarySearchValues(key); i >= 0 {
        cs.RUnlock()
        return &cs.values[i], true
    }
    cs.RUnlock()
    return nil, false
}
```
get的时候也同样要先加锁，然后先判断overflow map中 是否存在改key, 如果不存在二分搜索values数组，如果都不存在，则返回没有找到即可。

## 总结
整个CompactMap的设计思路还是比较清晰的，很像skiplist，通过分段(section)，可以大大的加快查找速度，不过如果key比较随机，那么就会带来很多的二分查找，插入排序的操作，性能就不是那么理想，不过总的来说这个数据结构还是比较高效的。
