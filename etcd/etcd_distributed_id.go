package main 
 
import ( 
    "context" 
    "fmt" 
    "go.etcd.io/etcd/clientv3" 
    "time" 
) 
 
var conf clientv3.Config 
 
// 锁结构体 
type EtcdMutex struct { 
    Ttl int64//租约时间 
 
    Conf   clientv3.Config    //etcd集群配置 
    Key    string//etcd的key 
    cancel context.CancelFunc //关闭续租的func 
 
    txn     clientv3.Txn 
    lease   clientv3.Lease 
    leaseID clientv3.LeaseID 
} 
 
// 初始化锁 
func (em *EtcdMutex) init() error { 
    var err error 
    var ctx context.Context 
 
    client, err := clientv3.New(em.Conf) 
    if err != nil { 
        return err 
    } 
 
    em.txn = clientv3.NewKV(client).Txn(context.TODO()) 
    em.lease = clientv3.NewLease(client) 
    leaseResp, err := em.lease.Grant(context.TODO(), em.Ttl) 
 
    if err != nil { 
        return err 
    } 
 
    ctx, em.cancel = context.WithCancel(context.TODO()) 
    em.leaseID = leaseResp.ID 
    _, err = em.lease.KeepAlive(ctx, em.leaseID) 
 
    return err 
} 
 
// 获取锁 
func (em *EtcdMutex) Lock() error { 
    err := em.init() 
    if err != nil { 
        return err 
    } 
 
    // LOCK 
    em.txn.If(clientv3.Compare(clientv3.CreateRevision(em.Key), "=", 0)). 
        Then(clientv3.OpPut(em.Key, "", clientv3.WithLease(em.leaseID))).Else() 
 
    txnResp, err := em.txn.Commit() 
    if err != nil { 
        return err 
    } 
 
    // 判断txn.if条件是否成立 
    if !txnResp.Succeeded { 
        return fmt.Errorf("抢锁失败") 
    } 
 
    returnnil 
} 
 
//释放锁 
func (em *EtcdMutex) UnLock() { 
    // 租约自动过期，立刻过期 
    // cancel取消续租，而revoke则是立即过期 
    em.cancel() 
    em.lease.Revoke(context.TODO(), em.leaseID) 
 
    fmt.Println("释放了锁") 
} 
 
// groutine1 
func try2lock1() { 
    eMutex1 := &EtcdMutex{ 
        Conf: conf, 
        Ttl:  10, 
        Key:  "lock", 
    } 
 
    err := eMutex1.Lock() 
    if err != nil { 
        fmt.Println("groutine1抢锁失败") 
        return 
    } 
    defer eMutex1.UnLock() 
 
    fmt.Println("groutine1抢锁成功") 
    time.Sleep(10 * time.Second) 
} 
 
// groutine2 
func try2lock2() { 
    eMutex2 := &EtcdMutex{ 
        Conf: conf, 
        Ttl:  10, 
        Key:  "lock", 
    } 
 
    err := eMutex2.Lock() 
    if err != nil { 
        fmt.Println("groutine2抢锁失败") 
        return 
    } 
 
    defer eMutex2.UnLock() 
    fmt.Println("groutine2抢锁成功") 
} 
 
// 测试代码 
func EtcdRunTester() { 
    conf = clientv3.Config{ 
        Endpoints:   []string{"127.0.0.1:2379"}, 
        DialTimeout: 5 * time.Second, 
    } 
 
    // 启动两个协程竞争锁 
    go try2lock1() 
    go try2lock2() 
 
    time.Sleep(300 * time.Second) 
} 
