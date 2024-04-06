// @title Car API
// @version 1.0

package http

import (
	"net/http"
	"tzsolution/configuration"
	"tzsolution/http/api"
	"tzsolution/postgresql"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Configuration *configuration.Configuration
	Engine        *gin.Engine
	Db            *postgresql.Postgresql
}

func NewHttpServer(configuration *configuration.Configuration, db *postgresql.Postgresql) *HttpServer {
	return &HttpServer{
		Engine:        gin.Default(),
		Configuration: configuration,
		Db:            db,
	}
}

func (server *HttpServer) router() {
	server.Engine.Use(gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	}))

	pepleApi := api.NewPeopleApi(server.Db)
	server.Engine.POST("/user", pepleApi.Add)

	carApi := api.NewCarApi(server.Db)
	server.Engine.GET("/info", carApi.GetInfoByRegNum)
	server.Engine.GET("/info/filter", carApi.GetFilteredAndPaginatedInfo)
	server.Engine.POST("/", carApi.Add)
	server.Engine.PUT("/", carApi.Update)
	server.Engine.DELETE("/", carApi.Delete)
}

func (server *HttpServer) Run(port string) {
	server.router()

	err := server.Engine.Run(port)
	if err != nil {
		panic(err)
	}
}
