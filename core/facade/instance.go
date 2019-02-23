package facade

import (
	"tianwei.pro/micro/di"
	"tianwei.pro/micro/di/single"
	"tianwei.pro/sam/agent"
)

var (
	RpcSamAgent *agent.SamAgentFacade
	RpcUser     *RpcUserFacade
	RpcSystem   *RpcSystemFacade
)

func init() {
	RpcSamAgent = &agent.SamAgentFacade{}
	RpcUser = &RpcUserFacade{}
	RpcSystem = &RpcSystemFacade{}
	single.Provide(di.NewRpcConsumerName("samAgentFacade"), RpcSamAgent)
	single.Provide(di.NewRpcConsumerName("userFacade"), RpcUser)
	single.Provide(di.NewRpcConsumerName("systemFacade"), RpcSystem)
}
