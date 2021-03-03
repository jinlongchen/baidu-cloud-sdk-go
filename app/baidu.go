package main

import (
	"github.com/brickman-source/golang-utilities/baidu"
	"github.com/jinlongchen/baidu/context"
	"github.com/jinlongchen/baidu/http"
	"runtime"
	"time"

	"github.com/brickman-source/golang-utilities/banner"
	"github.com/brickman-source/golang-utilities/cache"
	"github.com/brickman-source/golang-utilities/config"
	"github.com/brickman-source/golang-utilities/database"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/svc"
	"github.com/brickman-source/golang-utilities/sync"
	"github.com/brickman-source/golang-utilities/version"
)

var (
	BINARY = "1.0"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	svc.Execute(NewService())
}

type Service struct {
	startTime   time.Time
	ctx         *context.Context
	httpHandler *http.Handler
	waitGroup   sync.WaitGroupWrapper
}

func NewService() *Service {
	return &Service{
		startTime: time.Now(),
	}
}

func (d *Service) Main(cfg *config.Config) {
	banner.Print(cfg.GetString("application.name"))
	log.Infof(version.String(cfg.GetString("application.name"), BINARY))

	ctx := &context.Context{
		Config: cfg,
	}
	connStrCfg := cfg.GetString("database.pg.connStr")
	if connStrCfg != "" {
		connection := &database.SQLConnection{
			DSN: connStrCfg,
		}
		ctx.Database = connection.GetDatabase()
	}

	redisAddr := cfg.GetString("database.redis.address")
	if redisAddr != "" {
		redisCache := cache.NewRedisCache(
			map[string]string{
				"redis1": redisAddr,
			},
			cfg.GetString("database.redis.password"),
		)
		ctx.Cache = redisCache
	}
	ctx.Baidu = baidu.NewBaidu(ctx.Cache, cfg, func(s string) {
		log.Infof("%s", s)
	})

	httpHandler := http.NewHttpHandler(ctx)
	d.waitGroup.Wrap(func() {
		httpHandler.Serve()
	})


	d.ctx = ctx
	d.httpHandler = httpHandler
}

func (d *Service) Exit() {
	if d.ctx.Database != nil {
		_ = d.ctx.Database.Close()
		d.ctx.Database = nil
	}
	if d.ctx.Cache != nil {
		_ = d.ctx.Cache.Close()
		d.ctx.Cache = nil
	}
	if d.ctx.Baidu != nil {
		_ = d.ctx.Baidu.Exit()
		d.ctx.Baidu = nil
	}
	d.httpHandler.Exit()
	d.waitGroup.Wait()
}

func (d *Service) GetName() string {
	return "wechat-service"
}
