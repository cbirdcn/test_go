package main

import (
	"log"
	"fmt"
	"time"
	"github.com/tidwall/buntdb"
)

func main() {
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open("data.db")
	// db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var config buntdb.Config
	if err := db.ReadConfig(&config); err != nil{
		log.Fatal(err)
	}
	// fmt.Println(config)
	// 默认配置：{1 100 33554432 false <nil> <nil>}
	// type Config struct {
	// 	SyncPolicy SyncPolicy // 调整数据同步到磁盘的频率，取值为Never、EverySecond和Always。默认为EverySecond。
	// 	AutoShrinkPercentage int // AutoShrinkPercentage由后台进程触发，当aof文件的大小大于上次压缩文件结果的百分比。例如，如果此值为100，则最后一次收缩进程产生一个100mb的文件，那么新的aof文件必须是200mb之前触发收缩。
	// 	AutoShrinkMinSize int // AutoShrinkMinSize定义自动收缩之前aof文件的最小大小。
	// 	AutoShrinkDisabled bool // AutoShrinkDisabled关闭自动背景收缩
	// 	OnExpired func(keys []string) // OnExpired用于自定义处理已过期键时的删除选项
	// 	OnExpiredSync func(key, value string, tx *Tx) error // OnExpiredSync将在执行删除过期项的同一事务中调用。如果OnExpired存在，那么这个回调将不会被调用。如果存在此回调，则删除超时项是此回调的显式责任。
	// }
	// config.SyncPolicy = buntdb.Never
	// if err := db.SetConfig(config); err != nil{
	// 	log.Fatal(err)
	// }

	noConcurrence()

	count := 100
	withConcurrence(count)

}

func noConcurrence() {
	t := time.Now().UnixMilli()
	for i:=1;i<=10000; i++ {
		err = db.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(string(i), string(i), nil)
			return err
		})
	}
	fmt.Println(time.Now().UnixMilli()-t)
}

func withConcurrence(count int) {
	t := time.Now().UnixMilli()
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			for i:=1;i<=10000; i++ {
				err = db.Update(func(tx *buntdb.Tx) error {
					_, _, err := tx.Set(string(i), string(i), nil)
					return err
				})
			}
			wg.Done()
		}()
	}
	
	wg.Wait()
	fmt.Println(time.Now().UnixMilli()-t)
}

