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
		// CRUD básico - Nuevos endpoints
		doctors.POST("/", doctorInfoController.CreateDoctor)
		doctors.GET("/", doctorInfoController.GetDoctorsPaginated)
		doctors.GET("/all", doctorInfoController.GetAllDoctors)

		// Rutas específicas ANTES de las rutas con parámetros
		doctors.GET("/specialty/:specialty_id", doctorInfoController.GetDoctorsBySpecialtyID)

		// Rutas con parámetros al final
		doctors.GET("/:rut", doctorInfoController.GetDoctorByRUT)
		doctors.PUT("/:rut", doctorInfoController.UpdateDoctor)
		doctors.DELETE("/:rut", doctorInfoController.DeleteDoctor)

		// Endpoints de compatibilidad (deprecated)
		// Estos endpoints mantienen la compatibilidad con el frontend antiguo
		// pero redirigen a los nuevos métodos
		doctors.POST("/info", doctorInfoController.CreateDoctorInfo) // deprecated
		doctors.GET("/info/:rut", doctorInfoController.GetDoctorInfoByRUT) // deprecated
		doctors.GET("/info/all", doctorInfoController.GetAllDoctorInfo) // deprecated
		doctors.PUT("/info/:rut", doctorInfoController.UpdateDoctorInfo) // deprecated
		doctors.DELETE("/info/:rut", doctorInfoController.DeleteDoctorInfo) // deprecated
	}
}
