package container

import (
	"github.com/ahnlabio/tsm-controller/config"
	"github.com/ahnlabio/tsm-controller/handlers"
	"github.com/ahnlabio/tsm-controller/service"
)

var container *Container

type Container struct {
	AppConfig  *config.Config
	TsmService *service.TSMService
	Handlers   *handlers.Handlers
}

func GetInstnace() *Container {
	if container == nil {
		appConfig := config.GetConfig()
		tsmService := service.NewTSMService(appConfig)
		handers := handlers.NewHandler(tsmService)

		container = &Container{
			AppConfig:  appConfig,
			TsmService: tsmService,
			Handlers:   handers,
		}
	}
	return container
}

func (c *Container) GetHandlers() *handlers.Handlers {
	return c.Handlers
}
