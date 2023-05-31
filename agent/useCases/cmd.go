package useCases

import (
	"github.com/mitchellh/go-ps"
	"github.com/noahhai/vigil/agent/types"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetCmdExe(svc *types.Services) types.CommandExec {
	return types.CommandExec(func(program string, args ...string) error {
		if runtime.GOOS == "windows" {
			if _, err := exec.LookPath(program); err != nil {
				args = append([]string{"/c", program}, args...)
				if p, err := ps.FindProcess(os.Getppid()); err != nil {
					if strings.Contains(strings.ToLower(p.Executable()), "powershell") {
						program = "powershell"
					} else {
						program = "cmd"
					}
				} else {
					program = "cmd"
				}
			}
		}
		cmd := exec.Command(program, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			svc.LogSvc.GetErr().Println(err)
		}
		return err
	})
}
