package service

import (
	"fmt"
	"giligili/serializer"
	"os"
	"time"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// UploadTokenService 获取 oss 上传服务的 token
type UploadTokenService struct {
	Filename string `form:"filename" json:"filename"`
}

// Update 更新视频
func (service *UploadTokenService) Post() serializer.Response {
	appid := os.Getenv("TENCENT_OSS_APP_ID")
	bucket := os.Getenv("TENCENT_OSS_BUCKET")
	region := os.Getenv("TENCENT_OSS_REGION")
	c := sts.NewClient(
		// 通过环境变量获取密钥, os.Getenv 方法表示获取环境变量
		os.Getenv("TENCENT_OSS_SECRET_ID"),  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		os.Getenv("TENCENT_OSS_SECRET_KEY"), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		nil,
	)
	// 策略概述 https://cloud.tencent.com/document/product/436/18023
	opt := &sts.CredentialOptions{
		DurationSeconds: int64((5 * time.Minute).Seconds()), // 过期时间：5分钟
		Region:          region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						// 简单上传
						"name/cos:PostObject",
						"name/cos:PutObject",
					},
					Effect: "allow",
					Resource: []string{
						// 这里改成允许的路径前缀，可以根据自己网站的用户登录态判断允许上传的具体路径，例子： a.jpg 或者 a/* 或者 * (使用通配符*存在重大安全风险, 请谨慎评估使用)
						// 存储桶的命名格式为 BucketName-APPID，此处填写的 bucket 必须为此格式
						"qcs::cos:" + region + ":uid/" + appid + ":" + bucket + "/upload/*",
					},
				},
			},
		},
	}

	res, err := c.GetCredential(opt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res.Credentials)

	return serializer.Response{
		Data: res,
	}
}
