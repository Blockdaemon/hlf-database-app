package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Blockdaemon/config"

	"github.com/Blockdaemon/hlf-webapp/web"
	"github.com/Blockdaemon/hlf-webapp/web/controllers"

	"github.com/Blockdaemon/hlf-database-app/blockchain"
)

func initializeChannelAndCC(fSetup *blockchain.FabricSetup, force bool) {

	// Any one of these can fail if it was partially completed on last run,
	// so ignore errors for now, until this code is smarter.

	// FIXME: test if channel is already there and we joined it
	err := fSetup.CreateChannel()
	if err != nil {
		fmt.Printf("Unable to create channel: %v\n", err)
		if !force {
			return
		}
		fmt.Printf("IGNORING create channel error\n")
	}

	err = fSetup.UpdateChannel()
	if err != nil {
		fmt.Printf("Unable to update channel peers: %v\n", err)
		if !force {
			return
		}
		fmt.Printf("IGNORING update channel error\n")
	}

	err = fSetup.JoinChannel()
	if err != nil {
		fmt.Printf("Unable to join channel: %v\n", err)
		if !force {
			return
		}
		fmt.Printf("IGNORING join channel error\n")
	}

	// FIXME: test if CC is already installed
	err = fSetup.InstallCC()
	if err != nil {
		fmt.Printf("Unable to install the chaincode: %v\n", err)
		if !force {
			return
		}
		fmt.Printf("IGNORING install CC eror\n")
	}

	// FIXME: test if CC is already instantiated
	err = fSetup.InstantiateCC()
	if err != nil {
		fmt.Printf("Unable to instantiate the chaincode: %v\n", err)
	}
}

func usage() {
	fmt.Printf("%s: <init> (does create, update, join, install, instantiate)\n", os.Args[0])
	fmt.Printf("%s: <create | update | join | install | instantiate>\n", os.Args[0])
	fmt.Printf("%s: get <key>\n", os.Args[0])
	fmt.Printf("%s: set <key> <value>\n", os.Args[0])
	fmt.Printf("%s: store <key> <infile>\n", os.Args[0])
	fmt.Printf("%s: fetch <key> <outfile>\n", os.Args[0])
	fmt.Printf("%s: webapp\n", os.Args[0])
}

func newSetup(config *config.Config) (*blockchain.FabricSetup, error) {

	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Network parameters
		OrdererID: config.GetString("ORDERER_ID"),
		PeerOrg:   config.GetString("PEER_ORG"),

		// Channel parameters
		ChannelID:         config.GetString("CHANNEL"),
		ChannelConfig:     config.GetString("ARTIFACTS") + "/" + config.GetString("CHANNEL") + ".channel.tx",
		AnchorPeersConfig: config.GetString("ARTIFACTS") + "/" + config.GetString("CHANNEL") + ".anchor-peers.tx",

		// Chaincode parameters
		ChainCodeID:      "hlf-database-app",
		ChaincodeGoPath:  os.Getenv("GOPATH"),
		ChaincodePath:    "github.com/Blockdaemon/hlf-database-app/chaincode/",
		ChaincodeVersion: "0",
		OrgAdmin:         "Admin",
		OrgName:          os.Getenv("PEER_ORGNAME"),
		ConfigFile:       "config.yaml",

		// User parameters
		UserName: "Admin",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		return nil, err
	}

	return &fSetup, nil
}

func doInitCommand(fSetup *blockchain.FabricSetup, cmd string, config *config.Config) error {
	var err error
	switch cmd {
	case "init":
		initializeChannelAndCC(fSetup, true)
		return nil
	case "create":
		err = fSetup.CreateChannel()
	case "update":
		err = fSetup.UpdateChannel()
	case "join":
		err = fSetup.JoinChannel()
	case "install":
		err = fSetup.InstallCC()
	case "instantiate":
		err = fSetup.InstantiateCC()
	case "webapp":
		err = fSetup.CreateChannelAndEventClients()
		if err != nil {
			fmt.Printf("Unable to create channel and event clients: %v\n", err)
			return err
		}
		// Web app setup
		app := &controllers.Application{
			Fabric:  fSetup,
			WebRoot: config.GetString("WEBROOT"),
			WebPort: config.GetInt("WEBPORT"),
		}
		// GO GO GO!
		web.Serve(app)
	default:
		usage()
	}
	return err
}

