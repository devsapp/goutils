package fc

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	fc3 "github.com/alibabacloud-go/fc-20230330/v3/client"
	fc2 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// Client is the wrapper of fc client
type Client struct {
	clientV2 *fc2.Client
	clientV3 *fc3.Client
}

// NewClient creates a new fc client
func NewClient(ak, sk, sts, region, accountID string) (*Client, error) {
	config := new(openapi.Config)
	config.SetAccessKeyId(ak)
	config.SetAccessKeySecret(sk)
	config.Endpoint = tea.String(fmt.Sprintf("%s.%s.fc.aliyuncs.com", accountID, region))

	config.SetRegionId(region)
	config.SetProtocol("http")
	if sts != "" {
		config.SetSecurityToken(sts)
	}

	if IsV2() {
		client, err := fc2.NewClient(config)
		if err != nil {
			return nil, err
		}

		return &Client{clientV2: client}, nil
	}

	if IsV3() {
		client, err := fc3.NewClient(config)
		if err != nil {
			return nil, err
		}

		return &Client{clientV3: client}, nil
	}

	return nil, fmt.Errorf("Only support function-compute v2 and v3")
}
