package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID         = "Org1MSP"
	cryptoPath    = "../test-network/organizations/peerOrganizations/org1.example.com"
	certPath      = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	keyPath       = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath   = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint  = "localhost:7051"
	gatewayPeer   = "peer0.org1.example.com"
	channelName   = "mychannel"
	chaincodeName = "diploma"
)

// Credential struct mirrors chaincode
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

// FabricService wraps gateway connection
type FabricService struct {
	gw       *client.Gateway
	network  *client.Network
	contract *client.Contract
}

// ConnectGateway initializes Fabric gateway connection
func ConnectGateway() (*FabricService, error) {
	ccp := &FabricService{}

	// gRPC connection
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certificatePEM); !ok {
		panic("failed to add peer TLS certificate to cert pool")
	}

	tlsCreds := credentials.NewClientTLSFromCert(certPool, gatewayPeer)
	grpcConn, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(tlsCreds))
	if err != nil {
		return nil, err
	}

	id, err := newIdentity()
	if err != nil {
		return nil, err
	}
	sign, err := newSign()
	if err != nil {
		return nil, err
	}

	ccp.gw, err = client.Connect(id,
		client.WithSign(sign),
		client.WithClientConnection(grpcConn),
		client.WithHash(hash.SHA256),
	)
	if err != nil {
		return nil, err
	}

	ccp.network = ccp.gw.GetNetwork(channelName)
	ccp.contract = ccp.network.GetContract(chaincodeName)

	return ccp, nil
}

func newIdentity() (*identity.X509Identity, error) {
	certPEM, err := readFirstFile(certPath)
	if err != nil {
		return nil, err
	}
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		return nil, err
	}
	return identity.NewX509Identity(mspID, cert)
}

func newSign() (identity.Sign, error) {
	keyPEM, err := readFirstFile(keyPath)
	if err != nil {
		return nil, err
	}
	privKey, err := identity.PrivateKeyFromPEM(keyPEM)
	if err != nil {
		return nil, err
	}
	return identity.NewPrivateKeySign(privKey)
}

func readFirstFile(dir string) ([]byte, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files in %s", dir)
	}
	return os.ReadFile(path.Join(dir, files[0].Name()))
}

// ReadCredential queries the chaincode for a credential by ID
func (f *FabricService) ReadCredential(id string) (*Credential, error) {
	result, err := f.contract.EvaluateTransaction("ReadCredential", id)
	if err != nil {
		return nil, err
	}
	var cred Credential
	if err := json.Unmarshal(result, &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}

func main() {
	fs, err := ConnectGateway()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/credential/:id", func(c *gin.Context) {
		id := c.Param("id")
		cred, err := fs.ReadCredential(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cred)
	})

	router.Run(":8080")
}
