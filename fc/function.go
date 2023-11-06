package fc

import (
	"github.com/alibabacloud-go/tea/tea"

	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
)

// ListFunctions list all functions in a service
func (c *Client) ListFunctions(serviceName string) ([]*fc_open20210406.ListFunctionsResponseBodyFunctions, error) {
	funcs := make([]*fc_open20210406.ListFunctionsResponseBodyFunctions, 0)
	var token *string

	for {
		listFunctionsHeaders := &fc_open20210406.ListFunctionsHeaders{}
		listFunctionsRequest := &fc_open20210406.ListFunctionsRequest{NextToken: token}
		runtime := &util.RuntimeOptions{}

		resp, err := c.client.ListFunctionsWithOptions(tea.String(serviceName), listFunctionsRequest, listFunctionsHeaders, runtime)
		if err != nil {
			return nil, err
		}

		for _, function := range resp.Body.Functions {
			funcs = append(funcs, function)
		}

		if resp.Body.NextToken == nil {
			break
		}

		token = resp.Body.NextToken
	}

	return funcs, nil
}

// GetFunction get function info
func (c *Client) GetFunction(serviceName, functionName string) (*fc_open20210406.GetFunctionResponseBody, error) {
	getFunctionHeaders := &fc_open20210406.GetFunctionHeaders{}
	getFunctionRequest := &fc_open20210406.GetFunctionRequest{}
	runtime := &util.RuntimeOptions{}

	resp, err := c.client.GetFunctionWithOptions(tea.String(serviceName), tea.String(functionName), getFunctionRequest, getFunctionHeaders, runtime)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
