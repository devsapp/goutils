package fc

import (
	"fmt"

	fc3 "github.com/alibabacloud-go/fc-20230330/v3/client"
	fc2 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

type Function struct {
	FunctionName          *string
	Runtime               *string
	GpuConfig             *GPUConfig
	CustomContainerConfig *CustomContainerConfig
	Layers                []*string
	EnvironmentVariables  map[string]*string
}

type CustomContainerConfig struct {
	Entrypoint []*string `json:"entrypoint,omitempty"`
	Command    []*string `json:"command,omitempty"`
	Image      *string   `json:"image,omitempty"`
}

type GPUConfig struct {
	GpuMemorySize *int32  `json:"gpuMemorySize,omitempty"`
	GpuType       *string `json:"gpuType,omitempty"`
}

func (f *Function) FromV2(function interface{}) *Function {
	switch f2 := function.(type) {
	case *fc2.ListFunctionsResponseBodyFunctions:
		f.FunctionName = f2.FunctionName
		f.Runtime = f2.Runtime
		f.Layers = make([]*string, 0)

		f.CustomContainerConfig = nil
		if f2.CustomContainerConfig != nil {
			f.CustomContainerConfig = &CustomContainerConfig{
				Entrypoint: []*string{f2.CustomContainerConfig.Command},
				Command:    []*string{f2.CustomContainerConfig.Args},
				Image:      f2.CustomContainerConfig.Image,
			}
		}

		f.GpuConfig = nil
		if *f2.GpuMemorySize != 0 {
			f.GpuConfig = &GPUConfig{
				GpuMemorySize: f2.GpuMemorySize,
				GpuType:       f2.InstanceType,
			}
		}

		for _, l := range f2.Layers {
			f.Layers = append(f.Layers, l)
		}

		f.EnvironmentVariables = make(map[string]*string)
		for k, v := range f2.EnvironmentVariables {
			f.EnvironmentVariables[k] = v
		}
	case *fc2.GetFunctionResponseBody:
		f.FunctionName = f2.FunctionName
		f.Runtime = f2.Runtime

		f.CustomContainerConfig = nil
		if f2.CustomContainerConfig != nil {
			f.CustomContainerConfig = &CustomContainerConfig{
				Entrypoint: []*string{f2.CustomContainerConfig.Command},
				Command:    []*string{f2.CustomContainerConfig.Args},
				Image:      f2.CustomContainerConfig.Image,
			}
		}

		f.GpuConfig = nil
		if *f2.GpuMemorySize != 0 {
			f.GpuConfig = &GPUConfig{
				GpuMemorySize: f2.GpuMemorySize,
				GpuType:       f2.InstanceType,
			}
		}

		for _, l := range f2.Layers {
			f.Layers = append(f.Layers, l)
		}

		f.EnvironmentVariables = make(map[string]*string)
		for k, v := range f2.EnvironmentVariables {
			f.EnvironmentVariables[k] = v
		}

	default:
		panic(fmt.Sprintf("unsupported type %t %v", f2, f2))
	}

	return f
}

func (f *Function) FromV3(function *fc3.Function) *Function {
	f.FunctionName = tea.String(*function.FunctionName)
	f.Runtime = tea.String(*function.Runtime)

	f.CustomContainerConfig = nil
	if function.CustomContainerConfig != nil {
		f.CustomContainerConfig = &CustomContainerConfig{
			Entrypoint: make([]*string, 0),
			Command:    make([]*string, 0),
			Image:      tea.String(*function.CustomContainerConfig.Image),
		}

		for _, e := range function.CustomContainerConfig.Entrypoint {
			f.CustomContainerConfig.Entrypoint = append(f.CustomContainerConfig.Entrypoint, tea.String(*e))
		}

		for _, cmd := range function.CustomContainerConfig.Command {
			f.CustomContainerConfig.Command = append(f.CustomContainerConfig.Command, tea.String(*cmd))
		}
	}

	f.GpuConfig = nil
	if function.GpuConfig != nil {
		f.GpuConfig = &GPUConfig{
			GpuMemorySize: tea.Int32(*function.GpuConfig.GpuMemorySize),
			GpuType:       tea.String(*function.GpuConfig.GpuType),
		}
	}

	for _, l := range function.Layers {
		f.Layers = append(f.Layers, tea.String(*l.Arn))
	}

	f.EnvironmentVariables = make(map[string]*string)
	for k, v := range function.EnvironmentVariables {
		f.EnvironmentVariables[k] = v
	}

	return f
}
