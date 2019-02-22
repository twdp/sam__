package agent

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"net/http"
	"strconv"
	"tianwei.pro/micro/di/single"
)

const (
	SamUserInfoSessionKey = "__sam_user_info_key__"
	SamTokenCookieName    = "__sam_t__"
	SamTokenHeaderName    = "token"
)

// 查找路径和方法，是否需要登录或验证权限
// 如果需要登录或验证权限，则获取当前用户信息
// nil 则返回401  如果没权限，则返回403
// sam 的过滤器
var SamFilter = func(ctx *context.Context) {
	var urlStrategy int8 = Anonymous
	var id int64 = 0
	a := single.GetByName(samAgentName).(*agent)
	if _id, _strategy, err := a.CheckPermissionStrategy(ctx); err != nil {
		logs.Error("sam filter error: %v", err)
		ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
		ctx.ResponseWriter.Write([]byte(err.Error()))
		return
	} else {
		id = _id
		urlStrategy = _strategy
	}

	//if _, ok := ctx.Input.Session(SamUserInfoSessionKey).(*UserInfo); !ok {
	//	// 获取token信息
	//	token := ctx.Input.Header(SamTokenHeaderName)
	//	if token != "" {
	//		token = ctx.GetCookie(SamTokenCookieName)
	//	}
	//	if token != "" {
	//		// 根据token获取用户信息
	//		if us, err := a.verifyToken(token); err != nil || us.Err != nil {
	//			msg := ""
	//			if err != nil {
	//				msg = err.Error()
	//			} else {
	//				msg = us.Err.Error()
	//			}
	//			ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
	//			ctx.ResponseWriter.Write([]byte(msg))
	//			return
	//		} else {
	//			ctx.Output.Session(SamUserInfoSessionKey, us)
	//		}
	//	}
	//}

	if urlStrategy == Anonymous {
		return
	}

	var systemInfo *moduleInfo
	if s, err := a.loadSysInfo(); err != nil {
		logs.Error("load system info failed. strategy: Child")
		systemInfo = &moduleInfo{
			keepSign:       false,
			permissionType: Child,
			routes:         make(map[string][]*tree),
		}
	} else {
		systemInfo = s
	}

	var u *UserInfo

	if uu, ok := ctx.Input.Session(SamUserInfoSessionKey).(*UserInfo); !ok {

		// todo: del
		if !systemInfo.keepSign {
			ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			ctx.ResponseWriter.Write([]byte("请重新登录"))
			return
		}
		// 获取token信息
		token := ctx.Input.Header(SamTokenHeaderName)
		if token == "" {
			token = ctx.GetCookie(SamTokenCookieName)
		}

		if token == "" {
			ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			ctx.ResponseWriter.Write([]byte("请重新登录"))
			return
		}
		// 根据token获取用户信息
		if us, err := a.verifyToken(token); err != nil || us.Err != nil {
			msg := ""
			if err != nil {
				msg = err.Error()
			} else {
				msg = us.Err.Error()
			}
			ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			ctx.ResponseWriter.Write([]byte(msg))
			return
		} else {
			u = us
		}
		ctx.Output.Session(SamUserInfoSessionKey, u)
	} else {
		u = uu
	}

	if urlStrategy == OnlyNeedLogin {
		return
	}

	// owner 可以操作所有的操作
	if len(u.Permissions) > 0 && u.Permissions[0].RoleName == "owner" {
		return
	}

	permissionId := ctx.Input.Param("permissionId")
	if permissionId == "" {
		permissionId = ctx.Input.Param(":permissionId")
	}

	var ppId int64 = -1

	if ppid, err := strconv.ParseInt(permissionId, 10, 64); err == nil {
		ppId = ppid
	}

	hasPermission := false

	if u != nil {
		for _, p := range u.Permissions {
			if p.VerifyUrl(ppId, id, systemInfo.permissionType) {
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		// 403没权限
		ctx.ResponseWriter.WriteHeader(http.StatusForbidden)
		ctx.ResponseWriter.Write([]byte("暂无权限"))
	}

}