/*
benchmark测试
默认为持久化aof到文件的模式，主要看SET

准备：
go install github.com/tidwall/buntdb-benchmark@latest


官方环境：MacBook Pro 15" 2.8 GHz Intel Core i7
官方结果：
$ buntdb-benchmark -q
GET: 4609604.74 operations per second
SET: 248500.33 operations per second
ASCEND_100: 2268998.79 operations per second
ASCEND_200: 1178388.14 operations per second
ASCEND_400: 679134.20 operations per second
ASCEND_800: 348445.55 operations per second
DESCEND_100: 2313821.69 operations per second
DESCEND_200: 1292738.38 operations per second
DESCEND_400: 675258.76 operations per second
DESCEND_800: 337481.67 operations per second
SPATIAL_SET: 134824.60 operations per second
SPATIAL_INTERSECTS_100: 939491.47 operations per second
SPATIAL_INTERSECTS_200: 561590.40 operations per second
SPATIAL_INTERSECTS_400: 306951.15 operations per second
SPATIAL_INTERSECTS_800: 159673.91 operations per second


本地环境：macmini2014 2.6G双核I5 8G内存
Go: 1.21
结果：
test-Mac-mini:Documents mac$ /Users/mac/go/bin/buntdb-benchmark -q
GET: 3946212.96 operations per second
SET: 106740.65 operations per second
ASCEND_100: 1360166.18 operations per second
ASCEND_200: 858437.66 operations per second
ASCEND_400: 444723.60 operations per second
ASCEND_800: 190225.09 operations per second
DESCEND_100: 1325592.25 operations per second
DESCEND_200: 748689.35 operations per second
DESCEND_400: 415582.18 operations per second
DESCEND_800: 210034.51 operations per second
SPATIAL_SET: 60797.97 operations per second
SPATIAL_INTERSECTS_100: 709376.23 operations per second
SPATIAL_INTERSECTS_200: 415195.97 operations per second
SPATIAL_INTERSECTS_400: 184845.18 operations per second
SPATIAL_INTERSECTS_800: 122983.24 operations per second


docker容器：3CPU，3G内存
宿主机：Mac mini (2018) 16G内存 3 GHz Intel Core i5
结果：
sh-4.2# buntdb-benchmark -q
GET: 2916780.69 operations per second
SET: 2232.89 operations per second
ASCEND_100: 1334286.01 operations per second
ASCEND_200: 871977.72 operations per second
ASCEND_400: 526805.34 operations per second
ASCEND_800: 264287.10 operations per second
DESCEND_100: 1662264.13 operations per second
DESCEND_200: 951752.66 operations per second
DESCEND_400: 508810.00 operations per second
DESCEND_800: 270697.21 operations per second
SPATIAL_SET: 2227.42 operations per second
SPATIAL_INTERSECTS_100: 765721.06 operations per second
SPATIAL_INTERSECTS_200: 535045.84 operations per second
SPATIAL_INTERSECTS_400: 284266.68 operations per second
SPATIAL_INTERSECTS_800: 138849.64 operations per second


容器：3CPU，3G内存
宿主机：Mac mini (2018) 16G内存 3 GHz Intel Core i5
结果：
sh-4.2# buntdb-benchmark -q
GET: 2916780.69 operations per second
SET: 2232.89 operations per second
ASCEND_100: 1334286.01 operations per second
ASCEND_200: 871977.72 operations per second
ASCEND_400: 526805.34 operations per second
ASCEND_800: 264287.10 operations per second
DESCEND_100: 1662264.13 operations per second
DESCEND_200: 951752.66 operations per second
DESCEND_400: 508810.00 operations per second
DESCEND_800: 270697.21 operations per second
SPATIAL_SET: 2227.42 operations per second
SPATIAL_INTERSECTS_100: 765721.06 operations per second
SPATIAL_INTERSECTS_200: 535045.84 operations per second
SPATIAL_INTERSECTS_400: 284266.68 operations per second
SPATIAL_INTERSECTS_800: 138849.64 operations per second


容器：2CPU，4G内存
宿主机：Mac mini (2018) 16G内存 3 GHz Intel Core i5
结果：
sh-4.2# buntdb-benchmark -q
GET: 2335840.08 operations per second
SET: 2210.37 operations per second
ASCEND_100: 1121585.75 operations per second
ASCEND_200: 589481.35 operations per second
ASCEND_400: 334646.79 operations per second
ASCEND_800: 186428.44 operations per second
DESCEND_100: 1278634.03 operations per second
DESCEND_200: 697258.41 operations per second
DESCEND_400: 323260.85 operations per second
DESCEND_800: 191232.53 operations per second
SPATIAL_SET: 2158.86 operations per second
SPATIAL_INTERSECTS_100: 533128.24 operations per second
SPATIAL_INTERSECTS_200: 287342.85 operations per second
SPATIAL_INTERSECTS_400: 188346.55 operations per second
SPATIAL_INTERSECTS_800: 95478.62 operations per second


容器：4CPU，8G内存
宿主机：Mac mini (2018) 16G内存 3 GHz Intel Core i5
结果：
sh-4.2# buntdb-benchmark -q
GET: 2621750.42 operations per second
SET: 1803.15 operations per second
ASCEND_100: 1614269.13 operations per second
ASCEND_200: 991942.67 operations per second
ASCEND_400: 525048.16 operations per second
ASCEND_800: 338188.24 operations per second
DESCEND_100: 1771707.14 operations per second
DESCEND_200: 1194943.26 operations per second
DESCEND_400: 587932.44 operations per second
DESCEND_800: 326081.67 operations per second
SPATIAL_SET: 1613.28 operations per second
SPATIAL_INTERSECTS_100: 546104.04 operations per second
SPATIAL_INTERSECTS_200: 477107.69 operations per second
SPATIAL_INTERSECTS_400: 284697.63 operations per second
SPATIAL_INTERSECTS_800: 153132.31 operations per second


结论1：容器内做基准测试是不准确的，尤其是大量读写操作时，对本就是模拟的硬盘镜像性能更是考验。所以容器测试不可靠。真实性能就是，2CPU可以做每秒10万SET带落盘的水平。

*/

/*

业务逻辑测试(无并发)

容器3核3G：

条件1：:memory:模式，SyncPolicy=Never，Set10000
结果1：14毫秒

条件2：文件模式，SyncPolicy=Never，Set10000
结果2：4109毫秒

条件3：文件模式，SyncPolicy=EverySecond，Set10000
结果3：4611毫秒

真机：

内存-10ms
文件-不持久化-64ms
文件-每秒持久化-87ms

*/

/*

业务逻辑测试(并发)

真机2C4G:

10并发-文件模式-默认每秒-每个并发10000次SET-1098ms
100并发-文件模式-默认每秒-每个并发10000次SET-10557ms
100并发-文件模式-默认每秒-每个并发1次SET-2ms
100并发-文件模式-默认每秒-每个并发10次SET-14ms
1000并发-文件模式-默认每秒-每个并发10次SET-113ms
1000并发-文件模式-默认每秒-每个并发100次SET-1000ms
5000并发-文件模式-默认每秒-每个并发1次SET-114ms
5000并发-文件模式-默认每秒-每个并发10次SET-729ms

*/

// 结论：buntdb对串行的支持很强在2C就能达到10w次SET/s的性能，但是对并发的支持很弱，并发数量和耗时是1:1放大的关系。
// 不适合开放给多客户端，应该作为内部热存储开放rpc服务给其他服务访问，并且其他服务不能有大量的并发。或者给低并发服务比如日志系统加入一个串行化操作队列中间件，做持久化日志操作。
// 对游戏而言，理论上最大情况下，即使5000人/服同时做SET操作，也才耗时114ms。普通情况下100并发，每人1-10次SET也就是不到10ms的耗时。