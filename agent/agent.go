// +build !server

package agent

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"sync"
	"tianwei.pro/micro/conf"
	"tianwei.pro/micro/di/single"
	"time"
)

const (
	tokenKey   = "_token_"
	systemInfo = "_system_info_"
	samAgentName = "sam_agent"
)

func init() {
	appKey := beego.AppConfig.String("appKey")
	if appKey == "" {
		logs.Warn("app key is empty")
	}
	secret := beego.AppConfig.String("secret")
	if secret == "" {
		logs.Warn("secret is empty")
	}
	a := &agent{
		cacheManager: make(map[string]cache.Cache),
		appKey:       appKey,
		secret:       secret,
	}
	single.Provide(samAgentName, a)
}

type tree struct {
	beego.Tree
	id   int64
	Type int8
}

type moduleInfo struct {
	routes         map[string][]*tree
	permissionType int8
	keepSign       bool
}

type agent struct {
	sync.Mutex

	appKey string
	secret string

	cacheManager map[string]cache.Cache

	SamAgent *SamAgentFacade `inject:"rpc_c_samAgentFacade"`
}

func (a *agent) loadCacheByKey(key string) cache.Cache {
	if c, exist := a.cacheManager[key]; exist {
		return c
	} else {
		a.Lock()
		defer a.Unlock()
		if v, ok := a.cacheManager[key]; ok {
			return v
		}

		c, err := cache.NewCache(conf.Conf.DefaultString("cache.name", "memory"), conf.Conf.DefaultString("cache.conf", `{"interval":60}`))
		if err != nil {
			logs.Error("create cache fail. err: %v", err)
			panic(err)
		}
		a.cacheManager[key] = c
		return c
	}
}

func (a *agent) verifyToken(token string) (*UserInfo, error) {
	samAgentFacade := a.getAgentFacade()
	cache := a.loadCacheByKey(tokenKey)
	if cache.IsExist(token) {
		return cache.Get(token).(*UserInfo), nil
	} else {
		u := &UserInfo{}
		if user, err := samAgentFacade.VerifyToken(&VerifyTokenParam{
			SystemInfoParam: SystemInfoParam{
				AppKey: a.appKey,
				Secret: a.secret,
			},
			Token: token,
		}); err != nil {
			return u, err
		} else {
			u = user
		}

		//cache.Put(token, u, time.Minute)
		return u, nil
	}
}

func (a *agent) getAgentFacade() *SamAgentFacade {
	if a.SamAgent == nil {
		panic("sam agent facade not found")
	}
	return a.SamAgent
}


func (a *agent) loadSysInfo() (*moduleInfo, error) {
	samAgentFacade := a.getAgentFacade()

	cache := a.loadCacheByKey(systemInfo)
	if cache.IsExist("---") {
		return cache.Get("---").(*moduleInfo), nil
	} else {
		s := &SystemInfo{}
		if ss, err := samAgentFacade.LoadSystemInfo(&SystemInfoParam{
			AppKey: a.appKey,
			Secret: a.secret,
		}); err != nil {
			return nil, err
		} else {
			s = ss
		}

		routes := make(map[string][]*tree)
		for k := range beego.HTTPMETHOD {
			routes[k] = []*tree{}
		}

		for _, v := range s.Routers {
			tt := beego.NewTree()
			tt.AddRouter(v.Url, "sam")
			t := &tree{
				Tree: *tt,
				id:   v.Id,
				Type: v.Type,
			}
			routes[v.Method] = append(routes[v.Method], t)
		}
		ss := &moduleInfo{
			permissionType: s.PermissionType,
			keepSign:       s.KeepSign,
			routes:         routes,
		}

		cache.Put("---", ss, 10 * time.Minute)
		return ss, nil

	}
}

// beego.HTTPMETHOD
//
// @return int64  id
// @return string 正则表达的url
// @return string   method
// @return strategy Anonymous\OnlyNeedLogin\CheckRolePermission
// @return error  -> 验证时报错
func (a *agent) CheckPermissionStrategy(ctx *context.Context) (int64, int8, error) {
	method := ctx.Input.Method()
	path := ctx.Input.URL()

	if s, err := a.loadSysInfo(); err != nil {
		return 0, CheckRolePermission, err
	} else {
		var tree *tree
		routers := s.routes[method]
		for _, r := range routers {
			obj := r.Match(path, ctx)
			if obj != nil && obj.(string) == "sam" {
				tree = r
				break
			}
		}
		if tree == nil {
			return 0, Anonymous, nil
		}

		return tree.id, tree.Type, nil
	}

}

