package main

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"golang.org/x/crypto/openpgp"
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
	IssuerID          string          `json:"issuerId"`
	IssuerSignature   string          `json:"issuerSignature"`
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

// Issuer represents an authorized organization
type Issuer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIKey    string `json:"apiKey"`
	Signature string `json:"signature"`
}

// CreateCredentialRequest for POST /credential
type CreateCredentialRequest struct {
	DiplomaHash       string          `json:"diplomaHash" binding:"required"`
	GraduatePublicKey string          `json:"graduatePublicKey" binding:"required"`
	IssuerID          string          `json:"issuerId" binding:"required"`
	IssuerSignature   string          `json:"issuerSignature" binding:"required"`
	DiplomaMetadata   DiplomaMetadata `json:"diplomaMetadata" binding:"required"`
	CredentialType    string          `json:"credentialType" binding:"required"`
}

// VerifyHashRequest for POST /verify/hash
type VerifyHashRequest struct {
	DiplomaHash string `json:"diplomaHash" binding:"required"`
}

// VerifySignatureRequest for POST /verify/signature
type VerifySignatureRequest struct {
	CredentialID      string `json:"credentialId" binding:"required"`
	GraduateSignature string `json:"graduateSignature" binding:"required"`
}

var issuers = []Issuer{}

type User struct {
	IssuerID     string `json:"issuerId"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
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

// GetAllCredentials queries all credentials from the chaincode
func (f *FabricService) GetAllCredentials() ([]*Credential, error) {
	result, err := f.contract.EvaluateTransaction("GetAllCredentials")
	if err != nil {
		return nil, err
	}
	var creds []*Credential
	if err := json.Unmarshal(result, &creds); err != nil {
		return nil, err
	}
	return creds, nil
}

// CreateCredential submits a transaction to create a new credential
func (f *FabricService) CreateCredential(cred *Credential) error {
	credJSON, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	_, err = f.contract.SubmitTransaction("CreateCredential", string(credJSON))
	return err
}

func (f *FabricService) RevokeCredential(id string) error {
	_, err := f.contract.SubmitTransaction("RevokeCredential", id)

	if err != nil {
		return err
	}

	return nil
}

// LoadIssuers loads authorized issuers from JSON file
func LoadIssuers() error {
	data, err := os.ReadFile("../../backend/issuers.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &issuers)
}

func LoadUsers() ([]User, error) {
	data, err := os.ReadFile("../../backend/users.json")
	if err != nil {
		return nil, err
	}
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// ValidateAPIKey checks if the provided API key is valid for the issuer
func ValidateAPIKey(issuerID, apiKey string) bool {
	for _, issuer := range issuers {
		if issuer.ID == issuerID && issuer.APIKey == apiKey {
			return true
		}
	}
	return false
}

// GenerateCredentialID creates a unique ID for the credential
func GenerateCredentialID(diplomaHash string) string {
	h := sha256.New()
	h.Write([]byte(diplomaHash))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func main() {
	// Load authorized issuers
	if err := LoadIssuers(); err != nil {
		panic(fmt.Sprintf("Failed to load issuers: %v", err))
	}

	fs, err := ConnectGateway()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.POST("/auth/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		users, err := LoadUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
			return
		}

		// Hash the input password
		h := sha256.New()
		h.Write([]byte(req.Password))
		passwordHash := hex.EncodeToString(h.Sum(nil))

		for _, user := range users {
			if user.Username == req.Username && user.PasswordHash == passwordHash {
				var issuer Issuer
				for _, iss := range issuers {
					if iss.ID == user.IssuerID {
						issuer = iss
						break
					}
				}
				var userResponse = struct {
					Issuer    Issuer `json:"issuer"`
					Username  string `json:"username"`
					FirstName string `json:"firstName"`
					LastName  string `json:"lastName"`
				}{
					Issuer:    issuer,
					Username:  user.Username,
					FirstName: user.FirstName,
					LastName:  user.LastName,
				}
				c.JSON(http.StatusOK, userResponse)
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})

	// GET /credential/:id - Read credential by ID
	router.GET("/credential/:id", func(c *gin.Context) {
		id := c.Param("id")
		cred, err := fs.ReadCredential(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Credential not found", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cred)
	})

	// POST /credential - Create new credential (with API key validation)
	router.POST("/credential", func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
			return
		}

		var req CreateCredentialRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		// Validate API key
		if !ValidateAPIKey(req.IssuerID, apiKey) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key for issuer"})
			return
		}

		// Generate unique credential ID
		credentialID := GenerateCredentialID(req.DiplomaHash)

		// Create credential object
		credential := &Credential{
			ID:                credentialID,
			DiplomaHash:       req.DiplomaHash,
			GraduatePublicKey: req.GraduatePublicKey,
			IssuerID:          req.IssuerID,
			IssuerSignature:   req.IssuerSignature,
			DiplomaMetadata:   req.DiplomaMetadata,
			Status:            "Valid",
			CredentialType:    req.CredentialType,
		}

		// Submit to blockchain
		if err := fs.CreateCredential(credential); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create credential", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":      "Credential created successfully",
			"credentialId": credentialID,
			"credential":   credential,
		})
	})

	// POST /verify/hash - Verify diploma hash exists
	router.POST("/verify/hash", func(c *gin.Context) {
		var req VerifyHashRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		// Generate credential ID from hash
		credentialID := GenerateCredentialID(req.DiplomaHash)

		// Try to read the credential
		credential, err := fs.ReadCredential(credentialID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"verified": false,
				"message":  "Diploma hash not found in blockchain",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"verified":     true,
			"message":      "Diploma hash verified",
			"credentialId": credentialID,
			"status":       credential.Status,
			"issuerId":     credential.IssuerID,
		})
	})

	// POST /verify/signature - Verify graduate signature and return full diploma data
	router.POST("/verify/signature", func(c *gin.Context) {
		var req VerifySignatureRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		// Read credential from blockchain
		credential, err := fs.ReadCredential(req.CredentialID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"verified": false,
				"error":    "Credential not found",
			})
			return
		}

		// Verify PGP signature
		publicKeyArmored := credential.GraduatePublicKey
		signatureArmored := req.GraduateSignature
		messageToVerify := credential.ID

		// Decode public key
		keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(publicKeyArmored))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"verified": false,
				"error":    "Invalid public key format",
			})
			return
		}

		// Decode signature
		signatureReader := strings.NewReader(signatureArmored)
		messageReader := strings.NewReader(messageToVerify)

		// Verify the signature
		_, err = openpgp.CheckArmoredDetachedSignature(keyring, messageReader, signatureReader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"verified": false,
				"error":    "Signature verification failed",
				"details":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"verified":   true,
			"message":    "Graduate signature verified",
			"credential": credential,
		})
	})

	// GET /credentials - Get all credentials with optional university filter
	router.GET("/credentials", func(c *gin.Context) {
		universityFilter := c.Query("university")

		// Get all credentials from blockchain
		credentials, err := fs.GetAllCredentials()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve credentials",
				"details": err.Error(),
			})
			return
		}

		// Filter by university ID or name if provided
		if universityFilter != "" {
			var filtered []*Credential
			universityLower := strings.ToLower(universityFilter)
			for _, cred := range credentials {
				// Check both issuerID and university name in metadata
				if cred.IssuerID == universityFilter ||
					strings.ToLower(cred.DiplomaMetadata.UniversityName) == universityLower ||
					strings.Contains(strings.ToLower(cred.DiplomaMetadata.UniversityName), universityLower) {
					filtered = append(filtered, cred)
				}
			}
			credentials = filtered
		}

		c.JSON(http.StatusOK, gin.H{
			"credentials": credentials,
			"count":       len(credentials),
		})
	})

	// PATCH /credential/:id/revoke - Revoke credential by ID
	router.PATCH("/credential/:id/revoke", func(c *gin.Context) {
		id := c.Param("id")

		err := fs.RevokeCredential(id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Credential not found or could not be processed", "details": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	fmt.Println("Gateway running on http://0.0.0.0:8080")
	router.Run("0.0.0.0:8080")
}
