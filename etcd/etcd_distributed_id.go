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
 
    client, err := clientv3.New(em.Conf) // New creates a new etcdv3 client from a given configuration.
    if err != nil { 
        return err 
    } 
 
    em.txn = clientv3.NewKV(client).Txn(context.TODO()) // type KV interface contains functions like Put, Get, Txn, ... Txn creates a transaction.
    em.lease = clientv3.NewLease(client) // type Lease interface contains functions like Grant, ...
    // 请求服务器
    leaseResp, err := em.lease.Grant(context.TODO(), em.Ttl) // Grant creates a new lease. Return *LeaseGrantResponse
    // type LeaseGrantResponse contains *pb.ResponseHeader, ID LeaseID, TTL int64, Error string
    // 响应举例：cluster_id:10316109323310759371 member_id:12530464223947363063 revision:12 raft_term:3  <nil>
    fmt.Println("获得新租约", leaseResp.ID)
 
    if err != nil { 
        return err 
    } 
 
    ctx, em.cancel = context.WithCancel(context.TODO()) 
    em.leaseID = leaseResp.ID 
    _, err = em.lease.KeepAlive(ctx, em.leaseID) // KeepAlive attempts to keep the given lease alive forever. 可以理解为自动定时续约某个租约，返回只读Channel，每次自动续租成功后会向通道中发送信号。
 
    return err 
} 
  
// 获取锁 
func (em *EtcdMutex) Lock() error {  
    err := em.init() 
    if err != nil { 
        return err 
    } 
 
    // LOCK 
    // 事务：如果版本=0，则将Key的值置为""并绑定已获得租约
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

    fmt.Println("抢锁成功", em.leaseID)
 
    return nil 
} 
 
//释放锁方法1 
func (em *EtcdMutex) UnLockWithParam(id clientv3.LeaseID) { 
    // 租约自动过期，立刻过期 
    // cancel取消续租，而revoke则是立即过期 
    em.cancel() 
    em.lease.Revoke(context.TODO(), id) 
    // em.leaseID
 
    fmt.Println("释放租约", id) 
} 

//释放锁方法2 
func (em *EtcdMutex) UnLockWithoutParam() { 
    // 租约自动过期，立刻过期 
    // cancel取消续租，而revoke则是立即过期 
    em.cancel() 
    em.lease.Revoke(context.TODO(), em.leaseID) // 如果外部用defer多次调用有风险：因为em对象拿到的leaseID是变化的，但是多次defer执行时才从em取租约ID都是最后一次的值。
 
    fmt.Println("释放租约", em.leaseID) 
} 
 
// goroutine1 
func try2lock1() { 
    eMutex1 := &EtcdMutex{ 
        Conf: conf, 
        Ttl:  30, 
        Key:  "lock", 
    } 
 
    gotLock := false
    for i:=0; i<2; i++ {
        err := eMutex1.Lock() 
        // eMutex1.UnLockWithoutParam() 
        defer eMutex1.UnLockWithParam(eMutex1.leaseID) 
        if err != nil { 
            fmt.Println("goroutine1抢锁失败", eMutex1.leaseID) 
            // 如果此时退出，表示不处理抢锁失败的租约，到期自动退租，也不会再次抢锁。
            // return 
        } else {
            gotLock = true
            fmt.Println("goroutine1抢锁成功", eMutex1.leaseID) 
            fmt.Println("goroutine1执行业务逻辑")
            time.Sleep(1 * time.Second)
            break
        }
        fmt.Println("goroutine1等待一段时间重新抢锁")
        time.Sleep(1 * time.Second)
    }
    if !gotLock {
        fmt.Println("goroutine1可以在此做一些释放租约前的操作, 或者记录重试后仍然无法获得锁的日志")
    }

    // 如果此时退出，表示已主动退租(UnLock)
    return
} 
 
// goroutine2 
func try2lock2() {
    eMutex2 := &EtcdMutex{ 
        Conf: conf, 
        Ttl:  30, 
        Key:  "lock", 
    } 
 
    // 不带重试的做法
    // {
    //     err := eMutex2.Lock() 
    //     if err != nil { 
    //         fmt.Println("goroutine2抢锁失败") 
    //         // return 
    //     } 
    
    //     defer eMutex2.UnLock() 
    //     fmt.Println("goroutine2抢锁成功") 
    // }

    // 带重试的做法
    {
        gotLock := false
        for i:=0; i<2; i++ {
            err := eMutex2.Lock() 
            // eMutex2.UnLockWithoutParam() // 注意：无参数不带defer的原因是，UnLock中有释放租约的行为，如果重试了2次，在defer执行时才获取的租约ID都是最后一个。要么就将ID作为参数传入，要么就不用defer
            defer eMutex2.UnLockWithParam(eMutex2.leaseID) // defer更合适。因为多次重试时避免了自己选择解锁和执行业务逻辑的位置和顺序了，以及无论业务逻辑出现什么问题都基本能保障及时退租。
            if err != nil { 
                fmt.Println("goroutine2抢锁失败", eMutex2.leaseID) 
                // return // 此时return就不再重试抢锁
            } else {
                gotLock = true
                fmt.Println("goroutine2抢锁成功", eMutex2.leaseID) 
                fmt.Println("goroutine2执行业务逻辑")
                time.Sleep(1 * time.Second)
                break
            } 
            fmt.Println("goroutine2等待一段时间重新抢锁")
            time.Sleep(1 * time.Second)
        }

        if !gotLock {
            fmt.Println("goroutine2可以在此做一些释放租约前的操作, 或者记录重试后仍然无法获得锁的日志")
        }
    }

    return
} 
 
// 测试代码 
func EtcdRunTester() { 
    conf = clientv3.Config{ 
        Endpoints:   []string{"host.docker.internal:23791", "host.docker.internal:23792", "host.docker.internal:23793"}, // 集群或单点都可以
        DialTimeout: 3 * time.Second, 
    } 
 
    // 启动两个协程竞争锁 
    go try2lock1() 
    go try2lock2() 
 
    time.Sleep(5 * time.Second) 
} 

func main() {
    EtcdRunTester()
}

/*
输出：
获得新租约 1366716111406998522
获得新租约 1366716111406998524
抢锁成功 1366716111406998522
释放租约 1366716111406998524
goroutine2抢锁失败 1366716111406998524
goroutine2等待一段时间重新抢锁
释放租约 1366716111406998522
goroutine1抢锁成功 1366716111406998522
获得新租约 1366716111406998530
抢锁成功 1366716111406998530
释放租约 1366716111406998530
goroutine2抢锁成功 1366716111406998530
*/

/*
解释：
由于在goroutine中进行，所以没有严格的先后顺序
22(g1)和24(g2)分别是两个goroutine最先获得的租约ID
g1(22)抢锁成功，将值写入到"lock"中，并经过短暂的业务逻辑操作后，释放了租约，g1抢锁全过程已结束
同时，g2(24)抢锁失败，不得不主动释放租约，并等待一段时间后准备重新抢锁
一段时间后
g2获得了新租约30，并成功抢锁，经过业务逻辑后，释放了租约，g2抢锁全过程也结束
*/

// 参考：https://www.51cto.com/article/662978.html