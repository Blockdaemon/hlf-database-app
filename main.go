package main

import (
	"fmt"
	"os"

	"github.com/Blockdaemon/config"

	"github.com/Blockdaemon/hlf-database-app/blockchain"
)

func main() {
	config := new(config.Config)
	config.DescribeOptionalString("DOMAIN", "The domain to use in CAs", "blockdaemon.io")
	config.DescribeOptionalString("CHANNEL", "The channel to use", "blockdaemon")
	config.DescribeOptionalString("ARTIFACTS", "The artifact directory", os.Getenv("GOPATH")+"/src/github.com/Blockdaemon/hlf-service-network/artifacts")
	config.Parse()

	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		Domain:    config.GetString("DOMAIN"),
		OrdererID: "orderer.hf." + config.GetString("DOMAIN"),

		// Channel parameters
		ChannelID:     config.GetString("CHANNEL"),
		ChannelConfig: config.GetString("ARTIFACTS") + "/" + config.GetString("CHANNEL") + ".channel.tx",

		// Chaincode parameters
		ChainCodeID:      "hlf-database-app",
		ChaincodeGoPath:  os.Getenv("GOPATH"),
		ChaincodePath:    "github.com/Blockdaemon/hlf-database-app/chaincode/",
		ChaincodeVersion: "0",
		OrgAdmin:         "Admin",
		OrgName:          "org1",
		ConfigFile:       "config.yaml",

		// User parameters
		UserName: "Admin",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}
	// Close SDK
	defer fSetup.CloseSDK()

	err = fSetup.CreateAndJoinChannel()
	if err != nil {
		fmt.Printf("Unable to create and join channel: %v\n", err)
		//return
	}

	// Install the chaincode
	err = fSetup.InstallCC()
	if err != nil {
		fmt.Printf("Unable to install the chaincode: %v\n", err)
		//return
	}

	// Instantiate the chaincode
	err = fSetup.InstantiateCC()
	if err != nil {
		fmt.Printf("Unable to instantiate the chaincode: %v\n", err)
		//return
	}

	err = fSetup.CreateChannelAndEventClients()
	if err != nil {
		fmt.Printf("Unable to create channel and event clients: %v\n", err)
		return
	}

	// Launch the web application listening
	/*
		app := &controllers.Application{
			Fabric: &fSetup,
		}
		web.Serve(app)
	*/
}
