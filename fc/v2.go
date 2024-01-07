package fc

import (
	fc2 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func (c *Client) listFunctionV2(serviceName string) ([]*Function, error) {
	funcs := make([]*Function, 0)
	var token *string

	for {
		listFunctionsHeaders := &fc2.ListFunctionsHeaders{}
		listFunctionsRequest := &fc2.ListFunctionsRequest{NextToken: token}
		runtime := &util.RuntimeOptions{}

		resp, err := c.clientV2.ListFunctionsWithOptions(tea.String(serviceName), listFunctionsRequest, listFunctionsHeaders, runtime)
		if err != nil {
			return nil, err
		}

		for _, function := range resp.Body.Functions {
			funcs = append(funcs, new(Function).FromV2(function))
		}

		if resp.Body.NextToken == nil {
			break
		}

		token = resp.Body.NextToken
	}

	return funcs, nil
}

func (c *Client) GetFunctionV2Raw(serviceName, functionName string) (*fc2.GetFunctionResponseBody, error) {
	getFunctionHeaders := &fc2.GetFunctionHeaders{}
	getFunctionRequest := &fc2.GetFunctionRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := c.clientV2.GetFunctionWithOptions(tea.String(serviceName), tea.String(functionName), getFunctionRequest, getFunctionHeaders, runtime)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *Client) GetFunctionV2(serviceName, functionName string) (*Function, error) {
	f, err := c.GetFunctionV2Raw(serviceName, functionName)
	if err != nil {
		return nil, err
	}

	return new(Function).FromV2(f), nil
}
