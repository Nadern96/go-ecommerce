package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nadern96/go-ecommerce/controllers"
)

func UserRoutes(routes *gin.Engine) {
	routes.POST("/users/signup", controllers.SignUp())
	routes.POST("/users/login", controllers.Login())
	routes.POST("/admin/add-product", controllers.AddProduct())
	routes.GET("/users/view-product", controllers.ViewProduct())
	routes.GET("/users/search", controllers.SearchProduct())
	routes.POST("/add-address", controllers.AddAddress())
	routes.PUT("/edit-home-address", controllers.EditHomeAddress())
	routes.PUT("/edit-work-address", controllers.EditWorkAddress())
	routes.DELETE("/delete-addresses", controllers.DeleteAddress())
}
