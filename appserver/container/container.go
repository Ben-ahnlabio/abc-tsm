package container

import (
	"log"

	"github.com/ahnlabio/tsm-appserver/config"
	"github.com/ahnlabio/tsm-appserver/handlers"
	"github.com/ahnlabio/tsm-appserver/tsmcontroller"
)

var container *Container

type Container struct {
	AppConfig     *config.Config
	TSMController *tsmcontroller.TSMController
	Handlers      *handlers.Handlers
}

func GetInstnace() *Container {
	if container == nil {
		log.Print("Container is not initialized. Create new container.")
		appConfig := config.GetConfig()
		player1 := tsmcontroller.Player{Url: appConfig.Player1Url}
		player2 := tsmcontroller.Player{Url: appConfig.Player2Url}
		tsmController := tsmcontroller.NewTSMController(player1, player2)
		handlers := handlers.NewHandler(tsmController)

		container = &Container{
			AppConfig:     appConfig,
			TSMController: tsmController,
			Handlers:      handlers,
		}
	}
	return container
}

func (c *Container) GetHandlers() *handlers.Handlers {
	return c.Handlers
}
