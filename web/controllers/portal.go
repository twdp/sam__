package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
	"tianwei.pro/business"
	"tianwei.pro/business/controller"
	"tianwei.pro/sam/agent"
	"tianwei.pro/sam/core/dto/req"
	"tianwei.pro/sam/core/facade"
)

var cpt *captcha.Captcha

func init() {
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/api/v1/captcha", store)
}

type PortalController struct {
	controller.RestfulController
}

// @router /email [post]
func (p *PortalController) LoginByEmail() {
	email := p.GetString("email")
	password := p.GetString("password")

	if email == "" || password == "" {
		p.SetSession("needCpt", true)
		p.E500(business.H{
			"msg":     "账号或密码不能为空",
			"needCpt": true,
		})
		return
	}

	reply, err := facade.RpcUser.Login(&req.EmailLoginDto{
		Email:    email,
		Password: password,
	})
	//business.CheckError(&reply.Response, err)
	if err != nil {
		panic(err)
	}

	if reply.Err != "" {
		count := p.GetSession("passwordWrongCount")
		c := 1
		if count == nil {
			p.SetSession("passwordWrongCount", c)
		} else {
			c = count.(int)
			p.SetSession("passwordWrongCount", c+1)
		}
		result := make(map[string]interface{})
		result["msg"] = reply.Err
		if c > 2 {
			result["needCpt"] = true
		}
		p.E500(result)
		return
	}

	param := &agent.VerifyTokenParam{
		SystemInfoParam: agent.SystemInfoParam{
			AppKey: beego.AppConfig.String("appkey"),
			Secret: beego.AppConfig.String("secret"),
		},
		Token: reply.Token,
	}

	userInfo, err := facade.RpcSamAgent.VerifyToken(param)
	business.CheckError(&userInfo.Response, err)

	p.SetSession(agent.SamUserInfoSessionKey, userInfo)

	// 跨 域Secure必须为false
	p.Ctx.SetCookie(agent.SamTokenCookieName, reply.Token, 30*24*60*60, "/", beego.AppConfig.DefaultString("topDomain", "/"), false, true)

	p.ReturnJson(business.H{
		"token": reply.Token,
	})
}
