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
	config.DescribeOptionalString("ARTIFACTS", "The artifact directory", "artifacts")
	config.Parse()

	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		Domain: config.GetString("DOMAIN"),
		OrdererID: "orderer.hf." + config.GetString("DOMAIN"),

		// Channel parameters
		ChannelID:     config.GetString("CHANNEL"),
		ChannelConfig: config.GetString("ARTIFACTS") + "/" + config.GetString("CHANNEL") + ".channel.tx",

		// Chaincode parameters
		ChainCodeID:     "hlf-database-app",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/Blockdaemon/hlf-database-app/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "org1",
		ConfigFile:      "config.yaml",

		// User parameters
		UserName: os.Getenv("USER"),
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
	}
	// Close SDK
	defer fSetup.CloseSDK()

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
	}

	// Launch the web application listening
	/*
	app := &controllers.Application{
		Fabric: &fSetup,
	}
	web.Serve(app)
	*/
}
