package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nadern96/go-ecommerce/controllers"
	"github.com/nadern96/go-ecommerce/database"
	"github.com/nadern96/go-ecommerce/middleware"
	"github.com/nadern96/go-ecommerce/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "products"), database.UserData(database.Client, "users"))
	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.POST("/add-to-cart", app.AddToCart())
	router.POST("/remove-item", app.RemoveCartItem())
	router.GET("/list-cart", app.GetItemFromCart())
	router.POST("/cart-checkout", app.BuyFromCart())
	router.POST("/instant-buy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
