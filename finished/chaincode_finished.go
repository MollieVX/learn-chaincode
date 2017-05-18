/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//structure for account
type Account struct {
	Name     string `json:"name"`
	Balance  int    `json:"balance"`
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("No arguments required")
	}

	var account1, account2, account3 Account
	account1.Name = "Vatsala"
	account1.Balance = 1000
	account2.Name = "Harish"
	account2.Balance = 1000
	account3.Name = "Narayan"
	account3.Balance = 1000

	b, err := json.Marshal(account1)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for account 1")
	}

	err = stub.PutState("Vatsala", b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(account2)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for account 2")
	}

	err = stub.PutState("Harish", b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(account3)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for account 3")
	}

	err = stub.PutState("Narayan", b)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "transfer"{
		return t.transfer(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

//Function to transfer funds from one account to another
func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var accFrom Account
	var accTo Account
	var amount int

	accFromAsBytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return nil, errors.New(jsonResp)
	}
	err = json.Unmarshal(accFromAsBytes, &accFrom)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of accFrom")
	}

	accToAsBytes, err := stub.GetState(args[1])
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + args[1] + "\"}"
		return nil, errors.New(jsonResp)
	}
	err = json.Unmarshal(accToAsBytes, &accTo)
	if err != nil {
		return nil, errors.New("Failed to marshal string to struct of accTo")
	}

	amount, err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Enter an integer value in the 'Amount'")
	}

	if accFrom.Balance < amount {
		return nil, errors.New("Insufficient Balance")
	}

	accFrom.Balance = accFrom.Balance - amount
	accTo.Balance = accTo.Balance + amount

	b, err := json.Marshal(accFrom)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for account 1")
	}

	err = stub.PutState(args[0], b)
	if err != nil {
		return nil, err
	}

	b, err = json.Marshal(accTo)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Errors while creating json string for account 1")
	}

	err = stub.PutState(args[1], b)
	if err != nil {
		return nil, err
	}

	return nil, nil

}
