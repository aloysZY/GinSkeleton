package redis

//
// import (
// 	"context"
// 	"strconv"
// 	"time"
//
// 	"github.com/redis/go-redis/v9"
// )
//
// // 定义一个redis客户端结构体
// type RedisClient struct {
// 	Client *redis.Client
// }
//
// // // 释放连接到连接池 redis.Conn是从连接池中取出的单个连接，除非你有特殊的需要，否则尽量不要使用它。你可以使用它向 redis 发送任何数据并读取 redis 的响应，当你使用完毕时，应该把它返回给 go-redis，否则连接池会永远丢失一个连 接。
// // func (r *RedisClient) ReleaseOneRedisClient() {
// // 	_ = r.Client.Close()
// // }
//
// // Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中。
// // 如果某个成员已经是有序集的成员，那么更新这个成员的分数值，并通过重新插入这个成员元素，来保证该成员在正确的位置上
// // 如果有序集合 key 不存在，则创建一个空的有序集并执行 ZADD 操作。
// func (r *RedisClient) ZAdd(ctx context.Context, tokenKey string, tokenExpire int64, token string) (int64, error) {
// 	return r.Client.ZAdd(ctx, tokenKey, redis.Z{Score: float64(tokenExpire), Member: token}).Result()
// }
//
// // Zremrangebyscore 命令用于移除有序集中，指定分数（score）区间内的所有成员。
// func (r *RedisClient) ZRemRangeByScore(ctx context.Context, tokenKey string, min, max int64) (int64, error) {
// 	return r.Client.ZRemRangeByScore(ctx, tokenKey, strconv.FormatInt(min, 10), strconv.FormatInt(max, 10)).Result()
// }
//
// // Redis Zcard 命令用于计算集合中元素的数量。
// func (r *RedisClient) ZCard(ctx context.Context, tokenKey string) (int64, error) {
// 	return r.Client.ZCard(ctx, tokenKey).Result()
// }
//
// // Zremrangebyrank 命令用于移除有序集中，指定排名(rank)区间内的所有成员。
// func (r *RedisClient) ZRemRangeByRank(ctx context.Context, tokenKey string, min, max int64) (int64, error) {
// 	return r.Client.ZRemRangeByRank(ctx, tokenKey, min, max).Result()
// }
//
// // Redis Zrevrange 命令返回有序集中，指定区间内的成员。
// // 其中成员的位置按分数值递减(从大到小)来排列。
// func (r *RedisClient) ZRevRange(ctx context.Context, tokenKey string, min, max int64) ([]string, error) {
// 	return r.Client.ZRevRange(ctx, tokenKey, min, max).Result()
// }
//
// // Redis Zscore 命令返回有序集中，成员的分数值。 如果成员元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
// func (r *RedisClient) ZScore(ctx context.Context, tokenKey, token string) (float64, error) {
// 	return r.Client.ZScore(ctx, tokenKey, token).Result()
// }
//
// // Redis Expireat 命令用于以 UNIX 时间戳(unix timestamp)格式设置 key 的过期时间。key 过期后将不再可用。
// func (r *RedisClient) ExpireAt(ctx context.Context, tokenKey string, expireAt time.Time) (bool, error) {
// 	return r.Client.ExpireAt(ctx, tokenKey, expireAt).Result()
// }
//
// // Redis DEL 命令用于删除已存在的键。不存在的 key 会被忽略。
// func (r *RedisClient) Del(ctx context.Context, tokenKey string) (int64, error) {
// 	return r.Client.Del(ctx, tokenKey).Result()
// }
