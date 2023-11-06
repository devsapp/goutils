package project

import (
	"fmt"
	"strings"

	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	gr "github.com/awesome-fc/golang-runtime"
	"github.com/devsapp/goutils/fc"
	"github.com/sirupsen/logrus"
)

// T is a struct for project functions
type T struct {
	StableDiffusion *fc_open20210406.ListFunctionsResponseBodyFunctions `json:"stable_diffusion"`
	Filemgr         *fc_open20210406.ListFunctionsResponseBodyFunctions `json:"filemgr"`
	Lora            *fc_open20210406.ListFunctionsResponseBodyFunctions `json:"lora"`
}

// Get project functions
func Get(ctx *gr.FCContext) T {
	t := T{}

	client, err := fc.NewClient(ctx.Credentials.AccessKeyID, ctx.Credentials.AccessKeySecret, ctx.Credentials.SecurityToken, ctx.Region, ctx.AccountID)
	if err != nil {
		logrus.Errorf("fc client failed, due to %s", err)
	} else {
		funcs, err := client.ListFunctions(ctx.Service.ServiceName)
		if err != nil {
			logrus.Errorf("list function failed, due to %s", err)
		} else {
			mayFilemgr := make([]*fc_open20210406.ListFunctionsResponseBodyFunctions, 0)
			mayStableDiffusion := make([]*fc_open20210406.ListFunctionsResponseBodyFunctions, 0)

			for _, function := range funcs {
				fmt.Println(*function.FunctionName, ctx.Function.Name)
				if function.CustomContainerConfig != nil && function.CustomContainerConfig.Image != nil && strings.Contains(*function.CustomContainerConfig.Image, "kohya_ss") {
					fmt.Println(*function.FunctionName, "lora")
					t.Lora = function
				} else if *function.FunctionName == "sd" {
					t.StableDiffusion = function
				} else if function.CustomContainerConfig != nil && function.CustomContainerConfig.Image != nil && strings.Contains(*function.CustomContainerConfig.Image, "fc-stable-diffusion") && !strings.Contains(*function.CustomContainerConfig.Image, "kohya_ss") {
					mayStableDiffusion = append(mayStableDiffusion, function)
				} else if function.GpuMemorySize == nil && *function.Runtime == "custom" && len(function.Layers) > 0 {
					mayFilemgr = append(mayFilemgr, function)
					if *function.FunctionName == "admin" {
						t.Filemgr = function
					}
				}
			}

			if t.Filemgr == nil && len(mayFilemgr) > 0 {
				t.Filemgr = mayFilemgr[0]
			}
			if t.StableDiffusion == nil && len(mayStableDiffusion) > 0 {
				t.StableDiffusion = mayStableDiffusion[0]
			}
		}
	}

	return t
}
