package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"example.com/tsmutils"
	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
)

var mobile0PublicKey = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="
var mobile1PublicKey = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEkAzm+8yn+d0ypywEwtgNnjisUkXBH17HpOd9YqRDybobqmCuaZA8cqAyLFS/qlu6j7lKCDWBwTElXJgvG9nywQ=="
var tsmDynamicMob0 = tsm.Configuration{URL: "http://localhost:8510"}.WithAPIKeyAuthentication("apikey0")
var tsmDynamicMob1 = tsm.Configuration{URL: "http://localhost:8511"}.WithAPIKeyAuthentication("apikey0")

type TSMNode struct {
	Config    *tsm.Configuration
	PublicKey string
	KeyId     string
}

type GetKeyResult struct {
	KeyId         string `json:"keyId"`
	UserPublicKey string `json:"userPublicKey"`
}

type CopyKeyResult struct {
	NewKeyId      string `json:"keyId"`
	UserPublicKey string `json:"userPublicKey"`
}

func main() {
	fmt.Printf("Hello, tsm client.\n")

	var nodes = []TSMNode{
		{
			Config:    tsmDynamicMob0,
			PublicKey: mobile0PublicKey,
			KeyId:     "",
		},
		{
			Config:    tsmDynamicMob1,
			PublicKey: mobile1PublicKey,
			KeyId:     "",
		},
	}

	// dynamic0 TSM 키 생성
	genKeyResult := client0GenKey(nodes[0].PublicKey)
	nodes[0].KeyId = genKeyResult.KeyId
	log.Printf("genKeyResult: %v\n", genKeyResult)

	// dynamic1 로 TSM 키 복사
	copyKeyResult := client1CopyKey(mobile1PublicKey, genKeyResult.KeyId)
	nodes[1].KeyId = copyKeyResult.NewKeyId
	log.Printf("copyKeyResult: %v\n", copyKeyResult)

	// public key 가 같은지 확인
	if genKeyResult.UserPublicKey != copyKeyResult.UserPublicKey {
		panic("User public key mismatch")
	}

	message := "Hello, world!"

	// dynamic node1 message 에 서명
	presignatureIds := preSign(nodes[1], 1)
	log.Printf("presignatureIds: %v\n", presignatureIds)
	messageBytes := []byte(message)
	msgHash := sha256.Sum256(messageBytes)
	sig1 := finalizeSign(nodes[1], presignatureIds[0], msgHash[:])

	client0 := tsmutils.GetClientFromConfig(nodes[1].Config)
	pubKey0, err := client0.Schnorr().PublicKey(context.TODO(), nodes[1].KeyId, nil)
	if err != nil {
		panic(err)
	}

	// verify signature node1
	log.Printf("verify signature node1\n")
	node1Err := tsm.SchnorrVerifySignature(pubKey0, msgHash[:], sig1)
	if node1Err != nil {
		panic(node1Err)
	}

	// dynamic node0 message 에 서명
	presignatureIds = preSign(nodes[0], 1)
	log.Printf("presignatureIds: %v\n", presignatureIds)
	sig2 := finalizeSign(nodes[0], presignatureIds[0], msgHash[:])

	client1 := tsmutils.GetClientFromConfig(nodes[0].Config)
	pubKey1, err := client1.Schnorr().PublicKey(context.TODO(), nodes[0].KeyId, nil)
	if err != nil {
		panic(err)
	}

	// verify signature node0
	log.Printf("verify signature node0\n")
	node0Err := tsm.SchnorrVerifySignature(pubKey1, msgHash[:], sig2)
	if node0Err != nil {
		panic(node0Err)
	}

	log.Printf("All signatures are verified\n")
}

