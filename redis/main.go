package main

import (
	"context"
	"fmt"
	"time"
	"errors"
	"sync"

	"github.com/go-redis/redis/v8" // 注意导入的是新版本
)

var (
	rdb *redis.Client
)

const ADDR string = "host.docker.internal:6379"
const PASS string = "123456"
const TIMEOUT int = 5

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     ADDR,
		Password: PASS,  // no password set
		DB:       0,   // use default DB
		PoolSize: 100, // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(TIMEOUT)*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}

func main() {
	ctx := context.Background()
	GetSet(ctx)
	fmt.Println("---------------------------")
	Zset(ctx)
	fmt.Println("---------------------------")
	Do(ctx)
	fmt.Println("---------------------------")
	Pipeline(ctx)
	fmt.Println("---------------------------")
	TxPipiline(ctx)
	fmt.Println("---------------------------")
	TransactionDemo(ctx)
	fmt.Println("---------------------------")
	Set(ctx)
	fmt.Println("---------------------------")
	Hash(ctx)
	fmt.Println("---------------------------")
	List(ctx)
}

func GetSet(ctx context.Context) {
	if err := initClient(); err != nil {
		panic(err)
	}

	err := rdb.Set(ctx, "key", "value", 0).Err() // expiration=0表示没有过期时间，持久有效
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result() // Result()返回结果和错误
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func Zset(ctx context.Context) {
	zsetKey := "language_rank"
	languages := []*redis.Z{ // 指针切片
		&redis.Z{Score: 90.0, Member: "Golang"}, // 结构redis.Z是固定的
		&redis.Z{Score: 98.0, Member: "Java"},
		&redis.Z{Score: 95.0, Member: "Python"},
		&redis.Z{Score: 97.0, Member: "JavaScript"},
		&redis.Z{Score: 99.0, Member: "C/C++"},
	}
	// ZADD
	num, err := rdb.ZAdd(ctx, zsetKey, languages...).Result() // 要传入ctx
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Printf("zadd %d succ.\n", num)

	// 把Golang的分数加10
	newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10.0, "Golang").Result()
	if err != nil {
		fmt.Printf("zincrby failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	// 取分数最高的3个
	ret, err := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Result()
	if err != nil {
		fmt.Printf("zrevrange failed, err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}

	// 取95~100分的
	op := &redis.ZRangeBy{ // 指针
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}

func Do(ctx context.Context) {
	res, err := rdb.Do(ctx, "del", "key").Result()
	if err != nil {
		fmt.Printf("do failed, err:%v\n", err)
		return
	}
	fmt.Printf("do del result %d \n", res)
}

func Pipeline(ctx context.Context) {
	// 方法1
	pipe := rdb.Pipeline()

	incr := pipe.Incr(ctx, "pipeline_counter")
	pipe.Expire(ctx, "pipeline_counter", time.Hour)

	_, err := pipe.Exec(ctx)
	fmt.Println(incr.Val(), err)

	// 方法2
	// var incr *redis.IntCmd
	_, err = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(ctx, "pipelined_counter")
		pipe.Expire(ctx, "pipelined_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}

// 事务，对应于redis的MULTI包裹的事务
func TxPipiline(ctx context.Context) {
	// 方法1
	pipe := rdb.TxPipeline()

	incr := pipe.Incr(ctx, "tx_pipeline_counter")
	pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)

	_, err := pipe.Exec(ctx)
	fmt.Println(incr.Val(), err)

	// 方法2
	// var incr *redis.IntCmd
	_, err = rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		incr = pipe.Incr(ctx, "tx_pipelined_counter")
		pipe.Expire(ctx, "tx_pipelined_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}

func TransactionDemo(ctx context.Context) {
	var (
		maxRetries   = 1000 // 乐观锁最大尝试次数
		routineCount = 10 // 并发请求数量
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 带超时时间，超时退出协程
	defer cancel()

	// Increment 使用GET和SET命令以事务方式递增Key的值
	increment := func(key string) error {
		// 事务函数
		txf := func(tx *redis.Tx) error {
			// 获得key的当前值或零值
			n, err := tx.Get(ctx, key).Int()
			if err != nil && err != redis.Nil {
				return err
			}

			// 实际的操作代码（乐观锁定中的本地操作）
			n++

			// 操作仅在 Watch 的 Key 没发生变化的情况下提交
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, key, n, 0)
				return nil
			})
			return err
		}

		// 最多重试 maxRetries 次
		for i := 0; i < maxRetries; i++ {
			err := rdb.Watch(ctx, txf, key)
			if err == nil {
				// 成功
				return nil
			}
			if err == redis.TxFailedErr {
				// 乐观锁丢失 重试
				continue
			}
			// 返回其他的错误
			return err
		}

		return errors.New("increment reached maximum number of retries")
	}

	start, err := rdb.Get(context.TODO(), "counter3").Int()
	fmt.Println("started with", start, err)
	fmt.Printf("increment with %d \n", routineCount)

	// 模拟 routineCount 个并发同时去修改 counter3 的值
	var wg sync.WaitGroup
	wg.Add(routineCount)
	for i := 0; i < routineCount; i++ {
		go func() {
			defer wg.Done()
			if err := increment("counter3"); err != nil {
				fmt.Println("increment error:", err)
			}
		}()
	}
	wg.Wait()

	n, err := rdb.Get(context.TODO(), "counter3").Int()
	fmt.Println("ended with", n, err)
}

func Set(ctx context.Context) {
	addCmd := rdb.SAdd(ctx, "set", "s1", "s2", "s3")
	fmt.Println(addCmd)
	stringSliceCmd := rdb.SMembers(ctx, "set")
	for _,v := range stringSliceCmd.Val() {
		fmt.Println(v)
	}
}

func Hash(ctx context.Context) {
	intCmd := rdb.HSet(ctx, "hash", map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"})
	fmt.Println(intCmd)
	// 通过指定map中的key获取值
	getOne := rdb.HGet(ctx, "hash", "k1")
	fmt.Println(getOne.Val())
	// 获取所有的key-value键值对
	all := rdb.HGetAll(ctx, "hash")
	for key, value := range all.Val() {
		fmt.Println("key --> ", key, " value --> ", value)
	}
}

func List(ctx context.Context) {
	intCmd := rdb.LPush(ctx, "list", "l1", "l2", "l3")
	fmt.Println(intCmd)
	lRange := rdb.LRange(ctx, "list", 0, 3) // 从最左边开始取数据
	for _, v := range lRange.Val() {
		fmt.Println(v)
	}
}

// 参考
// https://www.liwenzhou.com/posts/Go/go_redis/
// https://juejin.cn/post/7101286398388338695