package project

import (
	"fmt"
	"strings"

	gr "github.com/awesome-fc/golang-runtime"
	"github.com/devsapp/goutils/fc"
	"github.com/sirupsen/logrus"
)

const PrefixDelimiter = "__"

// T is a struct for project functions
type T struct {
	StableDiffusion *fc.Function `json:"stable_diffusion"`
	Filemgr         *fc.Function `json:"filemgr"`
	Lora            *fc.Function `json:"lora"`
}

// Get project functions
func Get(ctx *gr.FCContext) T {
	t := T{}

	client, err := fc.NewClient(ctx.Credentials.AccessKeyID, ctx.Credentials.AccessKeySecret, ctx.Credentials.SecurityToken, ctx.Region, ctx.AccountID)
	if err != nil {
		logrus.Errorf("fc client failed, due to %s", err)
	} else {
		serviceNameOrPrefix := ctx.Service.ServiceName
		if serviceNameOrPrefix == "" {
			parts := strings.Split(ctx.Function.Name, PrefixDelimiter)
			if len(parts) >= 1 {
				serviceNameOrPrefix = parts[0]
			}
		}

		funcs, err := client.ListFunctions(serviceNameOrPrefix)
		if err != nil {
			logrus.Errorf("list function failed, due to %s", err)
		} else {
			mayFilemgr := make([]*fc.Function, 0)
			mayStableDiffusion := make([]*fc.Function, 0)

			for _, function := range funcs {
				if withImage("kohya_ss", function) {
					// 使用 kohya_ss 镜像, 说明是 lora 函数
					t.Lora = function
				} else if withFunctionName("sd", serviceNameOrPrefix, *function.FunctionName) {
					// 函数名为 sd，说明是 sd 函数
					t.StableDiffusion = function
				} else if withImage("fc-stable-diffusion", function) && !withImage("kohya_ss", function) {
					// 镜像名为 fc-stable-diffusion 并且不是 kohya_ss，说明可能是 sd 函数
					mayStableDiffusion = append(mayStableDiffusion, function)
				} else if function.GpuConfig == nil && *function.Runtime == "custom" && len(function.Layers) > 0 {
					// 不是 gpu，是 custom runtime 并且有 layers，说明可能是文件管理器
					mayFilemgr = append(mayFilemgr, function)
					if withFunctionName("admin", serviceNameOrPrefix, *function.FunctionName) {
						// 如果函数是 admin，则一定是文件管理器
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

// withFunctionName is match, support fc v2 and fc v3
func withFunctionName(wantFunctionName, prefix, functionName string) bool {
	if fc.IsV3() {
		return wantFunctionName == fmt.Sprintf("%s%s%s", prefix, PrefixDelimiter, functionName)
	}

	return wantFunctionName == functionName
}

func withImage(wantImage string, f *fc.Function) bool {
	return f != nil && f.CustomContainerConfig != nil && f.CustomContainerConfig.Image != nil && strings.Contains(*f.CustomContainerConfig.Image, wantImage)
}
