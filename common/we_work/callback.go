package we_work

import (
	"openscrm/common/log"
	"openscrm/conf"
	work_wx "openscrm/pkg/easywework"
)

var Callback *work_wx.CallBackHandler

func SetupWXCallback() {
	// 在开发环境或配置不完整时跳过企业微信回调初始化
	if conf.Settings.WeWork.CallbackToken == "todo" || conf.Settings.WeWork.CallbackAesKey == "todo" {
		log.Sugar.Warn("企业微信回调配置不完整，跳过回调处理器初始化")
		return
	}

	var err error
	Callback, err = work_wx.NewCBHandler(conf.Settings.WeWork.CallbackToken, conf.Settings.WeWork.CallbackAesKey)
	if err != nil {
		log.Sugar.Errorf("企业微信回调处理器初始化失败: %v", err)
		// 在开发环境不要panic，只记录错误
		if conf.Settings.App.Env == "DEMO" || conf.Settings.App.Env == "DEV" || conf.Settings.App.Env == "TEST" {
			log.Sugar.Warn("开发环境跳过企业微信回调初始化失败")
			return
		}
		panic("init callback handler failed")
	}
}
