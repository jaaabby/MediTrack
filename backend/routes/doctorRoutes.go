package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupDoctorRoutes configura las rutas de doctores
func SetupDoctorRoutes(router *gin.RouterGroup, doctorService services.DoctorService) {
	doctorController := controllers.NewDoctorController(doctorService)

	doctors := router.Group("/doctors")
	{
		doctors.POST("/", doctorController.CreateDoctor)
		doctors.GET("/", doctorController.GetAllDoctors)
		doctors.GET("/:id", doctorController.GetDoctorByID)
		doctors.GET("/specialty", doctorController.GetDoctorsBySpecialty)
		doctors.PUT("/:id", doctorController.UpdateDoctor)
		doctors.DELETE("/:id", doctorController.DeleteDoctor)
	}
}
