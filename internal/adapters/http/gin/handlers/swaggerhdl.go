package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerHandler struct {
}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) SetupSwagger(ctx *gin.Context) {
	ctx.File("./swagger.yaml")
}

func (h *SwaggerHandler) SetupSwaggerUI(ctx *gin.Context) {
	handler := ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger.yaml"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	handler(ctx)
}