func client0GenKey(nodePubKey string) *GetKeyResult {
	// appserver 에 요청하여 generate key session id 를 가져온다.
	// session id 가 발급되면 player1 과 player2 가 generate key 대기 상태가 된다.
	// player0 의 public key 를 player1, player2 에게 알려줘야 한다.
	sessionId := startGenerateKeySession(nodePubKey)
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(nodePubKey)
	if err != nil {
		panic(err)
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}
	players := []int{0, 1, 2} // The players (nodes) that should generate a sharing of the key
	sessionConfig := tsm.NewSessionConfig(sessionId, players, dynamicPublicKeys)
	ctx := context.Background()

	client := tsmutils.GetClientFromConfig(tsmDynamicMob0)
	threshold := 1

	// player1, player2 와 함께 generate key 를 실행한다.
	log.Printf("sessionConfig: %v\n", sessionConfig)
	log.Printf("Generating key. using client.Schnorr\n")
	keyId, err := client.Schnorr().GenerateKey(ctx, sessionConfig, threshold, "ED-25519", "")
	if err != nil {
		panic(err)
	}

	// 완료되면 player 0, 1, 2 가 key share 를 나눠 갖게 된다.
	log.Printf("keyId: %s\n", keyId)
	userPubkey := tsmutils.GetPubkeyStringFromClient(client, keyId)
	log.Printf("userPubkey: %s\n", userPubkey)
	return &GetKeyResult{KeyId: keyId, UserPublicKey: userPubkey}
}

func client1CopyKey(nodePubKey string, keyId string) *CopyKeyResult {
	// appserver 에 요청하여 copy key session id 를 가져온다.
	// session id 가 player1 과 player2 가 copy key 대기 상태가 된다.
	// player0 의 public key 를 player1, player2 에게 알려줘야 한다.
	sessionId := startCopyKeySession(nodePubKey, keyId)
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(nodePubKey)
	if err != nil {
		panic(err)
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	client := tsmutils.GetClientFromConfig(tsmDynamicMob1)
	newThreshold := 1
	newPlayers := []int{0, 1, 2} // The players (nodes) that should generate a sharing of the key
	keyCopySessionConfig := tsm.NewSessionConfig(sessionId, newPlayers, dynamicPublicKeys)

	ctx := context.Background()

	curveName := "ED-25519"
	// player1, player2 와 함께 copy key 를 실행한다.
	log.Printf("Coping key. using client.Schnorr\n")
	newKeyId, err := client.Schnorr().CopyKey(ctx, keyCopySessionConfig, "", curveName, newThreshold, "")
	if err != nil {
		panic(err)
	}

	// 완료되면 player 0, 1, 2 가 새로운 key share 를 나눠 갖게 된다.
	// 기존 키는 그대로 사용이 가능하다.
	// public key 는 이전에 만들었던 것과 동일하다.
	userPubkey := tsmutils.GetPubkeyStringFromClient(client, newKeyId)
	return &CopyKeyResult{
		NewKeyId:      newKeyId,
		UserPublicKey: userPubkey,
	}
}

func preSign(node TSMNode, presignatureCount uint64) []string {
	sessionId := startGeneratePreSignSignSession(node.PublicKey, node.KeyId)
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(node.PublicKey)
	if err != nil {
		panic(err)
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}
	var players []int = []int{0, 1}
	sessionConfig := tsm.NewSessionConfig(sessionId, players, dynamicPublicKeys)
	client := tsmutils.GetClientFromConfig(node.Config)
	preSignatureId, err := client.Schnorr().GeneratePresignatures(context.TODO(), sessionConfig, node.KeyId, presignatureCount)
	if err != nil {
		panic(err)
	}

	log.Printf("preSignatureId: %s\n", preSignatureId)
	return preSignatureId
}

func finalizeSign(node TSMNode, preSignatureId string, messageHash []byte) []byte {

	byteToStr := base64.StdEncoding.EncodeToString(messageHash)
	partialSigns := getPartialSignResult(preSignatureId, node.KeyId, byteToStr)
	client := tsmutils.GetClientFromConfig(node.Config)

	partialSignatures := make([][]byte, 0)
	partialSignBytes, err := base64.StdEncoding.DecodeString(partialSigns)
	if err != nil {
		panic(err)
	}
	partialSignatures = append(partialSignatures, partialSignBytes)

	partialSignResult, err := client.Schnorr().SignWithPresignature(context.TODO(), node.KeyId, preSignatureId, nil, messageHash)
	if err != nil {
		panic(err)
	}

	partialSignatures = append(partialSignatures, partialSignResult.PartialSignature)
	signature, err := tsm.SchnorrFinalizeSignature(messageHash, partialSignatures)
	if err != nil {
		panic(err)
	}

	return signature
}

type GenerateKeyRequestBody struct {
	PublicKey string `json:"publicKey"`
}

type GenerateKeyResponse struct {
	SessionId string `json:"sessionId"`
}

func startGenerateKeySession(publicKey string) string {

	url := "http://localhost:3000/v1/tsm/generateKey"
	addrReqBody := GenerateKeyRequestBody{
		PublicKey: publicKey,
	}
	value, _ := json.Marshal(addrReqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "ABC")

	client := &http.Client{Timeout: time.Duration(3000) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to get session id. status code: %d", resp.StatusCode))
	}

	var resObj GenerateKeyResponse
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		panic(err)
	}

	return resObj.SessionId
}

