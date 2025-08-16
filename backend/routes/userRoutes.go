package routes

import (
	"meditrack/controllers"
	"meditrack/services"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configura las rutas de usuarios
func SetupUserRoutes(router *gin.RouterGroup, userService services.UserService) {
	userController := controllers.NewUserController(userService)

	users := router.Group("/users")
	{
		users.POST("/", userController.CreateUser)
		users.GET("/", userController.GetAllUsers)
		users.GET("/:id", userController.GetUserByID)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}
}
