# OpenSCRM 企业微信配置指南

## 🎯 配置企业微信信息

### 第一步：获取企业微信信息
1. 登录企业微信管理后台：https://work.weixin.qq.com
2. 获取以下信息：

#### 企业基本信息
- **企业ID (ExtCorpID)**: 在【我的企业】→【企业信息】中查看

#### 应用信息
- **通讯录Secret**: 在【管理工具】→【通讯录同步】中查看
- **客户联系Secret**: 在【客户联系】→【API】中查看  
- **自建应用ID & Secret**: 在【应用管理】→【自建】中创建应用获取

#### 回调配置
- **CallbackToken**: 自定义32位字符串
- **CallbackAesKey**: 自定义43位字符串

### 第二步：修改配置文件

编辑 `api-server/conf/config.yaml`：

```yaml
WeWork:
  # 企业ID (必填)
  ExtCorpID: "你的企业ID"
  # 通讯录Secret (必填)  
  ContactSecret: "你的通讯录Secret"
  # 客户联系Secret (必填)
  CustomerSecret: "你的客户联系Secret"
  # 自建应用ID (必填)
  MainAgentID: 你的应用ID
  # 自建应用Secret (必填)
  MainAgentSecret: "你的应用Secret"
  # 回调Token (必填)
  CallbackToken: "你自定义的Token"
  # 回调AesKey (必填)  
  CallbackAesKey: "你自定义的AesKey"
  # 会话存档私钥路径
  PriKeyPath: /conf/private.key
```

### 第三步：重启API服务

```bash
# 如果是Docker部署
docker-compose restart api-server

# 如果是本地运行
./api-server
```

### 第四步：验证配置

访问Dashboard，在员工管理页面点击"同步数据"按钮，如果同步成功说明配置正确。

## ⚠️ 注意事项

1. **所有Secret都不能是"todo"** - 必须填入真实值
2. **企业ID格式**: 通常以"ww"开头，18位字符
3. **回调地址**: 需要在企业微信后台配置对应的回调URL
4. **网络访问**: 确保服务器能访问企业微信API

## 🔗 企业微信配置参考链接

- 企业微信管理后台: https://work.weixin.qq.com/wework_admin/frame#profile
- 通讯录API: https://work.weixin.qq.com/wework_admin/frame#apps/contactsApi  
- 客户联系API: https://work.weixin.qq.com/wework_admin/frame#customer/analysis
- 会话存档: https://work.weixin.qq.com/wework_admin/frame#financial/corpEncryptData

