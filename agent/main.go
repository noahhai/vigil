package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/noahhai/vigil/agent/helpers"
	"github.com/noahhai/vigil/agent/initial"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var cmdStr, statusStr, durationStr string
	var report bool
	flag.BoolVarP(&report, "r", "r", false, "report only")
	flag.StringVarP(&cmdStr, "command", "c", "", "command text for reporting")
	flag.StringVarP(&statusStr, "status", "s", "", "result text for reporting")
	flag.StringVarP(&durationStr, "duration", "d", "", "duration text for reporting")
	flag.Parse()


	svc := initial.InitServices()
	uc := initial.InitUseCases(svc)
	args := os.Args[1:]
	c, err  := helpers.LoadConfig()
	if err != nil {
		svc.GetErr().Println(err)
		os.Exit(1)
	}


	if report {
		if cmdStr == "" {
			svc.GetErr().Println("command text '-c' is required for report mode")
			os.Exit(1)
		}
		if durationStr == "" {
			durationStr = "<unspecified>"
		}
		if statusStr == "" {
			statusStr = "<unspecified>"
		}
		if err := postData(c, cmdStr, durationStr, statusStr); err != nil {
			svc.GetErr().Printf("Failed to forward task status to notifier: %v\n", err)
			os.Exit(1)
		} else {
			svc.GetOut().Println("Forwarded task status to notifier")
			os.Exit(0)
		}

	}

	if len(args) < 1 {
		svc.GetErr().Println("no command specified")
		os.Exit(1)
	}
	var subArgs []string
	if len(args) > 1 {
		subArgs = args[1:]
	}
	start := time.Now()
	err = uc.CommandExec(args[0], subArgs...)
	elapsed := time.Since(start)
	elapsedString := elapsed.String()
	svc.GetOut().Printf("Elapsed time: %s\n", elapsedString)
	taskName := strings.Join(append([]string{args[0]}, subArgs...), " ")
	status := "failure. " + err.Error()
	if err == nil {
		status = "success"
	}
	if err := postData(c, taskName, elapsedString, status); err != nil {
		svc.GetErr().Printf("Failed to forward task status to notifier: %v\n", err)
		os.Exit(1)
	} else {
		svc.GetOut().Println("Forwarded task status to notifier")
	}
}

type taskFinish struct {
	Name string `json:"name"`
	Duration string `json:"duration"`
	Status string `json:"status"`
	Username string `json:"username"`
	Token string `json:"token"`
}

func postData(c helpers.Config, taskName string, duration string, status string) error {
	requestBody, err := json.Marshal(taskFinish{
		Name: taskName,
		Duration: duration,
		Status: status,
		Username: c.Username,
		Token: c.Token,
	})
	if err != nil {
		return err
	}
	workPath := fmt.Sprintf("%s/%s", helpers.GetAppRootURL(), "work")
	resp, err := http.Post(workPath, "application/json", bytes.NewBuffer(requestBody))
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err == nil {
		err = errFromStatus(resp.StatusCode)
	}
	return err
}

func errFromStatus (status int) error {
	if status < 200 || status > 300 {
		return fmt.Errorf("error status code: %d", status)
	}
	return nil
}