type CopyKeyRequestBody struct {
	PublicKey string `json:"publicKey"`
	KeyId     string `json:"keyId"`
}

type CopyKeyResponse struct {
	SessionId string `json:"sessionId"`
}

func startCopyKeySession(publicKey string, existingKeyId string) string {
	url := "http://localhost:3000/v1/tsm/copyKey"
	addrReqBody := CopyKeyRequestBody{
		PublicKey: publicKey,
		KeyId:     existingKeyId,
	}
	value, _ := json.Marshal(addrReqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "ABC")

	client := &http.Client{Timeout: time.Duration(3000) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to get session id. status code: %d", resp.StatusCode))
	}

	var resObj CopyKeyResponse
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		panic(err)
	}

	return resObj.SessionId
}

type PreSignRequestBody struct {
	PublicKey string `json:"publicKey"`
	KeyId     string `json:"keyId"`
	Count     uint64 `json:"count"`
}

type PreSignResponse struct {
	SessionId string `json:"sessionId"`
}

func startGeneratePreSignSignSession(publicKey string, keyId string) string {
	url := "http://localhost:3000/v1/tsm/preSign"
	addrReqBody := PreSignRequestBody{
		PublicKey: publicKey,
		KeyId:     keyId,
		Count:     1,
	}
	value, _ := json.Marshal(addrReqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "ABC")

	client := &http.Client{Timeout: time.Duration(3000) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to get session id. status code: %d", resp.StatusCode))
	}

	var resObj PreSignResponse
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		panic(err)
	}

	return resObj.SessionId
}

type GetPartialSizeResultRequestBody struct {
	PreSignatureId string `json:"preSignatureId"`
	KeyId          string `json:"keyId"`
	MessageHash    string `json:"messageHash"`
}

type GetPartialSignResultResponse struct {
	PartialSignResult string `json:"partialSignResult" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

func getPartialSignResult(preSignatureId string, keyId string, messageHash string) string {
	url := "http://localhost:3000/v1/tsm/finalizeSign"
	addrReqBody := GetPartialSizeResultRequestBody{
		PreSignatureId: preSignatureId,
		KeyId:          keyId,
		MessageHash:    messageHash,
	}
	value, _ := json.Marshal(addrReqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "ABC")

	client := &http.Client{Timeout: time.Duration(3000) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("failed to get session id. status code: %d", resp.StatusCode))
	}

	var resObj GetPartialSignResultResponse
	err = json.Unmarshal(body, &resObj)
	if err != nil {
		panic(err)
	}

	return resObj.PartialSignResult
}
