package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// Query the chaincode to get the value of key
func (setup *FabricSetup) Query(key string) (string, error) {
	response, err := setup.QueryRaw(key)
	return string(response), err
}

// Query, but return raw payload
func (setup *FabricSetup) QueryRaw(key string) ([]byte, error) {

	response, err := setup.client.Query(channel.Request{
		ChaincodeID: setup.ChainCodeID,
		Fcn:         "invoke",
		Args: [][]byte{
			[]byte("query"),
			[]byte(key),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %v", err)
	}

	return response.Payload, nil
}