func doGetSetCommand(fSetup *blockchain.FabricSetup) {
	var getKey, setKey, setValue string
	var storeKey, fetchKey, filename string

	if os.Args[1] != "get" && len(os.Args) < 4 {
		usage()
		return
	}

	switch os.Args[1] {
	case "get":
		if len(os.Args) != 3 {
			usage()
			return
		}
		getKey = os.Args[2]
	case "set":
		setKey = os.Args[2]
		setValue = os.Args[3]
	case "store":
		storeKey = os.Args[2]
		filename = os.Args[3]
	case "fetch":
		fetchKey = os.Args[2]
		filename = os.Args[3]
	default:
		usage()
		return
	}

	err := fSetup.CreateChannelAndEventClients()
	if err != nil {
		fmt.Printf("Unable to create channel and event clients: %v\n", err)
		return
	}

	if getKey != "" {
		val, err := fSetup.Query(getKey)
		if err != nil {
			fmt.Printf("Query '%s' failed: %v\n", getKey, err)
		} else {
			fmt.Printf("'%s'='%s'\n", getKey, val)
		}
	} else if setKey != "" {
		txid, err := fSetup.InvokeString(setKey, setValue)
		if err != nil {
			fmt.Printf("Invoke '%s'='%s' failed: %v\n", setKey, setValue, err)
		} else {
			fmt.Printf("Transaction %s successful\n", txid)
		}
	} else if storeKey != "" {
		val, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Failed to read '%s': %v\n", filename, err)
		} else {
			txid, err := fSetup.Invoke(storeKey, val)
			if err != nil {
				fmt.Printf("InvokeRaw '%s'= contents of '%s' failed: %v\n", storeKey, filename, err)
			} else {
				fmt.Printf("Transaction %s successful\n", txid)
			}
		}
	} else if fetchKey != "" {
		val, err := fSetup.QueryRaw(fetchKey)
		if err != nil {
			fmt.Printf("QueryRaw '%s' failed: %v\n", fetchKey, err)
		} else {
			err := ioutil.WriteFile(filename, val, os.FileMode(int(0644)))
			if err != nil {
				fmt.Printf("Failed to write '%s': %v\n", filename, err)
			}
		}
	} else {
		usage()
	}
}

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	bdsrc := os.Getenv("GOPATH") + "/src/github.com/Blockdaemon"
	config := new(config.Config)
	config.DescribeOptionalString("ORDERER_ID", "the orderer to use", "orderer0.hlf.blockdaemon.io")
	config.DescribeOptionalString("PEER_ORG", "Peer organization", "PeerOrg")
	config.DescribeOptionalString("PEER_ORGNAME", "The name of the org the admin is in", "PeerOrgName")
	config.DescribeOptionalString("CHANNEL", "The channel to use", "blockdaemon")
	config.DescribeOptionalString("ARTIFACTS", "The artifact directory",
		bdsrc+"/hlf-service-network/artifacts")
	config.DescribeOptionalString("WEBROOT", "The hlf-webapp directory",
		bdsrc+"/hlf-webapp")
	config.DescribeOptionalInt("WEBPORT", "The listen port for hlf-webapp", 3001)
	config.Parse()

	fSetup, err := newSetup(config)
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
		return
	}

	// Close SDK
	defer fSetup.CloseSDK()

	// Simple init command
	if len(os.Args) == 2 && os.Args[1] != "get" {
		err = doInitCommand(fSetup, os.Args[1], config)
		if err != nil {
			fmt.Printf("%s failed: %v\n", os.Args[1], err)
		}
	} else {
		doGetSetCommand(fSetup)
	}
}
