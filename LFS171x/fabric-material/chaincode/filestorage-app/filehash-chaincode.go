/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */


package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

//// Define the Smart Contract structure
type SmartContract struct {
}


/* Define File storage  structure, with 2 properties.  
Structure tags are used by encoding/json library
*/

type Filestorage struct {
	Fileguid string `json:"fileguid"`
	Filehash string `json:"filehash"`
	Timestamp string `json:"timestamp"`
}

/*
 * The Init method *
 called when the Smart Contract "filehash-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "filehash-chaincode"
 The app also specifies the specific smart contract function to call with args
 */

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryHash" {
		return s.queryHash(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordHash" {
		return s.recordHash(APIstub, args)
	} else if function == "queryAllHash" {
		return s.queryAllHash(APIstub)
	} else if function == "changeFileHash" {
		return s.changeFileHash(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}


/*
 * The queryHash method *
Used to view the records of one particular hash
It takes one argument -- the key for the hash in question
 */

func (s *SmartContract) queryHash(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fileAsBytes, _ := APIstub.GetState(args[0])
	if fileAsBytes == nil {
		return shim.Error("Could not locate file")
	}
	return shim.Success(fileAsBytes)
}


/*
 * The initLedger method *
  initialize data to our network
 */

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	file := []File{
        }

	i := 0
	for i < len(file) {
		fmt.Println("i is ", i)
		fileAsBytes, _ := json.Marshal(file[i])
		APIstub.PutState(strconv.Itoa(i+1), fileAsBytes)
		fmt.Println("File Added", file[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordFile method * 
 user like Suman would use to record each of his  file hashes.
This method takes in three arguments (attributes to be saved in the ledger). 
 */

func (s *SmartContract) recordFile(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var file = File{ Filehash: args[1], Timestamp: args[2] }

	fileAsBytes, _ := json.Marshal(file)
	err := APIstub.PutState(args[0], fileAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record file hash: %s", args[0]))
	}

	return shim.Success(nil)
}


/*
 * The queryAllFile method *
allows for assessing all the records added to the ledger(all file hashes)
This method does not take any arguments. Returns JSON string containing results. 
 */

func (s *SmartContract) queryAllFile(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
	
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllFile:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


/*
 * The changeFileHash method *
The data in the world state can be updated with who has possession. 
This function takes in 2 arguments, file id and new file hash. 
 */

func (s *SmartContract) changeFileHash(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	fileAsBytes, _ := APIstub.GetState(args[0])
	if fileAsBytes == nil {
		return shim.Error("Could not locate file")
	}
	file := File{}

	json.Unmarshal(fileAsBytes, &file)

	file.Filehash = args[1]

	fileAsBytes, _ = json.Marshal(file)
	err := APIstub.PutState(args[0], fileAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change file hash: %s", args[0]))
	}

	return shim.Success(nil)
}


/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}