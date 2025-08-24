package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupMedicalCenterRoutes configura las rutas de medical center
func SetupMedicalCenterRoutes(router *gin.RouterGroup, medicalCenterService services.MedicalCenterService) {
	medicalCenterController := controllers.NewMedicalCenterController(medicalCenterService)
	medicalCenters := router.Group("/medical-centers")
	{
		medicalCenters.POST("/", medicalCenterController.CreateMedicalCenter)
		medicalCenters.GET("/", medicalCenterController.GetAllMedicalCenters)
		medicalCenters.GET("/:id", medicalCenterController.GetMedicalCenterByID)
		medicalCenters.PUT("/:id", medicalCenterController.UpdateMedicalCenter)
		medicalCenters.DELETE("/:id", medicalCenterController.DeleteMedicalCenter)
	}
}
