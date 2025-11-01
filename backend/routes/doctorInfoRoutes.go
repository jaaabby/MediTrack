package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupDoctorInfoRoutes configura las rutas de información de doctores
func SetupDoctorInfoRoutes(router *gin.RouterGroup, doctorInfoService *services.DoctorInfoService) {
	doctorInfoController := controllers.NewDoctorInfoController(doctorInfoService)
	doctors := router.Group("/doctors")
	{
		// CRUD básico
		doctors.POST("/", doctorInfoController.CreateDoctorInfo)
		doctors.GET("/", doctorInfoController.GetDoctorsPaginated)

		// Rutas específicas ANTES de las rutas con parámetros
		doctors.GET("/all", doctorInfoController.GetAllDoctorInfo)
		doctors.GET("/specialty/:specialty_id", doctorInfoController.GetDoctorsBySpecialtyID)
		doctors.GET("/specialty-code/:code", doctorInfoController.GetDoctorsBySpecialtyCode)

		// Rutas con parámetros al final
		doctors.GET("/:rut", doctorInfoController.GetDoctorInfoByRUT)
		doctors.PUT("/:rut", doctorInfoController.UpdateDoctorInfo)
		doctors.DELETE("/:rut", doctorInfoController.DeleteDoctorInfo)
	}
}

