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

var bookingsArr = []models.ClassBooking{}

type bookingHandler struct {
	Handler
	config config.Config
	logger zerolog.Logger
}

func NewBookingHandler(config config.Config, logger zerolog.Logger) Handler {
	return &bookingHandler{
		config: config,
		logger: logger,
	}
}

func (h *bookingHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	for _, class := range bookingsArr {
		if class.MemberName == id {
			ctx.IndentedJSON(http.StatusOK, class)
			return
		}
	}
	h.logger.Error().Msg(fmt.Sprintf("class member not found: %v", id))

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "member not found"})
}

func (h *bookingHandler) Save(ctx *gin.Context) {
	var newClassBooking models.ClassBooking
	if err := ctx.BindJSON(&newClassBooking); err != nil {
		h.logger.Error().Msg(fmt.Sprintf("error while unmarshaling booking details: %v", err.Error()))
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	classDetails, ok := Classes[newClassBooking.ClassName]
	if !ok {
		h.logger.Error().Msg(fmt.Sprintf("class: %v does not exist", newClassBooking.ClassName))
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "class does not exist"})
		return
	}

	if err := common.CheckClassAvailability(classDetails.StartDate,
		classDetails.EndDate, newClassBooking.Date); err != nil {
		h.logger.Error().Msg(fmt.Sprintf("invalid booking date: %v", err.Error()))
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	bookingsArr = append(bookingsArr, newClassBooking)
	h.logger.Info().Msg(fmt.Sprintf("successfully saved booking data for member: %v for class: %v",
		newClassBooking.MemberName, classDetails.ClassName))

	ctx.IndentedJSON(http.StatusCreated, gin.H{"booking done": newClassBooking})
}
