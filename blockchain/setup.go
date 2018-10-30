package blockchain

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile       string
	Domain           string
	OrgID            string
	OrdererID        string
	ChannelID        string
	ChainCodeID      string
	initialized      bool
	ChannelConfig    string
	ChaincodeGoPath  string
	ChaincodePath    string
	ChaincodeVersion string
	OrgAdmin         string
	OrgName          string
	UserName         string

	sdk           *fabsdk.FabricSDK
	resClient     *resmgmt.Client
	adminIdentity *msp.SigningIdentity
	client        *channel.Client
	event         *event.Client
}

// Initialize reads the configuration file and sets up the client, chain and event hub
func (setup *FabricSetup) Initialize() error {

	// Add parameters for the initialization
	if setup.initialized {
		return errors.New("sdk already initialized")
	}

	// Initialize the SDK with the configuration file
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create SDK")
	}
	setup.sdk = sdk
	//fmt.Println("SDK created")

	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to load Admin identity")
	}
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}
	setup.resClient = resMgmtClient
	//fmt.Println("Resource management client created")

	// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to create MSP client")
	}
	adminIdentity, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		return errors.WithMessage(err, "failed to get admin signing identity")
	}
	setup.adminIdentity = &adminIdentity

	//fmt.Println("Initialization Successful")
	setup.initialized = true
	return nil
}

func (setup *FabricSetup) CreateAndJoinChannel() error {
	req := resmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfigPath: setup.ChannelConfig, SigningIdentities: []msp.SigningIdentity{*setup.adminIdentity}}
	txID, err := setup.resClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "failed to save channel")
	}
	fmt.Println("Channel created")

	// Make admin user join the previously created channel
	if err = setup.resClient.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
		return errors.WithMessage(err, "failed to make admin join channel")
	}
	fmt.Println("Channel joined")
	return nil
}

func (setup *FabricSetup) InstallCC() error {

	// Create the chaincode package that will be sent to the peers
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("ccPkg created")

	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: setup.ChaincodeVersion, Package: ccPkg}
	_, err = setup.resClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "failed to install chaincode")
	}
	fmt.Println("Chaincode installed")
	return nil
}

func (setup *FabricSetup) InstantiateCC() error {
	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{setup.OrgName + ".hlf." + setup.Domain})
	req := resmgmt.InstantiateCCRequest{
		Name:    setup.ChainCodeID,
		Path:    setup.ChaincodeGoPath,
		Version: "0",
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy,
		//CollConfig:	collConfig
	}

	resp, err := setup.resClient.InstantiateCC(setup.ChannelID, req)
	if err != nil || resp.TransactionID == "" {
		// Seriously, hyperledger?
		if strings.Contains(err.Error(), "chaincode exists "+setup.ChainCodeID) {
			fmt.Println("Chaincode already instantiated")
			return nil
		}
		if strings.Contains(err.Error(), "chaincode with name '"+setup.ChainCodeID+"' already exists") {
			fmt.Println("Chaincode already instantiated")
			return nil
		}
		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	fmt.Println("Chaincode instantiated")
	return nil
}

func (setup *FabricSetup) CreateChannelAndEventClients() (err error) { // LOL https://github.com/golang/go/issues/6842
	// Channel client is used to query and execute transactions
	clientContext := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	setup.client, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	//fmt.Println("Channel client created")

	// Creation of the client which will enables access to our channel events
	setup.event, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	//fmt.Println("Event client created")

	//fmt.Println("Create channel and event clients successful")
	return nil
}

func (setup *FabricSetup) CloseSDK() {
	setup.sdk.Close()
}
