# cas

## autoInc.go

- 为并发请求分配自增整形ID

### 实现

- sync.WaitGroup控制协程完成
- 导入sync/atomic包，可以支持atomic.CompareAndSwapInt32(&counter, old, old+1)
- 使用自旋锁for{} + atomic.CompareAndSwapInt32()为并发请求分配自增ID
- 导入sync包，对sync.Map类型的全局变量sm，支持sm.Range()遍历和sm.Store(index, old+1)存储功能
- 分配的kv对保存到sync.Map类型的全局变量中

#### 输出

```Shell
final counter:1000
final map is:
key=26,val=13
...
key=864,val=877
key=948,val=963
key set len is: 1000
val set len is: 1000
```

## counter.go

- 与autoInc.go不同，只输出每个并发请求ID与CAS写入失败重试次数
- 可以在此基础上，开展有限制重试次数的CAS操作，超过N后宣告写入失败
