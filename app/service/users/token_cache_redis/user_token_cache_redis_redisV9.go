package token_cache_redis

//
// import (
// 	"context"
// 	"strconv"
// 	"strings"
// 	"time"
//
// 	"ginskeleton/app/global/variable"
// 	"ginskeleton/app/model/web"
// 	"ginskeleton/app/utils/md5_encrypt"
// 	"ginskeleton/app/utils/redis_factory" // 在导入的时候初始化 init
//
// 	"github.com/opentracing/opentracing-go"
//
// 	"go.uber.org/zap"
// )
//
// func CreateUsersTokenCacheFactory(ctx context.Context, userId int64) *userTokenCacheRedis {
// 	redCli := redis_factory.GetOneRedisClient(ctx)
// 	if redCli == nil {
// 		return nil
// 	}
// 	return &userTokenCacheRedis{redisClient: &web.RedisClient{Client: redCli},
// 		userTokenKey: "token_userid_" + strconv.FormatInt(userId, 10)}
// }
//
// type userTokenCacheRedis struct {
// 	redisClient  *web.RedisClient
// 	userTokenKey string
// }
//
// // SetTokenCache 设置缓存
// func (u *userTokenCacheRedis) SetTokenCache(ctx context.Context, tokenExpire int64, token string) bool {
// 	// 存储用户token时转为MD5，下一步比较的时候可以更加快速地比较是否一致
// 	span := opentracing.SpanFromContext(ctx)
// 	newCtx := opentracing.ContextWithSpan(context.Background(), span)
//
// 	_, err := u.redisClient.ZAdd(newCtx, u.userTokenKey, tokenExpire, md5_encrypt.MD5(token))
// 	if err == nil {
// 		return true
// 	} else {
// 		variable.ZapLog.Error("缓存用户token到redis出错", zap.Error(err))
// 	}
// 	return false
// }
//
// // DelOverMaxOnlineCache 删除缓存,删除超过系统允许最大在线数量之外的用户
// func (u *userTokenCacheRedis) DelOverMaxOnlineCache(ctx context.Context) bool {
// 	// 首先先删除过期的token
// 	_, _ = u.redisClient.ZRemRangeByScore(ctx, u.userTokenKey, 0, time.Now().Unix()-1)
//
// 	onlineUsers := variable.ConfigYml.GetInt("Token.JwtTokenOnlineUsers")
// 	alreadyCacheNum, err := u.redisClient.ZCard(ctx, u.userTokenKey)
// 	if err == nil && alreadyCacheNum > int64(onlineUsers) {
// 		// 删除超过最大在线数量之外的token
// 		if alreadyCacheNum, err = u.redisClient.ZRemRangeByRank(ctx, u.userTokenKey, 0, alreadyCacheNum-int64(onlineUsers)-1); err == nil {
// 			return true
// 		} else {
// 			variable.ZapLog.Error("删除超过系统允许之外的token出错：", zap.Error(err))
// 		}
// 	}
// 	return false
// }
//
// // TokenCacheIsExists 查询token是否在redis存在
// func (u *userTokenCacheRedis) TokenCacheIsExists(ctx context.Context, token string) (exists bool) {
//
// 	curTimestamp := float64(time.Now().Unix())
// 	onlineUsers := variable.ConfigYml.GetInt64("Token.JwtTokenOnlineUsers")
//
// 	token = md5_encrypt.MD5(token) // 存入的时候就进行了 md5 加密
//
// 	if strSlice, err := u.redisClient.ZRevRange(ctx, u.userTokenKey, 0, onlineUsers-1); err == nil {
// 		for _, val := range strSlice {
// 			if score, err := u.redisClient.ZScore(ctx, u.userTokenKey, token); err == nil {
// 				if score > curTimestamp {
// 					if strings.Compare(val, token) == 0 {
// 						exists = true
// 						break
// 					}
// 				}
// 			}
// 		}
// 	} else {
// 		variable.ZapLog.Error("获取用户在redis缓存的 token 值出错：", zap.Error(err))
// 	}
// 	return
// }
//
// // SetUserTokenExpire 设置用户的 usertoken 键过期时间
// // 参数： 时间戳
// func (u *userTokenCacheRedis) SetUserTokenExpire(ctx context.Context, ts time.Time) bool {
// 	span := opentracing.SpanFromContext(ctx)
// 	newCtx := opentracing.ContextWithSpan(context.Background(), span)
// 	if _, err := u.redisClient.ExpireAt(newCtx, u.userTokenKey, ts); err == nil {
// 		return true
// 	}
// 	return false
// }
//
// // ClearUserToken 清除某个用户的全部缓存，当用户更改密码或者用户被禁用则删除该用户的全部缓存
// func (u *userTokenCacheRedis) ClearUserToken(ctx context.Context) bool {
// 	if _, err := u.redisClient.Del(ctx, u.userTokenKey); err == nil {
// 		return true
// 	}
// 	return false
// }
//
// // ReleaseRedisConn 释放redis
// // func (u *userTokenCacheRedis) ReleaseRedisConn() {
// // 	u.redisClient.ReleaseOneRedisClient()
// // }
