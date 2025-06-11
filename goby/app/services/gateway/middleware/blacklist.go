package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/redis"

	// "github.com/bitdance-panic/gobuy/app/services/gateway/biz/dao"
	"github.com/bitdance-panic/gobuy/app/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	Redis "github.com/go-redis/redis/v8"
)

var (
	BlacklistCache sync.Map     // 内存缓存：map[identifier]expiryTime
	CacheMutex     sync.RWMutex // 缓存读写锁
	LastSyncTime   time.Time    // 最后同步时间
)

// 黑名单中间件
func BlacklistMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取用户标识（示例：优先取用户ID，未登录则取IP）
		var identifier string
		// if userID, exists := c.Get(app_consts.CONTEXT_UID_KEY); exists {
		// 	fmt.Println("tokenaa", userID)
		// 	identifier = fmt.Sprintf("user:%d", userID.(int))
		// } else {
		identifier = fmt.Sprintf("ip:%s", c.ClientIP())
		// }
		// 检查Redis黑名单
		if isBlocked, err := CheckBlockedInRedis(identifier); err != nil {
			hlog.Errorf("Redis查询失败: %v", err)
			c.AbortWithStatus(consts.StatusInternalServerError)
			return
		} else if isBlocked {
			hlog.Warnf("拒绝黑名单用户访问: %s", identifier)
			utils.FailFull(c, consts.StatusForbidden, "您的账户或设备已被封禁", nil)
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}

// 检查是否在黑名单中
func CheckBlockedInRedis(identifier string) (bool, error) {
	ctx := context.Background()

	// 1. 检查Hash是否存在
	data, err := redis.RedisClient.HGet(ctx, "blacklist:entries", identifier).Result()
	if err == Redis.Nil {
		return false, nil // 不在黑名单中
	} else if err != nil {
		return false, err
	}

	// 2. 解析数据
	var entry struct {
		Reason    string    `json:"reason"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return false, err
	}

	// 3. 检查是否过期
	if !entry.ExpiresAt.IsZero() && entry.ExpiresAt.Before(time.Now()) {
		// 异步删除过期条目
		go func() {
			_ = redis.RedisClient.HDel(ctx, "blacklist:entries", identifier).Err()
			_ = redis.RedisClient.ZRem(ctx, "blacklist:expiry", identifier).Err()
		}()
		return false, nil
	}

	return true, nil
}

// 启动定时任务清理过期条目（利用Redis的Sorted Set）
func StartRedisCleanupTask() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			ctx := context.Background()
			now := time.Now().UnixNano()

			// 查询所有过期的标识符
			expiredIDs, err := redis.RedisClient.ZRangeByScore(ctx, "blacklist:expiry", &Redis.ZRangeBy{
				Min: "0",
				Max: fmt.Sprintf("%d", now),
			}).Result()
			if err != nil {
				hlog.Errorf("Redis清理查询失败: %v", err)
				continue
			}

			// 批量删除
			if len(expiredIDs) > 0 {
				pipe := redis.RedisClient.Pipeline()
				pipe.HDel(ctx, "blacklist:entries", expiredIDs...)
				pipe.ZRemRangeByScore(ctx, "blacklist:expiry", "0", fmt.Sprintf("%d", now))
				if _, err := pipe.Exec(ctx); err != nil {
					hlog.Errorf("Redis清理失败: %v", err)
				} else {
					hlog.Infof("清理 %d 个过期黑名单条目", len(expiredIDs))
				}
			}
		}
	}()
}
