/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const IssuerKey = "ISSUER_"
const IssuerKeyRangeEnd = "ISSUER_\uffff"
const CredentialKey = "CREDENTIAL_"
const CredentialKeyRangeEnd = "CREDENTIAL_\uffff"

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
	IssuerID          string          `json:"issuerId"`
	IssuerSignature   string          `json:"issuerSignature"` // Issuer's signature of diplomaHash
	DiplomaMetadata   DiplomaMetadata `json:"diplomaMetadata"` // Non-sensitive metadata
	Status            string          `json:"status"`          // "Valid" or "Revoked"
	CredentialType    string          `json:"credentialType"`  // Type of credential
}

type DiplomaMetadata struct {
	UniversityName string `json:"universityName"`
	DegreeName     string `json:"degreeName"`
	IssueDate      string `json:"issueDate"`
	ExpiryDate     string `json:"expiryDate"`
}

type Issuer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`    // "Active" or "Revoked"
	PublicKey string `json:"publicKey"` // Issuer's public key for verification
}

// InitLedger adds a base set of issuers to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	err := s.addAuthorizedIssuers(ctx)
	if err != nil {
		return fmt.Errorf("failed to add authorized issuers: %v", err)
	}

	return nil
}

// CreateCredential issues a new credential to the world state with given details.
func (s *SmartContract) CreateCredential(ctx contractapi.TransactionContextInterface, credentialJSON string) (string, error) {
	var credential Credential

	err := json.Unmarshal([]byte(credentialJSON), &credential)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal credential: %v", err)
	}

	credential.ID = CredentialKey + credential.ID

	// Check if the issuer exists and is active
	issuerKey := "ISSUER_" + credential.IssuerID

	issuerBytes, err := ctx.GetStub().GetState(issuerKey)
	if err != nil {
		return "", fmt.Errorf("failed to read issuer from ledger: %v", err)
	}

	if issuerBytes == nil {
		return "", fmt.Errorf("issuer %s does not exist", credential.IssuerID)
	}

	var issuer Issuer
	if err := json.Unmarshal(issuerBytes, &issuer); err != nil {
		return "", fmt.Errorf("failed to unmarshal issuer data: %v", err)
	}

	if issuer.Status != "Active" {
		return "", fmt.Errorf("issuer %s is not active", credential.IssuerID)
	}

	exists, err := s.CredentialExists(ctx, credential.ID)
	if err != nil {
		return "", fmt.Errorf("failed to check credential existence: %v", err)
	}
	if exists {
		return "", fmt.Errorf("the credential %s already exists", credential.ID)
	}

	credentialBytes, err := json.Marshal(credential)
	if err != nil {
		return "", fmt.Errorf("failed to marshal credential data: %v", err)
	}

	return credential.ID, ctx.GetStub().PutState(credential.ID, credentialBytes)
}

// ReadCredential returns the credential stored in the world state with given id.
func (s *SmartContract) ReadCredential(ctx contractapi.TransactionContextInterface, id string) (*Credential, error) {
	key := CredentialKey + id
	credentialJSON, err := ctx.GetStub().GetState(key)
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

	credential.ID = CredentialKey + credential.ID

	// Check if the issuer exists and is active
	issuerKey := "ISSUER_" + credential.IssuerID

	issuerBytes, err := ctx.GetStub().GetState(issuerKey)
	if err != nil {
		return fmt.Errorf("failed to read issuer from ledger: %v", err)
	}

	if issuerBytes == nil {
		return fmt.Errorf("issuer %s does not exist", credential.IssuerID)
	}

	var issuer Issuer
	if err := json.Unmarshal(issuerBytes, &issuer); err != nil {
		return fmt.Errorf("failed to unmarshal issuer data: %v", err)
	}

	if issuer.Status != "Active" {
		return fmt.Errorf("issuer %s is not active", credential.IssuerID)
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
	resultsIterator, err := ctx.GetStub().GetStateByRange(CredentialKey, CredentialKeyRangeEnd)
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

func (s *SmartContract) RevokeCredential(ctx contractapi.TransactionContextInterface, id string) error {
	key := CredentialKey + id

	data, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}

	if data == nil {
		return fmt.Errorf("credential not found")
	}

	var credential Credential
	if err := json.Unmarshal(data, &credential); err != nil {
		return err
	}

	credential.Status = "Revoked"

	newData, _ := json.Marshal(credential)
	
	return ctx.GetStub().PutState(key, newData)
}

func (s *SmartContract) ReadIssuer(ctx contractapi.TransactionContextInterface, id string) (*Issuer, error) {
	key := IssuerKey + id

	state, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if state == nil {
		return nil, fmt.Errorf("the issuer %s does not exist", id)
	}

	var issuer Issuer
	err = json.Unmarshal(state, &issuer)
	if err != nil {
		return nil, err
	}

	return &issuer, nil
}

func (s *SmartContract) CreateIssuer(ctx contractapi.TransactionContextInterface, issuerJSON string) error {
	var issuer Issuer
	if err := json.Unmarshal([]byte(issuerJSON), &issuer); err != nil {
		return err
	}

	key := IssuerKey + issuer.ID
	exists, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}
	if exists != nil {
		return fmt.Errorf("issuer %s already exists", issuer.ID)
	}

	data, err := json.Marshal(issuer)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, data)
}

func (s *SmartContract) RevokeIssuer(ctx contractapi.TransactionContextInterface, id string) error {
	key := IssuerKey + id
	data, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("issuer not found")
	}

	var issuer Issuer
	if err := json.Unmarshal(data, &issuer); err != nil {
		return err
	}
	issuer.Status = "Revoked"

	newData, _ := json.Marshal(issuer)
	return ctx.GetStub().PutState(key, newData)
}

func (s *SmartContract) GetAllIssuers(ctx contractapi.TransactionContextInterface) ([]*Issuer, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange(IssuerKey, IssuerKeyRangeEnd)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var items []*Issuer
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var issuer Issuer
		err = json.Unmarshal(queryResponse.Value, &issuer)
		if err != nil {
			return nil, err
		}
		items = append(items, &issuer)
	}

	return items, nil
}

func (s *SmartContract) addAuthorizedIssuers(ctx contractapi.TransactionContextInterface) error {
	issuers := []Issuer{
		{
			ID:        "lu",
			Name:      "University of Latvia",
			Status:    "Active",
			PublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5678...",
		},
		{
			ID:        "rtu",
			Name:      "Riga Technical University",
			Status:    "Active",
			PublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5678...",
		},
	}

	for _, issuer := range issuers {
		issuerJSON, err := json.Marshal(issuer)
		if err != nil {
			return err
		}

		key := IssuerKey + issuer.ID

		err = ctx.GetStub().PutState(key, issuerJSON)
		if err != nil {
			return fmt.Errorf("failed to put issuer: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) AddMockCredentials(ctx contractapi.TransactionContextInterface) error {
	credentials := []Credential{
		{
			ID:                "credential1",
			DiplomaHash:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			GraduatePublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234...",
			IssuerID: "lu",
			IssuerSignature:   "3045022100abcd...",
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
			IssuerID: "lu",
			IssuerSignature:   "3046022100bcde...",
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
			IssuerID: "rtu",
			IssuerSignature:   "3045022100cdef...",
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

		key := CredentialKey + credential.ID
		err = ctx.GetStub().PutState(key, credentialJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}
