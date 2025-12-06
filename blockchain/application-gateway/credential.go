package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "../test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

var now = time.Now()
var credentialId = fmt.Sprintf("credential%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gw.Close()

	// Override default values for chaincode and channel name as they may differ in testing contexts.
	chaincodeName := "diploma"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	initLedger(contract)
	getAllCredentials(contract)
	CreateCredential(contract)
	readCredentialByID(contract)
	// transferAssetAsync(contract)
	// exampleErrorHandling(contract)
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

// This type of transaction would typically only be run once by an application the first time it was started after its
// initial deployment. A new version of the chaincode deployed later would likely not need to run an "init" function.
func initLedger(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger \n")

	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction to query ledger state.
func getAllCredentials(contract *client.Contract) {
	fmt.Println("\n--> Evaluate Transaction: GetAllCredentials, function returns all the current assets on the ledger")

	evaluateResult, err := contract.EvaluateTransaction("GetAllCredentials")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

type Credential struct {
	ID                string          `json:"id"`
	DiplomaHash       string          `json:"diplomaHash"`
	GraduatePublicKey string          `json:"graduatePublicKey"`
	IssuerSignature   string          `json:"issuerSignature"`
	IssuerPublicKey   string          `json:"issuerPublicKey"`
	DiplomaMetadata   DiplomaMetadata `json:"diplomaMetadata"`
	Status            string          `json:"status"`
	CredentialType    string          `json:"credentialType"`
}

type DiplomaMetadata struct {
	UniversityName string `json:"universityName"`
	DegreeName     string `json:"degreeName"`
	IssueDate      string `json:"issueDate"`
	ExpiryDate     string `json:"expiryDate"`
}

// Submit a transaction synchronously, blocking until it has been committed to the ledger.
func CreateCredential(contract *client.Contract) {
	fmt.Printf("\n--> Submit Transaction: CreateCredential, creates new credential with ID, Color, Size, Owner and AppraisedValue arguments \n")

	credential := Credential{
		ID:                credentialId,
		DiplomaHash:       "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
		GraduatePublicKey: "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4567...",
		IssuerSignature:   "3045022100defg...",
		IssuerPublicKey:   "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8901...",
		DiplomaMetadata: DiplomaMetadata{
			UniversityName: "Harvard University",
			DegreeName:     "Bachelor of Arts in Economics",
			IssueDate:      "2024-12-06",
			ExpiryDate:     "",
		},
		Status:         "Valid",
		CredentialType: "Diploma",
	}

	credentialJSON, err := json.Marshal(credential)
	if err != nil {
		panic(fmt.Errorf("failed to marshal credential: %w", err))
	}

	_, err = contract.SubmitTransaction("CreateCredential", string(credentialJSON))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

// Evaluate a transaction by assetID to query ledger state.
func readCredentialByID(contract *client.Contract) {
	fmt.Printf("\n--> Evaluate Transaction: ReadCredential, function returns asset attributes\n")

	evaluateResult, err := contract.EvaluateTransaction("ReadCredential", credentialId)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := formatJSON(evaluateResult)

	fmt.Printf("*** Result:%s\n", result)
}

// Submit transaction, passing in the wrong number of arguments ,expected to throw an error containing details of any error responses from the smart contract.
// func exampleErrorHandling(contract *client.Contract) {
// 	fmt.Println("\n--> Submit Transaction: UpdateAsset asset70, asset70 does not exist and should return an error")

// 	_, err := contract.SubmitTransaction("UpdateAsset", "asset70", "blue", "5", "Tomoko", "300")
// 	if err == nil {
// 		panic("******** FAILED to return an error")
// 	}

// 	fmt.Println("*** Successfully caught the error:")

// 	var commitStatusErr *client.CommitStatusError
// 	var transactionErr *client.TransactionError

// 	if errors.As(err, &commitStatusErr) {
// 		if errors.Is(err, context.DeadlineExceeded) {
// 			fmt.Printf("Timeout waiting for transaction %s commit status: %s\n", commitStatusErr.TransactionID, commitStatusErr)
// 		} else {
// 			fmt.Printf("Error obtaining commit status for transaction %s with gRPC status %v: %s\n", commitStatusErr.TransactionID, status.Code(commitStatusErr), commitStatusErr)
// 		}
// 	} else if errors.As(err, &transactionErr) {
// 		// The error could be an EndorseError, SubmitError or CommitError.
// 		fmt.Println(err)
// 		fmt.Printf("TransactionID: %s\n", transactionErr.TransactionID)
// 	} else {
// 		panic(fmt.Errorf("unexpected error type %T: %w", err, err))
// 	}
// }

// Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
