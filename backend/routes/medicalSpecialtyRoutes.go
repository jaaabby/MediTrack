package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupMedicalSpecialtyRoutes configura las rutas de especialidades médicas
func SetupMedicalSpecialtyRoutes(router *gin.RouterGroup, specialtyService *services.MedicalSpecialtyService) {
	specialtyController := controllers.NewMedicalSpecialtyController(specialtyService)
	specialties := router.Group("/medical-specialties")
	{
		// CRUD básico
		specialties.POST("/", specialtyController.CreateMedicalSpecialty)
		specialties.GET("/", specialtyController.GetMedicalSpecialtiesPaginated)

		// Rutas específicas ANTES de las rutas con parámetros
		specialties.GET("/all", specialtyController.GetAllMedicalSpecialties)
		specialties.GET("/active", specialtyController.GetActiveMedicalSpecialties)
		specialties.GET("/search", specialtyController.SearchMedicalSpecialties)

		// Rutas con parámetros al final
		specialties.GET("/:id", specialtyController.GetMedicalSpecialtyByID)
		specialties.PUT("/:id", specialtyController.UpdateMedicalSpecialty)
		specialties.DELETE("/:id", specialtyController.DeleteMedicalSpecialty)
	}
}
