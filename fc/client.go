package fc

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	fc_open20210406 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// Client is the wrapper of fc client
type Client struct {
	client *fc_open20210406.Client
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

	client, err := fc_open20210406.NewClient(config)
	return &Client{client: client}, err
}
