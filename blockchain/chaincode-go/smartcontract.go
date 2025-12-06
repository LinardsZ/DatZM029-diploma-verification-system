/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Credential describes basic details of what makes up a simple credential
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Credential struct {
	ID                string          `json:"id"`                // Transaction identifier
	DiplomaHash       string          `json:"diplomaHash"`       // SHA-256 hash of diploma file
	GraduatePublicKey string          `json:"graduatePublicKey"` // Graduate's public key
	IssuerSignature   string          `json:"issuerSignature"`   // Issuer's signature of diplomaHash
	IssuerPublicKey   string          `json:"issuerPublicKey"`   // Issuer's public key for verification
	DiplomaMetadata   DiplomaMetadata `json:"diplomaMetadata"`   // Non-sensitive metadata
	Status            string          `json:"status"`            // "Valid" or "Revoked"
	CredentialType    string          `json:"credentialType"`    // Type of credential
}

type DiplomaMetadata struct {
	UniversityName string `json:"universityName"`
	DegreeName     string `json:"degreeName"`
	IssueDate      string `json:"issueDate"`
	ExpiryDate     string `json:"expiryDate"`
}

// InitLedger adds a base set of credentials to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	credentials := []Credential{
		{
			ID:                "credential1",
			DiplomaHash:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			GraduatePublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234...",
			IssuerSignature:   "3045022100abcd...",
			IssuerPublicKey:   "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5678...",
			DiplomaMetadata: DiplomaMetadata{
				UniversityName: "MIT",
				DegreeName:     "Bachelor of Science in Computer Science",
				IssueDate:      "2024-06-15",
				ExpiryDate:     "",
			},
			Status:         "Valid",
			CredentialType: "Diploma",
		},
		{
			ID:                "credential2",
			DiplomaHash:       "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592",
			GraduatePublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2345...",
			IssuerSignature:   "3046022100bcde...",
			IssuerPublicKey:   "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6789...",
			DiplomaMetadata: DiplomaMetadata{
				UniversityName: "Stanford University",
				DegreeName:     "Master of Business Administration",
				IssueDate:      "2024-05-20",
				ExpiryDate:     "",
			},
			Status:         "Valid",
			CredentialType: "Diploma",
		},
		{
			ID:                "credential3",
			DiplomaHash:       "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			GraduatePublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3456...",
			IssuerSignature:   "3045022100cdef...",
			IssuerPublicKey:   "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7890...",
			DiplomaMetadata: DiplomaMetadata{
				UniversityName: "Seoul National University",
				DegreeName:     "Doctor of Philosophy in Physics",
				IssueDate:      "2024-08-10",
				ExpiryDate:     "",
			},
			Status:         "Valid",
			CredentialType: "Diploma",
		},
	}

	for _, credential := range credentials {
		credentialJSON, err := json.Marshal(credential)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(credential.ID, credentialJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateCredential issues a new credential to the world state with given details.
func (s *SmartContract) CreateCredential(ctx contractapi.TransactionContextInterface, credentialJSON string) error {
	var credential Credential
	err := json.Unmarshal([]byte(credentialJSON), &credential)
	if err != nil {
		return fmt.Errorf("failed to unmarshal credential: %v", err)
	}

	exists, err := s.CredentialExists(ctx, credential.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the credential %s already exists", credential.ID)
	}

	credentialBytes, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(credential.ID, credentialBytes)
}

// ReadCredential returns the credential stored in the world state with given id.
func (s *SmartContract) ReadCredential(ctx contractapi.TransactionContextInterface, id string) (*Credential, error) {
	credentialJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if credentialJSON == nil {
		return nil, fmt.Errorf("the credential %s does not exist", id)
	}

	var credential Credential
	err = json.Unmarshal(credentialJSON, &credential)
	if err != nil {
		return nil, err
	}

	return &credential, nil
}

// UpdateCredential updates an existing credential in the world state with provided parameters.
func (s *SmartContract) UpdateCredential(ctx contractapi.TransactionContextInterface, credentialJSON string) error {
	var credential Credential
	err := json.Unmarshal([]byte(credentialJSON), &credential)
	if err != nil {
		return fmt.Errorf("failed to unmarshal credential: %v", err)
	}

	exists, err := s.CredentialExists(ctx, credential.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the credential %s does not exist", credential.ID)
	}

	credentialBytes, err := json.Marshal(credential)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(credential.ID, credentialBytes)
}

// DeleteCredential deletes an given credential from the world state.
func (s *SmartContract) DeleteCredential(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.CredentialExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the credential %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// CredentialExists returns true when credential with given ID exists in world state
func (s *SmartContract) CredentialExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	credentialJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return credentialJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllCredentials(ctx contractapi.TransactionContextInterface) ([]*Credential, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var credentials []*Credential
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var credential Credential
		err = json.Unmarshal(queryResponse.Value, &credential)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, &credential)
	}

	return credentials, nil
}
