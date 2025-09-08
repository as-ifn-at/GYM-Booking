package handlers

import (
	"fmt"
	"net/http"

	"github.com/as-ifn-at/REST/common"
	"github.com/as-ifn-at/REST/config"
	"github.com/as-ifn-at/REST/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var Classes = make(map[string]models.Class, 0)

var _ Handler = new(classHandler)

type classHandler struct {
	Handler
	config config.Config
	logger zerolog.Logger
}

func NewClassHandler(config config.Config, logger zerolog.Logger) Handler {
	return &classHandler{
		config: config,
		logger: logger,
	}
}

func (h *classHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	if _, ok := Classes[id]; !ok {
		h.logger.Error().Msg(fmt.Sprintf("class not found: %v", id))
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, Classes[id])
}

func (h *classHandler) Save(ctx *gin.Context) {
	var newClass models.Class
	if err := ctx.BindJSON(&newClass); err != nil {
		h.logger.Error().Msg(fmt.Sprintf("error while unmarshaling class details: %v", err.Error()))
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := common.CheckValidStartEndDate(newClass.StartDate, newClass.EndDate); err != nil {
		h.logger.Error().Msg(fmt.Sprintf("invalid start or end date: %v", err.Error()))
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	Classes[newClass.ClassName] = newClass
	h.logger.Info().Msg(fmt.Sprintf("successfully saved class information with name: %v", newClass.ClassName))

	ctx.IndentedJSON(http.StatusCreated, gin.H{"class created": newClass})
}
