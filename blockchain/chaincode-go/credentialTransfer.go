/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

func main() {
	credentialChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating credential chaincode: %v", err)
	}

	if err := credentialChaincode.Start(); err != nil {
		log.Panicf("Error starting credential chaincode: %v", err)
	}
}
