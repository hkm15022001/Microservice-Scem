package message

import (
	"os"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

var (
	zbClient zbc.Client
)

// ConnectZeebeEngine function
func ConnectZeebeEngine() error {
	gatewayAddress := os.Getenv("BROKER_ADDRESS")
	newZbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         gatewayAddress,
		UsePlaintextConnection: true,
	})

	if err != nil {
		return err
	}

	zbClient = newZbClient
	return nil
}
