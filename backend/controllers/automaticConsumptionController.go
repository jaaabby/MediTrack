package controllers

import (
	"meditrack/pkg/response"
	"meditrack/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AutomaticConsumptionController struct {
	automaticConsumptionService *services.AutomaticConsumptionService
}

func NewAutomaticConsumptionController(automaticConsumptionService *services.AutomaticConsumptionService) *AutomaticConsumptionController {
	return &AutomaticConsumptionController{
		automaticConsumptionService: automaticConsumptionService,
	}
}

// ProcessAutomaticConsumption ejecuta manualmente el proceso de consumo automático
func (c *AutomaticConsumptionController) ProcessAutomaticConsumption(ctx *gin.Context) {
	err := c.automaticConsumptionService.ProcessAutomaticConsumption()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   "Error al procesar consumo automático: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Proceso de consumo automático ejecutado exitosamente",
	})
}

