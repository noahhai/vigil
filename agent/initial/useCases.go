package initial

import (
	"github.com/noahhai/vigil/agent/types"
	"github.com/noahhai/vigil/agent/useCases"
)

func InitUseCases(svc *types.Services) *types.UseCases {
	return &types.UseCases{
		CommandExec: useCases.GetCmdExe(svc),
	}
}
