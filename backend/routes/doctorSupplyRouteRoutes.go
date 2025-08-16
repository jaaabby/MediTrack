package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupDoctorSupplyRouteRoutes configura las rutas de la relación doctor-ruta de suministro
func SetupDoctorSupplyRouteRoutes(router *gin.RouterGroup, doctorSupplyRouteService services.DoctorSupplyRouteService) {
	doctorSupplyRouteController := controllers.NewDoctorSupplyRouteController(doctorSupplyRouteService)

	doctorSupplyRoutes := router.Group("/doctor-supply-routes")
	{
		doctorSupplyRoutes.POST("/", doctorSupplyRouteController.CreateDoctorSupplyRoute)
		doctorSupplyRoutes.GET("/", doctorSupplyRouteController.GetAllDoctorSupplyRoutes)
		doctorSupplyRoutes.GET("/:id", doctorSupplyRouteController.GetDoctorSupplyRouteByID)
		doctorSupplyRoutes.GET("/doctor", doctorSupplyRouteController.GetDoctorSupplyRoutesByDoctorID)
		doctorSupplyRoutes.GET("/supply-route", doctorSupplyRouteController.GetDoctorSupplyRoutesBySupplyRouteID)
		doctorSupplyRoutes.PUT("/:id", doctorSupplyRouteController.UpdateDoctorSupplyRoute)
		doctorSupplyRoutes.DELETE("/:id", doctorSupplyRouteController.DeleteDoctorSupplyRoute)
	}
}
