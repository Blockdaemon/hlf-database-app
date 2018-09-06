package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// DatabaseChaincode implementation of Chaincode
type DatabaseChaincode struct {
}

func Errorf(format string, args ...interface{}) pb.Response {
	return shim.Error(fmt.Sprintf(format, args))
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *DatabaseChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### DatabaseChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Put in the ledger the key/value hello/world
	err := stub.PutState("hello", []byte("world"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *DatabaseChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### DatabaseChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return Errorf("Unknown action '%s'", args[0])
}

// query
// Every readonly functions in the ledger will be here
func (t *DatabaseChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### DatabaseChaincode query ###########")

	// Check whether the number of arguments is sufficient
	if len(args) != 2 {
		return shim.Error("The number of arguments is wrong.")
	}

	// Like the Invoke function, we manage multiple type of query requests with the second argument.
	// Get the state of the value matching the key requested
	state, err := stub.GetState(args[1])
	if err != nil {
		return Errorf("Failed to get state of %s", args[1])
	}

	// Return this value in response
	return shim.Success(state)
}

// invoke
// Every functions that read and write in the ledger will be here
func (t *DatabaseChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### DatabaseChaincode invoke ###########")

	// Check if we have the right number of args
	if len(args) != 3 {
		return shim.Error("The number of arguments is wrong.")
	}

	// Write the new value in the ledger
	err := stub.PutState(args[1], []byte(args[2]))
	if err != nil {
		return Errorf("Failed to update state of %s", args[1])
	}

	// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
	err = stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return this value in response
	return shim.Success(nil)
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(DatabaseChaincode))
	if err != nil {
		fmt.Printf("Error starting database chaincode: %s", err)
	}
}
