package initial

import (
	"github.com/noahhai/vigil/agent/services"
	"github.com/noahhai/vigil/agent/types"
)

func InitServices() *types.Services {
	return &types.Services{
		LogSvc: services.NewLogger(),
	}
}
