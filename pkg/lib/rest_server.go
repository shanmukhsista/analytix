package lib

import (
	"github.com/gin-gonic/gin"
)

type RestServerConfig struct {
	BindAddress string `envconfig:"SERVER_BIND_ADDRESS" required:"true" default:"localhost:9181"`
	IsDebug bool `envconfig:"SERVER_IS_DEBUG" default:"false"`
}

type GinRestServer struct {
	engine *gin.Engine
	config *RestServerConfig
}

func  NewRestServer(serverConfig *RestServerConfig)  *GinRestServer{
	g := gin.New()
	errorMiddleWare := ErrorMiddleware{
		Debug:        serverConfig.IsDebug,
		GenericError: "Internal Server Error",
		LogFunc:      nil,
	}
	g.Use(errorMiddleWare.Handler())
	return &GinRestServer{
		engine: g,
		config: serverConfig,
	}
}

func (g GinRestServer) GetEngine() *gin.Engine{
	return g.engine;
}

func (g GinRestServer) Start() error {
	err := g.engine.Run(g.config.BindAddress)
	if ( err != nil){
		return err
	}
	return nil
}