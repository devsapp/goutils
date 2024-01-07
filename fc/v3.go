package fc

import (
	fc3 "github.com/alibabacloud-go/fc-20230330/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

func (c *Client) listFunctionV3(prefix string) ([]*Function, error) {
	funcs := make([]*Function, 0)
	var token *string

	for {
		listFunctionsRequest := &fc3.ListFunctionsRequest{NextToken: token, Prefix: tea.String(prefix)}

		resp, err := c.clientV3.ListFunctions(listFunctionsRequest)
		if err != nil {
			return nil, err
		}

		for _, function := range resp.Body.Functions {
			funcs = append(funcs, new(Function).FromV3(function))
		}

		if resp.Body.NextToken == nil {
			break
		}

		token = resp.Body.NextToken
	}

	return funcs, nil
}

func (c *Client) GetFunctionV3Raw(functionName string) (*fc3.Function, error) {
	getFunctionRequest := &fc3.GetFunctionRequest{}

	resp, err := c.clientV3.GetFunction(tea.String(functionName), getFunctionRequest)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) GetFunctionV3(functionName string) (*Function, error) {
	f, err := c.GetFunctionV3Raw(functionName)
	if err != nil {
		return nil, err
	}

	return new(Function).FromV3(f), nil
}
