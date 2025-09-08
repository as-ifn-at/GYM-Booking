package handlers

import "github.com/gin-gonic/gin"

type Handler interface {
	Get(ctx *gin.Context)
	Save(ctx *gin.Context)
	// Delete(ctx *gin.Context)
}
