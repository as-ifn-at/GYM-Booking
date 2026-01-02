package routes

import (
	"net/http"

	"github.com/as-ifn-at/REST/config"
	"github.com/as-ifn-at/REST/internal/db/gormdbwrapper"
	"github.com/as-ifn-at/REST/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type router struct {
	router    *gin.Engine
	appConfig config.Config
	logger    zerolog.Logger
	dbWrapper *gormdbwrapper.DBWrapper
}

func NewRouter(config *config.Config, logger zerolog.Logger, dbWrapper *gormdbwrapper.DBWrapper) *router {
	return &router{
		logger:    logger,
		router:    gin.Default(),
		appConfig: *config,
		dbWrapper: dbWrapper,
	}
}

func (r *router) SetRouters() http.Handler {
	attachMiddleWares(r.router)
	r.classesRoutes()
	r.attendClassesRoutes()

	return r.router.Handler()
}

func attachMiddleWares(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(middlewares.RateLimit())
}
