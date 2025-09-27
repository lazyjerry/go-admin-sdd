package file_store

import (
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Zone string

const (
	// HuaDong 华东
	HuaDong Zone = "HuaDong"
	// HuaBei 华北
	HuaBei Zone = "HuaBei"
	// HuaNan 华南
	HuaNan Zone = "HuaNan"
	// BeiMei 北美
	BeiMei Zone = "BeiMei"
	// XinJiaPo 新加坡
	XinJiaPo Zone = "XinJiaPo"
)

type QiNiuKODO struct {
	Client     interface{}
	BucketName string
	cfg        storage.Config
	options    []ClientOption
}

// QiNiuKODO 為七牛 KODO 的簡單封裝，用於處理上傳、取得臨時 Token 等操作。
// 方法說明：
// - getToken: 根據 PutPolicy 生成上傳 Token
// - Setup: 初始化 KODO 設定與憑證
// - setZoneORDefault: 設定機房 Zone，若未指定則使用預設值
// - UpLoad: 使用表單上傳將本地檔案上傳到指定對象名
// - GetTempToken: 返回臨時上傳 Token
func (e *QiNiuKODO) getToken() string {
	putPolicy := storage.PutPolicy{
		Scope: e.BucketName,
	}
	if len(e.options) > 0 && e.options[0]["Expires"] != nil {
		putPolicy.Expires = e.options[0]["Expires"].(uint64)
	}
	upToken := putPolicy.UploadToken(e.Client.(*qbox.Mac))
	return upToken
}

//Setup 裝載
//endpoint sss
func (e *QiNiuKODO) Setup(endpoint, accessKeyID, accessKeySecret, BucketName string, options ...ClientOption) error {

	mac := qbox.NewMac(accessKeyID, accessKeySecret)
	// 获取存储空间。
	cfg := storage.Config{}
	// 空间对应的机房
	e.setZoneORDefault(cfg, options...)
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	e.Client = mac
	e.BucketName = BucketName
	e.cfg = cfg
	e.options = options
	return nil
}

// setZoneORDefault 设置Zone或者默认华东
func (e *QiNiuKODO) setZoneORDefault(cfg storage.Config, options ...ClientOption) {
	if len(options) > 0 && options[0]["Zone"] != nil {
		if _, ok := options[0]["Zone"].(Zone); !ok {
			cfg.Zone = &storage.ZoneHuadong
		}
		switch options[0]["Zone"].(Zone) {
		case HuaDong:
			cfg.Zone = &storage.ZoneHuadong
		case HuaBei:
			cfg.Zone = &storage.ZoneHuabei
		case HuaNan:
			cfg.Zone = &storage.ZoneHuanan
		case BeiMei:
			cfg.Zone = &storage.ZoneBeimei
		case XinJiaPo:
			cfg.Zone = &storage.ZoneXinjiapo
		default:
			cfg.Zone = &storage.ZoneHuadong
		}
	}
}

// UpLoad 文件上传
func (e *QiNiuKODO) UpLoad(yourObjectName string, localFile interface{}) error {

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&e.cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, e.getToken(), yourObjectName, localFile.(string), &putExtra)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(ret.Key, ret.Hash)
	return nil
}

func (e *QiNiuKODO) GetTempToken() (string, error) {
	token := e.getToken()
	return token, nil
}
