package tsmutils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"

	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
)

type NodeConfig struct {
	Player0PublicKey     string // dynamic mobile node public key
	PlayerIndex          string // current server node index
	NodePubicKey         string // current server node public key
	AnotherNodePublicKey string // another server node public key
}

func GetClientFromConfig(config *tsm.Configuration) *tsm.Client {
	client, err := tsm.NewClient(config)
	if err != nil {
		panic(err)
	}
	return client
}

func CreateSignSessionConfig(sessionId string, nodeConfig NodeConfig) (*tsm.SessionConfig, error) {
	player0PublicKey := getPublicKeyBytesFromString(nodeConfig.Player0PublicKey)
	player1PublicKey := getPublicKeyBytesFromString(nodeConfig.NodePubicKey)
	dynamicPublicKeys := map[int][]byte{
		0: player0PublicKey,
		1: player1PublicKey,
	}

	var players []int = []int{0, 1}
	log.Printf("[tsmutils] CreateSignSessionConfig. dynamic public keys: %v", dynamicPublicKeys)
	sessionConfig := tsm.NewSessionConfig(sessionId, players, dynamicPublicKeys)
	return sessionConfig, nil
}

func CreateKeySessionConfig(sessionId string, nodeConfig NodeConfig) (*tsm.SessionConfig, error) {
	dynamicPublicKeys := getDynamicPublicKeys(nodeConfig)
	dumpPublicKeys(dynamicPublicKeys)

	var players []int = []int{0, 1, 2}
	log.Printf("[tsmutils] CreateKeySessionConfig. dynamic public keys: %v", dynamicPublicKeys)
	sessionConfig := tsm.NewSessionConfig(sessionId, players, dynamicPublicKeys)
	return sessionConfig, nil
}

func getDynamicPublicKeys(config NodeConfig) map[int][]byte {
	nodeIndex := getNodeIndex(config.PlayerIndex)
	anotherNodeIndex := getOtherNodeIndex(config.PlayerIndex)

	player0PublicKeyBytes := getPublicKeyBytesFromString(config.Player0PublicKey)
	nodePublicKeyBytes := getPublicKeyBytesFromString(config.NodePubicKey)
	anotherPublicKeyBytes := getPublicKeyBytesFromString(config.AnotherNodePublicKey)

	dynamicPublicKeys := map[int][]byte{
		0:                player0PublicKeyBytes,
		nodeIndex:        nodePublicKeyBytes,
		anotherNodeIndex: anotherPublicKeyBytes,
	}

	return dynamicPublicKeys
}

func getPublicKeyBytesFromString(publicKey string) []byte {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}

	return publicKeyBytes
}

func getNodeIndex(playerIndex string) int {
	switch playerIndex {
	case "1":
		return 1
	case "2":
		return 2
	}
	log.Printf("[ERROR] invalid playerIndex: %s", playerIndex)
	panic("invalid player index")
}

func getOtherNodeIndex(playerIndex string) int {
	switch playerIndex {
	case "1":
		return 2
	case "2":
		return 1
	}
	log.Printf("[ERROR] invalid playerIndex: %s", playerIndex)
	panic("invalid player index")
}

func dumpPublicKeys(publicKeys map[int][]byte) {
	for i, key := range publicKeys {
		log.Printf("key[%d]: %s", i, sha256Hex(key))
	}
}

func sha256Hex(data []byte) string {
	// Create a new SHA-256 hash
	hash := sha256.New()

	// Write data to the hash
	hash.Write(data)

	// Compute the SHA-256 checksum
	checksum := hash.Sum(nil)

	// Encode the checksum to a hexadecimal string
	hexString := hex.EncodeToString(checksum)

	// Return the hexadecimal string
	return hexString
}
