实现的cache拥有以下特点：
* [x] 指定最大的内存使用大小
* [x] 采用LRU淘汰策略
* [x] 使用锁提供单机并发访问控制
* [ ] 可用于单机缓存，并提供http 服务
* [ ] 拥有防止缓存穿透的机制
* [ ] 采用一致性hash选择节点，实现负载均衡
* [ ] 可用于基于HTTP协议的分布式缓存
* [ ] 使用protobuf优化节点间的二进制通信