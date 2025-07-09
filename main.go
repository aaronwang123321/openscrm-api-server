package main

import (
	"fmt"
	"log"
	"net/http"
	"openscrm/app/constants"
	"openscrm/app/consumers"
	"openscrm/app/models"
	"openscrm/app/services"
	"openscrm/app/tasks"
	"openscrm/common/app"
	"openscrm/common/delay_queue"
	"openscrm/common/id_generator"
	log2 "openscrm/common/log"
	"openscrm/common/redis"
	"openscrm/common/session"
	"openscrm/common/storage"
	"openscrm/common/validator"
	we_work2 "openscrm/common/we_work"
	"openscrm/conf"
	"openscrm/routers"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	val "github.com/go-playground/validator/v10"
	"golang.org/x/net/context"
)

func init() {
	err := conf.SetupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	validateConfig(conf.Settings)
	setupTimezone()
	setupValidator()
	log2.SetupLogger(conf.Settings.App.Env)
	id_generator.SetupIDGenerator()
	redis.Setup(conf.Settings.Redis.Host, conf.Settings.Redis.Password, conf.Settings.Redis.DBNumber)
	models.SetupDB()
	storage.Setup(conf.Settings.Storage)
	session.Setup(conf.Settings.Redis.Host, conf.Settings.Redis.Password, conf.Settings.App.Key)
	we_work2.SetupWXCallback()
	we_work2.SetupClient(we_work2.CorpConf{
		ExtCorpID:       conf.Settings.WeWork.ExtCorpID,
		ContactSecret:   conf.Settings.WeWork.ContactSecret,
		CustomerSecret:  conf.Settings.WeWork.CustomerSecret,
		MainAgentID:     conf.Settings.WeWork.MainAgentID,
		MainAgentSecret: conf.Settings.WeWork.MainAgentSecret})
	delay_queue.SetupDelayQueue()
	consumers.Start()
	tasks.Start()
	models.SetupIDConverter()

	if conf.Settings.App.AutoSyncWeWorkData {
		// 检查企业微信配置是否完整
		if conf.Settings.WeWork.ContactSecret == "todo" ||
			conf.Settings.WeWork.CustomerSecret == "todo" ||
			conf.Settings.WeWork.MainAgentSecret == "todo" {
			log2.Sugar.Warn("企业微信API配置不完整，跳过自动同步数据")
		} else {
			services.Syncs()
		}
	}
}

// @title OpenSCRM
// @version 1.0
// @description  企微私域流量管理
// @termsOfService https://ixj.cn
//
//go:generate swag init
func main() {
	gin.SetMode(conf.Settings.Server.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", conf.Settings.Server.HttpHost, conf.Settings.Server.HttpPort),
		Handler:        router,
		ReadTimeout:    conf.Settings.Server.ReadTimeout,
		WriteTimeout:   conf.Settings.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()

	log.Printf("Api Server Listen on %s", s.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go consumers.Stop()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupValidator() {
	v := validator.NewCustomValidator()
	binding.Validator = v
	app.NewBindingValidator(&v.Trans)
}

// setupTimezone设置time包默认时区为北京时间
func setupTimezone() {
	time.Local = constants.PRCLocation
}

func validateConfig(c interface{}) {
	if err := validator.NewCustomValidator().ValidateStruct(c); err != nil {
		panic(err.(val.ValidationErrors))
	}
}
