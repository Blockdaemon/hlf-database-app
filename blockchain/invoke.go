package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"time"
)

// Invoke - set key's value to string
func (setup *FabricSetup) Invoke(key string, value string) (string, error) {
	return setup.InvokeRaw(key, []byte(value))
}

// Invoke - set key's value to []byte array
func (setup *FabricSetup) InvokeRaw(key string, value []byte) (string, error) {

	eventID := "eventInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in invoke")

	reg, notifier, err := setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}
	defer setup.event.Unregister(reg)

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(channel.Request{
		ChaincodeID: setup.ChainCodeID,
		Fcn:         "invoke",
		Args: [][]byte{
			[]byte("invoke"),
			[]byte(key),
			value,
		},
		TransientMap: transientDataMap,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	return string(response.TransactionID), nil
}
