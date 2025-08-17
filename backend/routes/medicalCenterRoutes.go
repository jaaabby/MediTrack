package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupMedicalCenterRoutes configura las rutas de medical center
func SetupMedicalCenterRoutes(router *gin.RouterGroup, medicalCenterService services.MedicalCenterService) {
	medicalCenterController := controllers.NewMedicalCenterController(medicalCenterService)
	centers := router.Group("/medical-centers")
	{
		centers.POST("/", medicalCenterController.CreateMedicalCenter)
		centers.GET("/", medicalCenterController.GetAllMedicalCenters)
		centers.GET("/:id", medicalCenterController.GetMedicalCenterByID)
		centers.PUT("/:id", medicalCenterController.UpdateMedicalCenter)
		centers.DELETE("/:id", medicalCenterController.DeleteMedicalCenter)
	}
}
