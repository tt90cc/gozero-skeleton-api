package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"tt90.cc/ucenter/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	// UcenterRpc ucenter.Ucenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	//conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config: c,
		Redis: redis.New(c.RedisConf.Host, func(r *redis.Redis) {
			r.Type = c.RedisConf.Type
			r.Pass = c.RedisConf.Pass
		}),
		// UcenterRpc: ucenter.NewUcenter(zrpc.MustNewClient(c.UcenterRpc)),
	}
}

func (s *ServiceContext) TryLock(key string, second int) bool {
	if ok, err := s.Redis.SetnxEx(key, "lock", second); !ok || err != nil {
		return false
	}
	return true
}

func (s *ServiceContext) UnLock(key string) {
	s.Redis.Del(key)
}
