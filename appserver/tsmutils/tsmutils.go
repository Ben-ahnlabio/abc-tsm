package tsmutils

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/ahnlabio/tsm-appserver/config"
	"gitlab.com/Blockdaemon/go-tsm-sdkv2/v64/tsm"
)

func GetPubkeyStringFromClient(client *tsm.Client, keyId string) string {
	ctx := context.Background()
	var derivationPath []uint32 = nil
	publicKey, err := client.Schnorr().PublicKey(ctx, keyId, derivationPath)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(publicKey)
}

func GetClientsFromConfigs(configs []*tsm.Configuration) []*tsm.Client {
	clients := make([]*tsm.Client, len(configs))
	for i, config := range configs {
		var err error
		if clients[i], err = tsm.NewClient(config); err != nil {
			panic(err)
		}
	}
	return clients
}

func GenerateSessionConfig(players []int, pubkeyStr string) *tsm.SessionConfig {
	player0PublicTenantKey, err := base64.StdEncoding.DecodeString(pubkeyStr)
	if err != nil {
		panic(err)
	}

	dynamicPublicKeys := map[int][]byte{
		0: player0PublicTenantKey,
	}

	sessionID := tsm.GenerateSessionID()
	sessionConfig := tsm.NewSessionConfig(sessionID, players, dynamicPublicKeys)
	return sessionConfig
}

func GetClientFromConfig(config *tsm.Configuration) *tsm.Client {
	client, err := tsm.NewClient(config)
	if err != nil {
		panic(err)
	}
	return client
}

func KeyListing(configs []*tsm.Configuration) {
	clients := make([]*tsm.Client, len(configs))
	for i, config := range configs {
		var err error
		if clients[i], err = tsm.NewClient(config); err != nil {
			panic(err)
		}
	}

	ctx := context.Background()
	for idx, client := range clients {
		keyIDs, err := client.KeyManagement().ListKeys(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("node: %d, keyIDs %v\n", idx, keyIDs)
	}
}

func GetClients() []*tsm.Client {
	appConfig := config.GetConfig()
	configs := []*tsm.Configuration{
		tsm.Configuration{URL: appConfig.Node1Url}.WithAPIKeyAuthentication(appConfig.Node1ApiKey),
		tsm.Configuration{URL: appConfig.Node2Url}.WithAPIKeyAuthentication(appConfig.Node2ApiKey),
	}

	clients := make([]*tsm.Client, len(configs))
	for i, config := range configs {
		var err error
		if clients[i], err = tsm.NewClient(config); err != nil {
			panic(err)
		}
	}
	return clients
}
