package redis_factory

//
// import (
// 	"encoding/json"
// 	"net"
// 	"time"
//
// 	"ginskeleton/app/core/event_manage"
// 	"ginskeleton/app/global/my_errors"
// 	"ginskeleton/app/global/variable"
// 	"ginskeleton/app/utils/yml_config"
// 	"ginskeleton/app/utils/yml_config/ymlconfig_interf"
//
// 	redigo "github.com/gomodule/redigo/redis"
//
// 	"go.uber.org/zap"
//
// 	"github.com/go-redis/redis"
// )
//
// var redisPool *redis.Client
// var configYml ymlconfig_interf.YmlConfigInterf
//
// // 定义一个redis客户端结构体
// type RedisClient struct {
// 	client *redis.Client
// }
//
// // 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
// func init() {
// 	configYml = yml_config.CreateYamlFactory()
// 	redisPool = initRedisClientPool()
// }
//
// func initRedisClientPool() *redis.Client {
//
// 	redisPool = redis.NewClient(&redis.Options{
// 		// 连接信息
// 		Network:  "tcp",                                                                       // 网络类型，tcp or unix，默认tcp
// 		Addr:     configYml.GetString("Redis.Host") + ":" + configYml.GetString("Redis.Port"), // 主机名+冒号+端口，默认localhost:6379
// 		Password: configYml.GetString("Redis.Auth"),                                           // 密码
// 		DB:       configYml.GetInt("Redis.IndexDb"),                                           // redis数据库index
//
// 		// 连接池容量及闲置连接数量
// 		PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
// 		MinIdleConns: 10, // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
//
// 		DialTimeout:  5 * time.Second, // 连接建立超时时间，默认5秒。
// 		ReadTimeout:  3 * time.Second, // 读超时，默认3秒， -1表示取消读超时
// 		WriteTimeout: 3 * time.Second, // 写超时，默认等于读超时
// 		PoolTimeout:  4 * time.Second, // 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒
//
// 		// 闲置连接检查包括IdleTimeout，MaxConnAge
// 		IdleCheckFrequency: 60 * time.Second, // 闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
// 		IdleTimeout:        5 * time.Minute,  // 闲置超时，默认5分钟，-1表示取消闲置超时检查
// 		MaxConnAge:         0 * time.Second,  // 连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
//
// 		// 命令执行失败时的重试策略
// 		MaxRetries:      3,                      // 命令执行失败时，最多重试多少次，默认为0即不重试，默认 3
// 		MinRetryBackoff: 8 * time.Millisecond,   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
// 		MaxRetryBackoff: 512 * time.Millisecond, // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
//
// 		// 可自定义连接函数
// 		Dialer: func() (net.Conn, error) {
// 			netDialer := &net.Dialer{
// 				Timeout:   5 * time.Second,
// 				KeepAlive: 5 * time.Minute,
// 			}
// 			// 此处对应redis ip及端口号
// 			conn, err := netDialer.Dial("tcp", configYml.GetString("Redis.Host")+":"+configYml.GetString("Redis.Port"))
// 			if err != nil {
// 				variable.ZapLog.Error(my_errors.ErrorsRedisInitConnFail + err.Error())
// 				return nil, err
// 			}
// 			return conn, err
// 		},
//
// 		// 钩子函数
// 		OnConnect: func(conn *redis.Conn) error { // 仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
// 			// fmt.Printf("conn=%v\n", conn)
// 			return nil
// 		},
// 	})
//
// 	auth := configYml.GetString("Redis.Auth") // 通过配置项设置redis密码
// 	if len(auth) >= 1 {
// 		if cmd := redisPool.Do("AUTH", auth); cmd.Err() != nil {
// 			_ = redisPool.Close()
// 			variable.ZapLog.Error(my_errors.ErrorsRedisAuthFail + cmd.Err().Error())
// 		}
// 	}
// 	_ = redisPool.Do("select", configYml.GetInt("Redis.IndexDb"))
//
// 	// 将redis的关闭事件，注册在全局事件统一管理器，由程序退出时统一销毁
// 	eventManageFactory := event_manage.CreateEventManageFactory()
// 	if _, exists := eventManageFactory.Get(variable.EventDestroyPrefix + "Redis"); exists == false {
// 		eventManageFactory.Set(variable.EventDestroyPrefix+"Redis", func(args ...interface{}) {
// 			_ = redisPool.Close()
// 		})
// 	}
// 	return redisPool
// }
//
// // 从连接池获取一个redis连接
// func GetOneRedisClient() *RedisClient {
// 	maxRetryTimes := configYml.GetInt("Redis.ConnFailRetryTimes")
// 	for i := 1; i <= maxRetryTimes; i++ {
// 		// oneConn = redisPool.Get()
// 		// 首先通过执行一个获取时间的命令检测连接是否有效，如果已有的连接无法执行命令，则重新尝试连接到redis服务器获取新的连接池地址
// 		// 连接不可用可能会发生的场景主要有：服务端redis重启、客户端网络在有线和无线之间切换等
// 		if _, err := redisPool.Do("time").Result(); err != nil {
// 			// fmt.Printf("连接已经失效(出错)：%+v\n", replyErr.Error())
// 			// 如果已有的redis连接池获取连接出错(官方库的说法是连接不可用)，那么继续使用从新初始化连接池
// 			initRedisClientPool()
// 		} else if i == maxRetryTimes {
// 			// variable.ZapLog.Error("Redis：网络中断,开始重连进行中..." , zap.Error(oneConn.Err()))
// 			variable.ZapLog.Error(my_errors.ErrorsRedisGetConnFail, zap.Error(err))
// 			return nil
// 		} else if err == nil {
// 			break
// 		}
// 		// 如果出现网络短暂的抖动，短暂休眠后，支持自动重连
// 		time.Sleep(time.Second * configYml.GetDuration("Redis.ReConnectInterval"))
// 	}
// 	return &RedisClient{redisPool}
// }
//
// // 为redis-go 客户端封装统一操作函数入口
// func (r *RedisClient) Execute(cmd string, args ...interface{}) (interface{}, error) {
// 	marshal, err := json.Marshal(args)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return r.client.Do(cmd, marshal).Result()
// }
//
// // 释放连接到连接池
// func (r *RedisClient) ReleaseOneRedisClient() {
// 	_ = r.client.Close()
// }
//
// //  封装几个数据类型转换的函数
//
// // bool 类型转换
// func (r *RedisClient) Bool(reply interface{}, err error) (bool, error) {
// 	return redigo.Bool(reply, err)
// }
//
// // string 类型转换
// func (r *RedisClient) String(reply interface{}, err error) (string, error) {
// 	return redigo.String(reply, err)
// }
//
// // string map 类型转换
// func (r *RedisClient) StringMap(reply interface{}, err error) (map[string]string, error) {
// 	return redigo.StringMap(reply, err)
// }
//
// // strings 类型转换
// func (r *RedisClient) Strings(reply interface{}, err error) ([]string, error) {
// 	return redigo.Strings(reply, err)
// }
//
// // Float64 类型转换
// func (r *RedisClient) Float64(reply interface{}, err error) (float64, error) {
// 	return redigo.Float64(reply, err)
// }
//
// // int 类型转换
// func (r *RedisClient) Int(reply interface{}, err error) (int, error) {
// 	return r.client.int
// }
//
// // int64 类型转换
// func (r *RedisClient) Int64(reply interface{}, err error) (int64, error) {
// 	return redigo.Int64(reply, err)
// }
//
// // int map 类型转换
// func (r *RedisClient) IntMap(reply interface{}, err error) (map[string]int, error) {
// 	return redigo.IntMap(reply, err)
// }
//
// // Int64Map 类型转换
// func (r *RedisClient) Int64Map(reply interface{}, err error) (map[string]int64, error) {
// 	return redigo.Int64Map(reply, err)
// }
//
// // int64s 类型转换
// func (r *RedisClient) Int64s(reply interface{}, err error) ([]int64, error) {
// 	return redigo.Int64s(reply, err)
// }
//
// // uint64 类型转换
// func (r *RedisClient) Uint64(reply interface{}, err error) (uint64, error) {
// 	return redigo.Uint64(reply, err)
// }
//
// // Bytes 类型转换
// func (r *RedisClient) Bytes(reply interface{}, err error) ([]byte, error) {
// 	return redigo.Bytes(reply, err)
// }
//
// // 以上封装了很多最常见类型转换函数，其他您可以参考以上格式自行封装
