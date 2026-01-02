# DatZM029-diploma-verification-system
**DatZM029**: Group 4 - Blockchain based credential verification system

**Team members**:
  bliiva
  DeiranLv
  Niklavs-M
  LinardsZ

## Overview
A blockchain-based diploma verification system built on Hyperledger Fabric that allows universities to issue and verify educational credentials securely.

## Prerequisites (https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html)
- **Docker Desktop v4.51.0 or older** (with WSL2 integration enabled for Windows)
- **WSL Ubuntu distro** with the following installed:
  - **Go**
  - **git**
  - **jq**

## Installation

### 0. Login into WSL
```bash
wsl.exe -d Ubuntu
```

### 1. Clone the Repository
```bash
git clone https://github.com/LinardsZ/DatZM029-diploma-verification-system.git
cd DatZM029-diploma-verification-system
```

### 2. Install Hyperledger Fabric Binaries and Docker Images
```bash
cd blockchain
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh -o install-fabric.sh
chmod +x install-fabric.sh
./install-fabric.sh binary docker
```

> If latest does not work - last known working version: fabric v2.5.14, fabric-ca 1.5.15

This will create:
- `blockchain/bin/` - Fabric binaries (peer, orderer, cryptogen, etc.)
- `blockchain/config/` - Configuration files

### 3. Copy Test Network
```bash
# Clone fabric-samples temporarily to get test-network
git clone --depth 1 --branch main https://github.com/hyperledger/fabric-samples.git temp-samples
cp -r temp-samples/test-network ./
rm -rf temp-samples
```

Your structure should now look like:
```
blockchain/
├── bin/
├── config/
├── test-network/
├── chaincode-go/
└── application-gateway/
```

## Running the Network

### 1. Start the Fabric Network and Create Channel
```bash
cd blockchain/test-network
./network.sh up createChannel -c mychannel -ca
```

This will:
- Generate crypto material (certificates)
- Start Docker containers (peers, orderer, CA)
- Create a channel named `mychannel`

**Verify containers are running:**
```bash
docker ps
```

You should see containers for:
- `peer0.org1.example.com`
- `peer0.org2.example.com`
- `orderer.example.com`

### 2. Deploy the Chaincode
```bash
./network.sh deployCC -ccn diploma -ccp ../chaincode-go -ccl go
```

This deploys the diploma verification smart contract to the network.

### 3. Initialize the Ledger
Set up environment variables:
```bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
```

Initialize with sample data:
```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C mychannel -n diploma \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"function":"InitLedger","Args":[]}'
```

## Starting the gateway

### 1. Build and start the gateway service
```bash
cd ~/DatZM029-diploma-verification-system/blockchain/test-network/application-gateway
go build -o gateway
./gateway
```

### 2. Get WSL IP address
```bash
ip addr show eth0
```

> Find the row starting with `inet`, e.g. `inet 192.1.68.1/20 ... scope global eth0` -> the IP address is 192.1.68.1

### 3. Test gateway service availability

Invoke authorized issuer creation and a mock diploma creation.
```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C mychannel -n diploma \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"function":"InitLedger","Args":[]}'
```

Invoke a new diploma creation.
```bash
peer chaincode invoke -o localhost:7050   --ordererTLSHostnameOverride orderer.example.com   --tls   --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"   -C mychannel -n diploma   --peerAddresses localhost:7051   --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"   --peerAddresses localhost:9051   --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"   -c '{"Args":["CreateCredential", "{\"id\":\"1\",\"diplomaHash\":\"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\",\"graduatePublicKey\":\"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234...\",\"issuerId\":\"lu\",\"issuerSignature\":\"3045022100abcd...\",\"diplomaMetadata\":{\"universityName\":\"MIT\",\"degreeName\":\"Bachelor of Science in Computer Science\",\"issueDate\":\"2024-06-15\",\"expiryDate\":\"\"},\"status\":\"Valid\",\"credentialType\":\"Diploma\"}"]}'
```

Open `<WSL_IP_ADDRESS>:8080/credential/1` in your browser on Windows. You should get a response back containing the credential data.

## Testing the Chaincode

### Query All Credentials
```bash
peer chaincode query -C mychannel -n diploma -c '{"Args":["GetAllCredentials"]}'
```

### Query Specific Credential
```bash
peer chaincode query -C mychannel -n diploma -c '{"Args":["ReadCredential","credential1"]}'
```

### Create a New Credential
```bash
peer chaincode invoke -o localhost:7050   --ordererTLSHostnameOverride orderer.example.com   --tls   --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"   -C mychannel -n diploma   --peerAddresses localhost:7051   --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"   --peerAddresses localhost:9051   --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"   -c '{"Args":["CreateCredential", "{\"id\":\"2\",\"diplomaHash\":\"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\",\"graduatePublicKey\":\"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234...\",\"issuerId\":\"lu\",\"issuerSignature\":\"3045022100abcd...\",\"diplomaMetadata\":{\"universityName\":\"MIT\",\"degreeName\":\"Bachelor of Science in Computer Science\",\"issueDate\":\"2024-06-15\",\"expiryDate\":\"\"},\"status\":\"Valid\",\"credentialType\":\"Diploma\"}"]}'
```

### Update Credential Status
```bash
peer chaincode invoke -o localhost:7050   --ordererTLSHostnameOverride orderer.example.com   --tls   --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"   -C mychannel -n diploma   --peerAddresses localhost:7051   --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"   --peerAddresses localhost:9051   --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"   -c '{"Args":["UpdateCredential", "{\"id\":\"1\",\"diplomaHash\":\"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\",\"graduatePublicKey\":\"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234...\",\"issuerId\":\"lu\",\"issuerSignature\":\"3045022100abcd...\",\"diplomaMetadata\":{\"universityName\":\"MIT\",\"degreeName\":\"Bachelor of Science in Computer Science\",\"issueDate\":\"2024-06-15\",\"expiryDate\":\"\"},\"status\":\"Revoked\",\"credentialType\":\"Diploma\"}"]}'
```

## Stopping the Network

```bash
cd blockchain/test-network
./network.sh down
```

This will:
- Stop all Docker containers
- Remove all Docker volumes
- Clean up the network

