package queue

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/tt90cc/utils/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"log"
	"time"
	"tt90.cc/ucenter/internal/svc"
)

const Prefix string = "user:"
const QName string = "user."

const (
	//CancelOrder       = Prefix + "cancelOrder"
)

// CustomerProcessor queue下面所有的方法都需要实现的类
type CustomerProcessor interface {
	SendTask(params Payload) (*asynq.Task, error)
	ProcessTask(ctx context.Context, t *asynq.Task) error
}

// Payload 统一queue参数
type Payload struct {
	Params    interface{}   `json:"params,omitempty"`
	ProcessIn time.Duration `json:"process_in,omitempty"`
}

// Queue 队列实例
type Queue struct {
	serverCtx *svc.ServiceContext
	client    *asynq.Client
	opt       asynq.RedisConnOpt
	logger    logx.Logger
}

// NewQueue 实例化queue
func NewQueue(serverCtx *svc.ServiceContext, ctx context.Context) *Queue {
	var opt asynq.RedisConnOpt
	opt = asynq.RedisClientOpt{
		Addr:     serverCtx.Config.RedisConf.Host,
		Password: serverCtx.Config.RedisConf.Pass,
		DB:       0,
	}

	return &Queue{
		serverCtx: serverCtx,
		client:    asynq.NewClient(opt),
		opt:       opt,
		logger:    logx.WithContext(ctx),
	}
}

// StartServer 启动监听服务
func (s Queue) StartServer() {
	threading.GoSafe(func() {
		srv := asynq.NewServer(
			s.opt,
			asynq.Config{
				// 每个进程并发执行的worker数量
				Concurrency: 5,
				// Optionally specify multiple queues with different priority.
				Queues: map[string]int{
					QName + "default":  3,
				},
			},
		)
		mux := asynq.NewServeMux()
		//background := context.Background()
		//mux.Handle(CancelOrder, NewCancelOrderProcessor(s.serverCtx, &background))
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	})
}

// SendQueue 发送队列
func (s Queue) SendQueue(obj CustomerProcessor, params Payload) error {
	task, err := obj.SendTask(params)
	if err != nil {
		logx.Errorf("could not create task: %v", err)
	}
	info, err := s.client.Enqueue(task, asynq.Queue(QName+"default"))
	if err != nil {
		return errorx.NewCodeError(errorx.ERR_DEFAULT, err.Error())
	}
	s.logger.Infof("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return nil
}
