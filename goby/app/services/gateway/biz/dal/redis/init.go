package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/gateway/conf"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

type Z = redis.Z

// Initialization
// The InitRedis function initializes the Redis client and checks if the connection is successful using Ping.
func Init() {
	conf_ := conf.GetConf()
	dsn := fmt.Sprintf(conf_.Redis.DSN, conf_.Redis.Username, conf_.Redis.Password, conf_.Redis.Host, conf_.Redis.Port)
	opt, _ := redis.ParseURL(dsn)
	RedisClient = redis.NewClient(opt)
	// 测试连接
	if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
		hlog.Fatalf("Redis连接失败: %v", err)
	}
	hlog.Info("Redis连接成功")
}

// Graceful Shutdown
// The CloseRedis function ensures the Redis connection is properly closed when it's no longer needed.
func Close() {
	if RedisClient != nil {
		_ = RedisClient.Close()
	}
}

// 从数据库加载黑名单到Redis
func SyncBlacklistToRedis() {
	ctx := context.Background()

	// 删除旧数据（可选，根据需求决定）
	// _ = redisClient.Del(ctx, "blacklist:entries", "blacklist:expiry").Err()

	// 查询有效黑名单条目
	var entries []models.Blacklist
	if err := tidb.DB.Where("(expires_at > ? OR expires_at IS NULL) AND is_deleted = 0", time.Now()).Find(&entries).Error; err != nil {
		hlog.Errorf("黑名单同步失败: %v", err)
		return
	}

	// 使用Pipeline批量写入Redis
	pipe := RedisClient.Pipeline()
	for _, entry := range entries {
		// 写入Hash
		data, _ := json.Marshal(map[string]interface{}{
			"reason":     entry.Reason,
			"expires_at": entry.ExpiresAt,
		})
		pipe.HSet(ctx, "blacklist:entries", entry.Identifier, data)

		// 写入Sorted Set（仅有过期时间的条目）
		if !entry.ExpiresAt.IsZero() {
			pipe.ZAdd(ctx, "blacklist:expiry", &redis.Z{
				Score:  float64(entry.ExpiresAt.UnixNano()),
				Member: entry.Identifier,
			})
		}
	}

	// 执行Pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		hlog.Errorf("Redis同步失败: %v", err)
		return
	}

	hlog.Info("黑名单同步到Redis完成")
}
