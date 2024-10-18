package tsmcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
)

type Player struct {
	Url string `json:"url"`
}

type TSMController struct {
	Player1 Player
	Player2 Player
}

func NewTSMController(player1 Player, player2 Player) *TSMController {
	return &TSMController{
		Player1: player1,
		Player2: player2,
	}
}

type GenerateKeyRequestBody struct {
	SessionId string `json:"sessionId" binding:"required" example:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
}

func (t *TSMController) StartGenerateKeySession(publicKey string) string {
	sessionId := tsm.GenerateSessionID()

	requestBody := GenerateKeyRequestBody{SessionId: sessionId, PublicKey: publicKey}
	log.Printf("[StartGenerateKeySession] %v", requestBody)
	player1GenKeyUrl := fmt.Sprintf("%s/v1/generateKey", t.Player1.Url)
	go httpRequest(player1GenKeyUrl, "POST", requestBody)

	player2GenKeyUrl := fmt.Sprintf("%s/v1/generateKey", t.Player2.Url)
	go httpRequest(player2GenKeyUrl, "POST", requestBody)

	return sessionId
}

type CopyKeyRequestBody struct {
	SessionId     string `json:"sessionId" binding:"required" example:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
	PublicKey     string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	ExistingKeyId string `json:"existingKeyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

func (t *TSMController) StartCopyKeySession(publicKey string, existingKeyID string) string {
	/*
		/v1/copyKey
	*/
	sessionId := tsm.GenerateSessionID()

	player1CopyKeyUrl := fmt.Sprintf("%s/v1/copyKey", t.Player1.Url)
	go httpRequest(player1CopyKeyUrl, "POST", CopyKeyRequestBody{SessionId: sessionId, PublicKey: publicKey, ExistingKeyId: existingKeyID})

	player2CopyKeyUrl := fmt.Sprintf("%s/v1/copyKey", t.Player2.Url)
	go httpRequest(player2CopyKeyUrl, "POST", CopyKeyRequestBody{SessionId: sessionId, PublicKey: publicKey, ExistingKeyId: existingKeyID})

	return sessionId
}

type PresignRequestBody struct {
	SessionId string `json:"sessionId" binding:"required" example:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	KeyId     string `json:"keyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
	Count     uint64 `json:"count" binding:"required" example:"3"`
}

func (t *TSMController) StartPresignSession(publicKey string, keyId string, count uint64) string {
	/*
		/v1/preSign
	*/

	log.Printf("[StartPresignSession] publicKey: %s, keyId: %s, count: %d", publicKey, keyId, count)
	sessionId := tsm.GenerateSessionID()
	player1PresignUrl := fmt.Sprintf("%s/v1/preSign", t.Player1.Url)
	go httpRequest(player1PresignUrl, "POST", PresignRequestBody{SessionId: sessionId, PublicKey: publicKey, KeyId: keyId, Count: count})

	return sessionId
}

type PartialSignRequestBody struct {
	SignSignatureId string `json:"signSignatureId" binding:"required" example:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
	MessageHash     string `json:"messageHash" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	KeyId           string `json:"keyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

type PartialSignResponseBody struct {
	Signature string `json:"signature" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
}

func (t *TSMController) PartialSign(preSignatureId string, messageHash string, keyId string) (string, error) {
	/*
		/v1/partialSign
	*/

	player1PartialSignUrl := fmt.Sprintf("%s/v1/partialSign", t.Player1.Url)
	player1PartialSignResponseBody := httpRequest(player1PartialSignUrl, "POST", PartialSignRequestBody{SignSignatureId: preSignatureId, MessageHash: messageHash, KeyId: keyId})

	var responseBody PartialSignResponseBody
	err := json.Unmarshal(player1PartialSignResponseBody, &responseBody)
	if err != nil {
		log.Printf("[PartialSign] failed to json.Unmarshal. error: %s", err)
		return "", err
	}

	return responseBody.Signature, nil
}

func httpRequest(url string, method string, requestBody any) []byte {
	var requestBodyBytes []byte
	if method == "POST" {
		var err error
		requestBodyBytes, err = json.Marshal(requestBody)
		if err != nil {
			log.Printf("[httpRequest] failed to json.Marshal. error: %s", err)
		}
	}

	log.Printf("[httpRequest] url: %s, method: %s, requestBody: %s", url, method, string(requestBodyBytes))

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Printf("[httpRequest] failed to http.NewRequest. error: %s", err.Error())
	}
	req.Header.Set("User-Agent", "ABC")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[httpRequest] failed to http.NewRequest. error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[httpRequest] failed to http.NewRequest. error: %s", err.Error())
	}

	return body
}
