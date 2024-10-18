package handlers

import (
	"log"
	"net/http"

	"github.com/ahnlabio/tsm-controller/service"
	"github.com/gin-gonic/gin"
)

type GenerateKeyRequestBody struct {
	SessionId string `json:"sessionId" binding:"required" example:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
}

type Handlers struct {
	service *service.TSMService
}

func NewHandler(service *service.TSMService) *Handlers {
	return &Handlers{service: service}
}

// GenerateKeyHandler godoc
// @Summary Generate a session key
// @Description Generate a session key
// @Tags session
// @Accept json
// @Produce json
// @Param body body GenerateKeyRequestBody true "Public key"
// @Success 200
// @Router /v1/generateKey [post]
func (h *Handlers) GenerateKeyHandler(c *gin.Context) {
	var requestBody GenerateKeyRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[GenerateKeyHandler] c.ShouldBind Error: %v\n", err)
		errResp(c, err)
		return
	}

	err = h.service.StartGenerateKeySession(requestBody.SessionId, requestBody.PublicKey)
	if err != nil {
		log.Printf("[GenerateKeyHandler] service.GenerateKey Error: %v\n", err)
		errResp(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

type CopyKeyRequestBody struct {
	SessionId     string `json:"sessionId" binding:"required" example:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
	PublicKey     string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	ExistingKeyId string `json:"existingKeyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

// CopyKeyHandler godoc
// @Summary Copy a session key
// @Description Copy a session key
// @Tags session
// @Accept json
// @Produce json
// @Param body body CopyKeyRequestBody true "Public key"
// @Success 200
// @Router /v1/copyKey [post]
func (h *Handlers) CopyKeyHandler(c *gin.Context) {
	var requestBody CopyKeyRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[CopyKeyHandler] c.ShouldBind Error: %v\n", err)
		errResp(c, err)
		return
	}

	err = h.service.StartCopyKeySession(requestBody.SessionId, requestBody.PublicKey, requestBody.ExistingKeyId)
	if err != nil {
		log.Printf("[CopyKeyHandler] service.CopyKey Error: %v\n", err)
		errResp(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

type PresignRequestBody struct {
	SessionId string `json:"sessionId" binding:"required" example:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
	PublicKey string `json:"publicKey" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	KeyId     string `json:"keyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
	Count     uint64 `json:"count" binding:"required" example:"3"`
}

// PreSignHandler godoc
// @Summary Start a presign session
// @Description Start a presign session
// @Tags session
// @Accept json
// @Produce json
// @Param body body PresignRequestBody true "Public key"
// @Success 200
// @Router /v1/presign [post]
func (h *Handlers) PreSignHandler(c *gin.Context) {
	var requestBody PresignRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[PreSignHandler] c.ShouldBind Error: %v\n", err)
		errResp(c, err)
		return
	}

	err = h.service.StartPresignSession(requestBody.SessionId, requestBody.PublicKey, requestBody.KeyId, requestBody.Count)
	if err != nil {
		log.Printf("[PreSignHandler] service.PreSign Error: %v\n", err)
		errResp(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}

type SignRequestBody struct {
	SignSignatureId string `json:"signSignatureId" binding:"required" example:"923J-NNcZlScEGi1phSmDWO-eZsQLtBGHVWIIIWZ7Zw"`
	MessageHash     string `json:"messageHash" binding:"required" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
	KeyId           string `json:"keyId" binding:"required" example:"zUhWR7jvWJoplMyFf35NHSdZXbtx"`
}

type SignResponseBody struct {
	Signature string `json:"signature" example:"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE2Bk6ZSVUhIStsXZsqyYidPy8vEQvLDVQ/YRgfgowgWFualE748OFoGwuGgE8C7L2zV4gX+1Ow1x/OTjqSSlh5A=="`
}

// PartialSignHandler godoc
// @Summary Finalize a sign session
// @Description Finalize a sign session
// @Tags session
// @Accept json
// @Produce json
// @Param body body SignRequestBody true "Public key"
// @Success 200
// @Router /v1/sign [post]
func (h *Handlers) PartialSignHandler(c *gin.Context) {
	var requestBody SignRequestBody
	err := c.ShouldBind(&requestBody)
	if err != nil {
		log.Printf("[SignHandler] c.ShouldBind Error: %v\n", err)
		errResp(c, err)
		return
	}

	signature, err := h.service.PartialSign(requestBody.SignSignatureId, requestBody.MessageHash, requestBody.KeyId)
	if err != nil {
		log.Printf("[SignHandler] service.Sign Error: %v\n", err)
		errResp(c, err)
		return
	}
	c.JSON(http.StatusOK, SignResponseBody{Signature: signature})
}

func errResp(c *gin.Context, err error) {
	if errorInfo, ok := err.(*service.SvcErr); ok {
		res := CommonErrorObject{
			Message: errorInfo.Msg,
			Text:    errorInfo.Text,
		}

		status := http.StatusInternalServerError
		if errorInfo.Text == service.INVALID_INPUT {
			status = http.StatusBadRequest
		}

		log.Printf("[ERROR] err: %s, url: %s, status: %d\n", err.Error(), c.Request.URL, status)
		c.JSON(status, gin.H{"error": &res})
		return
	}
	res := CommonErrorObject{
		Message: err.Error(),
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": &res})
}
