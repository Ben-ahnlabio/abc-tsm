package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"

	"github.com/ahnlabio/tsm-controller/config"
	"github.com/ahnlabio/tsm-controller/tsmutils"
)

type TSMService struct {
	config *config.Config
}

func NewTSMService(config *config.Config) *TSMService {
	return &TSMService{config: config}
}

func (s *TSMService) StartGenerateKeySession(sessionId string, publicKey string) error {
	/*
		GenreateKey session 을 시작합니다.
		Generate Key session 은 모든 노드가 참여합니다.
		session 시작 요청을 하고난 다음 node0 가 session 에 참여해 key 를 생성합니다.
	*/
	log.Printf("[Service] GenerateKey. sessionId: %s, publicKey: %s", sessionId, publicKey)
	sessionConfig, err := s.createKeygenSessionConfig(sessionId, publicKey)
	if err != nil {
		// encoding error. bad request 처리
		log.Printf("GenerateKey Service Error creating session config: %v", err)
		return err
	}

	client := s.getClient()
	threshold := 1 // The security threshold of the key
	curveName := "ED-25519"

	// 아래 go routine 이 실행되고난 다음 node0 또한 session 을 시작해야 합니다.
	go func() error {
		ctx := context.Background()
		log.Printf("GenerateKey session started. playerIndex: %s", s.config.PlayerIndex)
		log.Printf("client.Schnorr().GenerateKey. curveName: %s", curveName)
		keyId, err := client.Schnorr().GenerateKey(ctx, sessionConfig, threshold, curveName, "")
		if err != nil {
			log.Printf("Error generating key: %v", err)
			return err
		}
		log.Printf("Generated key with ID: %s, playerIndex: %s", keyId, s.config.PlayerIndex)
		return err
	}()

	return nil
}

func (s *TSMService) StartCopyKeySession(sessionId string, publicKey string, existingKeyId string) error {
	log.Printf("[Service] CopyKey. sessionId: %s, publicKey: %s, existingKeyID: %s", sessionId, publicKey, existingKeyId)
	sessionConfig, err := s.createKeygenSessionConfig(sessionId, publicKey)
	if err != nil {
		// encoding error. bad request 처리
		log.Printf("GenerateKey Service Error creating session config: %v", err)
		return err
	}

	client := s.getClient()
	newThreshold := 1 // The security threshold of the key
	curveName := "ED-25519"

	go func() error {
		var err error
		ctx := context.Background()
		log.Printf("client.Schnorr().CopyKey(). curveName: %s", curveName)
		newKeyId, err := client.Schnorr().CopyKey(ctx, sessionConfig, existingKeyId, curveName, newThreshold, "")
		if err != nil {
			log.Printf("Error generating key: %v", err)
			return err
		}
		log.Printf("Copied existingKeyID: %s, newKeyId: %s, playerIndex: %s", existingKeyId, newKeyId, s.config.PlayerIndex)
		return err
	}()

	return nil
}

func (s *TSMService) StartPresignSession(sessionId string, publicKey string, keyId string, presignatureCount uint64) error {
	log.Printf("[Service] PreSign. sessionId: %s, publicKey: %s, keyId: %s, presignatureCount: %d", sessionId, publicKey, keyId, presignatureCount)
	sessionConfig, err := s.createSignSessionConfig(sessionId, publicKey)
	if err != nil {
		log.Printf("GenerateKey Service Error creating session config: %v", err)
		return nil
	}

	client := s.getClient()
	go func() error {
		var err error
		ctx := context.Background()
		log.Printf("client.Schnorr().GeneratePresignatures()")
		_, err = client.Schnorr().GeneratePresignatures(ctx, sessionConfig, keyId, presignatureCount)
		if err != nil {
			log.Printf("Error generating presignature: %v", err)
		}

		log.Printf("Generated presignature. playerIndex: %s", s.config.PlayerIndex)
		return err
	}()

	return nil
}

func (s *TSMService) PartialSign(preSignatureId string, messageHash string, keyId string) (string, error) {
	log.Printf("[Service] PartialSign. preSignatureId: %s, messageHash: %s, keyId: %s", preSignatureId, messageHash, keyId)

	client := s.getClient()
	messageHashBytes, err := base64.StdEncoding.DecodeString(messageHash)
	if err != nil {
		return "", err
	}

	log.Printf("client.Schnorr().SignWithPresignature()")
	partialSignResult, err := client.Schnorr().SignWithPresignature(context.TODO(), keyId, preSignatureId, nil, messageHashBytes[:])
	if err != nil {
		return "", err
	}
	partialSignature := base64.StdEncoding.EncodeToString(partialSignResult.PartialSignature)
	return partialSignature, nil
}

func (s *TSMService) createKeygenSessionConfig(sessionId string, player0PublicKey string) (*tsm.SessionConfig, error) {
	/*
		session config 를 생성합니다.
		key generate, copy 는 모든 노드가 참여합니다.
		따라서 players 를 0, 1, 2 로 지정합니다.

		public key 는 node0 (mobile node) 의 public key 입니다.
		mobile node 는 dynamic node 이기 때문에 public key 를 입력 받아야 합니다.
		node1,node2 의 public key 는 설정파일에 저장되어 실행 시점부터 정해져 있습니다.
	*/

	nodeConfig := tsmutils.NodeConfig{
		PlayerIndex:          s.config.PlayerIndex,
		Player0PublicKey:     player0PublicKey,
		NodePubicKey:         s.config.NodePubicKey,
		AnotherNodePublicKey: s.config.AnotherNodePublicKey,
	}

	sessionConfig, err := tsmutils.CreateKeySessionConfig(sessionId, nodeConfig)
	if err != nil {
		return nil, errHandler(err)
	}
	return sessionConfig, nil
}

func (s *TSMService) createSignSessionConfig(sessionId string, player0PublicKey string) (*tsm.SessionConfig, error) {
	/*
		sign session config 를 생성합니다.
		sign 은 노드 0,1 만 참여합니다.
		따라서 players 를 0, 1 로 지정합니다.

		public key 는 node0 (mobile node) 의 public key 입니다.
		mobile node 는 dynamic node 이기 때문에 public key 를 입력 받아야 합니다.
		node1,node2 의 public key 는 설정파일에 저장되어 실행 시점부터 정해져 있습니다.
	*/
	if s.config.PlayerIndex != "1" {
		panic(fmt.Errorf("player is not allowed for sign"))
	}

	nodeConfig := tsmutils.NodeConfig{
		PlayerIndex:      s.config.PlayerIndex,
		Player0PublicKey: player0PublicKey,
		NodePubicKey:     s.config.NodePubicKey,
	}
	sessionConfig, err := tsmutils.CreateSignSessionConfig(sessionId, nodeConfig)
	if err != nil {
		return nil, errHandler(err)
	}
	return sessionConfig, nil
}

func (s *TSMService) getClient() *tsm.Client {
	tsmConfig := tsm.Configuration{URL: s.config.NodeUrl}.WithAPIKeyAuthentication(s.config.NodeApiKey)
	return tsmutils.GetClientFromConfig(tsmConfig)
}

func errHandler(err error) error {
	if errorInfo, ok := err.(*tsmutils.TsmUtilsErr); ok {
		if errorInfo.Text == tsmutils.DECODING_ERROR {
			// error 변환
			return InvalidInputError(err)
		}

		//log.Printf("[ERROR] err: %s, url: %s, status: %d\n", err.Error(), c.Request.URL, status)
	}

	return err
}
