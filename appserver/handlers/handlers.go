package handlers

import (
	"log"
	"net/http"

	"github.com/ahnlabio/tsm-appserver/tsmcontroller"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	TSMController *tsmcontroller.TSMController
}

func NewHandler(t *tsmcontroller.TSMController) *Handlers {
	return &Handlers{
		TSMController: t,
	}
}

type GenerateKeyRequestBody struct {
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
}

type GenerateKeyResponseBody struct {
	SessionId string `json:"sessionId" binding:"required" exaple:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
}

// GenerateKeyHandler godoc
// @Summary Generate a session key
// @Description Generate a session key
// @Tags session
// @Accept json
// @Produce json
// @Param body body GenerateKeyRequestBody true "Public key"
// @Success 200 {object} GenerateKeyResponseBody
// @Router /v1/generateKey [post]
func (h *Handlers) GenerateKeyHandler(c *gin.Context) {
	var requestBody GenerateKeyRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[GenerateKeyHandler] c.ShouldBind Error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionId := h.TSMController.StartGenerateKeySession(requestBody.PublicKey)
	log.Printf("[GenerateKeyHandler] session id: %s", sessionId)

	c.JSON(http.StatusOK, GenerateKeyResponseBody{SessionId: sessionId})
}

type CopyKeyRequestBody struct {
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	KeyId     string `json:"keyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

type CopyResponseBody struct {
	SessionId string `json:"sessionId" binding:"required" exaple:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
}

// CopyKeyHandler godoc
// @Summary Copy a session key
// @Description Copy a session key
// @Tags session
// @Accept json
// @Produce json
// @Param body body CopyKeyRequestBody true "Public key and key ID"
// @Success 200 {object} CopyResponseBody
// @Router /copyKey [post]
func (h *Handlers) CopyKeyHandler(c *gin.Context) {
	var requestBody CopyKeyRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[CopyKeyHandler] c.ShouldBind Error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionId := h.TSMController.StartCopyKeySession(requestBody.PublicKey, requestBody.KeyId)
	log.Printf("[CopyKeyHandler] session id: %s", sessionId)

	c.JSON(http.StatusOK, GenerateKeyResponseBody{SessionId: sessionId})
}

type PreSignRequestBody struct {
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	KeyId     string `json:"keyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
	Count     uint64 `json:"count" binding:"required" example:"3"`
}

type PreSignReponseBody struct {
	SessionId string `json:"sessionId" binding:"required" exaple:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
}

// PreSignHandler godoc
// @Summary Pre-sign a message
// @Description Pre-sign a message
// @Tags session
// @Accept json
// @Produce json
// @Param body body PreSignRequestBody true "Public key and key ID"
// @Success 200 {object} PreSignReponseBody
// @Router /preSign [post]
func (h *Handlers) PreSignHandler(c *gin.Context) {
	var requestBody PreSignRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[PreSignHandler] c.ShouldBind Error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionId := h.TSMController.StartPresignSession(requestBody.PublicKey, requestBody.KeyId, requestBody.Count)
	log.Printf("[PreSignHandler] session id: %s", sessionId)

	c.JSON(http.StatusOK, GenerateKeyResponseBody{SessionId: sessionId})
}

type PartialSignRequestBody struct {
	PreSignatureId string `json:"preSignatureId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
	MessageHash    string `json:"messageHash" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
	KeyId          string `json:"keyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

type PartialSignResponseBody struct {
	PartialSignature string `json:"partialSignResult" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

// PartialSignHandler godoc
// @Summary Finalize a signature
// @Description Finalize a signature
// @Tags session
// @Accept json
// @Produce json
// @Param body body FinalizSignRequestBody true "Pre-signature ID, message hash, and key ID"
// @Success 200 {object} FinalizeSignResponseBody
// @Router /finalizeSign [post]
func (h *Handlers) PartialSignHandler(c *gin.Context) {
	var requestBody PartialSignRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[PartialSignHandler] c.ShouldBind Error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	signature, err := h.TSMController.PartialSign(requestBody.PreSignatureId, requestBody.MessageHash, requestBody.KeyId)
	if err != nil {
		log.Printf("[PartialSignHandler] service.GenerateKey Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[PartialSignHandler] partialSignResult: %v", signature)
	c.JSON(http.StatusOK, PartialSignResponseBody{PartialSignature: signature})
}
