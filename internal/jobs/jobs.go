package jobs

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/tt90cc/utils/globalkey"
	"github.com/zeromicro/go-zero/core/service"
	"tt90.cc/ucenter/internal/svc"
)

func RegisterJobs(serverCtx *svc.ServiceContext) {
	crontab(serverCtx)
	queue(serverCtx)
}

func crontab(serverCtx *svc.ServiceContext) {
	c := cron.New(cron.WithSeconds())

	// 设置定时时间 秒 分 时 日 月 周
	if serverCtx.Config.Mode != service.DevMode {
		c.AddFunc("*/30 * * * * *", func() {
			lockKey := fmt.Sprintf(globalkey.RedisLock, "takeout.syncLock")
			tryLock := serverCtx.TryLock(lockKey, 30)
			if !tryLock {
				return
			}
			defer serverCtx.UnLock(lockKey)
			// todo do something
		})
	}

	c.Start()
}

func queue(serverCtx *svc.ServiceContext) {
	// todo consumer
}
