package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupOperatingRoomRoutes configura las rutas de salas de operación
func SetupOperatingRoomRoutes(router *gin.RouterGroup, operatingRoomService services.OperatingRoomService) {
	operatingRoomController := controllers.NewOperatingRoomController(operatingRoomService)

	operatingRooms := router.Group("/operating-rooms")
	{
		operatingRooms.POST("/", operatingRoomController.CreateOperatingRoom)
		operatingRooms.GET("/", operatingRoomController.GetAllOperatingRooms)
		operatingRooms.GET("/:id", operatingRoomController.GetOperatingRoomByID)
		operatingRooms.PUT("/:id", operatingRoomController.UpdateOperatingRoom)
		operatingRooms.DELETE("/:id", operatingRoomController.DeleteOperatingRoom)
	}
}
