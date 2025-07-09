package session

import (
	"openscrm/common/log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

var Store sessions.Store

func Setup(redisHost, redisPassword, aesKey string) {
	var err error
	// 如果是占位符密码，则使用空字符串
	if redisPassword == "dummy" {
		redisPassword = ""
	}
	Store, err = redis.NewStore(10, "tcp", redisHost, redisPassword, []byte(aesKey))
	if err != nil {
		log.Sugar.Fatalw("setup session failed", "err", err)
		return
	}
}
