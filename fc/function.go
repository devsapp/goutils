package fc

import (
	"fmt"
)

// ListFunctions list all functions in a service
func (c *Client) ListFunctions(serviceOrPrefix string) ([]*Function, error) {
	if c.clientV2 != nil {
		return c.listFunctionV2(serviceOrPrefix)
	}
	if c.clientV3 != nil {
		return c.listFunctionV3(serviceOrPrefix)
	}

	return nil, fmt.Errorf("client is not initialized")
}

// GetFunction get function info
func (c *Client) GetFunction(serviceNameOrEmpty, functionName string) (*Function, error) {
	if c.clientV2 != nil {
		return c.GetFunctionV2(serviceNameOrEmpty, functionName)
	}
	if c.clientV3 != nil {
		return c.GetFunctionV3(functionName)
	}

	return nil, fmt.Errorf("client is not initialized")
}
