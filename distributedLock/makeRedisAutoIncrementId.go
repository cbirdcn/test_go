package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"sync"
	"time"
)

var pool = newPool()

func main(){
	// 使用redigo
	testRedisLock() //
	// 方式1：setnx+uuid
	// 2：setnx+getset+milliSecondTimestamp

}

// redis连接池
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle: 20,
		MaxActive: 1000, // max number of connections
		// 访问容器会降低效率，需要局域网访问
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "172.21.0.20:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

}

// 返回自增id
func incrId(p redis.Conn, idKey string) (id uint64, err error){
	id, err = redis.Uint64(p.Do("INCR", idKey))
	return id, err
}

/////////testRedisConcurrentLock///////////


/**
redis lock
1.无死锁，但不可重入：依靠锁ttl判断是否占用，重入表示获取一个新锁（可以手动添加uuid解决重入），死锁依靠过期判断和redis自动清理过期id
2.非分布式：即使分布式部署服务器，也要请求同一个redis服务器，否则无法顺序获取id。可以加雪花id做成分布式非连续id。
3.非阻塞锁：未拿到锁时客户端轮询sleep时间请求服务端
4.高可用：大规模请求过来后，只要可以获取redis连接，最差结果就是等待到超时时间后没拿到锁
5.性能：获取锁+释放锁获取自增id，并发150-200成功率。
使用：
1.手动获取锁，返回获取锁结果（T/F）和过期时间(ms/0)
2.获取锁成功，给自增idINCR。
3.无论获取成功还是失败
*/


func testRedisLock(){
	for j:=1;j<=10;j++ {
		successCount := 0

		//start := time.Now().UnixNano() / 1e6
		var wg sync.WaitGroup

		for i:=1;i<=200;i++{
			wg.Add(1)
			go req2(&wg, i, &successCount)
			//go req2(&wg, i, &successCount)
		}
		wg.Wait()
		//fmt.Println(time.Now().UnixNano() / 1e6 - start)
		fmt.Println(successCount)
	}
}

var waitTimeout = 1000 // 最长等待时间，ms
var keepTimeout = 1000 // 最长锁定时间，ms
var lockKey = "orderLockKey" // 锁名

// 从连接池获取inst
func getInstance() redis.Conn{
	return pool.Get()
}

// 当前毫秒时间戳
func currentMilliSecond() uint64{
	return uint64(time.Now().UnixNano() / 1e6)
}

////////////////setnx+uuid//////////////////////////
//逻辑：set k v ex px ttl获取锁，拿不到就循环重新获取，直到拿到锁或等待时间到期。只需要释放已拿到的锁，当锁值=uuid时释放。即使不释放redis也可以根据锁到期时间删除锁。
//性能：每秒150-200并发获取+释放自增锁
//todo:分布式，需要支持多客户端能索引到同一台服务器的同一个key


// setnx+uuid
func req2(wg *sync.WaitGroup, i int, successCount *int){
	p := getInstance()

	// 提供参数：连接、锁名、是否自增、并发id(uuid4)
	u4 := uuid.New().String()
	idKey := "id"

	// 实际使用时，需要设定临时lockKey，比如根据方法名
	locked, _ := Lock2(p, lockKey, u4,true, idKey)

	if locked {
		//fmt.Println("locking " + strconv.Itoa(i) + " " + u4)
		// 检查id结果
		//redis.Int(p.Do("hset", "result", i, id))
		// 获取到锁才需要释放，即使不释放也会被redis根据ttl到期自动回收
		Unlock2(p, lockKey, u4, i)
		// 成功拿锁数量
		*successCount++
	}else{
		//fmt.Println("can't get lock " + strconv.Itoa(i) + " " + u4)
	}

	wg.Done()
}


// 获取锁，获取到返回true，否则返回false
// 参数：incr表示是否需要自增，true表示自增（从1开始）,u4表示并发id
// 返回：locked=true表示获取锁成功，id是获取的自增id（incr=false时或locked=false时id=0，正常id>0)，expireAt是锁生效截至日期,idKey是自增键名（返回INCR后的值，可以避免返回0，todo：和db incr打通）
func Lock2(p redis.Conn, lockKey string, u4 string, incr bool, idKey string) (locked bool, id uint64){
	locked = false
	start := currentMilliSecond()

	// 如果等待时间还没到超时时间，循环
	for ;int(currentMilliSecond() - start) < waitTimeout; {
		// 官方：SET resource_name my_random_value NX PX 30000
		// The command will set the key only if it does not already exist (NX option), with an expire of 30000 milliseconds (PX option). The key is set to a value “my_random_value”.
		// https://redis.io/topics/distlock
		res, err := redis.String(p.Do("SET", lockKey, u4, "NX", "PX", keepTimeout))

		// 假设有ABC三个协程同时获取锁
		// A协程获取到锁
		if res == "OK" && err == nil{
			// 如果要返回自增id
			id := uint64(0)
			if incr {
				// INCR过程最好放到acquireLock()里，并且必须用INCR返回值
				id, err = incrId(p, idKey)
				// 只能认为id基本可靠
			}

			// 获取锁成功
			return true, id
		}

		// BC协程没获取到锁，加延迟后继续循环，不控制循环次数和续期时长
		// TODO：不稳定，1秒钟成功数量少则10+，多则70+
		time.Sleep(time.Duration(10)*time.Millisecond)
	}

	// 等待超时，获取锁失败
	return false, 0
}

// 释放锁
// 获取锁成功才需要调用
func Unlock2(p redis.Conn, lockKey string, u4 string, i int){
	newValue, _ := redis.String(p.Do("Get", lockKey))

	// TODO：如果支持lua，把判等和del放到一个lua脚本内执行更好。pika不支持lua所以保持原样
	if u4 == newValue {
		_, _ = p.Do("DEL", lockKey)
		//fmt.Println("del " + strconv.Itoa(i) + " " + u4)
	}else{
		//fmt.Println("not del " + strconv.Itoa(i) + " " + u4)
	}

	// 释放连接池
	p.Close()
}